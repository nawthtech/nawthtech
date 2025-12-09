package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"worker/src/utils"
)

// ==== Responses مساعدة ====

type ResponseData struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Error   string      `json:"error,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func JSONResponse(w http.ResponseWriter, status int, resp ResponseData) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(resp)
}

// ==== Users Handlers ====

func GetProfile(w http.ResponseWriter, r *http.Request, userID string) {
	db := utils.GetD1().DB
	ctx := r.Context()

	if userID == "" {
		JSONResponse(w, http.StatusUnauthorized, ResponseData{
			Success: false,
			Error:   "UNAUTHORIZED",
		})
		return
	}

	sql := "SELECT id, name, email FROM users WHERE id = ?"
	rows, err := db.Query(ctx, sql, userID)
	if err != nil || len(rows) == 0 {
		JSONResponse(w, http.StatusNotFound, ResponseData{
			Success: false,
			Error:   "USER_NOT_FOUND",
		})
		return
	}

	user := map[string]interface{}{
		"id":    rows[0][0],
		"name":  rows[0][1],
		"email": rows[0][2],
	}

	JSONResponse(w, http.StatusOK, ResponseData{
		Success: true,
		Data:    user,
	})
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	db := utils.GetD1().DB
	ctx := r.Context()

	sql := "SELECT id, name, email FROM users LIMIT 50"
	rows, err := db.Query(ctx, sql)
	if err != nil {
		JSONResponse(w, http.StatusInternalServerError, ResponseData{
			Success: false,
			Error:   "DATABASE_ERROR",
		})
		return
	}

	users := []map[string]interface{}{}
	for _, row := range rows {
		users = append(users, map[string]interface{}{
			"id":    row[0],
			"name":  row[1],
			"email": row[2],
		})
	}

	JSONResponse(w, http.StatusOK, ResponseData{
		Success: true,
		Data:    users,
	})
}

// ==== Health Handlers ====

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	dbHealth := "unknown"
	db := utils.GetD1().DB

	if db != nil {
		dbHealth = "healthy"
	}

	JSONResponse(w, http.StatusOK, ResponseData{
		Success: true,
		Message: fmt.Sprintf("Service is %s", dbHealth),
		Data: map[string]interface{}{
			"status":      dbHealth,
			"database":    "D1",
			"timestamp":   time.Now().UTC().Format(time.RFC3339),
			"environment": r.Header.Get("ENVIRONMENT"),
			"version":     r.Header.Get("API_VERSION"),
			"service":     "nawthtech-worker",
		},
	})
}

func HealthReady(w http.ResponseWriter, r *http.Request) {
	db := utils.GetD1().DB
	if db == nil {
		JSONResponse(w, http.StatusServiceUnavailable, ResponseData{
			Success: false,
			Error:   "SERVICE_NOT_READY",
			Message: "Database is not ready",
		})
		return
	}

	JSONResponse(w, http.StatusOK, ResponseData{
		Success: true,
		Message: "Service is ready",
		Data: map[string]interface{}{
			"status":    "ready",
			"database":  "D1",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		},
	})
}

// ==== Services Handlers ====

func GetServices(w http.ResponseWriter, r *http.Request) {
	db := utils.GetD1().DB
	ctx := r.Context()

	sql := "SELECT id, title, description, price FROM services LIMIT 50"
	rows, err := db.Query(ctx, sql)
	if err != nil {
		JSONResponse(w, http.StatusInternalServerError, ResponseData{
			Success: false,
			Error:   "DATABASE_ERROR",
		})
		return
	}

	services := []map[string]interface{}{}
	for _, row := range rows {
		services = append(services, map[string]interface{}{
			"id":          row[0],
			"title":       row[1],
			"description": row[2],
			"price":       row[3],
		})
	}

	JSONResponse(w, http.StatusOK, ResponseData{
		Success: true,
		Data:    services,
	})
}

func GetServiceByID(w http.ResponseWriter, r *http.Request, serviceID string) {
	db := utils.GetD1().DB
	ctx := r.Context()

	sql := "SELECT id, title, description, price FROM services WHERE id = ?"
	rows, err := db.Query(ctx, sql, serviceID)
	if err != nil || len(rows) == 0 {
		JSONResponse(w, http.StatusNotFound, ResponseData{
			Success: false,
			Error:   "SERVICE_NOT_FOUND",
		})
		return
	}

	service := map[string]interface{}{
		"id":          rows[0][0],
		"title":       rows[0][1],
		"description": rows[0][2],
		"price":       rows[0][3],
	}

	JSONResponse(w, http.StatusOK, ResponseData{
		Success: true,
		Data:    service,
	})
}

// ==== Test Handler ====

func TestHandler(w http.ResponseWriter, r *http.Request) {
	JSONResponse(w, http.StatusOK, ResponseData{
		Success: true,
		Message: "Test handler is working!",
	})
}