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
	Stdout       *slog.Logger
	Stderr       *slog.Logger
	globalLogger Logger
)

// ========== التهيئة والإعداد ==========

func Init(env string) {
	level := slog.LevelInfo
	if env == "development" {
		level = slog.LevelDebug
	}

	opts := &slog.HandlerOptions{
		Level: level,
	}

	if env == "development" {
		Stdout = slog.New(slog.NewTextHandler(os.Stdout, opts))
		Stderr = slog.New(slog.NewTextHandler(os.Stderr, opts))
	} else {
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

func Debug(ctx context.Context, msg string, args ...any) {
	if globalLogger == nil {
		Init("development")
	}
	globalLogger.Debug(ctx, msg, args...)
}

func Info(ctx context.Context, msg string, args ...any) {
	if globalLogger == nil {
		Init("development")
	}
	globalLogger.Info(ctx, msg, args...)
}

func Warn(ctx context.Context, msg string, args ...any) {
	if globalLogger == nil {
		Init("development")
	}
	globalLogger.Warn(ctx, msg, args...)
}

func Error(ctx context.Context, msg string, args ...any) {
	if globalLogger == nil {
		Init("development")
	}
	globalLogger.Error(ctx, msg, args...)
}

func With(args ...any) Logger {
	if globalLogger == nil {
		Init("development")
	}
	return globalLogger.With(args...)
}

// ========== دوال مساعدة أساسية ==========

func ErrAttr(err error) slog.Attr {
	if err == nil {
		return slog.String("error", "")
	}
	return slog.String("error", err.Error())
}

func DurationAttr(duration time.Duration) slog.Attr {
	return slog.Duration("duration", duration)
}

func TimestampAttr() slog.Attr {
	return slog.String("timestamp", time.Now().Format(time.RFC3339))
}

func RequestIDAttr(requestID string) slog.Attr {
	return slog.String("request_id", requestID)
}

func UserIDAttr(userID string) slog.Attr {
	return slog.String("user_id", userID)
}

// ========== دوال مساعدة للطلبات والشبكة ==========

func RequestAttr(method, path string, statusCode int, duration time.Duration) slog.Attr {
	return slog.Group("request",
		slog.String("method", method),
		slog.String("path", path),
		slog.Int("status_code", statusCode),
		slog.Duration("duration", duration),
	)
}

func CORSAttr(origin, method string, allowed bool) slog.Attr {
	return slog.Group("cors",
		slog.String("origin", origin),
		slog.String("method", method),
		slog.Bool("allowed", allowed),
	)
}

func DatabaseQueryAttr(operation, table string, duration time.Duration, rowsAffected int64) slog.Attr {
	return slog.Group("database",
		slog.String("operation", operation),
		slog.String("table", table),
		slog.Duration("duration", duration),
		slog.Int64("rows_affected", rowsAffected),
	)
}

// ========== دوال مساعدة للأداء والذاكرة ==========

func PerformanceAttr(operation string, duration time.Duration) slog.Attr {
	return slog.Group("performance",
		slog.String("operation", operation),
		slog.Duration("duration", duration),
	)
}

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

func GoroutineCountAttr() slog.Attr {
	return slog.Int("goroutines", runtime.NumGoroutine())
}

// ========== دوال تسجيل مخصصة ==========

func LogStartup(ctx context.Context, service, version, environment string) {
	Info(ctx, "🚀 بدء تشغيل الخدمة",
		slog.String("service", service),
		slog.String("version", version),
		slog.String("environment", environment),
		slog.String("database", "Cloudflare D1"),
	)
}

func LogShutdown(ctx context.Context, service string, reason string) {
	Info(ctx, "🛑 إيقاف تشغيل الخدمة",
		slog.String("service", service),
		slog.String("reason", reason),
	)
}

func LogHealthCheck(ctx context.Context, service, status string, duration time.Duration, details map[string]interface{}) {
	attrs := []any{
		slog.String("service", service),
		slog.String("status", status),
		slog.Duration("duration", duration),
	}

	for k, v := range details {
		attrs = append(attrs, slog.Any(k, v))
	}

	Info(ctx, "فحص صحة الخدمة", attrs...)
}

func LogDatabaseConnection(ctx context.Context, status string, duration time.Duration, err error) {
	if err != nil {
		Error(ctx, "❌ فشل اتصال قاعدة البيانات",
			DatabaseQueryAttr("connect", "Cloudflare D1", duration, 0),
			ErrAttr(err),
		)
	} else {
		Info(ctx, "✅ اتصال قاعدة البيانات ناجح",
			DatabaseQueryAttr("connect", "Cloudflare D1", duration, 0),
		)
	}
}

// ========== دوال تسجيل للعمليات العامة ==========

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

func LogRequest(ctx context.Context, method, path string, statusCode int, duration time.Duration, userID string) {
	attrs := []any{RequestAttr(method, path, statusCode, duration)}
	if userID != "" {
		attrs = append(attrs, UserIDAttr(userID))
	}

	if statusCode >= 500 {
		Error(ctx, "طلب HTTP فاشل", attrs...)
	} else if statusCode >= 400 {
		Warn(ctx, "طلب HTTP برفض", attrs...)
	} else {
		Info(ctx, "طلب HTTP ناجح", attrs...)
	}
}

func LogCORSRequest(ctx context.Context, origin, method, path string, allowed bool) {
	attrs := []any{CORSAttr(origin, method, allowed), slog.String("path", path)}

	if !allowed {
		Warn(ctx, "طلب CORS مرفوض", attrs...)
	} else {
		Debug(ctx, "طلب CORS مسموح", attrs...)
	}
}

// ========== دوال مساعدة إضافية ==========

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
	return string(rune(float64(bytes) / float64(div))) + " " + string("KMGTPE"[exp]) + "B"
}

func GetGlobalLogger() Logger {
	if globalLogger == nil {
		Init("development")
	}
	return globalLogger
}

func SetGlobalLogger(logger Logger) {
	globalLogger = logger
}

// ========== دوال رموز تعبيرية ==========

func WithSuccess(logger Logger) Logger { return logger.With(slog.String("status", "✅")) }
func WithWarning(logger Logger) Logger { return logger.With(slog.String("status", "⚠️")) }
func WithError(logger Logger) Logger   { return logger.With(slog.String("status", "❌")) }
func WithInfo(logger Logger) Logger    { return logger.With(slog.String("status", "ℹ️")) }
func WithDebug(logger Logger) Logger   { return logger.With(slog.String("status", "🐛")) }