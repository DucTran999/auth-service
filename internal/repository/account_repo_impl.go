package repository

import (
	"context"
	"errors"

	"github.com/DucTran999/auth-service/internal/model"
	"gorm.io/gorm"
)

type accountRepoImpl struct {
	db *gorm.DB
}

func NewAccountRepo(db *gorm.DB) *accountRepoImpl {
	return &accountRepoImpl{
		db: db,
	}
}

// Create inserts a new account record into the database.
func (r *accountRepoImpl) Create(ctx context.Context, account model.Account) (*model.Account, error) {
	if err := r.db.WithContext(ctx).Create(&account).Error; err != nil {
		return nil, err
	}

	return &account, nil
}

// FindByEmail looks up an account by its email address.
func (r *accountRepoImpl) FindByEmail(ctx context.Context, email string) (*model.Account, error) {
	var account model.Account

	err := r.db.WithContext(ctx).Table(account.TableName()).First(&account, "email = ?", email).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &account, nil
}
