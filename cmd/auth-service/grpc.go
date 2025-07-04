package main

import (
	"fmt"

	"github.com/DucTran999/auth-service/internal/container"
	"github.com/DucTran999/auth-service/internal/server/grpc"
	serverGRPC "github.com/DucTran999/auth-service/internal/server/grpc"
)

func startGRPCServer(c *container.Container) (grpc.GRPCServer, error) {
	// start grpc server
	grpcSrv, err := serverGRPC.NewGRPCServer(c.AppConfig, c.GRPCHandler, c.Logger)
	if err != nil {
		return nil, fmt.Errorf("failed to init grpc server: %w", err)
	}

	go func() {
		if err := grpcSrv.Start(); err != nil {
			c.Logger.Fatalf("[FATAL] grpc server crashed: %v", err)
		}
	}()

	return grpcSrv, nil
}
