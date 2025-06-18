package gateway

import (
	"github.com/DucTran999/auth-service/internal/gen"
	"github.com/DucTran999/auth-service/internal/handler"
	"github.com/gin-gonic/gin"
)

func NewRouter(h handler.AppHandler) *gin.Engine {
	router := gin.Default()

	gen.RegisterHandlers(router, h)

	return router
}
