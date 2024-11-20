package ini

import (
	"fmt"
	"log"

	"github.com/DucTran999/auth-service/config"
	"github.com/DucTran999/shared-pkg/v2/database"
)

func connectDatabase(config *config.EnvConfiguration) (database.IDBConnector, error) {
	dbConf := database.DBConfig{
		Driver:                config.DBDriver,
		Env:                   config.ServiceEnv,
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

func closeDBConnection(dbInst database.IDBConnector) func() error {
	return func() error {
		log.Println("Stop db connection...")
		if err := dbInst.Stop(); err != nil {
			return fmt.Errorf("stop db connection got err: %v", err)
		}

		log.Println("Stop db connection successfully!")
		return nil
	}
}
