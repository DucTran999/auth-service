package model

import (
	"time"

	"github.com/google/uuid"
)

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

func (Session) TableName() string { return "sessions" }
