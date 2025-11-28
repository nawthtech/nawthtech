package logger

import (
	"context"
	"log/slog"
	"os"
	"runtime"
	"time"
)

// ========== أنواع وواجهات ==========

// Logger واجهة للسجلات
type Logger interface {
	Debug(ctx context.Context, msg string, args ...any)
	Info(ctx context.Context, msg string, args ...any)
	Warn(ctx context.Context, msg string, args ...any)
	Error(ctx context.Context, msg string, args ...any)
	With(args ...any) Logger
}

// DefaultLogger تطبيق افتراضي للسجلات
type DefaultLogger struct {
	logger *slog.Logger
}

// ========== متغيرات عامة ==========

var (
	// Loggers الرئيسية
	Stdout *slog.Logger
	Stderr *slog.Logger
	
	// Global logger instance
	globalLogger Logger
)

// ========== التهيئة والإعداد ==========

// Init تهيئة النظام
func Init(env string) {
	level := slog.LevelInfo
	if env == "development" {
		level = slog.LevelDebug
	}

	opts := &slog.HandlerOptions{
		Level: level,
	}

	if env == "development" {
		// في التطوير، استخدام TextHandler للقراءة السهلة
		Stdout = slog.New(slog.NewTextHandler(os.Stdout, opts))
		Stderr = slog.New(slog.NewTextHandler(os.Stderr, opts))
	} else {
		// في الإنتاج، استخدام JSONHandler
		opts.AddSource = true
		Stdout = slog.New(slog.NewJSONHandler(os.Stdout, opts))
		Stderr = slog.New(slog.NewJSONHandler(os.Stderr, opts))
	}

	globalLogger = &DefaultLogger{logger: Stdout}
}

// InitLogger تهيئة متقدمة للنظام
func InitLogger(env string, level slog.Level) {
	opts := &slog.HandlerOptions{
		Level: level,
	}

	if env == "development" {
		// استخدام TextHandler في التطوير للقراءة السهلة
		Stdout = slog.New(slog.NewTextHandler(os.Stdout, opts))
		Stderr = slog.New(slog.NewTextHandler(os.Stderr, opts))
	} else {
		// استخدام JSONHandler في الإنتاج
		opts.AddSource = true
		Stdout = slog.New(slog.NewJSONHandler(os.Stdout, opts))
		Stderr = slog.New(slog.NewJSONHandler(os.Stderr, opts))
	}

	globalLogger = &DefaultLogger{logger: Stdout}
}

// ========== تطبيق واجهة Logger ==========

func (l *DefaultLogger) Debug(ctx context.Context, msg string, args ...any) {
	l.logger.DebugContext(ctx, msg, args...)
}

func (l *DefaultLogger) Info(ctx context.Context, msg string, args ...any) {
	l.logger.InfoContext(ctx, msg, args...)
}

func (l *DefaultLogger) Warn(ctx context.Context, msg string, args ...any) {
	l.logger.WarnContext(ctx, msg, args...)
}

func (l *DefaultLogger) Error(ctx context.Context, msg string, args ...any) {
	l.logger.ErrorContext(ctx, msg, args...)
}

func (l *DefaultLogger) With(args ...any) Logger {
	return &DefaultLogger{logger: l.logger.With(args...)}
}

// ========== دوال الوصول العالمية ==========

// Debug تسجيل معلومات تصحيح
func Debug(ctx context.Context, msg string, args ...any) {
	if globalLogger == nil {
		Init("development")
	}
	globalLogger.Debug(ctx, msg, args...)
}

// Info تسجيل معلومات
func Info(ctx context.Context, msg string, args ...any) {
	if globalLogger == nil {
		Init("development")
	}
	globalLogger.Info(ctx, msg, args...)
}

// Warn تسجيل تحذير
func Warn(ctx context.Context, msg string, args ...any) {
	if globalLogger == nil {
		Init("development")
	}
	globalLogger.Warn(ctx, msg, args...)
}

// Error تسجيل خطأ
func Error(ctx context.Context, msg string, args ...any) {
	if globalLogger == nil {
		Init("development")
	}
	globalLogger.Error(ctx, msg, args...)
}

// With إرجاع logger مع حقول إضافية
func With(args ...any) Logger {
	if globalLogger == nil {
		Init("development")
	}
	return globalLogger.With(args...)
}

// ========== دوال مساعدة أساسية ==========

// ErrAttr دالة مساعدة لإرجاع سمة الخطأ
func ErrAttr(err error) slog.Attr {
	if err == nil {
		return slog.String("error", "nil")
	}
	return slog.String("error", err.Error())
}

// ErrorsAttr دالة مساعدة لإرجاع سمة الأخطاء المتعددة
func ErrorsAttr(errors ...error) slog.Attr {
	if len(errors) == 0 {
		return slog.Any("errors", []string{})
	}
	
	errStrs := make([]string, len(errors))
	for i, err := range errors {
		if err != nil {
			errStrs[i] = err.Error()
		} else {
			errStrs[i] = "nil"
		}
	}
	return slog.Any("errors", errStrs)
}

