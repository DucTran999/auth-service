package rest_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	gen "github.com/DucTran999/auth-service/gen/http"
	"github.com/DucTran999/auth-service/internal/handler/rest"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCheckLiveness(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	router := gin.New()

	const version = "v1.2.3"
	handler := rest.NewHealthHandler(version)

	router.GET("/healthz", handler.CheckLiveness)

	// Make request
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	var res gen.HealthResponse
	err := json.Unmarshal(w.Body.Bytes(), &res)
	require.NoError(t, err)

	assert.Equal(t, gen.HealthResponseStatusHealthy, res.Status)
	assert.NotNil(t, res.Timestamp)
	assert.NotNil(t, res.Uptime)
	assert.NotNil(t, res.Version)
	assert.Equal(t, version, *res.Version)
	assert.GreaterOrEqual(t, *res.Uptime, int64(0))
}
