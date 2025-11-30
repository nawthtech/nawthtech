package services

import (
	"context"
	"fmt"
	"time"

	"github.com/nawthtech/nawthtech/backend/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// ================================
// هياكل المعاملات المحدثة - تبسيط الهياكل
// ================================

type (
	// ReviewQueryParams هيكل معاملات الاستعراضات
	ReviewQueryParams struct {
		Page   int    `json:"page"`
		Limit  int    `json:"limit"`
		Rating int    `json:"rating"`
		SortBy string `json:"sort_by"`
	}
)

// ================================
// الواجهات الرئيسية (Main Interfaces) - المحدثة والمبسطة
// ================================

type (
	// AuthService واجهة خدمة المصادقة
	AuthService interface {
		Register(ctx context.Context, req AuthRegisterRequest) (*AuthResponse, error)
		Login(ctx context.Context, req AuthLoginRequest) (*AuthResponse, error)
		Logout(ctx context.Context, token string) error
		RefreshToken(ctx context.Context, refreshToken string) (*AuthResponse, error)
		VerifyToken(ctx context.Context, token string) (*TokenClaims, error)
		ForgotPassword(ctx context.Context, email string) error
		ResetPassword(ctx context.Context, token string, newPassword string) error
		ChangePassword(ctx context.Context, userID string, req ChangePasswordRequest) error
	}

	// UserService واجهة خدمة المستخدمين
	UserService interface {
		GetProfile(ctx context.Context, userID string) (*models.User, error)
		UpdateProfile(ctx context.Context, userID string, req UserUpdateRequest) (*models.User, error)
		UpdateAvatar(ctx context.Context, userID string, avatarURL string) error
		DeleteAccount(ctx context.Context, userID string) error
		SearchUsers(ctx context.Context, query string, params UserQueryParams) ([]models.User, error)
		GetUserStats(ctx context.Context, userID string) (*UserStats, error)
	}

	// ServiceService واجهة خدمة الخدمات
	ServiceService interface {
		CreateService(ctx context.Context, req ServiceCreateRequest) (*models.Service, error)
		GetServiceByID(ctx context.Context, serviceID string) (*models.Service, error)
		UpdateService(ctx context.Context, serviceID string, req ServiceUpdateRequest) (*models.Service, error)
		DeleteService(ctx context.Context, serviceID string) error
		GetServices(ctx context.Context, params ServiceQueryParams) ([]models.Service, error)
		SearchServices(ctx context.Context, query string, params ServiceQueryParams) ([]models.Service, error)
		GetFeaturedServices(ctx context.Context) ([]models.Service, error)
		GetSimilarServices(ctx context.Context, serviceID string) ([]models.Service, error)
	}

	// CategoryService واجهة خدمة الفئات
	CategoryService interface {
		GetCategories(ctx context.Context, params CategoryQueryParams) ([]models.Category, error)
		GetCategoryByID(ctx context.Context, categoryID string) (*models.Category, error)
		CreateCategory(ctx context.Context, req CategoryCreateRequest) (*models.Category, error)
		UpdateCategory(ctx context.Context, categoryID string, req CategoryUpdateRequest) (*models.Category, error)
		DeleteCategory(ctx context.Context, categoryID string) error
		GetCategoryTree(ctx context.Context) ([]CategoryNode, error)
	}

	// OrderService واجهة خدمة الطلبات
	OrderService interface {
		CreateOrder(ctx context.Context, req OrderCreateRequest) (*models.Order, error)
		GetOrderByID(ctx context.Context, orderID string) (*models.Order, error)
		GetUserOrders(ctx context.Context, userID string, params OrderQueryParams) ([]models.Order, error)
		UpdateOrderStatus(ctx context.Context, orderID string, status string, notes string) (*models.Order, error)
		CancelOrder(ctx context.Context, orderID string, reason string) (*models.Order, error)
		GetOrderStats(ctx context.Context, timeframe string) (*OrderStats, error)
	}

	// PaymentService واجهة خدمة الدفع
	PaymentService interface {
		CreatePaymentIntent(ctx context.Context, req PaymentIntentRequest) (*PaymentIntent, error)
		ConfirmPayment(ctx context.Context, paymentID string, confirmationData map[string]interface{}) (*PaymentResult, error)
		GetPaymentHistory(ctx context.Context, userID string, params PaymentQueryParams) ([]models.Payment, error)
		ValidatePayment(ctx context.Context, paymentData map[string]interface{}) (*PaymentValidation, error)
	}

	// UploadService واجهة خدمة الرفع
	UploadService interface {
		UploadFile(ctx context.Context, req UploadRequest) (*UploadResult, error)
		DeleteFile(ctx context.Context, fileID string) error
		GetFile(ctx context.Context, fileID string) (*models.File, error)
		GetUserFiles(ctx context.Context, userID string, params FileQueryParams) ([]models.File, error)
		GeneratePresignedURL(ctx context.Context, req PresignedURLRequest) (*PresignedURL, error)
		ValidateFile(ctx context.Context, fileInfo models.File) (*FileValidation, error)
		GetUploadQuota(ctx context.Context, userID string) (*UploadQuota, error)
	}

	// NotificationService واجهة خدمة الإشعارات
	NotificationService interface {
		CreateNotification(ctx context.Context, req NotificationCreateRequest) (*models.Notification, error)
		GetUserNotifications(ctx context.Context, userID string, params NotificationQueryParams) ([]models.Notification, error)
		MarkAsRead(ctx context.Context, notificationID string) error
		MarkAllAsRead(ctx context.Context, userID string) error
		DeleteNotification(ctx context.Context, notificationID string) error
		GetUnreadCount(ctx context.Context, userID string) (int64, error)
	}

	// AdminService واجهة خدمة الإدارة
	AdminService interface {
		GetDashboardStats(ctx context.Context) (*DashboardStats, error)
		GetUsers(ctx context.Context, params UserQueryParams) ([]models.User, error)
		GetSystemLogs(ctx context.Context, params SystemLogQuery) ([]models.SystemLog, error)
		UpdateSystemSettings(ctx context.Context, settings []models.Setting) error
		BanUser(ctx context.Context, userID string, reason string) error
		UnbanUser(ctx context.Context, userID string) error
	}

	// CacheService واجهة خدمة التخزين المؤقت
	CacheService interface {
		Get(key string) (interface{}, error)
		Set(key string, value interface{}, expiration time.Duration) error
		Delete(key string) error
		Exists(key string) (bool, error)
		Flush() error
	}
)

// ================================
// هياكل المعاملات المحدثة والمبسطة
// ================================

type (
	// Auth Structures
	AuthRegisterRequest struct {
		Email     string `json:"email" binding:"required,email"`
		Username  string `json:"username" binding:"required,min=3,max=50"`
		Password  string `json:"password" binding:"required,min=6"`
		FirstName string `json:"first_name" binding:"required"`
		LastName  string `json:"last_name" binding:"required"`
		Phone     string `json:"phone,omitempty"`
	}

	AuthLoginRequest struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	ChangePasswordRequest struct {
		CurrentPassword string `json:"current_password" binding:"required"`
		NewPassword     string `json:"new_password" binding:"required,min=6"`
	}

	// Category Structures
	CategoryQueryParams struct {
		Page     int    `json:"page"`
		Limit    int    `json:"limit"`
		ParentID string `json:"parent_id"`
		Active   *bool  `json:"active"`
		SortBy   string `json:"sort_by"`
	}

	CategoryCreateRequest struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		ParentID    string `json:"parent_id"`
		Icon        string `json:"icon"`
		Color       string `json:"color"`
		Image       string `json:"image,omitempty"`
		SortOrder   int    `json:"sort_order"`
	}

	CategoryUpdateRequest struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
		Color       string `json:"color"`
		Image       string `json:"image"`
		SortOrder   int    `json:"sort_order"`
		Active      *bool  `json:"active"`
	}

	CategoryNode struct {
		Category models.Category `json:"category"`
		Children []CategoryNode  `json:"children"`
	}

	// Order Structures
	OrderCreateRequest struct {
		Items         []OrderItem   `json:"items" binding:"required"`
		ShippingInfo  ShippingInfo  `json:"shipping_info"`
		PaymentMethod string        `json:"payment_method" binding:"required"`
		CustomerNotes string        `json:"customer_notes"`
		CouponCode    string        `json:"coupon_code,omitempty"`
	}

	OrderQueryParams struct {
		Page   int    `json:"page"`
		Limit  int    `json:"limit"`
		Status string `json:"status"`
		SortBy string `json:"sort_by"`
		UserID string `json:"user_id"`
	}

	OrderItem struct {
		ServiceID   string  `json:"service_id" binding:"required"`
		ServiceName string  `json:"service_name" binding:"required"`
		Quantity    int     `json:"quantity" binding:"required,min=1"`
		Price       float64 `json:"price" binding:"required,min=0"`
		Image       string  `json:"image,omitempty"`
	}

	ShippingInfo struct {
		FirstName      string `json:"first_name" binding:"required"`
		LastName       string `json:"last_name" binding:"required"`
		Email          string `json:"email" binding:"required,email"`
		Phone          string `json:"phone" binding:"required"`
		Address        string `json:"address" binding:"required"`
		City           string `json:"city" binding:"required"`
		Country        string `json:"country" binding:"required"`
		PostalCode     string `json:"postal_code" binding:"required"`
		ShippingMethod string `json:"shipping_method" binding:"required"`
	}

	// Payment Structures
	PaymentIntentRequest struct {
		Amount      float64                `json:"amount" binding:"required"`
		Currency    string                 `json:"currency" binding:"required"`
		Description string                 `json:"description"`
		Metadata    map[string]interface{} `json:"metadata"`
		UserID      string                 `json:"user_id" binding:"required"`
	}

	PaymentQueryParams struct {
		Page   int    `json:"page"`
		Limit  int    `json:"limit"`
		Status string `json:"status"`
		UserID string `json:"user_id"`
	}

	// Upload Structures
	UploadRequest struct {
		File        []byte            `json:"file"`
		Filename    string            `json:"filename"`
		ContentType string            `json:"content_type"`
		Size        int64             `json:"size"`
		Metadata    map[string]string `json:"metadata"`
		UserID      string            `json:"user_id"`
		IsPublic    bool              `json:"is_public"`
	}

	FileQueryParams struct {
		Page   int    `json:"page"`
		Limit  int    `json:"limit"`
		Type   string `json:"type"`
		SortBy string `json:"sort_by"`
		UserID string `json:"user_id"`
	}

	PresignedURLRequest struct {
		Filename    string            `json:"filename"`
		ContentType string            `json:"content_type"`
		Size        int64             `json:"size"`
		Metadata    map[string]string `json:"metadata"`
		UserID      string            `json:"user_id"`
	}

	// Service Structures
	ServiceCreateRequest struct {
		Title       string   `json:"title" binding:"required"`
		Description string   `json:"description" binding:"required"`
		Price       float64  `json:"price" binding:"required"`
		Duration    int      `json:"duration" binding:"required"`
		CategoryID  string   `json:"category_id" binding:"required"`
		ProviderID  string   `json:"provider_id" binding:"required"`
		Images      []string `json:"images"`
		Tags        []string `json:"tags"`
	}

	ServiceUpdateRequest struct {
		Title       string   `json:"title"`
		Description string   `json:"description"`
		Price       float64  `json:"price"`
		Duration    int      `json:"duration"`
		CategoryID  string   `json:"category_id"`
		Images      []string `json:"images"`
		Tags        []string `json:"tags"`
		IsActive    *bool    `json:"is_active"`
		IsFeatured  *bool    `json:"is_featured"`
	}

	ServiceQueryParams struct {
		Page       int      `json:"page"`
		Limit      int      `json:"limit"`
		CategoryID string   `json:"category_id"`
		ProviderID string   `json:"provider_id"`
		MinPrice   float64  `json:"min_price"`
		MaxPrice   float64  `json:"max_price"`
		Tags       []string `json:"tags"`
		Featured   *bool    `json:"featured"`
		Active     *bool    `json:"active"`
		SortBy     string   `json:"sort_by"`
	}

	// User Structures
	UserUpdateRequest struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Phone     string `json:"phone"`
		Avatar    string `json:"avatar"`
	}

	UserQueryParams struct {
		Page   int    `json:"page"`
		Limit  int    `json:"limit"`
		Role   string `json:"role"`
		Status string `json:"status"`
		Search string `json:"search"`
	}

	// Notification Structures
	NotificationCreateRequest struct {
		UserID  string                 `json:"user_id" binding:"required"`
		Title   string                 `json:"title" binding:"required"`
		Message string                 `json:"message" binding:"required"`
		Type    string                 `json:"type" binding:"required"`
		Data    map[string]interface{} `json:"data"`
	}

	NotificationQueryParams struct {
		Page   int    `json:"page"`
		Limit  int    `json:"limit"`
		Type   string `json:"type"`
		Unread *bool  `json:"unread"`
	}

	// System Structures
	SystemLogQuery struct {
		Page   int    `json:"page"`
		Limit  int    `json:"limit"`
		Level  string `json:"level"`
		Module string `json:"module"`
		UserID string `json:"user_id"`
	}
)

