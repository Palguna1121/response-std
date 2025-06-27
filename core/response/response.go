package response

import (
	"response-std/config"
	"response-std/core/services"

	"github.com/gin-gonic/gin"
)

var log = services.AppLogger

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

// Error responses
func Error(c *gin.Context, code int, message string, err error, logPrefix string) {
	if log != nil {
		APP_NAME := config.ENV.APP_NAME
		Prefix := logPrefix
		mssg := "" + APP_NAME + " " + Prefix + " : " + message

		log.Error(mssg, err, map[string]interface{}{
			"status_code": code,
			"request":     c.Request.URL.Path,
			"method":      c.Request.Method,
			"client_ip":   c.ClientIP(),
		})
	}
	respond(c, code, message, nil)
}

func BadRequest(c *gin.Context, message string, err error, logPrefix ...string) {
	Error(c, 400, message, err, getLogPrefix(logPrefix, "Bad Request"))
}

func Unauthorized(c *gin.Context, message string, err error, logPrefix ...string) {
	Error(c, 401, message, err, getLogPrefix(logPrefix, "Unauthorized"))
}

func Forbidden(c *gin.Context, message string, err error, logPrefix ...string) {
	Error(c, 403, message, err, getLogPrefix(logPrefix, "Forbidden"))
}

func NotFound(c *gin.Context, message string, err error, logPrefix ...string) {
	Error(c, 404, message, err, getLogPrefix(logPrefix, "Not Found"))
}

func UnprocessableEntity(c *gin.Context, message string, err error, logPrefix ...string) {
	Error(c, 422, message, err, getLogPrefix(logPrefix, "Unprocessable Entity"))
}

func Conflict(c *gin.Context, message string, err error, logPrefix ...string) {
	Error(c, 409, message, err, getLogPrefix(logPrefix, "Conflict"))
}

func InternalServerError(c *gin.Context, message string, err error, logPrefix ...string) {
	Error(c, 500, message, err, getLogPrefix(logPrefix, "Internal Server Error"))
}

func ServiceUnavailable(c *gin.Context, message string, err error, logPrefix ...string) {
	Error(c, 503, message, err, getLogPrefix(logPrefix, "Service Unavailable"))
}

func TooManyRequests(c *gin.Context, message string, err error, logPrefix ...string) {
	Error(c, 429, message, err, getLogPrefix(logPrefix, "Too Many Requests"))
}

func Gone(c *gin.Context, message string, err error, logPrefix ...string) {
	Error(c, 410, message, err, getLogPrefix(logPrefix, "Gone"))
}

func PreconditionFailed(c *gin.Context, message string, err error, logPrefix ...string) {
	Error(c, 412, message, err, getLogPrefix(logPrefix, "Precondition Failed"))
}

func RequestTimeout(c *gin.Context, message string, err error, logPrefix ...string) {
	Error(c, 408, message, err, getLogPrefix(logPrefix, "Request Timeout"))
}

func getLogPrefix(logPrefix []string, defaultValue string) string {
	if len(logPrefix) > 0 {
		return logPrefix[0]
	}
	return defaultValue
}
