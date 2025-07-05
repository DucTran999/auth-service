package signer

import (
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt"
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

func (a SigningAlgorithm) ToJWTMethod() (jwt.SigningMethod, error) {
	cleaned := strings.ToUpper(strings.TrimSpace(string(a)))
	switch SigningAlgorithm(cleaned) {
	case HS256:
		return jwt.SigningMethodHS256, nil
	case HS384:
		return jwt.SigningMethodHS384, nil
	case HS512:
		return jwt.SigningMethodHS512, nil
	case RS256:
		return jwt.SigningMethodRS256, nil
	case RS384:
		return jwt.SigningMethodRS384, nil
	case RS512:
		return jwt.SigningMethodRS512, nil
	case ES256:
		return jwt.SigningMethodES256, nil
	case ES384:
		return jwt.SigningMethodES384, nil
	case ES512:
		return jwt.SigningMethodES512, nil
	case EdDSA:
		return jwt.SigningMethodEdDSA, nil
	default:
		return nil, fmt.Errorf("%w: %s", ErrInvalidAlgorithm, cleaned)
	}
}
