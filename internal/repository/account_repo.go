package repository

import (
	"context"

	"github.com/DucTran999/auth-service/internal/model"
)

// AccountRepo defines the data access methods for managing accounts in the persistence layer.
type AccountRepo interface {
	// FindByEmail retrieves an account by its unique email address.
	// Returns ErrAccountNotFound if no account is found.
	FindByEmail(ctx context.Context, email string) (*model.Account, error)

	// FindByID retrieves an account by its unique account ID.
	// Returns ErrAccountNotFound if no account is found.
	FindByID(ctx context.Context, accountID string) (*model.Account, error)

	// Create inserts a new account record into the underlying data store.
	// Returns the created account with its generated ID.
	Create(ctx context.Context, account model.Account) (*model.Account, error)

	// UpdatePasswordHash updates the password hash of the given account.
	// It does not validate the old password — that should be handled by the use case layer.
	UpdatePasswordHash(ctx context.Context, accountID string, passwordHash string) error
}
