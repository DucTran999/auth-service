package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type EnvConfiguration struct {
	ServiceEnv   string `mapstructure:"SERVICE_ENV"`
	ServiceName  string `mapstructure:"SERVICE_NAME"`
	ServiceID    string `mapstructure:"SERVICE_ID"`
	ShutdownTime int    `mapstructure:"SHUTDOWN_TIME"`

	Host string `mapstructure:"HOST"`
	Port int    `mapstructure:"PORT"`

	LogToFile bool `mapstructure:"LOG_TO_FILE"`

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

func LoadConfig(configPath, configFile, configType string) (*EnvConfiguration, error) {
	viper.AddConfigPath(configPath)
	viper.SetConfigFile(configFile)
	viper.SetConfigType(configType)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var conf EnvConfiguration
	if err := viper.Unmarshal(&conf); err != nil {
		return nil, fmt.Errorf("error unmarshal config: %w", err)
	}

	return &conf, nil
}
