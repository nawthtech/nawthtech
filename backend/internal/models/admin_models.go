package models

import (
	"time"
)

// ================================
// لوحة التحكم والإحصائيات
// ================================

// DashboardStats إحصائيات لوحة التحكم
type DashboardStats struct {
	TotalUsers        int     `json:"total_users"`
	TotalOrders       int     `json:"total_orders"`
	TotalRevenue      float64 `json:"total_revenue"`
	ActiveServices    int     `json:"active_services"`
	PendingOrders     int     `json:"pending_orders"`
	SupportTickets    int     `json:"support_tickets"`
	ConversionRate    float64 `json:"conversion_rate"`
	BounceRate        float64 `json:"bounce_rate"`
	StoreVisits       int     `json:"store_visits"`
	NewCustomers      int     `json:"new_customers"`
}

// StoreMetrics مقاييس المتجر
type StoreMetrics struct {
	TotalProducts          int     `json:"total_products"`
	LowStockItems          int     `json:"low_stock_items"`
	StoreRevenue           float64 `json:"store_revenue"`
	StoreOrders            int     `json:"store_orders"`
	AverageOrderValue      float64 `json:"average_order_value"`
	TopSellingCategory     string  `json:"top_selling_category"`
	CustomerSatisfaction   float64 `json:"customer_satisfaction"`
	ReturnRate             float64 `json:"return_rate"`
}

// SalesPerformance أداء المبيعات
type SalesPerformance struct {
	Date   string  `json:"date"`
	Sales  float64 `json:"sales"`
	Orders int     `json:"orders"`
}

// DashboardData بيانات لوحة التحكم
type DashboardData struct {
	Stats         DashboardStats      `json:"stats"`
	StoreMetrics  StoreMetrics        `json:"store_metrics"`
	RecentOrders  []Order             `json:"recent_orders"`  // من shared_models
	UserActivity  []UserActivity      `json:"user_activity"`  // من shared_models
	SystemAlerts  []SystemAlert       `json:"system_alerts"`  // من shared_models
	SalesTrends   []SalesPerformance  `json:"sales_trends"`
	Performance   []PerformanceMetric `json:"performance"`    // من shared_models
}

// ================================
// النظام والصيانة
// ================================

// SystemStatus حالة النظام
type SystemStatus struct {
	Disk        DiskStatus        `json:"disk"`
	Database    DatabaseStatus    `json:"database"`
	System      SystemInfo        `json:"system"`
	Performance PerformanceMetrics `json:"performance"`
	Services    map[string]string `json:"services"`
	AIAnalysis  AIHealthAnalysis  `json:"ai_analysis"`
	Security    SecurityStatus    `json:"security"`
	LastChecked time.Time         `json:"last_checked"`
}

// SystemInfo معلومات النظام
type SystemInfo struct {
	Version      string     `json:"version"`
	NodeVersion  string     `json:"node_version"`
	Environment  string     `json:"environment"`
	Platform     string     `json:"platform"`
	Arch         string     `json:"arch"`
	Uptime       float64    `json:"uptime"`
	Memory       MemoryUsage `json:"memory"`
	CPU          CPUUsage   `json:"cpu"`
	PID          int        `json:"pid"`
}

// DiskStatus حالة القرص
type DiskStatus struct {
	Free            uint64   `json:"free"`
	Size            uint64   `json:"size"`
	Used            uint64   `json:"used"`
	FreePercentage  float64  `json:"free_percentage"`
	Path            string   `json:"path"`
	Threshold       string   `json:"threshold"`
	Recommendations []string `json:"recommendations"`
}

// DatabaseStatus حالة قاعدة البيانات
type DatabaseStatus struct {
	Connected     bool                `json:"connected"`
	ReadyState    string              `json:"ready_state"`
	DBName        string              `json:"db_name"`
	Host          string              `json:"host"`
	Connections   int                 `json:"connections"`
	Performance   DatabasePerformance `json:"performance"`
	Collections   int                 `json:"collections"`
	Size          int64               `json:"size"`
	StorageEngine string              `json:"storage_engine"`
}

// DatabasePerformance أداء قاعدة البيانات
type DatabasePerformance struct {
	QueryTime    time.Duration `json:"query_time"`
	Connections  int           `json:"connections"`
	Operations   int64         `json:"operations"`
}

// MemoryUsage استخدام الذاكرة
type MemoryUsage struct {
	HeapUsed        uint64  `json:"heap_used"`
	HeapTotal       uint64  `json:"heap_total"`
	UsagePercentage float64 `json:"usage_percentage"`
}

// CPUUsage استخدام المعالج
type CPUUsage struct {
	User   uint64 `json:"user"`
	System uint64 `json:"system"`
}

