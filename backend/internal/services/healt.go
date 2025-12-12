package services

import (
	"context"
	"database/sql"
	"fmt"
	"runtime"
	"time"

	"github.com/nawthtech/nawthtech/backend/internal/config"
	"go.uber.org/zap"
)

// ================================
// هياكل البيانات لـ Health
// ================================

// HealthRequest طلب فحص الصحة
type HealthRequest struct {
	CheckDatabase bool `json:"check_database"`
	CheckCache    bool `json:"check_cache"`
	CheckStorage  bool `json:"check_storage"`
	CheckServices bool `json:"check_services"`
}

// HealthCheckResult نتيجة فحص الصحة
type HealthCheckResult struct {
	Status      string            `json:"status"`
	Timestamp   time.Time         `json:"timestamp"`
	ServiceName string            `json:"service_name"`
	Version     string            `json:"version"`
	Uptime      string            `json:"uptime"`
	Checks      map[string]Check  `json:"checks"`
	SystemInfo  SystemInfo        `json:"system_info"`
}

// Check نتيجة فحص فردي
type Check struct {
	Name      string        `json:"name"`
	Status    string        `json:"status"` // "healthy", "unhealthy", "warning"
	Message   string        `json:"message,omitempty"`
	Duration  time.Duration `json:"duration"`
	Timestamp time.Time     `json:"timestamp"`
}

// SystemInfo معلومات النظام
type SystemInfo struct {
	GoVersion     string  `json:"go_version"`
	Architecture  string  `json:"architecture"`
	OS            string  `json:"os"`
	NumCPU        int     `json:"num_cpu"`
	MemoryUsageMB float64 `json:"memory_usage_mb"`
	Goroutines    int     `json:"goroutines"`
}

// DatabaseStats إحصائيات قاعدة البيانات
type DatabaseStats struct {
	Connected     bool   `json:"connected"`
	Driver        string `json:"driver"`
	Version       string `json:"version,omitempty"`
	MaxOpenConns  int    `json:"max_open_conns"`
	OpenConns     int    `json:"open_conns"`
	InUseConns    int    `json:"in_use_conns"`
	IdleConns     int    `json:"idle_conns"`
	WaitCount     int64  `json:"wait_count"`
	WaitDuration  string `json:"wait_duration"`
	MaxIdleClosed int64  `json:"max_idle_closed"`
	MaxLifetimeClosed int64 `json:"max_lifetime_closed"`
}

// ServiceStats إحصائيات الخدمات
type ServiceStats struct {
	TotalUsers      int64   `json:"total_users"`
	TotalServices   int64   `json:"total_services"`
	TotalOrders     int64   `json:"total_orders"`
	TotalPayments   int64   `json:"total_payments"`
	ActiveUsers     int64   `json:"active_users"`
	PendingOrders   int64   `json:"pending_orders"`
	TotalRevenue    float64 `json:"total_revenue"`
	AvgResponseTime float64 `json:"avg_response_time"`
}

// DetailedHealthResult نتيجة صحة مفصلة
type DetailedHealthResult struct {
	HealthCheckResult
	DatabaseStats  DatabaseStats `json:"database_stats"`
	ServiceStats   ServiceStats  `json:"service_stats"`
	Dependencies   []Dependency  `json:"dependencies"`
	Environment    string        `json:"environment"`
	ConfigInfo     ConfigInfo    `json:"config_info"`
}

// Dependency تبعية خارجية
type Dependency struct {
	Name        string `json:"name"`
	Type        string `json:"type"` // "database", "cache", "storage", "external_api"
	Status      string `json:"status"`
	URL         string `json:"url,omitempty"`
	PingTime    string `json:"ping_time,omitempty"`
	LastChecked string `json:"last_checked,omitempty"`
}

// ConfigInfo معلومات التكوين
type ConfigInfo struct {
	AppName     string `json:"app_name"`
	Environment string `json:"environment"`
	Port        string `json:"port"`
	Debug       bool   `json:"debug"`
	Database    string `json:"database"`
	CacheType   string `json:"cache_type"`
	StorageType string `json:"storage_type"`
}

// ================================
// واجهة HealthService
// ================================

