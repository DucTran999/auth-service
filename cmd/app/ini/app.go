package ini

import (
	"log"

	"github.com/DucTran999/auth-service/config"
	"github.com/DucTran999/auth-service/internal/gateway"
	"github.com/DucTran999/auth-service/internal/handler"
	"github.com/DucTran999/auth-service/internal/registry"
	"github.com/DucTran999/auth-service/pkg/logger"
)

func InitApp(config *config.EnvConfiguration) {
	logger := logger.NewLogger()
	defer logger.Sync()

	pg, err := connectDatabase(config)
	if err != nil {
		log.Fatalf("connect db failed got err: %v", err)
	}
	log.Println("DB connect successfully!")

	registry := registry.NewRegistry(pg)
	handler := handler.NewAppHandler(registry)
	httpServer := NewHttpServer(gateway.NewRouter(handler), config)
	go startHttpServer(httpServer)

	gracefulShutdown(shutdownHttpServer(httpServer))
}
