package handler

import (
	"errors"
	"net/http"

	"github.com/DucTran999/auth-service/internal/common"
	"github.com/DucTran999/auth-service/internal/gen"
	"github.com/DucTran999/auth-service/internal/model"
	"github.com/DucTran999/auth-service/internal/service"
	"github.com/gin-gonic/gin"
)

type AccountHandler interface {
	CreateAccount(ctx *gin.Context)
}

type accountHandlerImpl struct {
	BaseHandler
	service service.AccountService
}

func NewAccountHandler(accountSvc service.AccountService) *accountHandlerImpl {
	return &accountHandlerImpl{
		service: accountSvc,
	}
}

// CreateAccount handles the HTTP request to register a new account.
func (h *accountHandlerImpl) CreateAccount(ctx *gin.Context) {
	// Parse request body
	var payload gen.CreateAccountJSONRequestBody
	if err := ctx.Bind(&payload); err != nil {
		h.BadRequestResponse(ctx, common.ApiVersion1, err)
		return
	}

	// Convert request to domain model
	accountInfo := model.Account{
		Email:    payload.Email,
		Password: payload.Password,
	}

	// Attempt registration
	account, err := h.service.Register(ctx, accountInfo)
	if errors.Is(err, common.ErrEmailExisted) {
		h.ResourceConflictResponse(ctx, common.ApiVersion1, err)
		return
	}
	if err != nil {
		h.ServerInternalErrResponse(ctx, common.ApiVersion1)
		return
	}

	// Prepare response
	respData := gen.AccountResponse{
		Version: common.ApiVersion1,
		Success: true,
		Data: gen.Account{
			Id:        account.ID,
			Email:     account.Email,
			CreatedAt: account.CreatedAt,
			UpdatedAt: account.UpdatedAt,
		},
	}
	ctx.JSON(http.StatusOK, respData)
}
