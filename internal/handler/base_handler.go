package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	respSuccessCode        = 0
	respErrorCode          = 1
	successMsg             = "OK"
	serverInternalErrMsg   = "Server Internal Error"
	resourceConflictErrMsg = "The resource already exists"
)

type BaseHandler struct{}

func (BaseHandler) SuccessResponse(ctx *gin.Context, version string, data any) {
	respBody := gin.H{
		"version":   version,
		"errorCode": respSuccessCode,
		"message":   successMsg,
		"data":      data,
	}

	ctx.JSON(http.StatusOK, respBody)
}

func (BaseHandler) BadRequestResponse(ctx *gin.Context, version string, err error) {
	respBody := gin.H{
		"version":   version,
		"errorCode": respErrorCode,
		"message":   err.Error(),
		"data":      nil,
	}

	ctx.JSON(http.StatusBadRequest, respBody)
}

func (BaseHandler) ResourceConflictResponse(ctx *gin.Context, version string) {
	respBody := gin.H{
		"version":   version,
		"errorCode": respErrorCode,
		"message":   resourceConflictErrMsg,
		"data":      nil,
	}

	ctx.JSON(http.StatusConflict, respBody)
}

func (BaseHandler) ServerInternalErrResponse(ctx *gin.Context, version string) {
	respBody := gin.H{
		"version":   version,
		"errorCode": respErrorCode,
		"message":   serverInternalErrMsg,
		"data":      nil,
	}

	ctx.JSON(http.StatusInternalServerError, respBody)
}
