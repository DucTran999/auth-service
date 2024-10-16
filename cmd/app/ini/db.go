package ini

import (
	"github.com/DucTran999/auth-service/config"
	"github.com/DucTran999/auth-service/pkg/database"
	"gorm.io/gorm"
)

func connectDatabase(config *config.EnvConfiguration) (*gorm.DB, error) {
	dbConf := database.DBConfig{
		Driver:                config.DBDriver,
		Env:                   config.Environment,
		Host:                  config.Host,
		Port:                  config.DBPort,
		Username:              config.DBUsername,
		Password:              config.DBPasswd,
		Database:              config.DBDatName,
		SslMode:               config.DBSslMode,
		Timezone:              config.DBTimezone,
		MaxOpenConnections:    config.DBMaxOpenConnections,
		MaxIdleConnections:    config.DBMaxIdleConnections,
		MaxConnectionIdleTime: config.DBMaxConnectionIdleTime,
	}

	db, err := database.NewDBConnector(dbConf).Connect()
	if err != nil {
		return nil, err
	}

	return db, nil
}
