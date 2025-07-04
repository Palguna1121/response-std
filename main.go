package main

import (
	"response-std/config"
	"response-std/core/router"
	"response-std/core/services"
	"response-std/core/services/hooks"

	"github.com/gin-gonic/gin"
)

func main() {
	config.InitConfig()
	config.LoadDB()
	// Initialize logger
	services.AppLogger = services.NewLogger(config.ENV.LogLevel, config.ENV.Environment)

	var log = services.AppLogger

	r := gin.Default()
	for _, version := range config.ENV.API_VERSION {
		setup, ok := router.RouteRegistry[version]
		if !ok {
			panic("Unsupported API version: " + version)
		}
		setup(r)
	}

	// Initialize Gin router
	if config.ENV.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	log.Info("Starting API Service", map[string]interface{}{
		"APP_PORT":       config.ENV.APP_PORT,
		"environment":    config.ENV.Environment,
		"log_level":      config.ENV.LogLevel,
		"version_active": config.ENV.API_VERSION,
	})
	hooks.SendDiscordMessage(
		config.ENV.DiscordWebhookURL,
		config.ENV.APP_NAME,
		config.ENV.LogLevel,
		"Service started!",
		map[string]interface{}{
			"APP_PORT":       config.ENV.APP_PORT,
			"environment":    config.ENV.Environment,
			"version_active": config.ENV.API_VERSION,
		},
	)

	// Send logs to Discord
	// dcLogs()

	err := r.Run(":" + config.ENV.APP_PORT)
	if err != nil {
		panic("Failed to run server: " + err.Error())
	}
}

// func dcLogs() {
// 	colors := []struct {
// 		level string
// 		color int
// 	}{
// 		{"fatal", 0xFF0000},
// 		{"panic", 0xFF0000},
// 		{"error", 0xFF6B6B},
// 		{"warn", 0xFFB347},
// 		{"warning", 0xFFB347},
// 		{"info", 0x4ECDC4},
// 		{"debug", 0x95A5A6},
// 	}

// 	for _, c := range colors {
// 		errdc := hooks.SendDiscordMessage(
// 			config.ENV.DiscordWebhookURL,
// 			"MyApp",
// 			c.level,
// 			"🚨 Terjadi kesalahan saat proses pembayaran",
// 			map[string]interface{}{
// 				"user_id": 1023,
// 				"order":   "#A202406",
// 				"error":   "Timeout saat koneksi ke server pembayaran",
// 			},
// 		)

// 		if errdc != nil {
// 			// fallback jika ingin log ke file atau stdout
// 			println("Gagal mengirim log ke Discord dengan level", c.level, ":", errdc.Error())
// 		}
// 	}
// }
