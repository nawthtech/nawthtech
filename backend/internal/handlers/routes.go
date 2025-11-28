package handlers

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nawthtech/nawthtech/backend/internal/config"
	"github.com/nawthtech/nawthtech/backend/internal/handlers/health"
	"github.com/nawthtech/nawthtech/backend/internal/handlers/sse"
	"github.com/nawthtech/nawthtech/backend/internal/middleware"
	"github.com/nawthtech/nawthtech/backend/internal/services"
	"gorm.io/gorm"
)

// RegisterAllRoutes تسجيل جميع المسارات
func RegisterAllRoutes(router *gin.Engine, serviceContainer *services.ServiceContainer, config *config.Config, db *gorm.DB) {
	// تطبيق middleware العام على مستوى التطبيق
	applyGlobalMiddleware(router, config)
	
	// إنشاء حاوية الوسائط
	middlewares := initializeMiddlewares(serviceContainer, config)
	
	// مجموعة API الرئيسية
	api := router.Group("/api/v1")
	
	// ========== مسارات الصحة ==========
	registerHealthRoutes(router, serviceContainer, config, db)
	
	// ========== المسارات العامة (لا تتطلب مصادقة) ==========
	registerPublicRoutes(api, serviceContainer, middlewares)
	
	// ========== المسارات المحمية (تتطلب مصادقة) ==========
	registerProtectedRoutes(api, serviceContainer, middlewares)
	
	// ========== مسارات المسؤولين ==========
	registerAdminRoutes(api, serviceContainer, middlewares)
	
	// ========== مسارات البائعين ==========
	registerSellerRoutes(api, serviceContainer, middlewares)
	
	// ========== مسارات الوقت الحقيقي (SSE) ==========
	registerSSERoutes(api, serviceContainer, middlewares)
	
	// ========== مسارات الويب هووك ==========
	registerWebhookRoutes(api, serviceContainer, middlewares)
}

// applyGlobalMiddleware تطبيق الوسائط العامة على مستوى التطبيق
func applyGlobalMiddleware(router *gin.Engine, config *config.Config) {
	// CORS middleware - يتم تطبيقه على مستوى التطبيق بالكامل
	router.Use(middleware.CORS())
	
	// Security headers middleware
	router.Use(middleware.SecurityHeaders())
	
	// Rate limiting middleware
	router.Use(middleware.RateLimit())
}

// initializeMiddlewares تهيئة جميع الوسائط
func initializeMiddlewares(services *services.ServiceContainer, config *config.Config) *middleware.MiddlewareContainer {
	return &middleware.MiddlewareContainer{
		AuthMiddleware:      middleware.AuthMiddleware(services.Auth),
		AdminMiddleware:     middleware.AdminMiddleware(),
		CORSMiddleware:      middleware.CORS(),
		SecurityMiddleware:  middleware.SecurityHeaders(),
		RateLimitMiddleware: middleware.RateLimit(),
	}
}

// registerHealthRoutes تسجيل مسارات الصحة
func registerHealthRoutes(router *gin.Engine, services *services.ServiceContainer, config *config.Config, db *gorm.DB) {
	// إنشاء معالج الصحة
	healthHandler := health.NewHealthHandler(db, services.Cache, config)
	
	// مسارات الصحة العامة (بدون بادئة api/v1)
	router.GET("/health", healthHandler.Check)
	router.GET("/health/live", healthHandler.Live)
	router.GET("/health/ready", healthHandler.Ready)
	router.GET("/health/info", healthHandler.Info)
	router.GET("/health/detailed", healthHandler.Detailed)
	router.GET("/health/metrics", healthHandler.Metrics)
	
	// مسارات الصحة للمسؤولين (مع المصادقة)
	adminHealth := router.Group("/health")
	adminHealth.Use(middleware.AuthMiddleware(services.Auth), middleware.AdminMiddleware())
	adminHealth.GET("/admin", healthHandler.AdminHealth)
}

// registerSSERoutes تسجيل مسارات SSE
func registerSSERoutes(api *gin.RouterGroup, services *services.ServiceContainer, middlewares *middleware.MiddlewareContainer) {
	// مسارات SSE (تتطلب مصادقة)
	sseGroup := api.Group("/sse")
	sseGroup.Use(middlewares.AuthMiddleware)
	
	// مسار SSE الرئيسي
	sseGroup.GET("/events", sse.Handler)
	
	// مسار إشعارات SSE
	sseGroup.GET("/notifications", sse.NotificationHandler)
	
	// مسار SSE للمسؤولين
	adminSSE := sseGroup.Group("/admin")
	adminSSE.Use(middlewares.AdminMiddleware)
	adminSSE.GET("/events", sse.AdminHandler)
	
	// مسار إحصائيات اتصالات SSE (للمسؤولين فقط)
	adminSSE.GET("/connections", func(c *gin.Context) {
		stats := sse.GetConnectionStats()
		c.JSON(200, gin.H{
			"success": true,
			"message": "إحصائيات اتصالات SSE",
			"data":    stats,
		})
	})
}

