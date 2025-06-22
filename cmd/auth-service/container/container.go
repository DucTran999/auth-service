package container

import (
	"fmt"
	"log"

	"github.com/DucTran999/auth-service/config"
	"github.com/DucTran999/auth-service/internal/gen"
	"github.com/DucTran999/auth-service/internal/handler"
	"github.com/DucTran999/auth-service/internal/repository"
	"github.com/DucTran999/auth-service/internal/usecase"
	"github.com/DucTran999/auth-service/pkg"
	"github.com/DucTran999/dbkit"
	"github.com/DucTran999/shared-pkg/cache"
	"github.com/DucTran999/shared-pkg/logger"
	"gorm.io/gorm"
)

type Container interface {
	AppConfig() *config.EnvConfiguration

	AuthDB() *gorm.DB
	Logger() logger.ILogger
	Close()

	AppHandler() gen.ServerInterface
}
type appHandler struct {
	handler.HealthHandler
	handler.AccountHandler
	handler.AuthHandler
}

type container struct {
	logger    logger.ILogger
	appConfig *config.EnvConfiguration

	authDBConn dbkit.Connection
	cache      cache.Cache

	appHandler *appHandler
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

	cache, err := newRedisCache(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to connect redis cache: %w", err)
	}
	log.Println("[INFO] connection redis successfully")

	c := &container{
		authDBConn: conn,
		cache:      cache,
		logger:     logger,
		appConfig:  cfg,
	}
	c.initAppHandler()

	// Create new dependencies container instance
	return c, nil
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
	if err := c.cache.Close(); err != nil {
		c.logger.Warnf("failed to close cache connection: %v", err)
	}
	c.logger.Info("db connection closed gracefully")
}

func (c *container) AppHandler() gen.ServerInterface {
	return c.appHandler
}

func (c *container) AppConfig() *config.EnvConfiguration {
	return c.appConfig
}

func (c *container) initAppHandler() {
	hasher := pkg.NewHasher()

	// Account module
	accountRepo := repository.NewAccountRepo(c.authDBConn.DB())
	accountUC := usecase.NewAccountUseCase(hasher, accountRepo)

	// Auth module
	sessionRepo := repository.NewSessionRepository(c.authDBConn.DB())
	authUC := usecase.NewAuthUseCase(hasher, accountRepo, sessionRepo)

	c.appHandler = &appHandler{
		HealthHandler:  handler.NewHealthHandler(c.appConfig.ServiceVersion),
		AccountHandler: handler.NewAccountHandler(accountUC),
		AuthHandler:    handler.NewAuthHandler(authUC),
	}
}
