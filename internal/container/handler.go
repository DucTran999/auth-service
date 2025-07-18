package container

import "github.com/DucTran999/auth-service/internal/handler/rest"

type handlers struct {
	sessionAuth rest.SessionAuthHandler
	jwtAuth     rest.JWTAuthHandler
	account     rest.AccountHandler
	health      rest.HealthHandler
}

func (c *Container) initHandlers() {
	c.handlers = &handlers{
		sessionAuth: rest.NewSessionAuthHandler(c.Logger, c.useCases.sessionAuth),
		jwtAuth:     rest.NewJWTAuthHandler(c.Logger, c.useCases.jwtAuth),
		account:     rest.NewAccountHandler(c.Logger, c.useCases.account, c.useCases.session),
		health:      rest.NewHealthHandler(c.AppConfig.ServiceEnv),
	}
}
