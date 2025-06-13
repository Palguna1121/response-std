package response

import (
	"github.com/gin-gonic/gin"
)

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
func Error(c *gin.Context, code int, message string) {
	respond(c, code, message, nil)
}

func BadRequest(c *gin.Context, message string) {
	Error(c, 400, message)
}

func Unauthorized(c *gin.Context, message string) {
	Error(c, 401, message)
}

func Forbidden(c *gin.Context, message string) {
	Error(c, 403, message)
}

func NotFound(c *gin.Context, message string) {
	Error(c, 404, message)
}

func UnprocessableEntity(c *gin.Context, message string) {
	Error(c, 422, message)
}

func Conflict(c *gin.Context, message string) {
	Error(c, 409, message)
}

func InternalServerError(c *gin.Context, message string) {
	Error(c, 500, message)
}

func ServiceUnavailable(c *gin.Context, message string) {
	Error(c, 503, message)
}

func TooManyRequests(c *gin.Context, message string) {
	Error(c, 429, message)
}

func Gone(c *gin.Context, message string) {
	Error(c, 410, message)
}

func PreconditionFailed(c *gin.Context, message string) {
	Error(c, 412, message)
}

func RequestTimeout(c *gin.Context, message string) {
	Error(c, 408, message)
}
