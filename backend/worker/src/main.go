package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"worker/src/handlers"
	"worker/src/middleware"
	"worker/src/utils"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware as chiMiddleware"
)

// ==========================
// Main
// ==========================

func main() {
	// تهيئة البيئة و D1
	utils.LoadEnv()
	utils.InitDB() // يهيئ اتصال D1

	// إنشاء Router
	r := chi.NewRouter()

	// Middleware عامة
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)
	r.Use(middleware.CORS)

	// ==========================
	// Health Endpoints
	// ==========================
	r.Get("/health", handlers.HealthCheck)
	r.Get("/health/live", handlers.HealthCheck)
	r.Get("/health/ready", handlers.HealthReady)

	// ==========================
	// Auth Endpoints
	// ==========================
	r.Post("/auth/register", handlers.AuthRegister)
	r.Post("/auth/login", handlers.AuthLogin)
	r.Post("/auth/refresh", handlers.AuthRefresh)
	r.Post("/auth/forgot-password", handlers.AuthForgotPassword)

	// ==========================
	// User Endpoints (Protected)
	// ==========================
	r.Group(func(r chi.Router) {
		r.Use(middleware.Auth) // حماية جميع المسارات داخل هذه المجموعة
		r.Get("/user/profile", handlers.GetProfile)
	})

	// ==========================
	// Services Endpoints
	// ==========================
	r.Get("/services", handlers.GetServices)
	r.Get("/services/{id}", handlers.GetServiceByID)

	// ==========================
	// Test Endpoint
	// ==========================
	r.Get("/test", handlers.TestEndpoint)

	// ==========================
	// 404 Handler
	// ==========================
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not Found"))
	})

	// ==========================
	// Start server
	// ==========================
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Handler:      r,
		Addr:         ":" + port,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("🚀 Server started on port %s", port)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}