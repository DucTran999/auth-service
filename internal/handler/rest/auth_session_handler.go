package rest

import (
	"errors"
	"net/http"

	gen "github.com/DucTran999/auth-service/gen/http"
	"github.com/DucTran999/auth-service/internal/errs"
	"github.com/DucTran999/auth-service/internal/model"
	"github.com/DucTran999/auth-service/internal/usecase/dto"
	"github.com/DucTran999/auth-service/internal/usecase/port"
	"github.com/DucTran999/shared-pkg/logger"
	"github.com/gin-gonic/gin"
)

const (
	SessionKey = "session_id"
)

type SessionAuthHandler interface {
	LoginAccount(ctx *gin.Context)
	LogoutAccount(ctx *gin.Context)
}

type sessionAuthHandler struct {
	BaseHandler
	logger logger.ILogger
	authUC port.AuthSessionUsecase
}

func NewSessionAuthHandler(logger logger.ILogger, authUC port.AuthSessionUsecase) SessionAuthHandler {
	return &sessionAuthHandler{
		logger: logger,
		authUC: authUC,
	}
}

func (hdl *sessionAuthHandler) LoginAccount(ctx *gin.Context) {
	// Parse request body
	payload, err := ParseAndValidateJSON[gen.LoginAccountJSONRequestBody](ctx)
	if err != nil {
		hdl.BadRequestResponse(ctx, ApiVersion1, err.Error())
		return
	}

	// Set to empty when cookie not found
	currentSessionID, err := ctx.Cookie(SessionKey)
	if err != nil {
		currentSessionID = ""
	}

	// Convert request to model
	loginInput := dto.LoginInput{
		CurrentSessionID: currentSessionID,
		Email:            payload.Email,
		Password:         payload.Password,
		IP:               ctx.ClientIP(),
		UserAgent:        ctx.Request.UserAgent(),
	}

	// Authenticate user and create session
	session, err := hdl.authUC.Login(ctx.Request.Context(), loginInput)
	if err != nil {
		if errors.Is(err, errs.ErrInvalidCredentials) || errors.Is(err, errs.ErrAccountNotFound) {
			hdl.UnauthorizeErrorResponse(ctx, ApiVersion1, err.Error())
			return
		}

		hdl.logger.Error(err.Error())
		hdl.ServerInternalErrResponse(ctx, ApiVersion1)
		return
	}

	hdl.responseLoginSuccess(ctx, session)
}

func (hdl *sessionAuthHandler) LogoutAccount(ctx *gin.Context) {
	// Try to get session ID from cookie
	sessionID, err := ctx.Cookie(SessionKey)
	if err == nil {
		// Best-effort logout
		if err := hdl.authUC.Logout(ctx.Request.Context(), sessionID); err != nil {
			hdl.logger.Warn(err.Error())
		}
	}

	// Always clear the cookie
	ctx.SetCookie(SessionKey, "", -1, "/", "", true, true)

	// Always respond with 204 No Content
	hdl.NoContentResponse(ctx)
}

func (hdl *sessionAuthHandler) responseLoginSuccess(ctx *gin.Context, session *model.Session) {
	// Determine environment is secure or not
	secure := ctx.Request.Header.Get("X-Forwarded-Proto") == "https" || ctx.Request.TLS != nil

	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:     SessionKey,
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
