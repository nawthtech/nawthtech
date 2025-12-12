package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nawthtech/nawthtech/backend/internal/models"
	"github.com/nawthtech/nawthtech/backend/internal/services"
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
}

// NewHandlerContainer إنشاء حاوية handlers جديدة
func NewHandlerContainer(serviceContainer *services.ServiceContainer) *HandlerContainer {
	container := &HandlerContainer{}

	if serviceContainer != nil {
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
	}

	return container
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
// دوال مساعدة
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