package main

import (
	"cmp"
	"context"
<<<<<<< HEAD
=======
	"fmt"
>>>>>>> 6ae4997234e28055b50691dc8128a6fe491e84d6
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

<<<<<<< HEAD
	"github.com/nawthtech/backend/internal/config"
	"github.com/nawthtech/backend/internal/handlers"
	"github.com/nawthtech/backend/internal/logger"
	"github.com/nawthtech/backend/internal/middleware"
	"github.com/nawthtech/backend/internal/services"
	"github.com/nawthtech/backend/internal/utils"

=======
>>>>>>> 6ae4997234e28055b50691dc8128a6fe491e84d6
	"github.com/go-chi/chi/v5"
)

func main() {
<<<<<<< HEAD
	// تحميل الإعدادات
	cfg := config.Load()

	// تهيئة النظام
	if err := initializeSystem(cfg); err != nil {
		logger.Stderr.Error("فشل في تهيئة النظام", logger.ErrAttr(err))
		os.Exit(1)
	}

	// إنشاء الموجه
	r := chi.NewRouter()

	// تسجيل الوسائط
	middleware.Register(r)

	// تهيئة الخدمات
	services := initializeServices()

	// تسجيل المسارات
	handlers.Register(r, services)
=======
	// إنشاء الموجه
	r := chi.NewRouter()

	// وسائط أساسية
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			next.ServeHTTP(w, r)
		})
	})

	// مسارات أساسية
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"message": "مرحباً بك في NawthTech API", "status": "success"}`))
	})

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"status": "healthy", "timestamp": "` + time.Now().Format(time.RFC3339) + `"}`))
	})

	r.Get("/api/version", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"version": "1.0.0", "name": "NawthTech Backend"}`))
	})
>>>>>>> 6ae4997234e28055b50691dc8128a6fe491e84d6

	// إعداد الخادم
	port := cmp.Or(os.Getenv("PORT"), "3000")
	server := &http.Server{
		Addr:              ":" + port,
		Handler:           r,
		ReadTimeout:       5 * time.Minute,
		WriteTimeout:      5 * time.Minute,
		ReadHeaderTimeout: 30 * time.Second,
		IdleTimeout:       120 * time.Second,
<<<<<<< HEAD
		MaxHeaderBytes:    1 << 20, // 1MB
	}

	// بدء الخادم في goroutine منفصلة
	go func() {
		logger.Stdout.Info("بدء تشغيل الخادم", 
			"port", port,
			"environment", cfg.Environment,
			"version", cfg.Version,
		)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Stderr.Error("فشل في بدء الخادم", logger.ErrAttr(err))
			os.Exit(1)
		}
	}()

	// انتظار إشارة الإغلاق
	gracefulShutdown(server)
}

// initializeSystem تهيئة مكونات النظام
func initializeSystem(cfg *config.Config) error {
	// تهيئة قاعدة البيانات إذا كان رابط قاعدة البيانات متوفراً
	if cfg.DatabaseURL != "" {
		if err := utils.InitDatabase(cfg.DatabaseURL); err != nil {
			return err
		}
	}

	// تهيئة المدقق
	if err := utils.InitValidator(); err != nil {
		return err
	}

	// التحقق من المساحات التخزينية
	if err := checkStorage(); err != nil {
		return err
	}

	logger.Stdout.Info("تم تهيئة النظام بنجاح")
	return nil
}

// initializeServices تهيئة جميع الخدمات
func initializeServices() *handlers.Services {
	return &handlers.Services{
		Admin:   services.NewAdminService(),
		User:    services.NewUserService(),
		Auth:    services.NewAuthService(),
		Store:   services.NewStoreService(),
		Cart:    services.NewCartService(),
		Payment: services.NewPaymentService(),
		AI:      services.NewAIService(),
		Email:   services.NewEmailService(),
		Upload:  services.NewUploadService(),
	}
}

// checkStorage التحقق من المساحات التخزينية
func checkStorage() error {
	// إنشاء المجلدات الضرورية إذا لم تكن موجودة
	dirs := []string{
		"./uploads",
		"./logs", 
		"./backups",
		"./temp",
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	return nil
}

// gracefulShutdown إغلاق النظام بشكل آمن
func gracefulShutdown(server *http.Server) {
	// إنشاء قناة للإشارات
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// انتظار الإشارة
	<-sigChan

	logger.Stdout.Info("استلام إشارة إغلاق، بدء الإغلاق الآمن")

	// إعطاء وقت للإغلاق الآمن
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// إيقاف الخادم
	if err := server.Shutdown(ctx); err != nil {
		logger.Stderr.Error("فشل في إيقاف الخادم بشكل آمن", logger.ErrAttr(err))
	} else {
		logger.Stdout.Info("تم إيقاف الخادم بشكل آمن")
	}

	// إغلاق اتصالات قاعدة البيانات
	utils.CloseDatabase()
=======
		MaxHeaderBytes:    1 << 20,
	}

	// بدء الخادم
	go func() {
		fmt.Printf("🚀 بدء تشغيل الخادم على port %s\n", port)
		fmt.Printf("📡 Health check: http://localhost:%s/health\n", port)
		fmt.Printf("🔗 API: http://localhost:%s/api/version\n", port)
		
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("❌ فشل في بدء الخادم: %v\n", err)
			os.Exit(1)
		}
	}()

	// انتظار إشارة الإغلاق
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	fmt.Println("🛑 استلام إشارة إغلاق، بدء الإغلاق الآمن...")
	
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	if err := server.Shutdown(ctx); err != nil {
		fmt.Printf("❌ فشل في إيقاف الخادم: %v\n", err)
	} else {
		fmt.Println("✅ تم إيقاف الخادم بنجاح")
	}
>>>>>>> 6ae4997234e28055b50691dc8128a6fe491e84d6
}