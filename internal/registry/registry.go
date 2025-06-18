package registry

import (
	"github.com/DucTran999/auth-service/config"
	"gorm.io/gorm"
)

type Registry struct {
	AppConfig  *config.EnvConfiguration
	PostgresDB *gorm.DB
}

func NewRegistry(appConf *config.EnvConfiguration, pg *gorm.DB) *Registry {
	return &Registry{
		AppConfig:  appConf,
		PostgresDB: pg,
	}
}
