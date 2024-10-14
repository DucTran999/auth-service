package repository

import (
	"context"

	"github.com/DucTran999/auth-service/internal/model"
	"gorm.io/gorm"
)

type IUserRepo interface {
	CreateUser(user model.User) (*model.User, error)
}

type userRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *userRepo {
	return &userRepo{
		db: db,
	}
}

func (r *userRepo) CreateUser(user model.User) (*model.User, error) {
	if err := r.db.WithContext(context.Background()).Create(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}
