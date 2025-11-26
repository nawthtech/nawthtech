package models

import (
	"time"
)

// SystemUpdateRequest طلب تحديث النظام
type SystemUpdateRequest struct {
	UpdateType    string `json:"updateType"`
	Version       string `json:"version"`
	Force         bool   `json:"force"`
	Backup        bool   `json:"backup"`
	AnalyzeImpact bool   `json:"analyzeImpact"`
}

// SystemUpdateResult نتيجة تحديث النظام
type SystemUpdateResult struct {
	UpdateID          string        `json:"updateId"`
	EstimatedDuration time.Duration `json:"estimatedDuration"`
	RequiresRestart   bool          `json:"requiresRestart"`
	Steps             []string      `json:"steps"`
	BackupCreated     bool          `json:"backupCreated"`
	ImpactAnalysis    interface{}   `json:"impactAnalysis"`
}

// SystemStatus حالة النظام
type SystemStatus struct {
	Disk        DiskStatus        `json:"disk"`
	Database    DatabaseStatus    `json:"database"`
	System      SystemInfo        `json:"system"`
	Performance PerformanceMetrics `json:"performance"`
	Services    map[string]string `json:"services"`
	AIAnalysis  AIHealthAnalysis  `json:"aiAnalysis"`
	Security    SecurityStatus    `json:"security"`
	LastChecked time.Time         `json:"lastChecked"`
}

// DiskStatus حالة القرص
type DiskStatus struct {
	Free           uint64   `json:"free"`
	Size           uint64   `json:"size"`
	Used           uint64   `json:"used"`
	FreePercentage float64  `json:"freePercentage"`
	Path           string   `json:"path"`
	Threshold      string   `json:"threshold"`
	Recommendations []string `json:"recommendations"`
}

// DatabaseStatus حالة قاعدة البيانات
type DatabaseStatus struct {
	Connected    bool                `json:"connected"`
	ReadyState   string              `json:"readyState"`
	DBName       string              `json:"dbName"`
	Host         string              `json:"host"`
	Connections  int                 `json:"connections"`
	Performance  DatabasePerformance `json:"performance"`
	Collections  int                 `json:"collections"`
	Size         int64               `json:"size"`
	StorageEngine string             `json:"storageEngine"`
}

// DatabasePerformance أداء قاعدة البيانات
type DatabasePerformance struct {
	QueryTime    time.Duration `json:"queryTime"`
	Connections  int           `json:"connections"`
	Operations   int64         `json:"operations"`
}

// SystemInfo معلومات النظام
type SystemInfo struct {
	Version      string     `json:"version"`
	NodeVersion  string     `json:"nodeVersion"`
	Environment  string     `json:"environment"`
	Platform     string     `json:"platform"`
	Arch         string     `json:"arch"`
	Uptime       float64    `json:"uptime"`
	Memory       MemoryUsage `json:"memory"`
	CPU          CPUUsage   `json:"cpu"`
	PID          int        `json:"pid"`
}

// MemoryUsage استخدام الذاكرة
type MemoryUsage struct {
	HeapUsed       uint64  `json:"heapUsed"`
	HeapTotal      uint64  `json:"heapTotal"`
	UsagePercentage float64 `json:"usagePercentage"`
}

// CPUUsage استخدام المعالج
type CPUUsage struct {
	User   uint64 `json:"user"`
	System uint64 `json:"system"`
}

// PerformanceMetrics مقاييس الأداء
type PerformanceMetrics struct {
	CPU            CPUUsage    `json:"cpu"`
	Memory         MemoryUsage `json:"memory"`
	Uptime         float64     `json:"uptime"`
	ActiveHandles  int         `json:"activeHandles"`
	ActiveRequests int         `json:"activeRequests"`
	HeapStatistics MemoryUsage `json:"heapStatistics"`
	ResponseTimes  APIResponseTimes `json:"responseTimes"`
	Throughput     SystemThroughput `json:"throughput"`
	ErrorRates     ErrorRates       `json:"errorRates"`
}

// APIResponseTimes أوقات استجابة API
type APIResponseTimes struct {
	Average float64 `json:"average"`
	P95     float64 `json:"p95"`
	P99     float64 `json:"p99"`
}

// SystemThroughput إنتاجية النظام
type SystemThroughput struct {
	RequestsPerMinute int    `json:"requestsPerMinute"`
	DataProcessed     string `json:"dataProcessed"`
}

// ErrorRates معدلات الخطأ
type ErrorRates struct {
	ErrorRate   string `json:"errorRate"`
	TotalErrors int    `json:"totalErrors"`
}

// AIHealthAnalysis تحليل صحة الذكاء الاصطناعي
type AIHealthAnalysis struct {
	HealthScore float64 `json:"healthScore"`
	Status      string  `json:"status"`
	RiskLevel   string  `json:"riskLevel"`
}

