package main

import (
	"cmp"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
)

func main() {
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
}