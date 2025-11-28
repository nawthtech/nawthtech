package models

import "time"

// ================================
// النماذج المشتركة بين الخدمات
// ================================

// PerformanceMetric مقياس الأداء (نموذج مشترك)
type PerformanceMetric struct {
	Value  float64 `json:"value"`
	Label  string  `json:"label"`
	Change float64 `json:"change"`
}

// Pagination الترقيم الصفحي
type Pagination struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
	Total int `json:"total"`
	Pages int `json:"pages"`
}

// DateRange نطاق التاريخ
type DateRange struct {
	From time.Time `json:"from"`
	To   time.Time `json:"to"`
}

// Filter عامل التصفية
type Filter struct {
	Field    string      `json:"field"`
	Operator string      `json:"operator"`
	Value    interface{} `json:"value"`
}

// Order طلب (نموذج مشترك)
type Order struct {
	ID          string    `json:"id"`
	User        string    `json:"user"`
	Service     string    `json:"service"`
	Amount      float64   `json:"amount"`
	Status      string    `json:"status"`
	Date        string    `json:"date"`
	Type        string    `json:"type"`
	Category    string    `json:"category"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// UserActivity نشاط المستخدم (نموذج مشترك)
type UserActivity struct {
	User      string    `json:"user"`
	Action    string    `json:"action"`
	Service   string    `json:"service,omitempty"`
	Time      string    `json:"time"`
	IP        string    `json:"ip"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
}

// SystemAlert تنبيه النظام (نموذج مشترك)
type SystemAlert struct {
	Type      string    `json:"type"`
	Title     string    `json:"title"`
	Message   string    `json:"message"`
	Priority  string    `json:"priority"` // low, medium, high, critical
	Action    string    `json:"action,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

// LogEntry إدخال السجل (نموذج مشترك)
type LogEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Level     string    `json:"level"`
	Message   string    `json:"message"`
	Service   string    `json:"service"`
	UserID    string    `json:"user_id,omitempty"`
	RequestID string    `json:"request_id,omitempty"`
}

// HealthCheckResult نتيجة فحص الصحة
type HealthCheckResult struct {
	Service     string      `json:"service"`
	Status      string      `json:"status"`
	ResponseTime string     `json:"response_time,omitempty"`
	Error       string      `json:"error,omitempty"`
	Usage       string      `json:"usage,omitempty"`
	Details     interface{} `json:"details,omitempty"`
}