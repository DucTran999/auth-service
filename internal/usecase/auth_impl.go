package usecase

import (
	"context"
	"time"

	"github.com/DucTran999/auth-service/internal/common"
	"github.com/DucTran999/auth-service/internal/model"
	"github.com/DucTran999/auth-service/internal/repository"
	"github.com/DucTran999/auth-service/pkg"
	"github.com/google/uuid"
)

const (
	sessionDuration = 60 * time.Minute
)

type authUseCaseImpl struct {
	hasher      pkg.Hasher
	cache       pkg.Cache
	accountRepo repository.AccountRepo
	sessionRepo repository.SessionRepository
}

func NewAuthUseCase(
	hasher pkg.Hasher,
	cache pkg.Cache,
	accountRepo repository.AccountRepo,
	sessionRepo repository.SessionRepository,
) *authUseCaseImpl {
	return &authUseCaseImpl{
		hasher:      hasher,
		cache:       cache,
		accountRepo: accountRepo,
		sessionRepo: sessionRepo,
	}
}

// Login authenticates a user using email and password.
// It verifies credentials, checks account status, and creates a new session on success.
func (uc *authUseCaseImpl) Login(ctx context.Context, input LoginInput) (*model.Session, error) {
	session, err := uc.tryReuseSession(ctx, input.CurrentSessionID)
	if err != nil {
		return nil, err
	}
	if session != nil {
		return session, nil
	}

	account, err := uc.findAccountByEmail(ctx, input.Email)
	if err != nil {
		return nil, err
	}

	if err := uc.checkAccountActive(account); err != nil {
		return nil, err
	}

	if err := uc.verifyPassword(input.Password, account.PasswordHash); err != nil {
		return nil, err
	}

	return uc.createSession(ctx, account, input)
}

// tryReuseSession checks if the current session is valid and updates its expiration.
// Returns the updated session if reusable; otherwise, returns nil.
func (uc *authUseCaseImpl) tryReuseSession(ctx context.Context, sessionID string) (*model.Session, error) {
	if sessionID == "" || sessionID == uuid.Nil.String() {
		return nil, nil
	}

	// Try get session from cache
	sessionKey := common.KeyFromSessionID(sessionID)
	cachedSession := uc.getSessionFromCache(ctx, sessionKey)
	if cachedSession != nil {
		ttl, _ := uc.cache.TTL(ctx, sessionKey)

		// Example: Refresh TTL if less than 30% of sessionDuration remains
		if ttl < int64(sessionDuration.Seconds())/3 {
			_ = uc.cache.Expire(ctx, sessionKey, sessionDuration)
		}
		return cachedSession, nil
	}

	// Cache missing try to retrieve from DB
	session, err := uc.sessionRepo.FindByID(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	if session == nil || session.IsExpired() {
		return nil, nil
	}

	// Re-cache the session after extend ttl
	_ = uc.cache.Set(ctx, sessionKey, session, sessionDuration)
	return session, nil
}

func (uc *authUseCaseImpl) findAccountByEmail(
	ctx context.Context,
	email string,
) (*model.Account, error) {

	account, err := uc.accountRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if account == nil {
		return nil, ErrInvalidCredentials
	}

	return account, nil
}

func (uc *authUseCaseImpl) checkAccountActive(account *model.Account) error {
	if !account.IsActive {
		return ErrAccountDisabled
	}
	return nil
}

func (uc *authUseCaseImpl) verifyPassword(plain, hashed string) error {
	match, err := uc.hasher.ComparePasswordAndHash(plain, hashed)
	if err != nil {
		return err
	}
	if !match {
		return ErrInvalidCredentials
	}
	return nil
}

func (uc *authUseCaseImpl) createSession(
	ctx context.Context,
	account *model.Account,
	input LoginInput,
) (*model.Session, error) {

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

func (uc *authUseCaseImpl) getSessionFromCache(
	ctx context.Context,
	sessionKey string,
) *model.Session {

	var session model.Session
	if err := uc.cache.GetInto(ctx, sessionKey, &session); err != nil {
		// If get session from cache got error just pass it.
		// Already has fallback form DB
		return nil
	}

	return &session
}
