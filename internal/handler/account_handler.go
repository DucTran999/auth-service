package handler

import (
	"errors"

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
	userInfo := model.Account{
		Email:    payload.Email,
		Password: payload.Password,
	}

	// Attempt registration
	account, err := h.service.Register(ctx, userInfo)
	if errors.Is(err, common.ErrEmailExisted) {
		h.ResourceConflictResponse(ctx, common.ApiVersion1)
		return
	}
	if err != nil {
		h.ServerInternalErrResponse(ctx, common.ApiVersion1)
		return
	}

	// Prepare response
	respData := gen.Account{
		Id:        account.ID,
		Email:     account.Email,
		CreatedAt: account.CreatedAt,
		UpdatedAt: account.UpdatedAt,
	}
	h.SuccessResponse(ctx, common.ApiVersion1, respData)
}
