package model

import (
	"github.com/golang-jwt/jwt/v5"
)

type TokenClaims struct {
	jwt.RegisteredClaims

	// Custom fields
	Email string `json:"email"`
	Role  string `json:"role"`
}
