package rest_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	gen "github.com/DucTran999/auth-service/gen/http"
	"github.com/DucTran999/auth-service/internal/handler/rest"
	mockbuilder "github.com/DucTran999/auth-service/test/mock-builder"
	"github.com/DucTran999/shared-pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func NewAuthSessionHandlerUT(t *testing.T, builder *mockbuilder.UsecaseBuilderContainer) rest.SessionAuthHandler {
	t.Helper()

	log, err := logger.NewLogger(logger.Config{
		Environment: "staging",
	})
	require.NoError(t, err)

	return rest.NewSessionAuthHandler(log, builder.AuthSessionUC.GetInstance())
}

func TestLoginAccount(t *testing.T) {
	type testcase struct {
		name           string
		setupUT        func(t *testing.T) rest.SessionAuthHandler
		setupPayload   func(t *testing.T) []byte
		tokenKey       string
		expectedStatus int
	}

	tests := []testcase{
		{
			name: "invalid json payload",
			setupUT: func(t *testing.T) rest.SessionAuthHandler {
				builder := mockbuilder.NewUsecaseBuilderContainer(t)
				return NewAuthSessionHandlerUT(t, builder)
			},
			setupPayload: func(t *testing.T) []byte {
				// Invalid JSON
				return []byte(`{invalid`)
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "invalid credentials",
			setupUT: func(t *testing.T) rest.SessionAuthHandler {
				builder := mockbuilder.NewUsecaseBuilderContainer(t)
				builder.AuthSessionUC.LoginInvalidCredentials()
				return NewAuthSessionHandlerUT(t, builder)
			},
			tokenKey: rest.RefreshTokenKey,
			setupPayload: func(t *testing.T) []byte {
				req := gen.LoginAccountJSONRequestBody{
					Email:    "wrong@example.com",
					Password: "wrongpass",
				}
				b, _ := json.Marshal(req)
				return b
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "internal error during login",
			setupUT: func(t *testing.T) rest.SessionAuthHandler {
				builder := mockbuilder.NewUsecaseBuilderContainer(t)
				builder.AuthSessionUC.LoginInternalError()
				return NewAuthSessionHandlerUT(t, builder)
			},
			tokenKey: rest.RefreshTokenKey,
			setupPayload: func(t *testing.T) []byte {
				req := gen.LoginAccountJSONRequestBody{
					Email:    "user@example.com",
					Password: "validPass123!",
				}
				b, _ := json.Marshal(req)
				return b
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "login first time cookie",
			setupUT: func(t *testing.T) rest.SessionAuthHandler {
				builder := mockbuilder.NewUsecaseBuilderContainer(t)
				builder.AuthSessionUC.LoginSuccess()
				return NewAuthSessionHandlerUT(t, builder)
			},
			setupPayload: func(t *testing.T) []byte {
				req := gen.LoginAccountJSONRequestBody{
					Email:    "user@example.com",
					Password: "validPass123!",
				}
				b, _ := json.Marshal(req)
				return b
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "login success",
			setupUT: func(t *testing.T) rest.SessionAuthHandler {
				builder := mockbuilder.NewUsecaseBuilderContainer(t)
				builder.AuthSessionUC.LoginSuccess()
				return NewAuthSessionHandlerUT(t, builder)
			},
			tokenKey: rest.RefreshTokenKey,
			setupPayload: func(t *testing.T) []byte {
				req := gen.LoginAccountJSONRequestBody{
					Email:    "user@example.com",
					Password: "validPass123!",
				}
				b, _ := json.Marshal(req)
				return b
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			handler := tc.setupUT(t)

			gin.SetMode(gin.TestMode)
			router := gin.New()
			router.POST("/api/v1/login", handler.LoginAccount)

			req := httptest.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewBuffer(tc.setupPayload(t)))
			req.Header.Set("Content-Type", "application/json")

			// Optional: simulate existing session cookie
			req.AddCookie(&http.Cookie{
				Name:  tc.tokenKey,
				Value: "mock",
			})

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatus, w.Code)
		})
	}
}

func TestLogoutAccount(t *testing.T) {
	type testcase struct {
		name           string
		mustSetCookie  bool
		setupUT        func(t *testing.T) rest.SessionAuthHandler
		expectedStatus int
	}

	tests := []testcase{
		{
			name: "no cookie provided",
			setupUT: func(t *testing.T) rest.SessionAuthHandler {
				builder := mockbuilder.NewUsecaseBuilderContainer(t)
				return NewAuthSessionHandlerUT(t, builder)
			},
			expectedStatus: http.StatusNoContent,
		},
		{
			name:          "logout success with cookie",
			mustSetCookie: true,
			setupUT: func(t *testing.T) rest.SessionAuthHandler {
				builder := mockbuilder.NewUsecaseBuilderContainer(t)
				builder.AuthSessionUC.LogoutSuccess()
				return NewAuthSessionHandlerUT(t, builder)
			},
			expectedStatus: http.StatusNoContent,
		},
		{
			name:          "logout error but still 204",
			mustSetCookie: true,
			setupUT: func(t *testing.T) rest.SessionAuthHandler {
				builder := mockbuilder.NewUsecaseBuilderContainer(t)
				builder.AuthSessionUC.LogoutError()
				return NewAuthSessionHandlerUT(t, builder)
			},
			expectedStatus: http.StatusNoContent,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			handler := tc.setupUT(t)

			gin.SetMode(gin.TestMode)
			router := gin.New()
			router.POST("/api/v1/logout", handler.LogoutAccount)

			req := httptest.NewRequest(http.MethodPost, "/api/v1/logout", nil)
			if tc.mustSetCookie {
				req.AddCookie(&http.Cookie{
					Name:  "session_id",
					Value: "mock-session-id",
					Path:  "/",
				})
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatus, w.Code)

			// Cookie must be cleared in response
			cookies := w.Result().Cookies()
			foundCleared := false
			for _, c := range cookies {
				if c.Name == "session_id" && c.MaxAge < 0 {
					foundCleared = true
				}
			}
			assert.True(t, foundCleared, "session_id cookie should be cleared")
		})
	}
}
