package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/DucTran999/auth-service/cmd/api_v1/container"
	"github.com/DucTran999/auth-service/config"
	"github.com/DucTran999/auth-service/internal/v1/server"
	"github.com/DucTran999/auth-service/internal/v1/worker"
	pkgServer "github.com/DucTran999/shared-pkg/server"
)

type App struct {
	appConf              *config.EnvConfiguration
	deps                 container.Container
	httpServer           pkgServer.HttpServer
	sessionCleanupWorker worker.SessionCleanupWorker
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

	scWorker := worker.NewSessionCleanupWorker(c.Logger(), appConf, c.SessionCleaner())

	app := &App{
		deps:                 c,
		appConf:              appConf,
		httpServer:           httpServer,
		sessionCleanupWorker: scWorker,
	}

	return app, nil
}

func (a *App) Run() error {
	appCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	defer a.deps.Close()

	// Track worker exit
	workerDone := make(chan struct{})

	// Start HTTP server in a goroutine
	go func() {
		if err := a.httpServer.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.deps.Logger().Fatalf("start http server got err: %v", err)
		}
	}()

	// Start session cleanup worker
	go func() {
		defer close(workerDone)

		if err := a.sessionCleanupWorker.Start(appCtx); err != nil && !errors.Is(err, context.Canceled) {
			a.deps.Logger().Errorf("session cleanup worker exited with error: %v", err)
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

	shutdownStart := time.Now()
	select {
	case <-workerDone:
		a.deps.Logger().Infof("session cleanup worker stopped in %s", time.Since(shutdownStart))
	case <-ctx.Done():
		a.deps.Logger().Warnf("timeout waiting for session cleanup worker (after %s)", time.Since(shutdownStart))
		return ctx.Err()
	}

	a.deps.Logger().Info("http server stopped successfully")
	return nil
}
