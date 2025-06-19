package server

import (
	"github.com/DucTran999/auth-service/internal/gen"
	"github.com/gin-gonic/gin"
)

func NewRouter(h gen.ServerInterface) *gin.Engine {
	router := gin.Default()

	gen.RegisterHandlers(router, h)

	return router
}
