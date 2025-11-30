package v1

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nawthtech/nawthtech/backend/internal/utils"
	"github.com/nawthtech/nawthtech/backend/internal/logger"
	"golang.org/x/crypto/bcrypt"
)

// ================================
// 🏷️ الأنواع والمخططات (Types)
// ================================

// APIResponse استجابة API الموحدة
type APIResponse struct {
	Success   bool        `json:"success"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Error     string      `json:"error,omitempty"`
	Meta      interface{} `json:"meta,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

// PaginatedResponse استجابة paginated
type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Pagination struct {
		Page       int   `json:"page"`
		Limit      int   `json:"limit"`
		Total      int64 `json:"total"`
		TotalPages int   `json:"total_pages"`
		HasNext    bool  `json:"has_next"`
		HasPrev    bool  `json:"has_prev"`
	} `json:"pagination"`
}

// ErrorResponse استجابة الخطأ
type ErrorResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Error   string      `json:"error"`
	Code    string      `json:"code,omitempty"`
	Details interface{} `json:"details,omitempty"`
}

// UploadResult نتيجة الرفع من Cloudinary
type UploadResult struct {
	PublicID     string `json:"public_id"`
	SecureURL    string `json:"secure_url"`
	Format       string `json:"format"`
	Bytes        int    `json:"bytes"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
	ResourceType string `json:"resource_type"`
}

// ================================
// ✅ المحققات (Validators)
// ================================

// RegisterRequest هيكل طلب التسجيل
type RegisterRequest struct {
	FirstName string `json:"first_name" validate:"required,min=2,max=50"`
	LastName  string `json:"last_name" validate:"required,min=2,max=50"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8"`
	Phone     string `json:"phone" validate:"required,phone"`
}

// ValidateRegisterRequest التحقق من صحة طلب التسجيل
func ValidateRegisterRequest(c *gin.Context) (*RegisterRequest, []utils.ValidationError) {
	var req RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, []utils.ValidationError{
			{Field: "request", Message: "بيانات الطلب غير صالحة"},
		}
	}

	errors := utils.ValidateStruct(&req)
	if len(errors) > 0 {
		return nil, errors
	}

	return &req, nil
}

// LoginRequest هيكل طلب تسجيل الدخول
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// ValidateLoginRequest التحقق من صحة طلب تسجيل الدخول
func ValidateLoginRequest(c *gin.Context) (*LoginRequest, []utils.ValidationError) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, []utils.ValidationError{
			{Field: "request", Message: "بيانات الطلب غير صالحة"},
		}
	}

	errors := utils.ValidateStruct(&req)
	if len(errors) > 0 {
		return nil, errors
	}

	return &req, nil
}

// UploadImageRequest هيكل طلب رفع الصورة
type UploadImageRequest struct {
	PublicID  string `form:"public_id"`
	Folder    string `form:"folder"`
	Overwrite bool   `form:"overwrite"`
}

// ValidateUploadImageRequest التحقق من طلب رفع الصورة
func ValidateUploadImageRequest(c *gin.Context) (*UploadImageRequest, []utils.ValidationError) {
	var req UploadImageRequest

	if err := c.ShouldBind(&req); err != nil {
		return nil, []utils.ValidationError{
			{Field: "request", Message: "بيانات الطلب غير صالحة"},
		}
	}

	if req.Folder == "" {
		req.Folder = "nawthtech/uploads"
	}

	return &req, nil
}

// CreateServiceRequest هيكل طلب إنشاء خدمة
type CreateServiceRequest struct {
	Title       string   `json:"title" validate:"required,min=5,max=100"`
	Description string   `json:"description" validate:"required,min=10,max=1000"`
	Price       float64  `json:"price" validate:"required,min=0"`
	CategoryID  string   `json:"category_id" validate:"required"`
	Tags        []string `json:"tags,omitempty"`
}

// ValidateCreateServiceRequest التحقق من طلب إنشاء خدمة
func ValidateCreateServiceRequest(c *gin.Context) (*CreateServiceRequest, []utils.ValidationError) {
	var req CreateServiceRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, []utils.ValidationError{
			{Field: "request", Message: "بيانات الطلب غير صالحة"},
		}
	}

	errors := utils.ValidateStruct(&req)
	if len(errors) > 0 {
		return nil, errors
	}

	return &req, nil
}

