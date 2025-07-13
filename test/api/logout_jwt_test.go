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

func TestJWTLogout(t *testing.T) {
	app, err := setup.NewTestApp()
	require.NoError(t, err)
	app.TruncateTables(t)

	t.Run("failed to logout with invalid token", func(t *testing.T) {
		// Simulate a request that contains a valid refresh_token cookie
		req := httptest.NewRequest(http.MethodPost, "/api/v2/logout", nil)
		req.Header.Set("Content-Type", "application/json")
		req.AddCookie(&http.Cookie{
			Name:     "refresh_token",
			Value:    "mock-refresh-token", // simulated token
			Path:     "/",
			HttpOnly: true,
		})

		// Send the request
		w := httptest.NewRecorder()
		app.Router.ServeHTTP(w, req)

		// Expect HTTP 204 No Content
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("missing refresh token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/v2/logout", nil)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		app.Router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("valid refresh token", func(t *testing.T) {
		// First, register and login to get a real refresh_token
		register := `{"email": "jwtlogout@example.com", "password": "Strong123!"}`
		req := httptest.NewRequest(http.MethodPost, "/api/v1/register", strings.NewReader(register))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		app.Router.ServeHTTP(w, req)
		require.Equal(t, http.StatusCreated, w.Code)

		// Login to get refresh_token
		login := `{"email": "jwtlogout@example.com", "password": "Strong123!"}`
		req = httptest.NewRequest(http.MethodPost, "/api/v2/login", strings.NewReader(login))
		req.Header.Set("Content-Type", "application/json")
		w = httptest.NewRecorder()
		app.Router.ServeHTTP(w, req)
		require.Equal(t, http.StatusOK, w.Code)

		// Extract refresh_token from Set-Cookie
		var refreshToken string
		for _, c := range w.Result().Cookies() {
			if c.Name == "refresh_token" {
				refreshToken = c.Value
			}
		}
		require.NotEmpty(t, refreshToken)

		// Call logout with valid token
		req = httptest.NewRequest(http.MethodPost, "/api/v2/logout", nil)
		req.Header.Set("Content-Type", "application/json")
		req.AddCookie(&http.Cookie{
			Name:  "refresh_token",
			Value: refreshToken,
		})

		w = httptest.NewRecorder()
		app.Router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)

		// Check cookie is cleared
		found := false
		for _, c := range w.Result().Cookies() {
			if c.Name == "refresh_token" && c.MaxAge < 0 {
				found = true
			}
		}
		assert.True(t, found, "refresh_token cookie should be cleared")
	})
}
