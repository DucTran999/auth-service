package service

import (
	"github.com/DucTran999/auth-service/internal/model"
)

func (b *userBiz) CreateUser() (*model.User, error) {
	user := model.User{Name: "Modi"}

	result, err := b.repo.CreateUser(user)
	if err != nil {
		return nil, err
	}

	return result, nil
}
