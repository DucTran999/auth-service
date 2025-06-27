package http

import (
	"errors"
	"net/http"

	"github.com/DucTran999/auth-service/internal/domain"
	"github.com/DucTran999/auth-service/internal/gen"
	"github.com/DucTran999/auth-service/internal/usecase"
	"github.com/DucTran999/shared-pkg/logger"
	"github.com/gin-gonic/gin"
)

const (
	sessionKey = "session_id"
)

type AuthHandler interface {
	LoginAccount(ctx *gin.Context)
	LogoutAccount(ctx *gin.Context)
}

type authHandlerImpl struct {
	BaseHandler
	logger logger.ILogger
	authUC domain.AuthUseCase
}

func NewAuthHandler(logger logger.ILogger, authUC domain.AuthUseCase) *authHandlerImpl {
	return &authHandlerImpl{
		logger: logger,
		authUC: authUC,
	}
}

func (hdl *authHandlerImpl) LoginAccount(ctx *gin.Context) {
	// Parse request body
	payload, err := ParseAndValidateJSON[gen.LoginAccountJSONRequestBody](ctx)
	if err != nil {
		hdl.BadRequestResponse(ctx, ApiVersion1, err.Error())
		return
	}

	// Set to empty when cookie not found
	currentSessionID, err := ctx.Cookie(sessionKey)
	if err != nil {
		currentSessionID = ""
	}

	// Convert request to domain model
	loginInput := domain.LoginInput{
		CurrentSessionID: currentSessionID,
		Email:            string(payload.Email),
		Password:         payload.Password,
		IP:               ctx.ClientIP(),
		UserAgent:        ctx.Request.UserAgent(),
	}

	// Authenticate user and create session
	session, err := hdl.authUC.Login(ctx.Request.Context(), loginInput)
	if err != nil {
		if errors.Is(err, usecase.ErrInvalidCredentials) {
			hdl.UnauthorizeErrorResponse(ctx, ApiVersion1, err.Error())
			return
		}

		hdl.logger.Error(err.Error())
		hdl.ServerInternalErrResponse(ctx, ApiVersion1)
		return
	}

	hdl.responseLoginSuccess(ctx, session)
}

func (hdl *authHandlerImpl) LogoutAccount(ctx *gin.Context) {
	// Try to get session ID from cookie
	sessionID, err := ctx.Cookie(sessionKey)
	if err == nil {
		// Best-effort logout
		if err := hdl.authUC.Logout(ctx.Request.Context(), sessionID); err != nil {
			hdl.logger.Warn(err.Error())
		}
	}

	// Always clear the cookie
	ctx.SetCookie(sessionKey, "", -1, "/", "", true, true)

	// Always respond with 204 No Content
	hdl.NoContentResponse(ctx)
}

func (hdl *authHandlerImpl) responseLoginSuccess(ctx *gin.Context, session *domain.Session) {
	// Determine environment is secure or not
	secure := ctx.Request.Header.Get("X-Forwarded-Proto") == "https" || ctx.Request.TLS != nil

	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:     sessionKey,
		Value:    session.ID.String(),
		Path:     "/",
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteStrictMode,
	})

	resp := gen.LoginResponse{
		Success: true,
		Version: ApiVersion1,
		Data: gen.Account{
			Id:    session.AccountID,
			Email: session.Account.Email,
			Role:  session.Account.Role,
		},
	}
	ctx.JSON(http.StatusOK, resp)
}
