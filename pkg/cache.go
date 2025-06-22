package pkg

import (
	"context"
	"time"
)

type Cache interface {
	// Get retrieves a value from the cache by its key.
	Get(ctx context.Context, key string) (string, error)

	// Set stores a value in the cache with an optional expiration time.
	Set(ctx context.Context, key string, value any, expiration time.Duration) error

	// Close client connection
	Close() error
}