// registerPublicRoutes تسجيل المسارات العامة
func registerPublicRoutes(api *gin.RouterGroup, services *services.ServiceContainer, middlewares *middleware.MiddlewareContainer) {
	// معالج المصادقة
	authHandler := NewAuthHandler(services.Auth)
	api.POST("/auth/register", authHandler.Register)
	api.POST("/auth/login", authHandler.Login)
	api.POST("/auth/refresh", authHandler.RefreshToken)
	api.POST("/auth/forgot-password", authHandler.ForgotPassword) // تم التصحيح هنا
	api.POST("/auth/reset-password", authHandler.ResetPassword)
	
	// معالج الخدمات (العامة)
	serviceHandler := NewServiceHandler(services.Service)
	api.GET("/services", serviceHandler.GetServices)
	api.GET("/services/search", serviceHandler.SearchServices)
	api.GET("/services/featured", serviceHandler.GetFeaturedServices)
	api.GET("/services/categories", serviceHandler.GetCategories)
	api.GET("/services/:id", serviceHandler.GetServiceByID)
	
	// معالج الفئات
	categoryHandler := NewCategoryHandler(services.Category)
	api.GET("/categories", categoryHandler.GetCategories)
	api.GET("/categories/tree", categoryHandler.GetCategoryTree)
	api.GET("/categories/:id", categoryHandler.GetCategoryByID)
	
	// معالج المتاجر
	storeHandler := NewStoreHandler(services.Store)
	api.GET("/stores", storeHandler.GetStores)
	api.GET("/stores/featured", storeHandler.GetFeaturedStores)
	api.GET("/stores/:id", storeHandler.GetStoreByID)
	api.GET("/stores/slug/:slug", storeHandler.GetStoreBySlug)
	
	// معالج المحتوى
	contentHandler := NewContentHandler(services.Content)
	api.GET("/content", contentHandler.GetContentList)
	api.GET("/content/:id", contentHandler.GetContentByID)
	api.GET("/content/slug/:slug", contentHandler.GetContentBySlug)
}

