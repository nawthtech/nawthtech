package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/nawthtech/nawthtech/backend/internal/models"
)

// ================================
// هياكل المعاملات المحدثة
// ================================

type (
	ReviewQueryParams struct {
		Page   int    `json:"page"`
		Limit  int    `json:"limit"`
		Rating int    `json:"rating"`
		SortBy string `json:"sort_by"`
	}

	AuthRegisterRequest struct {
		Email     string `json:"email"`
		Username  string `json:"username"`
		Password  string `json:"password"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Phone     string `json:"phone"`
	}

	AuthLoginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	AuthResponse struct {
		User         *models.User `json:"user"`
		AccessToken  string       `json:"access_token"`
		RefreshToken string       `json:"refresh_token"`
		ExpiresAt    time.Time    `json:"expires_at"`
	}

	TokenClaims struct {
		UserID string `json:"user_id"`
		Email  string `json:"email"`
		Role   string `json:"role"`
		Exp    int64  `json:"exp"`
	}

	ChangePasswordRequest struct {
		CurrentPassword string `json:"current_password"`
		NewPassword     string `json:"new_password"`
	}

	UserUpdateRequest struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Phone     string `json:"phone"`
		Avatar    string `json:"avatar"`
	}

	UserQueryParams struct {
		Page  int    `json:"page"`
		Limit int    `json:"limit"`
		Role  string `json:"role"`
	}

	UserStats struct {
		TotalOrders   int     `json:"total_orders"`
		TotalSpent    float64 `json:"total_spent"`
		ActiveSince   string  `json:"active_since"`
		ServicesCount int     `json:"services_count"`
	}

	ServiceCreateRequest struct {
		Title       string   `json:"title"`
		Description string   `json:"description"`
		Price       float64  `json:"price"`
		Duration    int      `json:"duration"`
		CategoryID  string   `json:"category_id"`
		ProviderID  string   `json:"provider_id"`
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
		IsActive    bool     `json:"is_active"`
		IsFeatured  bool     `json:"is_featured"`
	}

	ServiceQueryParams struct {
		Page       int     `json:"page"`
		Limit      int     `json:"limit"`
		CategoryID string  `json:"category_id"`
		ProviderID string  `json:"provider_id"`
		MinPrice   float64 `json:"min_price"`
		MaxPrice   float64 `json:"max_price"`
		IsActive   bool    `json:"is_active"`
		IsFeatured bool    `json:"is_featured"`
	}
 ServiceContainer struct {
	Auth         AuthService
	User         UserService
	Service      ServiceService
	Category     CategoryService
	Order        OrderService
	Payment      PaymentService  // تأكد من وجود هذا
	Upload       UploadService
	Notification NotificationService
	Admin        AdminService
	Cache        CacheService
}

func NewServiceContainer(d1db *sql.DB) *ServiceContainer {
	return &ServiceContainer{
		Auth:         NewAuthService(d1db),
		User:         NewUserService(d1db),
		Service:      NewServiceService(d1db),
		Category:     NewCategoryService(d1db),
		Order:        NewOrderService(d1db),
		Payment:      NewPaymentService(d1db),  // أضف هذا السطر
		Upload:       NewUploadService(d1db),
		Notification: NewNotificationService(d1db),
		Admin:        NewAdminService(d1db),
		Cache:        NewCacheService(),
	}

	CategoryCreateRequest struct {
		Name  string `json:"name"`
		Slug  string `json:"slug"`
		Image string `json:"image"`
	}

	CategoryUpdateRequest struct {
		Name     string `json:"name"`
		Slug     string `json:"slug"`
		Image    string `json:"image"`
		IsActive bool   `json:"is_active"`
	}

	CategoryQueryParams struct {
		Page     int  `json:"page"`
		Limit    int  `json:"limit"`
		IsActive bool `json:"is_active"`
	}

	CategoryNode struct {
		Category  *models.Category `json:"category"`
		Children  []CategoryNode   `json:"children"`
		Services  int              `json:"services_count"`
	}

	OrderCreateRequest struct {
		UserID    string  `json:"user_id"`
		ServiceID string  `json:"service_id"`
		Amount    float64 `json:"amount"`
		Notes     string  `json:"notes"`
	}

	OrderQueryParams struct {
		Page   int    `json:"page"`
		Limit  int    `json:"limit"`
		Status string `json:"status"`
		UserID string `json:"user_id"`
	}

	OrderStats struct {
		TotalOrders   int     `json:"total_orders"`
		PendingOrders int     `json:"pending_orders"`
		Completed     int     `json:"completed_orders"`
		Cancelled     int     `json:"cancelled_orders"`
		TotalRevenue  float64 `json:"total_revenue"`
		AvgOrderValue float64 `json:"avg_order_value"`
	}
)
// ================================
// هياكل طلبات الدفع
// ================================

type PaymentCreateRequest struct {
	OrderID string  `json:"order_id"`
	Amount  float64 `json:"amount"`
}

type PaymentIntentRequest struct {
	OrderID   string  `json:"order_id"`
	Amount    float64 `json:"amount"`
	Currency  string  `json:"currency"`
	Customer  string  `json:"customer,omitempty"`
	ReturnURL string  `json:"return_url,omitempty"`
}

type PaymentIntent struct {
	ID           string                 `json:"id"`
	ClientSecret string                 `json:"client_secret"`
	Status       string                 `json:"status"`
	Amount       float64                `json:"amount"`
	Currency     string                 `json:"currency"`
	CreatedAt    time.Time              `json:"created_at"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
}

type PaymentResult struct {
	Success    bool        `json:"success"`
	Message    string      `json:"message"`
	PaymentID  string      `json:"payment_id,omitempty"`
	OrderID    string      `json:"order_id,omitempty"`
	Amount     float64     `json:"amount,omitempty"`
	Currency   string      `json:"currency,omitempty"`
	Status     string      `json:"status,omitempty"`
	Data       interface{} `json:"data,omitempty"`
	Timestamp  time.Time   `json:"timestamp"`
}

type PaymentQueryParams struct {
	Page     int       `json:"page"`
	Limit    int       `json:"limit"`
	Status   string    `json:"status"`
	UserID   string    `json:"user_id"`
	OrderID  string    `json:"order_id"`
	FromDate time.Time `json:"from_date,omitempty"`
	ToDate   time.Time `json:"to_date,omitempty"`
}

type PaymentValidation struct {
	Valid    bool                   `json:"valid"`
	Reason   string                 `json:"reason,omitempty"`
	Payment  *models.Payment        `json:"payment,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// ================================
// هياكل أخرى متعلقة بالدفع
// ================================

type PaymentMethod struct {
	ID         string    `json:"id"`
	Type       string    `json:"type"` // card, bank_account, etc.
	LastFour   string    `json:"last_four,omitempty"`
	ExpMonth   int       `json:"exp_month,omitempty"`
	ExpYear    int       `json:"exp_year,omitempty"`
	Brand      string    `json:"brand,omitempty"`
	IsDefault  bool      `json:"is_default"`
	CreatedAt  time.Time `json:"created_at"`
}

type RefundRequest struct {
	PaymentID string  `json:"payment_id"`
	Amount    float64 `json:"amount"`
	Reason    string  `json:"reason,omitempty"`
}

type RefundResult struct {
	ID        string    `json:"id"`
	PaymentID string    `json:"payment_id"`
	Amount    float64   `json:"amount"`
	Status    string    `json:"status"`
	Reason    string    `json:"reason,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// ================================
// تطبيق PaymentService مع D1
// ================================

type paymentServiceImpl struct {
	db *sql.DB
}

func (s *paymentServiceImpl) CreatePaymentIntent(ctx context.Context, req PaymentIntentRequest) (*PaymentIntent, error) {
	paymentID := fmt.Sprintf("pi_%d", time.Now().UnixNano())
	
	paymentIntent := &PaymentIntent{
		ID:           paymentID,
		ClientSecret: fmt.Sprintf("secret_%s", paymentID),
		Status:       "requires_payment_method",
		Amount:       req.Amount,
		Currency:     req.Currency,
		CreatedAt:    time.Now(),
		Metadata: map[string]interface{}{
			"order_id": req.OrderID,
		},
	}

	// حفظ في قاعدة البيانات
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO payment_intents (id, order_id, amount, currency, status, client_secret, metadata, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		paymentIntent.ID, req.OrderID, req.Amount, req.Currency, paymentIntent.Status,
		paymentIntent.ClientSecret, "{}", paymentIntent.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create payment intent: %w", err)
	}

	return paymentIntent, nil
}

func (s *paymentServiceImpl) ConfirmPayment(ctx context.Context, paymentID string, confirmationData map[string]interface{}) (*PaymentResult, error) {
	// تحديث حالة الدفع
	_, err := s.db.ExecContext(ctx, `
		UPDATE payment_intents SET status = ?, updated_at = ? WHERE id = ?`,
		"succeeded", time.Now(), paymentID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to confirm payment: %w", err)
	}

	// الحصول على تفاصيل الدفع
	row := s.db.QueryRowContext(ctx, `
		SELECT order_id, amount, currency FROM payment_intents WHERE id = ?`, paymentID)
	
	var orderID string
	var amount float64
	var currency string
	err = row.Scan(&orderID, &amount, &currency)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment details: %w", err)
	}

	return &PaymentResult{
		Success:   true,
		Message:   "Payment confirmed successfully",
		PaymentID: paymentID,
		OrderID:   orderID,
		Amount:    amount,
		Currency:  currency,
		Status:    "succeeded",
		Timestamp: time.Now(),
	}, nil
}

func (s *paymentServiceImpl) GetPaymentHistory(ctx context.Context, userID string, params PaymentQueryParams) ([]models.Payment, error) {
	query := `
		SELECT p.id, p.order_id, p.amount, p.currency, p.status, p.payment_method, 
		       p.transaction_id, p.created_at, p.updated_at
		FROM payments p
		INNER JOIN orders o ON p.order_id = o.id
		WHERE o.user_id = ?`
	
	args := []interface{}{userID}
	
	// تطبيق الفلترة
	if params.Status != "" {
		query += " AND p.status = ?"
		args = append(args, params.Status)
	}
	
	if params.OrderID != "" {
		query += " AND p.order_id = ?"
		args = append(args, params.OrderID)
	}
	
	if !params.FromDate.IsZero() {
		query += " AND p.created_at >= ?"
		args = append(args, params.FromDate)
	}
	
	if !params.ToDate.IsZero() {
		query += " AND p.created_at <= ?"
		args = append(args, params.ToDate)
	}
	
	// الترتيب والتصفح
	query += " ORDER BY p.created_at DESC LIMIT ? OFFSET ?"
	args = append(args, params.Limit, (params.Page-1)*params.Limit)
	
	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query payment history: %w", err)
	}
	defer rows.Close()
	
	var payments []models.Payment
	for rows.Next() {
		var payment models.Payment
		err := rows.Scan(
			&payment.ID, &payment.OrderID, &payment.Amount, &payment.Currency,
			&payment.Status, &payment.PaymentMethod, &payment.TransactionID,
			&payment.CreatedAt, &payment.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan payment: %w", err)
		}
		payments = append(payments, payment)
	}
	
	return payments, nil
}

