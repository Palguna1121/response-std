package services

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"response-std/core/models"

	"github.com/sirupsen/logrus"
)

type Logger struct {
	logger *logrus.Logger
}

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

	// Create logs directory if not exists
	logsDir := "logs"
	if _, err := os.Stat(logsDir); os.IsNotExist(err) {
		os.MkdirAll(logsDir, 0755)
	}

	// Set output to file in production
	if environment == "production" {
		logFile := filepath.Join(logsDir, fmt.Sprintf("api-service-%s.log", time.Now().Format("2006-01-02")))
		file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err == nil {
			logger.SetOutput(file)
		}
	}

	return &Logger{logger: logger}
}

func (l *Logger) LogRequest(requestLog *models.RequestLog) {
	l.logger.WithFields(logrus.Fields{
		"request_id":  requestLog.ID,
		"method":      requestLog.Method,
		"url":         requestLog.URL,
		"status_code": requestLog.StatusCode,
		"duration":    requestLog.Duration.String(),
		"error":       requestLog.Error,
		"retries":     requestLog.Retries,
	}).Info("API Request")
}

func (l *Logger) Info(message string, fields map[string]interface{}) {
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
	l.logger.WithFields(fields).Debug(message)
}

func (l *Logger) Warn(message string, fields map[string]interface{}) {
	l.logger.WithFields(fields).Warn(message)
}
