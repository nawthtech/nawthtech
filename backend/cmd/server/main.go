package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	v1shared "github.com/nawthtech/nawthtech/backend/api/v1"
	"github.com/nawthtech/nawthtech/backend/api/v1/routes"
	"github.com/nawthtech/nawthtech/backend/internal/cloudflare"
	"github.com/nawthtech/nawthtech/backend/internal/cloudinary"
	"github.com/nawthtech/nawthtech/backend/internal/config"
	"github.com/nawthtech/nawthtech/backend/internal/email"
	"github.com/nawthtech/nawthtech/backend/internal/handlers"
	"github.com/nawthtech/nawthtech/backend/internal/middleware"
	"github.com/nawthtech/nawthtech/backend/internal/mongodb"
	"github.com/nawthtech/nawthtech/backend/internal/services"
)

// initLogger تهيئة logger
func initLogger() {
	// إذا كان logger الافتراضي ليس لديه handler، قم بتهيئته
	if slog.Default().Handler() == nil {
		handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
		slog.SetDefault(slog.New(handler))
	}
}

// Run تشغيل خادم API
func Run() error {
	// ✅ تهيئة logger أولاً
	initLogger()

	// تحميل الإعدادات
	cfg := config.Load()

	// تسجيل بدء التشغيل
	slog.Info("🚀 بدء تشغيل خادم نوذ تك",
		"environment", cfg.Environment,
		"version", cfg.Version,
		"port", cfg.Port,
	)

	// ================================
	// 🔄 تهيئة جميع الخدمات
	// ================================

	// 1. 📧 تهيئة خدمة البريد الإلكتروني
	emailService, err := email.NewEmailService()
	if err != nil {
		slog.Warn("⚠️ فشل في تهيئة خدمة البريد الإلكتروني", "error", err)
	} else {
		slog.Info("✅ خدمة البريد الإلكتروني جاهزة للاستخدام",
			"enabled", email.IsEnabled(),
		)
	}

	// 2. 🌐 تهيئة خدمة Cloudflare
	cloudflareService, err := cloudflare.InitCloudflareService()
	if err != nil {
		slog.Warn("⚠️ فشل في تهيئة Cloudflare", "error", err)
	} else {
		slog.Info("✅ Cloudflare جاهز للاستخدام",
			"enabled", cloudflare.IsEnabled(),
		)
	}

	// 3. 🗄️ تهيئة قاعدة بيانات MongoDB
	mongoService, err := mongodb.NewMongoDBService()
	if err != nil {
		slog.Error("❌ فشل في تهيئة قاعدة البيانات", "error", err)
		return err
	}
	defer mongoService.Close()

	// 4. ☁️ تهيئة خدمة Cloudinary
	cloudinaryService, err := cloudinary.NewCloudinaryService()
	if err != nil {
		slog.Warn("❌ فشل في تهيئة خدمة Cloudinary", "error", err)
		// لا نوقف التطبيق إذا فشل Cloudinary، يمكن أن يعمل بدونها
	} else {
		slog.Info("✅ تم تهيئة خدمة Cloudinary بنجاح")
	}

	// ================================
	// 🏗️ بناء التطبيق
	// ================================

	// إنشاء حاوية الخدمات مع MongoDB
	serviceContainer := services.NewServiceContainer(mongoService.GetClient(), mongoService.Config.DatabaseName)

	// إنشاء تطبيق Gin
	app := initGinApp(cfg)

	// تسجيل جميع الوسائط
	registerMiddlewares(app, cfg)

	// تسجيل جميع المسارات
	registerAllRoutes(app, serviceContainer, cfg, mongoService, cloudinaryService, cloudflareService, emailService)

	// بدء الخادم
	return startServer(app, cfg)
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
	app.Use(middleware.CORSMiddleware())

	// ✅ وسيط رؤوس الأمان
	app.Use(middleware.SecurityHeaders())

	// ✅ وسيط التسجيل المخصص
	app.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		slog.Info("طلب HTTP",
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
	app.Use(middleware.RateLimitMiddlewareFunc())

	// ✅ وسيط إصدار API
	app.Use(v1shared.APIVersionMiddleware())

	// ✅ وسيط استجابة API الموحدة
	app.Use(v1shared.APIResponseMiddleware())

	slog.Info("✅ تم تسجيل الوسائط الأساسية",
		"cors_enabled", true,
		"security_headers", true,
		"rate_limiting", true,
		"api_versioning", true,
		"max_upload_size", "10MB",
	)
}

// registerAllRoutes تسجيل جميع المسارات
func registerAllRoutes(
	app *gin.Engine,
	serviceContainer *services.ServiceContainer,
	cfg *config.Config,
	mongoService *mongodb.MongoDBService,
	cloudinaryService *cloudinary.CloudinaryService,
	cloudflareService *cloudflare.CloudflareConfig,
	emailService *email.Office365Config,
) {
	slog.Info("🛣️  تسجيل مسارات التطبيق...")

	// ✅ إنشاء حاوية المعاجل
	handlerContainer := &routes.HandlerContainer{
		Auth:         handlers.NewAuthHandler(serviceContainer.Auth),
		User:         handlers.NewUserHandler(serviceContainer.User),
		Service:      handlers.NewServiceHandler(serviceContainer.Service),
		Category:     handlers.NewCategoryHandler(serviceContainer.Category),
		Order:        handlers.NewOrderHandler(serviceContainer.Order),
		Payment:      handlers.NewPaymentHandler(serviceContainer.Payment),
		Notification: handlers.NewNotificationHandler(serviceContainer.Notification),
		Admin:        handlers.NewAdminHandler(serviceContainer.Admin),
	}

	// ✅ تهيئة معالج الرفع مع Cloudinary
	if cloudinaryService != nil {
		uploadHandler := handlers.NewUploadHandlerWithService(cloudinaryService)
		handlerContainer.Upload = uploadHandler
	} else {
		// إنشاء معالج رفع بدون Cloudinary (للحالات الطارئة)
		uploadHandler, err := handlers.NewUploadHandler()
		if err != nil {
			slog.Error("❌ فشل في إنشاء معالج الرفع الافتراضي", "error", err)
		} else {
			handlerContainer.Upload = uploadHandler
		}
	}

	// ✅ تسجيل مسارات API v1
	apiGroup := app.Group("/api")
	v1Group := apiGroup.Group("/v1")
	routes.RegisterV1Routes(v1Group, handlerContainer, v1shared.AuthMiddleware())

	// ✅ تسجيل مسارات الصحة والفحص
	registerHealthRoutes(app, mongoService, cloudinaryService, cloudflareService, emailService, cfg)

	// ✅ تسجيل المسارات العامة
	registerGeneralRoutes(app, cfg)

	slog.Info("✅ تم تسجيل جميع المسارات بنجاح",
		"api_version", "v1",
		"cloudinary_enabled", cloudinaryService != nil,
		"cloudflare_enabled", cloudflare.IsEnabled(),
		"email_enabled", email.IsEnabled(),
	)
}

// registerHealthRoutes تسجيل مسارات الصحة والفحص
func registerHealthRoutes(
	app *gin.Engine,
	mongoService *mongodb.MongoDBService,
	cloudinaryService *cloudinary.CloudinaryService,
	cloudflareService *cloudflare.CloudflareConfig,
	emailService *email.Office365Config,
	cfg *config.Config,
) {
	// ✅ مسار الصحة الشامل
	app.GET("/health", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
		defer cancel()

		// فحص جميع الخدمات
		mongoStatus := mongoService.HealthCheck(ctx)
		cloudinaryStatus := "not_configured"
		if cloudinaryService != nil {
			cloudinaryStatus = "connected"
		}
		cloudflareStatus := cloudflare.HealthCheck()
		emailStatus := email.HealthCheck()

		response := gin.H{
			"status":      "healthy",
			"timestamp":   time.Now().UTC(),
			"version":     cfg.Version,
			"environment": cfg.Environment,
			"services": gin.H{
				"mongodb": gin.H{
					"status":   mongoStatus["status"],
					"database": mongoService.Config.DatabaseName,
				},
				"cloudinary": gin.H{
					"status": cloudinaryStatus,
				},
				"cloudflare": cloudflareStatus,
				"email":      emailStatus,
			},
		}

		// تحديد الحالة العامة بناءً على الخدمات الأساسية
		if mongoStatus["status"] != "healthy" {
			response["status"] = "degraded"
			response["message"] = "بعض الخدمات غير متاحة"
		}

		c.JSON(http.StatusOK, response)
	})

	// ✅ مسار الصحة البسيط (ل Load Balancers)
	app.GET("/health/live", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "live",
			"timestamp": time.Now().UTC(),
		})
	})

	// ✅ مسار الجاهزية
	app.GET("/health/ready", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		// التحقق من اتصال MongoDB فقط (الخدمة الأساسية)
		mongoStatus := mongoService.HealthCheck(ctx)

		if mongoStatus["status"] != "healthy" {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status":    "not_ready",
				"timestamp": time.Now().UTC(),
				"error":     "قاعدة البيانات غير متاحة",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":    "ready",
			"timestamp": time.Now().UTC(),
		})
	})
}

