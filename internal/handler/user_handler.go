package handler

import (
	"errors"

	"github.com/DucTran999/auth-service/internal/common"
	"github.com/DucTran999/auth-service/internal/dto"
	"github.com/DucTran999/auth-service/internal/model"
	service "github.com/DucTran999/auth-service/internal/service/user"
	"github.com/gin-gonic/gin"
)

type UserHandler interface {
	CreateUser(ctx *gin.Context)
}

type userHandlerImpl struct {
	BaseHandler
	service service.IUserService
}

func newUserHandler(us service.IUserService) *userHandlerImpl {
	return &userHandlerImpl{
		service: us,
	}
}

func (h *userHandlerImpl) CreateUser(ctx *gin.Context) {
	payload := new(dto.CreateUserRequest)
	if err := ctx.Bind(payload); err != nil {
		h.BadRequestResponse(ctx, common.ApiVersion1, err)
		return
	}

	userInfo := model.User{
		Username:  payload.Username,
		Password:  payload.Password,
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Email:     payload.Email,
		Gender:    payload.Gender,
	}
	result, err := h.service.RegisterUser(ctx.Request.Context(), userInfo)
	if errors.Is(err, common.ErrEmailExisted) {
		h.ResourceConflictResponse(ctx, common.ApiVersion1)
		return
	}

	if err != nil {
		h.ServerInternalErrResponse(ctx, common.ApiVersion1)
		return
	}

	h.SuccessResponse(ctx, common.ApiVersion1, dto.CreateUserResp{ID: result.ID})
}
