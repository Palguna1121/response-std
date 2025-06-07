package routes

import (
	"response-std/config"
	"response-std/controllers"
	"response-std/handlers"
	"response-std/middleware"
	"response-std/services"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	authController := controllers.NewAuthController()

	// Initialize logger
	logger := services.NewLogger(config.ENV.LogLevel, config.ENV.Environment)

	// Initialize API client
	apiClient := services.NewAPIClient(config.ENV, logger)

	// Initialize handlers
	apiHandler := handlers.NewAPIHandler(apiClient, logger)

	r.Use(gin.Recovery())
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.LoggingMiddleware(logger))
	r.Use(middleware.ErrorHandlingMiddleware(logger))
	r.Use(middleware.RateLimitMiddleware())

	//routes
	api := r.Group("/api/v1")
	{
		// Health check
		api.GET("/health", apiHandler.HealthCheck)

		// Execute API request
		api.POST("/login", authController.Login(config.DB))
		api.POST("/register", authController.Register(config.DB))

		auth := api.Group("/")
		auth.Use(middleware.AuthMiddleware(config.DB))
		{
			auth.GET("/me", authController.Me)

			// group and check role with middleware
			// adminGroup := auth.Group("/admin")
			// adminGroup.Use(middleware.RoleMiddleware("admin"))
			// {
			// 	adminGroup.GET("/users", adminController.GetUsers)
			// 	adminGroup.POST("/users", adminController.CreateUser)
			// }
			auth.POST("/logout", authController.Logout(config.DB))
		}
	}

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"service": "API Service Starter",
			"version": "1.0.0",
			"status":  "running",
			"docs":    "/api/v1/health",
		})
	})
}