// registerGeneralRoutes تسجيل المسارات العامة
func registerGeneralRoutes(app *gin.Engine, cfg *config.Config) {
	// ✅ مسار الصفحة الرئيسية
	app.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message":       "مرحباً بك في نوذ تك - منصة الخدمات الإلكترونية",
			"version":       cfg.Version,
			"environment":   cfg.Environment,
			"timestamp":     time.Now().UTC(),
			"documentation": "/api/v1/docs",
			"health_check":  "/health",
			"services": gin.H{
				"database":       "MongoDB",
				"upload_service": "Cloudinary",
				"cdn":            "Cloudflare",
				"email":          "Office 365",
			},
		})
	})

	// ✅ مسار معلومات النظام
	app.GET("/info", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"name":        "NawthTech Backend API",
			"version":     cfg.Version,
			"environment": cfg.Environment,
			"status":      "running",
			"timestamp":   time.Now().UTC(),
			"endpoints": gin.H{
				"api_v1":        "/api/v1",
				"health":        "/health",
				"documentation": "/api/v1/docs",
			},
			"features": []string{
				"المصادقة الآمنة",
				"إدارة المستخدمين",
				"الخدمات الإلكترونية",
				"نظام الطلبات والدفع",
				"رفع الملفات مع Cloudinary",
				"CDN مع Cloudflare",
				"إرسال البريد مع Office 365",
			},
		})
	})
}

