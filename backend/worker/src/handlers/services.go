package handlers

import (
	"encoding/json"
	"net/http"

	"nawthtech/utils"
)

func GetServices(w http.ResponseWriter, r *http.Request) {
	db := utils.GetD1DB()
	services, _ := utils.GetAllServices(db)
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    services,
	})
}

func GetServiceByID(w http.ResponseWriter, r *http.Request) {
	db := utils.GetD1DB()
	id := r.URL.Path[len("/services/"):]
	service, err := utils.GetServiceByID(db, id)
	if err != nil {
		respondJSON(w, http.StatusNotFound, map[string]interface{}{
			"success": false,
			"error":   "SERVICE_NOT_FOUND",
		})
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    service,
	})
}

// ================= Helper
func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}