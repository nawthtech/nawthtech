// في ملف routes.go أضف هذه الوظيفة
import (
	// ... الاستيرادات الحالية
	"nawthtech/backend/internal/middleware"
)

func RegisterAdminRoutes(router *mux.Router, adminHandler *AdminHandler, authMiddleware mux.MiddlewareFunc) {
	adminRouter := router.PathPrefix("/api/v1/admin").Subrouter()
	
	// تطبيق وسيط المصادقة على جميع مسارات الإدارة
	adminRouter.Use(authMiddleware)
	adminRouter.Use(middleware.AdminAuth)
	
	// مسارات تحديث النظام
	adminRouter.HandleFunc("/system/update", 
		middleware.RateLimiter(5, 30*60*1000)(
			middleware.ValidateAdminAction(
				middleware.UpdateMiddleware(
					adminHandler.InitiateSystemUpdate,
				),
			),
		),
	).Methods("POST")
	
	// مسارات حالة النظام
	adminRouter.HandleFunc("/system/status", 
		middleware.RateLimiter(30, 2*60*1000)(
			adminHandler.GetSystemStatus,
		),
	).Methods("GET")
	
	adminRouter.HandleFunc("/system/health", 
		middleware.RateLimiter(10, 5*60*1000)(
			adminHandler.GetSystemHealth,
		),
	).Methods("GET")
	
	// مسارات تحليلات الذكاء الاصطناعي
	adminRouter.HandleFunc("/system/ai-analytics", 
		middleware.RateLimiter(15, 5*60*1000)(
			adminHandler.GetAIAnalytics,
		),
	).Methods("GET")
	
	// مسارات تحليلات المستخدمين
	adminRouter.HandleFunc("/users/analytics", 
		middleware.RateLimiter(20, 10*60*1000)(
			adminHandler.GetUserAnalytics,
		),
	).Methods("GET")
	
	// مسارات الصيانة
	adminRouter.HandleFunc("/system/maintenance", 
		middleware.RateLimiter(10, 10*60*1000)(
			middleware.ValidateAdminAction(
				adminHandler.SetMaintenanceMode,
			),
		),
	).Methods("POST")
	
	// مسارات السجلات
	adminRouter.HandleFunc("/system/logs", 
		middleware.RateLimiter(30, 2*60*1000)(
			adminHandler.GetSystemLogs,
		),
	).Methods("GET")
	
	// مسارات النسخ الاحتياطي
	adminRouter.HandleFunc("/system/backup", 
		middleware.RateLimiter(5, 60*60*1000)(
			middleware.ValidateAdminAction(
				adminHandler.CreateSystemBackup,
			),
		),
	).Methods("POST")
	
	// مسارات التحسين
	adminRouter.HandleFunc("/system/optimize", 
		middleware.RateLimiter(10, 30*60*1000)(
			middleware.ValidateAdminAction(
				adminHandler.PerformOptimization,
			),
		),
	).Methods("POST")
}
