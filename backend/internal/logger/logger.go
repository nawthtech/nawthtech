package logger

import (
	"log/slog"
	"os"
	"time"
)

// ========== أنواع وواجهات ==========

// Logger واجهة للسجلات
type Logger interface {
	Info(message string, fields map[string]interface{})
	Warn(message string, fields map[string]interface{})
	Error(message string, fields map[string]interface{})
}

// DefaultLogger تطبيق افتراضي للسجلات (للتوافق مع الكود القديم)
type DefaultLogger struct{}

// ========== متغيرات عامة ==========

var (
	// logInstance للواجهة القديمة
	logInstance Logger = &DefaultLogger{}

	// معالجات slog
	stdoutHandler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	
	stdoutHandlerWithSource = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelInfo,
	})

	stderrHandler = slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelWarn,
	})
	
	stderrHandlerWithSource = slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelWarn,
	})

	// Loggers الرئيسية
	Stdout = slog.New(stdoutHandler)
	StdoutWithSource = slog.New(stdoutHandlerWithSource)
	Stderr = slog.New(stderrHandler)
	StderrWithSource = slog.New(stderrHandlerWithSource)
)

// ========== التهيئة والإعداد ==========

// Init تهيئة النظام (للتوافق مع الكود القديم)
func Init(env string) {
	if env == "production" {
		// في الإنتاج، استخدام JSON handler
		logInstance = &DefaultLogger{}
	} else {
		// في التطوير، استخدام text handler
		logInstance = &DefaultLogger{}
	}
}

// InitLogger تهيئة متقدمة للنظام (مستحسن)
func InitLogger(env string, level slog.Level) {
	opts := &slog.HandlerOptions{
		Level: level,
	}

	if env == "development" {
		opts.AddSource = true
		// استخدام TextHandler في التطوير للقراءة السهلة
		Stdout = slog.New(slog.NewTextHandler(os.Stdout, opts))
		Stderr = slog.New(slog.NewTextHandler(os.Stderr, opts))
	} else {
		// استخدام JSONHandler في الإنتاج
		Stdout = slog.New(slog.NewJSONHandler(os.Stdout, opts))
		Stderr = slog.New(slog.NewJSONHandler(os.Stderr, opts))
	}
}

// ========== دوال الواجهة القديمة (للتوافق) ==========

func (l *DefaultLogger) Info(message string, fields map[string]interface{}) {
	attrs := make([]any, 0, len(fields)*2)
	for k, v := range fields {
		attrs = append(attrs, slog.Any(k, v))
	}
	Stdout.Info(message, attrs...)
}

func (l *DefaultLogger) Warn(message string, fields map[string]interface{}) {
	attrs := make([]any, 0, len(fields)*2)
	for k, v := range fields {
		attrs = append(attrs, slog.Any(k, v))
	}
	Stderr.Warn(message, attrs...)
}

func (l *DefaultLogger) Error(message string, fields map[string]interface{}) {
	attrs := make([]any, 0, len(fields)*2)
	for k, v := range fields {
		attrs = append(attrs, slog.Any(k, v))
	}
	Stderr.Error(message, attrs...)
}

// Info تسجيل معلومات (واجهة قديمة)
func Info(message string, fields map[string]interface{}) {
	logInstance.Info(message, fields)
}

// Warn تسجيل تحذير (واجهة قديمة)
func Warn(message string, fields map[string]interface{}) {
	logInstance.Warn(message, fields)
}

