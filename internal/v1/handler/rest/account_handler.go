package rest

import (
	"errors"
	"net/http"

	gen "github.com/DucTran999/auth-service/gen/http"
	"github.com/DucTran999/auth-service/internal/model"
	"github.com/DucTran999/auth-service/internal/v1/usecase"
	"github.com/DucTran999/auth-service/internal/v1/usecase/dto"
	"github.com/DucTran999/auth-service/internal/v1/usecase/port"
	"github.com/DucTran999/shared-pkg/logger"
	"github.com/gin-gonic/gin"
)

type AccountHandlerImpl struct {
	logger logger.ILogger

	BaseHandler
	accountUC port.AccountUseCase
	sessionUC port.SessionUsecase
}

func NewAccountHandler(
	logger logger.ILogger,
	accountUC port.AccountUseCase,
	sessionUC port.SessionUsecase,
) *AccountHandlerImpl {
	return &AccountHandlerImpl{
		logger:    logger,
		accountUC: accountUC,
		sessionUC: sessionUC,
	}
}

// CreateAccount handles the HTTP request to register a new account.
func (hdl *AccountHandlerImpl) CreateAccount(ctx *gin.Context) {
	payload, err := ParseAndValidateJSON[gen.CreateAccountJSONRequestBody](ctx)
	if err != nil {
		hdl.BadRequestResponse(ctx, ApiVersion1, err.Error())
		return
	}

	input := dto.RegisterInput{
		Email:    string(payload.Email),
		Password: payload.Password,
	}
	account, err := hdl.accountUC.Register(ctx.Request.Context(), input)
	if err != nil {
		hdl.handleRegisterError(ctx, err)
		return
	}

	hdl.sendRegisterSuccess(ctx, account)
}

func (hdl *AccountHandlerImpl) ChangePassword(ctx *gin.Context) {
	payload, err := ParseAndValidateJSON[gen.ChangePasswordJSONRequestBody](ctx)
	if err != nil {
		hdl.BadRequestResponse(ctx, ApiVersion1, err.Error())
		return
	}

	session, ok := hdl.validateSessionFromCookie(ctx)
	if !ok {
		return
	}

	input := dto.ChangePasswordInput{
		AccountID:   session.AccountID.String(),
		OldPassword: payload.OldPassword,
		NewPassword: payload.NewPassword,
	}
	if err := hdl.accountUC.ChangePassword(ctx.Request.Context(), input); err != nil {
		hdl.handleChangePasswordError(ctx, err)
		return
	}

	hdl.NoContentResponse(ctx)
}

func (hdl *AccountHandlerImpl) handleRegisterError(ctx *gin.Context, err error) {
	if errors.Is(err, model.ErrEmailExisted) {
		hdl.ResourceConflictResponse(ctx, ApiVersion1, err.Error())
		return
	}
	hdl.ServerInternalErrResponse(ctx, ApiVersion1)
}

func (hdl *AccountHandlerImpl) sendRegisterSuccess(ctx *gin.Context, account *model.Account) {
	resp := gen.RegisterResponse{
		Version: ApiVersion1,
		Success: true,
		Data: gen.Account{
			Id:        account.ID,
			Email:     account.Email,
			Role:      account.Role,
			CreatedAt: &account.CreatedAt,
			UpdatedAt: &account.UpdatedAt,
		},
	}
	ctx.JSON(http.StatusCreated, resp)
}

func (hdl *AccountHandlerImpl) validateSessionFromCookie(ctx *gin.Context) (*model.Session, bool) {
	sessionID, err := ctx.Cookie("session_id")
	if err != nil {
		hdl.UnauthorizeErrorResponse(ctx, ApiVersion1, http.StatusText(http.StatusUnauthorized))
		return nil, false
	}

	session, err := hdl.sessionUC.ValidateSession(ctx.Request.Context(), sessionID)
	if err != nil {
		if errors.Is(err, model.ErrInvalidSessionID) || errors.Is(err, model.ErrSessionNotFound) {
			hdl.UnauthorizeErrorResponse(ctx, ApiVersion1, http.StatusText(http.StatusUnauthorized))
		} else {
			hdl.logger.Errorf("failed to validate session: %v", err)
			hdl.ServerInternalErrResponse(ctx, ApiVersion1)
		}
		return nil, false
	}

	return session, true
}

func (hdl *AccountHandlerImpl) handleChangePasswordError(ctx *gin.Context, err error) {
	switch {
	case errors.Is(err, usecase.ErrInvalidCredentials):
		hdl.UnauthorizeErrorResponse(ctx, ApiVersion1, err.Error())
	case errors.Is(err, usecase.ErrNewPasswordMustChanged):
		hdl.BadRequestResponse(ctx, ApiVersion1, err.Error())
	default:
		hdl.logger.Errorf("failed to change password: %v", err)
		hdl.ServerInternalErrResponse(ctx, ApiVersion1)
	}
}
