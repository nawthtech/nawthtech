package models

// FinancialReport التقرير المالي
type FinancialReport struct {
	Timeframe string               `json:"timeframe"`
	Revenue   float64              `json:"revenue"`
	Expenses  float64              `json:"expenses"`
	Profit    float64              `json:"profit"`
	Growth    string               `json:"growth"`
	Breakdown []RevenueBreakdown   `json:"breakdown"`
}

// RevenueBreakdown تفصيل الإيرادات
type RevenueBreakdown struct {
	Category   string  `json:"category"`
	Amount     float64 `json:"amount"`
	Percentage float64 `json:"percentage"`
}

// PlatformAnalytics تحليلات المنصة
type PlatformAnalytics struct {
	TotalUsers       int     `json:"total_users"`
	ActiveUsers      int     `json:"active_users"`
	TotalServices    int     `json:"total_services"`
	ActiveServices   int     `json:"active_services"`
	TotalOrders      int     `json:"total_orders"`
	CompletedOrders  int     `json:"completed_orders"`
	TotalRevenue     float64 `json:"total_revenue"`
	AverageRating    float64 `json:"average_rating"`
	GrowthRate       string  `json:"growth_rate"`
	RetentionRate    string  `json:"retention_rate"`
	ConversionRate   string  `json:"conversion_rate"`
}