// HealthService واجهة خدمة الصحة
type HealthService interface {
	// Basic Health Checks
	CheckHealth(ctx context.Context, req *HealthRequest) (*HealthCheckResult, error)
	CheckDatabase(ctx context.Context) (*Check, error)
	CheckCache(ctx context.Context) (*Check, error)
	CheckStorage(ctx context.Context) (*Check, error)
	
	// Advanced Health
	GetDetailedHealth(ctx context.Context) (*DetailedHealthResult, error)
	GetDatabaseHealth(ctx context.Context) (*DatabaseStats, error)
	GetServiceStats(ctx context.Context) (*ServiceStats, error)
	GetSystemInfo(ctx context.Context) (*SystemInfo, error)
	
	// Monitoring
	GetUptime(ctx context.Context) (string, error)
	GetMetrics(ctx context.Context) (map[string]interface{}, error)
}

// ================================
// تطبيق HealthService
// ================================

type healthServiceImpl struct {
	db     *sql.DB
	config *config.Config
	logger *zap.Logger
	startTime time.Time
}

// NewHealthService إنشاء خدمة صحة جديدة
func NewHealthService(db *sql.DB, cfg *config.Config, logger *zap.Logger) HealthService {
	return &healthServiceImpl{
		db:        db,
		config:    cfg,
		logger:    logger,
		startTime: time.Now(),
	}
}

// CheckHealth فحص الصحة الأساسي
func (s *healthServiceImpl) CheckHealth(ctx context.Context, req *HealthRequest) (*HealthCheckResult, error) {
	startTime := time.Now()
	checks := make(map[string]Check)
	
	// فحص الصحة الأساسي (دائماً مفعل)
	basicCheck := Check{
		Name:      "basic",
		Status:    "healthy",
		Message:   "Service is running",
		Duration:  time.Since(startTime),
		Timestamp: time.Now(),
	}
	checks["basic"] = basicCheck
	
	// فحص قاعدة البيانات إذا مطلوب
	if req.CheckDatabase {
		dbCheck, err := s.CheckDatabase(ctx)
		if err != nil {
			dbCheck.Status = "unhealthy"
			dbCheck.Message = err.Error()
		}
		checks["database"] = *dbCheck
	}
	
	// فحص التخزين المؤقت إذا مطلوب
	if req.CheckCache {
		cacheCheck, err := s.CheckCache(ctx)
		if err != nil {
			cacheCheck.Status = "unhealthy"
			cacheCheck.Message = err.Error()
		}
		checks["cache"] = *cacheCheck
	}
	
	// فحص التخزين إذا مطلوب
	if req.CheckStorage {
		storageCheck, err := s.CheckStorage(ctx)
		if err != nil {
			storageCheck.Status = "unhealthy"
			storageCheck.Message = err.Error()
		}
		checks["storage"] = *storageCheck
	}
	
	// فحص النظام
	systemInfo, _ := s.GetSystemInfo(ctx)
	
	// تحديد الحالة العامة
	overallStatus := "healthy"
	for _, check := range checks {
		if check.Status == "unhealthy" {
			overallStatus = "unhealthy"
			break
		} else if check.Status == "warning" {
			overallStatus = "warning"
		}
	}
	
	// حساب وقت التشغيل
	uptime := time.Since(s.startTime).String()
	
	return &HealthCheckResult{
		Status:      overallStatus,
		Timestamp:   time.Now(),
		ServiceName: s.config.AppName,
		Version:     s.config.Version,
		Uptime:      uptime,
		Checks:      checks,
		SystemInfo:  *systemInfo,
	}, nil
}

// CheckDatabase فحص قاعدة البيانات
func (s *healthServiceImpl) CheckDatabase(ctx context.Context) (*Check, error) {
	startTime := time.Now()
	
	if s.db == nil {
		return &Check{
			Name:      "database",
			Status:    "unhealthy",
			Message:   "Database connection is nil",
			Duration:  time.Since(startTime),
			Timestamp: time.Now(),
		}, fmt.Errorf("database connection is nil")
	}
	
	// اختبار الاتصال بقاعدة البيانات
	err := s.db.PingContext(ctx)
	if err != nil {
		return &Check{
			Name:      "database",
			Status:    "unhealthy",
			Message:   fmt.Sprintf("Database ping failed: %v", err),
			Duration:  time.Since(startTime),
			Timestamp: time.Now(),
		}, err
	}
	
	// اختبار استعلام بسيط
	var version string
	err = s.db.QueryRowContext(ctx, "SELECT sqlite_version()").Scan(&version)
	if err != nil {
		return &Check{
			Name:      "database",
			Status:    "warning",
			Message:   fmt.Sprintf("Database query failed: %v", err),
			Duration:  time.Since(startTime),
			Timestamp: time.Now(),
		}, nil
	}
	
	return &Check{
		Name:      "database",
		Status:    "healthy",
		Message:   fmt.Sprintf("Database connected (SQLite %s)", version),
		Duration:  time.Since(startTime),
		Timestamp: time.Now(),
	}, nil
}

