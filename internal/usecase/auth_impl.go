package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/DucTran999/auth-service/internal/domain"
	"github.com/DucTran999/auth-service/pkg/cache"
	"github.com/DucTran999/auth-service/pkg/hasher"
	"github.com/google/uuid"
)

const (
	sessionDuration = 60 * time.Minute
)

type authUseCaseImpl struct {
	hasher      hasher.Hasher
	cache       cache.Cache
	accountRepo domain.AccountRepo
	sessionRepo domain.SessionRepository
}

func NewAuthUseCase(
	hasher hasher.Hasher,
	cache cache.Cache,
	accountRepo domain.AccountRepo,
	sessionRepo domain.SessionRepository,
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
func (uc *authUseCaseImpl) Login(ctx context.Context, input domain.LoginInput) (*domain.Session, error) {
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

func (uc *authUseCaseImpl) Logout(ctx context.Context, sessionID string) error {
	// Fast check sessionID must uuid
	if _, err := uuid.Parse(sessionID); err != nil {
		return fmt.Errorf("logout: %w session=%s error=%w", domain.ErrInvalidSessionID, sessionID, err)
	}

	// Remove session in cache
	_ = uc.cache.Del(ctx, cache.KeyFromSessionID(sessionID))

	// expire the session in db
	err := uc.sessionRepo.UpdateExpiresAt(ctx, sessionID, time.Now())
	if err != nil {
		return fmt.Errorf("logout: failed to update expires_at in DB for session=%s, error=%w", sessionID, err)
	}

	return nil
}

// tryReuseSession checks if the current session is valid and updates its expiration.
// Returns the updated session if reusable; otherwise, returns nil.
func (uc *authUseCaseImpl) tryReuseSession(ctx context.Context, sessionID string) (*domain.Session, error) {
	if sessionID == "" || sessionID == uuid.Nil.String() {
		return nil, nil
	}

	// Try get session from cache
	sessionKey := cache.KeyFromSessionID(sessionID)
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
) (*domain.Account, error) {

	account, err := uc.accountRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if account == nil {
		return nil, ErrInvalidCredentials
	}

	return account, nil
}

func (uc *authUseCaseImpl) checkAccountActive(account *domain.Account) error {
	if !account.IsActive {
		return domain.ErrAccountDisabled
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
	account *domain.Account,
	input domain.LoginInput,
) (*domain.Session, error) {

	session := &domain.Session{
		AccountID: account.ID,
		Account: domain.Account{
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

	_ = uc.cache.Set(ctx, cache.KeyFromSessionID(session.ID.String()), session, sessionDuration)

	return session, nil
}

func (uc *authUseCaseImpl) getSessionFromCache(
	ctx context.Context,
	sessionKey string,
) *domain.Session {

	var session domain.Session
	if err := uc.cache.GetInto(ctx, sessionKey, &session); err != nil {
		// If get session from cache got error just pass it.
		// Already has fallback form DB
		return nil
	}

	return &session
}
