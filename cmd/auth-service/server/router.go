package server

import (
	"github.com/DucTran999/auth-service/internal/gen"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func NewRouter(serviceEnv string, h gen.ServerInterface) *gin.Engine {
	if serviceEnv == ProductionEnv.String() {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// binding custom validator
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		registerStrongPasswordValidators(v)
	}

	gen.RegisterHandlers(router, h)

	return router
}
