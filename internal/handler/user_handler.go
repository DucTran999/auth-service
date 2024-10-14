package handler

import (
	"net/http"

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
	result, err := h.service.CreateUser()
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