// Error تسجيل خطأ (واجهة قديمة)
func Error(message string, fields map[string]interface{}) {
	logInstance.Error(message, fields)
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

// ========== دوال مساعدة للتخزين المؤقت ==========

// CacheOperationAttr سمات عملية التخزين المؤقت
func CacheOperationAttr(operation, key string, duration time.Duration) slog.Attr {
	return slog.Group("cache",
		slog.String("operation", operation),
		slog.String("key", key),
		slog.Duration("duration", duration),
		TimestampAttr(),
	)
}

// CacheHitAttr سمة نجاح التخزين المؤقت
func CacheHitAttr(key string, hit bool) slog.Attr {
	return slog.Group("cache",
		slog.String("key", key),
		slog.Bool("hit", hit),
		slog.String("operation", "get"),
		TimestampAttr(),
	)
}

// CacheErrorAttr سمة خطأ التخزين المؤقت
func CacheErrorAttr(operation, key string, err error) slog.Attr {
	return slog.Group("cache_error",
		slog.String("operation", operation),
		slog.String("key", key),
		ErrAttr(err),
		TimestampAttr(),
	)
}

// CacheStatsAttr سمة إحصائيات التخزين المؤقت
func CacheStatsAttr(keysCount int64, hitRate float64, memoryUsage string) slog.Attr {
	return slog.Group("cache_stats",
		slog.Int64("keys_count", keysCount),
		slog.Float64("hit_rate", hitRate),
		slog.String("memory_usage", memoryUsage),
		TimestampAttr(),
	)
}

// RedisConnectionAttr سمة اتصال Redis
func RedisConnectionAttr(status string, environment string, retryCount int) slog.Attr {
	return slog.Group("redis_connection",
		slog.String("status", status),
		slog.String("environment", environment),
		slog.Int("retry_count", retryCount),
		TimestampAttr(),
	)
}

// ========== دوال مساعدة للطلبات والشبكة ==========

// RequestAttr سمات الطلب
func RequestAttr(method, path, requestID string, statusCode int, duration time.Duration) slog.Attr {
	return slog.Group("request",
		slog.String("method", method),
		slog.String("path", path),
		slog.String("request_id", requestID),
		slog.Int("status_code", statusCode),
		slog.Duration("duration", duration),
		TimestampAttr(),
	)
}

// CORSAttr سمة CORS
func CORSAttr(origin, method string, allowed bool) slog.Attr {
	return slog.Group("cors",
		slog.String("origin", origin),
		slog.String("method", method),
		slog.Bool("allowed", allowed),
		TimestampAttr(),
	)
}

// UserActionAttr سمة إجراء المستخدم
func UserActionAttr(userID, action, resource string) slog.Attr {
	return slog.Group("user_action",
		slog.String("user_id", userID),
		slog.String("action", action),
		slog.String("resource", resource),
		TimestampAttr(),
	)
}

// DatabaseQueryAttr سمة استعلام قاعدة البيانات
func DatabaseQueryAttr(operation, table string, duration time.Duration, rowsAffected int64) slog.Attr {
	return slog.Group("database",
		slog.String("operation", operation),
		slog.String("table", table),
		slog.Duration("duration", duration),
		slog.Int64("rows_affected", rowsAffected),
		TimestampAttr(),
	)
}

// PerformanceAttr سمة الأداء
func PerformanceAttr(operation string, duration time.Duration, memoryUsage string) slog.Attr {
	return slog.Group("performance",
		slog.String("operation", operation),
		slog.Duration("duration", duration),
		slog.String("memory_usage", memoryUsage),
		TimestampAttr(),
	)
}

// ========== دوال تسجيل مخصصة ==========

// LogCacheOperation تسجيل عملية تخزين مؤقت
func LogCacheOperation(operation, key string, duration time.Duration, success bool) {
	if success {
		Stdout.Info("عملية التخزين المؤقت",
			CacheOperationAttr(operation, key, duration),
			slog.Bool("success", true),
		)
	} else {
		Stderr.Error("فشل عملية التخزين المؤقت",
			CacheOperationAttr(operation, key, duration),
			slog.Bool("success", false),
		)
	}
}

// LogRedisConnection تسجيل اتصال Redis
func LogRedisConnection(status, environment string, retryCount int, err error) {
	if err != nil {
		Stderr.Error("فشل اتصال Redis",
			RedisConnectionAttr(status, environment, retryCount),
			ErrAttr(err),
		)
	} else {
		Stdout.Info("اتصال Redis ناجح",
			RedisConnectionAttr(status, environment, retryCount),
		)
	}
}

// LogRateLimit تسجيل تحديد المعدل
func LogRateLimit(userID, endpoint string, attempts int, limited bool) {
	attrs := slog.Group("rate_limit",
		slog.String("user_id", userID),
		slog.String("endpoint", endpoint),
		slog.Int("attempts", attempts),
		slog.Bool("limited", limited),
		TimestampAttr(),
	)

	if limited {
		Stderr.Warn("تم تحديد معدل الطلبات", attrs)
	} else {
		Stdout.Debug("طلب ضمن المعدل المسموح", attrs)
	}
}

// LogCORSRequest تسجيل طلب CORS
func LogCORSRequest(origin, method, path string, allowed bool) {
	level := slog.LevelDebug
	if !allowed {
		level = slog.LevelWarn
	}

	Stdout.Log(nil, level, "طلب CORS",
		CORSAttr(origin, method, allowed),
		slog.String("path", path),
	)
}

// ========== دوال للمراقبة والصحة ==========

// LogStartup تسجيل بدء التشغيل
func LogStartup(service, version, environment string) {
	Stdout.Info("🚀 بدء تشغيل الخدمة",
		slog.String("service", service),
		slog.String("version", version),
		slog.String("environment", environment),
		slog.String("timestamp", time.Now().Format(time.RFC3339)),
	)
}

// LogShutdown تسجيل إيقاف التشغيل
func LogShutdown(service string, reason string) {
	Stdout.Info("🛑 إيقاف تشغيل الخدمة",
		slog.String("service", service),
		slog.String("reason", reason),
		slog.String("timestamp", time.Now().Format(time.RFC3339)),
	)
}

// LogHealthCheck تسجيل فحص الصحة
func LogHealthCheck(service, status string, duration time.Duration, details map[string]interface{}) {
	attrs := make([]any, 0, len(details)+3)
	attrs = append(attrs,
		slog.String("service", service),
		slog.String("status", status),
		slog.Duration("duration", duration),
		TimestampAttr(),
	)
	
	for k, v := range details {
		attrs = append(attrs, slog.Any(k, v))
	}
	
	Stdout.Info("فحص صحة الخدمة", attrs...)
}