// registerProtectedRoutes تسجيل المسارات المحمية
func registerProtectedRoutes(api *gin.RouterGroup, services *services.ServiceContainer, middlewares *middleware.MiddlewareContainer) {
	protected := api.Group("")
	protected.Use(middlewares.AuthMiddleware)
	
	// معالج المستخدم
	userHandler := NewUserHandler(services.User)
	protected.GET("/user/profile", userHandler.GetProfile)
	protected.PUT("/user/profile", userHandler.UpdateProfile)
	protected.PUT("/user/password", userHandler.ChangePassword)
	protected.GET("/user/stats", userHandler.GetUserStats)
	
	// معالج السلة
	cartHandler := NewCartHandler(services.Cart)
	protected.GET("/cart", cartHandler.GetCart)
	protected.POST("/cart/items", cartHandler.AddToCart)
	protected.PUT("/cart/items/:id", cartHandler.UpdateCartItem)
	protected.DELETE("/cart/items/:id", cartHandler.RemoveFromCart)
	protected.DELETE("/cart", cartHandler.ClearCart)
	protected.GET("/cart/summary", cartHandler.GetCartSummary)
	protected.POST("/cart/apply-coupon", cartHandler.ApplyCoupon)
	protected.DELETE("/cart/coupon", cartHandler.RemoveCoupon)
	
	// معالج الطلبات
	orderHandler := NewOrderHandler(services.Order)
	protected.GET("/orders", orderHandler.GetUserOrders)
	protected.GET("/orders/:id", orderHandler.GetOrderByID)
	protected.POST("/orders", orderHandler.CreateOrder)
	protected.PUT("/orders/:id/cancel", orderHandler.CancelOrder)
	protected.GET("/orders/:id/track", orderHandler.TrackOrder)
	
	// معالج الدفع
	paymentHandler := NewPaymentHandler(services.Payment)
	protected.GET("/payment/methods", paymentHandler.GetPaymentMethods)
	protected.POST("/payment/methods", paymentHandler.AddPaymentMethod)
	protected.DELETE("/payment/methods/:id", paymentHandler.RemovePaymentMethod)
	protected.GET("/payment/history", paymentHandler.GetPaymentHistory)
	protected.POST("/payment/intent", paymentHandler.CreatePaymentIntent)
	protected.POST("/payment/confirm", paymentHandler.ConfirmPayment)
	
	// معالج الرفع
	uploadHandler := NewUploadHandler(services.Upload)
	protected.POST("/upload", uploadHandler.UploadFile)
	protected.GET("/upload/files", uploadHandler.GetUserFiles)
	protected.GET("/upload/files/:id", uploadHandler.GetFile)
	protected.DELETE("/upload/files/:id", uploadHandler.DeleteFile)
	protected.POST("/upload/presigned-url", uploadHandler.GeneratePresignedURL)
	protected.GET("/upload/quota", uploadHandler.GetUploadQuota)
	
	// معالج الإشعارات
	notificationHandler := NewNotificationHandler(services.Notification)
	protected.GET("/notifications", notificationHandler.GetUserNotifications)
	protected.PUT("/notifications/:id/read", notificationHandler.MarkAsRead)
	protected.PUT("/notifications/read-all", notificationHandler.MarkAllAsRead)
	protected.DELETE("/notifications/:id", notificationHandler.DeleteNotification)
	protected.GET("/notifications/unread-count", notificationHandler.GetUnreadCount)
	
	// معالج قائمة الرغبات
	wishlistHandler := NewWishlistHandler(services.Wishlist)
	protected.GET("/wishlist", wishlistHandler.GetUserWishlist)
	protected.POST("/wishlist/:serviceId", wishlistHandler.AddToWishlist)
	protected.DELETE("/wishlist/:serviceId", wishlistHandler.RemoveFromWishlist)
	protected.GET("/wishlist/check/:serviceId", wishlistHandler.IsInWishlist)
	protected.GET("/wishlist/count", wishlistHandler.GetWishlistCount)
	
	// معالج الاشتراكات
	subscriptionHandler := NewSubscriptionHandler(services.Subscription)
	protected.GET("/subscription", subscriptionHandler.GetUserSubscription)
	protected.POST("/subscription", subscriptionHandler.CreateSubscription)
	protected.PUT("/subscription/cancel", subscriptionHandler.CancelSubscription)
	protected.PUT("/subscription/renew", subscriptionHandler.RenewSubscription)
	protected.GET("/subscription/plans", subscriptionHandler.GetSubscriptionPlans)
	
	// معالج الذكاء الاصطناعي
	aiHandler := NewAIHandler(services.AI)
	protected.POST("/ai/generate-text", aiHandler.GenerateText)
	protected.POST("/ai/analyze-sentiment", aiHandler.AnalyzeSentiment)
	protected.POST("/ai/classify-content", aiHandler.ClassifyContent)
	protected.POST("/ai/extract-keywords", aiHandler.ExtractKeywords)
	protected.POST("/ai/summarize-text", aiHandler.SummarizeText)
	protected.POST("/ai/translate", aiHandler.TranslateText)
	protected.POST("/ai/generate-image", aiHandler.GenerateImage)
	protected.POST("/ai/chat", aiHandler.ChatCompletion)
}

// registerAdminRoutes تسجيل مسارات المسؤولين
func registerAdminRoutes(api *gin.RouterGroup, services *services.ServiceContainer, middlewares *middleware.MiddlewareContainer) {
	admin := api.Group("/admin")
	admin.Use(middlewares.AuthMiddleware, middlewares.AdminMiddleware)
	
	// معالج الإدارة
	adminHandler := NewAdminHandler(services.Admin)
	admin.GET("/dashboard", adminHandler.GetDashboard)
	admin.GET("/dashboard/stats", adminHandler.GetDashboardStats)
	admin.GET("/users", adminHandler.GetUsers)
	admin.PUT("/users/:id/status", adminHandler.UpdateUserStatus)
	admin.PUT("/users/:id/role", adminHandler.UpdateUserRole)
	admin.GET("/system/logs", adminHandler.GetSystemLogs)
	admin.PUT("/system/settings", adminHandler.UpdateSystemSettings)
	
	// معالج التقارير
	reportHandler := NewReportHandler(services.Report)
	admin.GET("/reports/sales", reportHandler.GenerateSalesReport)
	admin.GET("/reports/users", reportHandler.GenerateUserReport)
	admin.GET("/reports/services", reportHandler.GenerateServiceReport)
	admin.GET("/reports/financial", reportHandler.GenerateFinancialReport)
	admin.GET("/reports/system", reportHandler.GenerateSystemReport)
	admin.GET("/reports/templates", reportHandler.GetReportTemplates)
	admin.POST("/reports/schedule", reportHandler.ScheduleReport)
	admin.GET("/reports/scheduled", reportHandler.GetScheduledReports)
	
	// معالج التحليلات
	analyticsHandler := NewAnalyticsHandler(services.Analytics)
	admin.GET("/analytics/user", analyticsHandler.GetUserAnalytics)
	admin.GET("/analytics/service", analyticsHandler.GetServiceAnalytics)
	admin.GET("/analytics/platform", analyticsHandler.GetPlatformAnalytics)
	
	// معالج الفئات (الإدارة)
	categoryHandler := NewCategoryHandler(services.Category)
	admin.POST("/categories", categoryHandler.CreateCategory)
	admin.PUT("/categories/:id", categoryHandler.UpdateCategory)
	admin.DELETE("/categories/:id", categoryHandler.DeleteCategory)
	admin.GET("/categories/stats", categoryHandler.GetCategoryStats)
	
	// معالج المتاجر (الإدارة)
	storeHandler := NewStoreHandler(services.Store)
	admin.POST("/stores", storeHandler.CreateStore)
	admin.PUT("/stores/:id", storeHandler.UpdateStore)
	admin.DELETE("/stores/:id", storeHandler.DeleteStore)
	admin.POST("/stores/:id/verify", storeHandler.VerifyStore)
	admin.GET("/stores/:id/stats", storeHandler.GetStoreStats)
	
	// معالج الطلبات (الإدارة)
	orderHandler := NewOrderHandler(services.Order)
	admin.GET("/orders", orderHandler.GetAllOrders)
	admin.PUT("/orders/:id/status", orderHandler.UpdateOrderStatus)
	admin.GET("/orders/stats", orderHandler.GetOrderStats)
	
	// معالج الكوبونات
	couponHandler := NewCouponHandler(services.Coupon)
	admin.POST("/coupons", couponHandler.CreateCoupon)
	admin.GET("/coupons", couponHandler.GetCoupons)
	admin.GET("/coupons/:id", couponHandler.GetCouponByID)
	admin.PUT("/coupons/:id", couponHandler.UpdateCoupon)
	admin.DELETE("/coupons/:id", couponHandler.DeleteCoupon)
	admin.POST("/coupons/validate", couponHandler.ValidateCoupon)
}

