package container

import (
	"fmt"
	"log"

	"github.com/DucTran999/auth-service/config"
	"github.com/DucTran999/auth-service/internal/gen"
	"github.com/DucTran999/auth-service/internal/handler/background"
	"github.com/DucTran999/auth-service/pkg/cache"
	"github.com/DucTran999/auth-service/pkg/hasher"
	"github.com/DucTran999/dbkit"
	"github.com/DucTran999/shared-pkg/logger"
)

type Container interface {
	APIHandler() gen.ServerInterface
	SessionCleaner() background.SessionCleaner
	Logger() logger.ILogger
	Close()
}

type container struct {
	appConfig *config.EnvConfiguration

	logger logger.ILogger
	hasher hasher.Hasher

	authDBConn dbkit.Connection
	cache      cache.Cache

	useCases     *useCases
	repositories *repositories
	handlers     *handlers
	jobs         *jobs

	apiHandler gen.ServerInterface
}

// NewContainer initializes and wires together all core dependencies of the application,
// including logger, database, cache, repositories, usecases, and handlers.
// It returns a fully constructed container instance ready for use in the application.
func NewContainer(cfg *config.EnvConfiguration) (*container, error) {
	// Initialize application logger
	logger, err := newLogger(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %w", err)
	}
	log.Println("[INFO] initialize logger successfully")

	// Establish database connection for the authentication domain
	conn, err := newAuthDBConnection(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to connect auth db: %w", err)
	}
	log.Println("[INFO] connection db successfully")

	// Initialize Redis-based cache system
	cache, err := newRedisCache(cfg, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to connect redis cache: %w", err)
	}
	log.Println("[INFO] connection redis successfully")

	// Construct the container with base-level services
	c := &container{
		authDBConn: conn,
		cache:      cache,
		logger:     logger,
		appConfig:  cfg,
		hasher:     hasher.NewHasher(), // Utility for password hashing and similar needs
	}

	// Initialize layered application components in dependency order
	c.initRepositories() // Data access layer (repositories)
	c.initUseCases()     // Application business logic layer (usecases)
	c.initHandlers()     // HTTP handlers for API endpoints
	c.initAPIHandler()   // Adapter for generated OpenAPI ServerInterface implementation
	c.initJobs()

	return c, nil
}

func (c *container) APIHandler() gen.ServerInterface {
	return c.apiHandler
}

func (c *container) SessionCleaner() background.SessionCleaner {
	return c.jobs.SessionCleaner
}

func (c *container) Logger() logger.ILogger {
	return c.logger
}

func (c *container) Close() {
	if err := c.authDBConn.Close(); err != nil {
		c.logger.Warnf("failed to close db connect: %v", err)
	}
	if err := c.cache.Close(); err != nil {
		c.logger.Warnf("failed to close cache connection: %v", err)
	}
	c.logger.Info("db connection closed gracefully")
}
