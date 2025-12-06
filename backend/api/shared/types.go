package shared

import "time"

// BaseModel النموذج الأساسي لجميع الكيانات
type BaseModel struct {
	ID        string     `json:"id" bson:"_id"`
	CreatedAt time.Time  `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" bson:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" bson:"deleted_at,omitempty"`
}

// User المستخدم
type User struct {
	BaseModel
	FirstName string    `json:"first_name" bson:"first_name"`
	LastName  string    `json:"last_name" bson:"last_name"`
	Email     string    `json:"email" bson:"email"`
	Phone     string    `json:"phone" bson:"phone"`
	Avatar    string    `json:"avatar,omitempty" bson:"avatar,omitempty"`
	Role      string    `json:"role" bson:"role"`     // user, admin
	Status    string    `json:"status" bson:"status"` // active, suspended
	LastLogin time.Time `json:"last_login,omitempty" bson:"last_login,omitempty"`
}

// Service الخدمة
type Service struct {
	BaseModel
	Title       string   `json:"title" bson:"title"`
	Description string   `json:"description" bson:"description"`
	Price       float64  `json:"price" bson:"price"`
	CategoryID  string   `json:"category_id" bson:"category_id"`
	ProviderID  string   `json:"provider_id" bson:"provider_id"`
	Tags        []string `json:"tags,omitempty" bson:"tags,omitempty"`
	Images      []string `json:"images,omitempty" bson:"images,omitempty"`
	Status      string   `json:"status" bson:"status"` // active, inactive, pending
	Rating      float64  `json:"rating" bson:"rating"`
	ReviewCount int      `json:"review_count" bson:"review_count"`
}

// Category الفئة
type Category struct {
	BaseModel
	Name         string `json:"name" bson:"name"`
	Description  string `json:"description,omitempty" bson:"description,omitempty"`
	Icon         string `json:"icon,omitempty" bson:"icon,omitempty"`
	Color        string `json:"color,omitempty" bson:"color,omitempty"`
	ServiceCount int    `json:"service_count" bson:"service_count"`
}

// Order الطلب
type Order struct {
	BaseModel
	UserID    string    `json:"user_id" bson:"user_id"`
	ServiceID string    `json:"service_id" bson:"service_id"`
	Quantity  int       `json:"quantity" bson:"quantity"`
	Total     float64   `json:"total" bson:"total"`
	Status    string    `json:"status" bson:"status"` // pending, confirmed, completed, cancelled
	Notes     string    `json:"notes,omitempty" bson:"notes,omitempty"`
	DueDate   time.Time `json:"due_date,omitempty" bson:"due_date,omitempty"`
}

// Payment الدفع
type Payment struct {
	BaseModel
	OrderID       string     `json:"order_id" bson:"order_id"`
	UserID        string     `json:"user_id" bson:"user_id"`
	Amount        float64    `json:"amount" bson:"amount"`
	Currency      string     `json:"currency" bson:"currency"`
	Status        string     `json:"status" bson:"status"` // pending, completed, failed, refunded
	PaymentMethod string     `json:"payment_method" bson:"payment_method"`
	TransactionID string     `json:"transaction_id,omitempty" bson:"transaction_id,omitempty"`
	PaidAt        *time.Time `json:"paid_at,omitempty" bson:"paid_at,omitempty"`
}

// Notification الإشعار
type Notification struct {
	BaseModel
	UserID  string                 `json:"user_id" bson:"user_id"`
	Title   string                 `json:"title" bson:"title"`
	Message string                 `json:"message" bson:"message"`
	Type    string                 `json:"type" bson:"type"` // info, success, warning, error
	Read    bool                   `json:"read" bson:"read"`
	Data    map[string]interface{} `json:"data,omitempty" bson:"data,omitempty"`
}

// UploadedFile الملف المرفوع
type UploadedFile struct {
	BaseModel
	UserID       string `json:"user_id" bson:"user_id"`
	PublicID     string `json:"public_id" bson:"public_id"`
	SecureURL    string `json:"secure_url" bson:"secure_url"`
	Format       string `json:"format" bson:"format"`
	Bytes        int    `json:"bytes" bson:"bytes"`
	Width        int    `json:"width,omitempty" bson:"width,omitempty"`
	Height       int    `json:"height,omitempty" bson:"height,omitempty"`
	ResourceType string `json:"resource_type" bson:"resource_type"`
	Folder       string `json:"folder" bson:"folder"`
}

// SystemStats إحصائيات النظام
type SystemStats struct {
	TotalUsers      int64   `json:"total_users"`
	TotalServices   int64   `json:"total_services"`
	TotalOrders     int64   `json:"total_orders"`
	TotalRevenue    float64 `json:"total_revenue"`
	ActiveUsers     int64   `json:"active_users"`
	PendingOrders   int64   `json:"pending_orders"`
	CompletedOrders int64   `json:"completed_orders"`
}

// APIResponse استجابة API قياسية
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Code    int         `json:"code,omitempty"`
}

// PaginationRequest طلب الصفحة
type PaginationRequest struct {
	Page    int `json:"page" form:"page" query:"page"`
	PerPage int `json:"per_page" form:"per_page" query:"per_page"`
}

// PaginationResponse استجابة الصفحة
type PaginationResponse struct {
	Page       int         `json:"page"`
	PerPage    int         `json:"per_page"`
	Total      int64       `json:"total"`
	TotalPages int         `json:"total_pages"`
	Data       interface{} `json:"data"`
}