package middleware

import (
	"net/http"
	"os"
	"strings"
	"nawthtech-worker/utils"
)

func CORS(w http.ResponseWriter, r *http.Request) bool {
	allowedOrigins := strings.Split(os.Getenv("CORS_ALLOWED_ORIGINS"), ",")
	origin := r.Header.Get("Origin")
	for _, o := range allowedOrigins {
		if strings.TrimSpace(o) == origin {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			break
		}
	}
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusNoContent)
		return false
	}
	return true
}

func Auth(w http.ResponseWriter, r *http.Request, jwtSecret string) bool {
	token := r.Header.Get("Authorization")
	if token == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return false
	}
	valid, claims := utils.ValidateJWT(token, jwtSecret)
	if !valid {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return false
	}
	r.Header.Set("user_id", claims["id"].(string))
	return true
}