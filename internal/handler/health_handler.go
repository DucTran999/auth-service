package handler

import (
	"time"

	"github.com/DucTran999/auth-service/internal/common"
	"github.com/DucTran999/auth-service/internal/gen"
	"github.com/gin-gonic/gin"
)

type HealthHandler struct {
	BaseHandler
}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) CheckHealth(ctx *gin.Context) {
	response := gen.HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now().UTC(),
	}

	h.SuccessResponse(ctx, common.ApiVersion1, response)
}
