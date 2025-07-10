package web

import (
	"response-std/app/http/controllers"
	"response-std/app/pkg/permissions"
	"response-std/config"

	"github.com/gin-gonic/gin"
)

func SetupWebRoutes(r *gin.Engine) {
	//controllers
	userController := controllers.NewUserController(config.DB, permissions.NewSpatie(config.DB))

	// Semua routing v2
	api := r.Group("/api/web")
	api.GET("/hello", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Hello from web API!"})
	})

	api.GET("/error-debug", userController.ErrorDebug)
	api.GET("/error-info", userController.ErrorInfo)
	api.GET("/error-warn", userController.ErrorWarn)
	api.GET("/error-error", userController.ErrorError)
	api.GET("/error-critical", userController.ErrorCritical)

	user := api.Group("/users")
	{
		user.GET("/", userController.ListUser)
		user.GET("/:id", userController.GetUserByID)
		user.POST("/", userController.CreateUser)
		user.PUT("/:id/update", userController.UpdateUser)
		user.DELETE("/:id/delete", userController.DeleteUser)

	}
}
