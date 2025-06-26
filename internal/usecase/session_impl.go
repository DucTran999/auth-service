package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/DucTran999/auth-service/internal/common"
	"github.com/DucTran999/auth-service/internal/model"
	"github.com/DucTran999/auth-service/internal/repository"
	"github.com/DucTran999/auth-service/pkg"
	"github.com/google/uuid"
)

type sessionUC struct {
	sessionRepo repository.SessionRepository
	cache       pkg.Cache
}

func NewSessionUC(cache pkg.Cache, sessionRepo repository.SessionRepository) *sessionUC {
	return &sessionUC{
		sessionRepo: sessionRepo,
		cache:       cache,
	}
}

// DeleteExpiredBefore removes all sessions from the repository that expired before the given cutoff time.
// Useful for periodic cleanup of stale session data.
func (uc *sessionUC) DeleteExpiredBefore(ctx context.Context, cutoff time.Time) error {
	if err := uc.sessionRepo.DeleteExpiredBefore(ctx, cutoff); err != nil {
		return fmt.Errorf("failed to delete sessions expired before %s: %w", cutoff.Format(time.RFC3339), err)
	}
	return nil
}

func (uc *sessionUC) MarkExpiredSessions(ctx context.Context) error {
	// Fetch all active sessions from the DB
	activeSessions, err := uc.sessionRepo.FindAllActiveSession(ctx)
	if err != nil {
		return fmt.Errorf("mark expired sessions: failed to fetch active sessions: %w", err)
	}

	// No active session return intermediately
	if len(activeSessions) == 0 {
		return nil
	}

	// Identify sessions missing in cache (timed out)
	sessionTimeoutIDs, err := uc.findSessionTimeout(ctx, activeSessions)
	if err != nil {
		return fmt.Errorf("mark expired sessions: failed to find timed-out sessions: %w", err)
	}

	//  All session are active return intermediately
	if len(sessionTimeoutIDs) == 0 {
		return nil
	}

	// Update expiration in DB
	err = uc.sessionRepo.MarkSessionsExpired(ctx, sessionTimeoutIDs, time.Now())
	if err != nil {
		return fmt.Errorf("mark expired sessions: failed to update DB: %w", err)
	}

	return nil
}

func (uc *sessionUC) ValidateSession(ctx context.Context, sessionID string) (*model.Session, error) {
	if _, err := uuid.Parse(sessionID); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidSessionID, err)
	}

	session, err := uc.findSessionByID(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("validate session failed for id=%s: %w", sessionID, err)
	}

	return session, nil
}

// findSessionTimeout returns the IDs of sessions that are not found in the cache.
// These sessions are considered "timed out" and may need to be expired.
func (uc *sessionUC) findSessionTimeout(ctx context.Context, activeSessions []model.Session) ([]string, error) {
	cacheKeys := make([]string, len(activeSessions))
	for i, session := range activeSessions {
		cacheKeys[i] = common.KeyFromSessionID(session.ID.String())
	}

	missingKeys, err := uc.cache.MissingKeys(ctx, cacheKeys...)
	if err != nil {
		return nil, err
	}

	timedOutSessionIDs := make([]string, 0, len(missingKeys))
	for _, key := range missingKeys {
		timedOutSessionIDs = append(timedOutSessionIDs, common.SessionIDFromKey(key))
	}

	return timedOutSessionIDs, nil
}

func (uc *sessionUC) findSessionByID(ctx context.Context, sessionID string) (*model.Session, error) {
	var session model.Session

	// Try lookup in cache first
	if err := uc.cache.GetInto(ctx, common.KeyFromSessionID(sessionID), &session); err == nil {
		return &session, nil
	}

	// Fallback to DB
	found, err := uc.sessionRepo.FindByID(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	// Session is expired or not existed
	if found == nil || found.IsExpired() {
		return nil, ErrSessionNotFound
	}

	_ = uc.cache.Set(ctx, common.KeyFromSessionID(sessionID), found, sessionDuration)

	return found, nil
}
