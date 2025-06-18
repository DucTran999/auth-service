package main

import (
	"log"

	"github.com/DucTran999/auth-service/cmd/auth-service/app"
	"github.com/DucTran999/auth-service/config"
)

func main() {
	// Load application configuration
	appConf, err := config.LoadConfig(".env")
	if err != nil {
		log.Fatalf("[FATAL] failed to load configuration: %v", err)
	}
	log.Println("[INFO] Configuration loaded successfully")

	// Initialize application with dependencies
	appInstance, err := app.NewApp(appConf)
	if err != nil {
		log.Fatalf("[FATAL] failed to initialize application: %v", err)
	}

	// Run the application (start server, wait for shutdown)
	if err := appInstance.Run(); err != nil {
		log.Fatalf("[FATAL] app got: %v", err)
	}
}
