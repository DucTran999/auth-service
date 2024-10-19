package ini

import (
	"log"

	"github.com/DucTran999/auth-service/config"
	"github.com/DucTran999/auth-service/internal/gateway"
	"github.com/DucTran999/auth-service/internal/handler"
	"github.com/DucTran999/auth-service/internal/registry"
	"github.com/DucTran999/shared-pkg/v2/server"
)

func InitApp(config *config.EnvConfiguration) {
	pg, err := connectDatabase(config)
	if err != nil {
		log.Fatalf("connect db failed got err: %v", err)
	}
	log.Println("DB connect successfully!")

	registry := registry.NewRegistry(pg)
	handler := handler.NewAppHandler(registry)
	httpServer := server.NewGinHttpServer(gateway.NewRouter(handler), config.Host, config.Port)

	go func() {
		if err := httpServer.Start(); err != nil {
			log.Fatalf("start http server got err: %v", err)
		}
	}()

	server.GracefulShutdown(httpServer.Stop)
}
