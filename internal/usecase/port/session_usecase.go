package port

import (
	"context"

	"github.com/DucTran999/auth-service/internal/model"
)

// SessionUsecase defines business logic for validating sessions in user-facing flows.
type SessionUsecase interface {
	// Validate checks if a session exists and is not expired.
	// It first looks in the cache, and if missing, checks persistent storage.
	Validate(ctx context.Context, sessionID string) (*model.Session, error)
}
