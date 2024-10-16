package database

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type postgresConnector struct {
	config DBConfig
	db     *gorm.DB
}

func newPostgresConnector(conf DBConfig) *postgresConnector {
	return &postgresConnector{config: conf}
}

func (c *postgresConnector) Connect() (*gorm.DB, error) {
	var err error

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
		c.config.Host, c.config.Username, c.config.Password, c.config.Database,
		c.config.Port, c.config.SslMode, c.config.Timezone,
	)

	c.db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err = c.configPool(); err != nil {
		return nil, err
	}

	// Ping to db to verify the connection
	if err = c.Ping(); err != nil {
		return nil, err
	}

	return c.db, nil
}

func (c *postgresConnector) configPool() error {
	sqlDB, err := c.db.DB()
	if err != nil {
		return err
	}

	if c.config.MaxOpenConnections > 0 {
		sqlDB.SetMaxOpenConns(c.config.MaxOpenConnections)
	}

	if c.config.MaxIdleConnections > 0 {
		sqlDB.SetMaxIdleConns(c.config.MaxIdleConnections)
	}

	if c.config.MaxConnectionIdleTime > 0 {
		sqlDB.SetConnMaxIdleTime(time.Duration(c.config.MaxConnectionIdleTime) * time.Minute)
	}

	return nil
}

func (c *postgresConnector) Ping() error {
	sqlDB, err := c.db.DB()
	if err != nil {
		return err
	}

	if err := sqlDB.Ping(); err != nil {
		return err
	}

	return nil
}

func (c *postgresConnector) Stop() error {
	if c.db == nil {
		return nil
	}

	sqlDB, err := c.db.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}
