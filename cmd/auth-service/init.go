package main

import (
	"fmt"
	"log"

	"github.com/DucTran999/auth-service/config"
	"github.com/DucTran999/auth-service/internal/container"
)

func loadConfig() (*config.EnvConfiguration, error) {
	cfg, err := config.LoadConfig(".env")
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}
	log.Println("[INFO] Configuration loaded successfully")
	return cfg, nil
}

func initContainer(cfg *config.EnvConfiguration) (*container.Container, error) {
	c, err := container.NewContainer(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize container: %w", err)
	}
	log.Println("[INFO] DI container initialized")
	return c, nil
}
