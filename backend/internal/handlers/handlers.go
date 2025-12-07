package handlers

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nawthtech/nawthtech/backend/internal/cloudinary"
	"github.com/nawthtech/nawthtech/backend/internal/services"
	"github.com/nawthtech/nawthtech/backend/internal/utils"
	"github.com/nawthtech/nawthtech/backend/internal/ai"
	"github.com/nawthtech/nawthtech/backend/internal/ai/video"
)

// ================================
// الواجهات الأساسية للمعالجات
// ================================

type (
	// AuthHandler معالج المصادقة
	AuthHandler interface {
		Register(c *gin.Context)
		Login(c *gin.Context)
		Logout(c *gin.Context)
		RefreshToken(c *gin.Context)
		ForgotPassword(c *gin.Context)
		ResetPassword(c *gin.Context)
		VerifyToken(c *gin.Context)
	}

	// UserHandler معالج المستخدم
	UserHandler interface {
		GetProfile(c *gin.Context)
		UpdateProfile(c *gin.Context)
		ChangePassword(c *gin.Context)
		GetUserStats(c *gin.Context)
	}

	// ServiceHandler معالج الخدمات
	ServiceHandler interface {
		GetServices(c *gin.Context)
		GetServiceByID(c *gin.Context)
		SearchServices(c *gin.Context)
		GetFeaturedServices(c *gin.Context)
		GetCategories(c *gin.Context)
		CreateService(c *gin.Context)
		UpdateService(c *gin.Context)
		DeleteService(c *gin.Context)
		GetMyServices(c *gin.Context)
		GetSellerOrders(c *gin.Context)  // أضفنا هذه الدالة
		GetAllOrders(c *gin.Context)     // أضفنا هذه الدالة
	}

	// CategoryHandler معالج الفئات
	CategoryHandler interface {
		GetCategories(c *gin.Context)
		GetCategoryByID(c *gin.Context)
		CreateCategory(c *gin.Context)
		UpdateCategory(c *gin.Context)
		DeleteCategory(c *gin.Context)
	}

	// OrderHandler معالج الطلبات
	OrderHandler interface {
		CreateOrder(c *gin.Context)
		GetOrderByID(c *gin.Context)
		GetUserOrders(c *gin.Context)
		UpdateOrderStatus(c *gin.Context)
		CancelOrder(c *gin.Context)
		GetSellerOrders(c *gin.Context) // أضفنا هذه الدالة
		GetAllOrders(c *gin.Context)    // أضفنا هذه الدالة
	}

	// PaymentHandler معالج الدفع
	PaymentHandler interface {
		CreatePaymentIntent(c *gin.Context)
		ConfirmPayment(c *gin.Context)
		GetPaymentHistory(c *gin.Context)
		HandleStripeWebhook(c *gin.Context)  // أضفنا هذه الدالة
		HandlePayPalWebhook(c *gin.Context)  // أضفنا هذه الدالة
	}

	// UploadHandler معالج الرفع
	UploadHandler interface {
		UploadImage(c *gin.Context)
		UploadMultipleImages(c *gin.Context)
		DeleteImage(c *gin.Context)
		GetImageInfo(c *gin.Context)
		GetUserImages(c *gin.Context)
		HandleCloudinaryWebhook(c *gin.Context)  // أضفنا هذه الدالة
	}

	// NotificationHandler معالج الإشعارات
	NotificationHandler interface {
		GetUserNotifications(c *gin.Context)
		MarkAsRead(c *gin.Context)
		MarkAllAsRead(c *gin.Context)
		GetUnreadCount(c *gin.Context)
	}

	// AdminHandler معالج الإدارة
	AdminHandler interface {
		GetDashboard(c *gin.Context)
		GetDashboardStats(c *gin.Context)
		GetUsers(c *gin.Context)
		UpdateUserStatus(c *gin.Context)
		GetSystemLogs(c *gin.Context)
		GetAllOrders(c *gin.Context)  // أضفنا هذه الدالة
	}

	// AIHandler معالج الذكاء الاصطناعي
	AIHandler interface {
		GenerateContentHandler(c *gin.Context)
		AnalyzeImageHandler(c *gin.Context)
		TranslateTextHandler(c *gin.Context)
		SummarizeTextHandler(c *gin.Context)
		GetAICapabilitiesHandler(c *gin.Context)
	}

	// VideoHandler معالج الفيديو
	VideoHandler interface {
		GenerateVideoHandler(c *gin.Context)
		GetVideoStatusHandler(c *gin.Context)
		ListVideoJobsHandler(c *gin.Context)
		CancelVideoJobHandler(c *gin.Context)
		DownloadVideoHandler(c *gin.Context)
		GetVideoCapabilitiesHandler(c *gin.Context)
		GetVideoStatsHandler(c *gin.Context)
	}
)

