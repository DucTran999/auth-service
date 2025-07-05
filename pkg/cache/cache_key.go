package cache

import (
	"strings"
)

const (
	SessionKeyPrefix      = "auth:session:"
	RefreshTokenKeyPrefix = "auth:refresh"
)

// KeyFromSessionID returns the Redis cache key for the given session ID.
// e.g. "auth:session:abc123".
func KeyFromSessionID(sessionID string) string {
	return SessionKeyPrefix + sessionID
}

// SessionIDFromKey extracts the session ID from a full Redis cache key.
// e.g. "auth:session:abc123" → "abc123".
func SessionIDFromKey(key string) string {
	return strings.TrimPrefix(key, SessionKeyPrefix)
}

// KeyRefreshToken returns the Redis key for a specific refresh token.
//
//   - format: "auth:refresh:<user_id>:<jti>"
//   - example: "auth:refresh:42:b7f3-xyz"
func KeyRefreshToken(userID, jti string) string {
	return RefreshTokenKeyPrefix + ":" + userID + ":" + jti
}
