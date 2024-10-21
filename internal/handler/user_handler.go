package handler

import (
	"net/http"

	"github.com/DucTran999/auth-service/internal/dto"
	"github.com/DucTran999/auth-service/internal/model"
	service "github.com/DucTran999/auth-service/internal/service/user"
	"github.com/gin-gonic/gin"
)

type IUserHandler interface {
	CreateUser(ctx *gin.Context)
}

type userHandler struct {
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
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  1,
			"message": err.Error(),
			"data":    nil,
		})
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
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"status":  1,
			"message": "failed when creating user",
			"data":    nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  0,
		"message": "new user created",
		"data":    result.ID,
	})
}
