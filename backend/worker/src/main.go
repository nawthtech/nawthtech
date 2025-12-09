package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/cors"

	"nawthtech/worker/src/config"
	"nawthtech/worker/src/handlers"
	"nawthtech/worker/src/middleware"
)

func main() {
	cfg := config.LoadConfig()

	router := chi.NewRouter()

	// ✅ Middleware شامل
	c := cors.New(cors.Options{
		AllowedOrigins:   cfg.CORSAllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Content-Length"},
		AllowCredentials: true,
	})
	router.Use(c.Handler)
	router.Use(middleware.AuthMiddleware)

	// ✅ Health routes
	router.Get("/health", handlers.HealthCheck)
	router.Get("/health/live", handlers.HealthLive)
	router.Get("/health/ready", handlers.HealthReady)

	// ✅ Auth routes
	router.Post("/auth/register", handlers.Register)
	router.Post("/auth/login", handlers.Login)
	router.Post("/auth/refresh", handlers.Refresh)
	router.Post("/auth/forgot-password", handlers.ForgotPassword)

	// ✅ User routes
	router.Get("/user/profile", handlers.GetProfile)
	router.Put("/user/profile", handlers.UpdateProfile)

	// ✅ Services routes
	router.Get("/services", handlers.GetServices)
	router.Get("/services/{id}", handlers.GetServiceByID)

	// ✅ Test route
	router.Get("/test", handlers.TestEndpoint)

	port := cfg.Port
	log.Printf("🚀 Server running on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}