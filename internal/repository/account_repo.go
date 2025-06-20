package repository

import (
	"context"

	"github.com/DucTran999/auth-service/internal/model"
)

// AccountRepo defines the data access methods for managing accounts in the persistence layer.
type AccountRepo interface {
	// FindByEmail retrieves an account by email.
	FindByEmail(ctx context.Context, email string) (*model.Account, error)

	// Create inserts a new account into the database.
	Create(ctx context.Context, account model.Account) (*model.Account, error)
}
