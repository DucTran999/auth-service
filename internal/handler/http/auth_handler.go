package http

import (
	"errors"
	"net/http"

	"github.com/DucTran999/auth-service/internal/gen"
	"github.com/DucTran999/auth-service/internal/model"
	"github.com/DucTran999/auth-service/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type AuthHandler interface {
	LoginAccount(ctx *gin.Context)
}

type authHandlerImpl struct {
	BaseHandler
	authUC usecase.AuthUseCase
}

func NewAuthHandler(authUC usecase.AuthUseCase) *authHandlerImpl {
	return &authHandlerImpl{
		authUC: authUC,
	}
}

func (hdl *authHandlerImpl) LoginAccount(ctx *gin.Context) {
	// Parse request body
	payload, err := hdl.parseAndValidateLoginCredentials(ctx)
	if err != nil {
		return
	}

	// Set to empty when cookie not found
	currentSessionID, err := ctx.Cookie("session_id")
	if err != nil {
		currentSessionID = ""
	}

	// Convert request to domain model
	loginInput := usecase.LoginInput{
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
		hdl.ServerInternalErrResponse(ctx, ApiVersion1)
		return
	}

	hdl.responseLoginSuccess(ctx, session)
}

func (hdl *authHandlerImpl) parseAndValidateLoginCredentials(ctx *gin.Context,
) (*gen.LoginAccountJSONRequestBody, error) {

	var payload gen.LoginAccountJSONRequestBody
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) && len(ve) > 0 {
			hdl.ValidateErrorResponse(ctx, ApiVersion1, validationErrorMessage(ve[0]))
			return nil, err
		}

		hdl.BadRequestResponse(ctx, ApiVersion1, err.Error())
		return nil, err
	}

	return &payload, nil
}

func (hdl *authHandlerImpl) responseLoginSuccess(ctx *gin.Context, session *model.Session) {
	// Determine environment is secure or not
	secure := ctx.Request.Header.Get("X-Forwarded-Proto") == "https" || ctx.Request.TLS != nil

	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:     "session_id",
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
