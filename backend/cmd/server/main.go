package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nawthtech/nawthtech/backend/internal/cloudinary"
 "github.com/nawthtech/nawthtech/backend/internal/mongodb"
 "github.com/nawthtech/nawthtech/backend/internal/cloudflare"
	"github.com/nawthtech/nawthtech/backend/internal/config"
	"github.com/nawthtech/nawthtech/backend/internal/handlers"
	"github.com/nawthtech/nawthtech/backend/internal/logger"
	"github.com/nawthtech/nawthtech/backend/internal/middleware"
	"github.com/nawthtech/nawthtech/backend/internal/services"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// تحميل الإعدادات
	cfg := config.Load()
	logger.Stdout.Info("🚀 بدء تشغيل تطبيق نوذ تك", 
		"environment", cfg.Environment,
		"version", cfg.Version,
	)

// في دالة main
cloudflareService, err := cloudflare.InitCloudflareService()
if err != nil {
    logger.Error(context.Background(), "⚠️ فشل في تهيئة Cloudflare", "error", err.Error())
} else {
    logger.Info(context.Background(), "✅ Cloudflare جاهز للاستخدام")
}

	// تهيئة قاعدة بيانات MongoDB
	mongoClient, err := initMongoDB(cfg)
	if err != nil {
		logger.Stderr.Error("❌ فشل في تهيئة قاعدة البيانات", logger.ErrAttr(err))
		os.Exit(1)
	}
	defer closeMongoDB(mongoClient)

	// تهيئة خدمة Cloudinary
	cloudinaryService, err := initCloudinary(cfg)
	if err != nil {
		logger.Stderr.Error("❌ فشل في تهيئة خدمة Cloudinary", logger.ErrAttr(err))
		// لا نوقف التطبيق إذا فشل Cloudinary، يمكن أن يعمل بدونها
	} else {
		logger.Stdout.Info("✅ تم تهيئة خدمة Cloudinary بنجاح")
	}

	// إنشاء حاوية الخدمات مع MongoDB
	serviceContainer := services.NewServiceContainer(mongoClient, cfg.Database.Name)

	// إنشاء تطبيق Gin
	app := initGinApp(cfg)

	// تسجيل جميع الوسائط
	registerMiddlewares(app, cfg)

	// تسجيل جميع المسارات مع تمرير Cloudinary service
	registerAllRoutes(app, serviceContainer, cfg, mongoClient, cloudinaryService)

	// بدء الخادم
	startServer(app, cfg)
}

// initMongoDB تهيئة اتصال MongoDB
func initMongoDB(cfg *config.Config) (*mongo.Client, error) {
	logger.Stdout.Info("🗄️  تهيئة اتصال MongoDB...")

	// استخدام رابط الاتصال من الإعدادات
	connectionString := cfg.Database.URL
	if cfg.IsDevelopment() && connectionString == "" {
		connectionString = "mongodb://localhost:27017/nawthtech"
		logger.Stdout.Info("🔧 استخدام إعدادات MongoDB افتراضية للتطوير")
	}

	// إعداد خيارات العميل
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().
		ApplyURI(connectionString).
		SetServerAPIOptions(serverAPI).
		SetMaxPoolSize(100).
		SetMinPoolSize(10).
		SetConnectTimeout(10 * time.Second).
		SetSocketTimeout(30 * time.Second).
		SetServerSelectionTimeout(10 * time.Second)

	// الاتصال بقاعدة البيانات
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	// اختبار الاتصال
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	logger.Stdout.Info("✅ تم الاتصال بـ MongoDB بنجاح",
		"database", cfg.Database.Name,
		"connection_string", maskConnectionString(connectionString),
	)
	return client, nil
}

// initCloudinary تهيئة خدمة Cloudinary
func initCloudinary(cfg *config.Config) (*cloudinary.CloudinaryService, error) {
	logger.Stdout.Info("☁️  تهيئة خدمة Cloudinary...")

	service, err := cloudinary.NewCloudinaryService()
	if err != nil {
		return nil, err
	}

	logger.Stdout.Info("✅ تم تهيئة Cloudinary بنجاح",
		"cloud_name", os.Getenv("CLOUDINARY_CLOUD_NAME"),
		"environment", cfg.Environment,
	)
	return service, nil
}

