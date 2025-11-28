package main

import (
	"context"
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

	// تهيئة قاعدة البيانات
	db, err := initDatabase(cfg)
	if err != nil {
		logger.Stderr.Error("❌ فشل في تهيئة قاعدة البيانات", logger.ErrAttr(err))
		os.Exit(1)
	}
	defer closeDatabase(db)

	// إنشاء حاوية الخدمات
	serviceContainer := services.NewServiceContainer(db)

	// إنشاء تطبيق Gin
	app := initGinApp(cfg)

	// تسجيل جميع الوسائط
	registerMiddlewares(app, cfg)

	// تسجيل جميع المسارات
	handlers.RegisterAllRoutes(app, serviceContainer, cfg, db)

	// بدء الخادم
	startServer(app, cfg)
}

// initDatabase تهيئة قاعدة البيانات
func initDatabase(cfg *config.Config) (*gorm.DB, error) {
	logger.Stdout.Info("🗄️  تهيئة اتصال قاعدة البيانات...")

	// استخدام GetDSN بدلاً من DSN مباشرة
	dsn := cfg.GetDSN()
	if cfg.IsDevelopment() && dsn == "" {
		dsn = "host=localhost user=postgres password=postgres dbname=nawthtech port=5432 sslmode=disable"
		logger.Stdout.Info("🔧 استخدام إعدادات قاعدة بيانات افتراضية للتطوير")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
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