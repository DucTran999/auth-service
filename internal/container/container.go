package container

import (
	"fmt"
	"log"

	"github.com/DucTran999/auth-service/config"
	gen "github.com/DucTran999/auth-service/gen/http"
	"github.com/DucTran999/auth-service/internal/handler/background"
	"github.com/DucTran999/auth-service/pkg/cache"
	"github.com/DucTran999/auth-service/pkg/hasher"
	"github.com/DucTran999/auth-service/pkg/signer"
	"github.com/DucTran999/dbkit"
	"github.com/DucTran999/shared-pkg/logger"
)

type Container struct {
	AppConfig *config.EnvConfiguration

	Logger logger.ILogger
	Hasher hasher.Hasher
	Signer signer.TokenSigner

	AuthDB dbkit.Connection
	Cache  cache.Cache

	repositories *repositories
	useCases     *useCases
	handlers     *handlers

	RestHandler           gen.ServerInterface
	CleanupSessionHandler background.SessionCleaner
}

// NewContainer initializes and wires together all core dependencies of the application,
// including logger, database, cache, repositories, usecases, and handlers.
// It returns a fully constructed container instance ready for use in the application.
func NewContainer(cfg *config.EnvConfiguration) (*Container, error) {
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

	signer, err := newSigner(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to init jwt signer: %w", err)
	}
	log.Println("[INFO] signer jwt initialized")

	// Construct the container with base-level services
	c := &Container{
		AppConfig: cfg,
		AuthDB:    conn,
		Cache:     cache,
		Logger:    logger,
		Hasher:    hasher.NewHasher(), // Utility for password hashing and similar needs
		Signer:    signer,
	}

	// Initialize layered application components in dependency order
	c.initRepositories() // Data access layer (repositories)
	c.initUseCases()     // Application business logic layer (usecases)
	c.initHandlers()     // HTTP handlers for API endpoints
	c.initRestHandler()  // Adapter for generated OpenAPI ServerInterface implementation
	c.initJobs()

	return c, nil
}

func (c *Container) Close() {
	if err := c.AuthDB.Close(); err != nil {
		c.Logger.Warnf("failed to close db connect: %v", err)
	}
	if err := c.Cache.Close(); err != nil {
		c.Logger.Warnf("failed to close cache connection: %v", err)
	}
	c.Logger.Info("db connection closed gracefully")
}
