package handler

import (
	"errors"
	"net/http"

	"github.com/DucTran999/auth-service/internal/common"
	"github.com/DucTran999/auth-service/internal/dto"
	"github.com/DucTran999/auth-service/internal/model"
	service "github.com/DucTran999/auth-service/internal/service/user"
	"github.com/gin-gonic/gin"
)

type IUserHandler interface {
	CreateUser(ctx *gin.Context)
}

type userHandler struct {
	baseHandler
	service service.IUserService
}

func newUserHandler(us service.IUserService) *userHandler {
	return &userHandler{
		service: us,
	}
}

func (h *userHandler) CreateUser(ctx *gin.Context) {
	payload := new(dto.CreateUserRequest)
	if err := ctx.Bind(payload); err != nil {
		h.JsonResponse(ctx, http.StatusBadRequest, nil, err.Error())
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
		h.JsonResponse(ctx, http.StatusConflict, nil, "email exited")
		return
	}

	if err != nil {
		h.JsonResponse(ctx, http.StatusInternalServerError, nil, common.MessageInternalErr)
		return
	}

	h.JsonResponse(ctx, http.StatusOK, dto.CreateUserResp{ID: result.ID}, "")
}
