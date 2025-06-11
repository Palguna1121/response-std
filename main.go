package main

import (
	"github.com/Palguna1121/response-std/config"
	"github.com/Palguna1121/response-std/routes"
	"github.com/Palguna1121/response-std/services"

	"github.com/gin-gonic/gin"
)

func main() {
	config.InitConfig()
	config.LoadDB()
	// Initialize logger
	logger := services.NewLogger(config.ENV.LogLevel, config.ENV.Environment)

	r := gin.Default()
	routes.SetupRoutes(r)

	// Initialize Gin router
	if config.ENV.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	logger.Info("Starting API Service Starter", map[string]interface{}{
		"APP_PORT":    config.ENV.APP_PORT,
		"environment": config.ENV.Environment,
		"log_level":   config.ENV.LogLevel,
	})

	err := r.Run(":" + config.ENV.APP_PORT)
	if err != nil {
		panic("Failed to run server: " + err.Error())
	}
}
