package main

import (
	"response-std/config"

	"github.com/gin-gonic/gin"
)

func main() {
	config.InitConfig()
	config.LoadDB()

	router := gin.Default()
	api := router.Group("/api/v2")
	{
		api.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		})
	}

	router.Run(":" + config.ENV.APP_PORT) // listen and serve on port in .env file
}
