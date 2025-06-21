package usecase

import (
	"context"

	"github.com/DucTran999/auth-service/internal/model"
	"github.com/DucTran999/auth-service/internal/repository"
	"github.com/DucTran999/auth-service/pkg"
)

type authUseCaseImpl struct {
	hasher      pkg.Hasher
	accountRepo repository.AccountRepo
	sessionRepo repository.SessionRepository
}

func NewAuthUseCase(
	hasher pkg.Hasher,
	accountRepo repository.AccountRepo,
	sessionRepo repository.SessionRepository,
) *authUseCaseImpl {
	return &authUseCaseImpl{
		hasher:      hasher,
		accountRepo: accountRepo,
		sessionRepo: sessionRepo,
	}
}

// Login authenticates a user using email and password.
// It verifies credentials, checks account status, and creates a new session on success.
func (uc *authUseCaseImpl) Login(ctx context.Context, input LoginInput) (*model.Session, error) {
	// Step 1: Retrieve account by email
	account, err := uc.findAccountByEmail(ctx, input.Email)
	if err != nil {
		return nil, err
	}

	// Step 2: Check if account is active (optional, depending on business rule)
	if !account.IsActive {
		return nil, ErrAccountDisabled
	}

	// Step 3: Verify password
	match, err := uc.hasher.ComparePasswordAndHash(input.Password, account.PasswordHash)
	if err != nil {
		return nil, err
	}
	if !match {
		return nil, ErrInvalidCredentials
	}

	// Step 4: Create session (session storage, expiry, metadata, etc.)
	session := &model.Session{
		AccountID: account.ID,
		Account: model.Account{
			ID:       account.ID,
			Email:    account.Email,
			Role:     account.Role,
			IsActive: account.IsActive,
		},
		IPAddress: input.IP,
		UserAgent: input.UserAgent,
	}
	if err := uc.sessionRepo.Create(ctx, session); err != nil {
		return nil, err
	}

	return session, nil
}

func (uc *authUseCaseImpl) findAccountByEmail(ctx context.Context, email string) (*model.Account, error) {
	account, err := uc.accountRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if account == nil {
		return nil, ErrInvalidCredentials
	}

	return account, nil
}
