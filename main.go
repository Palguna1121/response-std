package main

import (
	"response-std/config"
	"response-std/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	config.InitConfig()
	config.LoadDB()

	r := gin.Default()
	routes.SetupRoutes(r)

	err := r.Run(":" + config.ENV.APP_PORT)
	if err != nil {
		panic("Failed to run server: " + err.Error())
	}
}
