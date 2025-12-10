package models

import "time"

// ================================
// User
// ================================

type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Password  string    `json:"-"`
	Status    string    `json:"status"` // active / inactive
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ================================
// Service
// ================================

type Service struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	CategoryID  string    `json:"category_id"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ================================
// Category
// ================================

type Category struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	Image     string    `json:"image"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ================================
// Order
// ================================

type Order struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	ServiceID string    `json:"service_id"`
	Status    string    `json:"status"` // pending, completed, canceled
	Amount    float64   `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ================================
// Payment
// ================================

type Payment struct {
	ID        string    `json:"id"`
	OrderID   string    `json:"order_id"`
	Amount    float64   `json:"amount"`
	Status    string    `json:"status"` // pending, completed, failed
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
