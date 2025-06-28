package rest

import (
	"context"
	"errors"
	"net/http"

	"github.com/DucTran999/auth-service/internal/domain"
	"github.com/DucTran999/auth-service/internal/gen"
	"github.com/DucTran999/auth-service/internal/usecase"
	"github.com/DucTran999/auth-service/internal/usecase/dto"
	"github.com/DucTran999/shared-pkg/logger"
	"github.com/gin-gonic/gin"
)

// AccountUseCase defines the business logic for managing user accounts.
type AccountUseCase interface {
	// Register creates a new user account with the provided information.
	// It typically includes validation, password hashing, and persistence logic.
	Register(ctx context.Context, input dto.RegisterInput) (*domain.Account, error)

	// ChangePassword change password for user when old password are match
	ChangePassword(ctx context.Context, input dto.ChangePasswordInput) error
}

// SessionUsecase defines business logic operations related to session lifecycle management.
type SessionUsecase interface {
	// ValidateSession find session in cache first if not try to lookup in DB.
	// Return session only if it is existed and not expire
	ValidateSession(ctx context.Context, sessionID string) (*domain.Session, error)
}

type AccountHandlerImpl struct {
	logger logger.ILogger

	BaseHandler
	accountUC AccountUseCase
	sessionUC SessionUsecase
}

func NewAccountHandler(
	logger logger.ILogger,
	accountUC AccountUseCase,
	sessionUC SessionUsecase,
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
	if errors.Is(err, domain.ErrEmailExisted) {
		hdl.ResourceConflictResponse(ctx, ApiVersion1, err.Error())
		return
	}
	hdl.ServerInternalErrResponse(ctx, ApiVersion1)
}

func (hdl *AccountHandlerImpl) sendRegisterSuccess(ctx *gin.Context, account *domain.Account) {
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

func (hdl *AccountHandlerImpl) validateSessionFromCookie(ctx *gin.Context) (*domain.Session, bool) {
	sessionID, err := ctx.Cookie("session_id")
	if err != nil {
		hdl.UnauthorizeErrorResponse(ctx, ApiVersion1, http.StatusText(http.StatusUnauthorized))
		return nil, false
	}

	session, err := hdl.sessionUC.ValidateSession(ctx.Request.Context(), sessionID)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidSessionID) || errors.Is(err, domain.ErrSessionNotFound) {
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
