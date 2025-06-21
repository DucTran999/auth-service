package repository

import (
	"context"
	"time"

	"github.com/DucTran999/auth-service/internal/model"
)

// SessionRepository defines the methods required to manage session data in the persistence layer.
type SessionRepository interface {
	// Create stores a new session in the database.
	Create(ctx context.Context, session *model.Session) error

	// FindByID retrieves a session by its session ID.
	// Returns nil if the session is not found.
	FindByID(ctx context.Context, sessionID string) (*model.Session, error)

	// UpdateExpiresAt updates the expiration timestamp of a session.
	UpdateExpiresAt(ctx context.Context, sessionID string, expiresAt time.Time) error
}
