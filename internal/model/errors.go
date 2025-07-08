package model

import "errors"

var (
	ErrEmailExisted    = errors.New("email already registered")
	ErrAccountDisabled = errors.New("account is disabled")

	ErrInvalidSessionID = errors.New("invalid session id")
	ErrSessionNotFound  = errors.New("session not found or expired")

	// message: invalid credentials.
	ErrInvalidCredentials = errors.New("invalid credentials")

	// message: new password must be different.
	ErrNewPasswordMustChanged = errors.New("new password must be different")
)
