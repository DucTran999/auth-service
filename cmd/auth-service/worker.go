package main

import (
	"context"
	"errors"

	"github.com/DucTran999/auth-service/internal/container"
	"github.com/DucTran999/auth-service/internal/worker"
)

func runSessionCleanupWorker(ctx context.Context, c *container.Container) <-chan struct{} {
	done := make(chan struct{})
	worker := worker.NewSessionCleanupWorker(c.Logger, c.AppConfig, c.CleanupSessionHandler)

	go func() {
		defer close(done)

		if err := worker.Start(ctx); err != nil &&
			!errors.Is(err, context.Canceled) &&
			!errors.Is(err, context.DeadlineExceeded) {
			c.Logger.Fatalf("[FATAL] session cleanup worker failed: %v", err)
		}
		c.Logger.Info("session cleanup worker exited cleanly")
	}()

	return done
}