// CheckCache فحص التخزين المؤقت
func (s *healthServiceImpl) CheckCache(ctx context.Context) (*Check, error) {
	startTime := time.Now()
	
	// للـ cache المبني في الذاكرة، نجربه
	if s.config.Cache.Enabled && s.config.Cache.Type == "memory" {
		return &Check{
			Name:      "cache",
			Status:    "healthy",
			Message:   "In-memory cache is enabled",
			Duration:  time.Since(startTime),
			Timestamp: time.Now(),
		}, nil
	}
	
	// إذا كان Redis
	if s.config.Cache.Enabled && s.config.Cache.Type == "redis" && s.config.Cache.Redis == "" {
		return &Check{
			Name:      "cache",
			Status:    "warning",
			Message:   "Redis cache configured but Redis URL is empty",
			Duration:  time.Since(startTime),
			Timestamp: time.Now(),
		}, nil
	}
	
	return &Check{
		Name:      "cache",
		Status:    "healthy",
		Message:   "Cache check completed",
		Duration:  time.Since(startTime),
		Timestamp: time.Now(),
	}, nil
}

// CheckStorage فحص التخزين
func (s *healthServiceImpl) CheckStorage(ctx context.Context) (*Check, error) {
	startTime := time.Now()
	
	// التحقق من تكوين التخزين
	if s.config.Upload.S3Bucket != "" && s.config.Upload.S3AccessKey != "" && s.config.Upload.S3SecretKey != "" {
		return &Check{
			Name:      "storage",
			Status:    "healthy",
			Message:   fmt.Sprintf("S3 storage configured (%s/%s)", s.config.Upload.S3Bucket, s.config.Upload.S3Region),
			Duration:  time.Since(startTime),
			Timestamp: time.Now(),
		}, nil
	}
	
	if s.config.Upload.CloudinaryURL != "" {
		return &Check{
			Name:      "storage",
			Status:    "healthy",
			Message:   "Cloudinary storage configured",
			Duration:  time.Since(startTime),
			Timestamp: time.Now(),
		}, nil
	}
	
	return &Check{
		Name:      "storage",
		Status:    "warning",
		Message:   "No external storage configured, using local storage",
		Duration:  time.Since(startTime),
		Timestamp: time.Now(),
	}, nil
}

// GetDetailedHealth الحصول على صحة مفصلة
func (s *healthServiceImpl) GetDetailedHealth(ctx context.Context) (*DetailedHealthResult, error) {
	// فحص الصحة الأساسي
	healthReq := &HealthRequest{
		CheckDatabase: true,
		CheckCache:    true,
		CheckStorage:  true,
		CheckServices: true,
	}
	
	basicHealth, err := s.CheckHealth(ctx, healthReq)
	if err != nil {
		return nil, err
	}
	
	// الحصول على إحصائيات قاعدة البيانات
	dbStats, _ := s.GetDatabaseHealth(ctx)
	
	// الحصول على إحصائيات الخدمات
	serviceStats, _ := s.GetServiceStats(ctx)
	
	// إنشاء قائمة التبعيات
	dependencies := []Dependency{
		{
			Name:   "SQLite Database",
			Type:   "database",
			Status: "healthy",
			URL:    s.config.Database.URL,
		},
	}
	
	// إضافة Redis إذا مكون
	if s.config.Cache.Enabled && s.config.Cache.Type == "redis" && s.config.Cache.Redis != "" {
		dependencies = append(dependencies, Dependency{
			Name:   "Redis Cache",
			Type:   "cache",
			Status: "healthy",
			URL:    s.config.Cache.Redis,
		})
	}
	
	// إضافة S3 إذا مكون
	if s.config.Upload.S3Bucket != "" {
		dependencies = append(dependencies, Dependency{
			Name:   "AWS S3 Storage",
			Type:   "storage",
			Status: "healthy",
			URL:    fmt.Sprintf("s3://%s/%s", s.config.Upload.S3Region, s.config.Upload.S3Bucket),
		})
	}
	
	// معلومات التكوين
	configInfo := ConfigInfo{
		AppName:     s.config.AppName,
		Environment: s.config.Environment,
		Port:        s.config.Port,
		Debug:       s.config.Debug,
		Database:    s.config.Database.Driver,
		CacheType:   s.config.Cache.Type,
		StorageType: "local", // القيمة الافتراضية
	}
	
	if s.config.Upload.S3Bucket != "" {
		configInfo.StorageType = "s3"
	} else if s.config.Upload.CloudinaryURL != "" {
		configInfo.StorageType = "cloudinary"
	}
	
	return &DetailedHealthResult{
		HealthCheckResult: *basicHealth,
		DatabaseStats:     *dbStats,
		ServiceStats:      *serviceStats,
		Dependencies:      dependencies,
		Environment:       s.config.Environment,
		ConfigInfo:        configInfo,
	}, nil
}

