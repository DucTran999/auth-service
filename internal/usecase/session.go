package usecase

import (
	"context"
	"time"

	"github.com/DucTran999/auth-service/internal/model"
)

// SessionUsecase defines business logic operations related to session lifecycle management.
type SessionUsecase interface {
	// DeleteExpiredBefore removes all session records from storage
	// that have an expiration time earlier than the given cutoff timestamp.
	// Used typically during background purging of old sessions.
	DeleteExpiredBefore(ctx context.Context, cutoff time.Time) error

	// MarkExpiredSessions applies an expiration timestamp to sessions
	// that are not currently tracked in the cache (e.g., Redis).
	// This ensures untracked sessions do not remain active indefinitely.
	MarkExpiredSessions(ctx context.Context) error

	// ValidateSession find session in cache first if not try to lookup in DB.
	// Return session only if it is existed and not expire
	ValidateSession(ctx context.Context, sessionID string) (*model.Session, error)
}
