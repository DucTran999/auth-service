package container

import (
	"fmt"
	"log"

	"github.com/DucTran999/auth-service/config"
	"github.com/DucTran999/dbkit"
	"github.com/DucTran999/shared-pkg/logger"
	"gorm.io/gorm"
)

type Container interface {
	AuthDB() *gorm.DB
	Logger() logger.ILogger

	Close()
}

type container struct {
	authDBConn dbkit.Connection
	logger     logger.ILogger
}

func NewContainer(cfg *config.EnvConfiguration) (*container, error) {
	logger, err := newLogger(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %w", err)
	}
	log.Println("[INFO] initialize logger successfully")

	// Connection database
	conn, err := newAuthDBConnection(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to connect auth db: %w", err)
	}
	log.Println("[INFO] connection db successfully")

	// Create new dependencies container instance
	return &container{
		authDBConn: conn,
		logger:     logger,
	}, nil
}

func (c *container) AuthDB() *gorm.DB {
	return c.authDBConn.DB()
}

func (c *container) Logger() logger.ILogger {
	return c.logger
}

func (c *container) Close() {
	if err := c.authDBConn.Close(); err != nil {
		c.logger.Warnf("failed to close db connect: %v", err)
	}
	c.logger.Info("db connection closed gracefully")
}
