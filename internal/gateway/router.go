package gateway

import (
	"github.com/DucTran999/auth-service/cmd/auth-service/container"
	"github.com/DucTran999/auth-service/internal/gen"
	"github.com/gin-gonic/gin"
)

func NewRouter(h container.AppHandler) *gin.Engine {
	router := gin.Default()

	gen.RegisterHandlers(router, h)

	return router
}
