package container

import (
	"github.com/DucTran999/auth-service/internal/handler/rest"
	"github.com/gin-gonic/gin"
)

type HealthHandler interface {
	CheckLiveness(ctx *gin.Context)
}

type AccountHandler interface {
	CreateAccount(ctx *gin.Context)
	ChangePassword(ctx *gin.Context)
}

type RestHandler struct {
	rest.SessionAuthHandler
	rest.JWTAuthHandler
	AccountHandler
	HealthHandler
}

func (c *Container) initRestHandler() {
	c.RestHandler = &RestHandler{
		JWTAuthHandler:     c.handlers.jwtAuth,
		SessionAuthHandler: c.handlers.sessionAuth,
		AccountHandler:     c.handlers.account,
		HealthHandler:      c.handlers.health,
	}
}