// closeMongoDB إغلاق اتصال MongoDB
func closeMongoDB(client *mongo.Client) {
	if client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		
		err := client.Disconnect(ctx)
		if err != nil {
			logger.Stderr.Error("❌ فشل في إغلاق اتصال MongoDB", logger.ErrAttr(err))
		} else {
			logger.Stdout.Info("✅ تم إغلاق اتصال MongoDB")
		}
	}
}

// initGinApp تهيئة تطبيق Gin
func initGinApp(cfg *config.Config) *gin.Engine {
	// تعيين وضع Gin
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	app := gin.New()

	// إعدادات Gin الأساسية
	app.ForwardedByClientIP = true
	
	// تعيين الوكائل الموثوق بها بناءً على البيئة
	if cfg.IsProduction() {
		app.SetTrustedProxies([]string{
			"127.0.0.1",
			"::1",
			"10.0.0.0/8",
			"172.16.0.0/12", 
			"192.168.0.0/16",
		})
	} else {
		app.SetTrustedProxies([]string{"127.0.0.1", "::1"})
	}

	// زيادة حجم الرفع الافتراضي إلى 10MB لاستيعاب الصور
	app.MaxMultipartMemory = 10 << 20 // 10 MB

	return app
}

// registerMiddlewares تسجيل الوسائط
func registerMiddlewares(app *gin.Engine, cfg *config.Config) {
	// ✅ وسيط CORS - يتم تطبيقه أولاً
	app.Use(middleware.CORS())

	// ✅ وسيط رؤوس الأمان
	app.Use(middleware.SecurityHeaders())

	// ✅ وسيط التسجيل المخصص
	app.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		logger.Stdout.Info("طلب HTTP",
			"method", param.Method,
			"path", param.Path,
			"status", param.StatusCode,
			"latency", param.Latency,
			"client_ip", param.ClientIP,
			"user_agent", param.Request.UserAgent(),
		)
		return ""
	}))

	// ✅ وسيط الاستعادة من الأخطاء
	app.Use(gin.Recovery())

	// ✅ وسيط تحديد المعدل
	app.Use(middleware.RateLimit())

	logger.Stdout.Info("✅ تم تسجيل الوسائط الأساسية",
		"cors_enabled", true,
		"security_headers", true,
		"rate_limiting", true,
		"max_upload_size", "10MB",
	)
}

