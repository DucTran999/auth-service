package container

import "github.com/gin-gonic/gin"

type HealthHandler interface {
	CheckLiveness(ctx *gin.Context)
}

type AuthHandler interface {
	LoginAccount(ctx *gin.Context)
	LoginAccountJWT(ctx *gin.Context)
	LogoutAccount(ctx *gin.Context)
}

type AccountHandler interface {
	CreateAccount(ctx *gin.Context)
	ChangePassword(ctx *gin.Context)
}

type RestHandler struct {
	AuthHandler
	AccountHandler
	HealthHandler
}

func (c *Container) initRestHandler() {
	c.RestHandler = &RestHandler{
		AuthHandler:    c.handlers.auth,
		AccountHandler: c.handlers.account,
		HealthHandler:  c.handlers.health,
	}
}
