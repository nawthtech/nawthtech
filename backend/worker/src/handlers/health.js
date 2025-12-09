// worker/src/handlers/health.go
package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"nawthtech-worker/src/utils"
)

// HealthResponse تمثل هيكل استجابة الصحة
type HealthResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// DatabaseHealth تمثل حالة قاعدة البيانات
type DatabaseHealth struct {
	Status   string `json:"status"`
	Database string `json:"database"`
}

// CheckHealth فحص صحة الخدمة
func CheckHealth(w http.ResponseWriter, r *http.Request, env map[string]string) {
	dbManager := utils.GetDatabaseManager(env)
	dbHealth := dbManager.HealthCheck()

	healthData := map[string]interface{}{
		"status":      dbHealth.Status,
		"database":    dbHealth.Database,
		"timestamp":   time.Now().UTC().Format(time.RFC3339),
		"environment": env["ENVIRONMENT"],
		"version":     env["API_VERSION"],
		"service":     "nawthtech-worker",
	}

	resp := HealthResponse{
		Success: true,
		Message: "Service is " + dbHealth.Status,
		Data:    healthData,
	}

	writeJSON(w, http.StatusOK, resp)
}

// ReadyHealth فحص جاهزية الخدمة
func ReadyHealth(w http.ResponseWriter, r *http.Request, env map[string]string) {
	dbManager := utils.GetDatabaseManager(env)
	dbHealth := dbManager.HealthCheck()

	if dbHealth.Status != "healthy" {
		resp := HealthResponse{
			Success: false,
			Error:   "SERVICE_NOT_READY",
			Message: "Database is not ready",
		}
		writeJSON(w, http.StatusServiceUnavailable, resp)
		return
	}

	data := map[string]interface{}{
		"status":    "ready",
		"database":  dbHealth.Database,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	}

	resp := HealthResponse{
		Success: true,
		Message: "Service is ready",
		Data:    data,
	}

	writeJSON(w, http.StatusOK, resp)
}

// writeJSON مساعدة لكتابة JSON للـ ResponseWriter
func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}