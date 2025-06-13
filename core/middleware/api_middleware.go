package middleware

import (
	"fmt"
	"time"

	"response-std/core/response"
	"response-std/core/services"

	// "github.com/gin-contrib/cors"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// ---------------------------
// CORS MIDDLEWARE
// ---------------------------
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

// ---------------------------
// LOGGING MIDDLEWARE
// ---------------------------
func LoggingMiddleware(logger *services.Logger) gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	})
}

// ---------------------------
// ERROR HANDLING MIDDLEWARE
// ---------------------------
func ErrorHandlingMiddleware(logger *services.Logger) gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			logger.Error("Panic recovered: "+err, nil, map[string]interface{}{})
		}
		response.InternalServerError(c, "Internal server error occurred")
	})
}

// ---------------------------
// RATE LIMITING MIDDLEWARE
// ---------------------------
var limiter = rate.NewLimiter(10, 20) // 10 requests per second, burst of 20

func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !limiter.Allow() {
			response.TooManyRequests(c, "Rate limit exceeded")
			c.Abort()
			return
		}
		c.Next()
	}
}