func (s *paymentServiceImpl) ValidatePayment(ctx context.Context, paymentData map[string]interface{}) (*PaymentValidation, error) {
	// التحقق الأساسي
	amount, ok := paymentData["amount"].(float64)
	if !ok || amount <= 0 {
		return &PaymentValidation{
			Valid:  false,
			Reason: "Invalid or missing amount",
		}, nil
	}
	
	currency, ok := paymentData["currency"].(string)
	if !ok || currency == "" {
		return &PaymentValidation{
			Valid:  false,
			Reason: "Invalid or missing currency",
		}, nil
	}
	
	// يمكن إضافة المزيد من التحققات هنا
	// مثل التحقق من رقم البطاقة، تاريخ الانتهاء، إلخ.
	
	return &PaymentValidation{
		Valid: true,
		Metadata: map[string]interface{}{
			"validated_at": time.Now(),
			"amount":       amount,
			"currency":     currency,
		},
	}, nil
}

// ================================
// دالة إنشاء PaymentService
// ================================

func NewPaymentService(db *sql.DB) PaymentService {
	return &paymentServiceImpl{db: db}
}

// ================================
// الواجهات الرئيسية
// ================================

type (
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

	UserService interface {
		GetProfile(ctx context.Context, userID string) (*models.User, error)
		UpdateProfile(ctx context.Context, userID string, req UserUpdateRequest) (*models.User, error)
		UpdateAvatar(ctx context.Context, userID string, avatarURL string) error
		DeleteAccount(ctx context.Context, userID string) error
		SearchUsers(ctx context.Context, query string, params UserQueryParams) ([]models.User, error)
		GetUserStats(ctx context.Context, userID string) (*UserStats, error)
	}

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

	CategoryService interface {
		GetCategories(ctx context.Context, params CategoryQueryParams) ([]models.Category, error)
		GetCategoryByID(ctx context.Context, categoryID string) (*models.Category, error)
		CreateCategory(ctx context.Context, req CategoryCreateRequest) (*models.Category, error)
		UpdateCategory(ctx context.Context, categoryID string, req CategoryUpdateRequest) (*models.Category, error)
		DeleteCategory(ctx context.Context, categoryID string) error
		GetCategoryTree(ctx context.Context) ([]CategoryNode, error)
	}

	OrderService interface {
		CreateOrder(ctx context.Context, req OrderCreateRequest) (*models.Order, error)
		GetOrderByID(ctx context.Context, orderID string) (*models.Order, error)
		GetUserOrders(ctx context.Context, userID string, params OrderQueryParams) ([]models.Order, error)
		UpdateOrderStatus(ctx context.Context, orderID string, status string, notes string) (*models.Order, error)
		CancelOrder(ctx context.Context, orderID string, reason string) (*models.Order, error)
		GetOrderStats(ctx context.Context, timeframe string) (*OrderStats, error)
	}

	PaymentService interface {
		CreatePaymentIntent(ctx context.Context, req PaymentIntentRequest) (*PaymentIntent, error)
		ConfirmPayment(ctx context.Context, paymentID string, confirmationData map[string]interface{}) (*PaymentResult, error)
		GetPaymentHistory(ctx context.Context, userID string, params PaymentQueryParams) ([]models.Payment, error)
		ValidatePayment(ctx context.Context, paymentData map[string]interface{}) (*PaymentValidation, error)
	}

	UploadService interface {
		UploadFile(ctx context.Context, req UploadRequest) (*UploadResult, error)
		DeleteFile(ctx context.Context, fileID string) error
		GetFile(ctx context.Context, fileID string) (*models.File, error)
		GetUserFiles(ctx context.Context, userID string, params FileQueryParams) ([]models.File, error)
		GeneratePresignedURL(ctx context.Context, req PresignedURLRequest) (*PresignedURL, error)
		ValidateFile(ctx context.Context, fileInfo models.File) (*FileValidation, error)
		GetUploadQuota(ctx context.Context, userID string) (*UploadQuota, error)
	}

	NotificationService interface {
		CreateNotification(ctx context.Context, req NotificationCreateRequest) (*models.Notification, error)
		GetUserNotifications(ctx context.Context, userID string, params NotificationQueryParams) ([]models.Notification, error)
		MarkAsRead(ctx context.Context, notificationID string) error
		MarkAllAsRead(ctx context.Context, userID string) error
		DeleteNotification(ctx context.Context, notificationID string) error
		GetUnreadCount(ctx context.Context, userID string) (int64, error)
	}

	AdminService interface {
		GetDashboardStats(ctx context.Context) (*DashboardStats, error)
		GetUsers(ctx context.Context, params UserQueryParams) ([]models.User, error)
		GetSystemLogs(ctx context.Context, params SystemLogQuery) ([]models.SystemLog, error)
		UpdateSystemSettings(ctx context.Context, settings []models.Setting) error
		BanUser(ctx context.Context, userID string, reason string) error
		UnbanUser(ctx context.Context, userID string) error
	}

	CacheService interface {
		Get(key string) (interface{}, error)
		Set(key string, value interface{}, expiration time.Duration) error
		Delete(key string) error
		Exists(key string) (bool, error)
		Flush() error
	}
)

