package port

import (
	"context"

	"github.com/DucTran999/auth-service/internal/model"
	"github.com/DucTran999/auth-service/internal/v1/usecase/dto"
)

// AccountUseCase defines the business logic for managing user accounts.
type AccountUseCase interface {
	// Register creates a new user account with the provided information.
	// It typically includes validation, password hashing, and persistence logic.
	Register(ctx context.Context, input dto.RegisterInput) (*model.Account, error)

	// ChangePassword change password for user when old password are match
	ChangePassword(ctx context.Context, input dto.ChangePasswordInput) error
}
