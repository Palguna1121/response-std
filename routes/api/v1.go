package api

import (
	"response-std/app/http/controllers"
	"response-std/app/http/middleware"
	"response-std/app/pkg/permissions"
	"response-std/config"
	"response-std/libs/external/handlers"
	"response-std/libs/external/services"

	"github.com/Palguna1121/goupload"

	"github.com/gin-gonic/gin"
)

func SetupRoutesv1(r *gin.Engine) {
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

	// register upload routes
	registerUploadRoutes(r)

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
			auth.POST("/register", authController.Register(config.DB, permissions.NewSpatie(config.DB)))
		}

		// Protected routes (require authentication)
		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware(config.DB))
		{
			// Auth endpoints
			protected.POST("/auth/logout", authController.Logout(config.DB))
			protected.POST("/auth/refresh", authController.RefreshToken(config.DB))
			protected.GET("/auth/me", func(c *gin.Context) {
				authController.Me(c, permissions.NewSpatie(config.DB))
			})

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

				// Example API request validation
				// admin.Use(mware.ValidationMiddleware())
			}
		}
	}
}

func registerUploadRoutes(r *gin.Engine) {
	// Konfigurasi uploader
	uploader := goupload.NewImageUploader(goupload.UploadConfig{
		StoragePath:       "storage/app/public/uploads/images",
		BaseURL:           config.ENV.BASE_URL,
		EnableTimestamp:   true,
		CreateDateDir:     true,
		AllowedExtensions: []string{"jpg", "jpeg", "png", "webp"},
		MaxFileSize:       2 << 20, // 2 MB
	})

	// Daftarkan route ke root (global)
	uploader.RegisterRoutes(r, "/upload")
	uploader.ServeStaticFiles(r, "/storage")
}
