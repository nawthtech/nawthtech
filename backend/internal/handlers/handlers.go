package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nawthtech/nawthtech/backend/internal/config"
	"github.com/nawthtech/nawthtech/backend/internal/logger"
	"github.com/nawthtech/nawthtech/backend/internal/middleware"
)

// HandlerContainer يحتوي على جميع handlers
type HandlerContainer struct {
	cfg    *config.Config
	client *http.Client
}

// NewHandlerContainer ينشئ container جديد للـ handlers
func NewHandlerContainer(cfg *config.Config) *HandlerContainer {
	return &HandlerContainer{
		cfg: cfg,
		client: &http.Client{
			Timeout: 30 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 10,
				IdleConnTimeout:     90 * time.Second,
			},
		},
	}
}

// RegisterAllRoutes يسجل جميع routes
func (hc *HandlerContainer) RegisterAllRoutes(app *gin.Engine) {
	// Health check
	app.GET("/health", hc.HealthCheck)
	app.GET("/api/health", hc.HealthCheck)

	// API v1
	api := app.Group("/api/v1")
	{
		// Authentication
		auth := api.Group("/auth")
		{
			auth.POST("/register", hc.Register)
			auth.POST("/login", hc.Login)
			auth.GET("/me", middleware.AuthMiddleware(), hc.GetCurrentUser)
		}

		// Users
		users := api.Group("/users")
		users.Use(middleware.AuthMiddleware())
		{
			users.GET("/", middleware.AdminMiddleware(), hc.ListUsers)
			users.GET("/:id", hc.GetUser)
			users.PUT("/:id", hc.UpdateUser)
		}

		// Services
		services := api.Group("/services")
		services.Use(middleware.AuthMiddleware())
		{
			services.GET("/", hc.ListServices)
			services.POST("/", hc.CreateService)
			services.GET("/:id", hc.GetService)
			services.PUT("/:id", hc.UpdateService)
			services.DELETE("/:id", hc.DeleteService)
		}

		// AI
		ai := api.Group("/ai")
		ai.Use(middleware.AuthMiddleware())
		{
			ai.POST("/generate", hc.GenerateAI)
			ai.GET("/quota", hc.GetAIQuota)
			ai.GET("/requests", hc.GetAIRequests)
		}

		// Email (Admin only)
		email := api.Group("/email")
		email.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())
		{
			email.GET("/logs", hc.GetEmailLogs)
			email.GET("/config", hc.GetEmailConfig)
		}
	}

	// Webhooks
	webhooks := app.Group("/webhooks")
	{
		webhooks.POST("/email", hc.HandleEmailWebhook)
		webhooks.POST("/stripe", hc.HandleStripeWebhook)
	}

	// Worker proxy
	worker := app.Group("/worker")
	{
		worker.GET("/health", hc.WorkerHealthCheck)
		worker.Any("/*path", hc.WorkerProxy)
	}
}

// ==================== Utility Functions ====================

// callWorker يستدعي Worker API
func (hc *HandlerContainer) callWorker(c *gin.Context, method, path string, data interface{}) (map[string]interface{}, error) {
	workerURL := os.Getenv("WORKER_API_URL")
	if workerURL == "" {
		workerURL = "https://api.nawthtech.com"
	}

	apiKey := os.Getenv("WORKER_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("WORKER_API_KEY not configured")
	}

	url := fmt.Sprintf("%s%s", workerURL, path)

	var reqBody []byte
	if data != nil {
		var err error
		reqBody, err = json.Marshal(data)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %v", err)
		}
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Copy headers from original request
	req.Header.Set("Authorization", c.GetHeader("Authorization"))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Worker-API-Key", apiKey)
	req.Header.Set("X-Forwarded-For", c.ClientIP())
	req.Header.Set("X-Forwarded-Host", c.Request.Host)
	req.Header.Set("X-Forwarded-Proto", c.Request.Proto)

	// Copy query parameters
	q := req.URL.Query()
	for key, values := range c.Request.URL.Query() {
		for _, value := range values {
			q.Add(key, value)
		}
	}
	req.URL.RawQuery = q.Encode()

	resp, err := hc.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("worker API call failed: %v", err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	// Set response status code
	c.Status(resp.StatusCode)

	return result, nil
}

// ==================== Handler Functions ====================

// HealthCheck handler
func (hc *HandlerContainer) HealthCheck(c *gin.Context) {
	health := gin.H{
		"status":      "healthy",
		"timestamp":   time.Now().UTC().Format(time.RFC3339),
		"environment": hc.cfg.Environment,
		"version":     "1.0.0",
		"services": gin.H{
			"worker_api": os.Getenv("WORKER_API_KEY") != "",
			"cache":      true,
			"rate_limit": true,
		},
		"uptime": fmt.Sprintf("%v", time.Since(startTime)),
	}

	// Check worker health if configured
	if os.Getenv("WORKER_API_KEY") != "" {
		workerHealth, err := hc.callWorker(c, "GET", "/health", nil)
		if err != nil {
			health["worker_status"] = "unreachable"
			health["worker_error"] = err.Error()
		} else {
			health["worker_status"] = "healthy"
			health["worker_data"] = workerHealth
		}
	}

	c.JSON(http.StatusOK, health)
}

var startTime = time.Now()

// WorkerHealthCheck checks worker health directly
func (hc *HandlerContainer) WorkerHealthCheck(c *gin.Context) {
	result, err := hc.callWorker(c, "GET", "/health", nil)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error":   "Worker unavailable",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

// WorkerProxy proxies requests to worker
func (hc *HandlerContainer) WorkerProxy(c *gin.Context) {
	path := c.Param("path")
	if path == "" {
		path = "/"
	}

	// Get request body
	var data interface{}
	if c.Request.Body != nil && c.Request.ContentLength > 0 {
		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid JSON body",
				"details": err.Error(),
			})
			return
		}
	}

	result, err := hc.callWorker(c, c.Request.Method, path, data)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"error":   "Failed to proxy request",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

