package routes

import "github.com/gin-gonic/gin"

func SetupRoutes(r *gin.Engine) {
	// Semua routing v2
	api := r.Group("/api/v2")
	api.GET("/hello", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Hello from V2"})
	})
}
