package service

import "errors"

var (
	ErrEmailExisted = errors.New("email already registered")
)
