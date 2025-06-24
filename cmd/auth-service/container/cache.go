package container

import (
	"context"
	"time"

	"github.com/DucTran999/auth-service/config"
	"github.com/DucTran999/auth-service/pkg"
	"github.com/DucTran999/cachekit"
	"github.com/DucTran999/shared-pkg/logger"
)

type loggingCache struct {
	inner  pkg.Cache
	logger logger.ILogger
}

// Get retrieves a value from the cache by its key.
func (lc *loggingCache) GetInto(ctx context.Context, key string, dest any) error {
	err := lc.inner.GetInto(ctx, key, dest)
	if err != nil {
		lc.logger.Warnf("cache get failed: key=%s err=%v", key, err)
		return err
	}

	return err
}

// Set stores a value in the cache with an optional expiration time.
func (lc *loggingCache) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	err := lc.inner.Set(ctx, key, value, expiration)
	if err != nil {
		lc.logger.Warnf("cache set failed: key=%s err=%v", key, err)
	}
	return err
}

// Set Expire time for key
func (lc *loggingCache) Expire(ctx context.Context, key string, expiration time.Duration) error {
	err := lc.inner.Expire(ctx, key, expiration)
	if err != nil {
		lc.logger.Warnf("cache set failed: key=%s err=%v", key, err)
	}
	return err
}

func (lc *loggingCache) TTL(ctx context.Context, key string) (int64, error) {
	ttl, err := lc.inner.TTL(ctx, key)
	if err != nil {
		lc.logger.Warnf("cache set failed: key=%s err=%v", key, err)
	}
	return ttl, err
}

// Close client connection
func (lc *loggingCache) Close() error {
	return lc.inner.Close()
}

func newRedisCache(config *config.EnvConfiguration, logger logger.ILogger) (pkg.Cache, error) {
	cacheConf := cachekit.RedisConfig{
		Host:     config.RedisHost,
		Port:     config.RedisPort,
		Password: config.RedisPasswd,
	}

	cache, err := cachekit.NewRedisCache(cacheConf)
	if err != nil {
		return nil, err
	}

	return &loggingCache{
		inner:  cache,
		logger: logger,
	}, nil
}
