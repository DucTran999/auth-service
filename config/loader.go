package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

func LoadConfig(configFile string) (*EnvConfiguration, error) {
	// Locate the directory containing go.mod
	goModDir, err := findGoModDir()
	if err != nil {
		// If go.mod is not found, use the current working directory
		goModDir, err = os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("failed to get current working directory: %w", err)
		}
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

	conf.Normalize()
	if err = conf.Validate(); err != nil {
		return nil, fmt.Errorf("failed when validating %w", err)
	}

	return &conf, nil
}

// FindGoModDir returns the directory containing the nearest go.mod file.
func findGoModDir() (string, error) {
	// cwd
	startPath, err := os.Getwd()
	if err != nil {
		return "", err
	}

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