// registerSellerRoutes تسجيل مسارات البائعين
func registerSellerRoutes(api *gin.RouterGroup, services *services.ServiceContainer, middlewares *middleware.MiddlewareContainer) {
	seller := api.Group("/seller")
	seller.Use(middlewares.AuthMiddleware, middleware.SellerMiddleware())
	
	// معالج الخدمات (البائعين)
	serviceHandler := NewServiceHandler(services.Service)
	seller.POST("/services", serviceHandler.CreateService)
	seller.PUT("/services/:id", serviceHandler.UpdateService)
	seller.DELETE("/services/:id", serviceHandler.DeleteService)
	seller.GET("/services/my", serviceHandler.GetMyServices)
	seller.GET("/services/stats", serviceHandler.GetServiceStats)
	
	// معالج الطلبات (البائعين)
	orderHandler := NewOrderHandler(services.Order)
	seller.GET("/orders", orderHandler.GetSellerOrders)
	seller.PUT("/orders/:id/status", orderHandler.UpdateOrderStatus)
	seller.GET("/orders/stats", orderHandler.GetSellerOrderStats)
	
	// معالج المتاجر (البائعين)
	storeHandler := NewStoreHandler(services.Store)
	seller.POST("/stores", storeHandler.CreateStore)
	seller.PUT("/stores/my", storeHandler.UpdateStore)
	seller.GET("/stores/my", storeHandler.GetMyStore)
	seller.GET("/stores/my/stats", storeHandler.GetMyStoreStats)
	seller.GET("/stores/my/reviews", storeHandler.GetStoreReviews)
}

// registerWebhookRoutes تسجيل مسارات الويب هووك
func registerWebhookRoutes(api *gin.RouterGroup, services *services.ServiceContainer, middlewares *middleware.MiddlewareContainer) {
	webhook := api.Group("/webhook")
	{
		// ويب هووك الدفع
		paymentHandler := NewPaymentHandler(services.Payment)
		webhook.POST("/payment/stripe", paymentHandler.HandleStripeWebhook)
		webhook.POST("/payment/paypal", paymentHandler.HandlePayPalWebhook)
		
		// ويب هووك الرفع
		uploadHandler := NewUploadHandler(services.Upload)
		webhook.POST("/upload/cloudinary", uploadHandler.HandleCloudinaryWebhook)
		
		// ويب هووك التحليلات
		analyticsHandler := NewAnalyticsHandler(services.Analytics)
		webhook.POST("/analytics/plausible", analyticsHandler.HandlePlausibleWebhook)
	}
}

// HealthHandler معالج الصحة المبسط
type HealthHandler struct {
	config *config.Config
}

func NewHealthHandler(config *config.Config) *HealthHandler {
	return &HealthHandler{
		config: config,
	}
}

func (h *HealthHandler) Check(c *gin.Context) {
	response := gin.H{
		"status":    "healthy",
		"service":   "nawthtech-backend",
		"timestamp": time.Now().Format(time.RFC3339),
		"version":   h.config.Version,
	}
	c.JSON(200, response)
}