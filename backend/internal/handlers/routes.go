package handlers

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nawthtech/nawthtech/backend/internal/config"
	"github.com/nawthtech/nawthtech/backend/internal/middleware"
	"github.com/nawthtech/nawthtech/backend/internal/services"
	"github.com/nawthtech/nawthtech/backend/internal/ai"
	"github.com/nawthtech/nawthtech/backend/internal/ai/video"
	"go.mongodb.org/mongo-driver/mongo"
)

// RegisterAllRoutes تسجيل جميع المسارات
func RegisterAllRoutes(router *gin.Engine, serviceContainer *services.ServiceContainer, config *config.Config, 
	mongoClient *mongo.Client, aiClient *ai.Client, videoService *video.VideoService) {
	
	// تطبيق middleware العام على مستوى التطبيق
	applyGlobalMiddleware(router, config)

	// إنشاء حاوية الوسائط
	middlewares := initializeMiddlewares(serviceContainer, config)

	// مجموعة API الرئيسية
	api := router.Group("/api/v1")

	// ========== مسارات الصحة ==========
	registerHealthRoutes(router, config, mongoClient)

	// ========== المسارات العامة (لا تتطلب مصادقة) ==========
	registerPublicRoutes(api, serviceContainer, middlewares, aiClient, videoService)

	// ========== المسارات المحمية (تتطلب مصادقة) ==========
	registerProtectedRoutes(api, serviceContainer, middlewares, aiClient, videoService)

	// ========== مسارات المسؤولين ==========
	registerAdminRoutes(api, serviceContainer, middlewares)

	// ========== مسارات البائعين ==========
	registerSellerRoutes(api, serviceContainer, middlewares)

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

	// Logger middleware
	router.Use(middleware.Logger())

	// Recovery middleware
	router.Use(middleware.Recovery())
}

// initializeMiddlewares تهيئة جميع الوسائط
func initializeMiddlewares(services *services.ServiceContainer, config *config.Config) *middleware.MiddlewareContainer {
	return &middleware.MiddlewareContainer{
		AuthMiddleware:      middleware.AuthMiddleware(services.Auth),
		AdminMiddleware:     middleware.AdminMiddleware(),
		CORSMiddleware:      middleware.CORS(),
		SecurityMiddleware:  middleware.SecurityHeaders(),
		RateLimitMiddleware: middleware.RateLimit(),
		LoggerMiddleware:    middleware.Logger(),
		RecoveryMiddleware:  middleware.Recovery(),
	}
}

// registerHealthRoutes تسجيل مسارات الصحة
func registerHealthRoutes(router *gin.Engine, config *config.Config, mongoClient *mongo.Client) {
	// إنشاء معالج الصحة
	healthHandler := NewHealthHandler(config)

	// مسارات الصحة العامة (بدون بادئة api/v1)
	router.GET("/health", healthHandler.Check)
	router.GET("/health/live", healthHandler.Live)
	router.GET("/health/ready", healthHandler.Ready)
	router.GET("/health/info", healthHandler.Info)

	// مسارات الصحة للمسؤولين (مع المصادقة)
	adminHealth := router.Group("/health")
	// ملاحظة: سيتم تفعيل المصادقة لاحقاً
	// adminHealth.Use(middleware.AuthMiddleware(services.Auth), middleware.AdminMiddleware())
	adminHealth.GET("/admin", healthHandler.AdminHealth)
}

// registerPublicRoutes تسجيل المسارات العامة
func registerPublicRoutes(api *gin.RouterGroup, services *services.ServiceContainer, 
	middlewares *middleware.MiddlewareContainer, aiClient *ai.Client, videoService *video.VideoService) {
	
	// معالج المصادقة
	authHandler := NewAuthHandler(services.Auth)
	api.POST("/auth/register", authHandler.Register)
	api.POST("/auth/login", authHandler.Login)
	api.POST("/auth/refresh", authHandler.RefreshToken)
	api.POST("/auth/forgot-password", authHandler.ForgotPassword)
	api.POST("/auth/reset-password", authHandler.ResetPassword)
	api.POST("/auth/verify-token", authHandler.VerifyToken)

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
	api.GET("/categories/:id", categoryHandler.GetCategoryByID)

	// معالج الذكاء الاصطناعي (الميزات العامة)
	aiHandler := NewAIHandler(aiClient)
	api.GET("/ai/capabilities", aiHandler.GetAICapabilitiesHandler)

	// معالج الفيديو (الميزات العامة)
	videoHandler := NewVideoHandler(videoService)
	api.GET("/video/capabilities", videoHandler.GetVideoCapabilitiesHandler)
	api.GET("/video/types", videoHandler.ListVideoJobsHandler) // إعادة تسمية للتوافق
}

