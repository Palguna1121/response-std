package services

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"response-std/config"
	"response-std/core/models/requests"
	"response-std/core/services/hooks"

	"github.com/sirupsen/logrus"
)

type Logger struct {
	logger *logrus.Logger
}

var AppLogger *Logger

func NewLogger(logLevel, environment string) *Logger {
	logger := logrus.New()

	// Set log level
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		level = logrus.InfoLevel
	}
	logger.SetLevel(level)

	// Set formatter based on environment
	if environment == "production" {
		logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: time.RFC3339,
		})
	} else {
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: time.RFC3339,
		})
	}

	// Configure output based on LOG_CHANNEL setting
	outputs := []io.Writer{}

	// File logging
	if config.ENV.IsFileLoggingEnabled() {
		fileOutput := setupFileLogging(environment)
		if fileOutput != nil {
			outputs = append(outputs, fileOutput)
		}
	}

	// Console output (always for development, optional for production)
	if environment != "production" || len(outputs) == 0 {
		outputs = append(outputs, os.Stdout)
	}

	// Set multiple outputs if needed
	if len(outputs) > 1 {
		logger.SetOutput(io.MultiWriter(outputs...))
	} else if len(outputs) == 1 {
		logger.SetOutput(outputs[0])
	}

	// Add Discord hook if enabled
	if config.ENV.IsDiscordLoggingEnabled() {
		minLevel := hooks.ParseLogLevel(config.ENV.DiscordMinLogLevel)
		discordHook := hooks.NewDiscordHook(
			config.ENV.DiscordWebhookURL,
			config.ENV.APP_NAME,
			minLevel,
		)
		logger.AddHook(discordHook)
	}

	return &Logger{logger: logger}
}

func setupFileLogging(environment string) io.Writer {
	// Create logs directory if not exists
	logsDir := config.ENV.LogDir
	if logsDir == "" {
		logsDir = "logs"
	}

	if _, err := os.Stat(logsDir); os.IsNotExist(err) {
		if err := os.MkdirAll(logsDir, 0755); err != nil {
			return nil
		}
	}

	// Determine log file name based on environment
	var logFileName string
	switch environment {
	case "production":
		logFileName = fmt.Sprintf("%s-api-service.log", time.Now().Format("2006-01-02"))
	case "development":
		logFileName = fmt.Sprintf("%s-api-development.log", time.Now().Format("2006-01-02"))
	default:
		logFileName = fmt.Sprintf("%s-api-%s.log", time.Now().Format("2006-01-02"), environment)
	}

	logFile := filepath.Join(logsDir, logFileName)
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil
	}

	return file
}

func (l *Logger) LogRequest(requestLog *requests.RequestLog) {
	fields := logrus.Fields{
		"request_id":  requestLog.ID,
		"method":      requestLog.Method,
		"url":         requestLog.URL,
		"status_code": requestLog.StatusCode,
		"duration":    requestLog.Duration.String(),
		"error":       requestLog.Error,
		"retries":     requestLog.Retries,
	}

	// Log to Discord only for errors (status code >= 400)
	if config.ENV.ShouldLogToDiscord("info") && requestLog.StatusCode >= 400 {
		l.logger.WithFields(fields).Error("API Request Error")
	} else {
		l.logger.WithFields(fields).Info("API Request")
	}
}

func (l *Logger) Info(message string, fields map[string]interface{}) {
	if fields == nil {
		fields = make(map[string]interface{})
	}
	l.logger.WithFields(fields).Info(message)
}

func (l *Logger) Error(message string, err error, fields map[string]interface{}) {
	if fields == nil {
		fields = make(map[string]interface{})
	}
	if err != nil {
		fields["error"] = err.Error()
	}
	l.logger.WithFields(fields).Error(message)
}

func (l *Logger) Debug(message string, fields map[string]interface{}) {
	if fields == nil {
		fields = make(map[string]interface{})
	}
	l.logger.WithFields(fields).Debug(message)
}

func (l *Logger) Warn(message string, fields map[string]interface{}) {
	if fields == nil {
		fields = make(map[string]interface{})
	}
	l.logger.WithFields(fields).Warn(message)
}

func (l *Logger) Critical(message string, err error, fields map[string]interface{}) {
	if fields == nil {
		fields = make(map[string]interface{})
	}
	if err != nil {
		fields["error"] = err.Error()
	}
	l.logger.WithFields(fields).Fatal(message)
}