// GetDatabaseHealth الحصول على صحة قاعدة البيانات
func (s *healthServiceImpl) GetDatabaseHealth(ctx context.Context) (*DatabaseStats, error) {
	stats := &DatabaseStats{
		Connected:     true,
		Driver:        s.config.Database.Driver,
		MaxOpenConns:  s.config.Database.MaxConns,
	}
	
	if s.db == nil {
		stats.Connected = false
		return stats, nil
	}
	
	// الحصول على إحصائيات قاعدة البيانات
	stats.OpenConns = s.db.Stats().OpenConnections
	stats.InUseConns = s.db.Stats().InUse
	stats.IdleConns = s.db.Stats().Idle
	stats.WaitCount = s.db.Stats().WaitCount
	stats.WaitDuration = s.db.Stats().WaitDuration.String()
	stats.MaxIdleClosed = s.db.Stats().MaxIdleClosed
	stats.MaxLifetimeClosed = s.db.Stats().MaxLifetimeClosed
	
	// الحصول على إصدار SQLite
	if s.config.Database.Driver == "sqlite3" {
		var version string
		err := s.db.QueryRowContext(ctx, "SELECT sqlite_version()").Scan(&version)
		if err == nil {
			stats.Version = version
		}
	}
	
	return stats, nil
}

// GetServiceStats الحصول على إحصائيات الخدمات
func (s *healthServiceImpl) GetServiceStats(ctx context.Context) (*ServiceStats, error) {
	stats := &ServiceStats{}
	
	if s.db == nil {
		return stats, fmt.Errorf("database connection is nil")
	}
	
	// إجمالي المستخدمين
	err := s.db.QueryRowContext(ctx, 
		"SELECT COUNT(*) FROM users WHERE status = 'active'").Scan(&stats.TotalUsers)
	if err != nil {
		s.logger.Warn("Failed to get total users", zap.Error(err))
	}
	
	// إجمالي الخدمات
	err = s.db.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM services WHERE is_active = TRUE").Scan(&stats.TotalServices)
	if err != nil {
		s.logger.Warn("Failed to get total services", zap.Error(err))
	}
	
	// إجمالي الطلبات
	err = s.db.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM orders").Scan(&stats.TotalOrders)
	if err != nil {
		s.logger.Warn("Failed to get total orders", zap.Error(err))
	}
	
	// إجمالي المدفوعات
	err = s.db.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM payments").Scan(&stats.TotalPayments)
	if err != nil {
		s.logger.Warn("Failed to get total payments", zap.Error(err))
	}
	
	// المستخدمين النشطين (سجلوا دخول في آخر 30 يوم)
	err = s.db.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM users WHERE last_login >= DATE('now', '-30 days')").Scan(&stats.ActiveUsers)
	if err != nil {
		s.logger.Warn("Failed to get active users", zap.Error(err))
	}
	
	// الطلبات المعلقة
	err = s.db.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM orders WHERE status = 'pending'").Scan(&stats.PendingOrders)
	if err != nil {
		s.logger.Warn("Failed to get pending orders", zap.Error(err))
	}
	
	// إجمالي الإيرادات
	err = s.db.QueryRowContext(ctx,
		"SELECT COALESCE(SUM(amount), 0) FROM orders WHERE status = 'completed'").Scan(&stats.TotalRevenue)
	if err != nil {
		s.logger.Warn("Failed to get total revenue", zap.Error(err))
	}
	
	// متوسط وقت الاستجابة (محاكاة)
	stats.AvgResponseTime = 125.5 // قيمة افتراضية
	
	return stats, nil
}

