package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/nawthtech/nawthtech/backend/internal/handlers/sse"
)

// RegisterSSERoutes تسجيل مسارات SSE
func RegisterSSERoutes(router *gin.RouterGroup, authMiddleware gin.HandlerFunc, adminMiddleware gin.HandlerFunc) {
	sseRoutes := router.Group("/sse")
	
	// مسارات عامة (قد تحتاج مصادقة أساسية)
	publicSSE := sseRoutes.Group("")
	publicSSE.Use(authMiddleware)
	{
		publicSSE.GET("/events", sse.Handler)
		publicSSE.GET("/notifications", sse.NotificationHandler)
	}
	
	// مسارات المسؤولين
	adminSSE := sseRoutes.Group("/admin")
	adminSSE.Use(adminMiddleware)
	{
		adminSSE.GET("/events", sse.AdminHandler)
		adminSSE.GET("/monitoring", sse.AdminHandler)
	}
}