// ==================== Authentication Handlers ====================

// Register handler
func (hc *HandlerContainer) Register(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Username string `json:"username" binding:"required,min=3,max=30"`
		Password string `json:"password" binding:"required,min=8"`
		FullName string `json:"full_name,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"details": err.Error(),
		})
		return
	}

	result, err := hc.callWorker(c, "POST", "/api/v1/auth/register", req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Registration failed",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, result)
}

// Login handler
func (hc *HandlerContainer) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"details": err.Error(),
		})
		return
	}

	result, err := hc.callWorker(c, "POST", "/api/v1/auth/login", req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Login failed",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetCurrentUser handler
func (hc *HandlerContainer) GetCurrentUser(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	result, err := hc.callWorker(c, "GET", fmt.Sprintf("/api/v1/users/%s", userID), nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get user",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

// ==================== User Handlers ====================

// ListUsers handler (admin only)
func (hc *HandlerContainer) ListUsers(c *gin.Context) {
	// Build query string from request parameters
	queryParams := c.Request.URL.Query()
	path := "/api/v1/users"
	if len(queryParams) > 0 {
		path += "?"
		for key, values := range queryParams {
			for _, value := range values {
				path += fmt.Sprintf("%s=%s&", key, value)
			}
		}
		path = strings.TrimSuffix(path, "&")
	}

	result, err := hc.callWorker(c, "GET", path, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to list users",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetUser handler
func (hc *HandlerContainer) GetUser(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User ID is required",
		})
		return
	}

	// Check permissions
	currentUserID, _ := c.Get("userID")
	currentUserRole, _ := c.Get("userRole")

	if userID != currentUserID && currentUserRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Access denied",
		})
		return
	}

	result, err := hc.callWorker(c, "GET", fmt.Sprintf("/api/v1/users/%s", userID), nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get user",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

// UpdateUser handler
func (hc *HandlerContainer) UpdateUser(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User ID is required",
		})
		return
	}

	// Check permissions
	currentUserID, _ := c.Get("userID")
	currentUserRole, _ := c.Get("userRole")

	if userID != currentUserID && currentUserRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Access denied",
		})
		return
	}

	var updateData map[string]interface{}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"details": err.Error(),
		})
		return
	}

	result, err := hc.callWorker(c, "PUT", fmt.Sprintf("/api/v1/users/%s", userID), updateData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update user",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

// ==================== Service Handlers ====================

// ListServices handler
func (hc *HandlerContainer) ListServices(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	// Build query with user filter for non-admin users
	userRole, _ := c.Get("userRole")
	queryParams := c.Request.URL.Query()
	
	if userRole != "admin" {
		queryParams.Set("user_id", userID.(string))
	}

	path := "/api/v1/services"
	if len(queryParams) > 0 {
		path += "?"
		for key, values := range queryParams {
			for _, value := range values {
				path += fmt.Sprintf("%s=%s&", key, value)
			}
		}
		path = strings.TrimSuffix(path, "&")
	}

	result, err := hc.callWorker(c, "GET", path, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to list services",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

// CreateService handler
func (hc *HandlerContainer) CreateService(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	var serviceData map[string]interface{}
	if err := c.ShouldBindJSON(&serviceData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"details": err.Error(),
		})
		return
	}

	// Add user ID to service data
	serviceData["user_id"] = userID

	result, err := hc.callWorker(c, "POST", "/api/v1/services", serviceData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create service",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, result)
}

// GetService handler
func (hc *HandlerContainer) GetService(c *gin.Context) {
	serviceID := c.Param("id")
	if serviceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Service ID is required",
		})
		return
	}

	result, err := hc.callWorker(c, "GET", fmt.Sprintf("/api/v1/services/%s", serviceID), nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get service",
			"details": err.Error(),
		})
		return
	}

	// Check ownership for non-admin users
	userID, _ := c.Get("userID")
	userRole, _ := c.Get("userRole")
	
	if serviceData, ok := result["data"].(map[string]interface{}); ok {
		if serviceUserID, ok := serviceData["user_id"].(string); ok {
			if serviceUserID != userID && userRole != "admin" {
				c.JSON(http.StatusForbidden, gin.H{
					"error": "Access denied",
				})
				return
			}
		}
	}

	c.JSON(http.StatusOK, result)
}

// UpdateService handler
func (hc *HandlerContainer) UpdateService(c *gin.Context) {
	serviceID := c.Param("id")
	if serviceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Service ID is required",
		})
		return
	}

	var updateData map[string]interface{}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"details": err.Error(),
		})
		return
	}

	result, err := hc.callWorker(c, "PUT", fmt.Sprintf("/api/v1/services/%s", serviceID), updateData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update service",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

// DeleteService handler
func (hc *HandlerContainer) DeleteService(c *gin.Context) {
	serviceID := c.Param("id")
	if serviceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Service ID is required",
		})
		return
	}

	result, err := hc.callWorker(c, "DELETE", fmt.Sprintf("/api/v1/services/%s", serviceID), nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to delete service",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

// ==================== AI Handlers ====================

// GenerateAI handler
func (hc *HandlerContainer) GenerateAI(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	var aiRequest map[string]interface{}
	if err := c.ShouldBindJSON(&aiRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"details": err.Error(),
		})
		return
	}

	// Add user ID to request
	aiRequest["user_id"] = userID

	result, err := hc.callWorker(c, "POST", "/api/v1/ai/generate", aiRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to generate AI content",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetAIQuota handler
func (hc *HandlerContainer) GetAIQuota(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	result, err := hc.callWorker(c, "GET", fmt.Sprintf("/api/v1/ai/quota?user_id=%s", userID), nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get AI quota",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetAIRequests handler
func (hc *HandlerContainer) GetAIRequests(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	// Build query with user filter
	queryParams := c.Request.URL.Query()
	queryParams.Set("user_id", userID.(string))

	path := "/api/v1/ai/requests"
	if len(queryParams) > 0 {
		path += "?"
		for key, values := range queryParams {
			for _, value := range values {
				path += fmt.Sprintf("%s=%s&", key, value)
			}
		}
		path = strings.TrimSuffix(path, "&")
	}

	result, err := hc.callWorker(c, "GET", path, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get AI requests",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

// ==================== Email Handlers ====================

// GetEmailLogs handler
func (hc *HandlerContainer) GetEmailLogs(c *gin.Context) {
	// Only admin can access email logs
	userRole, exists := c.Get("userRole")
	if !exists || userRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Admin access required",
		})
		return
	}

	queryParams := c.Request.URL.Query()
	path := "/api/v1/email/logs"
	if len(queryParams) > 0 {
		path += "?"
		for key, values := range queryParams {
			for _, value := range values {
				path += fmt.Sprintf("%s=%s&", key, value)
			}
		}
		path = strings.TrimSuffix(path, "&")
	}

	result, err := hc.callWorker(c, "GET", path, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get email logs",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetEmailConfig handler
func (hc *HandlerContainer) GetEmailConfig(c *gin.Context) {
	// Only admin can access email config
	userRole, exists := c.Get("userRole")
	if !exists || userRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Admin access required",
		})
		return
	}

	result, err := hc.callWorker(c, "GET", "/api/v1/email/config", nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get email config",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

// ==================== Webhook Handlers ====================

// HandleEmailWebhook handler
func (hc *HandlerContainer) HandleEmailWebhook(c *gin.Context) {
	var webhookData map[string]interface{}
	if err := c.ShouldBindJSON(&webhookData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid webhook data",
			"details": err.Error(),
		})
		return
	}

	// Verify webhook signature if configured
	webhookSecret := os.Getenv("EMAIL_WEBHOOK_SECRET")
	if webhookSecret != "" {
		signature := c.GetHeader("X-Webhook-Signature")
		if signature == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Missing webhook signature",
			})
			return
		}
		// Add signature verification logic here
	}

	// Forward to worker
	result, err := hc.callWorker(c, "POST", "/webhooks/email", webhookData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to process webhook",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

// HandleStripeWebhook handler
func (hc *HandlerContainer) HandleStripeWebhook(c *gin.Context) {
	var webhookData map[string]interface{}
	if err := c.ShouldBindJSON(&webhookData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid webhook data",
			"details": err.Error(),
		})
		return
	}

	// Verify Stripe signature if configured
	stripeSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")
	if stripeSecret != "" {
		signature := c.GetHeader("Stripe-Signature")
		if signature == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Missing Stripe signature",
			})
			return
		}
		// Add Stripe signature verification logic here
	}

	// Forward to worker
	result, err := hc.callWorker(c, "POST", "/webhooks/stripe", webhookData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to process Stripe webhook",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

// ==================== Middleware Helper Functions ====================

// Helper function to extract user ID from JWT token
func extractUserIDFromToken(tokenString string) (string, string, error) {
	// Implement JWT token parsing logic
	// For now, return placeholder values
	return "user-123", "user", nil
}

// Helper function to validate API key
func validateAPIKey(apiKey string) (string, string, error) {
	// Implement API key validation logic
	// For now, return placeholder values
	return "user-123", "user", nil
}