package model

import (
	"time"

	"github.com/google/uuid"
)

// Account represents a user account in the system.
type Account struct {
	ID uuid.UUID `json:"id" gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()"`

	Email    string `json:"email" gorm:"column:email;type:text;unique;not null"`
	Password string `json:"password_hash" gorm:"column:password_hash;type:text;not null"`

	IsVerified bool   `json:"is_verified" gorm:"column:is_verified;type:boolean;default:false"`
	IsActive   bool   `json:"is_active" gorm:"column:is_active;type:boolean;default:true"`
	Role       string `json:"role" gorm:"column:role;type:varchar(255);default:'user'"`

	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;type:timestamptz;default:now()"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at;type:timestamptz;default:now()"`
}

func (a *Account) TableName() string { return "accounts" }
