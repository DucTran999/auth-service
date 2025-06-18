package server

import (
	"github.com/DucTran999/auth-service/cmd/auth-service/container"
	"github.com/DucTran999/auth-service/config"
	"github.com/DucTran999/auth-service/internal/gateway"
	"github.com/DucTran999/auth-service/internal/handler"
	"github.com/DucTran999/auth-service/internal/registry"
	"github.com/DucTran999/shared-pkg/server"
)

// NewHTTPServer creates a new HTTP server with injected dependencies
func NewHTTPServer(cfg *config.EnvConfiguration, deps container.Container) (server.HttpServer, error) {
	registry := registry.NewRegistry(cfg, deps.AuthDB())
	serverConf := server.ServerConfig{
		Host: cfg.Host,
		Port: cfg.Port,
	}

	handler := handler.NewAppHandler(registry)
	httpServer, err := server.NewGinHttpServer(gateway.NewRouter(handler), serverConf)
	if err != nil {
		return nil, err
	}

	return httpServer, nil
}