// PerformanceMetrics مقاييس الأداء
type PerformanceMetrics struct {
	CPU            CPUUsage          `json:"cpu"`
	Memory         MemoryUsage       `json:"memory"`
	Uptime         float64           `json:"uptime"`
	ActiveHandles  int               `json:"active_handles"`
	ActiveRequests int               `json:"active_requests"`
	HeapStatistics MemoryUsage       `json:"heap_statistics"`
	ResponseTimes  APIResponseTimes  `json:"response_times"`
	Throughput     SystemThroughput  `json:"throughput"`
	ErrorRates     ErrorRates        `json:"error_rates"`
}

// APIResponseTimes أوقات استجابة API
type APIResponseTimes struct {
	Average float64 `json:"average"`
	P95     float64 `json:"p95"`
	P99     float64 `json:"p99"`
}

// SystemThroughput إنتاجية النظام
type SystemThroughput struct {
	RequestsPerMinute int    `json:"requests_per_minute"`
	DataProcessed     string `json:"data_processed"`
}

// ErrorRates معدلات الخطأ
type ErrorRates struct {
	ErrorRate   string `json:"error_rate"`
	TotalErrors int    `json:"total_errors"`
}

// AIHealthAnalysis تحليل صحة الذكاء الاصطناعي
type AIHealthAnalysis struct {
	HealthScore float64 `json:"health_score"`
	Status      string  `json:"status"`
	RiskLevel   string  `json:"risk_level"`
}

// SecurityStatus حالة الأمان
type SecurityStatus struct {
	SSLEnabled       bool      `json:"ssl_enabled"`
	RateLimiting     bool      `json:"rate_limiting"`
	Authentication   bool      `json:"authentication"`
	LastSecurityScan time.Time `json:"last_security_scan"`
}

// ================================
// الطلبات والنتائج
// ================================

// SystemUpdateRequest طلب تحديث النظام
type SystemUpdateRequest struct {
	UpdateType    string `json:"update_type"`
	Version       string `json:"version"`
	Force         bool   `json:"force"`
	Backup        bool   `json:"backup"`
	AnalyzeImpact bool   `json:"analyze_impact"`
}

// SystemUpdateResult نتيجة تحديث النظام
type SystemUpdateResult struct {
	UpdateID          string        `json:"update_id"`
	EstimatedDuration time.Duration `json:"estimated_duration"`
	RequiresRestart   bool          `json:"requires_restart"`
	Steps             []string      `json:"steps"`
	BackupCreated     bool          `json:"backup_created"`
	ImpactAnalysis    interface{}   `json:"impact_analysis"`
}

// AIAnalyticsRequest طلب تحليلات الذكاء الاصطناعي
type AIAnalyticsRequest struct {
	Timeframe    string `json:"timeframe"`
	AnalysisType string `json:"analysis_type"`
}

// AIAnalyticsResult نتيجة تحليلات الذكاء الاصطناعي
type AIAnalyticsResult struct {
	Trends         interface{} `json:"trends"`
	Anomalies      interface{} `json:"anomalies"`
	Optimizations  interface{} `json:"optimizations"`
	Capacity       interface{} `json:"capacity"`
	Predictions    interface{} `json:"predictions"`
	RiskAssessment interface{} `json:"risk_assessment"`
	GeneratedAt    time.Time   `json:"generated_at"`
	AnalysisPeriod string      `json:"analysis_period"`
}

// UserAnalyticsRequest طلب تحليلات المستخدمين
type UserAnalyticsRequest struct {
	Timeframe    string `json:"timeframe"`
	UserSegment  string `json:"user_segment"`
	AnalysisDepth string `json:"analysis_depth"`
}

// UserAnalyticsResult نتيجة تحليلات المستخدمين
type UserAnalyticsResult struct {
	Overview       interface{} `json:"overview"`
	Behavior       interface{} `json:"behavior"`
	Segments       interface{} `json:"segments"`
	Predictions    interface{} `json:"predictions"`
	Retention      interface{} `json:"retention"`
	Recommendations interface{} `json:"recommendations"`
	GeneratedAt    time.Time   `json:"generated_at"`
}

// MaintenanceRequest طلب الصيانة
type MaintenanceRequest struct {
	Enabled  bool        `json:"enabled"`
	Message  string      `json:"message"`
	Schedule *time.Time  `json:"schedule,omitempty"`
	Duration *time.Duration `json:"duration,omitempty"`
}

// MaintenanceResult نتيجة الصيانة
type MaintenanceResult struct {
	Enabled     bool        `json:"enabled"`
	Message     string      `json:"message"`
	Schedule    *time.Time  `json:"schedule,omitempty"`
	Duration    *time.Duration `json:"duration,omitempty"`
	InitiatedBy string      `json:"initiated_by"`
}

// BackupRequest طلب النسخ الاحتياطي
type BackupRequest struct {
	Type        string `json:"type"`
	IncludeLogs bool   `json:"include_logs"`
	Optimize    bool   `json:"optimize"`
	Schedule    bool   `json:"schedule"`
}

