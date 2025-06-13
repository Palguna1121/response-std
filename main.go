package main

import (
	"response-std/config"
	"response-std/core/router"
	"response-std/core/services"

	"github.com/gin-gonic/gin"
)

func main() {
	config.InitConfig()
	config.LoadDB()
	// Initialize logger
	logger := services.NewLogger(config.ENV.LogLevel, config.ENV.Environment)

	r := gin.Default()
	setup, ok := router.RouteRegistry[config.ENV.API_VERSION]
	if !ok {
		panic("Unsupported API version: " + config.ENV.API_VERSION)
	}
	setup(r)

	// Initialize Gin router
	if config.ENV.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	logger.Info("Starting API Service", map[string]interface{}{
		"APP_PORT":    config.ENV.APP_PORT,
		"environment": config.ENV.Environment,
		"log_level":   config.ENV.LogLevel,
	})

	err := r.Run(":" + config.ENV.APP_PORT)
	if err != nil {
		panic("Failed to run server: " + err.Error())
	}
}
