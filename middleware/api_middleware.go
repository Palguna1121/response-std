package middleware

import (
	"time"

	"response-std/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization", "X-Requested-With"}
	config.ExposeHeaders = []string{"Content-Length"}
	config.AllowCredentials = true
	config.MaxAge = 12 * time.Hour

	return cors.New(config)
}

func LoggingMiddleware(logger *services.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		c.Next()

		duration := time.Since(start)
		statusCode := c.Writer.Status()

		logger.Info("HTTP Request", map[string]interface{}{
			"method":      method,
			"path":        path,
			"status_code": statusCode,
			"duration":    duration.String(),
			"client_ip":   c.ClientIP(),
			"user_agent":  c.Request.UserAgent(),
		})
	}
}

func ErrorHandlingMiddleware(logger *services.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			logger.Error("Request Error", err, map[string]interface{}{
				"path":   c.Request.URL.Path,
				"method": c.Request.Method,
			})

			c.JSON(500, gin.H{
				"success":   false,
				"error":     err.Error(),
				"message":   "Internal server error",
				"timestamp": time.Now(),
			})
		}
	}
}

func RateLimitMiddleware() gin.HandlerFunc {
	// Simple rate limiting
	return func(c *gin.Context) {
		// Implementasi rate limiting jika diperlukan
		// For now, just pass through
		c.Next()
	}
}
