package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/nawthtech/nawthtech/backend/internal/config"
	"github.com/nawthtech/nawthtech/backend/internal/models"
	"go.uber.org/zap"
)

// ================================
// هياكل الطلبات (Requests)
// ================================

type AuthRegisterRequest struct {
	FirstName string `json:"first_name" validate:"required,min=2,max=50"`
	LastName  string `json:"last_name" validate:"required,min=2,max=50"`
	Email     string `json:"email" validate:"required,email"`
	Username  string `json:"username" validate:"required,min=3,max=30,alphanum"`
	Phone     string `json:"phone" validate:"omitempty,min=10,max=20"`
	Password  string `json:"password" validate:"required,min=8"`
}

type AuthLoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=8"`
}

type TokenClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	Exp    int64  `json:"exp"`
}

type AuthResponse struct {
	User         *models.User `json:"user"`
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	ExpiresAt    time.Time    `json:"expires_at"`
}

type UserUpdateRequest struct {
	FirstName string `json:"first_name" validate:"omitempty,min=2,max=50"`
	LastName  string `json:"last_name" validate:"omitempty,min=2,max=50"`
	Phone     string `json:"phone" validate:"omitempty,min=10,max=20"`
	Avatar    string `json:"avatar" validate:"omitempty,url"`
}

type UserQueryParams struct {
	Page  int    `json:"page" validate:"min=1"`
	Limit int    `json:"limit" validate:"min=1,max=100"`
	Role  string `json:"role"`
	Email string `json:"email"`
}

type UserStats struct {
	TotalOrders   int     `json:"total_orders"`
	TotalSpent    float64 `json:"total_spent"`
	ActiveSince   string  `json:"active_since"`
	ServicesCount int     `json:"services_count"`
}

type ServiceCreateRequest struct {
	Title       string   `json:"title" validate:"required,min=5,max=200"`
	Description string   `json:"description" validate:"required,min=10,max=5000"`
	Price       float64  `json:"price" validate:"required,min=0"`
	Duration    int      `json:"duration" validate:"required,min=1"`
	CategoryID  string   `json:"category_id" validate:"required"`
	ProviderID  string   `json:"provider_id" validate:"required"`
	Images      []string `json:"images" validate:"max=10"`
	Tags        []string `json:"tags" validate:"max=20"`
}

type ServiceUpdateRequest struct {
	Title       string   `json:"title" validate:"omitempty,min=5,max=200"`
	Description string   `json:"description" validate:"omitempty,min=10,max=5000"`
	Price       float64  `json:"price" validate:"omitempty,min=0"`
	Duration    int      `json:"duration" validate:"omitempty,min=1"`
	CategoryID  string   `json:"category_id"`
	Images      []string `json:"images" validate:"max=10"`
	Tags        []string `json:"tags" validate:"max=20"`
	IsActive    bool     `json:"is_active"`
	IsFeatured  bool     `json:"is_featured"`
}

type ServiceQueryParams struct {
	Page       int     `json:"page" validate:"min=1"`
	Limit      int     `json:"limit" validate:"min=1,max=100"`
	CategoryID string  `json:"category_id"`
	ProviderID string  `json:"provider_id"`
	MinPrice   float64 `json:"min_price" validate:"min=0"`
	MaxPrice   float64 `json:"max_price" validate:"min=0"`
	IsActive   bool    `json:"is_active"`
	IsFeatured bool    `json:"is_featured"`
	Search     string  `json:"search"`
}

type CategoryCreateRequest struct {
	Name  string `json:"name" validate:"required,min=2,max=100"`
	Slug  string `json:"slug" validate:"required,min=2,max=100,slug"`
	Image string `json:"image" validate:"omitempty,url"`
}

type CategoryUpdateRequest struct {
	Name     string `json:"name" validate:"omitempty,min=2,max=100"`
	Slug     string `json:"slug" validate:"omitempty,min=2,max=100,slug"`
	Image    string `json:"image" validate:"omitempty,url"`
	IsActive bool   `json:"is_active"`
}

type CategoryQueryParams struct {
	Page     int  `json:"page" validate:"min=1"`
	Limit    int  `json:"limit" validate:"min=1,max=100"`
	IsActive bool `json:"is_active"`
}

type CategoryNode struct {
	Category  *models.Category `json:"category"`
	Children  []CategoryNode   `json:"children"`
	Services  int              `json:"services_count"`
}

type OrderCreateRequest struct {
	UserID    string  `json:"user_id" validate:"required"`
	ServiceID string  `json:"service_id" validate:"required"`
	Amount    float64 `json:"amount" validate:"required,min=0"`
	Notes     string  `json:"notes" validate:"max=500"`
}

type OrderQueryParams struct {
	Page   int    `json:"page" validate:"min=1"`
	Limit  int    `json:"limit" validate:"min=1,max=100"`
	Status string `json:"status"`
	UserID string `json:"user_id"`
}

type OrderStats struct {
	TotalOrders   int     `json:"total_orders"`
	PendingOrders int     `json:"pending_orders"`
	Completed     int     `json:"completed_orders"`
	Cancelled     int     `json:"cancelled_orders"`
	TotalRevenue  float64 `json:"total_revenue"`
	AvgOrderValue float64 `json:"avg_order_value"`
}

