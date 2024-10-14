package service

import (
	"github.com/DucTran999/auth-service/internal/model"
	"github.com/DucTran999/auth-service/internal/repository"
)

type IUserService interface {
	CreateUser() (*model.User, error)
}

type userBiz struct {
	repo repository.IUserRepo
}

func NewUserBiz(ur repository.IUserRepo) *userBiz {
	return &userBiz{
		repo: ur,
	}
}