// startServer بدء الخادم
func startServer(app *gin.Engine, cfg *config.Config) error {
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
		slog.Info("🌐 بدء تشغيل الخادم",
			"port", cfg.Port,
			"environment", cfg.Environment,
			"version", cfg.Version,
			"services", []string{
				"MongoDB",
				"Cloudinary",
				"Cloudflare",
				"Office 365",
			},
			"read_timeout", "30s",
			"write_timeout", "30s",
			"idle_timeout", "120s",
		)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("❌ فشل في بدء الخادم", "error", err)
			os.Exit(1)
		}
	}()

	// انتظار إشارة الإغلاق
	sig := <-sigChan
	slog.Info("🛑 استلام إشارة إغلاق",
		"signal", sig.String(),
	)

	// إيقاف الخادم بشكل أنيق
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	slog.Info("⏳ إيقاف الخادم بشكل أنيق...",
		"timeout", "30s",
	)

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("❌ فشل في إيقاف الخادم بشكل أنيق", "error", err)
		return err
	}

	slog.Info("✅ تم إيقاف الخادم بنجاح",
		"duration", "أنيق",
	)

	return nil
}

// main الدالة الرئيسية
func main() {
	// ✅ تهيئة logger أولاً
	initLogger()

	if err := Run(); err != nil {
		slog.Error("❌ فشل في تشغيل الخادم", "error", err)
		os.Exit(1)
	}
}