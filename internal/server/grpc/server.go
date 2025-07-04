package grpc

import (
	"fmt"
	"net"
	"strconv"

	"github.com/DucTran999/auth-service/config"
	"github.com/DucTran999/auth-service/gen/grpc/pb"
	"github.com/DucTran999/shared-pkg/logger"
	"google.golang.org/grpc"
)

type GRPCServer interface {
	Start() error
	GracefulStop()
}

type grpcServer struct {
	lis    net.Listener
	server *grpc.Server
	logger logger.ILogger
}

func (g *grpcServer) Start() error {
	g.logger.Infof("grpc running on: %s", g.lis.Addr().String())
	return g.server.Serve(g.lis)
}

func (g *grpcServer) GracefulStop() {
	g.server.GracefulStop()
}

func NewGRPCServer(
	cfg *config.EnvConfiguration,
	apiHandler pb.AuthServiceServer,
	logger logger.ILogger,
) (GRPCServer, error) {
	addr := net.JoinHostPort(cfg.GRPCHost, strconv.Itoa(cfg.GRPCPort))
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("failed to listen on %s: %w", addr, err)
	}

	srv := grpc.NewServer()

	pb.RegisterAuthServiceServer(srv, apiHandler)

	return &grpcServer{
		lis:    listener,
		server: srv,
		logger: logger,
	}, nil
}
