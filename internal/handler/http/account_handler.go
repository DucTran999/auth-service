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
	payload, err := hdl.parseAndValidateCreateAccount(ctx)
	if err != nil {
		return // response already handled
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

func (hdl *accountHandlerImpl) parseAndValidateCreateAccount(ctx *gin.Context,
) (*gen.CreateAccountJSONRequestBody, error) {

	var payload gen.CreateAccountJSONRequestBody
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
