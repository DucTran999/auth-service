package container

import (
	"context"
	"time"

	"github.com/DucTran999/auth-service/config"
	"github.com/DucTran999/auth-service/pkg"
	"github.com/DucTran999/shared-pkg/cache"
	"github.com/DucTran999/shared-pkg/logger"
)

type loggingCache struct {
	inner  cache.Cache
	logger logger.ILogger
}

// Get retrieves a value from the cache by its key.
func (lc *loggingCache) Get(ctx context.Context, key string) (string, error) {
	val, err := lc.inner.Get(ctx, key)
	if err != nil {
		lc.logger.Warnf("cache get failed: key=%s err=%v", key, err)
	}
	return val, err
}

// Set stores a value in the cache with an optional expiration time.
func (lc *loggingCache) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	err := lc.inner.Set(ctx, key, value, expiration)
	if err != nil {
		lc.logger.Warnf("cache set failed: key=%s err=%v", key, err)
	}
	return err
}

// Close client connection
func (lc *loggingCache) Close() error {
	return lc.inner.Close()
}

func newRedisCache(config *config.EnvConfiguration, logger logger.ILogger) (pkg.Cache, error) {
	cacheConf := cache.Config{
		IsCacheOnMemory: config.CacheInMem,
		Host:            config.RedisHost,
		Port:            config.RedisPort,
		Password:        config.RedisPasswd,
	}

	cache, err := cache.NewCache(cacheConf)
	if err != nil {
		return nil, err
	}

	return &loggingCache{
		inner:  cache,
		logger: logger,
	}, nil
}
