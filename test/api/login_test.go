package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DucTran999/auth-service/test/setup"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLogin(t *testing.T) {
	tests := []struct {
		name           string
		setupBefore    func(t *testing.T, app *setup.TestApp)
		payload        string
		expectedStatus int
	}{
		{
			name: "login success",
			setupBefore: func(t *testing.T, app *setup.TestApp) {
				app.TruncateTables(t)
				req := httptest.NewRequest(http.MethodPost, "/api/v1/register", strings.NewReader(`{
					"email": "login@example.com",
					"password": "Strong123!"
				}`))
				req.Header.Set("Content-Type", "application/json")
				w := httptest.NewRecorder()
				app.Router.ServeHTTP(w, req)

				require.Equal(t, http.StatusCreated, w.Code)
			},
			payload: `{
				"email": "login@example.com",
				"password": "Strong123!"
			}`,
			expectedStatus: http.StatusOK,
		},
		{
			name: "invalid credentials",
			setupBefore: func(t *testing.T, app *setup.TestApp) {
				app.TruncateTables(t)
			},
			payload: `{
				"email": "notfound@example.com",
				"password": "WrongPass123!"
			}`,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "missing fields",
			payload: `{
				"email": "",
				"password": ""
			}`,
			expectedStatus: http.StatusBadRequest,
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

			if tc.setupBefore != nil {
				tc.setupBefore(t, app)
			}

			req := httptest.NewRequest(http.MethodPost, "/api/v1/login", strings.NewReader(tc.payload))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			app.Router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatus, w.Code)
		})
	}
}
