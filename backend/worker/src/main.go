package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
	
	"worker/src/handlers"
	"worker/src/utils"
)

// ===== Middleware =====

func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		allowedOrigins := strings.Split(os.Getenv("CORS_ALLOWED_ORIGINS"), ",")
		origin := r.Header.Get("Origin")

		for _, o := range allowedOrigins {
			if strings.TrimSpace(o) == origin {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				break
			}
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With, X-API-Key")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// مثال: افحص وجود header Authorization
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			handlers.JSONResponse(w, http.StatusUnauthorized, handlers.ResponseData{
				Success: false,
				Error:   "UNAUTHORIZED",
			})
			return
		}
		// هنا يمكنك إضافة تحقق JWT إذا أردت
		next.ServeHTTP(w, r)
	})
}

// ===== Router =====

func main() {
	// تهيئة قاعدة D1
	if err := utils.GetD1().Connect(); err != nil {
		log.Fatalf("❌ Failed to connect to D1: %v", err)
	}
	defer utils.GetD1().Disconnect(context.Background())

	mux := http.NewServeMux()

	// ===== Health =====
	mux.HandleFunc("/health", handlers.HealthCheck)
	mux.HandleFunc("/health/ready", handlers.HealthReady)

	// ===== Users =====
	mux.Handle("/user/profile", AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// مثال: اجعل userID يأتي من Authorization
		userID := r.Header.Get("X-User-ID")
		handlers.GetProfile(w, r, userID)
	})))

	mux.Handle("/users", AuthMiddleware(http.HandlerFunc(handlers.GetUsers)))

	// ===== Services =====
	mux.HandleFunc("/services", handlers.GetServices)
	mux.HandleFunc("/services/", func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) < 3 || parts[2] == "" {
			handlers.JSONResponse(w, http.StatusBadRequest, handlers.ResponseData{
				Success: false,
				Error:   "INVALID_SERVICE_ID",
			})
			return
		}
		serviceID := parts[2]
		handlers.GetServiceByID(w, r, serviceID)
	})

	// ===== Test =====
	mux.HandleFunc("/test", handlers.TestHandler)

	// ===== 404 =====
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handlers.JSONResponse(w, http.StatusNotFound, handlers.ResponseData{
			Success: false,
			Error:   "NOT_FOUND",
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	fmt.Printf("🚀 Worker running on port %s\n", port)
	err := http.ListenAndServe(":"+port, CORSMiddleware(mux))
	if err != nil {
		log.Fatalf("❌ Server failed: %v", err)
	}
}