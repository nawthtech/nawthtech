package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nawthtech/nawthtech/backend/internal/ai"
	"github.com/nawthtech/nawthtech/backend/internal/models"
	"github.com/nawthtech/nawthtech/backend/internal/services"
	"github.com/nawthtech/nawthtech/backend/internal/utils"
)

// ================================
// تعريفات الـ Handler Structs
// ================================

// AuthHandler معالجة طلبات المصادقة
type AuthHandler struct {
	service services.AuthService
}

// UserHandler معالجة طلبات المستخدمين
type UserHandler struct {
	service services.UserService
}

// ServiceHandler معالجة طلبات الخدمات
type ServiceHandler struct {
	service services.ServiceService
}

// CategoryHandler معالجة طلبات الفئات
type CategoryHandler struct {
	service services.CategoryService
}

// OrderHandler معالجة طلبات الطلبات
type OrderHandler struct {
	service services.OrderService
}

// PaymentHandler معالجة طلبات الدفع
type PaymentHandler struct {
	service services.PaymentService
}

// UploadHandler معالجة طلبات الرفع
type UploadHandler struct {
	service services.UploadService
}

// NotificationHandler معالجة طلبات الإشعارات
type NotificationHandler struct {
	service services.NotificationService
}

// AdminHandler معالجة طلبات الإدارة
type AdminHandler struct {
	service services.AdminService
}

// HealthHandler معالجة طلبات الصحة
type HealthHandler struct {
	service services.HealthService
}

// AIHandler معالجة طلبات الذكاء الاصطناعي
type AIHandler struct {
	aiClient *ai.Client
}

// ================================
// HandlerContainer
// ================================

// HandlerContainer حاوية لجميع الـ handlers
type HandlerContainer struct {
	Auth         *AuthHandler
	User         *UserHandler
	Service      *ServiceHandler
	Category     *CategoryHandler
	Order        *OrderHandler
	Payment      *PaymentHandler
	Upload       *UploadHandler
	Notification *NotificationHandler
	Admin        *AdminHandler
	Health       *HealthHandler
	AI           *AIHandler
	Email        *EmailHandler
}

// NewHandlerContainer إنشاء حاوية handlers جديدة
func NewHandlerContainer(serviceContainer *services.ServiceContainer) *HandlerContainer {
	container := &HandlerContainer{}
 aiClient, err := ai.NewClient()
	if err == nil {
		container.AI = &AIHandler{aiClient: aiClient}
		log.Println("✅ AI Client initialized")
	} else {
		log.Printf("⚠️ Failed to initialize AI Client: %v", err)
	}
	if serviceContainer != nil {
	}
		if serviceContainer.Auth != nil {
			container.Auth = &AuthHandler{service: serviceContainer.Auth}
		}
		if serviceContainer.User != nil {
			container.User = &UserHandler{service: serviceContainer.User}
		}
		if serviceContainer.Service != nil {
			container.Service = &ServiceHandler{service: serviceContainer.Service}
		}
		if serviceContainer.Category != nil {
			container.Category = &CategoryHandler{service: serviceContainer.Category}
		}
		if serviceContainer.Order != nil {
			container.Order = &OrderHandler{service: serviceContainer.Order}
		}
		if serviceContainer.Payment != nil {
			container.Payment = &PaymentHandler{service: serviceContainer.Payment}
		}
		if serviceContainer.Upload != nil {
			container.Upload = &UploadHandler{service: serviceContainer.Upload}
		}
		if serviceContainer.Notification != nil {
			container.Notification = &NotificationHandler{service: serviceContainer.Notification}
		}
		if serviceContainer.Admin != nil {
			container.Admin = &AdminHandler{service: serviceContainer.Admin}
		}
		if serviceContainer.Health != nil {
			container.Health = &HealthHandler{service: serviceContainer.Health}
		}
		if serviceContainer.Email != nil {
			emailWorker, err := email.NewCloudflareEmailWorker()
		if err == nil {
				container.Email = &EmailHandler{
					service:     serviceContainer.Email,
					emailWorker: emailWorker,
				}
			} else {
				// Log error but continue without email worker
				fmt.Printf("Warning: Failed to initialize email worker: %v\n", err)
				container.Email = &EmailHandler{
					service: serviceContainer.Email,
				}
  if serviceContainer != nil {
		// ... rest of the handlers
	}
			}
		}
	}

	return container
}

