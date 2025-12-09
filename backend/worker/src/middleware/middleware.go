package middleware

import (
	"net/http"
	"strings"
	"github.com/nawthtech/nawthtech/backend/worker/src/config"
	"github.com/nawthtech/nawthtech/backend/worker/src/utils"
)

// Cors middleware
func Cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg := config.LoadConfig()
		origin := r.Header.Get("Origin")
		if contains(cfg.CORSOrigins, origin) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		}
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// Auth middleware
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.JSONResponse(w, http.StatusUnauthorized, map[string]interface{}{
				"success": false,
				"error":   "UNAUTHORIZED",
			})
			return
		}
		token := strings.TrimPrefix(authHeader, "Bearer ")
		userID, err := utils.ValidateJWT(token)
		if err != nil {
			utils.JSONResponse(w, http.StatusUnauthorized, map[string]interface{}{
				"success": false,
				"error":   "INVALID_TOKEN",
			})
			return
		}
		r.Header.Set("X-User-ID", userID)
		next.ServeHTTP(w, r)
	})
}

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}