package app

import (
	"context"
	"os/signal"
	"syscall"
	"time"

	"github.com/DucTran999/auth-service/cmd/auth-service/container"
	"github.com/DucTran999/auth-service/cmd/auth-service/server"
	"github.com/DucTran999/auth-service/config"
	pkgServer "github.com/DucTran999/shared-pkg/server"
)

type App struct {
	appConf    *config.EnvConfiguration
	deps       container.Container
	httpServer pkgServer.HttpServer
}

// InitApp initializes the application, setting up logging, database connection, HTTP server, and graceful shutdown handling.
func NewApp(appConf *config.EnvConfiguration) (*App, error) {
	c, err := container.NewContainer(appConf)
	if err != nil {
		return nil, err
	}

	httpServer, err := server.NewHTTPServer(appConf, c)
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

func (a *App) Run() {
	appCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	defer a.deps.Close()

	go func() {
		if err := a.httpServer.Start(); err != nil {
			a.deps.Logger().Fatalf("start http server got err: %v", err)
		}
	}()

	<-appCtx.Done()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := a.httpServer.Stop(ctx); err != nil {
		a.deps.Logger().Errorf("failed to stop http server: %v", err)
	}
	a.deps.Logger().Info("http server stopped successfully")
}
