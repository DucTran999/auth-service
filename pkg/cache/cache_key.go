package cache

import "strings"

const (
	SessionKeyPrefix = "auth:session:"
)

// KeyFromSessionID returns the Redis cache key for the given session ID.
// e.g. "auth:session:abc123"
func KeyFromSessionID(sessionID string) string {
	return SessionKeyPrefix + sessionID
}

// SessionIDFromKey extracts the session ID from a full Redis cache key.
// e.g. "auth:session:abc123" → "abc123"
func SessionIDFromKey(key string) string {
	return strings.TrimPrefix(key, SessionKeyPrefix)
}
