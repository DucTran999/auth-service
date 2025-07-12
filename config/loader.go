package config

import (
	"fmt"
	"path/filepath"

	"github.com/DucTran999/auth-service/pkg/projectpath"
	"github.com/spf13/viper"
)

func LoadConfig(configFile string) (*EnvConfiguration, error) {
	// Locate the directory containing go.mod
	root := projectpath.MustRoot()

	// Build the full path to the config file (e.g., .env)
	configPath := filepath.Join(root, configFile)

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

	conf.Normalize()
	if err := conf.Validate(); err != nil {
		return nil, fmt.Errorf("failed when validating %w", err)
	}

	return &conf, nil
}