// ================================
// تطبيقات D1 SQL
// ================================

type (
	authServiceImpl struct {
		db *sql.DB
	}

	userServiceImpl struct {
		db *sql.DB
	}

	serviceServiceImpl struct {
		db *sql.DB
	}

	categoryServiceImpl struct {
		db *sql.DB
	}

	orderServiceImpl struct {
		db *sql.DB
	}

	cacheServiceImpl struct {
		store map[string]interface{}
	}
)

// ================================
// تطبيقات AuthService باستخدام D1
// ================================

func (s *authServiceImpl) Register(ctx context.Context, req AuthRegisterRequest) (*AuthResponse, error) {
	userID := fmt.Sprintf("user_%d", time.Now().UnixNano())

	_, err := s.db.ExecContext(ctx,
		`INSERT INTO users (id, email, username, password, first_name, last_name, phone, role, status, email_verified, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		userID, req.Email, req.Username, "hashed_password", req.FirstName, req.LastName, req.Phone, "user", "active", false, time.Now(), time.Now(),
	)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		ID:            userID,
		Email:         req.Email,
		Username:      req.Username,
		Role:          "user",
		Status:        "active",
		EmailVerified: false,
	}

	return &AuthResponse{
		User:         user,
		AccessToken:  "access_" + userID,
		RefreshToken: "refresh_" + userID,
		ExpiresAt:    time.Now().Add(24 * time.Hour),
	}, nil
}

func (s *authServiceImpl) Login(ctx context.Context, req AuthLoginRequest) (*AuthResponse, error) {
	var user models.User
	row := s.db.QueryRowContext(ctx, "SELECT id, email, username, role, status, email_verified FROM users WHERE email = ?", req.Email)
	err := row.Scan(&user.ID, &user.Email, &user.Username, &user.Role, &user.Status, &user.EmailVerified)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("المستخدم غير موجود")
		}
		return nil, err
	}

	return &AuthResponse{
		User:         &user,
		AccessToken:  "access_" + user.ID,
		RefreshToken: "refresh_" + user.ID,
		ExpiresAt:    time.Now().Add(24 * time.Hour),
	}, nil
}

func (s *authServiceImpl) Logout(ctx context.Context, token string) error {
	// تنفيذ تسجيل الخروج
	return nil
}

func (s *authServiceImpl) RefreshToken(ctx context.Context, refreshToken string) (*AuthResponse, error) {
	// تنفيذ تحديث الرمز
	return nil, fmt.Errorf("not implemented")
}

func (s *authServiceImpl) VerifyToken(ctx context.Context, token string) (*TokenClaims, error) {
	// تنفيذ التحقق من الرمز
	return nil, fmt.Errorf("not implemented")
}

func (s *authServiceImpl) ForgotPassword(ctx context.Context, email string) error {
	// تنفيذ نسيان كلمة المرور
	return fmt.Errorf("not implemented")
}

func (s *authServiceImpl) ResetPassword(ctx context.Context, token string, newPassword string) error {
	// تنفيذ إعادة تعيين كلمة المرور
	return fmt.Errorf("not implemented")
}

func (s *authServiceImpl) ChangePassword(ctx context.Context, userID string, req ChangePasswordRequest) error {
	// تنفيذ تغيير كلمة المرور
	return fmt.Errorf("not implemented")
}

// ================================
// تطبيقات UserService
// ================================

func (s *userServiceImpl) GetProfile(ctx context.Context, userID string) (*models.User, error) {
	row := s.db.QueryRowContext(ctx, "SELECT id,email,username,first_name,last_name,phone,avatar,role,status,email_verified,created_at,updated_at FROM users WHERE id = ?", userID)

	user := &models.User{}
	err := row.Scan(&user.ID, &user.Email, &user.Username, &user.FirstName, &user.LastName,
		&user.Phone, &user.Avatar, &user.Role, &user.Status, &user.EmailVerified,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("المستخدم غير موجود")
		}
		return nil, err
	}

	return user, nil
}

func (s *userServiceImpl) UpdateProfile(ctx context.Context, userID string, req UserUpdateRequest) (*models.User, error) {
	_, err := s.db.ExecContext(ctx,
		`UPDATE users SET first_name=?, last_name=?, phone=?, avatar=?, updated_at=?
		 WHERE id=?`,
		req.FirstName, req.LastName, req.Phone, req.Avatar, time.Now(), userID,
	)
	if err != nil {
		return nil, err
	}
	return s.GetProfile(ctx, userID)
}

func (s *userServiceImpl) UpdateAvatar(ctx context.Context, userID string, avatarURL string) error {
	_, err := s.db.ExecContext(ctx,
		"UPDATE users SET avatar=?, updated_at=? WHERE id=?",
		avatarURL, time.Now(), userID,
	)
	return err
}

func (s *userServiceImpl) DeleteAccount(ctx context.Context, userID string) error {
	_, err := s.db.ExecContext(ctx, "DELETE FROM users WHERE id=?", userID)
	return err
}

func (s *userServiceImpl) SearchUsers(ctx context.Context, query string, params UserQueryParams) ([]models.User, error) {
	// تنفيذ البحث
	return []models.User{}, nil
}

func (s *userServiceImpl) GetUserStats(ctx context.Context, userID string) (*UserStats, error) {
	// تنفيذ الحصول على الإحصائيات
	return &UserStats{}, nil
}

// ================================
// تطبيقات ServiceService
// ================================

func (s *serviceServiceImpl) CreateService(ctx context.Context, req ServiceCreateRequest) (*models.Service, error) {
	serviceID := fmt.Sprintf("service_%d", time.Now().UnixNano())
	service := &models.Service{
		ID:          serviceID,
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
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	_, err := s.db.ExecContext(ctx,
		`INSERT INTO services (id,title,description,price,duration,category_id,provider_id,images,tags,is_active,is_featured,rating,created_at,updated_at)
		 VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?)`,
		service.ID, service.Title, service.Description, service.Price, service.Duration, service.CategoryID,
		service.ProviderID, serializeStrings(service.Images), serializeStrings(service.Tags),
		service.IsActive, service.IsFeatured, service.Rating, service.CreatedAt, service.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return service, nil
}

func (s *serviceServiceImpl) GetServiceByID(ctx context.Context, serviceID string) (*models.Service, error) {
	// تنفيذ الحصول على الخدمة
	return &models.Service{}, nil
}

func (s *serviceServiceImpl) UpdateService(ctx context.Context, serviceID string, req ServiceUpdateRequest) (*models.Service, error) {
	// تنفيذ تحديث الخدمة
	return &models.Service{}, nil
}

func (s *serviceServiceImpl) DeleteService(ctx context.Context, serviceID string) error {
	// تنفيذ حذف الخدمة
	return nil
}

func (s *serviceServiceImpl) GetServices(ctx context.Context, params ServiceQueryParams) ([]models.Service, error) {
	// تنفيذ الحصول على الخدمات
	return []models.Service{}, nil
}

func (s *serviceServiceImpl) SearchServices(ctx context.Context, query string, params ServiceQueryParams) ([]models.Service, error) {
	// تنفيذ البحث في الخدمات
	return []models.Service{}, nil
}

func (s *serviceServiceImpl) GetFeaturedServices(ctx context.Context) ([]models.Service, error) {
	// تنفيذ الحصول على الخدمات المميزة
	return []models.Service{}, nil
}

func (s *serviceServiceImpl) GetSimilarServices(ctx context.Context, serviceID string) ([]models.Service, error) {
	// تنفيذ الحصول على خدمات مشابهة
	return []models.Service{}, nil
}

// ================================
// تطبيقات CategoryService
// ================================

func (s *categoryServiceImpl) GetCategories(ctx context.Context, params CategoryQueryParams) ([]models.Category, error) {
	// تنفيذ الحصول على الفئات
	return []models.Category{}, nil
}

func (s *categoryServiceImpl) GetCategoryByID(ctx context.Context, categoryID string) (*models.Category, error) {
	// تنفيذ الحصول على الفئة
	return &models.Category{}, nil
}

func (s *categoryServiceImpl) CreateCategory(ctx context.Context, req CategoryCreateRequest) (*models.Category, error) {
	categoryID := fmt.Sprintf("category_%d", time.Now().UnixNano())
	category := &models.Category{
		ID:        categoryID,
		Name:      req.Name,
		Slug:      req.Slug,
		Image:     req.Image,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err := s.db.ExecContext(ctx,
		`INSERT INTO categories (id,name,slug,image,is_active,created_at,updated_at)
		 VALUES (?,?,?,?,?,?,?)`,
		category.ID, category.Name, category.Slug, category.Image, category.IsActive, category.CreatedAt, category.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return category, nil
}

func (s *categoryServiceImpl) UpdateCategory(ctx context.Context, categoryID string, req CategoryUpdateRequest) (*models.Category, error) {
	// تنفيذ تحديث الفئة
	return &models.Category{}, nil
}

func (s *categoryServiceImpl) DeleteCategory(ctx context.Context, categoryID string) error {
	// تنفيذ حذف الفئة
	return nil
}

func (s *categoryServiceImpl) GetCategoryTree(ctx context.Context) ([]CategoryNode, error) {
	// تنفيذ الحصول على شجرة الفئات
	return []CategoryNode{}, nil
}

// ================================
// تطبيقات OrderService
// ================================

func (s *orderServiceImpl) CreateOrder(ctx context.Context, req OrderCreateRequest) (*models.Order, error) {
	orderID := fmt.Sprintf("order_%d", time.Now().UnixNano())
	order := &models.Order{
		ID:         orderID,
		UserID:     req.UserID,
		ServiceID:  req.ServiceID,
		Status:     "pending",
		Amount:     req.Amount,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	_, err := s.db.ExecContext(ctx,
		`INSERT INTO orders (id,user_id,service_id,status,amount,created_at,updated_at)
		 VALUES (?,?,?,?,?,?,?)`,
		order.ID, order.UserID, order.ServiceID, order.Status, order.Amount, order.CreatedAt, order.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (s *orderServiceImpl) GetOrderByID(ctx context.Context, orderID string) (*models.Order, error) {
	row := s.db.QueryRowContext(ctx,
		"SELECT id,user_id,service_id,status,amount,created_at,updated_at FROM orders WHERE id=?",
		orderID,
	)
	order := &models.Order{}
	err := row.Scan(&order.ID, &order.UserID, &order.ServiceID, &order.Status, &order.Amount, &order.CreatedAt, &order.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("الطلب غير موجود")
		}
		return nil, err
	}
	return order, nil
}

func (s *orderServiceImpl) GetUserOrders(ctx context.Context, userID string, params OrderQueryParams) ([]models.Order, error) {
	// تنفيذ الحصول على طلبات المستخدم
	return []models.Order{}, nil
}

func (s *orderServiceImpl) UpdateOrderStatus(ctx context.Context, orderID string, status string, notes string) (*models.Order, error) {
	// تنفيذ تحديث حالة الطلب
	return &models.Order{}, nil
}

func (s *orderServiceImpl) CancelOrder(ctx context.Context, orderID string, reason string) (*models.Order, error) {
	// تنفيذ إلغاء الطلب
	return &models.Order{}, nil
}

func (s *orderServiceImpl) GetOrderStats(ctx context.Context, timeframe string) (*OrderStats, error) {
	// تنفيذ الحصول على إحصائيات الطلبات
	return &OrderStats{}, nil
}

// ================================
// تطبيقات CacheService
// ================================

func (c *cacheServiceImpl) Get(key string) (interface{}, error) {
	val, ok := c.store[key]
	if !ok {
		return nil, fmt.Errorf("key not found")
	}
	return val, nil
}

func (c *cacheServiceImpl) Set(key string, value interface{}, expiration time.Duration) error {
	c.store[key] = value
	// يمكن إضافة منطق انتهاء الصلاحية هنا
	return nil
}

func (c *cacheServiceImpl) Delete(key string) error {
	delete(c.store, key)
	return nil
}

func (c *cacheServiceImpl) Exists(key string) (bool, error) {
	_, ok := c.store[key]
	return ok, nil
}

func (c *cacheServiceImpl) Flush() error {
	c.store = make(map[string]interface{})
	return nil
}

// ================================
// دوال مساعدة
// ================================

func serializeStrings(arr []string) string {
	return strings.Join(arr, ",")
}

// ================================
// دوال الإنشاء
// ================================

func NewAuthService(db *sql.DB) AuthService {
	return &authServiceImpl{db: db}
}

func NewUserService(db *sql.DB) UserService {
	return &userServiceImpl{db: db}
}

func NewServiceService(db *sql.DB) ServiceService {
	return &serviceServiceImpl{db: db}
}

func NewCategoryService(db *sql.DB) CategoryService {
	return &categoryServiceImpl{db: db}
}

func NewOrderService(db *sql.DB) OrderService {
	return &orderServiceImpl{db: db}
}

func NewCacheService() CacheService {
	return &cacheServiceImpl{store: make(map[string]interface{})}
}