package handler

import "github.com/gin-gonic/gin"

type baseHandler struct{}

func (baseHandler) JsonResponse(ctx *gin.Context, httpCode int, data any, message string) {
	hasError := 0
	if message != "" {
		hasError = 1
	}

	ctx.JSON(httpCode, gin.H{"status": hasError, "message": message, "data": data})
}
