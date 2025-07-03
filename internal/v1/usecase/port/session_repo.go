package port

import (
	"context"
	"time"

	"github.com/DucTran999/auth-service/internal/model"
)

// SessionRepository defines the methods required to manage session data in the persistence layer.
type SessionRepository interface {
	// Create stores a new session in the database.
	Create(ctx context.Context, session *model.Session) error

	// DeleteExpiredBefore permanently deletes sessions that expired before the given cutoff time.
	DeleteExpiredBefore(ctx context.Context, cutoff time.Time) error

	// FindAllActiveSession retrieves all currently active sessions.
	FindAllActiveSession(ctx context.Context) ([]model.Session, error)

	// FindByID retrieves a session by its session ID.
	// Returns nil if the session is not found.
	FindByID(ctx context.Context, sessionID string) (*model.Session, error)

	// UpdateExpiresAt updates the expiration timestamp of a session by session ID.
	UpdateExpiresAt(ctx context.Context, sessionID string, expiresAt time.Time) error

	// MarkSessionsExpired sets the expiration timestamp for multiple sessions by their IDs.
	MarkSessionsExpired(ctx context.Context, sessionIDs []string, expiresAt time.Time) error
}