// registerProtectedRoutes تسجيل المسارات المحمية
func registerProtectedRoutes(api *gin.RouterGroup, services *services.ServiceContainer, 
	middlewares *middleware.MiddlewareContainer, aiClient *ai.Client, videoService *video.VideoService) {
	
	protected := api.Group("")
	protected.Use(middlewares.AuthMiddleware)

	// معالج المستخدم
	userHandler := NewUserHandler(services.User)
	protected.GET("/user/profile", userHandler.GetProfile)
	protected.PUT("/user/profile", userHandler.UpdateProfile)
	protected.PUT("/user/password", userHandler.ChangePassword)
	protected.GET("/user/stats", userHandler.GetUserStats)

	// معالج الطلبات
	orderHandler := NewOrderHandler(services.Order)
	protected.GET("/orders", orderHandler.GetUserOrders)
	protected.GET("/orders/:id", orderHandler.GetOrderByID)
	protected.POST("/orders", orderHandler.CreateOrder)
	protected.PUT("/orders/:id/cancel", orderHandler.CancelOrder)
	protected.PUT("/orders/:id/status", orderHandler.UpdateOrderStatus)

	// معالج الدفع
	paymentHandler := NewPaymentHandler(services.Payment)
	protected.GET("/payment/history", paymentHandler.GetPaymentHistory)
	protected.POST("/payment/intent", paymentHandler.CreatePaymentIntent)
	protected.POST("/payment/confirm", paymentHandler.ConfirmPayment)

	// معالج الرفع
	uploadHandler := NewUploadHandler()
	protected.POST("/upload/image", uploadHandler.UploadImage)
	protected.POST("/upload/images", uploadHandler.UploadMultipleImages)
	protected.GET("/upload/images", uploadHandler.GetUserImages)
	protected.GET("/upload/images/:public_id", uploadHandler.GetImageInfo)
	protected.DELETE("/upload/images/:public_id", uploadHandler.DeleteImage)

	// معالج الإشعارات
	notificationHandler := NewNotificationHandler(services.Notification)
	protected.GET("/notifications", notificationHandler.GetUserNotifications)
	protected.PUT("/notifications/:id/read", notificationHandler.MarkAsRead)
	protected.PUT("/notifications/read-all", notificationHandler.MarkAllAsRead)
	protected.GET("/notifications/unread-count", notificationHandler.GetUnreadCount)

	// معالج الذكاء الاصطناعي (المحمي)
	aiHandler := NewAIHandler(aiClient)
	protected.POST("/ai/generate", aiHandler.GenerateContentHandler)
	protected.POST("/ai/analyze-image", aiHandler.AnalyzeImageHandler)
	protected.POST("/ai/translate", aiHandler.TranslateTextHandler)
	protected.POST("/ai/summarize", aiHandler.SummarizeTextHandler)

	// معالج الفيديو (المحمي)
	videoHandler := NewVideoHandler(videoService)
	protected.POST("/video/generate", videoHandler.GenerateVideoHandler)
	protected.GET("/video/jobs", videoHandler.ListVideoJobsHandler)
	protected.GET("/video/jobs/:jobId", videoHandler.GetVideoStatusHandler)
	protected.DELETE("/video/jobs/:jobId", videoHandler.CancelVideoJobHandler)
	protected.GET("/video/jobs/:jobId/download", videoHandler.DownloadVideoHandler)
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
	admin.GET("/system/logs", adminHandler.GetSystemLogs)

	// معالج الفئات (الإدارة)
	categoryHandler := NewCategoryHandler(services.Category)
	admin.POST("/categories", categoryHandler.CreateCategory)
	admin.PUT("/categories/:id", categoryHandler.UpdateCategory)
	admin.DELETE("/categories/:id", categoryHandler.DeleteCategory)

	// معالج الطلبات (الإدارة)
	orderHandler := NewOrderHandler(services.Order)
	admin.GET("/orders", adminHandler.GetAllOrders)
	admin.PUT("/orders/:id/status", orderHandler.UpdateOrderStatus)
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

	// معالج الطلبات (البائعين)
	orderHandler := NewOrderHandler(services.Order)
	seller.GET("/orders", orderHandler.GetSellerOrders)
	seller.PUT("/orders/:id/status", orderHandler.UpdateOrderStatus)
}

