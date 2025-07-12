// Centralize errors in app
package errs

import "errors"

var (
	ErrInvalidSessionID = errors.New("invalid session id")
	ErrSessionNotFound  = errors.New("session not found or expired")

	ErrEmailExisted           = errors.New("email already registered")
	ErrInvalidCredentials     = errors.New("invalid credentials")
	ErrNewPasswordMustChanged = errors.New("new password must be different")

	ErrAccountDisabled = errors.New("account is disabled")
	ErrAccountNotFound = errors.New("account not found")
)
