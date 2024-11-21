package ini

import (
	"fmt"

	"github.com/DucTran999/auth-service/config"
	gormdb "github.com/DucTran999/shared-pkg/database"
	"github.com/DucTran999/shared-pkg/logger"
)

func connectDatabase(config *config.EnvConfiguration) (gormdb.IDBConnector, error) {
	dbConf := gormdb.DBConfig{
		Driver:                config.DBDriver,
		Env:                   config.ServiceEnv,
		Host:                  config.DBHost,
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

	db, err := gormdb.NewDBConnector(dbConf).Connect()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func closeDBConnection(log logger.ILogger, dbInst gormdb.IDBConnector) func() error {
	return func() error {
		log.Info("Stop db connection...")
		if err := dbInst.Stop(); err != nil {
			return fmt.Errorf("stop db connection got err: %v", err)
		}

		log.Info("Stop db connection successfully!")
		return nil
	}
}
