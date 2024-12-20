package service

import (
	"context"

	"github.com/DucTran999/auth-service/internal/common"
	"github.com/DucTran999/auth-service/internal/model"
)

func (b *userBiz) RegisterUser(ctx context.Context, userInfo model.User) (*model.User, error) {
	foundUser, err := b.userRepo.GetUserByEmail(ctx, userInfo.Email)
	if err != nil {
		return nil, err
	}

	if foundUser != nil {
		return nil, common.ErrEmailExisted
	}

	user, err := b.userRepo.CreateUser(ctx, userInfo)
	if err != nil {
		return nil, err
	}

	return user, nil
}
