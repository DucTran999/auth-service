package http

import (
	"fmt"

	gen "github.com/DucTran999/auth-service/gen/http"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type RunningEnvironment int

const (
	ProductionEnv RunningEnvironment = iota
	DevelopmentEnv
)

func (r RunningEnvironment) String() string {
	switch r {
	case DevelopmentEnv:
		return "dev"
	case ProductionEnv:
		return "prod"
	// Set to default value dev if env invalid
	default:
		return "dev"
	}
}

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
