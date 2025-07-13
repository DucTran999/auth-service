package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DucTran999/auth-service/internal/handler/rest"
	"github.com/DucTran999/auth-service/test/setup"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLogoutAccount(t *testing.T) {
	tests := []struct {
		name           string
		setupBefore    func(t *testing.T, app *setup.TestApp) *http.Cookie
		expectedStatus int
	}{
		{
			name: "logout with valid session",
			setupBefore: func(t *testing.T, app *setup.TestApp) *http.Cookie {
				app.TruncateTables(t)

				// Register user
				registerPayload := `{"email":"logout@example.com", "password":"Strong123!"}`
				req := httptest.NewRequest(http.MethodPost, "/api/v1/register", strings.NewReader(registerPayload))
				req.Header.Set("Content-Type", "application/json")
				w := httptest.NewRecorder()
				app.Router.ServeHTTP(w, req)
				require.Equal(t, http.StatusCreated, w.Code)

				// Login user
				loginPayload := `{"email":"logout@example.com", "password":"Strong123!"}`
				req = httptest.NewRequest(http.MethodPost, "/api/v1/login", strings.NewReader(loginPayload))
				req.Header.Set("Content-Type", "application/json")
				w = httptest.NewRecorder()
				app.Router.ServeHTTP(w, req)
				require.Equal(t, http.StatusOK, w.Code)

				// Extract session cookie
				for _, c := range w.Result().Cookies() {
					if c.Name == rest.SessionKey {
						return c
					}
				}
				require.FailNow(t, "session cookie not found after login")
				return nil
			},
			expectedStatus: http.StatusNoContent,
		},
		{
			name: "logout without session cookie",
			setupBefore: func(t *testing.T, app *setup.TestApp) *http.Cookie {
				app.TruncateTables(t)
				return nil
			},
			expectedStatus: http.StatusNoContent,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			app, err := setup.NewTestApp()
			require.NoError(t, err)

			var cookie *http.Cookie
			if tc.setupBefore != nil {
				cookie = tc.setupBefore(t, app)
			}

			req := httptest.NewRequest(http.MethodPost, "/api/v1/logout", nil)
			req.Header.Set("Content-Type", "application/json")
			if cookie != nil {
				req.AddCookie(cookie)
			}

			w := httptest.NewRecorder()
			app.Router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatus, w.Code)

			// If cookie was passed in, check it was cleared
			if cookie != nil {
				var cleared bool
				for _, c := range w.Result().Cookies() {
					if c.Name == rest.SessionKey && c.MaxAge < 0 {
						cleared = true
						break
					}
				}
				assert.True(t, cleared, "session cookie should be cleared")
			}
		})
	}
}
