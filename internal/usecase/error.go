package usecase

import "errors"

var (
	ErrEmailExisted       = errors.New("email already registered")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrAccountDisabled    = errors.New("account is disabled")

	ErrNewPasswordMustChanged = errors.New("new password must be different")

	ErrInvalidSessionID = errors.New("invalid session id")
	ErrSessionNotFound  = errors.New("session not found or expired")
)