// registerAllRoutes تسجيل جميع المسارات مع دعم Cloudinary
func registerAllRoutes(
	app *gin.Engine, 
	serviceContainer *services.ServiceContainer, 
	cfg *config.Config, 
	mongoClient *mongo.Client,
	cloudinaryService *cloudinary.CloudinaryService,
) {
	logger.Stdout.Info("🛣️  تسجيل مسارات التطبيق...")

	// ✅ مجموعة API الأساسية
	api := app.Group("/api/v1")

	// ✅ مسارات المصادقة
	authHandler := handlers.NewAuthHandler(serviceContainer.AuthService)
	authRoutes := api.Group("/auth")
	{
		authRoutes.POST("/register", authHandler.Register)
		authRoutes.POST("/login", authHandler.Login)
		authRoutes.POST("/logout", authHandler.Logout)
		authRoutes.POST("/refresh-token", authHandler.RefreshToken)
		authRoutes.POST("/forgot-password", authHandler.ForgotPassword)
		authRoutes.POST("/reset-password", authHandler.ResetPassword)
		authRoutes.GET("/verify-token", authHandler.VerifyToken)
	}

	// ✅ مسارات المستخدم
	userHandler := handlers.NewUserHandler(serviceContainer.UserService)
	userRoutes := api.Group("/users")
	{
		userRoutes.GET("/profile", userHandler.GetProfile)
		userRoutes.PUT("/profile", userHandler.UpdateProfile)
		userRoutes.PUT("/change-password", userHandler.ChangePassword)
		userRoutes.GET("/stats", userHandler.GetUserStats)
	}

	// ✅ مسارات الخدمات
	serviceHandler := handlers.NewServiceHandler(serviceContainer.ServiceService)
	serviceRoutes := api.Group("/services")
	{
		serviceRoutes.GET("/", serviceHandler.GetServices)
		serviceRoutes.GET("/search", serviceHandler.SearchServices)
		serviceRoutes.GET("/featured", serviceHandler.GetFeaturedServices)
		serviceRoutes.GET("/categories", serviceHandler.GetCategories)
		serviceRoutes.GET("/my-services", serviceHandler.GetMyServices)
		serviceRoutes.POST("/", serviceHandler.CreateService)
		serviceRoutes.GET("/:id", serviceHandler.GetServiceByID)
		serviceRoutes.PUT("/:id", serviceHandler.UpdateService)
		serviceRoutes.DELETE("/:id", serviceHandler.DeleteService)
	}

	// ✅ مسارات الفئات
	categoryHandler := handlers.NewCategoryHandler(serviceContainer.CategoryService)
	categoryRoutes := api.Group("/categories")
	{
		categoryRoutes.GET("/", categoryHandler.GetCategories)
		categoryRoutes.POST("/", categoryHandler.CreateCategory)
		categoryRoutes.GET("/:id", categoryHandler.GetCategoryByID)
		categoryRoutes.PUT("/:id", categoryHandler.UpdateCategory)
		categoryRoutes.DELETE("/:id", categoryHandler.DeleteCategory)
	}

	// ✅ مسارات الطلبات
	orderHandler := handlers.NewOrderHandler(serviceContainer.OrderService)
	orderRoutes := api.Group("/orders")
	{
		orderRoutes.GET("/", orderHandler.GetUserOrders)
		orderRoutes.POST("/", orderHandler.CreateOrder)
		orderRoutes.GET("/:id", orderHandler.GetOrderByID)
		orderRoutes.PUT("/:id/status", orderHandler.UpdateOrderStatus)
		orderRoutes.DELETE("/:id", orderHandler.CancelOrder)
	}

	// ✅ مسارات الدفع
	paymentHandler := handlers.NewPaymentHandler(serviceContainer.PaymentService)
	paymentRoutes := api.Group("/payments")
	{
		paymentRoutes.GET("/history", paymentHandler.GetPaymentHistory)
		paymentRoutes.POST("/create-intent", paymentHandler.CreatePaymentIntent)
		paymentRoutes.POST("/confirm", paymentHandler.ConfirmPayment)
	}

	// ✅ مسارات الرفع - Cloudinary Integration
	var uploadHandler handlers.UploadHandler
	if cloudinaryService != nil {
		// استخدام Cloudinary إذا كان متاحاً
		uploadHandler = handlers.NewUploadHandlerWithService(cloudinaryService)
		logger.Stdout.Info("✅ تم تسجيل مسارات الرفع مع Cloudinary")
	} else {
		// استخدام خدمة الرفع الأساسية إذا فشل Cloudinary
		uploadHandler = handlers.NewUploadHandlerWithService(nil)
		logger.Stdout.Warn("⚠️  تم تسجيل مسارات الرفع بدون Cloudinary - باستخدام وضع أساسي")
	}

	uploadRoutes := api.Group("/upload")
	{
		uploadRoutes.POST("/image", uploadHandler.UploadImage)
		uploadRoutes.POST("/images", uploadHandler.UploadMultipleImages)
		uploadRoutes.GET("/image/:public_id", uploadHandler.GetImageInfo)
		uploadRoutes.DELETE("/image/:public_id", uploadHandler.DeleteImage)
		uploadRoutes.GET("/my-images", uploadHandler.GetUserImages)
	}

	// ✅ مسارات الإشعارات
	notificationHandler := handlers.NewNotificationHandler(serviceContainer.NotificationService)
	notificationRoutes := api.Group("/notifications")
	{
		notificationRoutes.GET("/", notificationHandler.GetUserNotifications)
		notificationRoutes.PUT("/:id/read", notificationHandler.MarkAsRead)
		notificationRoutes.PUT("/read-all", notificationHandler.MarkAllAsRead)
		notificationRoutes.GET("/unread-count", notificationHandler.GetUnreadCount)
	}

	// ✅ مسارات الإدارة
	adminHandler := handlers.NewAdminHandler(serviceContainer.AdminService)
	adminRoutes := api.Group("/admin")
	{
		adminRoutes.GET("/dashboard", adminHandler.GetDashboard)
		adminRoutes.GET("/dashboard/stats", adminHandler.GetDashboardStats)
		adminRoutes.GET("/users", adminHandler.GetUsers)
		adminRoutes.PUT("/users/:id/status", adminHandler.UpdateUserStatus)
		adminRoutes.GET("/system-logs", adminHandler.GetSystemLogs)
	}

	// ✅ مسارات الصحة والفحص
	api.GET("/health", func(c *gin.Context) {
		// فحص اتصال MongoDB
		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		mongoStatus := "connected"
		if err := mongoClient.Ping(ctx, nil); err != nil {
			mongoStatus = "disconnected"
		}

		// فحص حالة Cloudinary
		cloudinaryStatus := "not_configured"
		if cloudinaryService != nil {
			cloudinaryStatus = "connected"
		}

		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"timestamp": time.Now().UTC(),
			"services": gin.H{
				"mongodb": gin.H{
					"status": mongoStatus,
					"database": cfg.Database.Name,
				},
				"cloudinary": gin.H{
					"status": cloudinaryStatus,
				},
			},
			"version":     cfg.Version,
			"environment": cfg.Environment,
		})
	})

	// ✅ مسار الصفحة الرئيسية
	app.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message":     "مرحباً بك في نوذ تك - منصة الخدمات الإلكترونية",
			"version":     cfg.Version,
			"environment": cfg.Environment,
			"timestamp":   time.Now().UTC(),
			"database":    "MongoDB",
			"upload_service": "Cloudinary",
			"status":      "running",
		})
	})

	logger.Stdout.Info("✅ تم تسجيل جميع المسارات بنجاح",
		"total_endpoints", 45, // تقديري لعدد النقاط الطرفية
		"cloudinary_enabled", cloudinaryService != nil,
		"api_version", "v1",
	)
}