// registerWebhookRoutes تسجيل مسارات الويب هووك
func registerWebhookRoutes(api *gin.RouterGroup, services *services.ServiceContainer, middlewares *middleware.MiddlewareContainer) {
	webhook := api.Group("/webhook")
	{
		// ويب هووك الرفع (Cloudinary)
		uploadHandler := NewUploadHandler()
		webhook.POST("/upload/cloudinary", uploadHandler.HandleCloudinaryWebhook)
		
		// ويب هووك الدفع (Stripe, PayPal, etc.)
		paymentHandler := NewPaymentHandler(services.Payment)
		webhook.POST("/payment/stripe", paymentHandler.HandleStripeWebhook)
		webhook.POST("/payment/paypal", paymentHandler.HandlePayPalWebhook)
	}
}

// HealthHandler معالج الصحة
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
		"database":  "MongoDB",
	}
	c.JSON(200, response)
}

func (h *HealthHandler) Live(c *gin.Context) {
	response := gin.H{
		"status":    "alive",
		"timestamp": time.Now().Format(time.RFC3339),
		"service":   "nawthtech-backend",
	}
	c.JSON(200, response)
}

func (h *HealthHandler) Ready(c *gin.Context) {
	response := gin.H{
		"status":    "ready",
		"timestamp": time.Now().Format(time.RFC3339),
		"service":   "nawthtech-backend",
		"database":  "MongoDB",
	}
	c.JSON(200, response)
}

func (h *HealthHandler) Info(c *gin.Context) {
	response := gin.H{
		"name":        "NawthTech Backend",
		"version":     h.config.Version,
		"environment": h.config.Environment,
		"timestamp":   time.Now().Format(time.RFC3339),
		"database":    "MongoDB",
		"features": []string{
			"Authentication",
			"User Management",
			"Services",
			"Categories",
			"Orders",
			"Payments",
			"File Upload",
			"Notifications",
			"AI Services",
			"Video Generation",
		},
		"endpoints": []string{
			"/api/v1/auth/*",
			"/api/v1/services/*",
			"/api/v1/categories/*",
			"/api/v1/ai/*",
			"/api/v1/video/*",
			"/api/v1/user/*",
			"/api/v1/orders/*",
			"/api/v1/payment/*",
			"/api/v1/upload/*",
			"/api/v1/notifications/*",
			"/api/v1/admin/*",
			"/api/v1/seller/*",
		},
	}
	c.JSON(200, response)
}

func (h *HealthHandler) AdminHealth(c *gin.Context) {
	response := gin.H{
		"status":    "healthy",
		"service":   "nawthtech-backend-admin",
		"timestamp": time.Now().Format(time.RFC3339),
		"version":   h.config.Version,
		"database":  "MongoDB",
		"admin":     true,
	}
	c.JSON(200, response)
}

// Placeholder handlers for compilation (يجب تنفيذها في ملفات منفصلة)

// NewAuthHandler placeholder
func NewAuthHandler(authService interface{}) *AuthHandler {
	return &AuthHandler{}
}

type AuthHandler struct{}

func (h *AuthHandler) Register(c *gin.Context)           {}
func (h *AuthHandler) Login(c *gin.Context)              {}
func (h *AuthHandler) RefreshToken(c *gin.Context)       {}
func (h *AuthHandler) ForgotPassword(c *gin.Context)     {}
func (h *AuthHandler) ResetPassword(c *gin.Context)      {}
func (h *AuthHandler) VerifyToken(c *gin.Context)        {}

// NewServiceHandler placeholder
func NewServiceHandler(serviceService interface{}) *ServiceHandler {
	return &ServiceHandler{}
}

type ServiceHandler struct{}

func (h *ServiceHandler) GetServices(c *gin.Context)         {}
func (h *ServiceHandler) SearchServices(c *gin.Context)      {}
func (h *ServiceHandler) GetFeaturedServices(c *gin.Context) {}
func (h *ServiceHandler) GetCategories(c *gin.Context)       {}
func (h *ServiceHandler) GetServiceByID(c *gin.Context)      {}
func (h *ServiceHandler) CreateService(c *gin.Context)       {}
func (h *ServiceHandler) UpdateService(c *gin.Context)       {}
func (h *ServiceHandler) DeleteService(c *gin.Context)       {}
func (h *ServiceHandler) GetMyServices(c *gin.Context)       {}

