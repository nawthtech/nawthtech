package handlers

import (
	"encoding/json"
	"net/http"
	"time"
)

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	resp := map[string]interface{}{
		"success":    true,
		"message":    "Service is healthy",
		"data": map[string]interface{}{
			"status":      "healthy",
			"database":    "D1",
			"timestamp":   time.Now().Format(time.RFC3339),
			"environment": "production",
			"version":     "v1",
			"service":     "nawthtech-worker",
		},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}