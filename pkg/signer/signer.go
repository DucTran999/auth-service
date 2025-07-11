package signer

import "github.com/golang-jwt/jwt/v5"

type TokenSigner interface {
	// Sign generates a signed JWT string from the given claims.
	// The claims must implement the jwt.Claims interface (e.g., jwt.MapClaims, custom claims structs).
	// Returns the signed JWT as a string.
	Sign(claims jwt.Claims) (string, error)

	// ParseInto parses the JWT string into the provided destination claims struct.
	// The destination must implement jwt.Claims (e.g., a custom claims struct with embedded jwt.RegisteredClaims).
	// Useful when working with strongly-typed custom claims.
	ParseInto(tokenStr string, dest jwt.Claims) error
}
