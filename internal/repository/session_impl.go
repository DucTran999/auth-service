package repository

import (
	"context"
	"errors"
	"time"

	"github.com/DucTran999/auth-service/internal/errs"
	"github.com/DucTran999/auth-service/internal/model"
	"gorm.io/gorm"
)

// sessionRepo implements the SessionRepository interface.
type sessionRepoImpl struct {
	db *gorm.DB
}

// NewSessionRepository returns a concrete implementation of SessionRepository.
func NewSessionRepository(db *gorm.DB) *sessionRepoImpl {
	return &sessionRepoImpl{db: db}
}

// Create inserts a new session record into the database.
func (r *sessionRepoImpl) Create(ctx context.Context, session *model.Session) error {
	return r.db.WithContext(ctx).Create(session).Error
}

func (r *sessionRepoImpl) FindByID(ctx context.Context, sessionID string) (*model.Session, error) {
	var session model.Session

	err := r.db.WithContext(ctx).
		Preload("Account", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "email", "role", "is_active")
		}).
		Where("id = ?", sessionID).
		First(&session).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrSessionNotFound
		}
		return nil, err
	}

	return &session, nil
}

func (r *sessionRepoImpl) UpdateExpiresAt(
	ctx context.Context,
	sessionID string,
	expiresAt time.Time,
) error {
	return r.db.WithContext(ctx).
		Model(&model.Session{}).
		Where("id = ?", sessionID).
		Update("expires_at", expiresAt).Error
}

func (r *sessionRepoImpl) DeleteExpiredBefore(ctx context.Context, cutoff time.Time) error {
	query := `DELETE FROM sessions WHERE expires_at < ?`
	return r.db.WithContext(ctx).Exec(query, cutoff).Error
}

func (r *sessionRepoImpl) FindAllActiveSession(ctx context.Context) ([]model.Session, error) {
	var activeSessions []model.Session

	err := r.db.WithContext(ctx).
		Select("id").
		Where("expires_at IS NULL").
		Find(&activeSessions).Error

	if err != nil {
		return nil, err
	}

	return activeSessions, nil
}

// MarkSessionsExpired sets the expiration timestamp for multiple sessions by their IDs.
func (r *sessionRepoImpl) MarkSessionsExpired(
	ctx context.Context,
	sessionIDs []string,
	expiresAt time.Time,
) error {
	if len(sessionIDs) == 0 {
		return nil // nothing to do
	}

	// Bulk update: set expires_at where session id is in sessionIDs
	err := r.db.WithContext(ctx).
		Model(&model.Session{}).
		Where("id IN ?", sessionIDs).
		Update("expires_at", expiresAt).Error
	if err != nil {
		return err
	}

	return nil
}
