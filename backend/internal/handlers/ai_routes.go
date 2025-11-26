package handlers

import (
	"github.com/nawthtech/nawthtech/backend/internal/middleware"
	"github.com/nawthtech/nawthtech/backend/internal/services"

	"github.com/go-chi/chi/v5"
)

// RegisterAIRoutes تسجيل جميع مسارات الذكاء الاصطناعي
func RegisterAIRoutes(router chi.Router, aiService *services.AIService) {
	aiHandler := NewAIHandler(aiService)

	router.Route("/api/v1/ai", func(r chi.Router) {
		// تطبيق المصادقة على جميع المسارات
		r.Use(middleware.AuthMiddleware)

		// ==================== التحليلات الأساسية ====================
		r.Route("/analyze", func(analyze chi.Router) {
			// تحليل احتياجات المستخدم باستخدام الذكاء الاصطناعي
			analyze.With(middleware.RateLimiter(30, 15*60*1000)).Post("/needs", aiHandler.AnalyzeUserNeeds)
			
			// التحقق من صحة الطلب باستخدام الذكاء الاصطناعي
			analyze.With(middleware.RateLimiter(50, 10*60*1000)).Post("/validate-order", aiHandler.ValidateOrder)
			
			// تحليل المحتوى باستخدام الذكاء الاصطناعي
			analyze.With(middleware.RateLimiter(40, 10*60*1000)).Post("/content", aiHandler.AnalyzeContent)
			
			// تحليل شامل متعدد الأبعاد باستخدام الذكاء الاصطناعي
			analyze.With(middleware.RateLimiter(20, 30*60*1000)).Post("/comprehensive", aiHandler.ComprehensiveAnalysis)
			
			// مقارنة نصوص متعددة باستخدام الذكاء الاصطناعي
			analyze.With(middleware.RateLimiter(30, 15*60*1000)).Post("/compare-texts", aiHandler.CompareTexts)
		})

		// ==================== توليد المحتوى ====================
		r.Route("/generate", func(generate chi.Router) {
			// توليد محتوى باستخدام الذكاء الاصطناعي
			generate.With(middleware.RateLimiter(25, 15*60*1000)).Post("/content", aiHandler.GenerateContent)
			
			// توليد صور باستخدام الذكاء الاصطناعي
			generate.With(middleware.RateLimiter(15, 30*60*1000)).Post("/images", aiHandler.GenerateImages)
		})

		// ==================== إدارة النمو والاستراتيجيات ====================
		r.Route("/orders/{orderId}", func(orders chi.Router) {
			// إنشاء تقرير نمو ذكي للطلب
			orders.With(middleware.RateLimiter(10, 60*60*1000)).Post("/growth-report", aiHandler.GenerateGrowthReport)
			
			// بدء استراتيجية نمو ذكية للطلب
			orders.With(middleware.RateLimiter(5, 60*60*1000)).Post("/start-strategy", aiHandler.StartGrowthStrategy)
		})

		// ==================== المساعدة في النماذج ====================
		r.Route("/assist", func(assist chi.Router) {
			// المساعدة في ملء النماذج باستخدام الذكاء الاصطناعي
			assist.With(middleware.RateLimiter(40, 10*60*1000)).Post("/form", aiHandler.AssistForm)
		})

		// ==================== التوصيات الذكية ====================
		r.Route("/recommendations", func(recommendations chi.Router) {
			// توليد توصيات ذكية مخصصة
			recommendations.With(middleware.RateLimiter(30, 15*60*1000)).Post("/", aiHandler.GenerateRecommendations)
		})

		// ==================== إدارة السجلات والملاحظات ====================
		r.Route("/logs", func(logs chi.Router) {
			// الحصول على سجلات الذكاء الاصطناعي للمستخدم
			logs.Get("/", aiHandler.GetAILogs)
			
			// إضافة ملاحظات على نتيجة الذكاء الاصطناعي
			logs.Post("/{logId}/feedback", aiHandler.AddAIFeedback)
		})

		// ==================== حالة النظام والإحصائيات ====================
		r.Route("/status", func(status chi.Router) {
			// الحصول على حالة نظام الذكاء الاصطناعي
			status.Get("/", aiHandler.GetAIStatus)
			
			// الحصول على إحصائيات استخدام الذكاء الاصطناعي
			status.Get("/usage", aiHandler.GetAIUsage)
		})

		// ==================== إدارة الذكاء الاصطناعي (للمسؤولين) ====================
		r.Route("/admin", func(admin chi.Router) {
			admin.Use(middleware.AdminAuth)

			// الحصول على إحصائيات الذكاء الاصطناعي الشاملة
			admin.Get("/stats", aiHandler.GetAIStats)
			
			// الحصول على معلومات نماذج الذكاء الاصطناعي
			admin.Get("/models", aiHandler.GetAIModels)
			
			// تحديث نموذج الذكاء الاصطناعي
			admin.Post("/models/{modelId}/update", aiHandler.UpdateAIModel)
			
			// إعادة تدريب نماذج الذكاء الاصطناعي
			admin.Post("/retrain", aiHandler.RetrainModels)
		})

		// ==================== اختبار وتقييم النماذج ====================
		r.Route("/test", func(test chi.Router) {
			test.Use(middleware.AdminAuth)

			// اختبار نماذج الذكاء الاصطناعي
			test.Post("/", aiHandler.TestModels)
		})
	})
}