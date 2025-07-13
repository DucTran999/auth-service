package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DucTran999/auth-service/test/setup"
	"github.com/go-playground/assert/v2"
	"github.com/stretchr/testify/require"
)

func TestRegister(t *testing.T) {
	tests := []struct {
		name           string
		payload        string
		beforeTest     func(t *testing.T, app *setup.TestApp)
		expectedStatus int
	}{
		{
			name: "success",
			payload: `{
				"email": "user1@example.com",
				"password": "StrongPass123!"
			}`,
			beforeTest:     func(t *testing.T, app *setup.TestApp) { app.TruncateTables(t) },
			expectedStatus: http.StatusCreated,
		},
		{
			name: "invalid email format",
			payload: `{
				"email": "invalid-email",
				"password": "StrongPass123!"
			}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "weak password",
			payload: `{
				"email": "user2@example.com",
				"password": "123"
			}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "empty fields",
			payload: `{
				"email": "",
				"password": ""
			}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "email already exists",
			payload: `{
				"email": "user3@example.com",
				"password": "StrongPass123!"
			}`,
			beforeTest: func(t *testing.T, app *setup.TestApp) {
				app.TruncateTables(t)

				// Register account to test duplicate
				req := httptest.NewRequest(http.MethodPost, "/api/v1/register", strings.NewReader(`{
					"email": "user3@example.com",
					"password": "StrongPass123!"
				}`))
				req.Header.Set("Content-Type", "application/json")
				w := httptest.NewRecorder()
				app.Router.ServeHTTP(w, req)

				require.Equal(t, http.StatusCreated, w.Code)
			},
			expectedStatus: http.StatusConflict,
		},
		{
			name:           "invalid json",
			payload:        `{invalid`,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			app, err := setup.NewTestApp()
			require.NoError(t, err)

			if tc.beforeTest != nil {
				tc.beforeTest(t, app)
			}

			req := httptest.NewRequest(http.MethodPost, "/api/v1/register", strings.NewReader(tc.payload))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			app.Router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatus, w.Code)
		})
	}
}