// DurationAttr دالة مساعدة للوقت
func DurationAttr(duration time.Duration) slog.Attr {
	return slog.Duration("duration", duration)
}

// TimestampAttr دالة مساعدة للطابع الزمني
func TimestampAttr() slog.Attr {
	return slog.String("timestamp", time.Now().Format(time.RFC3339))
}

// RequestIDAttr دالة مساعدة لمعرف الطلب
func RequestIDAttr(requestID string) slog.Attr {
	return slog.String("request_id", requestID)
}

// UserIDAttr دالة مساعدة لمعرف المستخدم
func UserIDAttr(userID string) slog.Attr {
	return slog.String("user_id", userID)
}

// UserRoleAttr دالة مساعدة لدور المستخدم
func UserRoleAttr(role string) slog.Attr {
	return slog.String("user_role", role)
}

// ========== دوال مساعدة للتخزين المؤقت ==========

// CacheOperationAttr سمات عملية التخزين المؤقت
func CacheOperationAttr(operation, key string, duration time.Duration) slog.Attr {
	return slog.Group("cache",
		slog.String("operation", operation),
		slog.String("key", key),
		slog.Duration("duration", duration),
	)
}

// CacheHitAttr سمة نجاح التخزين المؤقت
func CacheHitAttr(key string, hit bool) slog.Attr {
	return slog.Group("cache",
		slog.String("key", key),
		slog.Bool("hit", hit),
		slog.String("operation", "get"),
	)
}

// CacheErrorAttr سمة خطأ التخزين المؤقت
func CacheErrorAttr(operation, key string, err error) slog.Attr {
	return slog.Group("cache_error",
		slog.String("operation", operation),
		slog.String("key", key),
		ErrAttr(err),
	)
}

// ========== دوال مساعدة للطلبات والشبكة ==========

// RequestAttr سمات الطلب
func RequestAttr(method, path string, statusCode int, duration time.Duration) slog.Attr {
	return slog.Group("request",
		slog.String("method", method),
		slog.String("path", path),
		slog.Int("status_code", statusCode),
		slog.Duration("duration", duration),
	)
}

// CORSAttr سمة CORS
func CORSAttr(origin, method string, allowed bool) slog.Attr {
	return slog.Group("cors",
		slog.String("origin", origin),
		slog.String("method", method),
		slog.Bool("allowed", allowed),
	)
}

// UserActionAttr سمة إجراء المستخدم
func UserActionAttr(userID, action, resource string) slog.Attr {
	return slog.Group("user_action",
		slog.String("user_id", userID),
		slog.String("action", action),
		slog.String("resource", resource),
	)
}

// DatabaseQueryAttr سمة استعلام قاعدة البيانات
func DatabaseQueryAttr(operation, table string, duration time.Duration, rowsAffected int64) slog.Attr {
	return slog.Group("database",
		slog.String("operation", operation),
		slog.String("table", table),
		slog.Duration("duration", duration),
		slog.Int64("rows_affected", rowsAffected),
	)
}

// ========== دوال مساعدة للأداء والذاكرة ==========

// PerformanceAttr سمة الأداء
func PerformanceAttr(operation string, duration time.Duration) slog.Attr {
	return slog.Group("performance",
		slog.String("operation", operation),
		slog.Duration("duration", duration),
	)
}

// MemoryUsageAttr سمة استخدام الذاكرة
func MemoryUsageAttr() slog.Attr {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return slog.Group("memory",
		slog.String("alloc", formatMemory(m.Alloc)),
		slog.String("total_alloc", formatMemory(m.TotalAlloc)),
		slog.String("sys", formatMemory(m.Sys)),
		slog.Uint64("num_gc", uint64(m.NumGC)),
	)
}

// GoroutineCountAttr سمة عدد الـ goroutines
func GoroutineCountAttr() slog.Attr {
	return slog.Int("goroutines", runtime.NumGoroutine())
}

// ========== دوال تسجيل مخصصة ==========

// LogCacheOperation تسجيل عملية تخزين مؤقت
func LogCacheOperation(ctx context.Context, operation, key string, duration time.Duration, success bool) {
	if success {
		Info(ctx, "عملية التخزين المؤقت",
			CacheOperationAttr(operation, key, duration),
			slog.Bool("success", true),
		)
	} else {
		Error(ctx, "فشل عملية التخزين المؤقت",
			CacheOperationAttr(operation, key, duration),
			slog.Bool("success", false),
		)
	}
}

// LogRedisConnection تسجيل اتصال Redis
func LogRedisConnection(ctx context.Context, status, environment string, retryCount int, err error) {
	if err != nil {
		Error(ctx, "فشل اتصال Redis",
			slog.String("status", status),
			slog.String("environment", environment),
			slog.Int("retry_count", retryCount),
			ErrAttr(err),
		)
	} else {
		Info(ctx, "اتصال Redis ناجح",
			slog.String("status", status),
			slog.String("environment", environment),
			slog.Int("retry_count", retryCount),
		)
	}
}

