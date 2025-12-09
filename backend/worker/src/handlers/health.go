package handlers

import (
	"net/http"
	"time"

	"nawthtech/worker/src/utils"
)

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	db, _ := utils.ConnectD1()
	status, _ := db.HealthCheck()
	utils.JSONResponse(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Service is " + status,
		"data": map[string]interface{}{
			"status":      status,
			"database":    "D1",
			"timestamp":   time.Now(),
			"service":     "nawthtech-worker",
			"environment": "production",
		},
	})
}

func HealthLive(w http.ResponseWriter, r *http.Request) {
	utils.JSONResponse(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"status":  "live",
	})
}

func HealthReady(w http.ResponseWriter, r *http.Request) {
	db, _ := utils.ConnectD1()
	status, _ := db.HealthCheck()
	if status != "healthy" {
		utils.JSONResponse(w, http.StatusServiceUnavailable, map[string]interface{}{
			"success": false,
			"error":   "SERVICE_NOT_READY",
		})
		return
	}

	utils.JSONResponse(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"status":  "ready",
	})
}