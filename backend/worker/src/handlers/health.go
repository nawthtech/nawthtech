package handlers

import (
	"encoding/json"
	"net/http"
	"time"
)

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	resp := map[string]interface{}{
		"success": true,
		"message": "Service is healthy",
		"data": map[string]interface{}{
			"status":      "healthy",
			"database":    "D1",
			"timestamp":   time.Now().UTC().Format(time.RFC3339),
			"environment": os.Getenv("ENVIRONMENT"),
			"version":     os.Getenv("API_VERSION"),
		},
	}
	json.NewEncoder(w).Encode(resp)
}

func HealthReadyHandler(w http.ResponseWriter, r *http.Request) {
	resp := map[string]interface{}{
		"success": true,
		"message": "Service is ready",
		"data": map[string]interface{}{
			"status":    "ready",
			"database":  "D1",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		},
	}
	json.NewEncoder(w).Encode(resp)
}