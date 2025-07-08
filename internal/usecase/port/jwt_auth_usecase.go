package port

import (
	"context"

	"github.com/DucTran999/auth-service/internal/usecase/dto"
)

type JWTAuthUsecase interface {
	Login(ctx context.Context, input dto.LoginJWTInput) (*dto.TokenPairs, error)

	RefreshToken(ctx context.Context, refreshToken string) (*dto.TokenPairs, error)

	RevokeRefreshToken(ctx context.Context, refreshToken string) error
}
