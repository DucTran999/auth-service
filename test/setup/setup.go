package setup

import (
	"fmt"
	"testing"

	"github.com/DucTran999/auth-service/config"
	"github.com/DucTran999/auth-service/internal/container"
	"github.com/DucTran999/auth-service/internal/server/http"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

type TestApp struct {
	Router *gin.Engine
	DI     *container.Container
}

func NewTestApp() (*TestApp, error) {
	gin.SetMode(gin.TestMode)

	cfg, err := config.LoadConfig(".test.env")
	if err != nil {
		return nil, fmt.Errorf("failed to load test config: %w", err)
	}

	ctn, err := container.NewContainer(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to build test container: %w", err)
	}

	if err := http.SetupValidator(); err != nil {
		return nil, fmt.Errorf("failed to setup validator: %w", err)
	}

	router, err := http.NewRouter(cfg.ServiceEnv, ctn.RestHandler)
	if err != nil {
		return nil, fmt.Errorf("failed to init router: %w", err)
	}

	return &TestApp{
		Router: router,
		DI:     ctn,
	}, nil
}

func (app *TestApp) TruncateTables(t *testing.T) {
	t.Helper()
	db := app.DI.AuthDB.DB()
	err := db.Exec("TRUNCATE sessions, accounts CASCADE").Error
	require.NoError(t, err)
}