// ================================
// هياكل النتائج المحدثة والمبسطة
// ================================

type (
	AuthResponse struct {
		User         *models.User `json:"user"`
		AccessToken  string       `json:"access_token"`
		RefreshToken string       `json:"refresh_token"`
		ExpiresAt    time.Time    `json:"expires_at"`
	}

	TokenClaims struct {
		UserID    string    `json:"user_id"`
		Email     string    `json:"email"`
		Role      string    `json:"role"`
		ExpiresAt time.Time `json:"expires_at"`
	}

	PaymentIntent struct {
		ID           string    `json:"id"`
		ClientSecret string    `json:"client_secret"`
		Amount       float64   `json:"amount"`
		Currency     string    `json:"currency"`
		Status       string    `json:"status"`
		CreatedAt    time.Time `json:"created_at"`
	}

	PaymentResult struct {
		ID        string    `json:"id"`
		Status    string    `json:"status"`
		Amount    float64   `json:"amount"`
		Currency  string    `json:"currency"`
		PaidAt    time.Time `json:"paid_at"`
	}

	PaymentValidation struct {
		IsValid bool     `json:"is_valid"`
		Message string   `json:"message"`
		Errors  []string `json:"errors"`
	}

	UploadResult struct {
		ID          string            `json:"id"`
		URL         string            `json:"url"`
		Filename    string            `json:"filename"`
		Size        int64             `json:"size"`
		ContentType string            `json:"content_type"`
		Metadata    map[string]string `json:"metadata"`
		UploadedAt  time.Time         `json:"uploaded_at"`
	}

	PresignedURL struct {
		URL         string    `json:"url"`
		Method      string    `json:"method"`
		ExpiresAt   time.Time `json:"expires_at"`
	}

	FileValidation struct {
		IsValid  bool     `json:"is_valid"`
		Errors   []string `json:"errors"`
		Warnings []string `json:"warnings"`
	}

	UploadQuota struct {
		Used      int64 `json:"used"`
		Total     int64 `json:"total"`
		Remaining int64 `json:"remaining"`
	}

	OrderStats struct {
		TotalOrders    int     `json:"total_orders"`
		PendingOrders  int     `json:"pending_orders"`
		CompletedOrders int    `json:"completed_orders"`
		CanceledOrders int     `json:"canceled_orders"`
		TotalRevenue   float64 `json:"total_revenue"`
		AverageOrderValue float64 `json:"average_order_value"`
	}

	UserStats struct {
		TotalOrders    int       `json:"total_orders"`
		TotalSpent     float64   `json:"total_spent"`
		JoinedDate     time.Time `json:"joined_date"`
		LastOrderDate  time.Time `json:"last_order_date"`
		WishlistCount  int       `json:"wishlist_count"`
	}

	DashboardStats struct {
		TotalUsers      int     `json:"total_users"`
		TotalServices   int     `json:"total_services"`
		TotalOrders     int     `json:"total_orders"`
		TotalRevenue    float64 `json:"total_revenue"`
		PendingOrders   int     `json:"pending_orders"`
		ActiveStores    int     `json:"active_stores"`
	}
)