// startServer بدء الخادم
func startServer(app *gin.Engine, cfg *config.Config) {
	// إعداد الخادم
	server := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           app,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		IdleTimeout:       120 * time.Second,
		MaxHeaderBytes:    1 << 20, // 1MB
	}

	// قناة لاستقبال إشارات النظام
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	// بدء الخادم في goroutine
	go func() {
		logger.Stdout.Info("🌐 بدء تشغيل الخادم",
			"port", cfg.Port,
			"environment", cfg.Environment,
			"version", cfg.Version,
			"database", "MongoDB",
			"upload_service", "Cloudinary",
		)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Stderr.Error("❌ فشل في بدء الخادم", logger.ErrAttr(err))
			os.Exit(1)
		}
	}()

	// انتظار إشارة الإغلاق
	sig := <-sigChan
	logger.Stdout.Info("🛑 استلام إشارة إغلاق", 
		"signal", sig.String(),
	)

	// إيقاف الخادم بشكل أنيق
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Stderr.Error("❌ فشل في إيقاف الخادم بشكل أنيق", logger.ErrAttr(err))
	} else {
		logger.Stdout.Info("✅ تم إيقاف الخادم بنجاح")
	}
}

// maskConnectionString إخفاء كلمة السر في رابط الاتصال للأمان
func maskConnectionString(connectionString string) string {
	// إخفاء كلمة السر لعرض آمن في السجلات
	// مثال: mongodb://user:password@host -> mongodb://user:****@host
	if len(connectionString) > 50 {
		return connectionString[:30] + "****" + connectionString[len(connectionString)-20:]
	}
	return "****"
}