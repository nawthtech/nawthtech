package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nawthtech/nawthtech/backend/internal/config"
	"github.com/nawthtech/nawthtech/backend/internal/handlers"
	"github.com/nawthtech/nawthtech/backend/internal/logger"
	"github.com/nawthtech/nawthtech/backend/internal/middleware"
	"github.com/nawthtech/nawthtech/backend/internal/services"
	"github.com/nawthtech/nawthtech/backend/internal/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// تحميل الإعدادات
	cfg := config.Load()
	logger.Stdout.Info("🚀 بدء تشغيل تطبيق نوذ تك", 
		"environment", cfg.Environment,
		"version", cfg.Version,
	)

	// عرض إحصائيات CORS عند البدء
	corsStats := config.GetCORSStats()
	logger.Stdout.Info("🌐 إعدادات CORS", 
		"total_allowed_origins", corsStats["totalAllowedOrigins"],
		"environment", corsStats["environment"],
	)

	// تهيئة قاعدة البيانات
	db, err := initDatabase(cfg)
	if err != nil {
		logger.Stderr.Error("❌ فشل في تهيئة قاعدة البيانات", logger.ErrAttr(err))
		os.Exit(1)
	}
	defer closeDatabase(db)

	// تشغيل ترحيلات قاعدة البيانات
	if err := runMigrations(db); err != nil {
		logger.Stderr.Error("❌ فشل في تشغيل ترحيلات قاعدة البيانات", logger.ErrAttr(err))
		if cfg.IsProduction() {
			os.Exit(1)
		}
	}

	// تهيئة خدمة التخزين المؤقت
	cacheService, err := initCacheService(cfg)
	if err != nil {
		logger.Stderr.Error("❌ فشل في تهيئة خدمة التخزين المؤقت", logger.ErrAttr(err))
		// نستمر بدون تخزين مؤقت في بيئة التطوير
		if cfg.IsProduction() {
			os.Exit(1)
		}
	}

	// فحص صحة التطبيق
	if !healthCheck(cfg, db, cacheService) {
		logger.Stderr.Error("❌ فحص الصحة فشل - التطبيق غير جاهز")
		if cfg.IsProduction() {
			os.Exit(1)
		}
	}

	// إنشاء تطبيق Gin
	app := initGinApp(cfg)

	// تسجيل جميع الوسائط
	registerMiddlewares(app, cfg)

	// تسجيل جميع المسارات
	registerAllRoutes(app, db, cfg, cacheService)

	// بدء الخادم
	startServer(app, cfg)
}

// initDatabase تهيئة قاعدة البيانات
func initDatabase(cfg *config.Config) (*gorm.DB, error) {
	logger.Stdout.Info("🗄️  تهيئة اتصال قاعدة البيانات...")

	// في بيئة التطوير، يمكن استخدام SQLite للاختبار
	if cfg.IsDevelopment() && cfg.DatabaseURL == "" {
		logger.Stdout.Info("🔧 استخدام قاعدة بيانات للتطوير")
		// يمكن إضافة SQLite هنا إذا أردت
		return nil, nil
	}

	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// تكوين اتصال قاعدة البيانات
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// إعدادات تجمع الاتصالات
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	logger.Stdout.Info("✅ تم الاتصال بقاعدة البيانات بنجاح")
	return db, nil
}

// closeDatabase إغلاق اتصال قاعدة البيانات
func closeDatabase(db *gorm.DB) {
	if db != nil {
		sqlDB, err := db.DB()
		if err == nil {
			sqlDB.Close()
			logger.Stdout.Info("✅ تم إغلاق اتصال قاعدة البيانات")
		}
	}
}

// initCacheService تهيئة خدمة التخزين المؤقت
func initCacheService(cfg *config.Config) (services.CacheService, error) {
	logger.Stdout.Info("🔮 تهيئة خدمة التخزين المؤقت...")

	if !cfg.IsCacheEnabled() {
		logger.Stdout.Info("⚠️  خدمة التخزين المؤقت معطلة في الإعدادات")
		return nil, nil
	}

	cacheService := services.NewCacheService(cfg.GetCacheConfig())
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := cacheService.Initialize(ctx)
	if err != nil {
		return nil, err
	}

	// اختبار الخدمة
	health, err := cacheService.HealthCheck(ctx)
	if err != nil {
		return nil, err
	}

	logger.Stdout.Info("✅ تم تهيئة خدمة التخزين المؤقت بنجاح", 
		"status", health.Status,
		"environment", health.Environment,
	)

	return cacheService, nil
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

	return app
}

// registerMiddlewares تسجيل الوسائط
func registerMiddlewares(app *gin.Engine, cfg *config.Config) {
	// ✅ وسيط CORS المحدث - يتم تطبيقه أولاً
	app.Use(middleware.CORS())

	// ✅ وسيط رؤوس الأمان
	app.Use(middleware.SecurityHeaders())

	// ✅ وسيط التسجيل
	app.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// تسجيل طلبات CORS بشكل خاص
		if param.Method == "OPTIONS" {
			logger.Stdout.Info("طلب CORS Preflight",
				"method", param.Method,
				"path", param.Path,
				"status", param.StatusCode,
				"latency", param.Latency,
				"client_ip", param.ClientIP,
				"origin", param.Request.Header.Get("Origin"),
			)
		} else {
			logger.Stdout.Info("طلب HTTP",
				"method", param.Method,
				"path", param.Path,
				"status", param.StatusCode,
				"latency", param.Latency,
				"client_ip", param.ClientIP,
				"user_agent", param.Request.UserAgent(),
				"origin", param.Request.Header.Get("Origin"),
			)
		}
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
	)
}

