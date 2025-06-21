package handler

import (
	"errors"
	"net/http"

	"github.com/DucTran999/auth-service/internal/gen"
	"github.com/DucTran999/auth-service/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type AccountHandler interface {
	CreateAccount(ctx *gin.Context)
}

type accountHandlerImpl struct {
	BaseHandler
	accountUC usecase.AccountUseCase
}

func NewAccountHandler(accountUC usecase.AccountUseCase) *accountHandlerImpl {
	return &accountHandlerImpl{
		accountUC: accountUC,
	}
}

// CreateAccount handles the HTTP request to register a new account.
func (hdl *accountHandlerImpl) CreateAccount(ctx *gin.Context) {
	// Parse request body
	var payload gen.CreateAccountJSONRequestBody
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		var ve validator.ValidationErrors
		// Get the first validation error
		if errors.As(err, &ve) && len(ve) > 0 {
			hdl.ValidateErrorResponse(ctx, ApiVersion1, validationErrorMessage(ve[0]))
			return
		}

		// Fallback for other types of errors
		hdl.BadRequestResponse(ctx, ApiVersion1, err.Error())
		return
	}

	registerInfo := usecase.RegisterInput{
		Email:    string(payload.Email),
		Password: payload.Password,
	}

	// Attempt registration
	account, err := hdl.accountUC.Register(ctx, registerInfo)
	if errors.Is(err, usecase.ErrEmailExisted) {
		hdl.ResourceConflictResponse(ctx, ApiVersion1, err.Error())
		return
	}
	if err != nil {
		hdl.ServerInternalErrResponse(ctx, ApiVersion1)
		return
	}

	// Prepare response
	respData := gen.RegisterResponse{
		Version: ApiVersion1,
		Success: true,
		Data: gen.Account{
			Id:        account.Id,
			Email:     account.Email,
			Role:      account.Role,
			CreatedAt: &account.CreatedAt,
			UpdatedAt: &account.UpdatedAt,
		},
	}
	ctx.JSON(http.StatusCreated, respData)
}
