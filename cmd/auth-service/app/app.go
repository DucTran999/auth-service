package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/DucTran999/auth-service/cmd/auth-service/container"
	"github.com/DucTran999/auth-service/config"
	"github.com/DucTran999/auth-service/internal/server"
	pkgServer "github.com/DucTran999/shared-pkg/server"
)

type App struct {
	appConf    *config.EnvConfiguration
	deps       container.Container
	httpServer pkgServer.HttpServer
}

// InitApp initializes the application, setting up logging, database connection, HTTP server, etc.
func NewApp(appConf *config.EnvConfiguration) (*App, error) {
	c, err := container.NewContainer(appConf)
	if err != nil {
		return nil, err
	}

	httpServer, err := server.NewHTTPServer(appConf, c.APIHandler())
	if err != nil {
		return nil, err
	}

	app := &App{
		deps:       c,
		appConf:    appConf,
		httpServer: httpServer,
	}

	return app, nil
}

func (a *App) Run() error {
	appCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	defer a.deps.Close()

	// Start HTTP server in a goroutine
	go func() {
		if err := a.httpServer.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.deps.Logger().Fatalf("start http server got err: %v", err)
		}
	}()

	// Wait for termination signal
	<-appCtx.Done()

	// Shutdown with timeout
	shutdownTime := time.Duration(a.appConf.ShutdownTime) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTime)
	defer cancel()

	if err := a.httpServer.Stop(ctx); err != nil {
		a.deps.Logger().Errorf("failed to stop http server: %v", err)
		return fmt.Errorf("graceful shutdown failed: %w", err)
	}

	a.deps.Logger().Info("http server stopped successfully")
	return nil
}
