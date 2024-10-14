package logger

import (
	"log"

	"go.uber.org/zap"
)

// initLogger initializes and returns the zap logger.
func NewLogger() *zap.Logger {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	return logger
}