// ================================
// التطبيقات الفعلية للمعالجات
// ================================

type (
	authHandler struct {
		authService services.AuthService
	}

	userHandler struct {
		userService services.UserService
	}

	serviceHandler struct {
		serviceService services.ServiceService
	}

	categoryHandler struct {
		categoryService services.CategoryService
	}

	orderHandler struct {
		orderService services.OrderService
	}

	paymentHandler struct {
		paymentService services.PaymentService
	}

	uploadHandler struct {
		cloudinaryService *cloudinary.CloudinaryService
	}

	notificationHandler struct {
		notificationService services.NotificationService
	}

	adminHandler struct {
		adminService services.AdminService
	}

	aiHandler struct {
		aiClient *ai.Client
	}

	videoHandler struct {
		videoService *video.VideoService
	}
)

// ================================
// دوال إنشاء المعالجات
// ================================

func NewAuthHandler(authService services.AuthService) AuthHandler {
	return &authHandler{authService: authService}
}

func NewUserHandler(userService services.UserService) UserHandler {
	return &userHandler{userService: userService}
}

func NewServiceHandler(serviceService services.ServiceService) ServiceHandler {
	return &serviceHandler{serviceService: serviceService}
}

func NewCategoryHandler(categoryService services.CategoryService) CategoryHandler {
	return &categoryHandler{categoryService: categoryService}
}

func NewOrderHandler(orderService services.OrderService) OrderHandler {
	return &orderHandler{orderService: orderService}
}

func NewPaymentHandler(paymentService services.PaymentService) PaymentHandler {
	return &paymentHandler{paymentService: paymentService}
}

func NewUploadHandler() (UploadHandler, error) {
	cloudinaryService, err := cloudinary.NewCloudinaryService()
	if err != nil {
		return nil, err
	}
	return &uploadHandler{cloudinaryService: cloudinaryService}, nil
}

func NewUploadHandlerWithService(cloudinaryService *cloudinary.CloudinaryService) UploadHandler {
	return &uploadHandler{cloudinaryService: cloudinaryService}
}

func NewNotificationHandler(notificationService services.NotificationService) NotificationHandler {
	return &notificationHandler{notificationService: notificationService}
}

func NewAdminHandler(adminService services.AdminService) AdminHandler {
	return &adminHandler{adminService: adminService}
}

func NewAIHandler(aiClient *ai.Client) AIHandler {
	return &aiHandler{aiClient: aiClient}
}

func NewVideoHandler(videoService *video.VideoService) VideoHandler {
	return &videoHandler{videoService: videoService}
}

// ================================
// التطبيقات الأساسية للمعالجات
// ================================

// AuthHandler implementations
func (h *authHandler) Register(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Register endpoint - MongoDB Ready",
		"data":    gin.H{"database": "MongoDB"},
	})
}

func (h *authHandler) Login(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Login endpoint - MongoDB Ready",
		"data":    gin.H{"database": "MongoDB"},
	})
}

func (h *authHandler) Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Logout endpoint",
	})
}

func (h *authHandler) RefreshToken(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Refresh token endpoint",
	})
}

