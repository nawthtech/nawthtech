package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"nawthtech-worker/middleware"
	"nawthtech-worker/handlers"
	"nawthtech-worker/utils"
)

type RequestWithEnv struct {
	*http.Request
	Env map[string]string
	DB  *utils.D1Database
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8787"
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		env := loadEnv()
		db, err := utils.GetD1Database(env["D1_DATABASE_URL"])
		if err != nil {
			http.Error(w, "Database connection failed", 500)
			return
		}

		req := &RequestWithEnv{
			Request: r,
			Env:     env,
			DB:      db,
		}

		handleRoutes(w, req)
	})

	log.Printf("🚀 Worker running on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func loadEnv() map[string]string {
	env := map[string]string{
		"ENVIRONMENT":       os.Getenv("ENVIRONMENT"),
		"API_VERSION":       os.Getenv("API_VERSION"),
		"D1_DATABASE_URL":   os.Getenv("D1_DATABASE_URL"),
		"JWT_SECRET":        os.Getenv("JWT_SECRET"),
		"SESSION_SECRET":    os.Getenv("SESSION_SECRET"),
		"CORS_ALLOWED_ORIGINS": os.Getenv("CORS_ALLOWED_ORIGINS"),
	}
	return env
}

func handleRoutes(w http.ResponseWriter, r *RequestWithEnv) {
	path := r.URL.Path
	method := r.Method

	// Middleware CORS
	if !middleware.CORS(w, r) {
		return
	}

	// ✅ Health routes
	if path == "/health" && method == http.MethodGet {
		handlers.HealthCheckHandler(w, r)
		return
	}
	if path == "/health/ready" && method == http.MethodGet {
		handlers.HealthReadyHandler(w, r)
		return
	}

	// ✅ Auth routes
	if strings.HasPrefix(path, "/auth") {
		switch path {
		case "/auth/register":
			if method == http.MethodPost {
				handlers.RegisterHandler(w, r)
				return
			}
		case "/auth/login":
			if method == http.MethodPost {
				handlers.LoginHandler(w, r)
				return
			}
		case "/auth/refresh":
			if method == http.MethodPost {
				handlers.RefreshHandler(w, r)
				return
			}
		case "/auth/forgot-password":
			if method == http.MethodPost {
				handlers.ForgotPasswordHandler(w, r)
				return
			}
		}
	}

	// ✅ Protected routes
	if strings.HasPrefix(path, "/user") {
		if !middleware.Auth(w, r, r.Env["JWT_SECRET"]) {
			return
		}
		switch path {
		case "/user/profile":
			if method == http.MethodGet {
				handlers.GetProfileHandler(w, r)
				return
			}
			if method == http.MethodPut {
				handlers.UpdateProfileHandler(w, r)
				return
			}
		}
	}

	// ✅ Services routes
	if strings.HasPrefix(path, "/services") {
		switch path {
		case "/services":
			if method == http.MethodGet {
				handlers.GetServicesHandler(w, r)
				return
			}
		default:
			if method == http.MethodGet && strings.HasPrefix(path, "/services/") {
				handlers.GetServiceByIDHandler(w, r)
				return
			}
		}
	}

	// ❌ Route not found
	http.NotFound(w, r)
}

// Helper: write JSON response
func JSONResponse(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if payload != nil {
		if err := json.NewEncoder(w).Encode(payload); err != nil {
			log.Println("JSON encode error:", err)
		}
	}
}

// Optional: Timeout context
func withTimeout(ctx context.Context, duration time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, duration)
}