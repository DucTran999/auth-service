package rest

import (
	"errors"
	"net/http"
	"time"

	gen "github.com/DucTran999/auth-service/gen/http"
	"github.com/DucTran999/auth-service/internal/errs"
	"github.com/DucTran999/auth-service/internal/usecase/auth"
	"github.com/DucTran999/auth-service/internal/usecase/dto"
	"github.com/DucTran999/auth-service/internal/usecase/port"
	"github.com/DucTran999/shared-pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	refreshTokenKey = "refresh_token"
)

type JWTAuthHandler interface {
	LoginWithJWT(ctx *gin.Context)
	RefreshToken(ctx *gin.Context)
	LogoutJWT(ctx *gin.Context)
}

type jwtAuthHandler struct {
	BaseHandler
	logger logger.ILogger
	authUC port.AuthJWTUsecase
}

// NewJWTAuthHandler creates a new JWT authentication handler.
func NewJWTAuthHandler(logger logger.ILogger, authUC port.AuthJWTUsecase) JWTAuthHandler {
	return &jwtAuthHandler{
		logger: logger,
		authUC: authUC,
	}
}

func (hdl *jwtAuthHandler) LoginWithJWT(ctx *gin.Context) {
	// Parse request body
	payload, err := ParseAndValidateJSON[gen.LoginAccountJSONRequestBody](ctx)
	if err != nil {
		hdl.BadRequestResponse(ctx, APIVersion2, err.Error())
		return
	}

	// Prepare input
	input := dto.LoginJWTInput{
		Email:     string(payload.Email),
		Password:  payload.Password,
		IP:        ctx.ClientIP(),
		UserAgent: ctx.Request.UserAgent(),
	}

	// Authenticate
	tokens, err := hdl.authUC.Login(ctx, input)
	if err != nil {
		if errors.Is(err, errs.ErrInvalidCredentials) {
			hdl.UnauthorizeErrorResponse(ctx, APIVersion2, err.Error())
			return
		}
		hdl.logger.Error("failed to login with jwt", zap.String("error", err.Error()))
		hdl.ServerInternalErrResponse(ctx, APIVersion2)
		return
	}

	hdl.responseLoginJWTSuccess(ctx, tokens)
}

func (hdl *jwtAuthHandler) RefreshToken(ctx *gin.Context) {
	refreshToken, err := ctx.Cookie(refreshTokenKey)
	if err != nil {
		hdl.UnauthorizeErrorResponse(ctx, APIVersion2, http.StatusText(http.StatusUnauthorized))
		return
	}

	tokens, err := hdl.authUC.RefreshToken(ctx, refreshToken)
	if err != nil {
		hdl.UnauthorizeErrorResponse(ctx, APIVersion2, err.Error())
		return
	}

	hdl.responseLoginJWTSuccess(ctx, tokens)
}

func (hdl *jwtAuthHandler) LogoutJWT(ctx *gin.Context) {
	refreshToken, err := ctx.Cookie(refreshTokenKey)
	if err != nil {
		hdl.UnauthorizeErrorResponse(ctx, APIVersion2, http.StatusText(http.StatusUnauthorized))
		return
	}

	if err := hdl.authUC.RevokeRefreshToken(ctx, refreshToken); err != nil {
		hdl.UnauthorizeErrorResponse(ctx, APIVersion2, http.StatusText(http.StatusUnauthorized))
		return
	}

	// Always clear the cookie
	ctx.SetCookie(refreshTokenKey, "", -1, "/", "", true, true)

	// Always respond with 204 No Content
	hdl.NoContentResponse(ctx)
}

func (hdl *jwtAuthHandler) responseLoginJWTSuccess(ctx *gin.Context, tokens *dto.TokenPairs) {
	// Determine environment is secure or not
	secure := ctx.Request.Header.Get("X-Forwarded-Proto") == "https" || ctx.Request.TLS != nil

	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:     refreshTokenKey,
		Value:    tokens.RefreshToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(auth.RefreshTokenLifetime),
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
