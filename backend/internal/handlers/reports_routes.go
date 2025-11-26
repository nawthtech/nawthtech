package handlers

import (
	"github.com/gin-gonic/gin"
)

func RegisterReportsRoutes(router *gin.RouterGroup, reportsHandler *ReportsHandler, authMiddleware gin.HandlerFunc, adminMiddleware gin.HandlerFunc) {
	reportsRoutes := router.Group("/reports")
	reportsRoutes.Use(authMiddleware)
	
	// ==================== التقارير التلقائية ====================
	reportsRoutes.POST("/generate", adminMiddleware, reportsHandler.GenerateReport)
	reportsRoutes.POST("/compare", adminMiddleware, reportsHandler.GenerateComparisonReport)
	
	// ==================== إدارة التقارير ====================
	reportsRoutes.GET("", adminMiddleware, reportsHandler.GetReports)
	reportsRoutes.GET("/:id", adminMiddleware, reportsHandler.GetReportByID)
	reportsRoutes.PUT("/:id", adminMiddleware, reportsHandler.UpdateReport)
	reportsRoutes.DELETE("/:id", adminMiddleware, reportsHandler.DeleteReport)
	
	// ==================== تحليل التقارير ====================
	reportsRoutes.POST("/:id/analyze", adminMiddleware, reportsHandler.AnalyzeReport)
	reportsRoutes.GET("/:id/export", adminMiddleware, reportsHandler.ExportReport)
	
	// ==================== تقارير الأداء ====================
	reportsRoutes.GET("/performance/dashboard", adminMiddleware, reportsHandler.GetDashboardPerformance)
	reportsRoutes.POST("/performance/custom", adminMiddleware, reportsHandler.GenerateCustomPerformanceReport)
}