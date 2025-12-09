package main

import (
	"log"
	"net/http"
	"os"
	"worker/src/middleware"
	"worker/src/handlers"
)

func main() {
	mux := http.NewServeMux()

	// Middleware + Routes
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if !middleware.CORS(w, r) {
			return
		}
		handlers.HealthCheckHandler(w, r)
	})
	mux.HandleFunc("/health/ready", func(w http.ResponseWriter, r *http.Request) {
		if !middleware.CORS(w, r) {
			return
		}
		handlers.HealthReadyHandler(w, r)
	})

	// Auth Routes
	mux.HandleFunc("/auth/register", func(w http.ResponseWriter, r *http.Request) {
		if !middleware.CORS(w, r) {
			return
		}
		handlers.RegisterHandler(w, r)
	})
	mux.HandleFunc("/auth/login", func(w http.ResponseWriter, r *http.Request) {
		if !middleware.CORS(w, r) {
			return
		}
		handlers.LoginHandler(w, r)
	})
	mux.HandleFunc("/auth/refresh", func(w http.ResponseWriter, r *http.Request) {
		if !middleware.CORS(w, r) {
			return
		}
		handlers.RefreshHandler(w, r)
	})

	// Protected Routes
	mux.HandleFunc("/user/profile", func(w http.ResponseWriter, r *http.Request) {
		if !middleware.CORS(w, r) {
			return
		}
		jwtSecret := os.Getenv("JWT_SECRET")
		if !middleware.Auth(w, r, jwtSecret) {
			return
		}
		handlers.UserProfileHandler(w, r)
	})
	mux.HandleFunc("/user/list", func(w http.ResponseWriter, r *http.Request) {
		if !middleware.CORS(w, r) {
			return
		}
		jwtSecret := os.Getenv("JWT_SECRET")
		if !middleware.Auth(w, r, jwtSecret) {
			return
		}
		handlers.UserListHandler(w, r)
	})

	// Services
	mux.HandleFunc("/services", func(w http.ResponseWriter, r *http.Request) {
		if !middleware.CORS(w, r) {
			return
		}
		handlers.ServiceListHandler(w, r)
	})
	mux.HandleFunc("/services/", func(w http.ResponseWriter, r *http.Request) {
		if !middleware.CORS(w, r) {
			return
		}
		handlers.ServiceDetailHandler(w, r)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Println("Server running on port:", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}