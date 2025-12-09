package handlers

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nawthtech/nawthtech/backend/internal/config"
	"github.com/nawthtech/nawthtech/backend/internal/utils"
)

// HandlerContainer يحتفظ بالاعتماديات
type HandlerContainer struct {
	Cfg    *config.Config
	DB     *sql.DB
}

// NewHandlerContainer
func NewHandlerContainer(cfg *config.Config, db *sql.DB) *HandlerContainer {
	return &HandlerContainer{
		Cfg: cfg,
		DB:  db,
	}
}

// Health
func (h *HandlerContainer) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":      "healthy",
		"timestamp":   time.Now().UTC(),
		"version":     h.Cfg.Version,
		"environment": h.Cfg.Environment,
		"database":    "SQL (backend) / D1 (workers)",
	})
}

// Auth register (very simple stub)
func (h *HandlerContainer) Register(c *gin.Context) {
	type Req struct {
		Email     string `json:"email" binding:"required,email"`
		Password  string `json:"password" binding:"required,min=6"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	}
	var req Req
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	// create user in DB (example using SQL)
	now := time.Now()
	userID := "user_" + now.Format("20060102150405")

	// NOTE: you should insert into DB here. For now return success stub.
	access, refresh, _ := utils.GenerateJWT(h.Cfg, userID, req.Email, "user")

	c.JSON(http.StatusOK, gin.H{
		"success":       true,
		"access_token":  access,
		"refresh_token": refresh,
		"user": gin.H{
			"id":    userID,
			"email": req.Email,
			"name":  req.FirstName + " " + req.LastName,
		},
	})
}

// Login endpoint stub
func (h *HandlerContainer) Login(c *gin.Context) {
	type Req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}
	var req Req
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	// Authenticate against DB (omitted) -> return tokens
	userID := "user_12345"
	access, refresh, _ := utils.GenerateJWT(h.Cfg, userID, req.Email, "user")

	c.JSON(http.StatusOK, gin.H{
		"success":       true,
		"access_token":  access,
		"refresh_token": refresh,
		"user": gin.H{"id": userID, "email": req.Email},
	})
}

// Protected example: GetProfile
func (h *HandlerContainer) GetProfile(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "UNAUTHORIZED"})
		return
	}
	// Query user from DB (omitted). Return stub:
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"id":         userID,
			"email":      "user@example.com",
			"first_name": "Nawth",
			"last_name":  "Tech",
		},
	})
}

// GetServices example
func (h *HandlerContainer) GetServices(c *gin.Context) {
	// Query list of services from DB or D1 proxy
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"services": []gin.H{
				{"id": "srv1", "title": "خدمة تجريبية", "price": 99.0},
			},
		},
	})
}