// ================================
// NewAIHandler (دالة منفصلة لإضافة AI handler)
// ================================

// NewAIHandler إنشاء AI handler جديد
func NewAIHandler(aiClient *ai.Client) *AIHandler {
	return &AIHandler{aiClient: aiClient}
}

// ================================
// AuthHandler Methods
// ================================

// Register تسجيل مستخدم جديد
func (h *AuthHandler) Register(c *gin.Context) {
	var req services.AuthRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	response, err := h.service.Register(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, response)
}

// Login تسجيل الدخول
func (h *AuthHandler) Login(c *gin.Context) {
	var req services.AuthLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	response, err := h.service.Login(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// Logout تسجيل الخروج
func (h *AuthHandler) Logout(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token is required"})
		return
	}

	err := h.service.Logout(c.Request.Context(), token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

// RefreshToken تجديد التوكن
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	response, err := h.service.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// ================================
// UserHandler Methods
// ================================

// GetProfile الحصول على الملف الشخصي
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	user, err := h.service.GetProfile(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateProfile تحديث الملف الشخصي
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var updateData map[string]interface{}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	user, err := h.service.UpdateProfile(c.Request.Context(), userID, updateData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// ================================
// ServiceHandler Methods
// ================================

// CreateService إنشاء خدمة جديدة
func (h *ServiceHandler) CreateService(c *gin.Context) {
	var service models.Service
	if err := c.ShouldBindJSON(&service); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	createdService, err := h.service.CreateService(c.Request.Context(), service)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdService)
}

// GetServices الحصول على قائمة الخدمات
func (h *ServiceHandler) GetServices(c *gin.Context) {
	services, err := h.service.GetServices(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, services)
}

// ================================
// CategoryHandler Methods
// ================================

// CreateCategory إنشاء فئة جديدة
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var category models.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	createdCategory, err := h.service.CreateCategory(c.Request.Context(), category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdCategory)
}

// GetCategories الحصول على قائمة الفئات
func (h *CategoryHandler) GetCategories(c *gin.Context) {
	categories, err := h.service.GetCategories(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, categories)
}

// ================================
// OrderHandler Methods
// ================================

// CreateOrder إنشاء طلب جديد
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var order models.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	createdOrder, err := h.service.CreateOrder(c.Request.Context(), order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdOrder)
}

// GetUserOrders الحصول على طلبات المستخدم
func (h *OrderHandler) GetUserOrders(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	orders, err := h.service.GetUserOrders(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orders)
}

// ================================
// PaymentHandler Methods
// ================================

// CreatePaymentIntent إنشاء نية دفع
func (h *PaymentHandler) CreatePaymentIntent(c *gin.Context) {
	var req services.PaymentIntentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	intent, err := h.service.CreatePaymentIntent(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, intent)
}

// ConfirmPayment تأكيد الدفع
func (h *PaymentHandler) ConfirmPayment(c *gin.Context) {
	paymentID := c.Param("id")
	if paymentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Payment ID is required"})
		return
	}

	var confirmationData map[string]interface{}
	if err := c.ShouldBindJSON(&confirmationData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid confirmation data"})
		return
	}

	result, err := h.service.ConfirmPayment(c.Request.Context(), paymentID, confirmationData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// ================================
// UploadHandler Methods
// ================================

// UploadFile رفع ملف
func (h *UploadHandler) UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}

	uploadType := c.PostForm("type")
	if uploadType == "" {
		uploadType = "general"
	}

	result, err := h.service.UploadFile(c.Request.Context(), file, uploadType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// ================================
// NotificationHandler Methods
// ================================

// GetNotifications الحصول على الإشعارات
func (h *NotificationHandler) GetNotifications(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	notifications, err := h.service.GetUserNotifications(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, notifications)
}

// MarkAsRead تحديد الإشعار كمقروء
func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	notificationID := c.Param("id")
	if notificationID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Notification ID is required"})
		return
	}

	err := h.service.MarkAsRead(c.Request.Context(), notificationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notification marked as read"})
}

// ================================
// AdminHandler Methods
// ================================

// GetStatistics الحصول على إحصائيات النظام
func (h *AdminHandler) GetStatistics(c *gin.Context) {
	stats, err := h.service.GetSystemStatistics(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// GetAllUsers الحصول على جميع المستخدمين
func (h *AdminHandler) GetAllUsers(c *gin.Context) {
	users, err := h.service.GetAllUsers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

// ================================
// HealthHandler Methods
// ================================

// CheckHealth فحص صحة النظام
func (h *HealthHandler) CheckHealth(c *gin.Context) {
	healthStatus, err := h.service.CheckHealth(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":    "unhealthy",
			"error":     err.Error(),
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		})
		return
	}

	c.JSON(http.StatusOK, healthStatus)
}

// HealthCheck فحص صحة مبسط
func (h *HealthHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"service":   "nawthtech-backend",
		"version":   "1.0.0",
	})
}

// ================================
// AIHandler Methods (من ai.go)
// ================================

// GenerateContentHandler معالج توليد المحتوى
func (h *AIHandler) GenerateContentHandler(c *gin.Context) {
	var req struct {
		Prompt      string  `json:"prompt" binding:"required"`
		Language    string  `json:"language" default:"en"`
		ContentType string  `json:"content_type"`
		Tone        string  `json:"tone"`
		Length      string  `json:"length"`
		MaxTokens   int     `json:"max_tokens" default:"1000"`
		Temperature float64 `json:"temperature" default:"0.7"`
  Duration int    `json:"duration" default:"30"`
		Provider string `json:"provider" default:"auto"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	if len(req.Prompt) > 1000 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Prompt is too long (max 1000 characters)",
		})
		return
	}

	// التحقق من طول النص
	if len(req.Prompt) > 5000 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Prompt is too long (max 5000 characters)",
		})
		return
	}

	// بناء prompt محسن
	prompt := h.buildEnhancedPrompt(req.Prompt, req.ContentType, req.Tone, req.Length)

	// توليد المحتوى باستخدام الواجهة الصحيحة
	content, err := h.aiClient.GenerateText(prompt, req.Provider)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to generate content",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"content":      content,
			"language":     req.Language,
			"content_type": req.ContentType,
			"tone":         req.Tone,
			"provider":     "req.provider",
			"model_used":   "default",
			"cost":         0.0,
			"tokens_used":  len(content),
			"created_at":   time.Now().UTC().Format(time.RFC3339),
		},
	})
}

// AnalyzeImageHandler معالج تحليل الصور
func (h *AIHandler) AnalyzeImageHandler(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Image file is required",
		})
		return
	}

	// التحقق من نوع الملف
	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/jpg":  true,
		"image/png":  true,
		"image/gif":  true,
		"image/webp": true,
	}

	fileHeader, _ := file.Open()
	buffer := make([]byte, 512)
	_, err = fileHeader.Read(buffer)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Failed to read image file",
		})
		return
	}
	fileHeader.Close()

	contentType := http.DetectContentType(buffer)
	if !allowedTypes[contentType] {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid image format. Supported formats: JPEG, PNG, GIF, WebP",
		})
		return
	}

	// التحقق من حجم الملف (10MB كحد أقصى)
	if file.Size > 10*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Image size too large (max 10MB)",
		})
		return
	}

	prompt := c.PostForm("prompt")
	if prompt == "" {
		prompt = "Describe this image in detail"
	}

	provider := c.PostForm("provider")
	if provider == "" {
		provider = "auto"
	}

	// قراءة ملف الصورة
	uploadedFile, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to open image file",
			"details": err.Error(),
		})
		return
	}
	defer uploadedFile.Close()

	imageData, err := io.ReadAll(uploadedFile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to read image data",
			"details": err.Error(),
		})
		return
	}


	‎// تحليل الصورة باستخدام العميل الحالي
	analysis, err := h.aiClient.AnalyzeImage(imageData, prompt, provider)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to analyze image",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"analysis":      analysis.Text,
			"confidence":    0.85,
			"filename":      file.Filename,
			"size":          file.Size,
			"provider":      provider,
			"model_used":    "vision",
			"cost":          0.0,
			"created_at":    time.Now().UTC().Format(time.RFC3339),
		},
	})
}

‎// توليد الفيديو
	videoURL, err := h.aiClient.GenerateVideo(req.Prompt, req.Provider)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to generate video",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"success": true,
		"message": "Video generation started",
		"data": gin.H{
			"prompt":   req.Prompt,
			"duration": req.Duration,
			"provider": req.Provider,
			"status":   "processing",
			"note":     "Use the check_video endpoint to check status",
		},
	})
}

// CheckVideoStatusHandler معالب التحقق من حالة الفيديو
func (h *AIHandler) CheckVideoStatusHandler(c *gin.Context) {
	operationID := c.Param("id")
	if operationID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Operation ID is required",
		})
		return
	}

‎	// الحصول على حالة الفيديو
	status, err := h.aiClient.GetVideoStatus(operationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to get video status",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    status,
	})
}

// AnalyzeTextHandler معالج تحليل النص
func (h *AIHandler) AnalyzeTextHandler(c *gin.Context) {
	var req struct {
		Text     string `json:"text" binding:"required"`
		Provider string `json:"provider" default:"auto"`
		Type     string `json:"type" default:"general"` // sentiment, entities, keywords, etc.
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	if len(req.Text) > 5000 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Text is too long (max 5000 characters)",
		})
		return
	}

‎	// تحليل النص
	analysis, err := h.aiClient.AnalyzeText(req.Text, req.Provider)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to analyze text",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"text":        req.Text,
			"analysis":    analysis,
			"provider":    req.Provider,
			"analysis_type": req.Type,
			"created_at":  time.Now().UTC().Format(time.RFC3339),
		},
	})
}

// TranslateTextHandler معالج ترجمة النص
func (h *AIHandler) TranslateTextHandler(c *gin.Context) {
	var req struct {
		Text       string `json:"text" binding:"required"`
		SourceLang string `json:"source_lang" default:"auto"`
		TargetLang string `json:"target_lang" binding:"required"`
		Formal     bool   `json:"formal" default:"false"`
  Provider   string `json:"provider" default:"auto"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// التحقق من طول النص
	if len(req.Text) > 5000 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Text is too long (max 5000 characters)",
		})
		return
	}

	// التحقق من صحة اللغة الهدف
	if !h.isValidLanguage(req.TargetLang) {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid target language",
		})
		return
	}

	‎// ترجمة النص باستخدام العميل الحالي
	translatedText, err := h.aiClient.TranslateText(req.Text, req.SourceLang, req.TargetLang, req.Provider)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to translate text",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"original_text":   req.Text,
			"translated_text": translatedText,
			"source_lang":     req.SourceLang,
			"target_lang":     req.TargetLang,
			"detected_lang":   "auto",
			"confidence":      0.9,
			"provider":        "req.Provider",
			"model_used":      "translation",
			"cost":            0.0,
			"created_at":      time.Now().UTC().Format(time.RFC3339),
		},
	})
}

