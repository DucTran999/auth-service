package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/DucTran999/auth-service/internal/model"
	"github.com/DucTran999/auth-service/internal/usecase/dto"
	"github.com/DucTran999/auth-service/internal/usecase/port"
	"github.com/DucTran999/auth-service/internal/usecase/shared"
	"github.com/DucTran999/auth-service/pkg/cache"
	"github.com/google/uuid"
)

const (
	SessionLifetime = 60 * time.Minute
)

type authSessionUsecase struct {
	cache           cache.Cache
	accountVerifier shared.AccountVerifier
	accountRepo     port.AccountRepo
	sessionRepo     port.SessionRepository
}

func NewAuthSessionUsecase(
	cache cache.Cache,
	accountVerifier shared.AccountVerifier,
	accountRepo port.AccountRepo,
	sessionRepo port.SessionRepository,
) port.AuthSessionUsecase {
	return &authSessionUsecase{
		cache:           cache,
		accountVerifier: accountVerifier,
		accountRepo:     accountRepo,
		sessionRepo:     sessionRepo,
	}
}

// Login authenticates a user using email and password.
// It verifies credentials, checks account status, and creates a new session on success.
func (uc *authSessionUsecase) Login(ctx context.Context, input dto.LoginInput) (*model.Session, error) {
	session, err := uc.tryReuseSession(ctx, input.CurrentSessionID)
	if err != nil {
		return nil, err
	}
	if session != nil {
		return session, nil
	}

	account, err := uc.accountVerifier.Verify(ctx, input.Email, input.Password)
	if err != nil {
		return nil, err
	}

	return uc.createSession(ctx, account, input)
}

func (uc *authSessionUsecase) Logout(ctx context.Context, sessionID string) error {
	// Fast check sessionID must uuid
	if _, err := uuid.Parse(sessionID); err != nil {
		return fmt.Errorf("logout: %w session=%s error=%w", model.ErrInvalidSessionID, sessionID, err)
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
func (uc *authSessionUsecase) tryReuseSession(ctx context.Context, sessionID string) (*model.Session, error) {
	if sessionID == "" || sessionID == uuid.Nil.String() {
		return nil, nil
	}

	// Try get session from cache
	sessionKey := cache.KeyFromSessionID(sessionID)
	cachedSession := uc.getSessionFromCache(ctx, sessionKey)
	if cachedSession != nil {
		ttl, _ := uc.cache.TTL(ctx, sessionKey)

		// Example: Refresh TTL if less than 30% of SessionLifetime remains
		if ttl < int64(SessionLifetime.Seconds())/3 {
			_ = uc.cache.Expire(ctx, sessionKey, SessionLifetime)
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
	_ = uc.cache.Set(ctx, sessionKey, session, SessionLifetime)
	return session, nil
}

func (uc *authSessionUsecase) createSession(
	ctx context.Context,
	account *model.Account,
	input dto.LoginInput,
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

	_ = uc.cache.Set(ctx, cache.KeyFromSessionID(session.ID.String()), session, SessionLifetime)

	return session, nil
}

func (uc *authSessionUsecase) getSessionFromCache(
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
