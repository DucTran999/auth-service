package container

import (
	"time"

	"github.com/DucTran999/auth-service/internal/config"
	"github.com/DucTran999/dbkit"
	dbConfig "github.com/DucTran999/dbkit/config"
)

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