func (h *authHandler) ForgotPassword(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Forgot password endpoint",
	})
}

func (h *authHandler) ResetPassword(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Reset password endpoint",
	})
}

func (h *authHandler) VerifyToken(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Verify token endpoint",
	})
}

// UserHandler implementations
func (h *userHandler) GetProfile(c *gin.Context) {
	// الحصول على معرف المستخدم من السياق (من middleware المصادقة)
	userID := utils.GetUserIDFromContext(c)
	if userID == "" {
		utils.ErrorResponse(c, http.StatusUnauthorized, "يجب تسجيل الدخول", "UNAUTHORIZED")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Get profile endpoint - MongoDB Ready",
		"data": gin.H{
			"user": gin.H{
				"id":        userID,
				"name":      "نوذ تك",
				"email":     "info@nawthtech.com",
				"createdAt": time.Now().Format(time.RFC3339),
			},
			"database": "MongoDB",
		},
	})
}

func (h *userHandler) UpdateProfile(c *gin.Context) {
	userID := utils.GetUserIDFromContext(c)
	if userID == "" {
		utils.ErrorResponse(c, http.StatusUnauthorized, "يجب تسجيل الدخول", "UNAUTHORIZED")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Update profile endpoint",
		"user_id": userID,
	})
}

func (h *userHandler) ChangePassword(c *gin.Context) {
	userID := utils.GetUserIDFromContext(c)
	if userID == "" {
		utils.ErrorResponse(c, http.StatusUnauthorized, "يجب تسجيل الدخول", "UNAUTHORIZED")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Change password endpoint",
		"user_id": userID,
	})
}

func (h *userHandler) GetUserStats(c *gin.Context) {
	userID := utils.GetUserIDFromContext(c)
	if userID == "" {
		utils.ErrorResponse(c, http.StatusUnauthorized, "يجب تسجيل الدخول", "UNAUTHORIZED")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Get user stats endpoint",
		"data": gin.H{
			"user_id":         userID,
			"total_services":  15,
			"total_orders":    47,
			"joined_date":     "2023-01-15",
			"account_status":  "active",
		},
	})
}

// ServiceHandler implementations
func (h *serviceHandler) GetServices(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Get services endpoint - MongoDB Ready",
		"data": gin.H{
			"services": []gin.H{
				{
					"id":          "service1",
					"title":       "تطوير واجهات المستخدم",
					"description": "تصميم وتطوير واجهات مستخدم تفاعلية",
					"price":       299.99,
					"category":    "تطوير الويب",
				},
				{
					"id":          "service2",
					"title":       "تطبيقات الجوال",
					"description": "تطوير تطبيقات جوال مبتكرة",
					"price":       499.99,
					"category":    "تطبيقات الجوال",
				},
			},
			"database": "MongoDB",
		},
	})
}

func (h *serviceHandler) GetServiceByID(c *gin.Context) {
	serviceID := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Get service by ID endpoint",
		"data": gin.H{
			"id":          serviceID,
			"title":       "خدمة مثال",
			"description": "وصف الخدمة",
			"price":       199.99,
			"database":    "MongoDB",
		},
	})
}

func (h *serviceHandler) SearchServices(c *gin.Context) {
	query := c.Query("q")
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Search services endpoint",
		"data": gin.H{
			"query":    query,
			"results":  []gin.H{},
			"database": "MongoDB",
		},
	})
}

func (h *serviceHandler) GetFeaturedServices(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Get featured services endpoint",
		"data": gin.H{
			"featured_services": []gin.H{},
			"database":          "MongoDB",
		},
	})
}

func (h *serviceHandler) GetCategories(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Get categories endpoint",
		"data": gin.H{
			"categories": []string{"تطوير الويب", "تطبيقات الجوال", "تصميم جرافيك"},
			"database":   "MongoDB",
		},
	})
}