// ================================
// 🛡️ الوسائط (Middleware)
// ================================

// APIResponseMiddleware وسيط لتوحيد استجابات API
func APIResponseMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start)
		status := c.Writer.Status()

		if !strings.Contains(c.Request.URL.Path, "/docs") && !strings.Contains(c.Request.URL.Path, "/health") {
			logger.Info(c.Request.Context(), "استجابة API",
				"method", c.Request.Method,
				"path", c.Request.URL.Path,
				"status", status,
				"duration", duration,
				"client_ip", c.ClientIP(),
			)
		}
	}
}

// APIVersionMiddleware وسيط لإدارة إصدارات API
func APIVersionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-API-Version", "v1")
		c.Header("X-API-Service", "NawthTech Backend")
		c.Header("X-API-Timestamp", time.Now().UTC().Format(time.RFC3339))
		c.Next()
	}
}

// AuthMiddleware وسيط المصادقة المبسط
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			ErrorResponse(c, http.StatusUnauthorized, "مطلوب مصادقة", "UNAUTHORIZED")
			c.Abort()
			return
		}

		if strings.HasPrefix(token, "Bearer ") {
			token = strings.TrimPrefix(token, "Bearer ")
		}

		// التحقق من التوكن (تنفيذ مبسط)
		c.Set("userID", "user123")
		c.Set("userRole", "user")

		c.Next()
	}
}

// AdminMiddleware وسيط التحقق من صلاحيات المدير
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("userRole")
		if !exists || userRole != "admin" {
			ErrorResponse(c, http.StatusForbidden, "غير مسموح بالوصول", "FORBIDDEN")
			c.Abort()
			return
		}
		c.Next()
	}
}

// ================================
// 🛠️ دوال المساعدة (Helpers)
// ================================

// SuccessResponse إرسال استجابة ناجحة موحدة
func SuccessResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	response := APIResponse{
		Success:   true,
		Message:   message,
		Data:      data,
		Timestamp: time.Now().UTC(),
	}
	c.JSON(statusCode, response)
}

// ErrorResponse إرسال استجابة خطأ موحدة
func ErrorResponse(c *gin.Context, statusCode int, message string, errorCode string) {
	response := APIResponse{
		Success:   false,
		Message:   message,
		Error:     errorCode,
		Timestamp: time.Now().UTC(),
	}
	c.JSON(statusCode, response)
}

// ValidationErrorResponse إرسال استجابة أخطاء تحقق
func ValidationErrorResponse(c *gin.Context, errors []utils.ValidationError) {
	response := APIResponse{
		Success:   false,
		Message:   "أخطاء في التحقق من البيانات",
		Error:     "VALIDATION_ERROR",
		Data:      errors,
		Timestamp: time.Now().UTC(),
	}
	c.JSON(http.StatusBadRequest, response)
}

// HashPassword تشفير كلمة المرور
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPassword التحقق من كلمة المرور
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GetUserIDFromContext الحصول على معرف المستخدم من السياق
func GetUserIDFromContext(c *gin.Context) string {
	if userID, exists := c.Get("userID"); exists {
		if id, ok := userID.(string); ok {
			return id
		}
	}
	return ""
}

// ValidateFileUpload التحقق من ملف الرفع
func ValidateFileUpload(c *gin.Context, fieldName string, maxSize int64, allowedTypes []string) error {
	file, err := c.FormFile(fieldName)
	if err != nil {
		return err
	}

	if file.Size > maxSize {
		return fmt.Errorf("حجم الملف يتجاوز الحد المسموح به: %d bytes", maxSize)
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowed := false
	for _, allowedType := range allowedTypes {
		if strings.EqualFold(ext, allowedType) {
			allowed = true
			break
		}
	}

	if !allowed {
		return fmt.Errorf("نوع الملف غير مسموح به. الأنواع المسموحة: %v", allowedTypes)
	}

	return nil
}

// GetPaginationParams الحصول على معاملات الترقيم
func GetPaginationParams(c *gin.Context) (page, limit int) {
	page, _ = strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ = strconv.Atoi(c.DefaultQuery("limit", "20"))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	return page, limit
}