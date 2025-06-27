package server

import (
	"github.com/DucTran999/auth-service/config"
	"github.com/DucTran999/auth-service/internal/gen"
	"github.com/DucTran999/shared-pkg/server"
)

// NewHTTPServer creates a new HTTP server with injected dependencies.
func NewHTTPServer(cfg *config.EnvConfiguration, apiHandler gen.ServerInterface) (server.HttpServer, error) {
	serverConf := server.ServerConfig{
		Host: cfg.Host,
		Port: cfg.Port,
	}

	router, err := NewRouter(cfg.ServiceEnv, apiHandler)
	if err != nil {
		return nil, err
	}

	httpServer, err := server.NewGinHttpServer(router, serverConf)
	if err != nil {
		return nil, err
	}

	return httpServer, nil
}
