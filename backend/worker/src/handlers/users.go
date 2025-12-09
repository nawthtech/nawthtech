package handlers

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3" // D1 يستخدم SQLite تحت الغطاء
)

// GetUsers يسترجع قائمة المستخدمين من D1
func GetUsers(w http.ResponseWriter, r *http.Request) {
	db, err := connectD1()
	if err != nil {
		http.Error(w, "Database connection failed", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	type User struct {
		ID    string `db:"id" json:"id"`
		Name  string `db:"name" json:"name"`
		Email string `db:"email" json:"email"`
		// لا تعرض كلمة السر
	}

	var users []User
	err = db.Select(&users, "SELECT id, name, email FROM users LIMIT 50")
	if err != nil {
		http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
		return
	}

	resp := map[string]interface{}{
		"success": true,
		"data":    users,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// GetProfile يسترجع بيانات المستخدم الواحد من D1
func GetProfile(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	db, err := connectD1()
	if err != nil {
		http.Error(w, "Database connection failed", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	type User struct {
		ID    string `db:"id" json:"id"`
		Name  string `db:"name" json:"name"`
		Email string `db:"email" json:"email"`
	}

	var user User
	err = db.Get(&user, "SELECT id, name, email FROM users WHERE id = ? LIMIT 1", userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	resp := map[string]interface{}{
		"success": true,
		"data":    user,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// connectD1 ينشئ اتصال بـ D1 باستخدام SQLite (Cloudflare D1)
func connectD1() (*sqlx.DB, error) {
	d1Path := os.Getenv("D1_DB_NAME")
	if d1Path == "" {
		return nil, errMissingD1
	}
	db, err := sqlx.Connect("sqlite3", d1Path)
	if err != nil {
		return nil, err
	}
	return db, nil
}

var errMissingD1 = &customError{"D1_DB_NAME environment variable is not set"}

type customError struct {
	msg string
}

func (e *customError) Error() string {
	return e.msg
}