package server

import (
	"github.com/DucTran999/auth-service/cmd/auth-service/container"
	"github.com/DucTran999/shared-pkg/server"
)

// NewHTTPServer creates a new HTTP server with injected dependencies
func NewHTTPServer(deps container.Container) (server.HttpServer, error) {
	cfg := deps.AppConfig()

	serverConf := server.ServerConfig{
		Host: cfg.Host,
		Port: cfg.Port,
	}

	router := NewRouter(deps.AppConfig().ServiceEnv, deps.AppHandler())
	httpServer, err := server.NewGinHttpServer(router, serverConf)
	if err != nil {
		return nil, err
	}

	return httpServer, nil
}
