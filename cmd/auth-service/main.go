package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"
)

func main() {
	appCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Load application configuration
	cfg, err := loadConfig()
	if err != nil {
		log.Printf("[ERROR] %v", err)
		return
	}

	// Initialize container
	c, err := initContainer(cfg)
	if err != nil {
		log.Printf("[ERROR] %v", err)
		return
	}
	defer c.Close()

	// start Rest server
	restSrv, err := startHTTPServer(c)
	if err != nil {
		log.Printf("[ERROR] %v", err)
		return
	}

	// start grpc server
	grpcSrv, err := startGRPCServer(c)
	if err != nil {
		log.Printf("[ERROR] %v", err)
		return
	}

	workerDone := runSessionCleanupWorker(appCtx, c)

	// gracefully shutdown
	waitForShutdown(appCtx, restSrv, grpcSrv, workerDone, c)
}
