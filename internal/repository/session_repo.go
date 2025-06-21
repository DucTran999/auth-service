package repository

import (
	"context"

	"github.com/DucTran999/auth-service/internal/model"
)

type SessionRepository interface {
	Create(ctx context.Context, session *model.Session) error
}
