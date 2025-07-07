package container

import (
	"context"
	"time"

	"github.com/DucTran999/auth-service/config"
	"github.com/DucTran999/auth-service/pkg/cache"
	"github.com/DucTran999/auth-service/pkg/signer"
	"github.com/DucTran999/cachekit"
	"github.com/DucTran999/dbkit"
	dbConfig "github.com/DucTran999/dbkit/config"
	"github.com/DucTran999/shared-pkg/logger"
)

func newLogger(cfg *config.EnvConfiguration) (logger.ILogger, error) {
	return logger.NewLogger(logger.Config{
		Environment: cfg.ServiceEnv,
		LogToFile:   cfg.LogToFile,
		FilePath:    cfg.LogFilePath,
	})
}

func newAuthDBConnection(config *config.EnvConfiguration) (dbkit.Connection, error) {
	pgConf := dbConfig.PostgreSQLConfig{
		Config: dbConfig.Config{
			Host:     config.DBHost,
			Port:     config.DBPort,
			Username: config.DBUsername,
			Password: config.DBPasswd,
			Database: config.DBDatName,
			TimeZone: config.DBTimezone,
		},
		PoolConfig: dbConfig.PoolConfig{
			MaxOpenConnection: config.DBMaxOpenConnections,
			MaxIdleConnection: config.DBMaxIdleConnections,
			ConnMaxIdleTime:   time.Duration(config.DBMaxConnectionIdleTime) * time.Second,
		},
		SSLMode: dbConfig.PgSSLDisable,
	}

	conn, err := dbkit.NewPostgreSQLConnection(pgConf)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func newSigner(config *config.EnvConfiguration) (signer.TokenSigner, error) {
	cfg := signer.Config{
		Alg:     signer.SigningAlgorithm(config.SignMethod),
		PrivPem: config.PrivPem,
		PubPem:  config.PubPem,
	}

	return signer.NewTokenSigner(cfg)
}

type loggingCache struct {
	inner  cachekit.RemoteCache
	logger logger.ILogger
}

func newRedisCache(config *config.EnvConfiguration, logger logger.ILogger) (cache.Cache, error) {
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

// Get retrieves a value from the cache by its key.
func (lc *loggingCache) GetInto(ctx context.Context, key string, dest any) error {
	err := lc.inner.GetInto(ctx, key, dest)
	if err != nil {
		lc.logger.Warnf("cache get failed: key=%s err=%v", key, err)
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

// Expire sets the expiration time for a given key.
func (lc *loggingCache) Expire(ctx context.Context, key string, expiration time.Duration) error {
	err := lc.inner.Expire(ctx, key, expiration)
	if err != nil {
		lc.logger.Warnf("cache set expire failed: key=%s err=%v", key, err)
	}
	return err
}

// TTL returns the remaining time-to-live (in seconds) for a given key.
func (lc *loggingCache) TTL(ctx context.Context, key string) (int64, error) {
	ttl, err := lc.inner.TTL(ctx, key)
	if err != nil {
		lc.logger.Warnf("get TTL failed: key=%s err=%v", key, err)
	}
	return ttl, err
}

// MissingKeys returns the subset of keys that are not present in the cache.
func (lc *loggingCache) MissingKeys(ctx context.Context, keys ...string) ([]string, error) {
	missing, err := lc.inner.MissingKeys(ctx, keys...)
	if err != nil {
		lc.logger.Warnf("missing keys check failed: keys=%v err=%v", keys, err)
	}
	return missing, err
}

func (lc *loggingCache) Del(ctx context.Context, key string) error {
	err := lc.inner.Del(ctx, key)
	if err != nil {
		lc.logger.Warnf("cache: failed to delete key=%s err=%v", key, err)
	}
	return err
}

func (lc *loggingCache) Has(ctx context.Context, key string) (bool, error) {
	existed, err := lc.inner.Has(ctx, key)
	if err != nil {
		lc.logger.Warnf("cache: failed to check exist key=%s err=%v", key, err)
	}
	return existed, err
}

// Close client connection.
func (lc *loggingCache) Close() error {
	return lc.inner.Close()
}