// ================================
// التطبيقات الفعلية المحدثة للعمل مع MongoDB
// ================================

type (
	authServiceImpl struct {
		db *mongo.Database
	}

	userServiceImpl struct {
		db *mongo.Database
	}

	serviceServiceImpl struct {
		db *mongo.Database
	}

	categoryServiceImpl struct {
		db *mongo.Database
	}

	orderServiceImpl struct {
		db *mongo.Database
	}

	paymentServiceImpl struct {
		db *mongo.Database
	}

	uploadServiceImpl struct {
		db *mongo.Database
	}

	notificationServiceImpl struct {
		db *mongo.Database
	}

	adminServiceImpl struct {
		db *mongo.Database
	}

	cacheServiceImpl struct {
		// implementation details
	}
)

// ================================
// تطبيقات AuthService مع MongoDB
// ================================

func (s *authServiceImpl) Register(ctx context.Context, req AuthRegisterRequest) (*AuthResponse, error) {
	user := &models.User{
		ID:            primitive.NewObjectID(),
		Email:         req.Email,
		Username:      req.Username,
		Password:      "hashed_password", // يجب تشفير كلمة المرور
		FirstName:     req.FirstName,
		LastName:      req.LastName,
		Phone:         req.Phone,
		Role:          "user",
		Status:        "active",
		EmailVerified: false,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// حفظ المستخدم في MongoDB
	_, err := s.db.Collection("users").InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		User:         user,
		AccessToken:  "access_token_" + objectID, _ := primitive.ObjectIDFromHex(user.ID)
objectID.Hex(),
		RefreshToken: "refresh_token_" + objectID, _ := primitive.ObjectIDFromHex(user.ID)
objectID.Hex(),
		ExpiresAt:    time.Now().Add(24 * time.Hour),
	}, nil
}

