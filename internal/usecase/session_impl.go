package usecase

import (
	"context"
	"fmt"
	"time"

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

func (uc *sessionUC) SetExpirationIfNotCached(ctx context.Context) error {
	// Fetch all active sessions from the DB
	activeSessions, err := uc.sessionRepo.FindAllActiveSession(ctx)
	if err != nil {
		return fmt.Errorf("set expiration: failed to fetch active sessions: %w", err)
	}

	// For each session, check if it exists in cache
	for _, s := range activeSessions {
		cacheKey := uc.getCacheKey(s.ID.String())

		// Check if session is already cached
		exists, err := uc.cache.TTL(ctx, cacheKey)
		if err != nil {
			// uc.logger.Warnf("cache check failed for session %s: %v", s.ID, err)
			continue // Skip this session, try others
		}

		if exists != 0 {
			continue // Already cached, skip
		}
	}

	return nil
}

func (uc *sessionUC) getCacheKey(sessionID string) string {
	return "session-" + sessionID
}
