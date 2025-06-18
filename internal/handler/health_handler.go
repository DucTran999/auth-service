package handler

import (
	"net/http"
	"time"

	"github.com/DucTran999/auth-service/internal/gen"
	"github.com/gin-gonic/gin"
)

type HealthHandler interface {
	CheckLiveness(ctx *gin.Context)
}

type healthHandlerImpl struct {
	BaseHandler
	serviceVersion string // version release
	startTime      time.Time
}

func NewHealthHandler(serviceVersion string) *healthHandlerImpl {
	return &healthHandlerImpl{
		serviceVersion: serviceVersion,
		startTime:      time.Now(),
	}
}

func (h *healthHandlerImpl) CheckLiveness(ctx *gin.Context) {
	uptime := int64(time.Since(h.startTime).Seconds())

	response := gen.HealthResponse{
		Status:    gen.HealthResponseStatusHealthy,
		Timestamp: time.Now().UTC(),
		Uptime:    &uptime,
		Version:   &h.serviceVersion,
	}

	ctx.JSON(http.StatusOK, response)
}
