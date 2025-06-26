package worker

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/DucTran999/auth-service/internal/config"
	"github.com/DucTran999/auth-service/internal/handler/background"
	"github.com/DucTran999/shared-pkg/logger"
)

type SessionCleanupWorker interface {
	Start(ctx context.Context) error
}

type sessionCleanupWorker struct {
	expireEvery   time.Duration
	purgeEvery    time.Duration
	runningExpire atomic.Bool
	runningPurge  atomic.Bool

	wg      sync.WaitGroup
	logger  logger.ILogger
	cleaner background.SessionCleaner
}

func NewSessionCleanupWorker(
	logger logger.ILogger,
	cfg *config.EnvConfiguration,
	cleaner background.SessionCleaner,
) *sessionCleanupWorker {

	return &sessionCleanupWorker{
		logger:        logger,
		cleaner:       cleaner,
		expireEvery:   time.Duration(cfg.ExpireIntervalInMins) * time.Minute,
		purgeEvery:    time.Duration(cfg.PurgeIntervalInDays) * time.Hour * 24,
		runningExpire: atomic.Bool{},
		runningPurge:  atomic.Bool{},
		wg:            sync.WaitGroup{},
	}
}

func (w *sessionCleanupWorker) Start(ctx context.Context) error {
	w.logger.Info("start session cleanup worker")
	expireTicker := time.NewTicker(w.expireEvery)
	purgeTicker := time.NewTicker(w.purgeEvery)
	defer expireTicker.Stop()
	defer purgeTicker.Stop()

	for {
		select {
		case <-expireTicker.C:
			if !w.runningExpire.CompareAndSwap(false, true) {
				w.logger.Warn("expire still running, skipping tick")
				continue
			}

			w.wg.Add(1)
			go func() {
				defer w.wg.Done()
				defer w.runningExpire.Store(false)
				w.runExpire(ctx)
			}()
		case <-purgeTicker.C:
			if !w.runningPurge.CompareAndSwap(false, true) {
				w.logger.Warn("purge still running, skipping tick")
				continue
			}

			w.wg.Add(1)
			go func() {
				defer w.wg.Done()
				defer w.runningPurge.Store(false)
				w.runPurge(ctx)
			}()
		case <-ctx.Done():
			w.logger.Info("session cleanup worker stopping...")
			w.wg.Wait()
			return ctx.Err()
		}
	}
}

func (w *sessionCleanupWorker) runExpire(ctx context.Context) {
	w.logger.Info("running expire untracked sessions")
	start := time.Now()

	defer func() {
		if r := recover(); r != nil {
			w.logger.Errorf("panic in expire: %v", r)
		}
		w.logger.Infof("expire finished in %s", time.Since(start))
	}()

	// Error is already logged by cleaner — no need to log again here
	_ = w.cleaner.ExpireUntrackedSessions(ctx)
}

func (w *sessionCleanupWorker) runPurge(ctx context.Context) {
	w.logger.Info("running purge expired sessions")
	start := time.Now()

	defer func() {
		if r := recover(); r != nil {
			w.logger.Errorf("panic in purge: %v", r)
		}
		w.logger.Infof("purge finished in %s", time.Since(start))
	}()

	// Error is already logged by cleaner — no need to log again here
	_ = w.cleaner.PurgeExpiredSessions(ctx)
}
