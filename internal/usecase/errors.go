package usecase

import "errors"

var (
	// message: invalid credentials.
	ErrInvalidCredentials = errors.New("invalid credentials")

	// message: new password must be different.
	ErrNewPasswordMustChanged = errors.New("new password must be different")
)
