package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/DucTran999/auth-service/internal/container"
	httpServer "github.com/DucTran999/auth-service/internal/server/http"
	"github.com/DucTran999/shared-pkg/server"
)

func startHTTPServer(c *container.Container) (server.HttpServer, error) {
	srv, err := httpServer.NewHTTPServer(c.AppConfig, c.RestHandler)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize HTTP server: %w", err)
	}

	go func() {
		if err := srv.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			c.Logger.Fatalf("[FATAL] HTTP server crashed: %v", err)
		}
	}()

	return srv, nil
}
