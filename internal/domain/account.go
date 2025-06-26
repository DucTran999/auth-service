package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// AccountRepo defines the data access methods for managing accounts in the persistence layer.
type AccountRepo interface {
	// FindByEmail retrieves an account by its unique email address.
	// Returns ErrAccountNotFound if no account is found.
	FindByEmail(ctx context.Context, email string) (*Account, error)

	// FindByID retrieves an account by its unique account ID.
	// Returns ErrAccountNotFound if no account is found.
	FindByID(ctx context.Context, accountID string) (*Account, error)

	// Create inserts a new account record into the underlying data store.
	// Returns the created account with its generated ID.
	Create(ctx context.Context, account Account) (*Account, error)

	// UpdatePasswordHash updates the password hash of the given account.
	// It does not validate the old password — that should be handled by the use case layer.
	UpdatePasswordHash(ctx context.Context, accountID string, passwordHash string) error
}

// Account represents a user account in the system.
type Account struct {
	ID uuid.UUID `json:"id" gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()"`

	Email        string `json:"email" gorm:"column:email;type:text;unique;not null"`
	PasswordHash string `json:"-" gorm:"column:password_hash;type:text;not null"`

	IsVerified bool   `json:"is_verified" gorm:"column:is_verified;type:boolean;default:false"`
	IsActive   bool   `json:"is_active" gorm:"column:is_active;type:boolean;default:true"`
	Role       string `json:"role" gorm:"column:role;type:varchar(255);default:'user'"`

	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;type:timestamptz;default:now()"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at;type:timestamptz;default:now()"`
}

func (a *Account) TableName() string { return "accounts" }
