package container

import "github.com/DucTran999/auth-service/internal/handler/grpc"

func (c *Container) initGRPCHandler() {
	c.GRPCHandler = grpc.NewAuthHandler()
}
