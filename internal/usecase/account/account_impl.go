package account

import (
	"context"
	"fmt"

	"github.com/DucTran999/auth-service/internal/errs"
	"github.com/DucTran999/auth-service/internal/model"
	"github.com/DucTran999/auth-service/internal/usecase/dto"
	"github.com/DucTran999/auth-service/internal/usecase/port"
	"github.com/DucTran999/auth-service/pkg/hasher"
)

type accountUsecase struct {
	hasher      hasher.Hasher
	accountRepo port.AccountRepo
}

func NewAccountUseCase(hasher hasher.Hasher, accountRepo port.AccountRepo) port.AccountUsecase {
	return &accountUsecase{
		hasher:      hasher,
		accountRepo: accountRepo,
	}
}

// Register handles new account creation:
// 1. Checks if the email is already in use.
// 2. Hashes the password securely.
// 3. Persists the account to the repository.
func (uc *accountUsecase) Register(ctx context.Context, input dto.RegisterInput) (*model.Account, error) {
	taken, err := uc.isEmailTaken(ctx, input.Email)
	if err != nil {
		return nil, err
	}
	if taken {
		return nil, errs.ErrEmailExisted
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
	if err := uc.accountRepo.Create(ctx, &account); err != nil {
		return nil, err
	}

	return &account, nil
}

func (uc *accountUsecase) ChangePassword(ctx context.Context, input dto.ChangePasswordInput) error {
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
func (uc *accountUsecase) isEmailTaken(ctx context.Context, email string) (bool, error) {
	account, err := uc.accountRepo.FindByEmail(ctx, email)
	if err != nil {
		return false, err
	}

	return account != nil, nil
}

func (uc *accountUsecase) validatePassword(password, hashed string) error {
	match, err := uc.hasher.ComparePasswordAndHash(password, hashed)
	if err != nil {
		return err
	}
	if !match {
		return errs.ErrInvalidCredentials
	}

	return nil
}

func (uc *accountUsecase) hashIfChanged(oldPassword, newPassword string) (string, error) {
	if oldPassword == newPassword {
		return "", errs.ErrNewPasswordMustChanged
	}

	return uc.hasher.HashPassword(newPassword)
}
