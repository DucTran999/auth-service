package handler

import (
	"github.com/DucTran999/auth-service/internal/common"
	"github.com/DucTran999/auth-service/internal/gen"
	"github.com/DucTran999/auth-service/internal/service"
	"github.com/gin-gonic/gin"
)

type AccountHandler interface {
	CreateAccount(ctx *gin.Context)
}

type userHandlerImpl struct {
	BaseHandler
	service service.IUserService
}

func NewAccountHandler(us service.IUserService) *userHandlerImpl {
	return &userHandlerImpl{
		service: us,
	}
}

func (h *userHandlerImpl) CreateAccount(ctx *gin.Context) {
	payload := new(gen.CreateAccountJSONRequestBody)
	if err := ctx.Bind(payload); err != nil {
		h.BadRequestResponse(ctx, common.ApiVersion1, err)
		return
	}

	h.SuccessResponse(ctx, common.ApiVersion1, gen.AccountResponse{})
}