// SummarizeTextHandler معالج تلخيص النص
func (h *AIHandler) SummarizeTextHandler(c *gin.Context) {
	var req struct {
		Text        string `json:"text" binding:"required"`
		Language    string `json:"language" default:"en"`
		SummaryType string `json:"summary_type" default:"paragraph"`
		MaxLength   int    `json:"max_length" default:"200"`
		Provider    string `json:"provider" default:"auto"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

‎	// التحقق من طول النص
	if len(req.Text) > 10000 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Text is too long (max 10000 characters)",
		})
		return
	}

‎	// بناء prompt للتلخيص
	prompt := h.buildSummaryPrompt(req.Text, req.SummaryType, req.MaxLength)

‎	// توليد التلخيص باستخدام العميل الحالي
	summary, err := h.aiClient.GenerateText(prompt, req.Provider)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to generate summary",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"original_length": len(req.Text),
			"summary":         summary,
			"summary_type":    req.SummaryType,
			"language":        req.Language,
			"provider":        req.Provider,
			"model_used":      "summarization",
			"cost":            0.0,
			"created_at":      time.Now().UTC().Format(time.RFC3339),
		},
	})
}

// GetAICapabilitiesHandler معالج قدرات الذكاء الاصطناعي
func (h *AIHandler) GetAICapabilitiesHandler(c *gin.Context) {
	var providers map[string][]string
	var usageStats map[string]interface{}
	
	if h.aiClient != nil {
		providers = h.aiClient.GetAvailableProviders()
		usageStats = h.aiClient.GetUsageStatistics()
	}

	capabilities := gin.H{
		"features": []gin.H{
			{
				"name":                "text_generation",
				"description":         "Generate text content",
				"supported_languages": []string{"en", "ar", "fr", "es", "de", "zh"},
				"max_tokens":          4000,
			},
			{
				"name":              "image_analysis",
				"description":       "Analyze and describe images",
				"supported_formats": []string{"jpeg", "jpg", "png", "gif", "webp"},
				"max_size_mb":       10,
			},
			{
				"name":                "translation",
				"description":         "Translate text between languages",
				"supported_languages": h.getSupportedLanguagesSimple(),
				"max_text_length":     5000,
			},
			{
				"name":             "summarization",
				"description":      "Summarize long texts",
				"max_input_length": 10000,
				"summary_types":    []string{"paragraph", "bullet_points", "keywords"},
			},
			{
				"name":        "video_generation",
				"description": "Generate videos from text",
				"max_duration": 60,
			},
		},

		"providers": providers,
		"usage_stats": usageStats,
		
		"content_types": []string{
			"blog_post", "article", "social_media", "email",
			"product_description", "ad_copy", "story", "poem",
			"code", "technical_documentation", "business_plan",
		},

		"tones": []string{
			"professional", "casual", "friendly", "formal",
			"persuasive", "informative", "creative", "humorous",
		},

		"limits": gin.H{
			"daily_text_generation":   10000,
			"daily_image_analysis":    50,
			"daily_translations":      100,
			"max_concurrent_requests": 5,
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    capabilities,
	})
}

// ===== دوال مساعدة AI =====

func (h *AIHandler) buildEnhancedPrompt(basePrompt, contentType, tone, length string) string {
	var builder strings.Builder

	builder.WriteString(basePrompt)

	// إضافة توجيهات بناءً على نوع المحتوى
	if contentType != "" {
		builder.WriteString("\n\nPlease generate " + contentType + " content.")
	}

	// إضافة النبرة المطلوبة
	if tone != "" {
		builder.WriteString("\nTone: " + tone)
	}

	// إضافة الطول المطلوب
	if length != "" {
		switch length {
		case "short":
			builder.WriteString("\nPlease be concise (1-2 paragraphs).")
		case "medium":
			builder.WriteString("\nPlease provide detailed content (3-5 paragraphs).")
		case "long":
			builder.WriteString("\nPlease provide comprehensive content (6+ paragraphs).")
		}
	}

	return builder.String()
}

func (h *AIHandler) buildSummaryPrompt(text, summaryType string, maxLength int) string {
	var prompt strings.Builder

	prompt.WriteString("Please summarize the following text")

	switch summaryType {
	case "paragraph":
		prompt.WriteString(" in a single paragraph")
	case "bullet_points":
		prompt.WriteString(" using bullet points")
	case "keywords":
		prompt.WriteString(" by extracting key keywords and phrases")
	}

	prompt.WriteString(". Maximum length: ")
	prompt.WriteString(string(rune(maxLength)))
	prompt.WriteString(" characters.\n\nText to summarize:\n")
	prompt.WriteString(text)

	return prompt.String()
}

func (h *AIHandler) isValidLanguage(lang string) bool {
	supportedLangs := h.getSupportedLanguagesSimple()
	for _, supportedLang := range supportedLangs {
		if supportedLang == lang {
			return true
		}
	}
	return false
}

func (h *AIHandler) getSupportedLanguagesSimple() []string {
	return []string{
		"en", "ar", "fr", "es", "de", "zh",
		"ru", "ja", "ko", "pt", "it", "nl",
		"tr", "fa", "hi", "bn", "ur",
	}
}

// Helper functions
func (h *AIHandler) getUserID(c *gin.Context) string {
	return utils.GetUserIDFromContext(c)
}

// ================================
// دوال مساعدة عامة
// ================================

// getCurrentUserID الحصول على معرف المستخدم الحالي
func getCurrentUserID(c *gin.Context) string {
	// يمكن تعديل هذا بناءً على طريقة المصادقة المستخدمة
	if userID, exists := c.Get("user_id"); exists {
		if id, ok := userID.(string); ok {
			return id
		}
	}
	
	// أو من التوكن إذا كان مخزناً في السياق
	if userID := c.GetString("user_id"); userID != "" {
		return userID
	}
	
	return ""
}

// successResponse إرسال استجابة ناجحة
func successResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    data,
		"error":   nil,
	})
}

// errorResponse إرسال استجابة خطأ
func errorResponse(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, gin.H{
		"success": false,
		"data":    nil,
		"error":   message,
	})
}

// bindAndValidate ربط وتحقق من البيانات
func bindAndValidate(c *gin.Context, data interface{}) bool {
	if err := c.ShouldBindJSON(data); err != nil {
		errorResponse(c, http.StatusBadRequest, fmt.Sprintf("Invalid request data: %v", err))
		return false
	}
	return true
}