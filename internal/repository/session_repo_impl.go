package repository

import (
	"context"

	"github.com/DucTran999/auth-service/internal/model"
	"gorm.io/gorm"
)

// sessionRepo implements the SessionRepository interface.
type sessionRepoImpl struct {
	db *gorm.DB
}

// NewSessionRepository returns a concrete implementation of SessionRepository.
func NewSessionRepository(db *gorm.DB) SessionRepository {
	return &sessionRepoImpl{db: db}
}

// Create inserts a new session record into the database.
func (r *sessionRepoImpl) Create(ctx context.Context, session *model.Session) error {
	return r.db.WithContext(ctx).Create(session).Error
}