func (h *serviceHandler) CreateService(c *gin.Context) {
	userID := utils.GetUserIDFromContext(c)
	if userID == "" {
		utils.ErrorResponse(c, http.StatusUnauthorized, "يجب تسجيل الدخول", "UNAUTHORIZED")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Create service endpoint",
		"data": gin.H{
			"service_id": "new_service_123",
			"created_by": userID,
			"database":   "MongoDB",
		},
	})
}

func (h *serviceHandler) UpdateService(c *gin.Context) {
	serviceID := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Update service endpoint",
		"data": gin.H{
			"service_id": serviceID,
			"database":   "MongoDB",
		},
	})
}

func (h *serviceHandler) DeleteService(c *gin.Context) {
	serviceID := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Delete service endpoint",
		"data": gin.H{
			"deleted_id": serviceID,
			"database":   "MongoDB",
		},
	})
}

func (h *serviceHandler) GetMyServices(c *gin.Context) {
	userID := utils.GetUserIDFromContext(c)
	if userID == "" {
		utils.ErrorResponse(c, http.StatusUnauthorized, "يجب تسجيل الدخول", "UNAUTHORIZED")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Get my services endpoint",
		"data": gin.H{
			"my_services": []gin.H{},
			"user_id":     userID,
			"database":    "MongoDB",
		},
	})
}

func (h *serviceHandler) GetSellerOrders(c *gin.Context) {
	userID := utils.GetUserIDFromContext(c)
	if userID == "" {
		utils.ErrorResponse(c, http.StatusUnauthorized, "يجب تسجيل الدخول", "UNAUTHORIZED")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Get seller orders endpoint",
		"data": gin.H{
			"orders":   []gin.H{},
			"seller_id": userID,
			"database": "MongoDB",
		},
	})
}

func (h *serviceHandler) GetAllOrders(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Get all orders endpoint (admin)",
		"data": gin.H{
			"orders":   []gin.H{},
			"database": "MongoDB",
		},
	})
}

// CategoryHandler implementations
func (h *categoryHandler) GetCategories(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Get categories endpoint",
		"data": gin.H{
			"categories": []gin.H{
				{"id": "cat1", "name": "تطوير الويب", "service_count": 15},
				{"id": "cat2", "name": "تطبيقات الجوال", "service_count": 8},
				{"id": "cat3", "name": "تصميم جرافيك", "service_count": 12},
			},
			"database": "MongoDB",
		},
	})
}

func (h *categoryHandler) GetCategoryByID(c *gin.Context) {
	categoryID := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Get category by ID endpoint",
		"data": gin.H{
			"id":            categoryID,
			"name":          "تطوير الويب",
			"description":   "خدمات تطوير الويب المختلفة",
			"service_count": 15,
			"database":      "MongoDB",
		},
	})
}

func (h *categoryHandler) CreateCategory(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Create category endpoint",
		"data": gin.H{
			"category_id": "new_category_123",
			"database":    "MongoDB",
		},
	})
}

func (h *categoryHandler) UpdateCategory(c *gin.Context) {
	categoryID := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Update category endpoint",
		"data": gin.H{
			"category_id": categoryID,
			"database":    "MongoDB",
		},
	})
}

func (h *categoryHandler) DeleteCategory(c *gin.Context) {
	categoryID := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Delete category endpoint",
		"data": gin.H{
			"deleted_id": categoryID,
			"database":   "MongoDB",
		},
	})
}

// OrderHandler implementations
func (h *orderHandler) CreateOrder(c *gin.Context) {
	userID := utils.GetUserIDFromContext(c)
	if userID == "" {
		utils.ErrorResponse(c, http.StatusUnauthorized, "يجب تسجيل الدخول", "UNAUTHORIZED")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Create order endpoint",
		"data": gin.H{
			"order_id":  "order_123",
			"status":    "pending",
			"user_id":   userID,
			"database":  "MongoDB",
		},
	})
}

func (h *orderHandler) GetOrderByID(c *gin.Context) {
	orderID := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Get order by ID endpoint",
		"data": gin.H{
			"id":       orderID,
			"status":   "completed",
			"amount":   299.99,
			"database": "MongoDB",
		},
	})
}

