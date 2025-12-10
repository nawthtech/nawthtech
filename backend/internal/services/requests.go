package services

// ================================
// User Requests
// ================================

type UserCreateRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type UserUpdateRequest struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

// ================================
// Service Requests
// ================================

type ServiceCreateRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	CategoryID  string  `json:"category_id"`
}

type ServiceUpdateRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	IsActive    bool    `json:"is_active"`
}

// ================================
// Category Requests
// ================================

type CategoryCreateRequest struct {
	Name  string `json:"name"`
	Slug  string `json:"slug"`
	Image string `json:"image"`
}

// ================================
// Order Requests
// ================================

type OrderCreateRequest struct {
	UserID    string  `json:"user_id"`
	ServiceID string  `json:"service_id"`
	Amount    float64 `json:"amount"`
}

// ================================
// Payment Requests
// ================================

type PaymentCreateRequest struct {
	OrderID string  `json:"order_id"`
	Amount  float64 `json:"amount"`
}
