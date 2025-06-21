package usecase

import (
	"context"

	"github.com/alexedwards/argon2id"

	"github.com/DucTran999/auth-service/internal/model"
	"github.com/DucTran999/auth-service/internal/repository"
)

type accountUseCaseImpl struct {
	accountRepo repository.AccountRepo
}

func NewAccountUseCase(accountRepo repository.AccountRepo) *accountUseCaseImpl {
	return &accountUseCaseImpl{
		accountRepo: accountRepo,
	}
}

// Register handles new account creation:
// 1. Checks if the email is already in use.
// 2. Hashes the password securely.
// 3. Persists the account to the repository.
func (uc *accountUseCaseImpl) Register(ctx context.Context, input RegisterInput) (*model.Account, error) {
	// Check if the email is already in use
	existing, err := uc.accountRepo.FindByEmail(ctx, input.Email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrEmailExisted
	}

	// Hash the password
	hashedPassword, err := uc.hashPassword(input.Password)
	if err != nil {
		return nil, err
	}

	// Bind input to domain model
	account := model.Account{
		Email:    input.Email,
		Password: hashedPassword,
	}

	// Persist the account
	created, err := uc.accountRepo.Create(ctx, account)
	if err != nil {
		return nil, err
	}

	return created, nil
}

// hashPassword securely hashes a plain password using Argon2id.
func (uc *accountUseCaseImpl) hashPassword(password string) (string, error) {
	return argon2id.CreateHash(password, argon2id.DefaultParams)
}
