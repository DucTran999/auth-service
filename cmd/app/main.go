package main

import (
	"log"

	"github.com/DucTran999/auth-service/cmd/app/ini"
	"github.com/DucTran999/auth-service/config"
)

func main() {
	const (
		configPath = "."
		configFile = ".env"
		configType = "env"
	)

	config, err := config.LoadConfig(configPath, configFile, configType)
	if err != nil {
		log.Fatal("Failed to load configurations", err)
	}
	log.Println("Load config successfully!")

	ini.InitApp(config)
}
