package handlers

import (
	"github.com/nawthtech/nawthtech/backend/internal/middleware"
	"github.com/nawthtech/nawthtech/backend/internal/services"

	"github.com/go-chi/chi/v5"
)

// RegisterAuthRoutes تسجيل جميع مسارات المصادقة
func RegisterAuthRoutes(router chi.Router, authService *services.AuthService) {
	authHandler := NewAuthHandler(authService)

	router.Route("/api/v1/auth", func(r chi.Router) {
		// ==================== Routes العامة (لا تتطلب مصادقة) ====================
		// تسجيل مستخدم جديد مع تحليل الذكاء الاصطناعي
		r.With(middleware.RateLimiter(5, 15*60*1000)).Post("/register", authHandler.Register)
		
		// تسجيل الدخول مع كشف الاحتيال بالذكاء الاصطناعي
		r.With(middleware.RateLimiter(10, 15*60*1000)).Post("/login", authHandler.Login)
		
		// تجديد token مع تحليل الأمان
		r.With(middleware.RateLimiter(20, 10*60*1000)).Post("/refresh-token", authHandler.RefreshToken)
		
		// طلب إعادة تعيين كلمة المرور مع تحليل الذكاء الاصطناعي
		r.With(middleware.RateLimiter(5, 60*60*1000)).Post("/forgot-password", authHandler.ForgotPassword)
		
		// إعادة تعيين كلمة المرور مع تحليل القوة
		r.With(middleware.RateLimiter(5, 30*60*1000)).Post("/reset-password", authHandler.ResetPassword)
		
		// التحقق من البريد الإلكتروني
		r.With(middleware.RateLimiter(10, 30*60*1000)).Post("/verify-email", authHandler.VerifyEmail)
		
		// إعادة إرسال رابط التحقق
		r.With(middleware.RateLimiter(3, 15*60*1000)).Post("/resend-verification", authHandler.ResendVerificationEmail)

		// ==================== Routes المحمية (تتطلب مصادقة) ====================
		r.Route("/protected", func(protected chi.Router) {
			protected.Use(middleware.AuthMiddleware)

			// الحصول على بيانات المستخدم الحالي مع تحليل السلوك
			protected.With(middleware.RateLimiter(60, 5*60*1000)).Get("/me", authHandler.GetCurrentUser)
			
			// تسجيل الخروج مع تحليل الجلسة
			protected.Post("/logout", authHandler.Logout)
			
			// تحديث الملف الشخصي مع تحليل التغييرات
			protected.With(middleware.RateLimiter(20, 10*60*1000)).Put("/profile", authHandler.UpdateProfile)
			
			// تغيير كلمة المرور مع تحليل متقدم
			protected.With(middleware.RateLimiter(5, 60*60*1000)).Put("/change-password", authHandler.ChangePassword)

			// ==================== Routes جديدة مع الذكاء الاصطناعي ====================
			// الحصول على رؤى أمنية شخصية باستخدام الذكاء الاصطناعي
			protected.With(middleware.RateLimiter(10, 30*60*1000)).Get("/security-insights", authHandler.GetSecurityInsights)
			
			// تحليل سلوك المستخدم للكشف عن الشذوذ
			protected.With(middleware.RateLimiter(5, 60*60*1000)).Post("/behavior-analysis", authHandler.AnalyzeUserBehavior)
			
			// الحصول على تحليلات الجلسات باستخدام الذكاء الاصطناعي
			protected.With(middleware.RateLimiter(15, 10*60*1000)).Get("/session-analytics", authHandler.GetSessionAnalytics)
			
			// تقييم مخاطر حساب المستخدم
			protected.With(middleware.RateLimiter(5, 24*60*60*1000)).Post("/risk-assessment", authHandler.AssessUserRisk)

			// ==================== إدارة الجلسات المتقدمة ====================
			// الحصول على الجلسات النشطة مع تحليل الذكاء الاصطناعي
			protected.With(middleware.RateLimiter(20, 5*60*1000)).Get("/active-sessions", authHandler.GetActiveSessions)
			
			// إنهاء جلسة محددة
			protected.With(middleware.RateLimiter(10, 5*60*1000)).Post("/terminate-session", authHandler.TerminateSession)
		})
	})
}