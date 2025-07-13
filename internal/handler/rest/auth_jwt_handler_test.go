package rest_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DucTran999/auth-service/internal/handler/rest"
	mockbuilder "github.com/DucTran999/auth-service/test/mock-builder"
	"github.com/DucTran999/shared-pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/stretchr/testify/require"
)

func NewAuthJWTHandlerUT(t *testing.T, builder *mockbuilder.UsecaseBuilderContainer) rest.JWTAuthHandler {
	logger, err := logger.NewLogger(logger.Config{Environment: "staging"})
	require.NoError(t, err)

	return rest.NewJWTAuthHandler(
		logger,
		builder.AuthJwtUC.GetInstance(),
	)
}

func TestLoginWithJWT(t *testing.T) {
	type testcase struct {
		name           string
		setupUT        func(t *testing.T) rest.JWTAuthHandler
		setupPayload   func(t *testing.T) []byte
		expectedStatus int
	}

	tests := []testcase{
		{
			name: "missing email",
			setupUT: func(t *testing.T) rest.JWTAuthHandler {
				builder := mockbuilder.NewUsecaseBuilderContainer(t)
				return NewAuthJWTHandlerUT(t, builder)
			},
			setupPayload: func(t *testing.T) []byte {
				payload := map[string]any{
					"password": "p@ssG0rk1234!",
				}
				jsonPayload, err := json.Marshal(payload)
				require.NoError(t, err)
				return jsonPayload
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "wrong credentials",
			setupUT: func(t *testing.T) rest.JWTAuthHandler {
				builder := mockbuilder.NewUsecaseBuilderContainer(t)
				builder.AuthJwtUC.LoginErrWrongCredentials()
				return NewAuthJWTHandlerUT(t, builder)
			},
			setupPayload: func(t *testing.T) []byte {
				payload := map[string]any{
					"email":    "danial@example.com",
					"password": "p@ssG0rk1234!",
				}
				jsonPayload, err := json.Marshal(payload)
				require.NoError(t, err)
				return jsonPayload
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "login DB error",
			setupUT: func(t *testing.T) rest.JWTAuthHandler {
				builder := mockbuilder.NewUsecaseBuilderContainer(t)
				builder.AuthJwtUC.LoginErrDB()
				return NewAuthJWTHandlerUT(t, builder)
			},
			setupPayload: func(t *testing.T) []byte {
				payload := map[string]any{
					"email":    "danial@example.com",
					"password": "p@ssG0rk1234!",
				}
				jsonPayload, err := json.Marshal(payload)
				require.NoError(t, err)
				return jsonPayload
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "login success",
			setupUT: func(t *testing.T) rest.JWTAuthHandler {
				builder := mockbuilder.NewUsecaseBuilderContainer(t)
				builder.AuthJwtUC.LoginSuccess()
				return NewAuthJWTHandlerUT(t, builder)
			},
			setupPayload: func(t *testing.T) []byte {
				payload := map[string]any{
					"email":    "danial@example.com",
					"password": "p@ssG0rk1234!",
				}
				jsonPayload, err := json.Marshal(payload)
				require.NoError(t, err)
				return jsonPayload
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			handler := tc.setupUT(t)

			// Setup handler with mock
			gin.SetMode(gin.TestMode)
			router := gin.New()
			router.POST("/api/v2/login", handler.LoginWithJWT)

			// Make request
			req := httptest.NewRequest(
				http.MethodPost,
				"/api/v2/login",
				bytes.NewBuffer(tc.setupPayload(t)),
			)
			req.Header.Set("Content-Type", "application/json")

			// Setup response recorder
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tc.expectedStatus, w.Code)
		})
	}
}

func TestLogoutJWT(t *testing.T) {
	type testcase struct {
		name           string
		mustSetCookie  bool
		setupUT        func(t *testing.T) rest.JWTAuthHandler
		expectedStatus int
	}

	tests := []testcase{
		{
			name: "missing refresh token in cookie",
			setupUT: func(t *testing.T) rest.JWTAuthHandler {
				builder := mockbuilder.NewUsecaseBuilderContainer(t)
				return NewAuthJWTHandlerUT(t, builder)
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "revoke token got internal error",
			setupUT: func(t *testing.T) rest.JWTAuthHandler {
				builder := mockbuilder.NewUsecaseBuilderContainer(t)
				builder.AuthJwtUC.RevokeRefreshTokenErr()
				return NewAuthJWTHandlerUT(t, builder)
			},
			mustSetCookie:  true,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "logout success",
			setupUT: func(t *testing.T) rest.JWTAuthHandler {
				builder := mockbuilder.NewUsecaseBuilderContainer(t)
				builder.AuthJwtUC.RevokeRefreshTokenSuccess()
				return NewAuthJWTHandlerUT(t, builder)
			},
			mustSetCookie:  true,
			expectedStatus: http.StatusNoContent,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			handler := tc.setupUT(t)

			// Setup handler with mock
			gin.SetMode(gin.TestMode)
			router := gin.New()
			router.POST("/api/v2/logout", handler.LogoutJWT)

			// Make request
			req := httptest.NewRequest(http.MethodPost, "/api/v2/logout", nil)
			req.Header.Set("Content-Type", "application/json")
			if tc.mustSetCookie {
				req.AddCookie(&http.Cookie{
					Name:  "refresh_token",
					Value: "mock-refresh-token",
					Path:  "/",
				})
			}

			// Setup response recorder
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tc.expectedStatus, w.Code)
		})
	}
}

func TestRefreshToken(t *testing.T) {
	type testcase struct {
		name           string
		mustSetCookie  bool
		setupUT        func(t *testing.T) rest.JWTAuthHandler
		expectedStatus int
	}

	tests := []testcase{
		{
			name: "missing refresh token in cookie",
			setupUT: func(t *testing.T) rest.JWTAuthHandler {
				builder := mockbuilder.NewUsecaseBuilderContainer(t)
				return NewAuthJWTHandlerUT(t, builder)
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "invalid refresh token usecase returns error",
			setupUT: func(t *testing.T) rest.JWTAuthHandler {
				builder := mockbuilder.NewUsecaseBuilderContainer(t)
				builder.AuthJwtUC.RefreshTokenError()
				return NewAuthJWTHandlerUT(t, builder)
			},
			mustSetCookie:  true,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "refresh token success",
			setupUT: func(t *testing.T) rest.JWTAuthHandler {
				builder := mockbuilder.NewUsecaseBuilderContainer(t)
				builder.AuthJwtUC.RefreshTokenSuccess()
				return NewAuthJWTHandlerUT(t, builder)
			},
			mustSetCookie:  true,
			expectedStatus: http.StatusOK,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			handler := tc.setupUT(t)

			gin.SetMode(gin.TestMode)
			router := gin.New()
			router.POST("/api/v2/refresh", handler.RefreshToken)

			req := httptest.NewRequest(http.MethodPost, "/api/v2/refresh", nil)
			req.Header.Set("Content-Type", "application/json")
			if tc.mustSetCookie {
				req.AddCookie(&http.Cookie{
					Name:  "refresh_token",
					Value: "mock-refresh-token",
					Path:  "/",
				})
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatus, w.Code)
		})
	}
}
