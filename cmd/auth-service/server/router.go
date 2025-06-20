package server

import (
	"fmt"

	"github.com/DucTran999/auth-service/internal/gen"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func NewRouter(serviceEnv string, h gen.ServerInterface) (*gin.Engine, error) {
	if serviceEnv == ProductionEnv.String() {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// binding custom validator
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		if err := registerStrongPasswordValidators(v); err != nil {
			return nil, fmt.Errorf("failed to register strong validator: %w", err)
		}
	}

	gen.RegisterHandlers(router, h)

	return router, nil
}
