package services

import (
	"context"
	"fmt"
	"time"

	"github.com/nawthtech/nawthtech/backend/internal/models"
	"github.com/nawthtech/nawthtech/backend/internal/utils"
	"gorm.io/gorm"
)

// ================================
// الواجهات الرئيسية (Main Interfaces) - المحدثة
// ================================

type (
	// AnalyticsService واجهة خدمة التحليلات
	AnalyticsService interface {
		TrackEvent(ctx context.Context, event models.Analytics) error
		GetUserAnalytics(ctx context.Context, userID string, timeframe string) (*UserAnalytics, error)
		GetServiceAnalytics(ctx context.Context, serviceID string, timeframe string) (*ServiceAnalytics, error)
		GetPlatformAnalytics(ctx context.Context, timeframe string) (*PlatformAnalytics, error)
	}

	// AdminService واجهة خدمة الإدارة
	AdminService interface {
		GetDashboardStats(ctx context.Context) (*DashboardStats, error)
		GetUsers(ctx context.Context, params UserQueryParams) ([]models.User, *utils.Pagination, error)
		GetSystemLogs(ctx context.Context, params SystemLogQuery) ([]models.SystemLog, *utils.Pagination, error)
		UpdateSystemSettings(ctx context.Context, settings []models.Setting) error
		BanUser(ctx context.Context, userID string, reason string) error
		UnbanUser(ctx context.Context, userID string) error
	}

	// ContentService واجهة خدمة المحتوى
	ContentService interface {
		CreateContent(ctx context.Context, req ContentCreateRequest) (*models.Content, error)
		GetContentByID(ctx context.Context, contentID string) (*models.Content, error)
		GetContentBySlug(ctx context.Context, slug string) (*models.Content, error)
		UpdateContent(ctx context.Context, contentID string, req ContentUpdateRequest) (*models.Content, error)
		DeleteContent(ctx context.Context, contentID string) error
		GetContentList(ctx context.Context, params ContentQueryParams) ([]models.Content, *utils.Pagination, error)
		PublishContent(ctx context.Context, contentID string) error
		UnpublishContent(ctx context.Context, contentID string) error
	}

	// NotificationService واجهة خدمة الإشعارات
	NotificationService interface {
		CreateNotification(ctx context.Context, req NotificationCreateRequest) (*models.Notification, error)
		GetUserNotifications(ctx context.Context, userID string, params NotificationQueryParams) ([]models.Notification, *utils.Pagination, error)
		MarkAsRead(ctx context.Context, notificationID string) error
		MarkAllAsRead(ctx context.Context, userID string) error
		DeleteNotification(ctx context.Context, notificationID string) error
		GetUnreadCount(ctx context.Context, userID string) (int64, error)
		SendBulkNotification(ctx context.Context, req BulkNotificationRequest) error
	}

	// UserService واجهة خدمة المستخدمين
	UserService interface {
		GetProfile(ctx context.Context, userID string) (*models.User, error)
		UpdateProfile(ctx context.Context, userID string, req UserUpdateRequest) (*models.User, error)
		UpdateAvatar(ctx context.Context, userID string, avatarURL string) error
		DeleteAccount(ctx context.Context, userID string) error
		SearchUsers(ctx context.Context, query string, params UserQueryParams) ([]models.User, *utils.Pagination, error)
		GetUserStats(ctx context.Context, userID string) (*UserStats, error)
	}

	// ServiceService واجهة خدمة الخدمات
	ServiceService interface {
		CreateService(ctx context.Context, req ServiceCreateRequest) (*models.Service, error)
		GetServiceByID(ctx context.Context, serviceID string) (*models.Service, error)
		UpdateService(ctx context.Context, serviceID string, req ServiceUpdateRequest) (*models.Service, error)
		DeleteService(ctx context.Context, serviceID string) error
		GetServices(ctx context.Context, params ServiceQueryParams) ([]models.Service, *utils.Pagination, error)
		SearchServices(ctx context.Context, query string, params ServiceQueryParams) ([]models.Service, *utils.Pagination, error)
		GetFeaturedServices(ctx context.Context) ([]models.Service, error)
		GetSimilarServices(ctx context.Context, serviceID string) ([]models.Service, error)
	}

	// CacheService واجهة خدمة التخزين المؤقت
	CacheService interface {
		Get(key string) (interface{}, error)
		Set(key string, value interface{}, expiration time.Duration) error
		Delete(key string) error
		Exists(key string) (bool, error)
		Flush() error
	}

	// AIService واجهة خدمة الذكاء الاصطناعي
	AIService interface {
		GenerateText(ctx context.Context, params AIGenerateParams) (*AIGenerationResult, error)
		AnalyzeSentiment(ctx context.Context, text string, language string) (*SentimentAnalysis, error)
		ClassifyContent(ctx context.Context, content string, categories []string) (*ContentClassification, error)
		ExtractKeywords(ctx context.Context, text string, maxKeywords int) (*KeywordExtraction, error)
		SummarizeText(ctx context.Context, text string, maxLength int) (*TextSummary, error)
		TranslateText(ctx context.Context, text string, sourceLang string, targetLang string) (*TranslationResult, error)
		GenerateImage(ctx context.Context, params AIImageParams) (*AIImageResult, error)
		ChatCompletion(ctx context.Context, messages []AIChatMessage, model string) (*ChatCompletionResult, error)
	}

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
		ValidateSession(ctx context.Context, sessionID string) (*SessionInfo, error)
	}

	// CartService واجهة خدمة عربة التسوق
	CartService interface {
		GetCart(ctx context.Context, userID string) (*models.Cart, error)
		AddToCart(ctx context.Context, userID string, item models.CartItem) (*models.Cart, error)
		UpdateCartItem(ctx context.Context, userID string, itemID string, quantity int) (*models.Cart, error)
		RemoveFromCart(ctx context.Context, userID string, itemID string) (*models.Cart, error)
		ClearCart(ctx context.Context, userID string) error
		GetCartSummary(ctx context.Context, userID string) (*CartSummary, error)
		ApplyCoupon(ctx context.Context, userID string, couponCode string) (*models.Cart, error)
		RemoveCoupon(ctx context.Context, userID string) (*models.Cart, error)
	}

	// CategoryService واجهة خدمة الفئات
	CategoryService interface {
		GetCategories(ctx context.Context, params CategoryQueryParams) ([]models.Category, *utils.Pagination, error)
		GetCategoryByID(ctx context.Context, categoryID string) (*models.Category, error)
		CreateCategory(ctx context.Context, req CategoryCreateRequest) (*models.Category, error)
		UpdateCategory(ctx context.Context, categoryID string, req CategoryUpdateRequest) (*models.Category, error)
		DeleteCategory(ctx context.Context, categoryID string) error
		GetCategoryTree(ctx context.Context) ([]CategoryNode, error)
		GetCategoryStats(ctx context.Context) (*CategoryStats, error)
	}

	// OrderService واجهة خدمة الطلبات
	OrderService interface {
		CreateOrder(ctx context.Context, req OrderCreateRequest) (*models.Order, error)
		GetOrderByID(ctx context.Context, orderID string) (*models.Order, error)
		GetUserOrders(ctx context.Context, userID string, params OrderQueryParams) ([]models.Order, *utils.Pagination, error)
		UpdateOrderStatus(ctx context.Context, orderID string, status string, notes string) (*models.Order, error)
		CancelOrder(ctx context.Context, orderID string, reason string) (*models.Order, error)
		GetOrderStats(ctx context.Context, timeframe string) (*OrderStats, error)
		ProcessOrderPayment(ctx context.Context, orderID string, paymentInfo PaymentInfo) (*OrderPaymentResult, error)
		TrackOrder(ctx context.Context, orderID string) (*OrderTracking, error)
	}

	// PaymentService واجهة خدمة الدفع
	PaymentService interface {
		CreatePaymentIntent(ctx context.Context, req PaymentIntentRequest) (*PaymentIntent, error)
		ConfirmPayment(ctx context.Context, paymentID string, confirmationData map[string]interface{}) (*PaymentResult, error)
		RefundPayment(ctx context.Context, paymentID string, amount float64, reason string) (*RefundResult, error)
		GetPaymentMethods(ctx context.Context, userID string) ([]PaymentMethod, error)
		AddPaymentMethod(ctx context.Context, userID string, method PaymentMethod) error
		RemovePaymentMethod(ctx context.Context, userID string, methodID string) error
		GetPaymentHistory(ctx context.Context, userID string, params PaymentQueryParams) ([]models.Payment, *utils.Pagination, error)
		ValidatePayment(ctx context.Context, paymentData map[string]interface{}) (*PaymentValidation, error)
	}

	// ReportService واجهة خدمة التقارير
	ReportService interface {
		GenerateSalesReport(ctx context.Context, params ReportParams) (*SalesReport, error)
		GenerateUserReport(ctx context.Context, params ReportParams) (*UserReport, error)
		GenerateServiceReport(ctx context.Context, params ReportParams) (*ServiceReport, error)
		GenerateFinancialReport(ctx context.Context, params ReportParams) (*FinancialReport, error)
		GenerateSystemReport(ctx context.Context, params ReportParams) (*SystemReport, error)
		GetReportTemplates(ctx context.Context) ([]ReportTemplate, error)
		ScheduleReport(ctx context.Context, req ScheduleReportRequest) (*ScheduledReport, error)
		GetScheduledReports(ctx context.Context, params ScheduledReportQuery) ([]ScheduledReport, *utils.Pagination, error)
	}

	// StoreService واجهة خدمة المتجر
	StoreService interface {
		GetStoreByID(ctx context.Context, storeID string) (*models.Store, error)
		GetStoreBySlug(ctx context.Context, slug string) (*models.Store, error)
		CreateStore(ctx context.Context, req StoreCreateRequest) (*models.Store, error)
		UpdateStore(ctx context.Context, storeID string, req StoreUpdateRequest) (*models.Store, error)
		DeleteStore(ctx context.Context, storeID string) error
		GetStoreStats(ctx context.Context, storeID string) (*StoreStats, error)
		GetStoreReviews(ctx context.Context, storeID string, params ReviewQueryParams) ([]models.Review, *utils.Pagination, error)
		VerifyStore(ctx context.Context, storeID string) error
		GetFeaturedStores(ctx context.Context) ([]models.Store, error)
	}

	// StrategyService واجهة خدمة الاستراتيجيات
	StrategyService interface {
		CreateStrategy(ctx context.Context, req StrategyCreateRequest) (*models.Strategy, error)
		GetStrategyByID(ctx context.Context, strategyID string) (*models.Strategy, error)
		UpdateStrategy(ctx context.Context, strategyID string, req StrategyUpdateRequest) (*models.Strategy, error)
		DeleteStrategy(ctx context.Context, strategyID string) error
		ExecuteStrategy(ctx context.Context, strategyID string, params map[string]interface{}) (*StrategyExecutionResult, error)
		GetStrategyPerformance(ctx context.Context, strategyID string, timeframe string) (*StrategyPerformance, error)
		BacktestStrategy(ctx context.Context, req BacktestRequest) (*BacktestResult, error)
		GetStrategyTemplates(ctx context.Context) ([]StrategyTemplate, error)
	}

	// UploadService واجهة خدمة الرفع
	UploadService interface {
		UploadFile(ctx context.Context, req UploadRequest) (*UploadResult, error)
		DeleteFile(ctx context.Context, fileID string) error
		GetFile(ctx context.Context, fileID string) (*models.File, error)
		GetUserFiles(ctx context.Context, userID string, params FileQueryParams) ([]models.File, *utils.Pagination, error)
		GeneratePresignedURL(ctx context.Context, req PresignedURLRequest) (*PresignedURL, error)
		ValidateFile(ctx context.Context, fileInfo models.File) (*FileValidation, error)
		GetUploadQuota(ctx context.Context, userID string) (*UploadQuota, error)
	}

	// RepositoryService واجهة خدمة المستودع
	RepositoryService interface {
		Create(ctx context.Context, entity interface{}) error
		GetByID(ctx context.Context, id string, entity interface{}) error
		Update(ctx context.Context, entity interface{}) error
		Delete(ctx context.Context, id string, entity interface{}) error
		Find(ctx context.Context, query interface{}, results interface{}, pagination *utils.Pagination) error
		Count(ctx context.Context, query interface{}) (int64, error)
		Exists(ctx context.Context, query interface{}) (bool, error)
	}

	// CouponService واجهة خدمة الكوبونات
	CouponService interface {
		CreateCoupon(ctx context.Context, req CouponCreateRequest) (*models.Coupon, error)
		GetCouponByID(ctx context.Context, couponID string) (*models.Coupon, error)
		GetCouponByCode(ctx context.Context, code string) (*models.Coupon, error)
		UpdateCoupon(ctx context.Context, couponID string, req CouponUpdateRequest) (*models.Coupon, error)
		DeleteCoupon(ctx context.Context, couponID string) error
		ValidateCoupon(ctx context.Context, code string, amount float64) (*CouponValidation, error)
		GetCoupons(ctx context.Context, params CouponQueryParams) ([]models.Coupon, *utils.Pagination, error)
	}

	// WishlistService واجهة خدمة قائمة الرغبات
	WishlistService interface {
		AddToWishlist(ctx context.Context, userID string, serviceID string) error
		RemoveFromWishlist(ctx context.Context, userID string, serviceID string) error
		GetUserWishlist(ctx context.Context, userID string, params WishlistQueryParams) ([]models.Service, *utils.Pagination, error)
		IsInWishlist(ctx context.Context, userID string, serviceID string) (bool, error)
		GetWishlistCount(ctx context.Context, userID string) (int64, error)
	}

	// SubscriptionService واجهة خدمة الاشتراكات
	SubscriptionService interface {
		CreateSubscription(ctx context.Context, req SubscriptionCreateRequest) (*models.Subscription, error)
		GetSubscriptionByID(ctx context.Context, subscriptionID string) (*models.Subscription, error)
		GetUserSubscription(ctx context.Context, userID string) (*models.Subscription, error)
		CancelSubscription(ctx context.Context, subscriptionID string) error
		RenewSubscription(ctx context.Context, subscriptionID string) (*models.Subscription, error)
		GetSubscriptionPlans(ctx context.Context) ([]SubscriptionPlan, error)
	}
)

