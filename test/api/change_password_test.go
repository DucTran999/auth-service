package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DucTran999/auth-service/internal/handler/rest"
	"github.com/DucTran999/auth-service/test/setup"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func registerAndLogin(t *testing.T, app *setup.TestApp, email, password string) *http.Cookie {
	t.Helper()

	// Register
	payload := fmt.Sprintf(`{"email":"%s","password":"%s"}`, email, password)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/register", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	app.Router.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	// Login
	req = httptest.NewRequest(http.MethodPost, "/api/v1/login", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	app.Router.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	for _, c := range w.Result().Cookies() {
		if c.Name == rest.SessionKey {
			return c
		}
	}
	return nil
}

func TestChangePassword_Success(t *testing.T) {
	app, err := setup.NewTestApp()
	require.NoError(t, err)
	app.TruncateTables(t)

	t.Run("change password success", func(t *testing.T) {
		app, err := setup.NewTestApp()
		require.NoError(t, err)
		app.TruncateTables(t)

		// Register and login
		sessionCookie := registerAndLogin(t, app, "user@example.com", "OldPass123!")

		//  Change password using the session cookie
		changePayload := `{"old_password":"OldPass123!", "new_password":"NewStrong123!"}`
		req := httptest.NewRequest(http.MethodPatch, "/api/v1/account/password", strings.NewReader(changePayload))
		req.Header.Set("Content-Type", "application/json")
		req.AddCookie(sessionCookie)

		w := httptest.NewRecorder()
		app.Router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
	})

	t.Run("wrong old password", func(t *testing.T) {
		app, err := setup.NewTestApp()
		require.NoError(t, err)
		app.TruncateTables(t)

		// Register and login
		sessionCookie := registerAndLogin(t, app, "user@example.com", "OldPass123!")

		// Attempt to change password with incorrect old password
		payload := `{"old_password":"WrongPass!", "new_password":"NewStrong123!"}`
		req := httptest.NewRequest(http.MethodPatch, "/api/v1/account/password", strings.NewReader(payload))
		req.Header.Set("Content-Type", "application/json")
		req.AddCookie(sessionCookie)

		w := httptest.NewRecorder()
		app.Router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("new password same as old", func(t *testing.T) {
		app, err := setup.NewTestApp()
		require.NoError(t, err)
		app.TruncateTables(t)

		// Register and login
		sessionCookie := registerAndLogin(t, app, "user@example.com", "OldPass123!")

		// Attempt to change to the same password
		payload := `{"old_password":"OldPass123!", "new_password":"OldPass123!"}`
		req := httptest.NewRequest(http.MethodPatch, "/api/v1/account/password", strings.NewReader(payload))
		req.Header.Set("Content-Type", "application/json")
		req.AddCookie(sessionCookie)

		w := httptest.NewRecorder()
		app.Router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code) // or another status code depending on your validation logic
	})
}
