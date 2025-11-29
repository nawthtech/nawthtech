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
		return slog.String("error", "")
	}
	return slog.String("error", err.Error())
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

// DatabaseQueryAttr سمة استعلام قاعدة البيانات
func DatabaseQueryAttr(operation, collection string, duration time.Duration, documentsAffected int64) slog.Attr {
	return slog.Group("database",
		slog.String("operation", operation),
		slog.String("collection", collection),
		slog.Duration("duration", duration),
		slog.Int64("documents_affected", documentsAffected),
	)
}

// MongoDBConnectionAttr سمة اتصال MongoDB
func MongoDBConnectionAttr(status string, duration time.Duration, err error) slog.Attr {
	attrs := []slog.Attr{
		slog.String("status", status),
		slog.Duration("duration", duration),
		slog.String("database", "MongoDB"),
	}
	
	if err != nil {
		attrs = append(attrs, ErrAttr(err))
	}
	
	return slog.Group("mongodb_connection", attrs...)
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

// LogStartup تسجيل بدء التشغيل
func LogStartup(ctx context.Context, service, version, environment string) {
	Info(ctx, "🚀 بدء تشغيل الخدمة",
		slog.String("service", service),
		slog.String("version", version),
		slog.String("environment", environment),
		slog.String("database", "MongoDB"),
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
		Error(ctx, "❌ فشل اتصال قاعدة البيانات",
			MongoDBConnectionAttr(status, duration, err),
		)
	} else {
		Info(ctx, "✅ اتصال قاعدة البيانات ناجح",
			MongoDBConnectionAttr(status, duration, nil),
		)
	}
}

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
		Error(ctx, "❌ فشل عملية الخدمة", attrs...)
	} else if !success {
		Warn(ctx, "⚠️ عملية الخدمة لم تنجح", attrs...)
	} else {
		Info(ctx, "✅ عملية الخدمة ناجحة", attrs...)
	}
}

// LogMongoDBOperation تسجيل عملية MongoDB
func LogMongoDBOperation(ctx context.Context, operation, collection string, duration time.Duration, documentsAffected int64, err error) {
	attrs := []any{
		DatabaseQueryAttr(operation, collection, duration, documentsAffected),
		slog.String("database", "MongoDB"),
	}

	if err != nil {
		attrs = append(attrs, ErrAttr(err))
		Error(ctx, "❌ فشل عملية قاعدة البيانات", attrs...)
	} else {
		Debug(ctx, "عملية قاعدة البيانات ناجحة", attrs...)
	}
}

// LogCloudinaryOperation تسجيل عملية Cloudinary
func LogCloudinaryOperation(ctx context.Context, operation, filename string, duration time.Duration, success bool, err error) {
	attrs := []any{
		slog.String("service", "cloudinary"),
		slog.String("operation", operation),
		slog.String("filename", filename),
		slog.Duration("duration", duration),
		slog.Bool("success", success),
	}

	if err != nil {
		attrs = append(attrs, ErrAttr(err))
		Error(ctx, "❌ فشل عملية Cloudinary", attrs...)
	} else if !success {
		Warn(ctx, "⚠️ عملية Cloudinary لم تنجح", attrs...)
	} else {
		Info(ctx, "✅ عملية Cloudinary ناجحة", attrs...)
	}
}

// LogAuthentication تسجيل عملية المصادقة
func LogAuthentication(ctx context.Context, operation, userID string, success bool, err error) {
	attrs := []any{
		slog.String("operation", operation),
		slog.String("user_id", userID),
		slog.Bool("success", success),
	}

	if err != nil {
		attrs = append(attrs, ErrAttr(err))
		Warn(ctx, "🔐 فشل عملية المصادقة", attrs...)
	} else if !success {
		Warn(ctx, "🔐 عملية المصادقة لم تنجح", attrs...)
	} else {
		Info(ctx, "🔐 عملية المصادقة ناجحة", attrs...)
	}
}

// LogRequest تسجيل طلب HTTP
func LogRequest(ctx context.Context, method, path string, statusCode int, duration time.Duration, userID string) {
	attrs := []any{
		RequestAttr(method, path, statusCode, duration),
	}

	if userID != "" {
		attrs = append(attrs, UserIDAttr(userID))
	}

	// تسجيل بناءً على حالة الاستجابة
	if statusCode >= 500 {
		Error(ctx, "طلب HTTP فاشل", attrs...)
	} else if statusCode >= 400 {
		Warn(ctx, "طلب HTTP برفض", attrs...)
	} else {
		Info(ctx, "طلب HTTP ناجح", attrs...)
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

// ========== دوال مساعدة إضافية ==========

// formatMemory تنسيق حجم الذاكرة
func formatMemory(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return string(rune(bytes)) + " B"
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return string(rune(float64(bytes)/float64(div))) + " " + string("KMGTPE"[exp]) + "B"
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

// ========== دوال بادئات الرموز التعبيرية ==========

// WithSuccess إضافة رمز نجاح
func WithSuccess(logger Logger) Logger {
	return logger.With(slog.String("status", "✅"))
}

// WithWarning إضافة رمز تحذير
func WithWarning(logger Logger) Logger {
	return logger.With(slog.String("status", "⚠️"))
}

// WithError إضافة رمز خطأ
func WithError(logger Logger) Logger {
	return logger.With(slog.String("status", "❌"))
}

// WithInfo إضافة رمز معلومات
func WithInfo(logger Logger) Logger {
	return logger.With(slog.String("status", "ℹ️"))
}

// WithDebug إضافة رمز تصحيح
func WithDebug(logger Logger) Logger {
	return logger.With(slog.String("status", "🐛"))
}