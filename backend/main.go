package main

import (
	"cmp"
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
"log"

	"github.com/nawthtech/nawthtech/backend/internal/config"
	"github.com/nawthtech/nawthtech/backend/internal/handlers"
	"github.com/nawthtech/nawthtech/backend/internal/logger"
	"github.com/nawthtech/nawthtech/backend/internal/middleware"
	"github.com/nawthtech/nawthtech/backend/internal/services"

	"github.com/go-chi/chi/v5"
)


func main() {
	// تكوين خدمة التخزين المؤقت
	cacheConfig := services.DefaultCacheConfig()
	cacheService := services.NewCacheService(cacheConfig)

	// تهيئة الخدمة
	ctx := context.Background()
	if err := cacheService.Initialize(ctx); err != nil {
		log.Printf("⚠️  فشل في تهيئة خدمة التخزين المؤقت: %v", err)
	}

	// استخدام الخدمة
	err := cacheService.Set(ctx, "user:123", map[string]interface{}{
		"name":  "أحمد",
		"email": "ahmed@example.com",
	}, 30*time.Minute)
	
	if err != nil {
		log.Printf("❌ خطأ في التخزين المؤقت: %v", err)
	}

	// فحص الصحة
	health, err := cacheService.HealthCheck(ctx)
	if err != nil {
		log.Printf("❌ خطأ في فحص الصحة: %v", err)
	} else {
		log.Printf("📊 حالة التخزين المؤقت: %s", health.Status)
	}
}

func main() {
	// تحميل الإعدادات
	cfg := config.Load()

	// إنشاء الخدمات
	adminService := services.NewAdminService()
	// TODO: إنشاء باقي الخدمات عند الحاجة
	userService := services.NewUserService()
	authService := services.NewAuthService()

	// تجميع الخدمات
	appServices := &handlers.Services{
		Admin: adminService,
		User:  userService,
		Auth:  authService,
		// TODO: إضافة باقي الخدمات
	}

	// إنشاء الموجه
	r := chi.NewRouter()

	// تسجيل الوسائط
	middleware.Register(r)

	// تسجيل المسارات
	handlers.Register(r, appServices)

	// إعداد الخادم
	port := cmp.Or(os.Getenv("PORT"), "3000")
	server := &http.Server{
		Addr:              ":" + port,
		Handler:           r,
		ReadTimeout:       5 * time.Minute,
		WriteTimeout:      5 * time.Minute,
		ReadHeaderTimeout: 30 * time.Second,
		IdleTimeout:       120 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}

	// بدء الخادم
	go func() {
		logger.Stdout.Info("بدء تشغيل الخادم", "port", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Stderr.Error("فشل في بدء الخادم", logger.ErrAttr(err))
			os.Exit(1)
		}
	}()

	// انتظار إشارة الإغلاق
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	logger.Stdout.Info("استلام إشارة إغلاق")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	if err := server.Shutdown(ctx); err != nil {
		logger.Stderr.Error("فشل في إيقاف الخادم", logger.ErrAttr(err))
	} else {
		logger.Stdout.Info("تم إيقاف الخادم بنجاح")
	}

func main() {
	// تهيئة قاعدة البيانات (مثال)
	var db *gorm.DB
	
	// إنشاء المستودع والخدمة
	repo := services.NewServicesRepository(db)
	service := services.NewServicesService(repo)
	
	// تهيئة Router
	router := gin.Default()
	
	// تسجيل المسارات
	api := router.Group("/api/v1")
	services.RegisterRoutes(api, service)
	
	// تشغيل الخادم
	log.Println("Starting server on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}