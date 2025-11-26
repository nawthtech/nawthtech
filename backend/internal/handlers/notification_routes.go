package handlers

import (
	"github.com/gin-gonic/gin"
)

func RegisterNotificationRoutes(router *gin.RouterGroup, notificationHandler *NotificationHandler, authMiddleware gin.HandlerFunc, adminMiddleware gin.HandlerFunc) {
	notificationRoutes := router.Group("/notifications")
	notificationRoutes.Use(authMiddleware)
	
	// ==================== إدارة الإشعارات ====================
	notificationRoutes.GET("", notificationHandler.GetNotifications)
	notificationRoutes.GET("/stats", notificationHandler.GetNotificationStats)
	notificationRoutes.PUT("/:id/read", notificationHandler.MarkAsRead)
	notificationRoutes.PUT("/read-all", notificationHandler.MarkAllAsRead)
	notificationRoutes.DELETE("/:id", notificationHandler.DeleteNotification)
	notificationRoutes.DELETE("", notificationHandler.DeleteReadNotifications)
	
	// ==================== تفضيلات الإشعارات ====================
	notificationRoutes.GET("/preferences", notificationHandler.GetPreferences)
	notificationRoutes.PUT("/preferences", notificationHandler.UpdatePreferences)
	
	// ==================== الإشعارات الذكية ====================
	notificationRoutes.GET("/ai-recommendations", notificationHandler.GetAIRecommendations)
	
	// ==================== إشعارات النظام (للمشرفين فقط) ====================
	adminNotificationRoutes := notificationRoutes.Group("")
	adminNotificationRoutes.Use(adminMiddleware)
	adminNotificationRoutes.POST("/smart", notificationHandler.CreateSmartNotifications)
	adminNotificationRoutes.POST("/system", notificationHandler.CreateSystemNotification)
}