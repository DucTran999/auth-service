package main

import (
	"log"

	"github.com/DucTran999/auth-service/cmd/auth-service/app"
	"github.com/DucTran999/auth-service/config"
)

func main() {
	const (
		configPath = "."
		configFile = ".env"
		configType = "env"
	)

	appConf, err := config.LoadConfig(configPath, configFile, configType)
	if err != nil {
		log.Fatalln("failed to load configurations", err)
	}
	log.Println("[INFO] load config successfully!")

	app, err := app.NewApp(appConf)
	if err != nil {
		log.Fatalln("failed to start app")
	}

	app.Run()
}
