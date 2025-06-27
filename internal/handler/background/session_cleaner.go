package background

import (
	"context"
	"fmt"
	"time"

	"github.com/DucTran999/auth-service/internal/domain"
	"github.com/DucTran999/shared-pkg/logger"
)

const (
	sessionRetention = 30 // 30 days
)

type SessionCleaner interface {
	ExpireUntrackedSessions(ctx context.Context) error

	PurgeExpiredSessions(ctx context.Context) error
}

type sessionCleaner struct {
	logger    logger.ILogger
	sessionUC domain.SessionUsecase
}

func NewSessionCleaner(logger logger.ILogger, sessionUC domain.SessionUsecase) *sessionCleaner {
	return &sessionCleaner{
		logger:    logger,
		sessionUC: sessionUC,
	}
}

func (sc *sessionCleaner) ExpireUntrackedSessions(ctx context.Context) error {
	err := sc.sessionUC.MarkExpiredSessions(ctx)
	if err != nil {
		sc.logger.Errorf("failed to set expiration on untracked sessions: %v", err)
		return fmt.Errorf("expire untracked sessions failed: %w", err)
	}

	return nil
}

func (sc *sessionCleaner) PurgeExpiredSessions(ctx context.Context) error {
	cutoff := time.Now().AddDate(0, 0, -sessionRetention)

	if err := sc.sessionUC.DeleteExpiredBefore(ctx, cutoff); err != nil {
		sc.logger.Errorf("failed to purge expired sessions: %v", err)
		return fmt.Errorf("purge expired sessions: %w", err)
	}

	return nil
}