// BackupResult نتيجة النسخ الاحتياطي
type BackupResult struct {
	BackupID string      `json:"backup_id"`
	Size     int64       `json:"size"`
	Path     string      `json:"path"`
	Type     string      `json:"type"`
	Strategy interface{} `json:"strategy"`
}

// OptimizationRequest طلب التحسين
type OptimizationRequest struct {
	Areas     []string `json:"areas"`
	Intensity string   `json:"intensity"`
}

// OptimizationResult نتيجة التحسين
type OptimizationResult struct {
	Improvements []string      `json:"improvements"`
	Metrics      interface{}   `json:"metrics"`
	Duration     time.Duration `json:"duration"`
}

// ================================
// السجلات والمراقبة
// ================================

// LogsRequest طلب السجلات
type LogsRequest struct {
	Level   string     `json:"level"`
	Limit   int        `json:"limit"`
	Page    int        `json:"page"`
	From    *time.Time `json:"from,omitempty"`
	To      *time.Time `json:"to,omitempty"`
	Analyze bool       `json:"analyze"`
}

// LogsResult نتيجة السجلات
type LogsResult struct {
	Logs       []LogEntry `json:"logs"`      // من shared_models
	Analysis   interface{} `json:"analysis"`
	Pagination Pagination  `json:"pagination"` // من shared_models
}

// ================================
// نماذج التقارير
// ================================

// ReportRequest طلب التقرير
type ReportRequest struct {
	Type       string     `json:"type"`
	Format     string     `json:"format"` // pdf, excel, csv, json
	DateRange  DateRange  `json:"date_range"` // من shared_models
	Filters    []Filter   `json:"filters"`    // من shared_models
	IncludeCharts bool    `json:"include_charts"`
}

// ReportResult نتيجة التقرير
type ReportResult struct {
	ReportID   string      `json:"report_id"`
	Type       string      `json:"type"`
	Data       interface{} `json:"data"`
	GeneratedAt time.Time  `json:"generated_at"`
	DownloadURL string     `json:"download_url,omitempty"`
}

// ================================
// إدارة المستخدمين والأدوار
// ================================

// UserManagementRequest طلب إدارة المستخدم
type UserManagementRequest struct {
	UserID    string   `json:"user_id"`
	Action    string   `json:"action"` // suspend, activate, change_role, reset_password
	NewRole   string   `json:"new_role,omitempty"`
	Reason    string   `json:"reason,omitempty"`
	Duration  *time.Duration `json:"duration,omitempty"`
}

// UserManagementResult نتيجة إدارة المستخدم
type UserManagementResult struct {
	UserID    string    `json:"user_id"`
	Action    string    `json:"action"`
	PreviousStatus string `json:"previous_status"`
	NewStatus string    `json:"new_status"`
	Timestamp time.Time `json:"timestamp"`
	PerformedBy string `json:"performed_by"`
}

// RolePermission إذن الدور
type RolePermission struct {
	Role       string   `json:"role"`
	Permissions []string `json:"permissions"`
	Description string  `json:"description"`
}

// ================================
// الإشعارات والإخطارات
// ================================

// Notification إشعار
type Notification struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	Title     string    `json:"title"`
	Message   string    `json:"message"`
	Priority  string    `json:"priority"`
	Recipient string    `json:"recipient"`
	Read      bool      `json:"read"`
	CreatedAt time.Time `json:"created_at"`
}

// BroadcastRequest طلب البث
type BroadcastRequest struct {
	Title     string   `json:"title"`
	Message   string   `json:"message"`
	Audience  string   `json:"audience"` // all, sellers, buyers, admins
	Channels  []string `json:"channels"` // email, push, in_app
	Schedule  *time.Time `json:"schedule,omitempty"`
}

// BroadcastResult نتيجة البث
type BroadcastResult struct {
	BroadcastID string    `json:"broadcast_id"`
	SentTo      int       `json:"sent_to"`
	Failed      int       `json:"failed"`
	Scheduled   bool      `json:"scheduled"`
	SentAt      time.Time `json:"sent_at"`
}

// ================================
// الهياكل المساعدة
// ================================

// SystemSummary ملخص النظام
type SystemSummary struct {
	Overall         string   `json:"overall"`
	Issues          []string `json:"issues"`
	Recommendations []string `json:"recommendations"`
}

// SystemMetrics مقاييس النظام
type SystemMetrics struct {
	Timestamp   time.Time          `json:"timestamp"`
	System      SystemInfo         `json:"system"`
	Performance PerformanceMetrics `json:"performance"`
	Database    DatabaseStatus     `json:"database"`
	Services    map[string]string  `json:"services"`
}

// UpdateImpactAnalysis تحليل تأثير التحديث
type UpdateImpactAnalysis struct {
	RiskScore         float64  `json:"risk_score"`
	Risks             []string `json:"risks"`
	Recommendations   []string `json:"recommendations"`
	AffectedServices  []string `json:"affected_services"`
}