// NewCategoryHandler placeholder
func NewCategoryHandler(categoryService interface{}) *CategoryHandler {
	return &CategoryHandler{}
}

type CategoryHandler struct{}

func (h *CategoryHandler) GetCategories(c *gin.Context)      {}
func (h *CategoryHandler) GetCategoryByID(c *gin.Context)    {}
func (h *CategoryHandler) CreateCategory(c *gin.Context)     {}
func (h *CategoryHandler) UpdateCategory(c *gin.Context)     {}
func (h *CategoryHandler) DeleteCategory(c *gin.Context)     {}

// NewUserHandler placeholder
func NewUserHandler(userService interface{}) *UserHandler {
	return &UserHandler{}
}

type UserHandler struct{}

func (h *UserHandler) GetProfile(c *gin.Context)             {}
func (h *UserHandler) UpdateProfile(c *gin.Context)          {}
func (h *UserHandler) ChangePassword(c *gin.Context)         {}
func (h *UserHandler) GetUserStats(c *gin.Context)           {}

// NewOrderHandler placeholder
func NewOrderHandler(orderService interface{}) *OrderHandler {
	return &OrderHandler{}
}

type OrderHandler struct{}

func (h *OrderHandler) GetUserOrders(c *gin.Context)         {}
func (h *OrderHandler) GetOrderByID(c *gin.Context)          {}
func (h *OrderHandler) CreateOrder(c *gin.Context)           {}
func (h *OrderHandler) CancelOrder(c *gin.Context)           {}
func (h *OrderHandler) UpdateOrderStatus(c *gin.Context)     {}
func (h *OrderHandler) GetSellerOrders(c *gin.Context)       {}
func (h *OrderHandler) GetAllOrders(c *gin.Context)          {}

// NewPaymentHandler placeholder
func NewPaymentHandler(paymentService interface{}) *PaymentHandler {
	return &PaymentHandler{}
}

type PaymentHandler struct{}

func (h *PaymentHandler) GetPaymentHistory(c *gin.Context)   {}
func (h *PaymentHandler) CreatePaymentIntent(c *gin.Context) {}
func (h *PaymentHandler) ConfirmPayment(c *gin.Context)      {}
func (h *PaymentHandler) HandleStripeWebhook(c *gin.Context) {}
func (h *PaymentHandler) HandlePayPalWebhook(c *gin.Context) {}

// NewUploadHandler placeholder
func NewUploadHandler() *UploadHandler {
	return &UploadHandler{}
}

type UploadHandler struct{}

func (h *UploadHandler) UploadImage(c *gin.Context)          {}
func (h *UploadHandler) UploadMultipleImages(c *gin.Context) {}
func (h *UploadHandler) GetUserImages(c *gin.Context)        {}
func (h *UploadHandler) GetImageInfo(c *gin.Context)         {}
func (h *UploadHandler) DeleteImage(c *gin.Context)          {}
func (h *UploadHandler) HandleCloudinaryWebhook(c *gin.Context) {}

// NewNotificationHandler placeholder
func NewNotificationHandler(notificationService interface{}) *NotificationHandler {
	return &NotificationHandler{}
}

type NotificationHandler struct{}

func (h *NotificationHandler) GetUserNotifications(c *gin.Context)  {}
func (h *NotificationHandler) MarkAsRead(c *gin.Context)            {}
func (h *NotificationHandler) MarkAllAsRead(c *gin.Context)         {}
func (h *NotificationHandler) GetUnreadCount(c *gin.Context)        {}

// NewAdminHandler placeholder
func NewAdminHandler(adminService interface{}) *AdminHandler {
	return &AdminHandler{}
}

type AdminHandler struct{}

func (h *AdminHandler) GetDashboard(c *gin.Context)          {}
func (h *AdminHandler) GetDashboardStats(c *gin.Context)     {}
func (h *AdminHandler) GetUsers(c *gin.Context)              {}
func (h *AdminHandler) UpdateUserStatus(c *gin.Context)      {}
func (h *AdminHandler) GetSystemLogs(c *gin.Context)         {}
func (h *AdminHandler) GetAllOrders(c *gin.Context)          {}