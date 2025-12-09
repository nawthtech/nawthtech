package services

import (
	"context"
	"database/sql"
	"fmt"
 "errors"
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
)

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

	paymentServiceImpl struct {
		db *sql.DB
	}

	uploadServiceImpl struct {
		db *sql.DB
	}

	notificationServiceImpl struct {
		db *sql.DB
	}

	adminServiceImpl struct {
		db *sql.DB
	}

	cacheServiceImpl struct {
		store map[string]interface{}
	}
)

// ================================
// أمثلة على تنفيذ AuthService باستخدام D1
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

// ================================
// دوال الإنشاء باستخدام D1
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
	return &cacheServiceImpl{store: make(map[string]interface{})}
}

// ================================
// Service Container مع D1
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

func NewServiceContainer(d1db *sql.DB) *ServiceContainer {
	return &ServiceContainer{
		Auth:         NewAuthService(d1db),
		User:         NewUserService(d1db),
		Service:      NewServiceService(d1db),
		Category:     NewCategoryService(d1db),
		Order:        NewOrderService(d1db),
		Payment:      NewPaymentService(d1db),
		Upload:       NewUploadService(d1db),
		Notification: NewNotificationService(d1db),
		Admin:        NewAdminService(d1db),
		Cache:        NewCacheService(),
	}
}

package services

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/nawthtech/nawthtech/backend/internal/models"
)

// ================================
// تطبيقات AuthService مع D1
// ================================

type authServiceImpl struct {
	db *sql.DB
}

