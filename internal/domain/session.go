package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// SessionRepository defines the methods required to manage session data in the persistence layer.
type SessionRepository interface {
	// Create stores a new session in the database.
	Create(ctx context.Context, session *Session) error

	// DeleteExpiredBefore permanently deletes sessions that expired before the given cutoff time.
	DeleteExpiredBefore(ctx context.Context, cutoff time.Time) error

	// FindAllActiveSession retrieves all currently active sessions.
	FindAllActiveSession(ctx context.Context) ([]Session, error)

	// FindByID retrieves a session by its session ID.
	// Returns nil if the session is not found.
	FindByID(ctx context.Context, sessionID string) (*Session, error)

	// UpdateExpiresAt updates the expiration timestamp of a session by session ID.
	UpdateExpiresAt(ctx context.Context, sessionID string, expiresAt time.Time) error

	// MarkSessionsExpired sets the expiration timestamp for multiple sessions by their IDs.
	MarkSessionsExpired(ctx context.Context, sessionIDs []string, expiresAt time.Time) error
}

type Session struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"session_id"`
	AccountID uuid.UUID `gorm:"type:uuid;not null" json:"account_id"`
	Account   Account   `gorm:"foreignKey:AccountID;constraint:OnDelete:CASCADE" json:"-"`

	IPAddress string `gorm:"type:inet" json:"ip_address"`
	UserAgent string `json:"user_agent"`

	CreatedAt time.Time  `gorm:"type:timestamptz;default:now()" json:"created_at"`
	UpdatedAt time.Time  `gorm:"type:timestamptz;default:now()" json:"updated_at"`
	ExpiresAt *time.Time `gorm:"type:timestamptz" json:"expires_at,omitempty"`
}

func (s *Session) TableName() string { return "sessions" }

// IsExpired checks whether the session has expired based on the ExpiresAt field.
// Returns false if ExpiresAt is nil (no expiration).
func (s *Session) IsExpired() bool {
	if s.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*s.ExpiresAt)
}
