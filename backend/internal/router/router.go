package router

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/nawthtech/nawthtech/backend/internal/handlers"
	"github.com/nawthtech/nawthtech/backend/internal/middleware"
	"github.com/nawthtech/nawthtech/backend/internal/services"
)

// NewRouter creates and configures the main application router
func NewRouter(serviceContainer *services.ServiceContainer) *gin.Engine {
	// Set Gin mode
	gin.SetMode(gin.ReleaseMode)
	
	// Create router
	router := gin.New()
	
	// Global middleware
	router.Use(gin.Recovery()) // Recover from panics
	router.Use(middleware.Logger()) // Custom logging
	router.Use(CORSMiddleware()) // CORS middleware
	
	// Static files
	router.Static("/static", "./static")
	router.StaticFile("/favicon.ico", "./static/favicon.ico")
	
	// Create handlers
	authHandler := handlers.NewAuthHandler(serviceContainer.AuthService)
	userHandler := handlers.NewUserHandler(serviceContainer.UserService)
	serviceHandler := handlers.NewServiceHandler(serviceContainer.ServiceService)
	categoryHandler := handlers.NewCategoryHandler(serviceContainer.CategoryService)
	orderHandler := handlers.NewOrderHandler(serviceContainer.OrderService)
	paymentHandler := handlers.NewPaymentHandler(serviceContainer.PaymentService)
	
	// Create upload handler
	uploadHandler, err := handlers.NewUploadHandler()
	if err != nil {
		// Fallback to nil handler if upload service fails
		uploadHandler = nil
	}
	
	notificationHandler := handlers.NewNotificationHandler(serviceContainer.NotificationService)
	adminHandler := handlers.NewAdminHandler(serviceContainer.AdminService)
	
	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"time":    time.Now().UTC(),
			"service": "nawthtech-backend",
		})
	})
	
	// API routes
	api := router.Group("/api")
	{
		// API v1
		v1 := api.Group("/v1")
		{
			// Auth routes
			auth := v1.Group("/auth")
			{
				auth.POST("/register", authHandler.Register)
				auth.POST("/login", authHandler.Login)
				auth.POST("/logout", authHandler.Logout)
				auth.POST("/refresh-token", authHandler.RefreshToken)
				auth.POST("/forgot-password", authHandler.ForgotPassword)
				auth.POST("/reset-password", authHandler.ResetPassword)
				auth.POST("/verify-token", authHandler.VerifyToken)
			}
			
			// User routes
			users := v1.Group("/users")
			{
				users.GET("/profile", userHandler.GetProfile)
				users.PUT("/profile", userHandler.UpdateProfile)
				users.PUT("/change-password", userHandler.ChangePassword)
				users.GET("/stats", userHandler.GetUserStats)
			}
			
			// Service routes
			servicesGroup := v1.Group("/services")
			{
				servicesGroup.GET("/", serviceHandler.GetServices)
				servicesGroup.GET("/search", serviceHandler.SearchServices)
				servicesGroup.GET("/featured", serviceHandler.GetFeaturedServices)
				servicesGroup.GET("/categories", serviceHandler.GetCategories)
				servicesGroup.GET("/my-services", serviceHandler.GetMyServices)
				servicesGroup.GET("/seller-orders", serviceHandler.GetSellerOrders)
				servicesGroup.GET("/:id", serviceHandler.GetServiceByID)
				servicesGroup.POST("/", serviceHandler.CreateService)
				servicesGroup.PUT("/:id", serviceHandler.UpdateService)
				servicesGroup.DELETE("/:id", serviceHandler.DeleteService)
			}
			
			// Category routes
			categories := v1.Group("/categories")
			{
				categories.GET("/", categoryHandler.GetCategories)
				categories.GET("/:id", categoryHandler.GetCategoryByID)
				categories.POST("/", categoryHandler.CreateCategory)
				categories.PUT("/:id", categoryHandler.UpdateCategory)
				categories.DELETE("/:id", categoryHandler.DeleteCategory)
			}
			
			// Order routes
			orders := v1.Group("/orders")
			{
				orders.POST("/", orderHandler.CreateOrder)
				orders.GET("/my-orders", orderHandler.GetUserOrders)
				orders.GET("/seller-orders", orderHandler.GetSellerOrders)
				orders.GET("/:id", orderHandler.GetOrderByID)
				orders.PUT("/:id/status", orderHandler.UpdateOrderStatus)
				orders.PUT("/:id/cancel", orderHandler.CancelOrder)
			}
			
			// Payment routes
			payments := v1.Group("/payments")
			{
				payments.POST("/create-intent", paymentHandler.CreatePaymentIntent)
				payments.POST("/confirm", paymentHandler.ConfirmPayment)
				payments.GET("/history", paymentHandler.GetPaymentHistory)
				payments.POST("/stripe-webhook", paymentHandler.HandleStripeWebhook)
				payments.POST("/paypal-webhook", paymentHandler.HandlePayPalWebhook)
			}
			
			// Upload routes
			if uploadHandler != nil {
				uploads := v1.Group("/uploads")
				{
					uploads.POST("/image", uploadHandler.UploadImage)
					uploads.POST("/images", uploadHandler.UploadMultipleImages)
					uploads.GET("/my-images", uploadHandler.GetUserImages)
					uploads.GET("/image/:public_id", uploadHandler.GetImageInfo)
					uploads.DELETE("/image/:public_id", uploadHandler.DeleteImage)
					uploads.POST("/cloudinary-webhook", uploadHandler.HandleCloudinaryWebhook)
				}
			}
			
			// Notification routes
			notifications := v1.Group("/notifications")
			{
				notifications.GET("/", notificationHandler.GetUserNotifications)
				notifications.GET("/unread-count", notificationHandler.GetUnreadCount)
				notifications.PUT("/:id/read", notificationHandler.MarkAsRead)
				notifications.PUT("/read-all", notificationHandler.MarkAllAsRead)
			}
			
			// Admin routes
			admin := v1.Group("/admin")
			admin.Use(middleware.AdminRequired()) // Admin middleware
			{
				admin.GET("/dashboard", adminHandler.GetDashboard)
				admin.GET("/dashboard/stats", adminHandler.GetDashboardStats)
				admin.GET("/users", adminHandler.GetUsers)
				admin.GET("/orders", adminHandler.GetAllOrders)
				admin.GET("/system-logs", adminHandler.GetSystemLogs)
				admin.PUT("/users/:id/status", adminHandler.UpdateUserStatus)
			}
		}
		
		// API v2 (future version)
		v2 := api.Group("/v2")
		{
			v2.GET("/status", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"version": "2.0.0",
					"status":  "coming_soon",
				})
			})
		}
	}
	
	// Documentation
	router.GET("/docs", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"api": map[string]string{
				"v1": "/api/v1",
				"v2": "/api/v2",
			},
			"endpoints": map[string][]string{
				"auth": {
					"POST /api/v1/auth/register",
					"POST /api/v1/auth/login",
					"POST /api/v1/auth/logout",
				},
				"users": {
					"GET /api/v1/users/profile",
					"PUT /api/v1/users/profile",
				},
				"services": {
					"GET /api/v1/services/",
					"POST /api/v1/services/",
					"GET /api/v1/services/:id",
				},
				"health": {
					"GET /health",
				},
			},
		})
	})
	
	// 404 handler
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "endpoint_not_found",
			"message": "The requested endpoint does not exist",
			"path":    c.Request.URL.Path,
			"docs":    "/docs",
		})
	})
	
	return router
}

// CORSMiddleware provides CORS configuration
func CORSMiddleware() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:3000",
			"http://localhost:5173",
			"http://127.0.0.1:3000",
			"http://127.0.0.1:5173",
			"https://nawthtech.com",
			"https://*.nawthtech.com",
		},
		AllowMethods: []string{
			"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS", "HEAD",
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Content-Length",
			"Accept-Encoding",
			"X-CSRF-Token",
			"Authorization",
			"X-API-Key",
			"X-Requested-With",
			"Accept",
			"Cache-Control",
		},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}

// DevelopmentRouter creates a router with development settings
func DevelopmentRouter(serviceContainer *services.ServiceContainer) *gin.Engine {
	gin.SetMode(gin.DebugMode)
	router := NewRouter(serviceContainer)
	
	// Add development-only middleware
	router.Use(gin.Logger()) // Detailed logging in dev
	
	return router
}

// ProductionRouter creates a router with production settings
func ProductionRouter(serviceContainer *services.ServiceContainer) *gin.Engine {
	router := NewRouter(serviceContainer)
	
	// Add production-only middleware
	router.Use(middleware.RateLimitMiddleware()) // Rate limiting
	router.Use(middleware.SecurityMiddleware())  // Security headers
	
	return router
}