// ================================
// هياكل المعاملات المحدثة
// ================================

type (
	// AI Structures
	AIGenerateParams struct {
		Prompt      string
		MaxTokens   int
		Temperature float64
		Model       string
		UserID      string
	}

	AIImageParams struct {
		Prompt     string
		Size       string
		Style      string
		Quality    string
		UserID     string
	}

	AIChatMessage struct {
		Role    string
		Content string
	}

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

	// Cart Structures
	CartSummary struct {
		TotalItems    int     `json:"total_items"`
		Subtotal      float64 `json:"subtotal"`
		Tax           float64 `json:"tax"`
		Shipping      float64 `json:"shipping"`
		Discount      float64 `json:"discount"`
		Total         float64 `json:"total"`
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

	CategoryStats struct {
		TotalCategories int            `json:"total_categories"`
		ActiveCategories int           `json:"active_categories"`
		TopCategories   []CategoryStat `json:"top_categories"`
	}

	CategoryStat struct {
		CategoryID   string `json:"category_id"`
		CategoryName string `json:"category_name"`
		ServiceCount int    `json:"service_count"`
		TotalSales   int    `json:"total_sales"`
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

	PaymentInfo struct {
		PaymentMethodID string                 `json:"payment_method_id"`
		PaymentIntent   string                 `json:"payment_intent"`
		Metadata        map[string]interface{} `json:"metadata"`
	}

	// Payment Structures
	PaymentIntentRequest struct {
		Amount      float64                `json:"amount" binding:"required"`
		Currency    string                 `json:"currency" binding:"required"`
		Description string                 `json:"description"`
		Metadata    map[string]interface{} `json:"metadata"`
		UserID      string                 `json:"user_id" binding:"required"`
	}

	PaymentMethod struct {
		ID          string `json:"id"`
		Type        string `json:"type"`
		Last4       string `json:"last4"`
		Brand       string `json:"brand"`
		ExpMonth    int    `json:"exp_month"`
		ExpYear     int    `json:"exp_year"`
		IsDefault   bool   `json:"is_default"`
	}

	PaymentQueryParams struct {
		Page   int    `json:"page"`
		Limit  int    `json:"limit"`
		Status string `json:"status"`
		UserID string `json:"user_id"`
	}

	// Report Structures
	ReportParams struct {
		StartDate time.Time              `json:"start_date"`
		EndDate   time.Time              `json:"end_date"`
		Format    string                 `json:"format"`
		Filters   map[string]interface{} `json:"filters"`
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

	// Strategy Structures
	StrategyCreateRequest struct {
		Name        string                 `json:"name" binding:"required"`
		Description string                 `json:"description"`
		Type        string                 `json:"type" binding:"required"`
		Parameters  map[string]interface{} `json:"parameters"`
		Rules       []StrategyRule         `json:"rules"`
		CreatedBy   string                 `json:"created_by" binding:"required"`
	}

	StrategyUpdateRequest struct {
		Name        string                 `json:"name"`
		Description string                 `json:"description"`
		Parameters  map[string]interface{} `json:"parameters"`
		Active      *bool                  `json:"active"`
	}

	StrategyRule struct {
		ID        string      `json:"id"`
		Condition string      `json:"condition"`
		Action    string      `json:"action"`
		Value     interface{} `json:"value"`
		Priority  int         `json:"priority"`
	}

	BacktestRequest struct {
		StrategyID string                 `json:"strategy_id"`
		StartDate  time.Time              `json:"start_date"`
		EndDate    time.Time              `json:"end_date"`
		Parameters map[string]interface{} `json:"parameters"`
	}

	// Store Structures
	StoreCreateRequest struct {
		Name         string `json:"name" binding:"required"`
		Slug         string `json:"slug" binding:"required"`
		Description  string `json:"description"`
		ContactEmail string `json:"contact_email" binding:"required,email"`
		Phone        string `json:"phone,omitempty"`
		Address      string `json:"address,omitempty"`
		Banner       string `json:"banner,omitempty"`
		Logo         string `json:"logo,omitempty"`
		OwnerID      string `json:"owner_id" binding:"required"`
	}

	StoreUpdateRequest struct {
		Name         string `json:"name"`
		Description  string `json:"description"`
		ContactEmail string `json:"contact_email"`
		Phone        string `json:"phone"`
		Address      string `json:"address"`
		Banner       string `json:"banner"`
		Logo         string `json:"logo"`
		IsActive     *bool  `json:"is_active"`
	}

	// Coupon Structures
	CouponCreateRequest struct {
		Code          string    `json:"code" binding:"required"`
		Description   string    `json:"description"`
		DiscountType  string    `json:"discount_type" binding:"required"`
		DiscountValue float64   `json:"discount_value" binding:"required"`
		MinAmount     float64   `json:"min_amount"`
		MaxDiscount   float64   `json:"max_discount"`
		UsageLimit    int       `json:"usage_limit"`
		StartDate     time.Time `json:"start_date" binding:"required"`
		EndDate       time.Time `json:"end_date" binding:"required"`
	}

	CouponUpdateRequest struct {
		Description string    `json:"description"`
		UsageLimit  int       `json:"usage_limit"`
		StartDate   time.Time `json:"start_date"`
		EndDate     time.Time `json:"end_date"`
		IsActive    *bool     `json:"is_active"`
	}

	CouponQueryParams struct {
		Page   int    `json:"page"`
		Limit  int    `json:"limit"`
		Active *bool  `json:"active"`
	}

	// Wishlist Structures
	WishlistQueryParams struct {
		Page  int    `json:"page"`
		Limit int    `json:"limit"`
		SortBy string `json:"sort_by"`
	}

	// Subscription Structures
	SubscriptionCreateRequest struct {
		UserID   string    `json:"user_id" binding:"required"`
		PlanID   string    `json:"plan_id" binding:"required"`
		StartDate time.Time `json:"start_date" binding:"required"`
		EndDate   time.Time `json:"end_date" binding:"required"`
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

	// Content Structures
	ContentCreateRequest struct {
		Title   string   `json:"title" binding:"required"`
		Content string   `json:"content" binding:"required"`
		Type    string   `json:"type" binding:"required"`
		AuthorID string  `json:"author_id" binding:"required"`
		Slug    string   `json:"slug" binding:"required"`
		Image   string   `json:"image,omitempty"`
		Tags    []string `json:"tags"`
	}

	ContentUpdateRequest struct {
		Title   string   `json:"title"`
		Content string   `json:"content"`
		Image   string   `json:"image"`
		Tags    []string `json:"tags"`
	}

	ContentQueryParams struct {
		Page     int    `json:"page"`
		Limit    int    `json:"limit"`
		Type     string `json:"type"`
		AuthorID string `json:"author_id"`
		Published *bool  `json:"published"`
		SortBy   string `json:"sort_by"`
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

	BulkNotificationRequest struct {
		UserIDs []string               `json:"user_ids" binding:"required"`
		Title   string                 `json:"title" binding:"required"`
		Message string                 `json:"message" binding:"required"`
		Type    string                 `json:"type" binding:"required"`
		Data    map[string]interface{} `json:"data"`
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
// هياكل النتائج المحدثة
// ================================

type (
	AIGenerationResult struct {
		Text         string    `json:"text"`
		Tokens       int       `json:"tokens"`
		Model        string    `json:"model"`
		FinishReason string    `json:"finish_reason"`
		GeneratedAt  time.Time `json:"generated_at"`
	}

	SentimentAnalysis struct {
		Sentiment  string  `json:"sentiment"`
		Confidence float64 `json:"confidence"`
		Positive   float64 `json:"positive"`
		Negative   float64 `json:"negative"`
		Neutral    float64 `json:"neutral"`
	}

	ContentClassification struct {
		Category    string  `json:"category"`
		Confidence  float64 `json:"confidence"`
		Categories  []Class `json:"categories"`
	}

	Class struct {
		Name       string  `json:"name"`
		Confidence float64 `json:"confidence"`
	}

	KeywordExtraction struct {
		Keywords []Keyword `json:"keywords"`
	}

	Keyword struct {
		Word       string  `json:"word"`
		Score      float64 `json:"score"`
		Frequency  int     `json:"frequency"`
	}

	TextSummary struct {
		Summary    string `json:"summary"`
		OriginalLength int `json:"original_length"`
		SummaryLength  int `json:"summary_length"`
		CompressionRatio float64 `json:"compression_ratio"`
	}

	TranslationResult struct {
		Text         string `json:"text"`
		SourceLang   string `json:"source_lang"`
		TargetLang   string `json:"target_lang"`
		Translations []Translation `json:"translations"`
	}

	Translation struct {
		Text string `json:"text"`
		Confidence float64 `json:"confidence"`
	}

	AIImageResult struct {
		URL         string    `json:"url"`
		Width       int       `json:"width"`
		Height      int       `json:"height"`
		Format      string    `json:"format"`
		GeneratedAt time.Time `json:"generated_at"`
	}

	ChatCompletionResult struct {
		Message     AIChatMessage `json:"message"`
		Tokens      int           `json:"tokens"`
		Model       string        `json:"model"`
		GeneratedAt time.Time     `json:"generated_at"`
	}

	AuthResponse struct {
		User         *models.User `json:"user"`
		AccessToken  string       `json:"access_token"`
		RefreshToken string       `json:"refresh_token"`
		ExpiresAt    time.Time    `json:"expires_at"`
		Session      *models.Session `json:"session,omitempty"`
	}

	TokenClaims struct {
		UserID    string    `json:"user_id"`
		Email     string    `json:"email"`
		Role      string    `json:"role"`
		ExpiresAt time.Time `json:"expires_at"`
	}

	SessionInfo struct {
		SessionID string    `json:"session_id"`
		UserID    string    `json:"user_id"`
		ExpiresAt time.Time `json:"expires_at"`
		IPAddress string    `json:"ip_address"`
		UserAgent string    `json:"user_agent"`
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

	RefundResult struct {
		ID        string    `json:"id"`
		Status    string    `json:"status"`
		Amount    float64   `json:"amount"`
		Currency  string    `json:"currency"`
		RefundedAt time.Time `json:"refunded_at"`
	}

	PaymentValidation struct {
		IsValid bool   `json:"is_valid"`
		Message string `json:"message"`
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

	StrategyExecutionResult struct {
		StrategyID string                 `json:"strategy_id"`
		Success    bool                   `json:"success"`
		Output     map[string]interface{} `json:"output"`
		Metrics    map[string]float64     `json:"metrics"`
		ExecutedAt time.Time              `json:"executed_at"`
	}

	StrategyPerformance struct {
		StrategyID     string            `json:"strategy_id"`
		TotalExecutions int              `json:"total_executions"`
		SuccessRate    float64           `json:"success_rate"`
		AverageMetrics map[string]float64 `json:"average_metrics"`
		LastExecuted   time.Time         `json:"last_executed"`
	}

	BacktestResult struct {
		StrategyID  string            `json:"strategy_id"`
		Period      string            `json:"period"`
		TotalTrades int               `json:"total_trades"`
		WinRate     float64           `json:"win_rate"`
		ProfitLoss  float64           `json:"profit_loss"`
		Metrics     map[string]float64 `json:"metrics"`
		ExecutedAt  time.Time         `json:"executed_at"`
	}

	StrategyTemplate struct {
		ID          string                 `json:"id"`
		Name        string                 `json:"name"`
		Description string                 `json:"description"`
		Type        string                 `json:"type"`
		Parameters  map[string]interface{} `json:"parameters"`
		Rules       []StrategyRule         `json:"rules"`
	}

	StoreStats struct {
		TotalSales    int     `json:"total_sales"`
		TotalRevenue  float64 `json:"total_revenue"`
		AverageRating float64 `json:"average_rating"`
		TotalReviews  int     `json:"total_reviews"`
		ActiveServices int    `json:"active_services"`
	}

	OrderStats struct {
		TotalOrders    int     `json:"total_orders"`
		PendingOrders  int     `json:"pending_orders"`
		CompletedOrders int    `json:"completed_orders"`
		CanceledOrders int     `json:"canceled_orders"`
		TotalRevenue   float64 `json:"total_revenue"`
		AverageOrderValue float64 `json:"average_order_value"`
	}

	OrderPaymentResult struct {
		OrderID   string    `json:"order_id"`
		PaymentID string    `json:"payment_id"`
		Status    string    `json:"status"`
		PaidAt    time.Time `json:"paid_at"`
	}

	OrderTracking struct {
		OrderID       string           `json:"order_id"`
		Status        string           `json:"status"`
		TrackingNumber string          `json:"tracking_number"`
		Events        []TrackingEvent  `json:"events"`
		EstimatedDelivery time.Time    `json:"estimated_delivery"`
	}

	TrackingEvent struct {
		Timestamp time.Time `json:"timestamp"`
		Status    string    `json:"status"`
		Location  string    `json:"location"`
		Message   string    `json:"message"`
	}

	CouponValidation struct {
		IsValid      bool    `json:"is_valid"`
		DiscountType string  `json:"discount_type"`
		DiscountValue float64 `json:"discount_value"`
		Message      string  `json:"message"`
	}

	SubscriptionPlan struct {
		ID          string    `json:"id"`
		Name        string    `json:"name"`
		Description string    `json:"description"`
		Price       float64   `json:"price"`
		Duration    int       `json:"duration"`
		Features    []string  `json:"features"`
		IsActive    bool      `json:"is_active"`
	}

	UserStats struct {
		TotalOrders    int     `json:"total_orders"`
		TotalSpent     float64 `json:"total_spent"`
		JoinedDate     time.Time `json:"joined_date"`
		LastOrderDate  time.Time `json:"last_order_date"`
		WishlistCount  int     `json:"wishlist_count"`
	}

	UserAnalytics struct {
		UserID         string            `json:"user_id"`
		SessionCount   int               `json:"session_count"`
		PageViews      int               `json:"page_views"`
		ConversionRate float64           `json:"conversion_rate"`
		FavoriteCategories []string       `json:"favorite_categories"`
	}

	ServiceAnalytics struct {
		ServiceID      string            `json:"service_id"`
		Views          int               `json:"views"`
		Conversions    int               `json:"conversions"`
		Revenue        float64           `json:"revenue"`
		Rating         float64           `json:"rating"`
		PopularTimes   map[string]int    `json:"popular_times"`
	}

	PlatformAnalytics struct {
		TotalUsers     int               `json:"total_users"`
		ActiveUsers    int               `json:"active_users"`
		TotalOrders    int               `json:"total_orders"`
		TotalRevenue   float64           `json:"total_revenue"`
		PopularServices []ServiceStat    `json:"popular_services"`
	}

	ServiceStat struct {
		ServiceID   string  `json:"service_id"`
		ServiceName string  `json:"service_name"`
		Orders      int     `json:"orders"`
		Revenue     float64 `json:"revenue"`
	}

	DashboardStats struct {
		TotalUsers      int     `json:"total_users"`
		TotalServices   int     `json:"total_services"`
		TotalOrders     int     `json:"total_orders"`
		TotalRevenue    float64 `json:"total_revenue"`
		PendingOrders   int     `json:"pending_orders"`
		ActiveStores    int     `json:"active_stores"`
	}

	SalesReport struct {
		Period         string           `json:"period"`
		TotalSales     int              `json:"total_sales"`
		TotalRevenue   float64          `json:"total_revenue"`
		TopServices    []ServiceStat    `json:"top_services"`
		SalesByDay     map[string]int   `json:"sales_by_day"`
		RevenueByDay   map[string]float64 `json:"revenue_by_day"`
	}

	UserReport struct {
		Period         string          `json:"period"`
		NewUsers       int             `json:"new_users"`
		ActiveUsers    int             `json:"active_users"`
		UserGrowth     float64         `json:"user_growth"`
		UserDemographics map[string]int `json:"user_demographics"`
	}

	ServiceReport struct {
		Period          string           `json:"period"`
		TotalServices   int              `json:"total_services"`
		NewServices     int              `json:"new_services"`
		PopularServices []ServiceStat    `json:"popular_services"`
		ServiceCategories map[string]int `json:"service_categories"`
	}

	FinancialReport struct {
		Period         string           `json:"period"`
		TotalRevenue   float64          `json:"total_revenue"`
		TotalExpenses  float64          `json:"total_expenses"`
		NetProfit      float64          `json:"net_profit"`
		RevenueSources map[string]float64 `json:"revenue_sources"`
		ExpenseBreakdown map[string]float64 `json:"expense_breakdown"`
	}

	SystemReport struct {
		Period          string           `json:"period"`
		SystemUptime    float64          `json:"system_uptime"`
		ErrorCount      int              `json:"error_count"`
		ActiveSessions  int              `json:"active_sessions"`
		PerformanceMetrics map[string]float64 `json:"performance_metrics"`
	}

	ReportTemplate struct {
		ID          string    `json:"id"`
		Name        string    `json:"name"`
		Description string    `json:"description"`
		Type        string    `json:"type"`
		Format      string    `json:"format"`
		Fields      []string  `json:"fields"`
	}

	ScheduleReportRequest struct {
		TemplateID  string    `json:"template_id" binding:"required"`
		Recipients  []string  `json:"recipients" binding:"required"`
		Schedule    string    `json:"schedule" binding:"required"`
		Parameters  ReportParams `json:"parameters"`
	}

	ScheduledReport struct {
		ID          string      `json:"id"`
		TemplateID  string      `json:"template_id"`
		Status      string      `json:"status"`
		NextRun     time.Time   `json:"next_run"`
		LastRun     time.Time   `json:"last_run"`
		CreatedAt   time.Time   `json:"created_at"`
	}

	ScheduledReportQuery struct {
		Page   int    `json:"page"`
		Limit  int    `json:"limit"`
		Status string `json:"status"`
	}
)

// ================================
// التطبيقات الفعلية المحدثة
// ================================

type (
	aiServiceImpl struct {
		db *gorm.DB
	}

	authServiceImpl struct {
		db *gorm.DB
	}

	cartServiceImpl struct {
		db *gorm.DB
	}

	categoryServiceImpl struct {
		db *gorm.DB
	}

	orderServiceImpl struct {
		db *gorm.DB
	}

	paymentServiceImpl struct {
		db *gorm.DB
	}

	reportServiceImpl struct {
		db *gorm.DB
	}

	storeServiceImpl struct {
		db *gorm.DB
	}

	strategyServiceImpl struct {
		db *gorm.DB
	}

	uploadServiceImpl struct {
		db *gorm.DB
	}

	repositoryServiceImpl struct {
		db *gorm.DB
	}

	couponServiceImpl struct {
		db *gorm.DB
	}

	wishlistServiceImpl struct {
		db *gorm.DB
	}

	subscriptionServiceImpl struct {
		db *gorm.DB
	}

	analyticsServiceImpl struct {
		db *gorm.DB
	}

	adminServiceImpl struct {
		db *gorm.DB
	}

	contentServiceImpl struct {
		db *gorm.DB
	}

	notificationServiceImpl struct {
		db *gorm.DB
	}

	userServiceImpl struct {
		db *gorm.DB
	}

	serviceServiceImpl struct {
		db *gorm.DB
	}

	cacheServiceImpl struct {
		// implementation details
	}
)

// ================================
// دوال الإنشاء المحدثة
// ================================

func NewAIService(db *gorm.DB) AIService {
	return &aiServiceImpl{db: db}
}

func NewAuthService(db *gorm.DB) AuthService {
	return &authServiceImpl{db: db}
}

func NewCartService(db *gorm.DB) CartService {
	return &cartServiceImpl{db: db}
}

func NewCategoryService(db *gorm.DB) CategoryService {
	return &categoryServiceImpl{db: db}
}

func NewOrderService(db *gorm.DB) OrderService {
	return &orderServiceImpl{db: db}
}

func NewPaymentService(db *gorm.DB) PaymentService {
	return &paymentServiceImpl{db: db}
}

func NewReportService(db *gorm.DB) ReportService {
	return &reportServiceImpl{db: db}
}

func NewStoreService(db *gorm.DB) StoreService {
	return &storeServiceImpl{db: db}
}

func NewStrategyService(db *gorm.DB) StrategyService {
	return &strategyServiceImpl{db: db}
}

func NewUploadService(db *gorm.DB) UploadService {
	return &uploadServiceImpl{db: db}
}

func NewRepositoryService(db *gorm.DB) RepositoryService {
	return &repositoryServiceImpl{db: db}
}

func NewCouponService(db *gorm.DB) CouponService {
	return &couponServiceImpl{db: db}
}

func NewWishlistService(db *gorm.DB) WishlistService {
	return &wishlistServiceImpl{db: db}
}

func NewSubscriptionService(db *gorm.DB) SubscriptionService {
	return &subscriptionServiceImpl{db: db}
}

func NewAnalyticsService(db *gorm.DB) AnalyticsService {
	return &analyticsServiceImpl{db: db}
}

func NewAdminService(db *gorm.DB) AdminService {
	return &adminServiceImpl{db: db}
}

func NewContentService(db *gorm.DB) ContentService {
	return &contentServiceImpl{db: db}
}

func NewNotificationService(db *gorm.DB) NotificationService {
	return &notificationServiceImpl{db: db}
}

func NewUserService(db *gorm.DB) UserService {
	return &userServiceImpl{db: db}
}

func NewServiceService(db *gorm.DB) ServiceService {
	return &serviceServiceImpl{db: db}
}

func NewCacheService(/* parameters */) CacheService {
	return &cacheServiceImpl{}
}

// ================================
// Service Container المحدث
// ================================

type ServiceContainer struct {
	Analytics     AnalyticsService
	Admin         AdminService
	Content       ContentService
	Notification  NotificationService
	User          UserService
	Service       ServiceService
	AI            AIService
	Auth          AuthService
	Cart          CartService
	Category      CategoryService
	Order         OrderService
	Payment       PaymentService
	Report        ReportService
	Store         StoreService
	Strategy      StrategyService
	Upload        UploadService
	Repository    RepositoryService
	Cache         CacheService
	Coupon        CouponService
	Wishlist      WishlistService
	Subscription  SubscriptionService
}

func NewServiceContainer(db *gorm.DB) *ServiceContainer {
	return &ServiceContainer{
		Analytics:     NewAnalyticsService(db),
		Admin:         NewAdminService(db),
		Content:       NewContentService(db),
		Notification:  NewNotificationService(db),
		User:          NewUserService(db),
		Service:       NewServiceService(db),
		AI:            NewAIService(db),
		Auth:          NewAuthService(db),
		Cart:          NewCartService(db),
		Category:      NewCategoryService(db),
		Order:         NewOrderService(db),
		Payment:       NewPaymentService(db),
		Report:        NewReportService(db),
		Store:         NewStoreService(db),
		Strategy:      NewStrategyService(db),
		Upload:        NewUploadService(db),
		Repository:    NewRepositoryService(db),
		Cache:         NewCacheService(),
		Coupon:        NewCouponService(db),
		Wishlist:      NewWishlistService(db),
		Subscription:  NewSubscriptionService(db),
	}
}

// ================================
// تطبيقات أساسية للخدمات المحدثة
// ================================

func (s *aiServiceImpl) GenerateText(ctx context.Context, params AIGenerateParams) (*AIGenerationResult, error) {
	return &AIGenerationResult{
		Text:        "نص تم إنشاؤه بواسطة الذكاء الاصطناعي: " + params.Prompt,
		Tokens:      50,
		Model:       params.Model,
		GeneratedAt: time.Now(),
	}, nil
}

func (s *authServiceImpl) Register(ctx context.Context, req AuthRegisterRequest) (*AuthResponse, error) {
	user := &models.User{
		ID:           fmt.Sprintf("user_%d", time.Now().Unix()),
		Email:        req.Email,
		Username:     req.Username,
		Password:     "hashed_password", // يجب تشفير كلمة المرور
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Phone:        req.Phone,
		Role:         "user",
		Status:       "active",
		EmailVerified: false,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	session := &models.Session{
		ID:        fmt.Sprintf("session_%d", time.Now().Unix()),
		UserID:    user.ID,
		Token:     "session_token_" + fmt.Sprintf("%d", time.Now().Unix()),
		ExpiresAt: time.Now().Add(24 * time.Hour),
		CreatedAt: time.Now(),
	}

	return &AuthResponse{
		User:         user,
		AccessToken:  "access_token_" + fmt.Sprintf("%d", time.Now().Unix()),
		RefreshToken: "refresh_token_" + fmt.Sprintf("%d", time.Now().Unix()),
		ExpiresAt:    time.Now().Add(24 * time.Hour),
		Session:      session,
	}, nil
}

func (s *cartServiceImpl) GetCart(ctx context.Context, userID string) (*models.Cart, error) {
	return &models.Cart{
		ID:          "cart_" + userID,
		UserID:      userID,
		Items:       []models.CartItem{},
		TotalAmount: 0,
		Discount:    0,
		Tax:         0,
		Shipping:    0,
		FinalAmount: 0,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}, nil
}

func (s *orderServiceImpl) CreateOrder(ctx context.Context, req OrderCreateRequest) (*models.Order, error) {
	// حساب المبلغ الإجمالي
	var totalAmount float64
	for _, item := range req.Items {
		totalAmount += item.Price * float64(item.Quantity)
	}

	return &models.Order{
		ID:           fmt.Sprintf("order_%d", time.Now().Unix()),
		UserID:       "user_id_from_context", // سيتم تعيينه من السياق
		SellerID:     "seller_id_from_items", // سيتم استخلاصه من العناصر
		Items:        req.Items,
		Status:       "pending",
		TotalAmount:  totalAmount,
		Discount:     0,
		Tax:          0,
		Shipping:     0,
		FinalAmount:  totalAmount,
		PaymentStatus: "pending",
		PaymentMethod: req.PaymentMethod,
		ShippingInfo: req.ShippingInfo,
		CustomerNotes: req.CustomerNotes,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}, nil
}

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

// ... تطبيقات مماثلة للخدمات الأخرى

// تطبيقات أساسية للخدمات الجديدة
func (s *couponServiceImpl) CreateCoupon(ctx context.Context, req CouponCreateRequest) (*models.Coupon, error) {
	return &models.Coupon{
		ID:           fmt.Sprintf("coupon_%d", time.Now().Unix()),
		Code:         req.Code,
		Description:  req.Description,
		DiscountType: req.DiscountType,
		DiscountValue: req.DiscountValue,
		MinAmount:    req.MinAmount,
		MaxDiscount:  req.MaxDiscount,
		UsageLimit:   req.UsageLimit,
		UsedCount:    0,
		StartDate:    req.StartDate,
		EndDate:      req.EndDate,
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}, nil
}

func (s *wishlistServiceImpl) AddToWishlist(ctx context.Context, userID string, serviceID string) error {
	// تنفيذ إضافة إلى قائمة الرغبات
	return nil
}

func (s *subscriptionServiceImpl) CreateSubscription(ctx context.Context, req SubscriptionCreateRequest) (*models.Subscription, error) {
	return &models.Subscription{
		ID:          fmt.Sprintf("sub_%d", time.Now().Unix()),
		UserID:      req.UserID,
		PlanID:      req.PlanID,
		Status:      "active",
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		RenewalDate: req.EndDate.AddDate(0, 1, 0), // تجديد بعد شهر
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}, nil
}