package response

import (
	"response-std/config"
	"response-std/libs/external/services"

	"github.com/gin-gonic/gin"
)

var log *services.Logger

func InitLogger() {
	log = services.NewLogger(config.ENV.LogLevel, config.ENV.Environment)
}

type Response struct {
	Status  string `json:"status"`
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func respond(c *gin.Context, statusCode int, message string, data any) {
	var status string
	if statusCode >= 200 && statusCode <= 299 {
		status = "success"
	} else {
		status = "error"
	}
	response := Response{
		Status:  status,
		Code:    statusCode,
		Message: message,
		Data:    data,
	}

	c.JSON(statusCode, response)
}

func Success(c *gin.Context, message string, data any) {
	respond(c, 200, message, data)
}

func Created(c *gin.Context, message string, data any) {
	respond(c, 201, message, data)
}

func Accepted(c *gin.Context, message string, data any) {
	respond(c, 202, message, data)
}

func NoContent(c *gin.Context) {
	c.Status(204)
}

// Error responses with log level checking
func Error(c *gin.Context, code int, message string, err error, logPrefix string, level ...string) {
	if log != nil {
		APP_NAME := config.ENV.APP_NAME
		Prefix := logPrefix

		// Determine log level (default to "error")
		logLevel := "error"
		if len(level) > 0 && level[0] != "" {
			logLevel = level[0]
		}

		// Add level prefix for debug
		if logLevel == "debug" {
			Prefix = "DEBUG: " + Prefix
		}

		mssg := "" + APP_NAME + " " + Prefix + " : " + message

		fields := map[string]interface{}{
			"status_code": code,
			"request":     c.Request.URL.Path,
			"method":      c.Request.Method,
			"client_ip":   c.ClientIP(),
		}

		// Log based on specified level
		switch logLevel {
		case "debug":
			log.Debug(mssg, fields)
		case "info":
			log.Info(mssg, fields)
		case "warn":
			log.Warn(mssg, fields)
		case "error":
			log.Error(mssg, err, fields)
		case "critical", "fatal":
			log.Critical(mssg, err, fields)
		default:
			// Default to error level
			log.Error(mssg, err, fields)
		}
	}
	respond(c, code, message, nil)
}

// Client Error Responses (4xx) - typically logged as warnings
func BadRequest(c *gin.Context, message string, err error, logPrefix ...string) {
	Error(c, 400, message, err, getLogPrefix(logPrefix, "Bad Request"), "critical")
}

func Unauthorized(c *gin.Context, message string, err error, logPrefix ...string) {
	Error(c, 401, message, err, getLogPrefix(logPrefix, "Unauthorized"), "error")
}

func Forbidden(c *gin.Context, message string, err error, logPrefix ...string) {
	Error(c, 403, message, err, getLogPrefix(logPrefix, "Forbidden"), "critical")
}

func NotFound(c *gin.Context, message string, err error, logPrefix ...string) {
	Error(c, 404, message, err, getLogPrefix(logPrefix, "Not Found"), "error")
}

func UnprocessableEntity(c *gin.Context, message string, err error, logPrefix ...string) {
	Error(c, 422, message, err, getLogPrefix(logPrefix, "Unprocessable Entity"), "warn")
}

func UnprocessableValidation(c *gin.Context, message string, err error, errInterface map[string]interface{}, logPrefix ...string) {
	validationErrorRespond(c, message, errInterface)
	if log != nil {
		log.Warn(message, map[string]interface{}{
			"error":      err,
			"message":    errInterface,
			"request":    c.Request.URL.Path,
			"method":     c.Request.Method,
			"client_ip":  c.ClientIP(),
			"log_prefix": getLogPrefix(logPrefix, "Unprocessable Entity"),
		})
	}
}

func Conflict(c *gin.Context, message string, err error, logPrefix ...string) {
	Error(c, 409, message, err, getLogPrefix(logPrefix, "Conflict"), "warn")
}

func Gone(c *gin.Context, message string, err error, logPrefix ...string) {
	Error(c, 410, message, err, getLogPrefix(logPrefix, "Gone"), "info")
}

func PreconditionFailed(c *gin.Context, message string, err error, logPrefix ...string) {
	Error(c, 412, message, err, getLogPrefix(logPrefix, "Precondition Failed"), "warn")
}

func RequestTimeout(c *gin.Context, message string, err error, logPrefix ...string) {
	Error(c, 408, message, err, getLogPrefix(logPrefix, "Request Timeout"), "warn")
}

func TooManyRequests(c *gin.Context, message string, err error, logPrefix ...string) {
	Error(c, 429, message, err, getLogPrefix(logPrefix, "Too Many Requests"), "warn")
}

// Server Error Responses (5xx) - logged as errors or critical
func InternalServerError(c *gin.Context, message string, err error, logPrefix ...string) {
	Error(c, 500, message, err, getLogPrefix(logPrefix, "Internal Server Error"), "critical")
}

func ServiceUnavailable(c *gin.Context, message string, err error, logPrefix ...string) {
	Error(c, 503, message, err, getLogPrefix(logPrefix, "Service Unavailable"), "critical")
}

func getLogPrefix(logPrefix []string, defaultValue string) string {
	if len(logPrefix) > 0 {
		return logPrefix[0]
	}
	return defaultValue
}

// ValidationErrorResponse is a custom error response for validation errors
type ErrorResponse struct {
	Status  string                 `json:"status"`
	Code    int                    `json:"code"`
	Message string                 `json:"message"`
	Error   map[string]interface{} `json:"error"`
}

func validationErrorRespond(c *gin.Context, message string, err map[string]interface{}) {
	response := ErrorResponse{
		Status:  "error",
		Code:    422,
		Message: message,
		Error:   err,
	}

	c.JSON(422, response)
}
