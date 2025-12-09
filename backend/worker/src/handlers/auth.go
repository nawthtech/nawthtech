package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"time"
	"github.com/nawthtech/nawthtech/backend/worker/src/utils"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	db, _ := utils.GetD1Database()

	var payload map[string]string
	json.NewDecoder(r.Body).Decode(&payload)
	username := payload["username"]
	email := payload["email"]
	password := payload["password"]

	// تخزين المستخدم في D1
	db.DB.Exec("INSERT INTO users (username,email,password) VALUES (?,?,?)", username, email, password)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "User registered",
	})
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	db, _ := utils.GetD1Database()
	var payload map[string]string
	json.NewDecoder(r.Body).Decode(&payload)
	email := payload["email"]
	password := payload["password"]

	row := db.DB.QueryRow("SELECT id FROM users WHERE email=? AND password=?", email, password)
	var userID string
	err := row.Scan(&userID)
	if err != nil {
		http.Error(w, "Invalid credentials", 401)
		return
	}

	jwtToken, _ := utils.GenerateJWT(map[string]interface{}{"id": userID}, os.Getenv("JWT_SECRET"), time.Hour*24)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"token":   jwtToken,
	})
}

func RefreshHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Refresh endpoint",
	})
}