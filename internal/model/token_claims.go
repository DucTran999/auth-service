package model

import (
	"encoding/json"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type TokenClaims struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
	Role  string    `json:"role"`

	IssuedAt  int64 `json:"iat"`
	ExpiresAt int64 `json:"exp"`

	// only refresh token have
	JTI string `json:"jti,omitempty"`
}

func (c TokenClaims) ToMapClaims() jwt.MapClaims {
	claims := jwt.MapClaims{
		"id":    c.ID,
		"email": c.Email,
		"role":  c.Role,
		"iat":   c.IssuedAt,
		"exp":   c.ExpiresAt,
	}
	if c.JTI != "" {
		claims["jti"] = c.JTI
	}
	return claims
}

// MapClaimsToTokenClaims converts jwt.MapClaims into your TokenClaims struct.
func MapClaimsToTokenClaims(mc jwt.MapClaims) (*TokenClaims, error) {
	data, err := json.Marshal(mc)
	if err != nil {
		return nil, err
	}

	var claims TokenClaims
	if err := json.Unmarshal(data, &claims); err != nil {
		return nil, err
	}

	return &claims, nil
}
