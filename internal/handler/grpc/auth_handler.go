package grpc

import (
	"context"

	"github.com/DucTran999/auth-service/gen/grpc/pb"
)

type AuthHandler struct {
	pb.UnimplementedAuthServiceServer
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

func (hdl *AuthHandler) Login(context.Context, *pb.LoginRequest) (*pb.LoginResponse, error) {
	return &pb.LoginResponse{AccessToken: "done"}, nil
}
