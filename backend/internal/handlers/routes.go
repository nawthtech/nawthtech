package handlers

import (
	"backend-app/internal/middleware"
	"backend-app/internal/services"

	"github.com/go-chi/chi/v5"
)

// Services يحتوي على جميع الخدمات
type Services struct {
	Admin   *services.AdminService
	User    *services.UserService
	Auth    *services.AuthService
	Store   *services.StoreService
	Cart    *services.CartService
	Payment *services.PaymentService
	AI      *services.AIService
	Email   *services.EmailService
	Upload  *services.UploadService
}

// Register تسجيل جميع المسارات
func Register(router chi.Router, services *Services) {
	// إنشاء المعالجات
	adminHandler := NewAdminHandler(services.Admin)
	authHandler := NewAuthHandler(services.Auth)
	userHandler := NewUserHandler(services.User)
	storeHandler := NewStoreHandler(services.Store)
	cartHandler := NewCartHandler(services.Cart)
	paymentHandler := NewPaymentHandler(services.Payment)
	aiHandler := NewAIHandler(services.AI)
	uploadHandler := NewUploadHandler(services.Upload)
	healthHandler := NewHealthHandler()

	// المسارات العامة
	router.Route("/api", func(r chi.Router) {
		// الصحة
		r.Get("/health", healthHandler.HealthCheck)
		
		// المصادقة
		r.Route("/auth", func(auth chi.Router) {
			auth.Post("/login", authHandler.Login)
			auth.Post("/register", authHandler.Register)
			auth.Post("/refresh", authHandler.RefreshToken)
			auth.Post("/forgot-password", authHandler.ForgotPassword)
			auth.Post("/reset-password", authHandler.ResetPassword)
		})

		// المتجر (عام)
		r.Route("/store", func(store chi.Router) {
			store.Get("/services", storeHandler.GetServices)
			store.Get("/services/{id}", storeHandler.GetService)
			store.Get("/categories", storeHandler.GetCategories)
		})

		// المسارات المحمية
		r.Group(func(protected chi.Router) {
			protected.Use(middleware.AuthMiddleware)

			// المستخدم
			protected.Route("/users", func(user chi.Router) {
				user.Get("/profile", userHandler.GetProfile)
				user.Put("/profile", userHandler.UpdateProfile)
				user.Post("/change-password", userHandler.ChangePassword)
			})

			// المتجر المحمي
			protected.Route("/store", func(store chi.Router) {
				store.Post("/orders", storeHandler.CreateOrder)
				store.Get("/orders", storeHandler.GetUserOrders)
				store.Get("/orders/{id}", storeHandler.GetOrder)
			})

			// السلة
			protected.Route("/cart", func(cart chi.Router) {
				cart.Get("/", cartHandler.GetCart)
				cart.Post("/items", cartHandler.AddToCart)
				cart.Put("/items/{id}", cartHandler.UpdateCartItem)
				cart.Delete("/items/{id}", cartHandler.RemoveFromCart)
				cart.Delete("/clear", cartHandler.ClearCart)
			})

			// المدفوعات
			protected.Route("/payments", func(payment chi.Router) {
				payment.Post("/create", paymentHandler.CreatePayment)
				payment.Get("/{id}", paymentHandler.GetPayment)
				payment.Post("/{id}/verify", paymentHandler.VerifyPayment)
			})

			// الذكاء الاصطناعي
			protected.Route("/ai", func(ai chi.Router) {
				ai.Post("/analyze", aiHandler.AnalyzeContent)
				ai.Get("/recommend", aiHandler.GetRecommendations)
				ai.Post("/generate", aiHandler.GenerateContent)
			})

			// الرفع
			protected.Route("/upload", func(upload chi.Router) {
				upload.Post("/image", uploadHandler.UploadImage)
				upload.Post("/file", uploadHandler.UploadFile)
			})
		})

		// مسارات الإدارة (تحتاج صلاحيات إدارة)
		r.Route("/v1/admin", func(admin chi.Router) {
			admin.Use(middleware.AuthMiddleware)
			admin.Use(middleware.AdminAuth)

			// تسجيل مسارات الإدارة
			RegisterAdminRoutes(admin, adminHandler)
		})
	})
}

// RegisterAdminRoutes تسجيل مسارات الإدارة
func RegisterAdminRoutes(router chi.Router, adminHandler *AdminHandler) {
	// مسارات تحديث النظام
	router.With(
		middleware.RateLimiter(5, 30*60*1000),
		middleware.ValidateAdminAction,
		middleware.UpdateMiddleware,
	).Post("/system/update", adminHandler.InitiateSystemUpdate)
	
	// مسارات حالة النظام
	router.With(
		middleware.RateLimiter(30, 2*60*1000),
	).Get("/system/status", adminHandler.GetSystemStatus)
	
	router.With(
		middleware.RateLimiter(10, 5*60*1000),
	).Get("/system/health", adminHandler.GetSystemHealth)
	
	// مسارات تحليلات الذكاء الاصطناعي
	router.With(
		middleware.RateLimiter(15, 5*60*1000),
	).Get("/system/ai-analytics", adminHandler.GetAIAnalytics)
	
	// مسارات تحليلات المستخدمين
	router.With(
		middleware.RateLimiter(20, 10*60*1000),
	).Get("/users/analytics", adminHandler.GetUserAnalytics)
	
	// مسارات الصيانة
	router.With(
		middleware.RateLimiter(10, 10*60*1000),
		middleware.ValidateAdminAction,
	).Post("/system/maintenance", adminHandler.SetMaintenanceMode)
	
	// مسارات السجلات
	router.With(
		middleware.RateLimiter(30, 2*60*1000),
	).Get("/system/logs", adminHandler.GetSystemLogs)
	
	// مسارات النسخ الاحتياطي
	router.With(
		middleware.RateLimiter(5, 60*60*1000),
		middleware.ValidateAdminAction,
	).Post("/system/backup", adminHandler.CreateSystemBackup)
	
	// مسارات التحسين
	router.With(
		middleware.RateLimiter(10, 30*60*1000),
		middleware.ValidateAdminAction,
	).Post("/system/optimize", adminHandler.PerformOptimization)

	// إدارة المستخدمين
	router.Get("/users", adminHandler.GetUsers)
	router.Get("/users/{id}", adminHandler.GetUser)
	router.Put("/users/{id}", adminHandler.UpdateUser)
	router.Delete("/users/{id}", adminHandler.DeleteUser)

	// إدارة الطلبات
	router.Get("/orders", adminHandler.GetOrders)
	router.Get("/orders/{id}", adminHandler.GetOrder)
	router.Put("/orders/{id}/status", adminHandler.UpdateOrderStatus)

	// إدارة الخدمات
	router.Get("/services", adminHandler.GetServices)
	router.Post("/services", adminHandler.CreateService)
	router.Put("/services/{id}", adminHandler.UpdateService)
	router.Delete("/services/{id}", adminHandler.DeleteService)
}
