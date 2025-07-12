package repository

import (
	"context"
	"errors"

	"github.com/DucTran999/auth-service/internal/errs"
	"github.com/DucTran999/auth-service/internal/model"
	"github.com/DucTran999/auth-service/internal/usecase/port"
	"gorm.io/gorm"
)

type accountRepo struct {
	db *gorm.DB
}

func NewAccountRepo(db *gorm.DB) port.AccountRepo {
	return &accountRepo{db: db}
}

// Create inserts a new account record into the database.
func (r *accountRepo) Create(ctx context.Context, account *model.Account) error {
	return r.db.WithContext(ctx).Create(&account).Error
}

// FindByEmail looks up an account by its email address.
func (r *accountRepo) FindByEmail(ctx context.Context, email string) (*model.Account, error) {
	var account model.Account

	err := r.db.WithContext(ctx).First(&account, "email = ?", email).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrAccountNotFound
		}
		return nil, err
	}

	return &account, nil
}

func (r *accountRepo) FindByID(ctx context.Context, id string) (*model.Account, error) {
	var account model.Account

	err := r.db.WithContext(ctx).First(&account, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &account, nil
}

func (r *accountRepo) UpdatePasswordHash(ctx context.Context, id, passwordHash string) error {
	return r.db.WithContext(ctx).
		Model(&model.Account{}).
		Where("id = ?", id).
		Update("password_hash", passwordHash).
		Error
}
