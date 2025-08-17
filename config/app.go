package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	APP_NAME       string        `mapstructure:"app_name" default:"response-std"`
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
	API_VERSION    []string      `mapstructure:"api_version" default:"v1"`
	API_BASE_URL   string        `mapstructure:"api_base_url" default:"http://localhost:5220/api/v1"`
	BASE_URL       string        `mapstructure:"base_url" default:"http://localhost:5220"`

	// Log Channel Configuration
	LogChannel         string `mapstructure:"log_channel" default:"file"`
	DiscordWebhookURL  string `mapstructure:"discord_webhook_url" default:""`
	DiscordMinLogLevel string `mapstructure:"discord_min_log_level" default:"error"`
	LogToFile          bool   `mapstructure:"log_to_file" default:"true"`
	LogDir             string `mapstructure:"log_dir" default:"logs"`
}

var ENV *Config

func InitConfig() {
	viper.SetConfigFile(".env")
	viper.BindEnv("app_name", "APP_NAME")
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
	viper.BindEnv("base_url", "BASE_URL")

	// Log Channel bindings
	viper.BindEnv("log_channel", "LOG_CHANNEL")
	viper.BindEnv("discord_webhook_url", "DISCORD_WEBHOOK_URL")
	viper.BindEnv("discord_min_log_level", "DISCORD_MIN_LOG_LEVEL")
	viper.BindEnv("log_to_file", "LOG_TO_FILE")
	viper.BindEnv("log_dir", "LOG_DIR")

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	rawVersions := viper.GetString("api_version")
	viper.Set("api_version", strings.Split(rawVersions, ","))

	if err := viper.Unmarshal(&ENV); err != nil {
		panic(fmt.Errorf("unable to decode into struct: %w", err))
	}
}

// Helper methods for log channel configuration
// IsFileLoggingEnabled checks if file logging is enabled
func (c *Config) IsFileLoggingEnabled() bool {
	// First check the config struct value
	if c.LogToFile {
		return true
	}

	// Then check environment variables
	logChannel := strings.ToLower(c.LogChannel)
	logToFile := strings.ToLower(os.Getenv("LOG_TO_FILE"))

	// File logging is enabled if:
	// 1. LOG_CHANNEL contains "file" or "both"
	// 2. LOG_TO_FILE is "true"
	return strings.Contains(logChannel, "file") ||
		strings.Contains(logChannel, "both") ||
		logToFile == "true"
}

// IsDiscordLoggingEnabled checks if Discord logging is enabled
func (c *Config) IsDiscordLoggingEnabled() bool {
	// First check the config struct value
	if c.DiscordWebhookURL == "" {
		return false
	}

	// Then check environment variables
	logChannel := strings.ToLower(c.LogChannel)

	// Discord logging is enabled if:
	// 1. LOG_CHANNEL contains "discord" or "both"
	// 2. DISCORD_WEBHOOK_URL is not empty
	return (strings.Contains(logChannel, "discord") ||
		strings.Contains(logChannel, "both")) &&
		c.DiscordWebhookURL != ""
}

// ShouldLogToDiscord checks if a specific log level should be sent to Discord
func (c *Config) ShouldLogToDiscord(level string) bool {
	if !c.IsDiscordLoggingEnabled() {
		return false
	}

	// Parse minimum level for Discord
	minLevel := strings.ToLower(c.DiscordMinLogLevel)
	if minLevel == "" {
		minLevel = "error" // default
	}

	currentLevel := strings.ToLower(level)

	// Level hierarchy: debug < info < warn < error < critical/fatal
	levelPriority := map[string]int{
		"debug":    1,
		"info":     2,
		"warn":     3,
		"warning":  3,
		"error":    4,
		"critical": 5,
		"fatal":    5,
		"panic":    5,
	}

	minPriority := levelPriority[minLevel]
	currentPriority := levelPriority[currentLevel]

	return currentPriority >= minPriority
}

// GetLogLevel returns the general log level for file logging
func (c *Config) GetLogLevel() string {
	if c.LogLevel != "" {
		return c.LogLevel
	}

	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		return "info" // default
	}
	return logLevel
}

// GetDiscordMinLogLevel returns the minimum log level for Discord notifications
func (c *Config) GetDiscordMinLogLevel() string {
	if c.DiscordMinLogLevel != "" {
		return c.DiscordMinLogLevel
	}

	level := os.Getenv("DISCORD_MIN_LOG_LEVEL")
	if level == "" {
		return "error" // default
	}
	return level
}