type PaymentIntentRequest struct {
	OrderID   string  `json:"order_id" validate:"required"`
	Amount    float64 `json:"amount" validate:"required,min=0"`
	Currency  string  `json:"currency" validate:"required,len=3"`
	Customer  string  `json:"customer,omitempty"`
	ReturnURL string  `json:"return_url" validate:"omitempty,url"`
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
	Page     int       `json:"page" validate:"min=1"`
	Limit    int       `json:"limit" validate:"min=1,max=100"`
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

type UploadRequest struct {
	FileName string `json:"file_name" validate:"required"`
	FileType string `json:"file_type" validate:"required"`
	FileSize int64  `json:"file_size" validate:"required,min=1"`
}

type UploadResult struct {
	URL      string    `json:"url"`
	FileName string    `json:"file_name"`
	FileType string    `json:"file_type"`
	FileSize int64     `json:"file_size"`
	Uploaded time.Time `json:"uploaded"`
}

type NotificationCreateRequest struct {
	UserID  string `json:"user_id" validate:"required"`
	Title   string `json:"title" validate:"required,min=2,max=200"`
	Message string `json:"message" validate:"required,min=2,max=1000"`
	Type    string `json:"type" validate:"required,oneof=info success warning error"`
}

type NotificationQueryParams struct {
	Page   int    `json:"page" validate:"min=1"`
	Limit  int    `json:"limit" validate:"min=1,max=100"`
	UserID string `json:"user_id"`
	Type   string `json:"type"`
	Read   *bool  `json:"read"`
}

type DashboardStats struct {
	TotalUsers      int64   `json:"total_users"`
	TotalServices   int64   `json:"total_services"`
	TotalOrders     int64   `json:"total_orders"`
	TotalRevenue    float64 `json:"total_revenue"`
	ActiveUsers     int64   `json:"active_users"`
	PendingOrders   int64   `json:"pending_orders"`
	CompletedOrders int64   `json:"completed_orders"`
}

type SystemLogQuery struct {
	Page     int       `json:"page" validate:"min=1"`
	Limit    int       `json:"limit" validate:"min=1,max=100"`
	Level    string    `json:"level"`
	UserID   string    `json:"user_id"`
	FromDate time.Time `json:"from_date,omitempty"`
	ToDate   time.Time `json:"to_date,omitempty"`
}

// ================================
// الواجهات (Interfaces)
// ================================

type AuthService interface {
	Register(ctx context.Context, req AuthRegisterRequest) (*AuthResponse, error)
	Login(ctx context.Context, req AuthLoginRequest) (*AuthResponse, error)
	Logout(ctx context.Context, token string) error
	RefreshToken(ctx context.Context, refreshToken string) (*AuthResponse, error)
	VerifyToken(ctx context.Context, token string) (*TokenClaims, error)
	ForgotPassword(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, token string, newPassword string) error
	ChangePassword(ctx context.Context, userID string, req ChangePasswordRequest) error
}

type UserService interface {
	GetProfile(ctx context.Context, userID string) (*models.User, error)
	UpdateProfile(ctx context.Context, userID string, req UserUpdateRequest) (*models.User, error)
	UpdateAvatar(ctx context.Context, userID string, avatarURL string) error
	DeleteAccount(ctx context.Context, userID string) error
	SearchUsers(ctx context.Context, query string, params UserQueryParams) ([]models.User, error)
	GetUserStats(ctx context.Context, userID string) (*UserStats, error)
}

type ServiceService interface {
	CreateService(ctx context.Context, req ServiceCreateRequest) (*models.Service, error)
	GetServiceByID(ctx context.Context, serviceID string) (*models.Service, error)
	UpdateService(ctx context.Context, serviceID string, req ServiceUpdateRequest) (*models.Service, error)
	DeleteService(ctx context.Context, serviceID string) error
	GetServices(ctx context.Context, params ServiceQueryParams) ([]models.Service, error)
	SearchServices(ctx context.Context, query string, params ServiceQueryParams) ([]models.Service, error)
	GetFeaturedServices(ctx context.Context) ([]models.Service, error)
	GetSimilarServices(ctx context.Context, serviceID string) ([]models.Service, error)
}

type CategoryService interface {
	GetCategories(ctx context.Context, params CategoryQueryParams) ([]models.Category, error)
	GetCategoryByID(ctx context.Context, categoryID string) (*models.Category, error)
	CreateCategory(ctx context.Context, req CategoryCreateRequest) (*models.Category, error)
	UpdateCategory(ctx context.Context, categoryID string, req CategoryUpdateRequest) (*models.Category, error)
	DeleteCategory(ctx context.Context, categoryID string) error
	GetCategoryTree(ctx context.Context) ([]CategoryNode, error)
}

type OrderService interface {
	CreateOrder(ctx context.Context, req OrderCreateRequest) (*models.Order, error)
	GetOrderByID(ctx context.Context, orderID string) (*models.Order, error)
	GetUserOrders(ctx context.Context, userID string, params OrderQueryParams) ([]models.Order, error)
	UpdateOrderStatus(ctx context.Context, orderID string, status string, notes string) (*models.Order, error)
	CancelOrder(ctx context.Context, orderID string, reason string) (*models.Order, error)
	GetOrderStats(ctx context.Context, timeframe string) (*OrderStats, error)
}

type PaymentService interface {
	CreatePaymentIntent(ctx context.Context, req PaymentIntentRequest) (*PaymentIntent, error)
	ConfirmPayment(ctx context.Context, paymentID string, confirmationData map[string]interface{}) (*PaymentResult, error)
	GetPaymentHistory(ctx context.Context, userID string, params PaymentQueryParams) ([]models.Payment, error)
	ValidatePayment(ctx context.Context, paymentData map[string]interface{}) (*PaymentValidation, error)
}

type UploadService interface {
	UploadFile(ctx context.Context, req UploadRequest, fileData []byte) (*UploadResult, error)
	DeleteFile(ctx context.Context, fileID string) error
	GetFile(ctx context.Context, fileID string) (*models.File, error)
	GetUserFiles(ctx context.Context, userID string) ([]models.File, error)
	GeneratePresignedURL(ctx context.Context, fileName, fileType string) (string, error)
}

type NotificationService interface {
	CreateNotification(ctx context.Context, req NotificationCreateRequest) (*models.Notification, error)
	GetUserNotifications(ctx context.Context, userID string, params NotificationQueryParams) ([]models.Notification, error)
	MarkAsRead(ctx context.Context, notificationID string) error
	MarkAllAsRead(ctx context.Context, userID string) error
	DeleteNotification(ctx context.Context, notificationID string) error
	GetUnreadCount(ctx context.Context, userID string) (int64, error)
}

type AdminService interface {
	GetDashboardStats(ctx context.Context) (*DashboardStats, error)
	GetUsers(ctx context.Context, params UserQueryParams) ([]models.User, error)
	GetSystemLogs(ctx context.Context, params SystemLogQuery) ([]models.SystemLog, error)
	UpdateSystemSettings(ctx context.Context, settings map[string]string) error
	BanUser(ctx context.Context, userID string, reason string) error
	UnbanUser(ctx context.Context, userID string) error
}

type CacheService interface {
	Get(key string) (interface{}, error)
	Set(key string, value interface{}, expiration time.Duration) error
	Delete(key string) error
	Exists(key string) (bool, error)
	Flush() error
}

// ================================
// التطبيقات (Implementations)
// ================================

type authServiceImpl struct {
	db *sql.DB
}

type userServiceImpl struct {
	db *sql.DB
}

type serviceServiceImpl struct {
	db *sql.DB
}

type categoryServiceImpl struct {
	db *sql.DB
}

type orderServiceImpl struct {
	db *sql.DB
}

type paymentServiceImpl struct {
	db *sql.DB
}

type uploadServiceImpl struct {
	db *sql.DB
}

type notificationServiceImpl struct {
	db *sql.DB
}

type adminServiceImpl struct {
	db *sql.DB
}

type cacheServiceImpl struct {
	store map[string]interface{}
	mu    sync.RWMutex
}

// ================================
// Service Container
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
	
	db     *sql.DB
	config *config.Config
	logger *zap.Logger
}

// ================================
// دوال الإنشاء (Factory Functions)
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

func NewPaymentService(db *sql.DB) PaymentService {
	return &paymentServiceImpl{db: db}
}

func NewUploadService(db *sql.DB) UploadService {
	return &uploadServiceImpl{db: db}
}

func NewNotificationService(db *sql.DB) NotificationService {
	return &notificationServiceImpl{db: db}
}

func NewAdminService(db *sql.DB) AdminService {
	return &adminServiceImpl{db: db}
}

func NewCacheService() CacheService {
	return &cacheServiceImpl{
		store: make(map[string]interface{}),
		mu:    sync.RWMutex{},
	}
}

func NewServiceContainer(db *sql.DB) *ServiceContainer {
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
		db:           db,
	}
}

func NewServiceContainerWithConfig(db *sql.DB, cfg *config.Config, logger *zap.Logger) *ServiceContainer {
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
		db:           db,
		config:       cfg,
		logger:       logger,
	}
}

// ================================
// دوال SQL للإنشاء
// ================================

