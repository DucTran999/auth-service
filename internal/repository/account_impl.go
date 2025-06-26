package repository

import (
	"context"
	"errors"

	"github.com/DucTran999/auth-service/internal/domain"
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
func (r *accountRepoImpl) Create(ctx context.Context, account domain.Account) (*domain.Account, error) {
	if err := r.db.WithContext(ctx).Create(&account).Error; err != nil {
		return nil, err
	}

	return &account, nil
}

// FindByEmail looks up an account by its email address.
func (r *accountRepoImpl) FindByEmail(ctx context.Context, email string) (*domain.Account, error) {
	var account domain.Account

	err := r.db.WithContext(ctx).First(&account, "email = ?", email).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &account, nil
}

func (r *accountRepoImpl) FindByID(ctx context.Context, id string) (*domain.Account, error) {
	var account domain.Account

	err := r.db.WithContext(ctx).First(&account, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &account, nil
}

func (r *accountRepoImpl) UpdatePasswordHash(ctx context.Context, id, passwordHash string) error {
	return r.db.WithContext(ctx).
		Model(&domain.Account{}).
		Where("id = ?", id).
		Update("password_hash", passwordHash).
		Error
}
