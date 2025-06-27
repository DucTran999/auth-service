package domain

import "context"

// AuthUseCase defines the authentication-related business logic.
type AuthUseCase interface {
	// Login verifies the provided credentials and returns the authenticated account.
	// Returns an error if authentication fails.
	Login(ctx context.Context, input LoginInput) (*Session, error)

	// Logout terminates the session associated with the given session ID.
	// It removes the session from cache (best-effort) and marks it as expired in the database.
	// Returns an error only if the database update fails.
	Logout(ctx context.Context, sessionID string) error
}
