package container

import (
	"fmt"
	"log"

	"github.com/DucTran999/auth-service/config"
	"github.com/DucTran999/auth-service/internal/gen"
	"github.com/DucTran999/auth-service/internal/handler/background"
	"github.com/DucTran999/auth-service/internal/handler/http"
	"github.com/DucTran999/auth-service/internal/repository"
	"github.com/DucTran999/auth-service/internal/usecase"
	"github.com/DucTran999/auth-service/pkg"
	"github.com/DucTran999/dbkit"
	"github.com/DucTran999/shared-pkg/logger"
)

type Container interface {
	APIHandler() gen.ServerInterface
	SessionCleaner() background.SessionCleaner
	Logger() logger.ILogger
	Close()
}

type repositories struct {
	account repository.AccountRepo
	session repository.SessionRepository
}

type useCases struct {
	auth    usecase.AuthUseCase
	account usecase.AccountUseCase
	session usecase.SessionUsecase
}

type handlers struct {
	auth    http.AuthHandler
	account http.AccountHandler
	health  http.HealthHandler
}

type jobs struct {
	SessionCleaner background.SessionCleaner
}

type container struct {
	appConfig *config.EnvConfiguration

	logger logger.ILogger
	hasher pkg.Hasher

	authDBConn dbkit.Connection
	cache      pkg.Cache

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
		hasher:     pkg.NewHasher(), // Utility for password hashing and similar needs
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

func (c *container) initRepositories() {
	c.repositories = &repositories{
		account: repository.NewAccountRepo(c.authDBConn.DB()),
		session: repository.NewSessionRepository(c.authDBConn.DB()),
	}
}

func (c *container) initUseCases() {
	accountUC := usecase.NewAccountUseCase(
		c.hasher,
		c.repositories.account,
	)

	authUC := usecase.NewAuthUseCase(
		c.hasher,
		c.cache,
		c.repositories.account,
		c.repositories.session,
	)

	sessionUC := usecase.NewSessionUC(c.cache, c.repositories.session)

	c.useCases = &useCases{
		account: accountUC,
		auth:    authUC,
		session: sessionUC,
	}
}

func (c *container) initHandlers() {
	c.handlers = &handlers{
		auth:    http.NewAuthHandler(c.logger, c.useCases.auth),
		account: http.NewAccountHandler(c.useCases.account),
		health:  http.NewHealthHandler(c.appConfig.ServiceEnv),
	}
}

type apiHandler struct {
	http.AuthHandler
	http.AccountHandler
	http.HealthHandler
}

func (c *container) initAPIHandler() {
	c.apiHandler = &apiHandler{
		AuthHandler:    c.handlers.auth,
		AccountHandler: c.handlers.account,
		HealthHandler:  c.handlers.health,
	}
}

func (c *container) initJobs() {
	c.jobs = &jobs{
		SessionCleaner: background.NewSessionCleaner(c.logger, c.useCases.session),
	}
}
