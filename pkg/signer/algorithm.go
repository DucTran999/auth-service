package signer

import (
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type SigningAlgorithm string

const (
	HS256 SigningAlgorithm = "HS256"
	HS384 SigningAlgorithm = "HS384"
	HS512 SigningAlgorithm = "HS512"

	RS256 SigningAlgorithm = "RS256"
	RS384 SigningAlgorithm = "RS384"
	RS512 SigningAlgorithm = "RS512"

	ES256 SigningAlgorithm = "ES256"
	ES384 SigningAlgorithm = "ES384"
	ES512 SigningAlgorithm = "ES512"

	EdDSA SigningAlgorithm = "EdDSA"
)

var jwtMethods = map[SigningAlgorithm]jwt.SigningMethod{
	HS256: jwt.SigningMethodHS256,
	HS384: jwt.SigningMethodHS384,
	HS512: jwt.SigningMethodHS512,
	RS256: jwt.SigningMethodRS256,
	RS384: jwt.SigningMethodRS384,
	RS512: jwt.SigningMethodRS512,
	ES256: jwt.SigningMethodES256,
	ES384: jwt.SigningMethodES384,
	ES512: jwt.SigningMethodES512,
	EdDSA: jwt.SigningMethodEdDSA,
}

func (a SigningAlgorithm) ToJWTMethod() (jwt.SigningMethod, error) {
	cleaned := SigningAlgorithm(strings.ToUpper(strings.TrimSpace(string(a))))
	if method, ok := jwtMethods[cleaned]; ok {
		return method, nil
	}
	return nil, fmt.Errorf("%w: %s", ErrInvalidAlgorithm, cleaned)
}