func (s *authServiceImpl) Register(ctx context.Context, req AuthRegisterRequest) (*AuthResponse, error) {
	userID := fmt.Sprintf("user_%d", time.Now().UnixNano())
	user := &models.User{
		ID:            userID,
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

	_, err := s.db.ExecContext(ctx, `
		INSERT INTO users (id,email,username,password,first_name,last_name,phone,role,status,email_verified,created_at,updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		user.ID, user.Email, user.Username, user.Password, user.FirstName, user.LastName, user.Phone,
		user.Role, user.Status, user.EmailVerified, user.CreatedAt, user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		User:         user,
		AccessToken:  "access_token_" + userID,
		RefreshToken: "refresh_token_" + userID,
		ExpiresAt:    time.Now().Add(24 * time.Hour),
	}, nil
}

func (s *authServiceImpl) Login(ctx context.Context, req AuthLoginRequest) (*AuthResponse, error) {
	row := s.db.QueryRowContext(ctx, "SELECT id,email,username,role,status,email_verified FROM users WHERE email = ?", req.Email)

	user := &models.User{}
	err := row.Scan(&user.ID, &user.Email, &user.Username, &user.Role, &user.Status, &user.EmailVerified)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("المستخدم غير موجود")
		}
		return nil, err
	}

	return &AuthResponse{
		User:         user,
		AccessToken:  "access_token_" + user.ID,
		RefreshToken: "refresh_token_" + user.ID,
		ExpiresAt:    time.Now().Add(24 * time.Hour),
	}, nil
}

// ================================
// تطبيقات UserService مع D1
// ================================

type userServiceImpl struct {
	db *sql.DB
}

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
	_, err := s.db.ExecContext(ctx, `
		UPDATE users SET first_name=?, last_name=?, phone=?, avatar=?, updated_at=?
		WHERE id=?`,
		req.FirstName, req.LastName, req.Phone, req.Avatar, time.Now(), userID,
	)
	if err != nil {
		return nil, err
	}
	return s.GetProfile(ctx, userID)
}

func (s *userServiceImpl) DeleteAccount(ctx context.Context, userID string) error {
	_, err := s.db.ExecContext(ctx, "DELETE FROM users WHERE id=?", userID)
	return err
}

// ================================
// تطبيقات ServiceService مع D1
// ================================

type serviceServiceImpl struct {
	db *sql.DB
}

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

	_, err := s.db.ExecContext(ctx, `
		INSERT INTO services (id,title,description,price,duration,category_id,provider_id,images,tags,is_active,is_featured,rating,created_at,updated_at)
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

// ================================
// دوال مساعدة لتحويل slice إلى نص لتخزينه في D1
// ================================

func serializeStrings(arr []string) string {
	result := ""
	for i, s := range arr {
		if i > 0 {
			result += ","
		}
		result += s
	}
	return result
}

func deserializeStrings(s string) []string {
	if s == "" {
		return []string{}
	}
	return split(s, ",")
}

func split(s string, sep string) []string {
	var result []string
	for _, v := range []rune(s) {
		result = append(result, string(v))
	}
	return result
}

// ================================
// دوال إنشاء الخدمات الجديدة مع D1
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

// ================================
// OrderService مع D1
// ================================

type orderServiceImpl struct {
	db *sql.DB
}

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

	_, err := s.db.ExecContext(ctx, `
		INSERT INTO orders (id,user_id,service_id,status,amount,created_at,updated_at)
		VALUES (?,?,?,?,?,?,?)`,
		order.ID, order.UserID, order.ServiceID, order.Status, order.Amount, order.CreatedAt, order.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (s *orderServiceImpl) GetOrder(ctx context.Context, orderID string) (*models.Order, error) {
	row := s.db.QueryRowContext(ctx, "SELECT id,user_id,service_id,status,amount,created_at,updated_at FROM orders WHERE id=?", orderID)
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

func (s *orderServiceImpl) UpdateOrderStatus(ctx context.Context, orderID string, status string) error {
	_, err := s.db.ExecContext(ctx, "UPDATE orders SET status=?, updated_at=? WHERE id=?", status, time.Now(), orderID)
	return err
}

func NewOrderService(db *sql.DB) OrderService {
	return &orderServiceImpl{db: db}
}

// ================================
// CategoryService مع D1
// ================================

type categoryServiceImpl struct {
	db *sql.DB
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

	_, err := s.db.ExecContext(ctx, `
		INSERT INTO categories (id,name,slug,image,is_active,created_at,updated_at)
		VALUES (?,?,?,?,?,?,?)`,
		category.ID, category.Name, category.Slug, category.Image, category.IsActive, category.CreatedAt, category.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return category, nil
}

func (s *categoryServiceImpl) ListCategories(ctx context.Context) ([]*models.Category, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT id,name,slug,image,is_active,created_at,updated_at FROM categories WHERE is_active=1")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*models.Category
	for rows.Next() {
		c := &models.Category{}
		err := rows.Scan(&c.ID, &c.Name, &c.Slug, &c.Image, &c.IsActive, &c.CreatedAt, &c.UpdatedAt)
		if err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, nil
}

func NewCategoryService(db *sql.DB) CategoryService {
	return &categoryServiceImpl{db: db}
}

// ================================
// PaymentService مع D1 (مثال أساسي)
// ================================

type paymentServiceImpl struct {
	db *sql.DB
}

func (s *paymentServiceImpl) CreatePayment(ctx context.Context, req PaymentCreateRequest) (*models.Payment, error) {
	paymentID := fmt.Sprintf("payment_%d", time.Now().UnixNano())
	payment := &models.Payment{
		ID:        paymentID,
		OrderID:   req.OrderID,
		Amount:    req.Amount,
		Status:    "pending",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err := s.db.ExecContext(ctx, `
		INSERT INTO payments (id,order_id,amount,status,created_at,updated_at)
		VALUES (?,?,?,?,?,?)`,
		payment.ID, payment.OrderID, payment.Amount, payment.Status, payment.CreatedAt, payment.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return payment, nil
}

func (s *paymentServiceImpl) UpdatePaymentStatus(ctx context.Context, paymentID string, status string) error {
	_, err := s.db.ExecContext(ctx, "UPDATE payments SET status=?, updated_at=? WHERE id=?", status, time.Now(), paymentID)
	return err
}

func NewPaymentService(db *sql.DB) PaymentService {
	return &paymentServiceImpl{db: db}
}

// ================================
// UploadService (تخزين الملفات) - مجرد مثال، يعتمد على رابط URL
// ================================

type uploadServiceImpl struct{}

func (s *uploadServiceImpl) UploadFile(ctx context.Context, fileName string, fileData []byte) (string, error) {
	// هنا يمكن دمج مع S3 أو Cloudflare R2، الآن مجرد مثال
	url := fmt.Sprintf("https://cdn.nawthtech.com/%s", fileName)
	return url, nil
}

func NewUploadService() UploadService {
	return &uploadServiceImpl{}
}

// ================================
// NotificationService - مجرد مثال
// ================================

type notificationServiceImpl struct{}

func (s *notificationServiceImpl) SendNotification(ctx context.Context, userID string, message string) error {
	fmt.Printf("Notification to %s: %s\n", userID, message)
	return nil
}

func NewNotificationService() NotificationService {
	return &notificationServiceImpl{}
}

// ================================
// AdminService - إدارة المستخدمين والخدمات
// ================================

type adminServiceImpl struct {
	db *sql.DB
}

func (s *adminServiceImpl) DeactivateUser(ctx context.Context, userID string) error {
	_, err := s.db.ExecContext(ctx, "UPDATE users SET status=?, updated_at=? WHERE id=?", "inactive", time.Now(), userID)
	return err
}

func (s *adminServiceImpl) DeactivateService(ctx context.Context, serviceID string) error {
	_, err := s.db.ExecContext(ctx, "UPDATE services SET is_active=?, updated_at=? WHERE id=?", false, time.Now(), serviceID)
	return err
}

func NewAdminService(db *sql.DB) AdminService {
	return &adminServiceImpl{db: db}
}

// ================================
// CacheService - مجرد مثال باستخدام map محلي
// ================================

type cacheServiceImpl struct {
	store map[string]interface{}
}

func (c *cacheServiceImpl) Set(key string, value interface{}) {
	c.store[key] = value
}

func (c *cacheServiceImpl) Get(key string) (interface{}, bool) {
	val, ok := c.store[key]
	return val, ok
}

func NewCacheService() CacheService {
	return &cacheServiceImpl{store: make(map[string]interface{})}
}

// Example: AuthService -------------------------------------

type AuthService interface {
	Register(ctx context.Context, email, password string) (string, error)
	Login(ctx context.Context, email, password string) (string, error)
	RefreshToken(ctx context.Context, token string) (string, error)
}

type authService struct{}

func NewAuthService() AuthService {
	return &authService{}
}

func (s *authService) Register(ctx context.Context, email, password string) (string, error) {
	if email == "" || password == "" {
		return "", errors.New("invalid email or password")
	}
	return "user_created_successfully", nil
}

func (s *authService) Login(ctx context.Context, email, password string) (string, error) {
	if email == "" || password == "" {
		return "", errors.New("invalid credentials")
	}
	return "jwt_token_here", nil
}

func (s *authService) RefreshToken(ctx context.Context, token string) (string, error) {
	if token == "" {
		return "", errors.New("invalid token")
	}
	return "refreshed_jwt_token", nil
}

// Example: UserService -------------------------------------

type UserService interface {
	GetProfile(ctx context.Context, userID string) (*UserProfile, error)
	UpdateProfile(ctx context.Context, userID string, data UpdateProfileDTO) error
}

type userService struct{}

func NewUserService() UserService {
	return &userService{}
}

type UserProfile struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type UpdateProfileDTO struct {
	Name string `json:"name"`
}

func (s *userService) GetProfile(ctx context.Context, userID string) (*UserProfile, error) {
	if userID == "" {
		return nil, errors.New("missing user ID")
	}
	return &UserProfile{
		ID:        userID,
		Email:     "example@mail.com",
		Name:      "Test User",
		CreatedAt: time.Now().Add(-24 * time.Hour),
	}, nil
}

func (s *userService) UpdateProfile(ctx context.Context, userID string, data UpdateProfileDTO) error {
	if userID == "" {
		return errors.New("missing user ID")
	}
	return nil
}

// Example: ProductService -------------------------------------

type ProductService interface {
	GetAll(ctx context.Context) ([]Product, error)
	GetByID(ctx context.Context, id string) (*Product, error)
	Create(ctx context.Context, p ProductDTO) (string, error)
	Delete(ctx context.Context, id string) error
}

type productService struct{}

func NewProductService() ProductService {
	return &productService{}
}

type Product struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

type ProductDTO struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

func (s *productService) GetAll(ctx context.Context) ([]Product, error) {
	return []Product{
		{ID: "1", Name: "Item A", Price: 20.5},
		{ID: "2", Name: "Item B", Price: 14.0},
	}, nil
}

func (s *productService) GetByID(ctx context.Context, id string) (*Product, error) {
	if id == "" {
		return nil, errors.New("missing product ID")
	}
	return &Product{
		ID:    id,
		Name:  "Example Item",
		Price: 33.0,
	}, nil
}

func (s *productService) Create(ctx context.Context, p ProductDTO) (string, error) {
	if p.Name == "" || p.Price <= 0 {
		return "", errors.New("invalid product data")
	}
	return "new-product-id", nil
}

func (s *productService) Delete(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("missing product ID")
	}
	return nil
}

// Example: OrderService -------------------------------------

type OrderService interface {
	CreateOrder(ctx context.Context, userID string, items []OrderItemDTO) (string, error)
	GetOrdersByUser(ctx context.Context, userID string) ([]Order, error)
}

type orderService struct{}

func NewOrderService() OrderService {
	return &orderService{}
}

type Order struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Items     []OrderItemDTO `json:"items"`
	CreatedAt time.Time `json:"created_at"`
}

type OrderItemDTO struct {
	ProductID string  `json:"product_id"`
	Qty       int     `json:"qty"`
	UnitPrice float64 `json:"unit_price"`
}

func (s *orderService) CreateOrder(ctx context.Context, userID string, items []OrderItemDTO) (string, error) {
	if userID == "" {
		return "", errors.New("missing user ID")
	}
	if len(items) == 0 {
		return "", errors.New("empty order")
	}
	return "order-id-12345", nil
}

func (s *orderService) GetOrdersByUser(ctx context.Context, userID string) ([]Order, error) {
	if userID == "" {
		return nil, errors.New("missing user ID")
	}
	return []Order{
		{ID: "O-1", UserID: userID, Items: []OrderItemDTO{}, CreatedAt: time.Now().Add(-10 * time.Hour)},
	}, nil
}