package pkg

import (
	"context"
	"time"
)

type Cache interface {
	GetInto(ctx context.Context, key string, dest any) error

	// Set stores a value in the cache with an optional expiration time.
	Set(ctx context.Context, key string, value any, expiration time.Duration) error

	// Set key expire time
	Expire(ctx context.Context, key string, expiration time.Duration) error

	// Check TTL
	TTL(ctx context.Context, key string) (int64, error)

	// Close client connection
	Close() error
}
