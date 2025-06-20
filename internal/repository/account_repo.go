package repository

import (
	"context"

	"github.com/DucTran999/auth-service/internal/model"
)

type AccountRepo interface {
	FindByEmail(ctx context.Context, email string) (*model.Account, error)

	Create(ctx context.Context, account model.Account) (*model.Account, error)
}
