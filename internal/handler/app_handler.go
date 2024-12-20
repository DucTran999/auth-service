package handler

import (
	"github.com/DucTran999/auth-service/internal/registry"
	"github.com/DucTran999/auth-service/internal/repository"
	service "github.com/DucTran999/auth-service/internal/service/user"
)

type AppHandler struct {
	IUserHandler
}

func NewAppHandler(reg *registry.Registry) AppHandler {
	userRepo := repository.NewUserRepo(reg.PostgresDB)
	userBiz := service.NewUserBiz(userRepo)

	return AppHandler{
		IUserHandler: newUserHandler(userBiz),
	}
}