func (h *orderHandler) GetUserOrders(c *gin.Context) {
	userID := utils.GetUserIDFromContext(c)
	if userID == "" {
		utils.ErrorResponse(c, http.StatusUnauthorized, "يجب تسجيل الدخول", "UNAUTHORIZED")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Get user orders endpoint",
		"data": gin.H{
			"orders":   []gin.H{},
			"user_id":  userID,
			"database": "MongoDB",
		},
	})
}

func (h *orderHandler) UpdateOrderStatus(c *gin.Context) {
	orderID := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Update order status endpoint",
		"data": gin.H{
			"order_id": orderID,
			"status":   "updated",
			"database": "MongoDB",
		},
	})
}

func (h *orderHandler) CancelOrder(c *gin.Context) {
	orderID := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Cancel order endpoint",
		"data": gin.H{
			"order_id": orderID,
			"status":   "cancelled",
			"database": "MongoDB",
		},
	})
}

func (h *orderHandler) GetSellerOrders(c *gin.Context) {
	userID := utils.GetUserIDFromContext(c)
	if userID == "" {
		utils.ErrorResponse(c, http.StatusUnauthorized, "يجب تسجيل الدخول", "UNAUTHORIZED")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Get seller orders endpoint",
		"data": gin.H{
			"orders":    []gin.H{},
			"seller_id": userID,
			"database":  "MongoDB",
		},
	})
}

func (h *orderHandler) GetAllOrders(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Get all orders endpoint (admin)",
		"data": gin.H{
			"orders":   []gin.H{},
			"database": "MongoDB",
		},
	})
}

// PaymentHandler implementations
func (h *paymentHandler) CreatePaymentIntent(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Create payment intent endpoint",
		"data": gin.H{
			"client_secret": "pi_123_secret_456",
			"database":      "MongoDB",
		},
	})
}

func (h *paymentHandler) ConfirmPayment(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Confirm payment endpoint",
		"data": gin.H{
			"status":   "succeeded",
			"database": "MongoDB",
		},
	})
}

func (h *paymentHandler) GetPaymentHistory(c *gin.Context) {
	userID := utils.GetUserIDFromContext(c)
	if userID == "" {
		utils.ErrorResponse(c, http.StatusUnauthorized, "يجب تسجيل الدخول", "UNAUTHORIZED")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Get payment history endpoint",
		"data": gin.H{
			"payments": []gin.H{},
			"user_id":  userID,
			"database": "MongoDB",
		},
	})
}

func (h *paymentHandler) HandleStripeWebhook(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Stripe webhook endpoint",
		"data": gin.H{
			"received": true,
			"database": "MongoDB",
		},
	})
}

func (h *paymentHandler) HandlePayPalWebhook(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "PayPal webhook endpoint",
		"data": gin.H{
			"received": true,
			"database": "MongoDB",
		},
	})
}

// UploadHandler implementations - Cloudinary Integration
func (h *uploadHandler) UploadImage(c *gin.Context) {
	// التحقق من وجود الملف
	file, err := c.FormFile("image")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "لم يتم توفير ملف صورة", "NO_FILE_PROVIDED")
		return
	}

	// التحقق من صحة الملف
	if err := h.cloudinaryService.ValidateImage(file); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), "INVALID_FILE")
		return
	}

	// الحصول على public_id من الطلب أو إنشاء واحد تلقائي
	publicID := c.PostForm("public_id")
	if publicID == "" {
		publicID = h.cloudinaryService.GeneratePublicID("img")
	}

	// رفع الصورة إلى Cloudinary
	result, err := h.cloudinaryService.UploadImageFromGinFile(c, "image", cloudinary.UploadOptions{
		PublicID:  publicID,
		Folder:    "nawthtech/uploads",
		Overwrite: true,
	})
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في رفع الصورة", "UPLOAD_FAILED")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "تم رفع الصورة بنجاح", gin.H{
		"image_url":     result.SecureURL,
		"public_id":     result.PublicID,
		"format":        result.Format,
		"size_bytes":    result.Bytes,
		"width":         result.Width,
		"height":        result.Height,
		"resource_type": result.ResourceType,
	})
}

