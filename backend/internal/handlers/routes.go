package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/nawthtech/nawthtech/backend/internal/config"
)

// RegisterAllRoutes
func RegisterAllRoutes(app *gin.Engine, cfg *config.Config, hc *HandlerContainer) {
	// Public
	api := app.Group("/api/v1")
	api.POST("/auth/register", hc.Register)
	api.POST("/auth/login", hc.Login)

	// Health
	app.GET("/health", hc.Health)
	app.GET("/health/live", func(c *gin.Context){ c.JSON(200, gin.H{"status":"live"}) })
	app.GET("/health/ready", func(c *gin.Context){ c.JSON(200, gin.H{"status":"ready"}) })

	// Protected (attach auth middleware)
	protected := api.Group("")
	protected.Use(AuthMiddleware(cfg))
	protected.GET("/user/profile", hc.GetProfile)
	protected.GET("/services", hc.GetServices)
}