func CreateTablesSQL() []string {
	return []string{
		// جدول المستخدمين
		`CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			email TEXT UNIQUE NOT NULL,
			username TEXT UNIQUE NOT NULL,
			password_hash TEXT NOT NULL,
			first_name TEXT NOT NULL,
			last_name TEXT NOT NULL,
			phone TEXT,
			avatar TEXT,
			role TEXT DEFAULT 'user',
			status TEXT DEFAULT 'active',
			email_verified BOOLEAN DEFAULT FALSE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			last_login TIMESTAMP
		)`,

		// جدول الفئات
		`CREATE TABLE IF NOT EXISTS categories (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			slug TEXT UNIQUE NOT NULL,
			image TEXT,
			description TEXT,
			parent_id TEXT,
			is_active BOOLEAN DEFAULT TRUE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (parent_id) REFERENCES categories(id) ON DELETE SET NULL
		)`,

		// جدول الخدمات
		`CREATE TABLE IF NOT EXISTS services (
			id TEXT PRIMARY KEY,
			title TEXT NOT NULL,
			description TEXT NOT NULL,
			price REAL NOT NULL,
			duration INTEGER NOT NULL,
			category_id TEXT NOT NULL,
			provider_id TEXT NOT NULL,
			images TEXT DEFAULT '[]',
			tags TEXT DEFAULT '[]',
			is_active BOOLEAN DEFAULT TRUE,
			is_featured BOOLEAN DEFAULT FALSE,
			rating REAL DEFAULT 0,
			review_count INTEGER DEFAULT 0,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE,
			FOREIGN KEY (provider_id) REFERENCES users(id) ON DELETE CASCADE
		)`,

		// جدول الطلبات
		`CREATE TABLE IF NOT EXISTS orders (
			id TEXT PRIMARY KEY,
			user_id TEXT NOT NULL,
			service_id TEXT NOT NULL,
			status TEXT DEFAULT 'pending',
			amount REAL NOT NULL,
			notes TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
			FOREIGN KEY (service_id) REFERENCES services(id) ON DELETE CASCADE
		)`,

		// جدول المدفوعات
		`CREATE TABLE IF NOT EXISTS payments (
			id TEXT PRIMARY KEY,
			order_id TEXT NOT NULL,
			amount REAL NOT NULL,
			currency TEXT DEFAULT 'USD',
			status TEXT DEFAULT 'pending',
			payment_method TEXT,
			transaction_id TEXT UNIQUE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE
		)`,

		// جدول الإشعارات
		`CREATE TABLE IF NOT EXISTS notifications (
			id TEXT PRIMARY KEY,
			user_id TEXT NOT NULL,
			title TEXT NOT NULL,
			message TEXT NOT NULL,
			type TEXT DEFAULT 'info',
			is_read BOOLEAN DEFAULT FALSE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		)`,

		// جدول الملفات
		`CREATE TABLE IF NOT EXISTS files (
			id TEXT PRIMARY KEY,
			user_id TEXT NOT NULL,
			name TEXT NOT NULL,
			url TEXT NOT NULL,
			size INTEGER,
			type TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		)`,

		// جدول السجلات
		`CREATE TABLE IF NOT EXISTS system_logs (
			id TEXT PRIMARY KEY,
			user_id TEXT,
			level TEXT NOT NULL,
			action TEXT NOT NULL,
			resource TEXT,
			details TEXT,
			ip_address TEXT,
			user_agent TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
	}
}

// ================================
// دوال مساعدة (Helper Functions)
// ================================

func serializeStrings(arr []string) string {
	if len(arr) == 0 {
		return "[]"
	}
	data, _ := json.Marshal(arr)
	return string(data)
}

func deserializeStrings(s string) ([]string, error) {
	if s == "" || s == "[]" {
		return []string{}, nil
	}
	var arr []string
	err := json.Unmarshal([]byte(s), &arr)
	return arr, err
}

func generateID(prefix string) string {
	return fmt.Sprintf("%s_%d", prefix, time.Now().UnixNano())
}

func buildSearchQuery(baseQuery string, search string, args []interface{}) (string, []interface{}) {
	if search == "" {
		return baseQuery, args
	}
	
	search = strings.TrimSpace(search)
	searchTerms := strings.Fields(search)
	
	var conditions []string
	for _, term := range searchTerms {
		conditions = append(conditions, "(title LIKE ? OR description LIKE ? OR tags LIKE ?)")
		args = append(args, "%"+term+"%", "%"+term+"%", "%"+term+"%")
	}
	
	whereClause := strings.Join(conditions, " AND ")
	if strings.Contains(baseQuery, "WHERE") {
		return baseQuery + " AND (" + whereClause + ")", args
	}
	return baseQuery + " WHERE " + whereClause, args
}

func validatePaginationParams(page, limit int) (int, int) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	return page, limit
}

func calculateOffset(page, limit int) int {
	return (page - 1) * limit
}

// ================================
// تطبيقات الخدمات (Service Implementations)
// ================================

// AuthService Implementation
func (s *authServiceImpl) Register(ctx context.Context, req AuthRegisterRequest) (*AuthResponse, error) {
	userID := generateID("user")
	
	// تشفير كلمة المرور (يجب استخدام bcrypt في الواقع)
	passwordHash := fmt.Sprintf("hash_%s", req.Password)
	
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO users (id, email, username, password_hash, first_name, last_name, phone, role, status, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		userID, req.Email, req.Username, passwordHash, req.FirstName, req.LastName, req.Phone,
		"user", "active", time.Now(), time.Now(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to register user: %w", err)
	}
	
	user := &models.User{
		ID:        userID,
		Email:     req.Email,
		Username:  req.Username,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Phone:     req.Phone,
		Role:      "user",
		Status:    "active",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	
	return &AuthResponse{
		User:         user,
		AccessToken:  generateID("access"),
		RefreshToken: generateID("refresh"),
		ExpiresAt:    time.Now().Add(24 * time.Hour),
	}, nil
}

func (s *authServiceImpl) Login(ctx context.Context, req AuthLoginRequest) (*AuthResponse, error) {
	var user models.User
	row := s.db.QueryRowContext(ctx,
		`SELECT id, email, username, first_name, last_name, phone, avatar, role, status, created_at, updated_at
		 FROM users WHERE email = ? AND status = 'active'`,
		req.Email)
	
	err := row.Scan(
		&user.ID, &user.Email, &user.Username, &user.FirstName, &user.LastName,
		&user.Phone, &user.Avatar, &user.Role, &user.Status, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("invalid credentials")
		}
		return nil, fmt.Errorf("failed to login: %w", err)
	}
	
	// تحديث آخر تسجيل دخول
	_, err = s.db.ExecContext(ctx,
		"UPDATE users SET last_login = ? WHERE id = ?",
		time.Now(), user.ID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update last login: %w", err)
	}
	
	return &AuthResponse{
		User:         &user,
		AccessToken:  generateID("access"),
		RefreshToken: generateID("refresh"),
		ExpiresAt:    time.Now().Add(24 * time.Hour),
	}, nil
}

func (s *authServiceImpl) Logout(ctx context.Context, token string) error {
	// في الواقع، يجب إضافة الرمز إلى القائمة السوداء
	return nil
}

func (s *authServiceImpl) RefreshToken(ctx context.Context, refreshToken string) (*AuthResponse, error) {
	return nil, errors.New("not implemented")
}

func (s *authServiceImpl) VerifyToken(ctx context.Context, token string) (*TokenClaims, error) {
	return nil, errors.New("not implemented")
}

func (s *authServiceImpl) ForgotPassword(ctx context.Context, email string) error {
	return errors.New("not implemented")
}

func (s *authServiceImpl) ResetPassword(ctx context.Context, token string, newPassword string) error {
	return errors.New("not implemented")
}

func (s *authServiceImpl) ChangePassword(ctx context.Context, userID string, req ChangePasswordRequest) error {
	return errors.New("not implemented")
}

// UserService Implementation
func (s *userServiceImpl) GetProfile(ctx context.Context, userID string) (*models.User, error) {
	var user models.User
	row := s.db.QueryRowContext(ctx,
		`SELECT id, email, username, first_name, last_name, phone, avatar, role, status, email_verified, created_at, updated_at, last_login
		 FROM users WHERE id = ?`,
		userID)
	
	err := row.Scan(
		&user.ID, &user.Email, &user.Username, &user.FirstName, &user.LastName,
		&user.Phone, &user.Avatar, &user.Role, &user.Status, &user.EmailVerified,
		&user.CreatedAt, &user.UpdatedAt, &user.LastLogin,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to get profile: %w", err)
	}
	
	return &user, nil
}

func (s *userServiceImpl) UpdateProfile(ctx context.Context, userID string, req UserUpdateRequest) (*models.User, error) {
	_, err := s.db.ExecContext(ctx,
		`UPDATE users SET first_name = ?, last_name = ?, phone = ?, avatar = ?, updated_at = ?
		 WHERE id = ?`,
		req.FirstName, req.LastName, req.Phone, req.Avatar, time.Now(), userID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update profile: %w", err)
	}
	
	return s.GetProfile(ctx, userID)
}

func (s *userServiceImpl) UpdateAvatar(ctx context.Context, userID string, avatarURL string) error {
	_, err := s.db.ExecContext(ctx,
		"UPDATE users SET avatar = ?, updated_at = ? WHERE id = ?",
		avatarURL, time.Now(), userID,
	)
	return err
}

func (s *userServiceImpl) DeleteAccount(ctx context.Context, userID string) error {
	_, err := s.db.ExecContext(ctx,
		"UPDATE users SET status = 'deleted', updated_at = ? WHERE id = ?",
		time.Now(), userID,
	)
	return err
}

func (s *userServiceImpl) SearchUsers(ctx context.Context, query string, params UserQueryParams) ([]models.User, error) {
	page, limit := validatePaginationParams(params.Page, params.Limit)
	offset := calculateOffset(page, limit)
	
	sqlQuery := `SELECT id, email, username, first_name, last_name, phone, avatar, role, status, created_at
				 FROM users WHERE status != 'deleted'`
	args := []interface{}{}
	
	if query != "" {
		sqlQuery += " AND (email LIKE ? OR username LIKE ? OR first_name LIKE ? OR last_name LIKE ?)"
		args = append(args, "%"+query+"%", "%"+query+"%", "%"+query+"%", "%"+query+"%")
	}
	
	if params.Role != "" {
		sqlQuery += " AND role = ?"
		args = append(args, params.Role)
	}
	
	if params.Email != "" {
		sqlQuery += " AND email = ?"
		args = append(args, params.Email)
	}
	
	sqlQuery += " LIMIT ? OFFSET ?"
	args = append(args, limit, offset)
	
	rows, err := s.db.QueryContext(ctx, sqlQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to search users: %w", err)
	}
	defer rows.Close()
	
	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.ID, &user.Email, &user.Username, &user.FirstName, &user.LastName,
			&user.Phone, &user.Avatar, &user.Role, &user.Status, &user.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}
	
	return users, nil
}

func (s *userServiceImpl) GetUserStats(ctx context.Context, userID string) (*UserStats, error) {
	// الحصول على عدد الطلبات
	var totalOrders int
	err := s.db.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM orders WHERE user_id = ?",
		userID,
	).Scan(&totalOrders)
	if err != nil {
		return nil, fmt.Errorf("failed to get total orders: %w", err)
	}
	
	// الحصول على المبلغ الإجمالي
	var totalSpent float64
	err = s.db.QueryRowContext(ctx,
		"SELECT COALESCE(SUM(amount), 0) FROM orders WHERE user_id = ? AND status = 'completed'",
		userID,
	).Scan(&totalSpent)
	if err != nil {
		return nil, fmt.Errorf("failed to get total spent: %w", err)
	}
	
	// الحصول على تاريخ الإنشاء
	var activeSince string
	err = s.db.QueryRowContext(ctx,
		"SELECT DATE(created_at) FROM users WHERE id = ?",
		userID,
	).Scan(&activeSince)
	if err != nil {
		activeSince = time.Now().Format("2006-01-02")
	}
	
	// الحصول على عدد الخدمات
	var servicesCount int
	err = s.db.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM services WHERE provider_id = ?",
		userID,
	).Scan(&servicesCount)
	if err != nil {
		servicesCount = 0
	}
	
	return &UserStats{
		TotalOrders:   totalOrders,
		TotalSpent:    totalSpent,
		ActiveSince:   activeSince,
		ServicesCount: servicesCount,
	}, nil
}

// ServiceService Implementation
func (s *serviceServiceImpl) CreateService(ctx context.Context, req ServiceCreateRequest) (*models.Service, error) {
	serviceID := generateID("service")
	imagesJSON := serializeStrings(req.Images)
	tagsJSON := serializeStrings(req.Tags)
	
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
		ReviewCount: 0,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO services (id, title, description, price, duration, category_id, provider_id, images, tags, is_active, is_featured, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		 serviceID, req.Title, req.Description, req.Price, req.Duration, req.CategoryID, req.ProviderID,
		 imagesJSON, tagsJSON, true, false, time.Now(), time.Now(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create service: %w", err)
	}
	
	return service, nil
}

func (s *serviceServiceImpl) GetServiceByID(ctx context.Context, serviceID string) (*models.Service, error) {
	var service models.Service
	var imagesJSON, tagsJSON string
	
	row := s.db.QueryRowContext(ctx,
		`SELECT id, title, description, price, duration, category_id, provider_id, images, tags, is_active, is_featured, rating, review_count, created_at, updated_at
		 FROM services WHERE id = ?`,
		serviceID)
	
	err := row.Scan(
		&service.ID, &service.Title, &service.Description, &service.Price, &service.Duration,
		&service.CategoryID, &service.ProviderID, &imagesJSON, &tagsJSON,
		&service.IsActive, &service.IsFeatured, &service.Rating, &service.ReviewCount,
		&service.CreatedAt, &service.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("service not found")
		}
		return nil, fmt.Errorf("failed to get service: %w", err)
	}
	
	service.Images, _ = deserializeStrings(imagesJSON)
	service.Tags, _ = deserializeStrings(tagsJSON)
	
	return &service, nil
}

func (s *serviceServiceImpl) UpdateService(ctx context.Context, serviceID string, req ServiceUpdateRequest) (*models.Service, error) {
	imagesJSON := serializeStrings(req.Images)
	tagsJSON := serializeStrings(req.Tags)
	
	_, err := s.db.ExecContext(ctx,
		`UPDATE services SET title = ?, description = ?, price = ?, duration = ?, category_id = ?, images = ?, tags = ?, is_active = ?, is_featured = ?, updated_at = ?
		 WHERE id = ?`,
		req.Title, req.Description, req.Price, req.Duration, req.CategoryID, imagesJSON, tagsJSON,
		req.IsActive, req.IsFeatured, time.Now(), serviceID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update service: %w", err)
	}
	
	return s.GetServiceByID(ctx, serviceID)
}

func (s *serviceServiceImpl) DeleteService(ctx context.Context, serviceID string) error {
	_, err := s.db.ExecContext(ctx,
		"DELETE FROM services WHERE id = ?",
		serviceID,
	)
	return err
}

func (s *serviceServiceImpl) GetServices(ctx context.Context, params ServiceQueryParams) ([]models.Service, error) {
	page, limit := validatePaginationParams(params.Page, params.Limit)
	offset := calculateOffset(page, limit)
	
	sqlQuery := `SELECT id, title, description, price, duration, category_id, provider_id, images, tags, is_active, is_featured, rating, review_count, created_at
				 FROM services WHERE 1=1`
	args := []interface{}{}
	
	if params.CategoryID != "" {
		sqlQuery += " AND category_id = ?"
		args = append(args, params.CategoryID)
	}
	
	if params.ProviderID != "" {
		sqlQuery += " AND provider_id = ?"
		args = append(args, params.ProviderID)
	}
	
	if params.IsActive {
		sqlQuery += " AND is_active = TRUE"
	}
	
	if params.IsFeatured {
		sqlQuery += " AND is_featured = TRUE"
	}
	
	if params.MinPrice > 0 {
		sqlQuery += " AND price >= ?"
		args = append(args, params.MinPrice)
	}
	
	if params.MaxPrice > 0 {
		sqlQuery += " AND price <= ?"
		args = append(args, params.MaxPrice)
	}
	
	sqlQuery, args = buildSearchQuery(sqlQuery, params.Search, args)
	sqlQuery += " ORDER BY created_at DESC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)
	
	rows, err := s.db.QueryContext(ctx, sqlQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get services: %w", err)
	}
	defer rows.Close()
	
	var services []models.Service
	for rows.Next() {
		var service models.Service
		var imagesJSON, tagsJSON string
		
		err := rows.Scan(
			&service.ID, &service.Title, &service.Description, &service.Price, &service.Duration,
			&service.CategoryID, &service.ProviderID, &imagesJSON, &tagsJSON,
			&service.IsActive, &service.IsFeatured, &service.Rating, &service.ReviewCount,
			&service.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan service: %w", err)
		}
		
		service.Images, _ = deserializeStrings(imagesJSON)
		service.Tags, _ = deserializeStrings(tagsJSON)
		services = append(services, service)
	}
	
	return services, nil
}

func (s *serviceServiceImpl) SearchServices(ctx context.Context, query string, params ServiceQueryParams) ([]models.Service, error) {
	params.Search = query
	return s.GetServices(ctx, params)
}

func (s *serviceServiceImpl) GetFeaturedServices(ctx context.Context) ([]models.Service, error) {
	params := ServiceQueryParams{
		Page:       1,
		Limit:      10,
		IsFeatured: true,
		IsActive:   true,
	}
	return s.GetServices(ctx, params)
}

func (s *serviceServiceImpl) GetSimilarServices(ctx context.Context, serviceID string) ([]models.Service, error) {
	// الحصول على خدمة الحالية لمعرفة فئتها
	service, err := s.GetServiceByID(ctx, serviceID)
	if err != nil {
		return nil, err
	}
	
	// البحث عن خدمات مشابهة في نفس الفئة
	params := ServiceQueryParams{
		Page:       1,
		Limit:      5,
		CategoryID: service.CategoryID,
		IsActive:   true,
	}
	
	services, err := s.GetServices(ctx, params)
	if err != nil {
		return nil, err
	}
	
	// إزالة الخدمة الحالية من النتائج
	var similarServices []models.Service
	for _, svc := range services {
		if svc.ID != serviceID {
			similarServices = append(similarServices, svc)
		}
	}
	
	return similarServices, nil
}

// CategoryService Implementation
func (s *categoryServiceImpl) GetCategories(ctx context.Context, params CategoryQueryParams) ([]models.Category, error) {
	page, limit := validatePaginationParams(params.Page, params.Limit)
	offset := calculateOffset(page, limit)
	
	sqlQuery := `SELECT id, name, slug, image, description, parent_id, is_active, created_at, updated_at
				 FROM categories WHERE 1=1`
	args := []interface{}{}
	
	if params.IsActive {
		sqlQuery += " AND is_active = TRUE"
	}
	
	sqlQuery += " ORDER BY name ASC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)
	
	rows, err := s.db.QueryContext(ctx, sqlQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get categories: %w", err)
	}
	defer rows.Close()
	
	var categories []models.Category
	for rows.Next() {
		var category models.Category
		err := rows.Scan(
			&category.ID, &category.Name, &category.Slug, &category.Image, &category.Description,
			&category.ParentID, &category.IsActive, &category.CreatedAt, &category.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan category: %w", err)
		}
		categories = append(categories, category)
	}
	
	return categories, nil
}

func (s *categoryServiceImpl) GetCategoryByID(ctx context.Context, categoryID string) (*models.Category, error) {
	var category models.Category
	row := s.db.QueryRowContext(ctx,
		`SELECT id, name, slug, image, description, parent_id, is_active, created_at, updated_at
		 FROM categories WHERE id = ?`,
		categoryID)
	
	err := row.Scan(
		&category.ID, &category.Name, &category.Slug, &category.Image, &category.Description,
		&category.ParentID, &category.IsActive, &category.CreatedAt, &category.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("category not found")
		}
		return nil, fmt.Errorf("failed to get category: %w", err)
	}
	
	return &category, nil
}

func (s *categoryServiceImpl) CreateCategory(ctx context.Context, req CategoryCreateRequest) (*models.Category, error) {
	categoryID := generateID("category")
	
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
		`INSERT INTO categories (id, name, slug, image, is_active, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		categoryID, req.Name, req.Slug, req.Image, true, time.Now(), time.Now(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create category: %w", err)
	}
	
	return category, nil
}

func (s *categoryServiceImpl) UpdateCategory(ctx context.Context, categoryID string, req CategoryUpdateRequest) (*models.Category, error) {
	_, err := s.db.ExecContext(ctx,
		`UPDATE categories SET name = ?, slug = ?, image = ?, is_active = ?, updated_at = ?
		 WHERE id = ?`,
		req.Name, req.Slug, req.Image, req.IsActive, time.Now(), categoryID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update category: %w", err)
	}
	
	return s.GetCategoryByID(ctx, categoryID)
}

func (s *categoryServiceImpl) DeleteCategory(ctx context.Context, categoryID string) error {
	_, err := s.db.ExecContext(ctx,
		"DELETE FROM categories WHERE id = ?",
		categoryID,
	)
	return err
}

func (s *categoryServiceImpl) GetCategoryTree(ctx context.Context) ([]CategoryNode, error) {
	// الحصول على جميع الفئات
	categories, err := s.GetCategories(ctx, CategoryQueryParams{
		Page:     1,
		Limit:    1000,
		IsActive: true,
	})
	if err != nil {
		return nil, err
	}
	
	// إنشاء خريطة للفئات
	categoryMap := make(map[string]*models.Category)
	for i := range categories {
		categoryMap[categories[i].ID] = &categories[i]
	}
	
	// بناء الشجرة
	rootNodes := []CategoryNode{}
	childrenMap := make(map[string][]*models.Category)
	
	// تجميع الفئات حسب parent_id
	for _, category := range categories {
		if category.ParentID == "" {
			// فئة جذرية
			node := CategoryNode{
				Category: &category,
				Children: []CategoryNode{},
			}
			rootNodes = append(rootNodes, node)
		} else {
			// فئة فرعية
			childrenMap[category.ParentID] = append(childrenMap[category.ParentID], &category)
		}
	}
	
	// إضافة الفئات الفرعية
	for i := range rootNodes {
		s.addChildrenToNode(&rootNodes[i], childrenMap)
	}
	
	return rootNodes, nil
}

func (s *categoryServiceImpl) addChildrenToNode(node *CategoryNode, childrenMap map[string][]*models.Category) {
	children := childrenMap[node.Category.ID]
	for _, child := range children {
		childNode := CategoryNode{
			Category: child,
			Children: []CategoryNode{},
		}
		s.addChildrenToNode(&childNode, childrenMap)
		node.Children = append(node.Children, childNode)
	}
}

// OrderService Implementation
func (s *orderServiceImpl) CreateOrder(ctx context.Context, req OrderCreateRequest) (*models.Order, error) {
	orderID := generateID("order")
	
	order := &models.Order{
		ID:        orderID,
		UserID:    req.UserID,
		ServiceID: req.ServiceID,
		Status:    "pending",
		Amount:    req.Amount,
		Notes:     req.Notes,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO orders (id, user_id, service_id, status, amount, notes, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		orderID, req.UserID, req.ServiceID, "pending", req.Amount, req.Notes, time.Now(), time.Now(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}
	
	return order, nil
}

func (s *orderServiceImpl) GetOrderByID(ctx context.Context, orderID string) (*models.Order, error) {
	var order models.Order
	row := s.db.QueryRowContext(ctx,
		`SELECT id, user_id, service_id, status, amount, notes, created_at, updated_at
		 FROM orders WHERE id = ?`,
		orderID)
	
	err := row.Scan(
		&order.ID, &order.UserID, &order.ServiceID, &order.Status, &order.Amount,
		&order.Notes, &order.CreatedAt, &order.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("order not found")
		}
		return nil, fmt.Errorf("failed to get order: %w", err)
	}
	
	return &order, nil
}

func (s *orderServiceImpl) GetUserOrders(ctx context.Context, userID string, params OrderQueryParams) ([]models.Order, error) {
	page, limit := validatePaginationParams(params.Page, params.Limit)
	offset := calculateOffset(page, limit)
	
	sqlQuery := `SELECT id, user_id, service_id, status, amount, notes, created_at, updated_at
				 FROM orders WHERE user_id = ?`
	args := []interface{}{userID}
	
	if params.Status != "" {
		sqlQuery += " AND status = ?"
		args = append(args, params.Status)
	}
	
	sqlQuery += " ORDER BY created_at DESC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)
	
	rows, err := s.db.QueryContext(ctx, sqlQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get user orders: %w", err)
	}
	defer rows.Close()
	
	var orders []models.Order
	for rows.Next() {
		var order models.Order
		err := rows.Scan(
			&order.ID, &order.UserID, &order.ServiceID, &order.Status, &order.Amount,
			&order.Notes, &order.CreatedAt, &order.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}
		orders = append(orders, order)
	}
	
	return orders, nil
}

func (s *orderServiceImpl) UpdateOrderStatus(ctx context.Context, orderID string, status string, notes string) (*models.Order, error) {
	_, err := s.db.ExecContext(ctx,
		`UPDATE orders SET status = ?, notes = COALESCE(?, notes), updated_at = ? WHERE id = ?`,
		status, notes, time.Now(), orderID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update order status: %w", err)
	}
	
	return s.GetOrderByID(ctx, orderID)
}

func (s *orderServiceImpl) CancelOrder(ctx context.Context, orderID string, reason string) (*models.Order, error) {
	return s.UpdateOrderStatus(ctx, orderID, "cancelled", reason)
}

func (s *orderServiceImpl) GetOrderStats(ctx context.Context, timeframe string) (*OrderStats, error) {
	// حساب الإحصائيات بناءً على timeframe
	var whereClause string
	switch timeframe {
	case "today":
		whereClause = "DATE(created_at) = DATE('now')"
	case "week":
		whereClause = "created_at >= DATE('now', '-7 days')"
	case "month":
		whereClause = "created_at >= DATE('now', '-1 month')"
	case "year":
		whereClause = "created_at >= DATE('now', '-1 year')"
	default:
		whereClause = "1=1"
	}
	
	stats := &OrderStats{}
	
	// إجمالي الطلبات
	err := s.db.QueryRowContext(ctx,
		fmt.Sprintf("SELECT COUNT(*) FROM orders WHERE %s", whereClause),
	).Scan(&stats.TotalOrders)
	if err != nil {
		return nil, fmt.Errorf("failed to get total orders: %w", err)
	}
	
	// الطلبات المعلقة
	err = s.db.QueryRowContext(ctx,
		fmt.Sprintf("SELECT COUNT(*) FROM orders WHERE status = 'pending' AND %s", whereClause),
	).Scan(&stats.PendingOrders)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending orders: %w", err)
	}
	
	// الطلبات المكتملة
	err = s.db.QueryRowContext(ctx,
		fmt.Sprintf("SELECT COUNT(*) FROM orders WHERE status = 'completed' AND %s", whereClause),
	).Scan(&stats.Completed)
	if err != nil {
		return nil, fmt.Errorf("failed to get completed orders: %w", err)
	}
	
	// الطلبات الملغاة
	err = s.db.QueryRowContext(ctx,
		fmt.Sprintf("SELECT COUNT(*) FROM orders WHERE status = 'cancelled' AND %s", whereClause),
	).Scan(&stats.Cancelled)
	if err != nil {
		return nil, fmt.Errorf("failed to get cancelled orders: %w", err)
	}
	
	// إجمالي الإيرادات
	err = s.db.QueryRowContext(ctx,
		fmt.Sprintf("SELECT COALESCE(SUM(amount), 0) FROM orders WHERE status = 'completed' AND %s", whereClause),
	).Scan(&stats.TotalRevenue)
	if err != nil {
		return nil, fmt.Errorf("failed to get total revenue: %w", err)
	}
	
	// متوسط قيمة الطلب
	if stats.Completed > 0 {
		stats.AvgOrderValue = stats.TotalRevenue / float64(stats.Completed)
	} else {
		stats.AvgOrderValue = 0
	}
	
	return stats, nil
}

// PaymentService Implementation
func (s *paymentServiceImpl) CreatePaymentIntent(ctx context.Context, req PaymentIntentRequest) (*PaymentIntent, error) {
	paymentID := generateID("pi")
	
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

func (s *paymentServiceImpl) ConfirmPayment(ctx context.Context, paymentID string, confirmationData map[string]interface{}) (*PaymentResult, error) {
	_, err := s.db.ExecContext(ctx,
		"UPDATE payments SET status = 'completed', updated_at = ? WHERE id = ?",
		time.Now(), paymentID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to confirm payment: %w", err)
	}
	
	// الحصول على تفاصيل الدفع
	var orderID, currency string
	var amount float64
	err = s.db.QueryRowContext(ctx,
		"SELECT order_id, amount, currency FROM payments WHERE id = ?",
		paymentID,
	).Scan(&orderID, &amount, &currency)
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
		Status:    "completed",
		Timestamp: time.Now(),
	}, nil
}

func (s *paymentServiceImpl) GetPaymentHistory(ctx context.Context, userID string, params PaymentQueryParams) ([]models.Payment, error) {
	page, limit := validatePaginationParams(params.Page, params.Limit)
	offset := calculateOffset(page, limit)
	
	sqlQuery := `SELECT p.id, p.order_id, p.amount, p.currency, p.status, p.payment_method, p.transaction_id, p.created_at, p.updated_at
				 FROM payments p
				 INNER JOIN orders o ON p.order_id = o.id
				 WHERE o.user_id = ?`
	args := []interface{}{userID}
	
	if params.Status != "" {
		sqlQuery += " AND p.status = ?"
		args = append(args, params.Status)
	}
	
	if params.OrderID != "" {
		sqlQuery += " AND p.order_id = ?"
		args = append(args, params.OrderID)
	}
	
	if !params.FromDate.IsZero() {
		sqlQuery += " AND p.created_at >= ?"
		args = append(args, params.FromDate)
	}
	
	if !params.ToDate.IsZero() {
		sqlQuery += " AND p.created_at <= ?"
		args = append(args, params.ToDate)
	}
	
	sqlQuery += " ORDER BY p.created_at DESC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)
	
	rows, err := s.db.QueryContext(ctx, sqlQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment history: %w", err)
	}
	defer rows.Close()
	
	var payments []models.Payment
	for rows.Next() {
		var payment models.Payment
		err := rows.Scan(
			&payment.ID, &payment.OrderID, &payment.Amount, &payment.Currency, &payment.Status,
			&payment.PaymentMethod, &payment.TransactionID, &payment.CreatedAt, &payment.UpdatedAt,
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
	
	return &PaymentValidation{
		Valid: true,
		Metadata: map[string]interface{}{
			"validated_at": time.Now(),
			"amount":       amount,
			"currency":     currency,
		},
	}, nil
}

// UploadService Implementation
func (s *uploadServiceImpl) UploadFile(ctx context.Context, req UploadRequest, fileData []byte) (*UploadResult, error) {
	fileID := generateID("file")
	
	// هنا يمكنك رفع الملف إلى S3 أو Cloudflare R2 أو تخزين محلي
	fileURL := fmt.Sprintf("https://storage.nawthtech.com/files/%s/%s", fileID, req.FileName)
	
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO files (id, name, url, size, type, created_at)
		 VALUES (?, ?, ?, ?, ?, ?)`,
		fileID, req.FileName, fileURL, req.FileSize, req.FileType, time.Now(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to save file metadata: %w", err)
	}
	
	return &UploadResult{
		URL:      fileURL,
		FileName: req.FileName,
		FileType: req.FileType,
		FileSize: req.FileSize,
		Uploaded: time.Now(),
	}, nil
}

func (s *uploadServiceImpl) DeleteFile(ctx context.Context, fileID string) error {
	_, err := s.db.ExecContext(ctx,
		"DELETE FROM files WHERE id = ?",
		fileID,
	)
	return err
}

func (s *uploadServiceImpl) GetFile(ctx context.Context, fileID string) (*models.File, error) {
	var file models.File
	row := s.db.QueryRowContext(ctx,
		`SELECT id, user_id, name, url, size, type, created_at
		 FROM files WHERE id = ?`,
		fileID)
	
	err := row.Scan(
		&file.ID, &file.UserID, &file.Name, &file.URL, &file.Size,
		&file.Type, &file.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("file not found")
		}
		return nil, fmt.Errorf("failed to get file: %w", err)
	}
	
	return &file, nil
}

func (s *uploadServiceImpl) GetUserFiles(ctx context.Context, userID string) ([]models.File, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, user_id, name, url, size, type, created_at
		 FROM files WHERE user_id = ? ORDER BY created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get user files: %w", err)
	}
	defer rows.Close()
	
	var files []models.File
	for rows.Next() {
		var file models.File
		err := rows.Scan(
			&file.ID, &file.UserID, &file.Name, &file.URL, &file.Size,
			&file.Type, &file.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan file: %w", err)
		}
		files = append(files, file)
	}
	
	return files, nil
}

func (s *uploadServiceImpl) GeneratePresignedURL(ctx context.Context, fileName, fileType string) (string, error) {
	fileID := generateID("file")
	return fmt.Sprintf("https://storage.nawthtech.com/upload/%s?filename=%s&type=%s", fileID, fileName, fileType), nil
}

// NotificationService Implementation
func (s *notificationServiceImpl) CreateNotification(ctx context.Context, req NotificationCreateRequest) (*models.Notification, error) {
	notificationID := generateID("notif")
	
	notification := &models.Notification{
		ID:        notificationID,
		UserID:    req.UserID,
		Title:     req.Title,
		Message:   req.Message,
		Type:      req.Type,
		IsRead:    false,
		CreatedAt: time.Now(),
	}
	
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO notifications (id, user_id, title, message, type, is_read, created_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		notificationID, req.UserID, req.Title, req.Message, req.Type, false, time.Now(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create notification: %w", err)
	}
	
	return notification, nil
}

func (s *notificationServiceImpl) GetUserNotifications(ctx context.Context, userID string, params NotificationQueryParams) ([]models.Notification, error) {
	page, limit := validatePaginationParams(params.Page, params.Limit)
	offset := calculateOffset(page, limit)
	
	sqlQuery := `SELECT id, user_id, title, message, type, is_read, created_at
				 FROM notifications WHERE user_id = ?`
	args := []interface{}{userID}
	
	if params.Type != "" {
		sqlQuery += " AND type = ?"
		args = append(args, params.Type)
	}
	
	if params.Read != nil {
		sqlQuery += " AND is_read = ?"
		args = append(args, *params.Read)
	}
	
	sqlQuery += " ORDER BY created_at DESC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)
	
	rows, err := s.db.QueryContext(ctx, sqlQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get notifications: %w", err)
	}
	defer rows.Close()
	
	var notifications []models.Notification
	for rows.Next() {
		var notification models.Notification
		err := rows.Scan(
			&notification.ID, &notification.UserID, &notification.Title, &notification.Message,
			&notification.Type, &notification.IsRead, &notification.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan notification: %w", err)
		}
		notifications = append(notifications, notification)
	}
	
	return notifications, nil
}

func (s *notificationServiceImpl) MarkAsRead(ctx context.Context, notificationID string) error {
	_, err := s.db.ExecContext(ctx,
		"UPDATE notifications SET is_read = TRUE WHERE id = ?",
		notificationID,
	)
	return err
}

func (s *notificationServiceImpl) MarkAllAsRead(ctx context.Context, userID string) error {
	_, err := s.db.ExecContext(ctx,
		"UPDATE notifications SET is_read = TRUE WHERE user_id = ? AND is_read = FALSE",
		userID,
	)
	return err
}

func (s *notificationServiceImpl) DeleteNotification(ctx context.Context, notificationID string) error {
	_, err := s.db.ExecContext(ctx,
		"DELETE FROM notifications WHERE id = ?",
		notificationID,
	)
	return err
}

func (s *notificationServiceImpl) GetUnreadCount(ctx context.Context, userID string) (int64, error) {
	var count int64
	err := s.db.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM notifications WHERE user_id = ? AND is_read = FALSE",
		userID,
	).Scan(&count)
	return count, err
}

// AdminService Implementation
func (s *adminServiceImpl) GetDashboardStats(ctx context.Context) (*DashboardStats, error) {
	stats := &DashboardStats{}
	
	// إجمالي المستخدمين
	err := s.db.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM users WHERE status = 'active'",
	).Scan(&stats.TotalUsers)
	if err != nil {
		return nil, fmt.Errorf("failed to get total users: %w", err)
	}
	
	// إجمالي الخدمات
	err = s.db.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM services WHERE is_active = TRUE",
	).Scan(&stats.TotalServices)
	if err != nil {
		return nil, fmt.Errorf("failed to get total services: %w", err)
	}
	
	// إجمالي الطلبات
	err = s.db.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM orders",
	).Scan(&stats.TotalOrders)
	if err != nil {
		return nil, fmt.Errorf("failed to get total orders: %w", err)
	}
	
	// إجمالي الإيرادات
	err = s.db.QueryRowContext(ctx,
		"SELECT COALESCE(SUM(amount), 0) FROM orders WHERE status = 'completed'",
	).Scan(&stats.TotalRevenue)
	if err != nil {
		return nil, fmt.Errorf("failed to get total revenue: %w", err)
	}
	
	// المستخدمين النشطين (سجلوا دخول في آخر 30 يوم)
	err = s.db.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM users WHERE last_login >= DATE('now', '-30 days')",
	).Scan(&stats.ActiveUsers)
	if err != nil {
		return nil, fmt.Errorf("failed to get active users: %w", err)
	}
	
	// الطلبات المعلقة
	err = s.db.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM orders WHERE status = 'pending'",
	).Scan(&stats.PendingOrders)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending orders: %w", err)
	}
	
	// الطلبات المكتملة
	err = s.db.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM orders WHERE status = 'completed'",
	).Scan(&stats.CompletedOrders)
	if err != nil {
		return nil, fmt.Errorf("failed to get completed orders: %w", err)
	}
	
	return stats, nil
}

func (s *adminServiceImpl) GetUsers(ctx context.Context, params UserQueryParams) ([]models.User, error) {
	userService := NewUserService(s.db)
	return userService.SearchUsers(ctx, "", params)
}

func (s *adminServiceImpl) GetSystemLogs(ctx context.Context, params SystemLogQuery) ([]models.SystemLog, error) {
	page, limit := validatePaginationParams(params.Page, params.Limit)
	offset := calculateOffset(page, limit)
	
	sqlQuery := `SELECT id, user_id, level, action, resource, details, ip_address, user_agent, created_at
				 FROM system_logs WHERE 1=1`
	args := []interface{}{}
	
	if params.Level != "" {
		sqlQuery += " AND level = ?"
		args = append(args, params.Level)
	}
	
	if params.UserID != "" {
		sqlQuery += " AND user_id = ?"
		args = append(args, params.UserID)
	}
	
	if !params.FromDate.IsZero() {
		sqlQuery += " AND created_at >= ?"
		args = append(args, params.FromDate)
	}
	
	if !params.ToDate.IsZero() {
		sqlQuery += " AND created_at <= ?"
		args = append(args, params.ToDate)
	}
	
	sqlQuery += " ORDER BY created_at DESC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)
	
	rows, err := s.db.QueryContext(ctx, sqlQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get system logs: %w", err)
	}
	defer rows.Close()
	
	var logs []models.SystemLog
	for rows.Next() {
		var log models.SystemLog
		err := rows.Scan(
			&log.ID, &log.UserID, &log.Level, &log.Action, &log.Resource,
			&log.Details, &log.IPAddress, &log.UserAgent, &log.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan system log: %w", err)
		}
		logs = append(logs, log)
	}
	
	return logs, nil
}

func (s *adminServiceImpl) UpdateSystemSettings(ctx context.Context, settings map[string]string) error {
	// هنا يمكنك تحديث الإعدادات في قاعدة بيانات أو ملف
	return nil
}

func (s *adminServiceImpl) BanUser(ctx context.Context, userID string, reason string) error {
	_, err := s.db.ExecContext(ctx,
		"UPDATE users SET status = 'banned', updated_at = ? WHERE id = ?",
		time.Now(), userID,
	)
	return err
}

func (s *adminServiceImpl) UnbanUser(ctx context.Context, userID string) error {
	_, err := s.db.ExecContext(ctx,
		"UPDATE users SET status = 'active', updated_at = ? WHERE id = ?",
		time.Now(), userID,
	)
	return err
}

// CacheService Implementation
func (c *cacheServiceImpl) Get(key string) (interface{}, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	val, ok := c.store[key]
	if !ok {
		return nil, errors.New("key not found")
	}
	return val, nil
}

func (c *cacheServiceImpl) Set(key string, value interface{}, expiration time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	c.store[key] = value
	return nil
}

func (c *cacheServiceImpl) Delete(key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	delete(c.store, key)
	return nil
}

func (c *cacheServiceImpl) Exists(key string) (bool, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	_, ok := c.store[key]
	return ok, nil
}

func (c *cacheServiceImpl) Flush() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	c.store = make(map[string]interface{})
	return nil
}

// ================================
// دوال ServiceContainer
// ================================

func (sc *ServiceContainer) InitializeDatabase(ctx context.Context) error {
	for _, query := range CreateTablesSQL() {
		_, err := sc.db.ExecContext(ctx, query)
		if err != nil {
			return fmt.Errorf("failed to execute query: %s, error: %w", query, err)
		}
	}
	return nil
}

func (sc *ServiceContainer) Close() error {
	var errors []string
	
	if sc.db != nil {
		if err := sc.db.Close(); err != nil {
			errors = append(errors, fmt.Sprintf("database: %v", err))
		}
	}
	
	if sc.logger != nil {
		if err := sc.logger.Sync(); err != nil {
			errors = append(errors, fmt.Sprintf("logger: %v", err))
		}
	}
	
	if len(errors) > 0 {
		return fmt.Errorf("errors closing service container: %s", strings.Join(errors, "; "))
	}
	
	return nil
}

// ================================
// تعريفات الأخطاء (Errors)
// ================================

var (
	ErrServiceNotFound   = errors.New("service not found")
	ErrUserNotFound      = errors.New("user not found")
	ErrCategoryNotFound  = errors.New("category not found")
	ErrOrderNotFound     = errors.New("order not found")
	ErrPaymentNotFound   = errors.New("payment not found")
	ErrFileNotFound      = errors.New("file not found")
	ErrNotificationNotFound = errors.New("notification not found")
	ErrInvalidRequest    = errors.New("invalid request")
	ErrUnauthorized      = errors.New("unauthorized")
	ErrInsufficientFunds = errors.New("insufficient funds")
	ErrDuplicateEntry    = errors.New("duplicate entry")
	ErrDatabase          = errors.New("database error")
	ErrValidation        = errors.New("validation error")
	ErrNotImplemented    = errors.New("not implemented")
)

// ================================
// دوال Init و Initialization
// ================================

func InitServiceContainer(db *sql.DB) *ServiceContainer {
	container := NewServiceContainer(db)
	
	// تهيئة قاعدة البيانات
	ctx := context.Background()
	if err := container.InitializeDatabase(ctx); err != nil {
		fmt.Printf("Warning: failed to initialize database: %v\n", err)
	}
	
	return container
}

func InitServiceContainerWithConfig(db *sql.DB, cfg *config.Config, logger *zap.Logger) *ServiceContainer {
	container := NewServiceContainerWithConfig(db, cfg, logger)
	
	// تهيئة قاعدة البيانات
	ctx := context.Background()
	if err := container.InitializeDatabase(ctx); err != nil {
		logger.Warn("Failed to initialize database", zap.Error(err))
	}
	
	return container
}