package signer

import "errors"

var (
	ErrMissingKey       = errors.New("jwtkit: missing key")
	ErrInvalidAlgorithm = errors.New("jwtkit: unsupported signing method")
)
