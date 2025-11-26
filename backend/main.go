package main

import (
	"cmp"
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nawthtech/nawthtech/backend/internal/config"
	"github.com/nawthtech/nawthtech/backend/internal/handlers"
	"github.com/nawthtech/nawthtech/backend/internal/logger"
	"github.com/nawthtech/nawthtech/backend/internal/middleware"
	"github.com/nawthtech/nawthtech/backend/internal/services"

	"github.com/go-chi/chi/v5"
)

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
}