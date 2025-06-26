package domain

import "errors"

var (
	ErrEmailExisted    = errors.New("email already registered")
	ErrAccountDisabled = errors.New("account is disabled")

	ErrInvalidSessionID = errors.New("invalid session id")
	ErrSessionNotFound  = errors.New("session not found or expired")
)
