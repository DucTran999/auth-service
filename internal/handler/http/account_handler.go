package http

import (
	"errors"
	"net/http"

	"github.com/DucTran999/auth-service/internal/gen"
	"github.com/DucTran999/auth-service/internal/model"
	"github.com/DucTran999/auth-service/internal/usecase"
	"github.com/DucTran999/shared-pkg/logger"
	"github.com/gin-gonic/gin"
)

type AccountHandler interface {
	CreateAccount(ctx *gin.Context)
	ChangePassword(ctx *gin.Context)
}

type accountHandlerImpl struct {
	logger logger.ILogger

	BaseHandler
	accountUC usecase.AccountUseCase
	sessionUC usecase.SessionUsecase
}

func NewAccountHandler(
	logger logger.ILogger,
	accountUC usecase.AccountUseCase,
	sessionUC usecase.SessionUsecase,
) *accountHandlerImpl {
	return &accountHandlerImpl{
		logger:    logger,
		accountUC: accountUC,
		sessionUC: sessionUC,
	}
}

// CreateAccount handles the HTTP request to register a new account.
func (hdl *accountHandlerImpl) CreateAccount(ctx *gin.Context) {
	payload, err := ParseAndValidateJSON[gen.CreateAccountJSONRequestBody](ctx)
	if err != nil {
		hdl.BadRequestResponse(ctx, ApiVersion1, err.Error())
		return
	}

	input := usecase.RegisterInput{
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

func (hdl *accountHandlerImpl) ChangePassword(ctx *gin.Context) {
	payload, err := ParseAndValidateJSON[gen.ChangePasswordJSONRequestBody](ctx)
	if err != nil {
		hdl.BadRequestResponse(ctx, ApiVersion1, err.Error())
		return
	}

	session, ok := hdl.validateSessionFromCookie(ctx)
	if !ok {
		return
	}

	input := usecase.ChangePasswordInput{
		AccountID:   session.AccountID.String(),
		OldPassword: payload.OldPassword,
		NewPassword: payload.NewPassword,
	}
	if err := hdl.accountUC.ChangePassword(ctx.Request.Context(), input); err != nil {
		hdl.logger.Errorf("failed to change password: %v", err)
		hdl.ServerInternalErrResponse(ctx, ApiVersion1)
		return
	}

	// Change password successfully
	hdl.NoContentResponse(ctx)
}

func (hdl *accountHandlerImpl) handleRegisterError(ctx *gin.Context, err error) {
	if errors.Is(err, usecase.ErrEmailExisted) {
		hdl.ResourceConflictResponse(ctx, ApiVersion1, err.Error())
		return
	}
	hdl.ServerInternalErrResponse(ctx, ApiVersion1)
}

func (hdl *accountHandlerImpl) sendRegisterSuccess(ctx *gin.Context, account *model.Account) {
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

func (hdl *accountHandlerImpl) validateSessionFromCookie(ctx *gin.Context) (*model.Session, bool) {
	sessionID, err := ctx.Cookie("session_id")
	if err != nil {
		hdl.UnauthorizeErrorResponse(ctx, ApiVersion1, http.StatusText(http.StatusUnauthorized))
		return nil, false
	}

	session, err := hdl.sessionUC.ValidateSession(ctx.Request.Context(), sessionID)
	if err != nil {
		if errors.Is(err, usecase.ErrInvalidSessionID) || errors.Is(err, usecase.ErrSessionNotFound) {
			hdl.UnauthorizeErrorResponse(ctx, ApiVersion1, http.StatusText(http.StatusUnauthorized))
		} else {
			hdl.logger.Errorf("failed to validate session: %v", err)
			hdl.ServerInternalErrResponse(ctx, ApiVersion1)
		}
		return nil, false
	}

	return session, true
}
