package service

import (
	"context"

	"github.com/DucTran999/auth-service/internal/model"
	"github.com/DucTran999/auth-service/internal/repository"
)

type IUserService interface {
	RegisterUser(ctx context.Context, userInfo model.User) (*model.User, error)
}

type userBiz struct {
	userRepo repository.IUserRepo
}

func NewUserBiz(ur repository.IUserRepo) *userBiz {
	return &userBiz{
		userRepo: ur,
	}
}