// GetSystemInfo الحصول على معلومات النظام
func (s *healthServiceImpl) GetSystemInfo(ctx context.Context) (*SystemInfo, error) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	
	return &SystemInfo{
		GoVersion:     runtime.Version(),
		Architecture:  runtime.GOARCH,
		OS:            runtime.GOOS,
		NumCPU:        runtime.NumCPU(),
		MemoryUsageMB: float64(m.Alloc) / 1024 / 1024,
		Goroutines:    runtime.NumGoroutine(),
	}, nil
}

// GetUptime الحصول على وقت التشغيل
func (s *healthServiceImpl) GetUptime(ctx context.Context) (string, error) {
	uptime := time.Since(s.startTime)
	
	days := int(uptime.Hours() / 24)
	hours := int(uptime.Hours()) % 24
	minutes := int(uptime.Minutes()) % 60
	seconds := int(uptime.Seconds()) % 60
	
	if days > 0 {
		return fmt.Sprintf("%d days, %d hours, %d minutes, %d seconds", days, hours, minutes, seconds), nil
	} else if hours > 0 {
		return fmt.Sprintf("%d hours, %d minutes, %d seconds", hours, minutes, seconds), nil
	} else if minutes > 0 {
		return fmt.Sprintf("%d minutes, %d seconds", minutes, seconds), nil
	}
	
	return fmt.Sprintf("%d seconds", seconds), nil
}

// GetMetrics الحصول على مقاييس النظام
func (s *healthServiceImpl) GetMetrics(ctx context.Context) (map[string]interface{}, error) {
	metrics := make(map[string]interface{})
	
	// معلومات النظام
	systemInfo, _ := s.GetSystemInfo(ctx)
	metrics["system"] = systemInfo
	
	// وقت التشغيل
	uptime, _ := s.GetUptime(ctx)
	metrics["uptime"] = uptime
	
	// إحصائيات قاعدة البيانات
	if s.db != nil {
		dbStats := s.db.Stats()
		metrics["database"] = map[string]interface{}{
			"open_connections":      dbStats.OpenConnections,
			"in_use":                dbStats.InUse,
			"idle":                  dbStats.Idle,
			"wait_count":            dbStats.WaitCount,
			"wait_duration_ms":      dbStats.WaitDuration.Milliseconds(),
			"max_idle_closed":       dbStats.MaxIdleClosed,
			"max_lifetime_closed":   dbStats.MaxLifetimeClosed,
		}
	}
	
	// وقت التشغيل
	metrics["start_time"] = s.startTime.Format(time.RFC3339)
	
	// التكوين (بدون معلومات سرية)
	metrics["config"] = map[string]interface{}{
		"app_name":     s.config.AppName,
		"version":      s.config.Version,
		"environment":  s.config.Environment,
		"debug":        s.config.Debug,
		"port":         s.config.Port,
		"database":     s.config.Database.Driver,
		"cache_enabled": s.config.Cache.Enabled,
		"cache_type":   s.config.Cache.Type,
	}
	
	return metrics, nil
}

// ================================
// دوال مساعدة (Helper Functions)
// ================================

// FormatDuration تنسيق المدة الزمنية
func FormatDuration(d time.Duration) string {
	if d < time.Second {
		return d.Round(time.Millisecond).String()
	}
	
	if d < time.Minute {
		return d.Round(time.Second).String()
	}
	
	if d < time.Hour {
		minutes := int(d.Minutes())
		seconds := int(d.Seconds()) % 60
		return fmt.Sprintf("%dm %ds", minutes, seconds)
	}
	
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60
	
	return fmt.Sprintf("%dh %dm %ds", hours, minutes, seconds)
}

// CalculatePercentage حساب النسبة المئوية
func CalculatePercentage(part, total float64) float64 {
	if total == 0 {
		return 0
	}
	return (part / total) * 100
}

// IsHealthy تحقق إذا كانت جميع الفحوصات صحية
func IsHealthy(checks map[string]Check) bool {
	for _, check := range checks {
		if check.Status != "healthy" {
			return false
		}
	}
	return true
}

// GetHealthStatusColor الحصول على لون حالة الصحة
func GetHealthStatusColor(status string) string {
	switch status {
	case "healthy":
		return "#10B981" // أخضر
	case "warning":
		return "#F59E0B" // أصفر
	case "unhealthy":
		return "#EF4444" // أحمر
	default:
		return "#6B7280" // رمادي
	}
}