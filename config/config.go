package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	APP_PORT       string        `mapstructure:"app_port" default:"5220"`
	DB_HOST        string        `mapstructure:"db_host" default:"localhost"`
	DB_PORT        string        `mapstructure:"db_port" default:"3306"`
	DB_USER        string        `mapstructure:"db_user" default:"root"`
	DB_PASSWORD    string        `mapstructure:"db_password" default:""`
	DB_NAME        string        `mapstructure:"db_name" default:"golang_api"`
	JWT_SECRET     string        `mapstructure:"jwt_secret" default:"supersecretkey"`
	RequestTimeout time.Duration `mapstructure:"request_timeout" default:"5s"`
	MaxRetries     int           `mapstructure:"max_retries" default:"3"`
	RetryDelay     time.Duration `mapstructure:"retry_delay" default:"200ms"`
	EnableLogging  bool          `mapstructure:"enable_logging" default:"true"`
	LogLevel       string        `mapstructure:"log_level" default:"info"`
	Environment    string        `mapstructure:"environment" default:"development"`
	API_VERSION    string        `mapstructure:"api_version" default:"v1"`
	API_BASE_URL   string        `mapstructure:"api_base_url" default:"http://localhost:5220/api/v1"`
}

var ENV *Config

func InitConfig() {
	viper.SetConfigFile(".env")
	viper.BindEnv("port", "APP_PORT")
	viper.BindEnv("db_host", "DB_HOST")
	viper.BindEnv("db_port", "DB_PORT")
	viper.BindEnv("db_user", "DB_USER")
	viper.BindEnv("db_password", "DB_PASSWORD")
	viper.BindEnv("db_name", "DB_NAME")
	viper.BindEnv("jwt_secret", "JWT_SECRET")
	viper.BindEnv("request_timeout", "REQUEST_TIMEOUT")
	viper.BindEnv("max_retries", "MAX_RETRIES")
	viper.BindEnv("retry_delay", "RETRY_DELAY")
	viper.BindEnv("enable_logging", "ENABLE_LOGGING")
	viper.BindEnv("log_level", "LOG_LEVEL")
	viper.BindEnv("environment", "ENVIRONMENT")
	viper.BindEnv("api_version", "API_VERSION")
	viper.BindEnv("api_base_url", "API_BASE_URL")

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	if err := viper.Unmarshal(&ENV); err != nil {
		panic(fmt.Errorf("unable to decode into struct: %w", err))
	}
}
