package ini

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/DucTran999/auth-service/config"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func NewHttpServer(router *gin.Engine, config *config.EnvConfiguration) *http.Server {
	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)
	server := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	return server
}

func startHttpServer(server *http.Server) {
	log.Printf("Server listening on: %v", server.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal("Server failed to start", zap.Error(err))
	}
}

// gracefulShutdown handles OS signals and performs a graceful shutdown of the server.
func gracefulShutdown(gracefulHandler ...func() error) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Create a context with a timeout to ensure the server shuts down gracefully
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	isExitClean := true
	for _, f := range gracefulHandler {
		if err := f(); err != nil {
			if isExitClean {
				isExitClean = false
			}
			log.Println(err)
		}
	}

	<-ctx.Done()
	log.Println("Server shutdown timed out")
	if isExitClean {
		log.Println("Server exited cleanly")
	} else {
		log.Println("Server exited. Got some error while shutting down")
	}
}

func shutdownHttpServer(server *http.Server) func() error {
	return func() error {
		if err := server.Shutdown(context.Background()); err != nil {
			return fmt.Errorf("graceful shutdown server got err: %v", err)
		}
		log.Printf("http server on %v shutdown successfully!", server.Addr)

		return nil
	}
}
