package service

import (
	"context"

	"github.com/alexedwards/argon2id"

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

// Register handles account creation logic:
//  1. Checks for existing email
//  2. Hashes password securely
//  3. Creates the account in the repository
func (svc *accountServiceImpl) Register(ctx context.Context, accountInfo model.Account) (*model.Account, error) {
	// Step 1: Check if the email already exists
	foundAccount, err := svc.accountRepo.FindByEmail(ctx, accountInfo.Email)
	if err != nil {
		return nil, err
	}
	if foundAccount != nil {
		return nil, ErrEmailExisted
	}

	// Step 2: Hash the user's password before saving
	hashedPassword, err := svc.hashPassword(accountInfo.Password)
	if err != nil {
		return nil, err
	}
	accountInfo.Password = hashedPassword

	// Step 3: Persist the account
	createdAccount, err := svc.accountRepo.Create(ctx, accountInfo)
	if err != nil {
		return nil, err
	}

	return createdAccount, nil
}

// hashPassword securely hashes a plain password using Argon2id.
func (svc *accountServiceImpl) hashPassword(password string) (string, error) {
	return argon2id.CreateHash(password, argon2id.DefaultParams)
}