func (h *uploadHandler) UploadMultipleImages(c *gin.Context) {
	// الحصول على النموذج متعدد الأجزاء
	form, err := c.MultipartForm()
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "طلب غير صالح", "INVALID_FORM")
		return
	}

	files := form.File["images"]
	if len(files) == 0 {
		utils.ErrorResponse(c, http.StatusBadRequest, "لم يتم توفير ملفات", "NO_FILES_PROVIDED")
		return
	}

	var results []gin.H
	var errors []string

	for i, file := range files {
		// التحقق من صحة كل ملف
		if err := h.cloudinaryService.ValidateImage(file); err != nil {
			errors = append(errors, fmt.Sprintf("الملف %s: %s", file.Filename, err.Error()))
			continue
		}

		// إنشاء public_id فريد
		publicID := h.cloudinaryService.GeneratePublicID("img")

		// رفع الصورة
		result, err := h.cloudinaryService.UploadImageFromGinFile(c, "images", cloudinary.UploadOptions{
			PublicID:  publicID,
			Folder:    "nawthtech/uploads",
			Overwrite: true,
		})
		if err != nil {
			errors = append(errors, fmt.Sprintf("الملف %s: %s", file.Filename, err.Error()))
			continue
		}

		results = append(results, gin.H{
			"filename":     file.Filename,
			"image_url":    result.SecureURL,
			"public_id":    result.PublicID,
			"format":       result.Format,
			"size_bytes":   result.Bytes,
			"width":        result.Width,
			"height":       result.Height,
			"upload_index": i,
		})
	}

	response := gin.H{
		"uploaded":       results,
		"total_uploaded": len(results),
		"total_failed":   len(errors),
	}

	if len(errors) > 0 {
		response["errors"] = errors
		utils.SuccessResponse(c, http.StatusPartialContent, "تم رفع بعض الملفات بنجاح مع وجود أخطاء", response)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "تم رفع جميع الملفات بنجاح", response)
}

func (h *uploadHandler) DeleteImage(c *gin.Context) {
	publicID := c.Param("public_id")

	if publicID == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "معرف الصورة مطلوب", "MISSING_PUBLIC_ID")
		return
	}

	err := h.cloudinaryService.DeleteImage(publicID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في حذف الصورة", "DELETE_FAILED")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "تم حذف الصورة بنجاح", gin.H{
		"public_id": publicID,
		"deleted":   true,
	})
}

func (h *uploadHandler) GetImageInfo(c *gin.Context) {
	publicID := c.Param("public_id")

	if publicID == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "معرف الصورة مطلوب", "MISSING_PUBLIC_ID")
		return
	}

	// في التطبيق الحقيقي، قد تحتاج إلى تنفيذ دالة GetImageInfo في CloudinaryService
	// حالياً سنستخدم دالة بسيطة للاستجابة
	utils.SuccessResponse(c, http.StatusOK, "معلومات الصورة", gin.H{
		"public_id":     publicID,
		"message":       "معلومات الصورة - هذه الدالة تحتاج إلى تنفيذ في CloudinaryService",
		"resource_type": "image",
	})
}

func (h *uploadHandler) GetUserImages(c *gin.Context) {
	userID := utils.GetUserIDFromContext(c)

	if userID == "" {
		utils.ErrorResponse(c, http.StatusUnauthorized, "يجب تسجيل الدخول", "UNAUTHORIZED")
		return
	}

	// في التطبيق الحقيقي، يمكن جلب الصور من قاعدة البيانات
	// هذا مثال بسيط للاستجابة
	utils.SuccessResponse(c, http.StatusOK, "صور المستخدم", gin.H{
		"user_id": userID,
		"images":  []gin.H{},
		"total":   0,
	})
}

