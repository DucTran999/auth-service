package usecase

import (
	"context"

	"github.com/DucTran999/auth-service/internal/model"
	"github.com/DucTran999/auth-service/internal/repository"
	"github.com/DucTran999/auth-service/pkg"
)

type accountUseCaseImpl struct {
	hasher      pkg.Hasher
	accountRepo repository.AccountRepo
}

func NewAccountUseCase(hasher pkg.Hasher, accountRepo repository.AccountRepo) *accountUseCaseImpl {
	return &accountUseCaseImpl{
		hasher:      hasher,
		accountRepo: accountRepo,
	}
}

// Register handles new account creation:
// 1. Checks if the email is already in use.
// 2. Hashes the password securely.
// 3. Persists the account to the repository.
func (uc *accountUseCaseImpl) Register(ctx context.Context, input RegisterInput) (*model.Account, error) {
	taken, err := uc.isEmailTaken(ctx, input.Email)
	if err != nil {
		return nil, err
	}
	if taken {
		return nil, ErrEmailExisted
	}

	// Hash the password
	hashedPassword, err := uc.hasher.HashPassword(input.Password)
	if err != nil {
		return nil, err
	}

	// Bind input to domain model
	account := model.Account{
		Email:        input.Email,
		PasswordHash: hashedPassword,
	}

	// Persist the account
	created, err := uc.accountRepo.Create(ctx, account)
	if err != nil {
		return nil, err
	}

	return created, nil
}

// isEmailTaken checks if the provided email already exists in the system.
// Returns ErrEmailExisted if a duplicate is found, or a repository error if any occurs.
func (uc *accountUseCaseImpl) isEmailTaken(ctx context.Context, email string) (bool, error) {
	account, err := uc.accountRepo.FindByEmail(ctx, email)
	if err != nil {
		return false, err
	}

	return account != nil, nil
}
