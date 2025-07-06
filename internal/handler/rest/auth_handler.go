package rest

import (
	"errors"
	"net/http"
	"time"

	gen "github.com/DucTran999/auth-service/gen/http"
	"github.com/DucTran999/auth-service/internal/model"
	"github.com/DucTran999/auth-service/internal/usecase"
	"github.com/DucTran999/auth-service/internal/usecase/dto"
	"github.com/DucTran999/auth-service/internal/usecase/port"
	"github.com/DucTran999/shared-pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	sessionKey      = "session_id"
	refreshTokenKey = "refresh_token"
)

type AuthHandler struct {
	BaseHandler
	logger logger.ILogger
	authUC port.AuthUseCase
}

func NewAuthHandler(logger logger.ILogger, authUC port.AuthUseCase) *AuthHandler {
	return &AuthHandler{
		logger: logger,
		authUC: authUC,
	}
}

func (hdl *AuthHandler) LoginAccount(ctx *gin.Context) {
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

	// Convert request to model
	loginInput := dto.LoginInput{
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

func (hdl *AuthHandler) LoginAccountJWT(ctx *gin.Context) {
	// Parse request body
	payload, err := ParseAndValidateJSON[gen.LoginAccountJSONRequestBody](ctx)
	if err != nil {
		hdl.BadRequestResponse(ctx, APIVersion2, err.Error())
		return
	}

	input := dto.LoginJWTInput{
		Email:     string(payload.Email),
		Password:  payload.Password,
		IP:        ctx.ClientIP(),
		UserAgent: ctx.Request.UserAgent(),
	}

	tokens, err := hdl.authUC.LoginJWT(ctx, input)
	if err != nil {
		if errors.Is(err, usecase.ErrInvalidCredentials) {
			hdl.UnauthorizeErrorResponse(ctx, APIVersion2, err.Error())
			return
		}
		hdl.logger.Error("failed to login with jwt", zap.String("error", err.Error()))
		hdl.ServerInternalErrResponse(ctx, APIVersion2)
		return
	}

	hdl.responseLoginJWTSuccess(ctx, tokens)
}

func (hdl *AuthHandler) LogoutAccount(ctx *gin.Context) {
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

func (hdl *AuthHandler) responseLoginSuccess(ctx *gin.Context, session *model.Session) {
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

func (hdl *AuthHandler) responseLoginJWTSuccess(ctx *gin.Context, tokens *dto.TokenPairs) {
	// Determine environment is secure or not
	secure := ctx.Request.Header.Get("X-Forwarded-Proto") == "https" || ctx.Request.TLS != nil

	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:     refreshTokenKey,
		Value:    tokens.RefreshToken,
		Path:     "/refresh-token",
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(usecase.RefreshTokenLifetime),
	})

	resp := gen.LoginJWTResponse{
		Success: true,
		Version: APIVersion2,
		Data: gen.AccessToken{
			AccessToken: tokens.AccessToken,
		},
	}

	ctx.JSON(http.StatusOK, resp)
}
