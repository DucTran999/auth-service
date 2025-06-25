package usecase

import "errors"

var (
	ErrEmailExisted       = errors.New("email already registered")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrAccountDisabled    = errors.New("account is disabled")
	ErrInvalidSessionID   = errors.New("invalid session id")
)
