package routes

import (
	"response-std/config"
	"response-std/core/handlers"
	"response-std/core/middleware"
	"response-std/core/services"
	"response-std/v1/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	// Initialize controllers
	authController := controllers.NewAuthController()

	// Initialize services
	logger := services.NewLogger(config.ENV.LogLevel, config.ENV.Environment)
	apiClient := services.NewAPIClient(config.ENV, logger)
	apiHandler := handlers.NewAPIHandler(apiClient, logger)

	// Global middlewares
	r.Use(gin.Recovery())
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.LoggingMiddleware(logger))
	r.Use(middleware.ErrorHandlingMiddleware(logger))
	r.Use(middleware.RateLimitMiddleware())

	// Root endpoint
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"service": "API Service Starter",
			"version": "1.0.0",
			"status":  "running",
			"docs":    "/api/v1/health",
		})
	})

	// API routes
	api := r.Group("/api/v1")
	{
		api.GET("/hello", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Hello from V1"})
		})
		// Health check
		api.GET("/health", apiHandler.HealthCheck)

		// Authentication routes (public)
		auth := api.Group("/auth")
		{
			auth.POST("/login", authController.Login(config.DB))
			auth.POST("/register", authController.Register(config.DB))
		}

		// Protected routes (require authentication)
		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware(config.DB))
		{
			// Auth endpoints
			protected.POST("/auth/logout", authController.Logout(config.DB))
			protected.POST("/auth/refresh", authController.RefreshToken(config.DB))
			protected.GET("/auth/me", authController.Me)

			// Admin routes (require admin role)
			admin := protected.Group("/admin")
			admin.Use(middleware.RoleMiddleware("admin"))
			{
				// Example admin routes
				// admin.GET("/dashboard", adminController.Dashboard)
				// admin.GET("/users", adminController.GetAllUsers)
				// admin.POST("/users", adminController.CreateUser)
				// admin.GET("/users/:id", adminController.GetUserByID)
				// admin.PUT("/users/:id", adminController.UpdateUser)
				// admin.DELETE("/users/:id", adminController.DeleteUser)
			}
		}
	}
}
