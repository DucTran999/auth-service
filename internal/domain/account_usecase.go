package domain

import "context"

// AccountUseCase defines the business logic for managing user accounts.
type AccountUseCase interface {
	// Register creates a new user account with the provided information.
	// It typically includes validation, password hashing, and persistence logic.
	Register(ctx context.Context, input RegisterInput) (*Account, error)

	// ChangePassword change password for user when old password are match
	ChangePassword(ctx context.Context, input ChangePasswordInput) error
}
