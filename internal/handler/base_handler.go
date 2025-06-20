package handler

import (
	"net/http"

	"github.com/DucTran999/auth-service/internal/gen"
	"github.com/gin-gonic/gin"
)

const (
	ApiVersion1 = "v1.0"

	// General Errors
	InternalErrorCode      = "INTERNAL_ERROR"
	ValidateErrorCode      = "VALIDATE_ERROR"
	BadRequestErrorCode    = "BAD_REQUEST"
	UnauthorizedErrorCode  = "UNAUTHORIZED"
	ForbiddenErrorCode     = "FORBIDDEN"
	NotFoundErrorCode      = "NOT_FOUND"
	MethodNotAllowedCode   = "METHOD_NOT_ALLOWED"
	ConflictErrorCode      = "CONFLICT"
	TooManyRequestsCode    = "TOO_MANY_REQUESTS"
	ServiceUnavailableCode = "SERVICE_UNAVAILABLE"
)

type BaseHandler struct{}

func (BaseHandler) ValidateErrorResponse(ctx *gin.Context, version string, err string) {
	respBody := gen.BadRequest{
		Version: version,
		Error: gen.ErrorDetail{
			Code:    ValidateErrorCode,
			Message: err,
		},
	}

	ctx.JSON(http.StatusBadRequest, respBody)
}

func (BaseHandler) BadRequestResponse(ctx *gin.Context, version, errMsg string) {
	respBody := gen.BadRequest{
		Version: version,
		Error: gen.ErrorDetail{
			Code:    BadRequestErrorCode,
			Message: errMsg,
		},
	}

	ctx.JSON(http.StatusBadRequest, respBody)
}

func (BaseHandler) ResourceConflictResponse(ctx *gin.Context, version, errMsg string) {
	respBody := gen.Conflict{
		Version: version,
		Error: gen.ErrorDetail{
			Code:    ConflictErrorCode,
			Message: errMsg,
		},
	}

	ctx.JSON(http.StatusConflict, respBody)
}

func (BaseHandler) ServerInternalErrResponse(ctx *gin.Context, version string) {
	respBody := gen.InternalServerError{
		Version: version,
		Error: gen.ErrorDetail{
			Code:    InternalErrorCode,
			Message: http.StatusText(http.StatusInternalServerError),
		},
	}

	ctx.JSON(http.StatusInternalServerError, respBody)
}
