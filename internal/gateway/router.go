package gateway

import (
	"github.com/DucTran999/auth-service/internal/handler"
	"github.com/gin-gonic/gin"
)

func NewRouter(h handler.AppHandler) *gin.Engine {
	router := gin.Default()

	v1 := router.Group("/api/v1")
	{
		v1.POST("/users", h.CreateUser)
	}

	return router
}
