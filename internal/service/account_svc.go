package service

import (
	"context"

	"github.com/DucTran999/auth-service/internal/common"
	"github.com/DucTran999/auth-service/internal/model"
	"github.com/DucTran999/auth-service/internal/repository"
)

// AccountService defines the business logic for managing user accounts.
type AccountService interface {
	// Register creates a new user account with the provided information.
	// It typically includes validation, password hashing, and persistence logic.
	Register(ctx context.Context, info model.Account) (*model.Account, error)
}

type accountServiceImpl struct {
	accountRepo repository.AccountRepo
}

func NewAccountService(accountRepo repository.AccountRepo) *accountServiceImpl {
	return &accountServiceImpl{
		accountRepo: accountRepo,
	}
}

func (svc *accountServiceImpl) Register(ctx context.Context, userInfo model.Account) (*model.Account, error) {
	foundAccount, err := svc.accountRepo.FindByEmail(ctx, userInfo.Email)
	if err != nil {
		return nil, err
	}

	if foundAccount != nil {
		return nil, common.ErrEmailExisted
	}

	user, err := svc.accountRepo.Create(ctx, userInfo)
	if err != nil {
		return nil, err
	}

	return user, nil
}
