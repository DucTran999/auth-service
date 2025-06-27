package usecase

import (
	"context"
	"fmt"

	"github.com/DucTran999/auth-service/internal/domain"
	"github.com/DucTran999/auth-service/pkg/hasher"
)

type accountUseCaseImpl struct {
	hasher      hasher.Hasher
	accountRepo domain.AccountRepo
}

func NewAccountUseCase(hasher hasher.Hasher, accountRepo domain.AccountRepo) *accountUseCaseImpl {
	return &accountUseCaseImpl{
		hasher:      hasher,
		accountRepo: accountRepo,
	}
}

// Register handles new account creation:
// 1. Checks if the email is already in use.
// 2. Hashes the password securely.
// 3. Persists the account to the repository.
func (uc *accountUseCaseImpl) Register(ctx context.Context, input domain.RegisterInput) (*domain.Account, error) {
	taken, err := uc.isEmailTaken(ctx, input.Email)
	if err != nil {
		return nil, err
	}
	if taken {
		return nil, domain.ErrEmailExisted
	}

	// Hash the password
	hashedPassword, err := uc.hasher.HashPassword(input.Password)
	if err != nil {
		return nil, err
	}

	// Bind input to domain model
	account := domain.Account{
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

func (uc *accountUseCaseImpl) ChangePassword(ctx context.Context, input domain.ChangePasswordInput) error {
	account, err := uc.accountRepo.FindByID(ctx, input.AccountID)
	if err != nil {
		return err
	}

	if err = uc.validatePassword(input.OldPassword, account.PasswordHash); err != nil {
		return fmt.Errorf("validate password: %w", err)
	}

	hashedPassword, err := uc.hashIfChanged(input.OldPassword, input.NewPassword)
	if err != nil {
		return fmt.Errorf("hash if changed: %w", err)
	}

	err = uc.accountRepo.UpdatePasswordHash(ctx, account.ID.String(), hashedPassword)
	if err != nil {
		return fmt.Errorf("update new password: %w", err)
	}

	return nil
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

func (uc *accountUseCaseImpl) validatePassword(password, hashed string) error {
	match, err := uc.hasher.ComparePasswordAndHash(password, hashed)
	if err != nil {
		return err
	}
	if !match {
		return ErrInvalidCredentials
	}

	return nil
}

func (uc *accountUseCaseImpl) hashIfChanged(oldPassword, newPassword string) (string, error) {
	if oldPassword == newPassword {
		return "", ErrNewPasswordMustChanged
	}

	return uc.hasher.HashPassword(newPassword)
}
