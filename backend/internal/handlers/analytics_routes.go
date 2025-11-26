package handlers

import (
	"github.com/gin-gonic/gin"
)

func RegisterAnalyticsRoutes(router *gin.RouterGroup, analyticsHandler *AnalyticsHandler, authMiddleware gin.HandlerFunc, adminMiddleware gin.HandlerFunc) {
	analyticsRoutes := router.Group("/analytics")
	analyticsRoutes.Use(authMiddleware)
	
	// ==================== التحليلات الأساسية ====================
	analyticsRoutes.GET("/overview", adminMiddleware, analyticsHandler.GetOverview)
	analyticsRoutes.GET("/performance", adminMiddleware, analyticsHandler.GetPerformance)
	
	// ==================== تحليلات الذكاء الاصطناعي ====================
	analyticsRoutes.GET("/ai-insights", adminMiddleware, analyticsHandler.GetAIInsights)
	
	// ==================== تحليلات المحتوى ====================
	analyticsRoutes.GET("/content", adminMiddleware, analyticsHandler.GetContentAnalytics)
	
	// ==================== تحليلات الجمهور ====================
	analyticsRoutes.GET("/audience", adminMiddleware, analyticsHandler.GetAudienceAnalytics)
	
	// ==================== التقارير المخصصة ====================
	analyticsRoutes.POST("/custom-report", adminMiddleware, analyticsHandler.GenerateCustomReport)
	analyticsRoutes.GET("/custom-reports", adminMiddleware, analyticsHandler.GetCustomReports)
	
	// ==================== التوقعات والتنبؤات ====================
	analyticsRoutes.GET("/predictions", adminMiddleware, analyticsHandler.GetPredictions)
}