func (s *authServiceImpl) Login(ctx context.Context, req AuthLoginRequest) (*AuthResponse, error) {
	var user models.User
	err := s.db.Collection("users").FindOne(ctx, bson.M{"email": req.Email}).Decode(&user)
	if err != nil {
		return nil, fmt.Errorf("المستخدم غير موجود")
	}

	// يجب التحقق من كلمة المرور هنا

	return &AuthResponse{
		User:         &user,
		AccessToken:  "access_token_" + objectID, _ := primitive.ObjectIDFromHex(user.ID)
objectID.Hex(),
		RefreshToken: "refresh_token_" + objectID, _ := primitive.ObjectIDFromHex(user.ID)
objectID.Hex(),
		ExpiresAt:    time.Now().Add(24 * time.Hour),
	}, nil
}

func (s *authServiceImpl) Logout(ctx context.Context, token string) error {
	// يمكن تنفيذ إدارة الجلسات هنا
	return nil
}

func (s *authServiceImpl) RefreshToken(ctx context.Context, refreshToken string) (*AuthResponse, error) {
	// تنفيذ تجديد التوكن
	return &AuthResponse{
		User: &models.User{
			ID:           primitive.NewObjectID(),
			Email:        "user@example.com",
			Username:     "user123",
			Role:         "user",
			Status:       "active",
			EmailVerified: true,
		},
		AccessToken:  "new_access_token_123",
		RefreshToken: "new_refresh_token_123",
		ExpiresAt:    time.Now().Add(24 * time.Hour),
	}, nil
}

func (s *authServiceImpl) VerifyToken(ctx context.Context, token string) (*TokenClaims, error) {
	return &TokenClaims{
		UserID:    "user_123",
		Email:     "user@example.com",
		Role:      "user",
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}, nil
}

func (s *authServiceImpl) ForgotPassword(ctx context.Context, email string) error {
	// تنفيذ إعادة تعيين كلمة المرور
	return nil
}

func (s *authServiceImpl) ResetPassword(ctx context.Context, token string, newPassword string) error {
	// تنفيذ إعادة تعيين كلمة المرور
	return nil
}

func (s *authServiceImpl) ChangePassword(ctx context.Context, userID string, req ChangePasswordRequest) error {
	// تنفيذ تغيير كلمة المرور
	return nil
}

// ================================
// تطبيقات UserService مع MongoDB
// ================================

func (s *userServiceImpl) GetProfile(ctx context.Context, userID string) (*models.User, error) {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, fmt.Errorf("معرف مستخدم غير صالح")
	}

	var user models.User
	err = s.db.Collection("users").FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		return nil, fmt.Errorf("المستخدم غير موجود")
	}

	return &user, nil
}

