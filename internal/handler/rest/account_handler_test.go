package rest_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	gen "github.com/DucTran999/auth-service/gen/http"
	httpServer "github.com/DucTran999/auth-service/internal/server/http"
	mockbuilder "github.com/DucTran999/auth-service/test/mock-builder"

	"github.com/DucTran999/auth-service/internal/handler/rest"
	"github.com/DucTran999/shared-pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/stretchr/testify/require"
)

func NewAccountHandlerUT(t *testing.T, builder *mockbuilder.UsecaseBuilderContainer) rest.AccountHandler {
	logger, err := logger.NewLogger(logger.Config{Environment: "staging"})
	require.NoError(t, err)

	return rest.NewAccountHandler(
		logger,
		builder.AccountUC.GetInstance(),
		builder.SessionUC.GetInstance(),
	)
}

func TestCreateAccount(t *testing.T) {
	type testcase struct {
		name           string
		setupUT        func(t *testing.T) rest.AccountHandler
		setupPayload   func(t *testing.T) []byte
		expectedStatus int
	}

	tests := []testcase{
		{
			name: "missing email",
			setupUT: func(t *testing.T) rest.AccountHandler {
				t.Helper()
				b := mockbuilder.NewUsecaseBuilderContainer(t)
				return NewAccountHandlerUT(t, b)
			},
			setupPayload: func(t *testing.T) []byte {
				t.Helper()
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
			name: "invalid email",
			setupUT: func(t *testing.T) rest.AccountHandler {
				t.Helper()
				b := mockbuilder.NewUsecaseBuilderContainer(t)
				return NewAccountHandlerUT(t, b)
			},
			setupPayload: func(t *testing.T) []byte {
				t.Helper()
				payload := map[string]any{
					"email":    "invalidEmail.com",
					"password": "p@ssG0rk1234!",
				}
				jsonPayload, err := json.Marshal(payload)
				require.NoError(t, err)
				return jsonPayload
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "invalid payload",
			setupUT: func(t *testing.T) rest.AccountHandler {
				t.Helper()
				b := mockbuilder.NewUsecaseBuilderContainer(t)
				return NewAccountHandlerUT(t, b)
			},
			setupPayload: func(t *testing.T) []byte {
				t.Helper()
				payload := map[string]any{
					"email":    1234,
					"password": "p@ssG0rk1234!",
				}
				jsonPayload, err := json.Marshal(payload)
				require.NoError(t, err)
				return jsonPayload
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "missing password",
			setupUT: func(t *testing.T) rest.AccountHandler {
				t.Helper()
				b := mockbuilder.NewUsecaseBuilderContainer(t)
				return NewAccountHandlerUT(t, b)
			},
			setupPayload: func(t *testing.T) []byte {
				t.Helper()
				payload := map[string]any{
					"email": "daniel@example.com",
				}
				jsonPayload, err := json.Marshal(payload)
				require.NoError(t, err)
				return jsonPayload
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "weak password",
			setupUT: func(t *testing.T) rest.AccountHandler {
				t.Helper()
				b := mockbuilder.NewUsecaseBuilderContainer(t)
				return NewAccountHandlerUT(t, b)
			},
			setupPayload: func(t *testing.T) []byte {
				t.Helper()
				payload := map[string]any{
					"email":    "daniel@example.com",
					"password": "weak",
				}
				jsonPayload, err := json.Marshal(payload)
				require.NoError(t, err)
				return jsonPayload
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "email already registered",
			setupUT: func(t *testing.T) rest.AccountHandler {
				t.Helper()
				b := mockbuilder.NewUsecaseBuilderContainer(t)
				b.AccountUC.RegisterConflictEmail()
				return NewAccountHandlerUT(t, b)
			},
			setupPayload: func(t *testing.T) []byte {
				t.Helper()
				payload := gen.CreateAccountRequest{
					Email:    "test@example.com",
					Password: "p@ssG0rk1234!",
				}
				jsonPayload, err := json.Marshal(payload)
				require.NoError(t, err)
				return jsonPayload
			},
			expectedStatus: http.StatusConflict,
		},
		{
			name: "register failed",
			setupUT: func(t *testing.T) rest.AccountHandler {
				t.Helper()
				b := mockbuilder.NewUsecaseBuilderContainer(t)
				b.AccountUC.RegisterError()
				return NewAccountHandlerUT(t, b)
			},
			setupPayload: func(t *testing.T) []byte {
				t.Helper()
				payload := gen.CreateAccountRequest{
					Email:    "test@example.com",
					Password: "p@ssG0rk1234!",
				}
				jsonPayload, err := json.Marshal(payload)
				require.NoError(t, err)
				return jsonPayload
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "register success",
			setupUT: func(t *testing.T) rest.AccountHandler {
				t.Helper()
				b := mockbuilder.NewUsecaseBuilderContainer(t)
				b.AccountUC.RegisterSuccess()
				return NewAccountHandlerUT(t, b)
			},
			setupPayload: func(t *testing.T) []byte {
				t.Helper()
				payload := gen.CreateAccountRequest{
					Email:    "test@example.com",
					Password: "p@ssG0rk1234!",
				}
				jsonPayload, err := json.Marshal(payload)
				require.NoError(t, err)
				return jsonPayload
			},
			expectedStatus: http.StatusCreated,
		},
	}

	err := httpServer.SetupValidator()
	require.NoError(t, err)
	gin.SetMode(gin.TestMode)

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// t.Parallel()
			handler := tc.setupUT(t)

			// Setup handler with mock
			gin.SetMode(gin.TestMode)
			router := gin.New()
			router.POST("/api/v1/register", handler.CreateAccount)

			// Setup response recorder
			req := httptest.NewRequest(
				http.MethodPost,
				"/api/v1/register",
				bytes.NewBuffer(tc.setupPayload(t)),
			)
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatus, w.Code)
		})
	}
}

func TestChangePassword(t *testing.T) {
	type testcase struct {
		name           string
		mustSetCookie  bool
		setupUT        func(t *testing.T) rest.AccountHandler
		setupPayload   func(t *testing.T) []byte
		expectedStatus int
	}

	tests := []testcase{
		{
			name: "missing cookie",
			setupUT: func(t *testing.T) rest.AccountHandler {
				t.Helper()
				b := mockbuilder.NewUsecaseBuilderContainer(t)
				return NewAccountHandlerUT(t, b)
			},
			setupPayload:   func(t *testing.T) []byte { return []byte{} },
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "session cookie invalid cause DB Err",
			setupUT: func(t *testing.T) rest.AccountHandler {
				t.Helper()
				b := mockbuilder.NewUsecaseBuilderContainer(t)
				b.SessionUC.ValidateError()
				return NewAccountHandlerUT(t, b)
			},
			mustSetCookie:  true,
			setupPayload:   func(t *testing.T) []byte { return []byte{} },
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "session cookie invalid",
			setupUT: func(t *testing.T) rest.AccountHandler {
				t.Helper()
				b := mockbuilder.NewUsecaseBuilderContainer(t)
				b.SessionUC.ValidateInvalidSession()
				return NewAccountHandlerUT(t, b)
			},
			mustSetCookie:  true,
			setupPayload:   func(t *testing.T) []byte { return []byte{} },
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "invalid password",
			setupUT: func(t *testing.T) rest.AccountHandler {
				t.Helper()
				b := mockbuilder.NewUsecaseBuilderContainer(t)
				b.SessionUC.ValidateSessionSuccess()
				return NewAccountHandlerUT(t, b)
			},
			mustSetCookie: true,
			setupPayload: func(t *testing.T) []byte {
				t.Helper()
				payload := gen.ChangePasswordRequest{
					OldPassword: "0ldP@ssGrork1234",
					NewPassword: "",
				}
				jsonPayload, err := json.Marshal(payload)
				require.NoError(t, err)
				return jsonPayload
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "new password is same as old password",
			setupUT: func(t *testing.T) rest.AccountHandler {
				t.Helper()
				b := mockbuilder.NewUsecaseBuilderContainer(t)
				b.SessionUC.ValidateSessionSuccess()
				b.AccountUC.ChangePasswordGotErrorSamePass()
				return NewAccountHandlerUT(t, b)
			},
			mustSetCookie: true,
			setupPayload: func(t *testing.T) []byte {
				t.Helper()
				payload := gen.ChangePasswordRequest{
					OldPassword: "0ldP@ssGrork1234",
					NewPassword: "newP@ssGok1234!",
				}
				jsonPayload, err := json.Marshal(payload)
				require.NoError(t, err)
				return jsonPayload
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "failed to change password cause invalid credentials",
			setupUT: func(t *testing.T) rest.AccountHandler {
				t.Helper()
				b := mockbuilder.NewUsecaseBuilderContainer(t)
				b.SessionUC.ValidateSessionSuccess()
				b.AccountUC.ChangePassErrGotWrongCredentials()
				return NewAccountHandlerUT(t, b)
			},
			mustSetCookie: true,
			setupPayload: func(t *testing.T) []byte {
				t.Helper()
				payload := gen.ChangePasswordRequest{
					OldPassword: "0ldP@ssGrork1234",
					NewPassword: "newP@ssGok1234!",
				}
				jsonPayload, err := json.Marshal(payload)
				require.NoError(t, err)
				return jsonPayload
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "failed to change password got DB error",
			setupUT: func(t *testing.T) rest.AccountHandler {
				t.Helper()
				b := mockbuilder.NewUsecaseBuilderContainer(t)
				b.SessionUC.ValidateSessionSuccess()
				b.AccountUC.ChangePassErrGotErrorDB()
				return NewAccountHandlerUT(t, b)
			},
			mustSetCookie: true,
			setupPayload: func(t *testing.T) []byte {
				t.Helper()
				payload := gen.ChangePasswordRequest{
					OldPassword: "0ldP@ssGrork1234",
					NewPassword: "newP@ssGok1234!",
				}
				jsonPayload, err := json.Marshal(payload)
				require.NoError(t, err)
				return jsonPayload
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "change password success",
			setupUT: func(t *testing.T) rest.AccountHandler {
				t.Helper()
				b := mockbuilder.NewUsecaseBuilderContainer(t)
				b.SessionUC.ValidateSessionSuccess()
				b.AccountUC.ChangePasswordSuccess()
				return NewAccountHandlerUT(t, b)
			},
			mustSetCookie: true,
			setupPayload: func(t *testing.T) []byte {
				t.Helper()
				payload := gen.ChangePasswordRequest{
					OldPassword: "0ldP@ssGrork1234",
					NewPassword: "newP@ssGok1234!",
				}
				jsonPayload, err := json.Marshal(payload)
				require.NoError(t, err)
				return jsonPayload
			},
			expectedStatus: http.StatusNoContent,
		},
	}

	err := httpServer.SetupValidator()
	require.NoError(t, err)
	gin.SetMode(gin.TestMode)

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			handler := tc.setupUT(t)

			// Setup handler with mock
			gin.SetMode(gin.TestMode)
			router := gin.New()
			router.POST("/api/v1/account/password", handler.ChangePassword)

			// Make request
			req := httptest.NewRequest(
				http.MethodPost,
				"/api/v1/account/password",
				bytes.NewBuffer(tc.setupPayload(t)),
			)
			req.Header.Set("Content-Type", "application/json")
			if tc.mustSetCookie {
				req.AddCookie(&http.Cookie{
					Name:  "session_id",
					Value: "mock-session-id",
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
