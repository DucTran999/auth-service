package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/DucTran999/auth-service/internal/common"
	"github.com/DucTran999/auth-service/internal/model"
	"github.com/DucTran999/auth-service/internal/repository"
	"github.com/DucTran999/auth-service/pkg"
)

type sessionUC struct {
	sessionRepo repository.SessionRepository
	cache       pkg.Cache
}

func NewSessionUC(sessionRepo repository.SessionRepository) *sessionUC {
	return &sessionUC{
		sessionRepo: sessionRepo,
	}
}

func (uc *sessionUC) DeleteExpiredBefore(ctx context.Context, cutoff time.Time) error {
	return uc.sessionRepo.DeleteExpiredBefore(ctx, cutoff)
}

func (uc *sessionUC) MarkExpiredSessions(ctx context.Context) error {
	// Fetch all active sessions from the DB
	activeSessions, err := uc.sessionRepo.FindAllActiveSession(ctx)
	if err != nil {
		return fmt.Errorf("mark expired sessions: failed to fetch active sessions: %w", err)
	}

	// Identify sessions missing in cache (timed out)
	sessionTimeoutIDs, err := uc.findSessionTimeout(ctx, activeSessions)
	if err != nil {
		return fmt.Errorf("mark expired sessions: failed to find timed-out sessions: %w", err)
	}

	// Update expiration in DB
	err = uc.sessionRepo.MarkSessionsExpired(ctx, sessionTimeoutIDs, time.Now())
	if err != nil {
		return fmt.Errorf("mark expired sessions: failed to update DB: %w", err)
	}

	return nil
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
		return nil, fmt.Errorf("findSessionTimeout: failed to get missing keys from cache: %w", err)
	}

	timedOutSessionIDs := make([]string, 0, len(missingKeys))
	for _, key := range missingKeys {
		timedOutSessionIDs = append(timedOutSessionIDs, common.SessionIDFromKey(key))
	}

	return timedOutSessionIDs, nil
}
