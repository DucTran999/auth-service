package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	gen "github.com/DucTran999/auth-service/gen/http"
	"github.com/DucTran999/auth-service/internal/handler/rest"
	"github.com/DucTran999/auth-service/test/setup"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJWTLogin(t *testing.T) {
	app, err := setup.NewTestApp()
	require.NoError(t, err)
	app.TruncateTables(t)

	tests := []struct {
		name           string
		setupBefore    func(t *testing.T, app *setup.TestApp)
		payload        string
		expectedStatus int
		checkToken     bool
	}{
		{
			name: "login success",
			setupBefore: func(t *testing.T, app *setup.TestApp) {
				req := httptest.NewRequest(http.MethodPost, "/api/v1/register", strings.NewReader(`{
					"email": "jwt@example.com",
					"password": "StrongPass123!"
				}`))
				req.Header.Set("Content-Type", "application/json")
				w := httptest.NewRecorder()
				app.Router.ServeHTTP(w, req)
				require.Equal(t, http.StatusCreated, w.Code)
			},
			payload: `{
				"email": "jwt@example.com",
				"password": "StrongPass123!"
			}`,
			expectedStatus: http.StatusOK,
			checkToken:     true,
		},
		{
			name:           "invalid json",
			payload:        `{invalid`,
			expectedStatus: http.StatusBadRequest,
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
			name: "wrong credentials",
			setupBefore: func(t *testing.T, app *setup.TestApp) {
				// No register needed
			},
			payload: `{
				"email": "notfound@example.com",
				"password": "Wrong123!"
			}`,
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupBefore != nil {
				tc.setupBefore(t, app)
			}

			req := httptest.NewRequest(http.MethodPost, "/api/v2/login", strings.NewReader(tc.payload))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("User-Agent", "test-agent")
			req.RemoteAddr = "127.0.0.1:12345"

			w := httptest.NewRecorder()
			app.Router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatus, w.Code)

			if tc.checkToken {
				// check that access_token exist
				resp := gen.LoginJWTResponse{}
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				require.NoError(t, err)
				assert.NotEmpty(t, resp.Data.AccessToken)

				// Also check refresh token set in cookie
				found := false
				for _, c := range w.Result().Cookies() {
					if c.Name == rest.RefreshTokenKey {
						found = true
						break
					}
				}
				assert.True(t, found, "refresh_token cookie must be set")

				// Check api version is 2
				assert.Equal(t, rest.APIVersion2, resp.Version)
			}
		})
	}
}
