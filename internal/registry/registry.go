package registry

import (
	"gorm.io/gorm"
)

type Registry struct {
	PostgresDB *gorm.DB
}

func NewRegistry(pg *gorm.DB) *Registry {
	return &Registry{
		PostgresDB: pg,
	}
}