func (h *uploadHandler) HandleCloudinaryWebhook(c *gin.Context) {
	// معالجة webhook من Cloudinary
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Cloudinary webhook received",
		"data":    gin.H{"received": true},
	})
}

// NotificationHandler implementations
func (h *notificationHandler) GetUserNotifications(c *gin.Context) {
	userID := utils.GetUserIDFromContext(c)
	if userID == "" {
		utils.ErrorResponse(c, http.StatusUnauthorized, "يجب تسجيل الدخول", "UNAUTHORIZED")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Get user notifications endpoint",
		"data": gin.H{
			"notifications": []gin.H{},
			"user_id":       userID,
			"database":      "MongoDB",
		},
	})
}

func (h *notificationHandler) MarkAsRead(c *gin.Context) {
	notificationID := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Mark as read endpoint",
		"data": gin.H{
			"notification_id": notificationID,
			"read":            true,
			"database":        "MongoDB",
		},
	})
}

func (h *notificationHandler) MarkAllAsRead(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Mark all as read endpoint",
		"data": gin.H{
			"marked_all": true,
			"database":   "MongoDB",
		},
	})
}

func (h *notificationHandler) GetUnreadCount(c *gin.Context) {
	userID := utils.GetUserIDFromContext(c)
	if userID == "" {
		utils.ErrorResponse(c, http.StatusUnauthorized, "يجب تسجيل الدخول", "UNAUTHORIZED")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Get unread count endpoint",
		"data": gin.H{
			"unread_count": 0,
			"user_id":      userID,
			"database":     "MongoDB",
		},
	})
}

// AdminHandler implementations
func (h *adminHandler) GetDashboard(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Get dashboard endpoint",
		"data": gin.H{
			"stats": gin.H{
				"total_users":    150,
				"total_services": 89,
				"total_orders":   234,
				"revenue":        15499.99,
			},
			"database": "MongoDB",
		},
	})
}

func (h *adminHandler) GetDashboardStats(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Get dashboard stats endpoint",
		"data": gin.H{
			"users": gin.H{
				"total":    150,
				"active":   132,
				"inactive": 18,
			},
			"services": gin.H{
				"total":   89,
				"active":  76,
				"pending": 13,
			},
			"database": "MongoDB",
		},
	})
}

func (h *adminHandler) GetUsers(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Get users endpoint",
		"data": gin.H{
			"users":    []gin.H{},
			"database": "MongoDB",
		},
	})
}

func (h *adminHandler) UpdateUserStatus(c *gin.Context) {
	userID := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Update user status endpoint",
		"data": gin.H{
			"user_id":  userID,
			"status":   "updated",
			"database": "MongoDB",
		},
	})
}

func (h *adminHandler) GetSystemLogs(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Get system logs endpoint",
		"data": gin.H{
			"logs":     []gin.H{},
			"database": "MongoDB",
		},
	})
}

func (h *adminHandler) GetAllOrders(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Get all orders endpoint (admin)",
		"data": gin.H{
			"orders":   []gin.H{},
			"database": "MongoDB",
		},
	})
}

// AIHandler implementations (ستكون في ai_handler.go، هنا فقط للتوافق)
func (h *aiHandler) GenerateContentHandler(c *gin.Context) {
	// سيتم تنفيذها في ai_handler.go
	c.JSON(http.StatusNotImplemented, gin.H{
		"success": false,
		"error":   "AI handler not implemented in this file",
	})
}

func (h *aiHandler) AnalyzeImageHandler(c *gin.Context) {
	// سيتم تنفيذها في ai_handler.go
	c.JSON(http.StatusNotImplemented, gin.H{
		"success": false,
		"error":   "AI handler not implemented in this file",
	})
}

func (h *aiHandler) TranslateTextHandler(c *gin.Context) {
	// سيتم تنفيذها في ai_handler.go
	c.JSON(http.StatusNotImplemented, gin.H{
		"success": false,
		"error":   "AI handler not implemented in this file",
	})
}

