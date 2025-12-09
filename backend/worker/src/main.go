package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"nawthtech-worker/src/handlers"
	"nawthtech-worker/src/utils"
)

func main() {
	// تهيئة قاعدة البيانات D1
	if err := utils.InitDatabase(); err != nil {
		log.Fatalf("❌ Failed to initialize database: %v", err)
	}
	defer utils.CloseDatabase()

	// إنشاء Router
	r := chi.NewRouter()

	// Middleware عامة
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(handlers.CorsMiddleware(next.ServeHTTP))
	})

	// ==========================
	// Health Routes
	// ==========================
	r.Get("/health", handlers.HealthCheckHandler)
	r.Get("/health/ready", handlers.HealthReadyHandler)

	// ==========================
	// Users Routes
	// ==========================
	r.Get("/user/profile", handlers.GetUserProfileHandler)
	r.Get("/users", handlers.GetUsersHandler)

	// ==========================
	// Services Routes
	// ==========================
	r.Get("/services", handlers.GetServicesHandler)
	r.Get("/services/{id}", func(w http.ResponseWriter, r *http.Request) {
		serviceID := chi.URLParam(r, "id")
		handlers.GetServiceByIDHandler(w, r, serviceID)
	})

	// ==========================
	// Test Route
	// ==========================
	r.Get("/test", handlers.TestHandler)

	// ==========================
	// Not Found Handler
	// ==========================
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		handlers.RespondJSON(w, http.StatusNotFound, handlers.JSONResponse{
			Success: false,
			Error:   "NOT_FOUND",
			Message: "Route does not exist",
		})
	})

	// بدء السيرفر
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("🚀 Server running on port %s in %s mode", port, os.Getenv("ENVIRONMENT"))
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("❌ Server failed: %v", err)
	}
}