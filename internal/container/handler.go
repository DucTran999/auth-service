package container

import "github.com/DucTran999/auth-service/internal/handler/rest"

type handlers struct {
	auth    AuthHandler
	account AccountHandler
	health  HealthHandler
}

func (c *Container) initHandlers() {
	c.handlers = &handlers{
		auth:    rest.NewAuthHandler(c.Logger, c.useCases.auth),
		account: rest.NewAccountHandler(c.Logger, c.useCases.account, c.useCases.restSession),
		health:  rest.NewHealthHandler(c.AppConfig.ServiceEnv),
	}
}
