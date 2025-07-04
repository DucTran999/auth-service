package container

import "github.com/DucTran999/auth-service/internal/handler/rest"

type handlers struct {
	auth    AuthHandler
	account AccountHandler
	health  HealthHandler
}

func (c *container) initHandlers() {
	c.handlers = &handlers{
		auth:    rest.NewAuthHandler(c.logger, c.useCases.auth),
		account: rest.NewAccountHandler(c.logger, c.useCases.account, c.useCases.restSession),
		health:  rest.NewHealthHandler(c.appConfig.ServiceEnv),
	}
}