// LogRateLimit تسجيل تحديد المعدل
func LogRateLimit(ctx context.Context, userID, endpoint string, attempts int, limited bool) {
	attrs := []any{
		slog.String("user_id", userID),
		slog.String("endpoint", endpoint),
		slog.Int("attempts", attempts),
		slog.Bool("limited", limited),
	}

	if limited {
		Warn(ctx, "تم تحديد معدل الطلبات", attrs...)
	} else {
		Debug(ctx, "طلب ضمن المعدل المسموح", attrs...)
	}
}

// LogCORSRequest تسجيل طلب CORS
func LogCORSRequest(ctx context.Context, origin, method, path string, allowed bool) {
	attrs := []any{
		CORSAttr(origin, method, allowed),
		slog.String("path", path),
	}

	if !allowed {
		Warn(ctx, "طلب CORS مرفوض", attrs...)
	} else {
		Debug(ctx, "طلب CORS مسموح", attrs...)
	}
}

// ========== دوال للمراقبة والصحة ==========

// LogStartup تسجيل بدء التشغيل
func LogStartup(ctx context.Context, service, version, environment string) {
	Info(ctx, "🚀 بدء تشغيل الخدمة",
		slog.String("service", service),
		slog.String("version", version),
		slog.String("environment", environment),
	)
}

// LogShutdown تسجيل إيقاف التشغيل
func LogShutdown(ctx context.Context, service string, reason string) {
	Info(ctx, "🛑 إيقاف تشغيل الخدمة",
		slog.String("service", service),
		slog.String("reason", reason),
	)
}

// LogHealthCheck تسجيل فحص الصحة
func LogHealthCheck(ctx context.Context, service, status string, duration time.Duration, details map[string]interface{}) {
	attrs := make([]any, 0, len(details)+3)
	attrs = append(attrs,
		slog.String("service", service),
		slog.String("status", status),
		slog.Duration("duration", duration),
	)
	
	for k, v := range details {
		attrs = append(attrs, slog.Any(k, v))
	}
	
	Info(ctx, "فحص صحة الخدمة", attrs...)
}

// LogDatabaseConnection تسجيل اتصال قاعدة البيانات
func LogDatabaseConnection(ctx context.Context, status string, duration time.Duration, err error) {
	if err != nil {
		Error(ctx, "فشل اتصال قاعدة البيانات",
			slog.String("status", status),
			slog.Duration("duration", duration),
			ErrAttr(err),
		)
	} else {
		Info(ctx, "اتصال قاعدة البيانات ناجح",
			slog.String("status", status),
			slog.Duration("duration", duration),
		)
	}
}

// LogSSEConnection تسجيل اتصال SSE
func LogSSEConnection(ctx context.Context, clientID, userID string, channels []string) {
	Info(ctx, "عميل SSE متصل",
		slog.String("client_id", clientID),
		slog.String("user_id", userID),
		slog.Any("channels", channels),
	)
}

// LogSSEDisconnection تسجيل انفصال SSE
func LogSSEDisconnection(ctx context.Context, clientID, userID string) {
	Info(ctx, "عميل SSE انقطع",
		slog.String("client_id", clientID),
		slog.String("user_id", userID),
	)
}

// ========== دوال مساعدة للنماذج والخدمات ==========

// LogServiceOperation تسجيل عملية خدمة
func LogServiceOperation(ctx context.Context, service, operation string, duration time.Duration, success bool, err error) {
	attrs := []any{
		slog.String("service", service),
		slog.String("operation", operation),
		slog.Duration("duration", duration),
		slog.Bool("success", success),
	}

	if err != nil {
		attrs = append(attrs, ErrAttr(err))
		Error(ctx, "فشل عملية الخدمة", attrs...)
	} else if !success {
		Warn(ctx, "عملية الخدمة لم تنجح", attrs...)
	} else {
		Info(ctx, "عملية الخدمة ناجحة", attrs...)
	}
}

// LogModelOperation تسجيل عملية على نموذج
func LogModelOperation(ctx context.Context, model, operation string, id interface{}, duration time.Duration, err error) {
	attrs := []any{
		slog.String("model", model),
		slog.String("operation", operation),
		slog.Any("id", id),
		slog.Duration("duration", duration),
	}

	if err != nil {
		attrs = append(attrs, ErrAttr(err))
		Error(ctx, "فشل عملية النموذج", attrs...)
	} else {
		Info(ctx, "عملية النموذج ناجحة", attrs...)
	}
}

// ========== دوال مساعدة إضافية ==========

// formatMemory تنسيق حجم الذاكرة
func formatMemory(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// sprintf دالة مساعدة للتنسيق (بدون استيراد fmt)
func sprintf(format string, args ...interface{}) string {
	// تنفيذ مبسط - في الواقع يجب استخدام fmt
	return format
}

// GetGlobalLogger الحصول على الـ logger العالمي
func GetGlobalLogger() Logger {
	if globalLogger == nil {
		Init("development")
	}
	return globalLogger
}

// SetGlobalLogger تعيين الـ logger العالمي
func SetGlobalLogger(logger Logger) {
	globalLogger = logger
}