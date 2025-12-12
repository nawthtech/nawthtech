package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/nawthtech/nawthtech/backend/internal/config"
	"github.com/nawthtech/nawthtech/backend/internal/middleware"
)

// RegisterAllRoutes تسجيل جميع المسارات
func RegisterAllRoutes(app *gin.Engine, cfg *config.Config, hc *HandlerContainer) {
	// ==================== Public Routes ====================
	api := app.Group("/api/v1")
	
	// Authentication
	auth := api.Group("/auth")
	{
		if hc.Auth != nil {
			auth.POST("/register", hc.Auth.Register)
			auth.POST("/login", hc.Auth.Login)
			auth.POST("/logout", hc.Auth.Logout)
			auth.POST("/refresh", hc.Auth.RefreshToken)
		}
	}
	
	// Health endpoints
	health := app.Group("/health")
	{
		if hc.Health != nil {
			health.GET("", hc.Health.CheckHealth)
			health.GET("/live", hc.Health.HealthCheck)
			health.GET("/ready", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ready"}) })
		} else {
			health.GET("/live", func(c *gin.Context) { c.JSON(200, gin.H{"status": "live"}) })
			health.GET("/ready", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ready"}) })
		}
	}
	
	// AI endpoints
 ai := api.Group("/ai")
 {
	if hc.AI != nil {
		ai.GET("/capabilities", hc.AI.GetAICapabilitiesHandler)
		ai.POST("/generate", hc.AI.GenerateContentHandler)
		ai.POST("/translate", hc.AI.TranslateTextHandler)
		ai.POST("/summarize", hc.AI.SummarizeTextHandler)
		ai.POST("/analyze-image", hc.AI.AnalyzeImageHandler)
		ai.POST("/analyze-text", hc.AI.AnalyzeTextHandler)
		ai.POST("/generate-video", hc.AI.GenerateVideoHandler)
		ai.GET("/video-status/:id", hc.AI.CheckVideoStatusHandler)
		ai.GET("/providers", func(c *gin.Context) {
			if hc.AI != nil && hc.AI.aiClient != nil {
				providers := hc.AI.aiClient.GetAvailableProviders()
				c.JSON(200, gin.H{"providers": providers})
			} else {
				c.JSON(200, gin.H{"providers": {}})
			}
		})
		ai.GET("/usage-stats", func(c *gin.Context) {
			if hc.AI != nil && hc.AI.aiClient != nil {
				stats := hc.AI.aiClient.GetUsageStatistics()
				c.JSON(200, gin.H{"stats": stats})
			} else {
				c.JSON(200, gin.H{"stats": {}})
			}
		})
	} else {
		ai.GET("/capabilities", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "AI service not available"})
		})
	}
}
	
	// Email endpoints (public for setup)
	email := api.Group("/email")
	{
		if hc.Email != nil {
			email.GET("/config", hc.Email.GetEmailConfig)
			email.POST("/validate", hc.Email.ValidateEmail)
			email.POST("/test", hc.Email.TestEmailRouting)
			
			// Admin endpoints for email management
			adminEmail := email.Group("")
			adminEmail.Use(middleware.AdminMiddleware())
			{
				adminEmail.POST("/deploy-worker", hc.Email.DeployEmailWorker)
				adminEmail.POST("/setup-dns", hc.Email.SetupEmailDNS)
				adminEmail.GET("/allow-list", hc.Email.GetEmailAllowList)
				adminEmail.POST("/allow-list/add", hc.Email.AddToEmailAllowList)
				adminEmail.POST("/allow-list/remove", hc.Email.RemoveFromEmailAllowList)
			}
		} else {
			email.GET("/status", func(c *gin.Context) {
				c.JSON(200, gin.H{"status": "Email service not configured"})
			})
		}
	}

	// ==================== Protected Routes ====================
	// Apply authentication middleware to all protected routes
	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware(cfg))
	
	// User routes
	user := protected.Group("/user")
	{
		if hc.User != nil {
			user.GET("/profile", hc.User.GetProfile)
			user.PUT("/profile", hc.User.UpdateProfile)
		}
	}
	
	// Email sending (protected)
	emailProtected := protected.Group("/email")
	{
		if hc.Email != nil {
			emailProtected.POST("/send", hc.Email.SendEmail)
		}
	}
	
	// Service routes
	service := protected.Group("/services")
	{
		if hc.Service != nil {
			service.GET("", hc.Service.GetServices)
			service.POST("", hc.Service.CreateService)
		}
	}
	
	// Category routes
	category := protected.Group("/categories")
	{
		if hc.Category != nil {
			category.GET("", hc.Category.GetCategories)
			category.POST("", hc.Category.CreateCategory)
		}
	}
	
	// Order routes
	order := protected.Group("/orders")
	{
		if hc.Order != nil {
			order.GET("", hc.Order.GetUserOrders)
			order.POST("", hc.Order.CreateOrder)
		}
	}
	
	// Payment routes
	payment := protected.Group("/payments")
	{
		if hc.Payment != nil {
			payment.POST("/intent", hc.Payment.CreatePaymentIntent)
			payment.POST("/:id/confirm", hc.Payment.ConfirmPayment)
		}
	}
	
	// Upload routes
	upload := protected.Group("/upload")
	{
		if hc.Upload != nil {
			upload.POST("", hc.Upload.UploadFile)
		}
	}
	
	// Notification routes
	notification := protected.Group("/notifications")
	{
		if hc.Notification != nil {
			notification.GET("", hc.Notification.GetNotifications)
			notification.PUT("/:id/read", hc.Notification.MarkAsRead)
		}
	}
	
	// Admin routes (admin only)
	admin := protected.Group("/admin")
	admin.Use(middleware.AdminMiddleware())
	{
		if hc.Admin != nil {
			admin.GET("/stats", hc.Admin.GetStatistics)
			admin.GET("/users", hc.Admin.GetAllUsers)
		}
		if hc.Email != nil {
			admin.GET("/email/reports", func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "Email reports endpoint"})
			})
		}
	}
	
	// ==================== Root/Home ====================
	app.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to NawthTech API",
			"version": "1.0.0",
			"status":  "running",
			"services": gin.H{
				"auth":      true,
				"email":     hc.Email != nil,
				"ai":        hc.AI != nil,
				"payments":  hc.Payment != nil,
				"services":  hc.Service != nil,
				"health":    hc.Health != nil,
			},
			"endpoints": gin.H{
				"docs":         "/api/v1/docs",
				"health":       "/health",
				"api_version":  "v1",
			},
		})
	})
	
	// API documentation
	api.GET("/docs", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "API Documentation",
			"version": "v1",
			"endpoints": []gin.H{
				{
					"path": "/api/v1/auth/*",
					"description": "Authentication endpoints",
				},
				{
					"path": "/api/v1/email/*",
					"description": "Email management endpoints",
				},
				{
					"path": "/api/v1/ai/*",
					"description": "AI services endpoints",
				},
				{
					"path": "/api/v1/user/*",
					"description": "User profile endpoints (protected)",
				},
				{
					"path": "/api/v1/services/*",
					"description": "Service management endpoints",
				},
			},
		})
	})
	
	// ==================== 404 Handler ====================
	app.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{
			"error":   "Endpoint not found",
			"message": "The requested resource does not exist",
			"path":    c.Request.URL.Path,
			"suggestions": []string{
				"/api/v1/docs",
				"/health",
				"/api/v1/auth/login",
			},
		})
	})
}