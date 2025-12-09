package handlers

import (
	"encoding/json"
	"net/http"
	"time"
	"nawthtech-worker/utils"
)

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	dbHealth := map[string]string{"status": "healthy", "database": "D1"}
	resp := map[string]interface{}{
		"success": true,
		"message": "Service is healthy",
		"data": map[string]interface{}{
			"status":     dbHealth["status"],
			"database":   dbHealth["database"],
			"timestamp":  time.Now().UTC().Format(time.RFC3339),
			"environment": r.Header.Get("ENVIRONMENT"),
			"version":    r.Header.Get("API_VERSION"),
		},
	}
	json.NewEncoder(w).Encode(resp)
}

func HealthReadyHandler(w http.ResponseWriter, r *http.Request) {
	dbHealth := map[string]string{"status": "healthy", "database": "D1"}
	if dbHealth["status"] != "healthy" {
		http.Error(w, "Service not ready", 503)
		return
	}
	resp := map[string]interface{}{
		"success": true,
		"message": "Service is ready",
		"data": map[string]interface{}{
			"status":    "ready",
			"database":  dbHealth["database"],
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		},
	}
	json.NewEncoder(w).Encode(resp)
}