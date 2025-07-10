package container

import (
	"github.com/DucTran999/auth-service/internal/handler/rest"
)

type RestHandler struct {
	rest.SessionAuthHandler
	rest.JWTAuthHandler
	rest.AccountHandler
	rest.HealthHandler
}

func (c *Container) initRestHandler() {
	c.RestHandler = &RestHandler{
		JWTAuthHandler:     c.handlers.jwtAuth,
		SessionAuthHandler: c.handlers.sessionAuth,
		AccountHandler:     c.handlers.account,
		HealthHandler:      c.handlers.health,
	}
}
