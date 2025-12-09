package handlers

import (
	"encoding/json"
	"net/http"

	"nawthtech/utils"
)

// GetProfile يعيد بيانات المستخدم بعد التحقق من JWT
func GetProfile(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		respondJSON(w, http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"error":   "UNAUTHORIZED",
		})
		return
	}

	db := utils.GetD1DB()
	user, err := utils.GetUserByID(db, userID)
	if err != nil {
		respondJSON(w, http.StatusNotFound, map[string]interface{}{
			"success": false,
			"error":   "USER_NOT_FOUND",
		})
		return
	}

	// إزالة الحقول الحساسة
	delete(user, "password")

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    user,
	})
}

// ================= Helper
func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}