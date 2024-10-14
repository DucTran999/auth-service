package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type EnvConfiguration struct {
	Environment string `mapstructure:"ENVIRONMENT"`

	Host string `mapstructure:"HOST"`
	Port int    `mapstructure:"PORT"`

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
		return nil, fmt.Errorf("error unmarshalling config: %w", err)
	}

	return &conf, nil
}
