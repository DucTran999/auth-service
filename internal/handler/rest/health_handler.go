package rest

import (
	"net/http"
	"time"

	gen "github.com/DucTran999/auth-service/gen/http"
	"github.com/gin-gonic/gin"
)

type HealthHandler interface {
	CheckLiveness(ctx *gin.Context)
}

type healthHandler struct {
	serviceVersion string // version release
	startTime      time.Time
}

func NewHealthHandler(serviceVersion string) HealthHandler {
	return &healthHandler{
		serviceVersion: serviceVersion,
		startTime:      time.Now(),
	}
}

func (h *healthHandler) CheckLiveness(ctx *gin.Context) {
	uptime := int64(time.Since(h.startTime).Seconds())

	response := gen.HealthResponse{
		Status:    gen.HealthResponseStatusHealthy,
		Timestamp: time.Now().UTC(),
		Uptime:    &uptime,
		Version:   &h.serviceVersion,
	}

	ctx.JSON(http.StatusOK, response)
}
