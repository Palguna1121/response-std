package routes

import (
	"response-std/controllers"
	"response-std/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	authController := controllers.NewAuthController()

	api := r.Group("/api/v1")
	{
		api.POST("/login", authController.Login)
		api.POST("/register", authController.Register)

		auth := api.Group("/")
		auth.Use(middleware.AuthMiddleware())
		{
			auth.GET("/me", authController.Me)

			// group and check role with middleware
			// adminGroup := auth.Group("/admin")
			// adminGroup.Use(middleware.RoleMiddleware("admin"))
			// {
			// 	adminGroup.GET("/users", adminController.GetUsers)
			// 	adminGroup.POST("/users", adminController.CreateUser)
			// }
		}
	}
}
