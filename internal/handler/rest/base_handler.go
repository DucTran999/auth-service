package rest

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/DucTran999/auth-service/internal/gen"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

const (
	ApiVersion1 = "v1.0"

	// General Errors.
	InternalErrorCode      = "INTERNAL_ERROR"
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

func ParseAndValidateJSON[T any](ctx *gin.Context) (*T, error) {
	var payload T
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) && len(ve) > 0 {
			// You must inject a common error response function
			return nil, fmt.Errorf("validation error: %s", validationErrorMessage(ve[0]))
		}
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}
	return &payload, nil
}

func validationErrorMessage(fe validator.FieldError) string {
	switch fe.Field() {
	case "Email":
		switch fe.Tag() {
		case "required":
			return "Email is required"
		case "email":
			return "Email must be valid"
		}
	case "Password":
		switch fe.Tag() {
		case "required":
			return "Password is required"
		case "password":
			return "Password must include at least 1 uppercase, 1 lowercase, 1 number, and 1 special character"
		}
	}

	return "Invalid input"
}

func (BaseHandler) UnauthorizeErrorResponse(ctx *gin.Context, version string, err string) {
	respBody := gen.Unauthorized{
		Version: version,
		Error: gen.ErrorDetail{
			Code:    UnauthorizedErrorCode,
			Message: err,
		},
	}

	ctx.JSON(http.StatusUnauthorized, respBody)
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

func (BaseHandler) NoContentResponse(ctx *gin.Context) {
	ctx.Status(http.StatusNoContent)
}
