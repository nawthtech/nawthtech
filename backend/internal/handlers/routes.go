package handlers

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nawthtech/nawthtech/backend/internal/config"
	"github.com/nawthtech/nawthtech/backend/internal/middleware"
	"github.com/nawthtech/nawthtech/backend/internal/services"
	"github.com/nawthtech/nawthtech/backend/internal/ai
	"go.mongodb.org/mongo-driver/mongo"
)

// RegisterAllRoutes تسجيل جميع المسارات
func RegisterAllRoutes(router *gin.Engine, serviceContainer *services.ServiceContainer, config *config.Config, mongoClient *mongo.Client) {
	// تطبيق middleware العام على مستوى التطبيق
	applyGlobalMiddleware(router, config)

	// إنشاء حاوية الوسائط
	middlewares := initializeMiddlewares(serviceContainer, config)

	// مجموعة API الرئيسية
	api := router.Group("/api/v1")

	// ========== مسارات الصحة ==========
	registerHealthRoutes(router, config, mongoClient)

	// ========== المسارات العامة (لا تتطلب مصادقة) ==========
	registerPublicRoutes(api, serviceContainer, middlewares)

	// ========== المسارات المحمية (تتطلب مصادقة) ==========
	registerProtectedRoutes(api, serviceContainer, middlewares)

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
func registerPublicRoutes(api *gin.RouterGroup, services *services.ServiceContainer, middlewares *middleware.MiddlewareContainer) {
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
}

// registerProtectedRoutes تسجيل المسارات المحمية
func registerProtectedRoutes(api *gin.RouterGroup, services *services.ServiceContainer, middlewares *middleware.MiddlewareContainer) {
	protected := api.Group("")
	// ملاحظة: سيتم تفعيل المصادقة لاحقاً
	// protected.Use(middlewares.AuthMiddleware)

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
	uploadHandler, _ := NewUploadHandler() // تجاهل الخطأ مؤقتاً
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
}

// registerAdminRoutes تسجيل مسارات المسؤولين
func registerAdminRoutes(api *gin.RouterGroup, services *services.ServiceContainer, middlewares *middleware.MiddlewareContainer) {
	admin := api.Group("/admin")
	// ملاحظة: سيتم تفعيل المصادقة لاحقاً
	// admin.Use(middlewares.AuthMiddleware, middlewares.AdminMiddleware)

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
	admin.GET("/orders", orderHandler.GetUserOrders) // سيتم تحديثها لاحقاً
	admin.PUT("/orders/:id/status", orderHandler.UpdateOrderStatus)
}

// registerSellerRoutes تسجيل مسارات البائعين
func registerSellerRoutes(api *gin.RouterGroup, services *services.ServiceContainer, middlewares *middleware.MiddlewareContainer) {
	seller := api.Group("/seller")
	// ملاحظة: سيتم تفعيل المصادقة لاحقاً
	// seller.Use(middlewares.AuthMiddleware, middleware.SellerMiddleware())

	// معالج الخدمات (البائعين)
	serviceHandler := NewServiceHandler(services.Service)
	seller.POST("/services", serviceHandler.CreateService)
	seller.PUT("/services/:id", serviceHandler.UpdateService)
	seller.DELETE("/services/:id", serviceHandler.DeleteService)
	seller.GET("/services/my", serviceHandler.GetMyServices)

	// معالج الطلبات (البائعين)
	orderHandler := NewOrderHandler(services.Order)
	seller.GET("/orders", orderHandler.GetUserOrders) // سيتم تحديثها لاحقاً
	seller.PUT("/orders/:id/status", orderHandler.UpdateOrderStatus)
}

// registerWebhookRoutes تسجيل مسارات الويب هووك
func registerWebhookRoutes(api *gin.RouterGroup, services *services.ServiceContainer, middlewares *middleware.MiddlewareContainer) {
	webhook := api.Group("/webhook")
	{
		// ويب هووك الرفع (Cloudinary)
		uploadHandler, _ := NewUploadHandler() // تجاهل الخطأ مؤقتاً
		webhook.POST("/upload/cloudinary", uploadHandler.UploadImage)
	}
}

func SetupRouter(videoHandler *handlers.VideoHandler) *gin.Engine {
    r := gin.Default()
    
    // Video routes
    video := api.Group("/video")
    video.Use(middleware.Auth()) // تأمين endpoints الفيديو
    {
        video.POST("/generate", videoHandler.GenerateVideoHandler)
        video.GET("/status/:jobId", videoHandler.GetVideoStatusHandler)
        video.GET("/download/:jobId", videoHandler.DownloadVideoHandler)
        video.GET("/types", videoHandler.ListVideoTypesHandler)
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
