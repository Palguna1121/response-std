package routes

import "github.com/gin-gonic/gin"

func SetupWebRoutes(r *gin.Engine) {
	// Semua routing v2
	api := r.Group("/api/web")
	api.GET("/hello", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Hello from web API!"})
	})
}
