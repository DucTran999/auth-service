package main

import (
	"context"
	"time"

	"github.com/DucTran999/auth-service/internal/container"
	"github.com/DucTran999/shared-pkg/server"
)

func waitForShutdown(appCtx context.Context, srv server.HttpServer, workerDone <-chan struct{}, c *container.Container) {
	<-appCtx.Done()
	c.Logger.Info("shutdown signal received")

	shutdownStart := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.AppConfig.ShutdownTime)*time.Second)
	defer cancel()

	if err := srv.Stop(ctx); err != nil {
		c.Logger.Errorf("failed to stop http server: %v", err)
	}

	select {
	case <-workerDone:
		c.Logger.Infof("session cleanup worker stopped in %s", time.Since(shutdownStart))
	case <-ctx.Done():
		c.Logger.Warnf("timeout waiting for session cleanup worker (after %s)", time.Since(shutdownStart))
	}
}
