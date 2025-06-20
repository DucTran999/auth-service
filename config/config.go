package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type EnvConfiguration struct {
	ServiceEnv     string `mapstructure:"SERVICE_ENV"`
	ServiceName    string `mapstructure:"SERVICE_NAME"`
	ServiceID      string `mapstructure:"SERVICE_ID"`
	ServiceVersion string `mapstructure:"SERVICE_VERSION"`
	ShutdownTime   int    `mapstructure:"SHUTDOWN_TIME"`

	Host string `mapstructure:"HOST"`
	Port int    `mapstructure:"PORT"`

	LogToFile   bool   `mapstructure:"LOG_TO_FILE"`
	LogFilePath string `mapstructure:"LOG_FILE_PATH"`

	DBDriver                string `mapstructure:"DB_DRIVER"`
	DBHost                  string `mapstructure:"DB_HOST"`
	DBPort                  int    `mapstructure:"DB_PORT"`
	DBUsername              string `mapstructure:"DB_USERNAME"`
	DBPasswd                string `mapstructure:"DB_PASSWORD"`
	DBDatName               string `mapstructure:"DB_DATABASE"`
	DBSslMode               string `mapstructure:"DB_SSL_MODE"`
	DBMaxOpenConnections    int    `mapstructure:"DB_MAX_OPEN_CONNECTIONS"`
	DBMaxIdleConnections    int    `mapstructure:"DB_MAX_IDLE_CONNECTIONS"`
	DBMaxConnectionIdleTime int    `mapstructure:"DB_MAX_CONNECTION_IDLE_TIME"`
	DBTimezone              string `mapstructure:"DB_TIMEZONE"`

	RedisHost   string `mapstructure:"REDIS_HOST"`
	RedisPort   int    `mapstructure:"REDIS_PORT"`
	RedisPasswd string `mapstructure:"REDIS_PASSWORD"`
	RedisDB     int    `mapstructure:"REDIS_DATABASE"`
}

func LoadConfig(configFile string) (*EnvConfiguration, error) {
	// Locate the directory containing go.mod
	goModDir, err := findGoModDir()
	if err != nil {
		return nil, fmt.Errorf("go.mod not found: %w", err)
	}

	// Build the full path to the config file (e.g., .env)
	configPath := filepath.Join(goModDir, configFile)

	// Load the config file using Viper
	viper.SetConfigFile(configPath)
	viper.AutomaticEnv() // override with environment variables if present

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", configPath, err)
	}

	// Unmarshal into custom config struct
	var conf EnvConfiguration
	if err := viper.Unmarshal(&conf); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &conf, nil
}

// FindGoModDir returns the directory containing the nearest go.mod file
func findGoModDir() (string, error) {
	// cwd
	startPath := "."

	dir := startPath
	if fi, err := os.Stat(startPath); err == nil && !fi.IsDir() {
		dir = filepath.Dir(startPath)
	}

	for {
		goMod := filepath.Join(dir, "go.mod")
		if _, err := os.Stat(goMod); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break // root reached
		}
		dir = parent
	}

	return "", fmt.Errorf("go.mod not found from path: %s", startPath)
}