// SecurityStatus حالة الأمان
type SecurityStatus struct {
	SSLEnabled       bool      `json:"sslEnabled"`
	RateLimiting     bool      `json:"rateLimiting"`
	Authentication   bool      `json:"authentication"`
	LastSecurityScan time.Time `json:"lastSecurityScan"`
}

// AIAnalyticsRequest طلب تحليلات الذكاء الاصطناعي
type AIAnalyticsRequest struct {
	Timeframe    string `json:"timeframe"`
	AnalysisType string `json:"analysisType"`
}

// AIAnalyticsResult نتيجة تحليلات الذكاء الاصطناعي
type AIAnalyticsResult struct {
	Trends      interface{} `json:"trends"`
	Anomalies   interface{} `json:"anomalies"`
	Optimizations interface{} `json:"optimizations"`
	Capacity    interface{} `json:"capacity"`
	Predictions interface{} `json:"predictions"`
	RiskAssessment interface{} `json:"riskAssessment"`
	GeneratedAt time.Time   `json:"generatedAt"`
	AnalysisPeriod string    `json:"analysisPeriod"`
}

// UserAnalyticsRequest طلب تحليلات المستخدمين
type UserAnalyticsRequest struct {
	Timeframe   string `json:"timeframe"`
	UserSegment string `json:"userSegment"`
	AnalysisDepth string `json:"analysisDepth"`
}

// UserAnalyticsResult نتيجة تحليلات المستخدمين
type UserAnalyticsResult struct {
	Overview       interface{} `json:"overview"`
	Behavior       interface{} `json:"behavior"`
	Segments       interface{} `json:"segments"`
	Predictions    interface{} `json:"predictions"`
	Retention      interface{} `json:"retention"`
	Recommendations interface{} `json:"recommendations"`
	GeneratedAt    time.Time   `json:"generatedAt"`
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
	InitiatedBy string      `json:"initiatedBy"`
}

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
	Logs      []LogEntry `json:"logs"`
	Analysis  interface{} `json:"analysis"`
	Pagination Pagination `json:"pagination"`
}

// LogEntry إدخال السجل
type LogEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Level     string    `json:"level"`
	Message   string    `json:"message"`
	Service   string    `json:"service"`
	UserID    string    `json:"userId,omitempty"`
	RequestID string    `json:"requestId,omitempty"`
}

// Pagination الترقيم
type Pagination struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
	Total int `json:"total"`
	Pages int `json:"pages"`
}

// BackupRequest طلب النسخ الاحتياطي
type BackupRequest struct {
	Type        string `json:"type"`
	IncludeLogs bool   `json:"includeLogs"`
	Optimize    bool   `json:"optimize"`
	Schedule    bool   `json:"schedule"`
}

// BackupResult نتيجة النسخ الاحتياطي
type BackupResult struct {
	BackupID string `json:"backupId"`
	Size     int64  `json:"size"`
	Path     string `json:"path"`
	Type     string `json:"type"`
	Strategy interface{} `json:"strategy"`
}

// OptimizationRequest طلب التحسين
type OptimizationRequest struct {
	Areas     []string `json:"areas"`
	Intensity string   `json:"intensity"`
}

// OptimizationResult نتيجة التحسين
type OptimizationResult struct {
	Improvements []string    `json:"improvements"`
	Metrics      interface{} `json:"metrics"`
	Duration     time.Duration `json:"duration"`
}

// HealthCheckResult نتيجة فحص الصحة
type HealthCheckResult struct {
	Service     string      `json:"service"`
	Status      string      `json:"status"`
	ResponseTime string     `json:"responseTime,omitempty"`
	Error       string      `json:"error,omitempty"`
	Usage       string      `json:"usage,omitempty"`
	Details     interface{} `json:"details,omitempty"`
}

// SystemSummary ملخص النظام
type SystemSummary struct {
	Overall         string   `json:"overall"`
	Issues          []string `json:"issues"`
	Recommendations []string `json:"recommendations"`
}

// SystemMetrics مقاييس النظام
type SystemMetrics struct {
	Timestamp time.Time              `json:"timestamp"`
	System    SystemInfo             `json:"system"`
	Performance PerformanceMetrics   `json:"performance"`
	Database  DatabaseStatus         `json:"database"`
	Services  map[string]string      `json:"services"`
}

// UpdateImpactAnalysis تحليل تأثير التحديث
type UpdateImpactAnalysis struct {
	RiskScore    float64     `json:"riskScore"`
	Risks        []string    `json:"risks"`
	Recommendations []string `json:"recommendations"`
	AffectedServices []string `json:"affectedServices"`
}
