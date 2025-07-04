package config

import (
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
)

type EnvConfiguration struct {
	ServiceEnv     string `mapstructure:"SERVICE_ENV" validate:"required,oneof=dev staging prod"`
	ServiceName    string `mapstructure:"SERVICE_NAME" validate:"required"`
	ServiceID      string `mapstructure:"SERVICE_ID" validate:"required"`
	ServiceVersion string `mapstructure:"SERVICE_VERSION" validate:"required"`
	ShutdownTime   int    `mapstructure:"SHUTDOWN_TIME" validate:"gte=0"`

	Host string `mapstructure:"HOST" validate:"required"`
	Port int    `mapstructure:"PORT" validate:"required,min=1,max=65535"`

	GRPCHost string `mapstructure:"GRPC_HOST" validate:"required"`
	GRPCPort int    `mapstructure:"GRPC_PORT" validate:"required"`

	LogToFile   bool   `mapstructure:"LOG_TO_FILE"`
	LogFilePath string `mapstructure:"LOG_FILE_PATH"` // validate if LogToFile is true (custom)

	DBDriver                string `mapstructure:"DB_DRIVER" validate:"required"`
	DBHost                  string `mapstructure:"DB_HOST" validate:"required"`
	DBPort                  int    `mapstructure:"DB_PORT" validate:"required,min=1,max=65535"`
	DBUsername              string `mapstructure:"DB_USERNAME" validate:"required"`
	DBPasswd                string `mapstructure:"DB_PASSWORD" validate:"required"`
	DBDatName               string `mapstructure:"DB_DATABASE" validate:"required"`
	DBSslMode               string `mapstructure:"DB_SSL_MODE" validate:"omitempty,oneof=disable require verify-ca verify-full"`
	DBMaxOpenConnections    int    `mapstructure:"DB_MAX_OPEN_CONNECTIONS" validate:"gte=0"`
	DBMaxIdleConnections    int    `mapstructure:"DB_MAX_IDLE_CONNECTIONS" validate:"gte=0"`
	DBMaxConnectionIdleTime int    `mapstructure:"DB_MAX_CONNECTION_IDLE_TIME" validate:"gte=0"`
	DBTimezone              string `mapstructure:"DB_TIMEZONE" validate:"omitempty"`

	RedisHost   string `mapstructure:"REDIS_HOST" validate:"required"`
	RedisPort   int    `mapstructure:"REDIS_PORT" validate:"required,min=1,max=65535"`
	RedisPasswd string `mapstructure:"REDIS_PASSWORD"`
	RedisDB     int    `mapstructure:"REDIS_DATABASE" validate:"gte=0"`

	PurgeIntervalInDays  int `mapstructure:"PURGE_INTERVAL_IN_DAYS" validate:"gte=0"`
	ExpireIntervalInMins int `mapstructure:"EXPIRE_INTERVAL_IN_MINS" validate:"gte=0"`
}

func (cfg *EnvConfiguration) Normalize() {
	cfg.ServiceEnv = strings.TrimSpace(cfg.ServiceEnv)
	cfg.ServiceName = strings.TrimSpace(cfg.ServiceName)
	cfg.ServiceID = strings.TrimSpace(cfg.ServiceID)
	cfg.ServiceVersion = strings.TrimSpace(cfg.ServiceVersion)

	cfg.Host = strings.TrimSpace(cfg.Host)
	cfg.LogFilePath = strings.TrimSpace(cfg.LogFilePath)

	cfg.DBDriver = strings.TrimSpace(cfg.DBDriver)
	cfg.DBHost = strings.TrimSpace(cfg.DBHost)
	cfg.DBUsername = strings.TrimSpace(cfg.DBUsername)
	cfg.DBPasswd = strings.TrimSpace(cfg.DBPasswd)
	cfg.DBDatName = strings.TrimSpace(cfg.DBDatName)
	cfg.DBSslMode = strings.TrimSpace(cfg.DBSslMode)
	cfg.DBTimezone = strings.TrimSpace(cfg.DBTimezone)

	cfg.RedisHost = strings.TrimSpace(cfg.RedisHost)
	cfg.RedisPasswd = strings.TrimSpace(cfg.RedisPasswd)
}

func (cfg *EnvConfiguration) Validate() error {
	validate := validator.New()

	// Basic tag-based validation
	if err := validate.Struct(cfg); err != nil {
		return err
	}

	if cfg.LogToFile && cfg.LogFilePath == "" {
		return errors.New("LOG_FILE_PATH is required when LOG_TO_FILE is true")
	}

	return nil
}