func (h *aiHandler) SummarizeTextHandler(c *gin.Context) {
	// سيتم تنفيذها في ai_handler.go
	c.JSON(http.StatusNotImplemented, gin.H{
		"success": false,
		"error":   "AI handler not implemented in this file",
	})
}

func (h *aiHandler) GetAICapabilitiesHandler(c *gin.Context) {
	// سيتم تنفيذها في ai_handler.go
	c.JSON(http.StatusNotImplemented, gin.H{
		"success": false,
		"error":   "AI handler not implemented in this file",
	})
}

// VideoHandler implementations (ستكون في video_handler.go، هنا فقط للتوافق)
func (h *videoHandler) GenerateVideoHandler(c *gin.Context) {
	// سيتم تنفيذها في video_handler.go
	c.JSON(http.StatusNotImplemented, gin.H{
		"success": false,
		"error":   "Video handler not implemented in this file",
	})
}

func (h *videoHandler) GetVideoStatusHandler(c *gin.Context) {
	// سيتم تنفيذها في video_handler.go
	c.JSON(http.StatusNotImplemented, gin.H{
		"success": false,
		"error":   "Video handler not implemented in this file",
	})
}

func (h *videoHandler) ListVideoJobsHandler(c *gin.Context) {
	// سيتم تنفيذها في video_handler.go
	c.JSON(http.StatusNotImplemented, gin.H{
		"success": false,
		"error":   "Video handler not implemented in this file",
	})
}

func (h *videoHandler) CancelVideoJobHandler(c *gin.Context) {
	// سيتم تنفيذها في video_handler.go
	c.JSON(http.StatusNotImplemented, gin.H{
		"success": false,
		"error":   "Video handler not implemented in this file",
	})
}

func (h *videoHandler) DownloadVideoHandler(c *gin.Context) {
	// سيتم تنفيذها في video_handler.go
	c.JSON(http.StatusNotImplemented, gin.H{
		"success": false,
		"error":   "Video handler not implemented in this file",
	})
}

func (h *videoHandler) GetVideoCapabilitiesHandler(c *gin.Context) {
	// سيتم تنفيذها في video_handler.go
	c.JSON(http.StatusNotImplemented, gin.H{
		"success": false,
		"error":   "Video handler not implemented in this file",
	})
}

func (h *videoHandler) GetVideoStatsHandler(c *gin.Context) {
	// سيتم تنفيذها في video_handler.go
	c.JSON(http.StatusNotImplemented, gin.H{
		"success": false,
		"error":   "Video handler not implemented in this file",
	})
}

// ================================
// دوال مساعدة إضافية
// ================================

// Helper function for multiple uploads
func (h *uploadHandler) processMultipleUploads(c *gin.Context, files []*multipart.FileHeader) ([]gin.H, []string) {
	var results []gin.H
	var errors []string

	for i, file := range files {
		// التحقق من صحة كل ملف
		if err := h.cloudinaryService.ValidateImage(file); err != nil {
			errors = append(errors, fmt.Sprintf("الملف %s: %s", file.Filename, err.Error()))
			continue
		}

		// إنشاء public_id فريد
		publicID := h.cloudinaryService.GeneratePublicID("img")

		// رفع الصورة
		result, err := h.cloudinaryService.UploadImageFromGinFile(c, "images", cloudinary.UploadOptions{
			PublicID:  publicID,
			Folder:    "nawthtech/uploads",
			Overwrite: true,
		})
		if err != nil {
			errors = append(errors, fmt.Sprintf("الملف %s: %s", file.Filename, err.Error()))
			continue
		}

		results = append(results, gin.H{
			"filename":     file.Filename,
			"image_url":    result.SecureURL,
			"public_id":    result.PublicID,
			"format":       result.Format,
			"size_bytes":   result.Bytes,
			"width":        result.Width,
			"height":       result.Height,
			"upload_index": i,
		})
	}

	return results, errors
}