package container

import (
	"github.com/DucTran999/auth-service/config"
	"github.com/DucTran999/shared-pkg/cache"
)

func newRedisCache(config *config.EnvConfiguration) (cache.Cache, error) {
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

	return cache, nil
}
