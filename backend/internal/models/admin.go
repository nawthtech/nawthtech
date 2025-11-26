package models

import "time"

// DashboardStats إحصائيات لوحة التحكم
type DashboardStats struct {
	TotalUsers        int     `json:"totalUsers"`
	TotalOrders       int     `json:"totalOrders"`
	TotalRevenue      float64 `json:"totalRevenue"`
	ActiveServices    int     `json:"activeServices"`
	PendingOrders     int     `json:"pendingOrders"`
	SupportTickets    int     `json:"supportTickets"`
	ConversionRate    float64 `json:"conversionRate"`
	BounceRate        float64 `json:"bounceRate"`
	StoreVisits       int     `json:"storeVisits"`
	NewCustomers      int     `json:"newCustomers"`
}

// StoreMetrics مقاييس المتجر
type StoreMetrics struct {
	TotalProducts           int     `json:"totalProducts"`
	LowStockItems          int     `json:"lowStockItems"`
	StoreRevenue           float64 `json:"storeRevenue"`
	StoreOrders            int     `json:"storeOrders"`
	AverageOrderValue      float64 `json:"averageOrderValue"`
	TopSellingCategory     string  `json:"topSellingCategory"`
	CustomerSatisfaction   float64 `json:"customerSatisfaction"`
	ReturnRate             float64 `json:"returnRate"`
}

// Order طلب
type Order struct {
	ID          string    `json:"id"`
	User        string    `json:"user"`
	Service     string    `json:"service"`
	Amount      float64   `json:"amount"`
	Status      string    `json:"status"`
	Date        string    `json:"date"`
	Type        string    `json:"type"`
	Category    string    `json:"category"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// UserActivity نشاط المستخدم
type UserActivity struct {
	User     string    `json:"user"`
	Action   string    `json:"action"`
	Service  string    `json:"service,omitempty"`
	Time     string    `json:"time"`
	IP       string    `json:"ip"`
	Type     string    `json:"type"`
	CreatedAt time.Time `json:"createdAt"`
}

// SystemAlert تنبيه النظام
type SystemAlert struct {
	Type      string `json:"type"`
	Title     string `json:"title"`
	Message   string `json:"message"`
	Priority  string `json:"priority"` // low, medium, high, critical
	Action    string `json:"action,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

// DashboardData بيانات لوحة التحكم
type DashboardData struct {
	Stats        DashboardStats `json:"stats"`
	StoreMetrics StoreMetrics   `json:"storeMetrics"`
	RecentOrders []Order        `json:"recentOrders"`
	UserActivity []UserActivity `json:"userActivity"`
	SystemAlerts []SystemAlert  `json:"systemAlerts"`
}

// SalesPerformance أداء المبيعات
type SalesPerformance struct {
	Date  string  `json:"date"`
	Sales float64 `json:"sales"`
	Orders int    `json:"orders"`
}

// PerformanceMetric مقياس الأداء
type PerformanceMetric struct {
	Value  float64 `json:"value"`
	Label  string  `json:"label"`
	Change float64 `json:"change"`
}
