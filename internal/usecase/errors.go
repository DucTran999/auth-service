package usecase

import "errors"

var (
	ErrInvalidCredentials     = errors.New("invalid credentials")
	ErrNewPasswordMustChanged = errors.New("new password must be different")
)
