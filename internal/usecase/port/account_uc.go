package port

import (
	"context"

	"github.com/DucTran999/auth-service/internal/domain"
	"github.com/DucTran999/auth-service/internal/usecase/dto"
)

// AccountUseCase defines the business logic for managing user accounts.
type AccountUseCase interface {
	// Register creates a new user account with the provided information.
	// It typically includes validation, password hashing, and persistence logic.
	Register(ctx context.Context, input dto.RegisterInput) (*domain.Account, error)

	// ChangePassword change password for user when old password are match
	ChangePassword(ctx context.Context, input dto.ChangePasswordInput) error
}
