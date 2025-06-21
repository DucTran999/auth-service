package usecase

import "errors"

var (
	ErrEmailExisted = errors.New("email already registered")
)
