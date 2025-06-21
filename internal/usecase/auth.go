package usecase

import (
	"context"

	"github.com/DucTran999/auth-service/internal/model"
)

// AuthUseCase defines the authentication-related business logic.
type AuthUseCase interface {
	// Login verifies the provided credentials and returns the authenticated account.
	// Returns an error if authentication fails.
	Login(ctx context.Context, input LoginInput) (*model.Account, error)
}

// LoginInput represents the input required to authenticate a user using email and password.
// It also includes optional request metadata for logging, auditing, or session management.
type LoginInput struct {
	Email     string `json:"email"`    // User's email address
	Password  string `json:"password"` // Plain-text password from the login form
	IP        string `json:"-"`        // Client IP address (injected by handler, not from JSON)
	UserAgent string `json:"-"`        // User-Agent header string (injected by handler)
}
