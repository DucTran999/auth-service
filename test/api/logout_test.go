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

func TestLogoutAccount(t *testing.T) {
	tests := []struct {
		name           string
		setupBefore    func(t *testing.T, app *setup.TestApp)
		cookie         *http.Cookie
		expectedStatus int
	}{
		{
			name: "logout with valid session",
			setupBefore: func(t *testing.T, app *setup.TestApp) {
				// Reset database and create an account
				app.TruncateTables(t)

				// Register a new account to simulate a valid session later
				registerPayload := `{"email":"logout@example.com", "password":"Strong123!"}`
				req := httptest.NewRequest(http.MethodPost, "/api/v1/register", strings.NewReader(registerPayload))
				req.Header.Set("Content-Type", "application/json")
				w := httptest.NewRecorder()
				app.Router.ServeHTTP(w, req)
				require.Equal(t, http.StatusCreated, w.Code)
			},
			// You may extract a real session_id from login if needed
			cookie: &http.Cookie{
				Name:  "session_id",
				Value: "valid-session-id", // Replace with actual session ID if applicable
			},
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "logout without session cookie",
			setupBefore:    func(t *testing.T, app *setup.TestApp) { app.TruncateTables(t) },
			cookie:         nil,
			expectedStatus: http.StatusNoContent,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Setup app and dependencies
			app, err := setup.NewTestApp()
			require.NoError(t, err)

			if tc.setupBefore != nil {
				tc.setupBefore(t, app)
			}

			// Prepare HTTP request
			req := httptest.NewRequest(http.MethodPost, "/api/v1/logout", nil)
			req.Header.Set("Content-Type", "application/json")
			if tc.cookie != nil {
				req.AddCookie(tc.cookie)
			}

			// Send request
			w := httptest.NewRecorder()
			app.Router.ServeHTTP(w, req)

			// Check response status
			assert.Equal(t, tc.expectedStatus, w.Code)

			// If a session cookie was passed in, ensure it was cleared
			if tc.cookie != nil {
				found := false
				for _, c := range w.Result().Cookies() {
					if c.Name == "session_id" && c.MaxAge < 0 {
						found = true
						break
					}
				}
				assert.True(t, found, "session cookie should be cleared")
			}
		})
	}
}