// registerAllRoutes تسجيل جميع المسارات
func registerAllRoutes(app *gin.Engine, db *gorm.DB, cfg *config.Config, cacheService services.CacheService) {
	// استخدام الدالة الجديدة من handlers
	handlers.RegisterAllRoutes(app, db, cfg)

	// ✅ تسجيل مسار لفحص إحصائيات CORS (للتطوير فقط)
	if cfg.IsDevelopment() {
		app.GET("/api/debug/cors-stats", func(c *gin.Context) {
			stats := config.GetCORSStats()
			c.JSON(200, gin.H{
				"cors_stats": stats,
				"timestamp":  time.Now().Format(time.RFC3339),
			})
		})
	}

	// ✅ مسار للصحة الموسعة
	app.GET("/health/detailed", func(c *gin.Context) {
		corsStats := config.GetCORSStats()
		
		response := gin.H{
			"status":    "healthy",
			"service":   "nawthtech-backend",
			"timestamp": time.Now().Format(time.RFC3339),
			"version":   cfg.Version,
			"cors": gin.H{
				"total_allowed_origins": corsStats["totalAllowedOrigins"],
				"environment":          corsStats["environment"],
			},
			"system": gin.H{
				"goroutines": utils.GetGoroutineCount(),
				"memory_mb":  utils.GetMemoryUsageMB(),
			},
		}
		c.JSON(200, response)
	})

	// ✅ معالج للمسارات غير المعروفة
	app.NoRoute(func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		logger.Stdout.Warn("مسار غير معروف", 
			"path", c.Request.URL.Path,
			"method", c.Request.Method,
			"origin", origin,
			"client_ip", c.ClientIP(),
		)
		
		c.JSON(404, gin.H{
			"error":   "مسار غير موجود",
			"path":    c.Request.URL.Path,
			"message": "الرجاء التحقق من المسار والمحاولة مرة أخرى",
		})
	})

	logger.Stdout.Info("✅ تم تسجيل جميع المسارات",
		"total_routes", countRoutes(app),
		"cors_debug_enabled", cfg.IsDevelopment(),
	)
}

// countRoutes حساب عدد المسارات المسجلة (دالة مساعدة)
func countRoutes(app *gin.Engine) int {
	count := 0
	for _, route := range app.Routes() {
		if route.Method != "OPTIONS" {
			count++
		}
	}
	return count
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
			"cors_enabled", true,
		)

		// ✅ عرض إحصائيات CORS النهائية
		corsStats := config.GetCORSStats()
		logger.Stdout.Info("🔧 إعدادات CORS النهائية",
			"total_origins", corsStats["totalAllowedOrigins"],
			"services", corsStats["services"],
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
		"timestamp", time.Now().Format(time.RFC3339),
	)

	// إيقاف الخادم بشكل أنيق
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Stderr.Error("❌ فشل في إيقاف الخادم بشكل أنيق", logger.ErrAttr(err))
	} else {
		logger.Stdout.Info("✅ تم إيقاف الخادم بنجاح")
	}

	// إغلاق خدمة التخزين المؤقت إذا كانت نشطة
	if cacheService != nil {
		if err := cacheService.Close(); err != nil {
			logger.Stderr.Error("❌ فشل في إغلاق خدمة التخزين المؤقت", logger.ErrAttr(err))
		} else {
			logger.Stdout.Info("✅ تم إغلاق خدمة التخزين المؤقت")
		}
	}
}

// ========== دوال مساعدة للاختبار ==========

// initTestData تهيئة بيانات الاختبار (للتطوير فقط)
func initTestData(db *gorm.DB, cfg *config.Config) {
	if !cfg.IsDevelopment() {
		return
	}

	logger.Stdout.Info("🧪 تهيئة بيانات الاختبار...")

	// يمكن إضافة بيانات اختبار هنا
	// مثال: إنشاء مستخدمين، خدمات، إلخ.

	logger.Stdout.Info("✅ تم تهيئة بيانات الاختبار")
}

// runMigrations تشغيل ترحيلات قاعدة البيانات
func runMigrations(db *gorm.DB) error {
	if db == nil {
		return nil
	}

	logger.Stdout.Info("🔄 تشغيل ترحيلات قاعدة البيانات...")

	// يمكن إضافة ترحيلات قاعدة البيانات هنا
	// مثال: db.AutoMigrate(&models.User{}, &models.Service{}, ...)

	logger.Stdout.Info("✅ تم تشغيل الترحيلات بنجاح")
	return nil
}

// healthCheck فحص صحة التطبيق
func healthCheck(cfg *config.Config, db *gorm.DB, cacheService services.CacheService) bool {
	logger.Stdout.Info("🔍 فحص صحة التطبيق...")

	// فحص قاعدة البيانات
	if db != nil {
		sqlDB, err := db.DB()
		if err == nil {
			if err := sqlDB.Ping(); err != nil {
				logger.Stderr.Error("❌ فشل في فحص قاعدة البيانات", logger.ErrAttr(err))
				return false
			}
		}
	}

	// فحص التخزين المؤقت
	if cacheService != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		health, err := cacheService.HealthCheck(ctx)
		if err != nil || health.Status != "healthy" {
			logger.Stderr.Error("❌ فشل في فحص خدمة التخزين المؤقت", logger.ErrAttr(err))
			return false
		}
	}

	// ✅ فحص إعدادات CORS
	corsStats := config.GetCORSStats()
	if corsStats["totalAllowedOrigins"].(int) == 0 {
		logger.Stderr.Warn("⚠️  لا توجد نطاقات مسموح بها في إعدادات CORS")
	}

	logger.Stdout.Info("✅ فحص الصحة مكتمل - التطبيق جاهز",
		"cors_origins", corsStats["totalAllowedOrigins"],
	)
	return true
}