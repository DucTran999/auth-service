package port

import (
	"context"

	"github.com/DucTran999/auth-service/internal/model"
)

// SessionUsecase defines business logic operations related to session lifecycle management.
type SessionUsecase interface {
	// ValidateSession find session in cache first if not try to lookup in DB.
	// Return session only if it is existed and not expire
	ValidateSession(ctx context.Context, sessionID string) (*model.Session, error)
}
