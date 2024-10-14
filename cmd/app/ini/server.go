package ini

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/DucTran999/auth-service/config"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func startServer(router *gin.Engine, config *config.EnvConfiguration, logger *zap.Logger) {
	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	go func() {
		logger.Info("Starting server", zap.String("address", addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Server failed to start", zap.Error(err))
		}
	}()

	gracefulShutdown(srv, logger)
}

// gracefulShutdown handles OS signals and performs a graceful shutdown of the server.
func gracefulShutdown(srv *http.Server, logger *zap.Logger) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...")

	// Create a context with a timeout to ensure the server shuts down gracefully
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	select {
	case <-ctx.Done():
		logger.Info("Server shutdown timed out", zap.Duration("timeout", 5*time.Second))

	}

	logger.Info("Server exited cleanly")
}
