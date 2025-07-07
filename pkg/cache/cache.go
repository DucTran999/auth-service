package cache

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

	// filter list keys are missing
	MissingKeys(ctx context.Context, keys ...string) ([]string, error)

	// Delete key
	Del(ctx context.Context, key string) error

	// Check key existed
	Has(ctx context.Context, key string) (bool, error)

	// Close client connection
	Close() error
}
