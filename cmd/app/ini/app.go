package ini

import (
	"log"

	"github.com/DucTran999/auth-service/config"
	"github.com/DucTran999/auth-service/internal/common"
	"github.com/DucTran999/auth-service/internal/gateway"
	"github.com/DucTran999/auth-service/internal/handler"
	"github.com/DucTran999/auth-service/internal/registry"
	"github.com/DucTran999/shared-pkg/logger"
	"github.com/DucTran999/shared-pkg/server"
)

func InitApp(config *config.EnvConfiguration) {
	logInst := initLogger(config)
	logInst.Info("Logger instance initialize successfully!")

	dbInst, err := connectDatabase(config)
	if err != nil {
		logInst.Fatalf("connect db failed got err: %v", err)
	}
	logInst.Info("DB connect successfully!")

	registry := registry.NewRegistry(dbInst.GetConn())
	handler := handler.NewAppHandler(registry)
	httpServer := server.NewGinHttpServer(gateway.NewRouter(handler), config.Host, config.Port)

	go func() {
		if err := httpServer.Start(); err != nil {
			logInst.Fatalf("start http server got err: %v", err)
		}
	}()

	server.GracefulShutdown(httpServer.Stop, closeDBConnection(dbInst))
}

func initLogger(appConf *config.EnvConfiguration) logger.ILogger {
	conf := logger.Config{
		Environment: appConf.ServiceEnv,
		LogToFile:   appConf.LogToFile,
		FilePath:    common.LogFilePath,
	}

	logger, err := logger.NewLogger(conf)
	if err != nil {
		log.Fatalln("Init logger ERR", err)
	}

	return logger
}
