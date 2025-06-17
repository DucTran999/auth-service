package container

import (
	"github.com/DucTran999/auth-service/config"
	"github.com/DucTran999/auth-service/internal/common"
	"github.com/DucTran999/shared-pkg/logger"
)

func newLogger(cfg *config.EnvConfiguration) (logger.ILogger, error) {
	return logger.NewLogger(logger.Config{
		Environment: cfg.ServiceEnv,
		LogToFile:   cfg.LogToFile,
		FilePath:    common.LogFilePath,
	})
}
