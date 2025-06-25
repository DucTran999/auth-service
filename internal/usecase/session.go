package usecase

import (
	"context"
	"time"
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
}
