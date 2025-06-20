package handler

import (
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

func (h *accountHandlerImpl) CreateAccount(ctx *gin.Context) {
	payload := new(gen.CreateAccountJSONRequestBody)
	if err := ctx.Bind(payload); err != nil {
		h.BadRequestResponse(ctx, common.ApiVersion1, err)
		return
	}

	// Convert request to model
	userInfo := model.Account{
		Email:    payload.Email,
		Password: payload.Password,
	}

	// Register user
	user, err := h.service.Register(ctx, userInfo)
	if err != nil {
		h.ServerInternalErrResponse(ctx, common.ApiVersion1)
		return
	}

	data := gen.Account{
		Id:    user.ID,
		Email: user.Email,
	}

	h.SuccessResponse(ctx, common.ApiVersion1, data)
}
