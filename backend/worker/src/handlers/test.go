package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"worker/src/utils"
)

// TestHandler يقوم بالتحقق من الاتصال بقاعدة البيانات ويعيد رد تجريبي
func TestHandler(w http.ResponseWriter, r *http.Request, env map[string]string) {
	db := utils.GetDatabase()
	dbStatus := "healthy"

	// محاولة تنفيذ استعلام تجريبي على D1
	if err := db.Ping(); err != nil {
		dbStatus = "unhealthy"
	}

	response := map[string]interface{}{
		"success":     true,
		"message":     "Test endpoint working",
		"database":    "D1 Cloudflare",
		"db_status":   dbStatus,
		"timestamp":   time.Now().UTC().Format(time.RFC3339),
		"environment": env["ENVIRONMENT"],
		"version":     env["API_VERSION"],
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}