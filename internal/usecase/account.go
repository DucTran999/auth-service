package usecase

import (
	"context"

	"github.com/DucTran999/auth-service/internal/model"
)

// AccountUseCase defines the business logic for managing user accounts.
type AccountUseCase interface {
	// Register creates a new user account with the provided information.
	// It typically includes validation, password hashing, and persistence logic.
	Register(ctx context.Context, input RegisterInput) (*model.Account, error)
}

type RegisterInput struct {
	Email    string
	Password string
}
