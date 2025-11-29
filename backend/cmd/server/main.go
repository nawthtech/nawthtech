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

	// تهيئة قاعدة بيانات MongoDB
	mongoClient, err := initMongoDB(cfg)
	if err != nil {
		logger.Stderr.Error("❌ فشل في تهيئة قاعدة البيانات", logger.ErrAttr(err))
		os.Exit(1)
	}
	defer closeMongoDB(mongoClient)

	// إنشاء حاوية الخدمات مع MongoDB
	serviceContainer := services.NewServiceContainer(mongoClient, cfg.Database.Name)

	// إنشاء تطبيق Gin
	app := initGinApp(cfg)

	// تسجيل جميع الوسائط
	registerMiddlewares(app, cfg)

	// تسجيل جميع المسارات
	handlers.RegisterAllRoutes(app, serviceContainer, cfg, mongoClient)

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
			"database", "MongoDB",
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