func (s *userServiceImpl) UpdateProfile(ctx context.Context, userID string, req UserUpdateRequest) (*models.User, error) {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, fmt.Errorf("معرف مستخدم غير صالح")
	}

	updateFields := bson.M{
		"first_name": req.FirstName,
		"last_name":  req.LastName,
		"phone":      req.Phone,
		"avatar":     req.Avatar,
		"updated_at": time.Now(),
	}

	_, err = s.db.Collection("users").UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		bson.M{"$set": updateFields},
	)
	if err != nil {
		return nil, err
	}

	// جلب المستخدم المحدث
	var user models.User
	err = s.db.Collection("users").FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *userServiceImpl) UpdateAvatar(ctx context.Context, userID string, avatarURL string) error {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return fmt.Errorf("معرف مستخدم غير صالح")
	}

	_, err = s.db.Collection("users").UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		bson.M{"$set": bson.M{"avatar": avatarURL, "updated_at": time.Now()}},
	)
	return err
}

func (s *userServiceImpl) DeleteAccount(ctx context.Context, userID string) error {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return fmt.Errorf("معرف مستخدم غير صالح")
	}

	_, err = s.db.Collection("users").DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}

func (s *userServiceImpl) SearchUsers(ctx context.Context, query string, params UserQueryParams) ([]models.User, error) {
	filter := bson.M{}
	if query != "" {
		filter["$or"] = []bson.M{
			{"first_name": bson.M{"$regex": query, "$options": "i"}},
			{"last_name": bson.M{"$regex": query, "$options": "i"}},
			{"email": bson.M{"$regex": query, "$options": "i"}},
			{"username": bson.M{"$regex": query, "$options": "i"}},
		}
	}

	cursor, err := s.db.Collection("users").Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []models.User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, err
	}

	return users, nil
}

func (s *userServiceImpl) GetUserStats(ctx context.Context, userID string) (*UserStats, error) {
	return &UserStats{
		TotalOrders:    5,
		TotalSpent:     1500.0,
		JoinedDate:     time.Now().AddDate(0, -6, 0),
		LastOrderDate:  time.Now().AddDate(0, 0, -5),
		WishlistCount:  3,
	}, nil
}

// ================================
// تطبيقات ServiceService مع MongoDB
// ================================

func (s *serviceServiceImpl) CreateService(ctx context.Context, req ServiceCreateRequest) (*models.Service, error) {
	service := &models.Service{
		ID:          primitive.NewObjectID(),
		Title:       req.Title,
		Description: req.Description,
		Price:       req.Price,
		Duration:    req.Duration,
		CategoryID:  req.CategoryID,
		ProviderID:  req.ProviderID,
		Images:      req.Images,
		Tags:        req.Tags,
		IsActive:    true,
		IsFeatured:  false,
		Rating:      0,
		ReviewCount: 0,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	_, err := s.db.Collection("services").InsertOne(ctx, service)
	if err != nil {
		return nil, err
	}

	return service, nil
}

func (s *serviceServiceImpl) GetServiceByID(ctx context.Context, serviceID string) (*models.Service, error) {
	objectID, err := primitive.ObjectIDFromHex(serviceID)
	if err != nil {
		return nil, fmt.Errorf("معرف خدمة غير صالح")
	}

	var service models.Service
	err = s.db.Collection("services").FindOne(ctx, bson.M{"_id": objectID}).Decode(&service)
	if err != nil {
		return nil, fmt.Errorf("الخدمة غير موجودة")
	}

	return &service, nil
}

func (s *serviceServiceImpl) UpdateService(ctx context.Context, serviceID string, req ServiceUpdateRequest) (*models.Service, error) {
	objectID, err := primitive.ObjectIDFromHex(serviceID)
	if err != nil {
		return nil, fmt.Errorf("معرف خدمة غير صالح")
	}

	updateFields := bson.M{
		"title":       req.Title,
		"description": req.Description,
		"price":       req.Price,
		"duration":    req.Duration,
		"category_id": req.CategoryID,
		"images":      req.Images,
		"tags":        req.Tags,
		"updated_at":  time.Now(),
	}

	if req.IsActive != nil {
		updateFields["is_active"] = *req.IsActive
	}
	if req.IsFeatured != nil {
		updateFields["is_featured"] = *req.IsFeatured
	}

	_, err = s.db.Collection("services").UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		bson.M{"$set": updateFields},
	)
	if err != nil {
		return nil, err
	}

	// جلب الخدمة المحدثة
	var service models.Service
	err = s.db.Collection("services").FindOne(ctx, bson.M{"_id": objectID}).Decode(&service)
	if err != nil {
		return nil, err
	}

	return &service, nil
}

func (s *serviceServiceImpl) DeleteService(ctx context.Context, serviceID string) error {
	objectID, err := primitive.ObjectIDFromHex(serviceID)
	if err != nil {
		return fmt.Errorf("معرف خدمة غير صالح")
	}

	_, err = s.db.Collection("services").DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}

func (s *serviceServiceImpl) GetServices(ctx context.Context, params ServiceQueryParams) ([]models.Service, error) {
	filter := bson.M{"is_active": true}

	if params.CategoryID != "" {
		filter["category_id"] = params.CategoryID
	}
	if params.ProviderID != "" {
		filter["provider_id"] = params.ProviderID
	}
	if params.Featured != nil {
		filter["is_featured"] = *params.Featured
	}

	cursor, err := s.db.Collection("services").Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var services []models.Service
	if err := cursor.All(ctx, &services); err != nil {
		return nil, err
	}

	return services, nil
}

func (s *serviceServiceImpl) SearchServices(ctx context.Context, query string, params ServiceQueryParams) ([]models.Service, error) {
	filter := bson.M{
		"$and": []bson.M{
			{"is_active": true},
			{"$or": []bson.M{
				{"title": bson.M{"$regex": query, "$options": "i"}},
				{"description": bson.M{"$regex": query, "$options": "i"}},
				{"tags": bson.M{"$in": []string{query}}},
			}},
		},
	}

	cursor, err := s.db.Collection("services").Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var services []models.Service
	if err := cursor.All(ctx, &services); err != nil {
		return nil, err
	}

	return services, nil
}

func (s *serviceServiceImpl) GetFeaturedServices(ctx context.Context) ([]models.Service, error) {
	cursor, err := s.db.Collection("services").Find(ctx, bson.M{
		"is_active":  true,
		"is_featured": true,
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var services []models.Service
	if err := cursor.All(ctx, &services); err != nil {
		return nil, err
	}

	return services, nil
}

func (s *serviceServiceImpl) GetSimilarServices(ctx context.Context, serviceID string) ([]models.Service, error) {
	// جلب الخدمة الحالية لمعرفة فئتها
	service, err := s.GetServiceByID(ctx, serviceID)
	if err != nil {
		return nil, err
	}

	// جلب خدمات مشابهة من نفس الفئة
	cursor, err := s.db.Collection("services").Find(ctx, bson.M{
		"category_id": service.CategoryID,
		"is_active":   true,
		"_id":         bson.M{"$ne": service.ID},
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var services []models.Service
	if err := cursor.All(ctx, &services); err != nil {
		return nil, err
	}

	return services, nil
}

// ================================
// تطبيقات CategoryService مع MongoDB
// ================================

func (s *categoryServiceImpl) GetCategories(ctx context.Context, params CategoryQueryParams) ([]models.Category, error) {
	filter := bson.M{}
	if params.Active != nil {
		filter["is_active"] = *params.Active
	}

	cursor, err := s.db.Collection("categories").Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var categories []models.Category
	if err := cursor.All(ctx, &categories); err != nil {
		return nil, err
	}

	return categories, nil
}

func (s *categoryServiceImpl) GetCategoryByID(ctx context.Context, categoryID string) (*models.Category, error) {
	objectID, err := primitive.ObjectIDFromHex(categoryID)
	if err != nil {
		return nil, fmt.Errorf("معرف فئة غير صالح")
	}

	var category models.Category
	err = s.db.Collection("categories").FindOne(ctx, bson.M{"_id": objectID}).Decode(&category)
	if err != nil {
		return nil, fmt.Errorf("الفئة غير موجودة")
	}

	return &category, nil
}

func (s *categoryServiceImpl) CreateCategory(ctx context.Context, req CategoryCreateRequest) (*models.Category, error) {
	category := &models.Category{
		ID:          primitive.NewObjectID(),
		Name:        req.Name,
		Description: req.Description,
		ParentID:    req.ParentID,
		Icon:        req.Icon,
		Color:       req.Color,
		Image:       req.Image,
		SortOrder:   req.SortOrder,
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	_, err := s.db.Collection("categories").InsertOne(ctx, category)
	if err != nil {
		return nil, err
	}

	return category, nil
}

func (s *categoryServiceImpl) UpdateCategory(ctx context.Context, categoryID string, req CategoryUpdateRequest) (*models.Category, error) {
	objectID, err := primitive.ObjectIDFromHex(categoryID)
	if err != nil {
		return nil, fmt.Errorf("معرف فئة غير صالح")
	}

	updateFields := bson.M{
		"name":        req.Name,
		"description": req.Description,
		"icon":        req.Icon,
		"color":       req.Color,
		"image":       req.Image,
		"sort_order":  req.SortOrder,
		"updated_at":  time.Now(),
	}

	if req.Active != nil {
		updateFields["is_active"] = *req.Active
	}

	_, err = s.db.Collection("categories").UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		bson.M{"$set": updateFields},
	)
	if err != nil {
		return nil, err
	}

	// جلب الفئة المحدثة
	var category models.Category
	err = s.db.Collection("categories").FindOne(ctx, bson.M{"_id": objectID}).Decode(&category)
	if err != nil {
		return nil, err
	}

	return &category, nil
}

func (s *categoryServiceImpl) DeleteCategory(ctx context.Context, categoryID string) error {
	objectID, err := primitive.ObjectIDFromHex(categoryID)
	if err != nil {
		return fmt.Errorf("معرف فئة غير صالح")
	}

	_, err = s.db.Collection("categories").DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}

func (s *categoryServiceImpl) GetCategoryTree(ctx context.Context) ([]CategoryNode, error) {
	// تنفيذ شجرة الفئات
	return []CategoryNode{}, nil
}

// ================================
// تطبيقات OrderService مع MongoDB
// ================================

func (s *orderServiceImpl) CreateOrder(ctx context.Context, req OrderCreateRequest) (*models.Order, error) {
	order := &models.Order{
		ID:           primitive.NewObjectID(),
		UserID:       "user_id_from_context", // سيتم تعيينه من السياق
		Items:        []models.OrderItem{},
		Status:       "pending",
		TotalAmount:  0,
		Discount:     0,
		Tax:          0,
		Shipping:     0,
		FinalAmount:  0,
		PaymentStatus: "pending",
		PaymentMethod: req.PaymentMethod,
		CustomerNotes: req.CustomerNotes,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// حساب المبلغ الإجمالي
	for _, item := range req.Items {
		order.TotalAmount += item.Price * float64(item.Quantity)
		order.Items = append(order.Items, models.OrderItem{
			ID:          primitive.NewObjectID(),
			ServiceID:   item.ServiceID,
			ServiceName: item.ServiceName,
			Quantity:    item.Quantity,
			Price:       item.Price,
			Image:       item.Image,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		})
	}

	order.FinalAmount = order.TotalAmount - order.Discount + order.Tax + order.Shipping

	_, err := s.db.Collection("orders").InsertOne(ctx, order)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (s *orderServiceImpl) GetOrderByID(ctx context.Context, orderID string) (*models.Order, error) {
	objectID, err := primitive.ObjectIDFromHex(orderID)
	if err != nil {
		return nil, fmt.Errorf("معرف طلب غير صالح")
	}

	var order models.Order
	err = s.db.Collection("orders").FindOne(ctx, bson.M{"_id": objectID}).Decode(&order)
	if err != nil {
		return nil, fmt.Errorf("الطلب غير موجود")
	}

	return &order, nil
}

func (s *orderServiceImpl) GetUserOrders(ctx context.Context, userID string, params OrderQueryParams) ([]models.Order, error) {
	filter := bson.M{"user_id": userID}
	if params.Status != "" {
		filter["status"] = params.Status
	}

	cursor, err := s.db.Collection("orders").Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var orders []models.Order
	if err := cursor.All(ctx, &orders); err != nil {
		return nil, err
	}

	return orders, nil
}

func (s *orderServiceImpl) UpdateOrderStatus(ctx context.Context, orderID string, status string, notes string) (*models.Order, error) {
	objectID, err := primitive.ObjectIDFromHex(orderID)
	if err != nil {
		return nil, fmt.Errorf("معرف طلب غير صالح")
	}

	_, err = s.db.Collection("orders").UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		bson.M{"$set": bson.M{
			"status":     status,
			"updated_at": time.Now(),
		}},
	)
	if err != nil {
		return nil, err
	}

	return s.GetOrderByID(ctx, orderID)
}

func (s *orderServiceImpl) CancelOrder(ctx context.Context, orderID string, reason string) (*models.Order, error) {
	return s.UpdateOrderStatus(ctx, orderID, "cancelled", reason)
}

func (s *orderServiceImpl) GetOrderStats(ctx context.Context, timeframe string) (*OrderStats, error) {
	return &OrderStats{
		TotalOrders:    50,
		PendingOrders:  5,
		CompletedOrders: 40,
		CanceledOrders: 5,
		TotalRevenue:   15000,
		AverageOrderValue: 300,
	}, nil
}

// ================================
// تطبيقات الخدمات الأخرى المبسطة
// ================================

func (s *paymentServiceImpl) CreatePaymentIntent(ctx context.Context, req PaymentIntentRequest) (*PaymentIntent, error) {
	return &PaymentIntent{
		ID:           fmt.Sprintf("pi_%d", time.Now().Unix()),
		ClientSecret: "secret_" + fmt.Sprintf("%d", time.Now().Unix()),
		Amount:       req.Amount,
		Currency:     req.Currency,
		Status:       "requires_payment_method",
		CreatedAt:    time.Now(),
	}, nil
}

func (s *paymentServiceImpl) ConfirmPayment(ctx context.Context, paymentID string, confirmationData map[string]interface{}) (*PaymentResult, error) {
	return &PaymentResult{
		ID:       paymentID,
		Status:   "succeeded",
		Amount:   100,
		Currency: "USD",
		PaidAt:   time.Now(),
	}, nil
}

func (s *paymentServiceImpl) GetPaymentHistory(ctx context.Context, userID string, params PaymentQueryParams) ([]models.Payment, error) {
	return []models.Payment{}, nil
}

func (s *paymentServiceImpl) ValidatePayment(ctx context.Context, paymentData map[string]interface{}) (*PaymentValidation, error) {
	return &PaymentValidation{
		IsValid: true,
		Message: "الدفع صالح",
		Errors:  []string{},
	}, nil
}

func (s *uploadServiceImpl) UploadFile(ctx context.Context, req UploadRequest) (*UploadResult, error) {
	return &UploadResult{
		ID:          primitive.NewObjectID().Hex(),
		URL:         "https://res.cloudinary.com/nawthtech/image/upload/v123/example.jpg",
		Filename:    req.Filename,
		Size:        req.Size,
		ContentType: req.ContentType,
		Metadata:    req.Metadata,
		UploadedAt:  time.Now(),
	}, nil
}

func (s *uploadServiceImpl) DeleteFile(ctx context.Context, fileID string) error {
	return nil
}

func (s *uploadServiceImpl) GetFile(ctx context.Context, fileID string) (*models.File, error) {
	return &models.File{
		ID:       fileID,
		Filename: "example.jpg",
		Size:     1024,
		URL:      "https://example.com/file.jpg",
		UserID:   "user_123",
	}, nil
}

func (s *uploadServiceImpl) GetUserFiles(ctx context.Context, userID string, params FileQueryParams) ([]models.File, error) {
	return []models.File{}, nil
}

func (s *uploadServiceImpl) GeneratePresignedURL(ctx context.Context, req PresignedURLRequest) (*PresignedURL, error) {
	return &PresignedURL{
		URL:       "https://example.com/presigned-url",
		Method:    "PUT",
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}, nil
}

func (s *uploadServiceImpl) ValidateFile(ctx context.Context, fileInfo models.File) (*FileValidation, error) {
	return &FileValidation{
		IsValid:  true,
		Errors:   []string{},
		Warnings: []string{},
	}, nil
}

func (s *uploadServiceImpl) GetUploadQuota(ctx context.Context, userID string) (*UploadQuota, error) {
	return &UploadQuota{
		Used:      500,
		Total:     1000,
		Remaining: 500,
	}, nil
}

func (s *notificationServiceImpl) CreateNotification(ctx context.Context, req NotificationCreateRequest) (*models.Notification, error) {
	return &models.Notification{
		ID:        primitive.NewObjectID(),
		UserID:    req.UserID,
		Title:     req.Title,
		Message:   req.Message,
		Type:      req.Type,
		IsRead:    false,
		CreatedAt: time.Now(),
	}, nil
}

func (s *notificationServiceImpl) GetUserNotifications(ctx context.Context, userID string, params NotificationQueryParams) ([]models.Notification, error) {
	return []models.Notification{}, nil
}

func (s *notificationServiceImpl) MarkAsRead(ctx context.Context, notificationID string) error {
	return nil
}

func (s *notificationServiceImpl) MarkAllAsRead(ctx context.Context, userID string) error {
	return nil
}

func (s *notificationServiceImpl) DeleteNotification(ctx context.Context, notificationID string) error {
	return nil
}

func (s *notificationServiceImpl) GetUnreadCount(ctx context.Context, userID string) (int64, error) {
	return 0, nil
}

func (s *adminServiceImpl) GetDashboardStats(ctx context.Context) (*DashboardStats, error) {
	return &DashboardStats{
		TotalUsers:    150,
		TotalServices: 89,
		TotalOrders:   234,
		TotalRevenue:  15499.99,
		PendingOrders: 15,
		ActiveStores:  45,
	}, nil
}

func (s *adminServiceImpl) GetUsers(ctx context.Context, params UserQueryParams) ([]models.User, error) {
	return []models.User{}, nil
}

func (s *adminServiceImpl) GetSystemLogs(ctx context.Context, params SystemLogQuery) ([]models.SystemLog, error) {
	return []models.SystemLog{}, nil
}

func (s *adminServiceImpl) UpdateSystemSettings(ctx context.Context, settings []models.Setting) error {
	return nil
}

func (s *adminServiceImpl) BanUser(ctx context.Context, userID string, reason string) error {
	return nil
}

func (s *adminServiceImpl) UnbanUser(ctx context.Context, userID string) error {
	return nil
}

// ================================
// تطبيقات CacheService
// ================================

func (s *cacheServiceImpl) Get(key string) (interface{}, error) {
	return nil, nil
}

func (s *cacheServiceImpl) Set(key string, value interface{}, expiration time.Duration) error {
	return nil
}

func (s *cacheServiceImpl) Delete(key string) error {
	return nil
}

func (s *cacheServiceImpl) Exists(key string) (bool, error) {
	return false, nil
}

func (s *cacheServiceImpl) Flush() error {
	return nil
}

// ================================
// دوال الإنشاء المحدثة للعمل مع MongoDB
// ================================

func NewAuthService(db *mongo.Database) AuthService {
	return &authServiceImpl{db: db}
}

func NewUserService(db *mongo.Database) UserService {
	return &userServiceImpl{db: db}
}

func NewServiceService(db *mongo.Database) ServiceService {
	return &serviceServiceImpl{db: db}
}

func NewCategoryService(db *mongo.Database) CategoryService {
	return &categoryServiceImpl{db: db}
}

func NewOrderService(db *mongo.Database) OrderService {
	return &orderServiceImpl{db: db}
}

func NewPaymentService(db *mongo.Database) PaymentService {
	return &paymentServiceImpl{db: db}
}

func NewUploadService(db *mongo.Database) UploadService {
	return &uploadServiceImpl{db: db}
}

func NewNotificationService(db *mongo.Database) NotificationService {
	return &notificationServiceImpl{db: db}
}

func NewAdminService(db *mongo.Database) AdminService {
	return &adminServiceImpl{db: db}
}

func NewCacheService() CacheService {
	return &cacheServiceImpl{}
}

// ================================
// Service Container المحدث للعمل مع MongoDB
// ================================

type ServiceContainer struct {
	Auth         AuthService
	User         UserService
	Service      ServiceService
	Category     CategoryService
	Order        OrderService
	Payment      PaymentService
	Upload       UploadService
	Notification NotificationService
	Admin        AdminService
	Cache        CacheService
}

func NewServiceContainer(mongoClient *mongo.Client, databaseName string) *ServiceContainer {
	db := mongoClient.Database(databaseName)
	
	return &ServiceContainer{
		Auth:         NewAuthService(db),
		User:         NewUserService(db),
		Service:      NewServiceService(db),
		Category:     NewCategoryService(db),
		Order:        NewOrderService(db),
		Payment:      NewPaymentService(db),
		Upload:       NewUploadService(db),
		Notification: NewNotificationService(db),
		Admin:        NewAdminService(db),
		Cache:        NewCacheService(),
	}
}