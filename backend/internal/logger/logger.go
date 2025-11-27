package logger

import (
	"os"
	"log/slog"
	"time"
)

var (
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

	// sends logs to stdout
	Stdout = slog.New(stdoutHandler)
	// sends logs to stdout with source info
	StdoutWithSource = slog.New(stdoutHandlerWithSource)

	// sends logs to stderr
	Stderr = slog.New(stderrHandler)
	// sends logs to stderr with source info
	StderrWithSource = slog.New(stderrHandlerWithSource)
)

// ErrAttr دالة مساعدة لإرجاع سمة الخطأ
func ErrAttr(err error) slog.Attr {
	return slog.Any("error", err)
}

// ErrorsAttr دالة مساعدة لإرجاع سمة الأخطاء المتعددة
func ErrorsAttr(errors ...error) slog.Attr {
	return slog.Any("errors", errors)
}

// ========== دوال مساعدة للتخزين المؤقت ==========

// CacheOperationAttr سمات عملية التخزين المؤقت
func CacheOperationAttr(operation, key string, duration time.Duration) slog.Attr {
	return slog.Group("cache",
		slog.String("operation", operation),
		slog.String("key", key),
		slog.Duration("duration", duration),
		slog.String("timestamp", time.Now().Format(time.RFC3339)),
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
		slog.String("error", err.Error()),
		slog.String("timestamp", time.Now().Format(time.RFC3339)),
	)
}

// CacheStatsAttr سمة إحصائيات التخزين المؤقت
func CacheStatsAttr(keysCount int64, hitRate float64, memoryUsage string) slog.Attr {
	return slog.Group("cache_stats",
		slog.Int64("keys_count", keysCount),
		slog.Float64("hit_rate", hitRate),
		slog.String("memory_usage", memoryUsage),
		slog.String("timestamp", time.Now().Format(time.RFC3339)),
	)
}

// RedisConnectionAttr سمة اتصال Redis
func RedisConnectionAttr(status string, environment string, retryCount int) slog.Attr {
	return slog.Group("redis_connection",
		slog.String("status", status),
		slog.String("environment", environment),
		slog.Int("retry_count", retryCount),
		slog.String("timestamp", time.Now().Format(time.RFC3339)),
	)
}

// ========== دوال مساعدة للخدمات ==========

// ServiceOperationAttr سمات عملية الخدمة
func ServiceOperationAttr(operation, serviceID, sellerID string) slog.Attr {
	return slog.Group("service",
		slog.String("operation", operation),
		slog.String("service_id", serviceID),
		slog.String("seller_id", sellerID),
		slog.String("timestamp", time.Now().Format(time.RFC3339)),
	)
}

// ServiceCreationAttr سمة إنشاء خدمة
func ServiceCreationAttr(serviceID, title, category string, price float64) slog.Attr {
	return slog.Group("service_creation",
		slog.String("service_id", serviceID),
		slog.String("title", title),
		slog.String("category", category),
		slog.Float64("price", price),
		slog.String("timestamp", time.Now().Format(time.RFC3339)),
	)
}

// ServiceSearchAttr سمة بحث الخدمات
func ServiceSearchAttr(query, category string, resultsCount int, duration time.Duration) slog.Attr {
	return slog.Group("service_search",
		slog.String("query", query),
		slog.String("category", category),
		slog.Int("results_count", resultsCount),
		slog.Duration("duration", duration),
		slog.String("timestamp", time.Now().Format(time.RFC3339)),
	)
}

// ServiceRatingAttr سمة تقييم الخدمة
func ServiceRatingAttr(serviceID, userID string, rating int, previousRating float64) slog.Attr {
	return slog.Group("service_rating",
		slog.String("service_id", serviceID),
		slog.String("user_id", userID),
		slog.Int("rating", rating),
		slog.Float64("previous_rating", previousRating),
		slog.String("timestamp", time.Now().Format(time.R3339)),
	)
}

// ServiceAnalyticsAttr سمة تحليلات الخدمة
func ServiceAnalyticsAttr(serviceID, period string, views, orders int, revenue float64) slog.Attr {
	return slog.Group("service_analytics",
		slog.String("service_id", serviceID),
		slog.String("period", period),
		slog.Int("views", views),
		slog.Int("orders", orders),
		slog.Float64("revenue", revenue),
		slog.String("timestamp", time.Now().Format(time.RFC3339)),
	)
}

// ========== دوال مساعدة عامة ==========

// RequestAttr سمات الطلب
func RequestAttr(method, path, requestID string, statusCode int, duration time.Duration) slog.Attr {
	return slog.Group("request",
		slog.String("method", method),
		slog.String("path", path),
		slog.String("request_id", requestID),
		slog.Int("status_code", statusCode),
		slog.Duration("duration", duration),
		slog.String("timestamp", time.Now().Format(time.RFC3339)),
	)
}

// UserActionAttr سمة إجراء المستخدم
func UserActionAttr(userID, action, resource string) slog.Attr {
	return slog.Group("user_action",
		slog.String("user_id", userID),
		slog.String("action", action),
		slog.String("resource", resource),
		slog.String("timestamp", time.Now().Format(time.RFC3339)),
	)
}

// DatabaseQueryAttr سمة استعلام قاعدة البيانات
func DatabaseQueryAttr(operation, table string, duration time.Duration, rowsAffected int64) slog.Attr {
	return slog.Group("database",
		slog.String("operation", operation),
		slog.String("table", table),
		slog.Duration("duration", duration),
		slog.Int64("rows_affected", rowsAffected),
		slog.String("timestamp", time.Now().Format(time.RFC3339)),
	)
}

// PerformanceAttr سمة الأداء
func PerformanceAttr(operation string, duration time.Duration, memoryUsage string) slog.Attr {
	return slog.Group("performance",
		slog.String("operation", operation),
		slog.Duration("duration", duration),
		slog.String("memory_usage", memoryUsage),
		slog.String("timestamp", time.Now().Format(time.RFC3339)),
	)
}

// ========== دوال تسجيل مخصصة ==========

// LogCacheOperation تسجيل عملية تخزين مؤقت
func LogCacheOperation(logger *slog.Logger, operation, key string, duration time.Duration, success bool) {
	if success {
		logger.Info("عملية التخزين المؤقت",
			CacheOperationAttr(operation, key, duration),
			slog.Bool("success", true),
		)
	} else {
		logger.Error("فشل عملية التخزين المؤقت",
			CacheOperationAttr(operation, key, duration),
			slog.Bool("success", false),
		)
	}
}

// LogServiceCreation تسجيل إنشاء خدمة
func LogServiceCreation(logger *slog.Logger, serviceID, title, category string, price float64, sellerID string) {
	logger.Info("تم إنشاء خدمة جديدة",
		ServiceCreationAttr(serviceID, title, category, price),
		slog.String("seller_id", sellerID),
	)
}

// LogServiceSearch تسجيل بحث الخدمات
func LogServiceSearch(logger *slog.Logger, query, category string, resultsCount int, duration time.Duration, userID string) {
	logger.Info("بحث في الخدمات",
		ServiceSearchAttr(query, category, resultsCount, duration),
		slog.String("user_id", userID),
	)
}

// LogRedisConnection تسجيل اتصال Redis
func LogRedisConnection(logger *slog.Logger, status, environment string, retryCount int, err error) {
	if err != nil {
		logger.Error("فشل اتصال Redis",
			RedisConnectionAttr(status, environment, retryCount),
			ErrAttr(err),
		)
	} else {
		logger.Info("اتصال Redis ناجح",
			RedisConnectionAttr(status, environment, retryCount),
		)
	}
}

// LogRateLimit تسجيل تحديد المعدل
func LogRateLimit(logger *slog.Logger, userID, endpoint string, attempts int, limited bool) {
	attrs := slog.Group("rate_limit",
		slog.String("user_id", userID),
		slog.String("endpoint", endpoint),
		slog.Int("attempts", attempts),
		slog.Bool("limited", limited),
		slog.String("timestamp", time.Now().Format(time.RFC3339)),
	)

	if limited {
		logger.Warn("تم تحديد معدل الطلبات", attrs)
	} else {
		logger.Debug("طلب ضمن المعدل المسموح", attrs)
	}
}

// ========== دوال للمستويات المختلفة ==========

// DebugCache تسجيل تصحيح للتخزين المؤقت
func DebugCache(logger *slog.Logger, message string, key string, value interface{}) {
	logger.Debug(message,
		slog.String("key", key),
		slog.Any("value", value),
		slog.String("timestamp", time.Now().Format(time.RFC3339)),
	)
}

// InfoService تسجيل معلومات الخدمة
func InfoService(logger *slog.Logger, message, serviceID string, additionalAttrs ...slog.Attr) {
	attrs := make([]any, 0, len(additionalAttrs)+2)
	attrs = append(attrs,
		slog.String("service_id", serviceID),
		slog.String("timestamp", time.Now().Format(time.RFC3339)),
	)
	
	for _, attr := range additionalAttrs {
		attrs = append(attrs, attr)
	}
	
	logger.Info(message, attrs...)
}

// WarnCache تسجيل تحذير للتخزين المؤقت
func WarnCache(logger *slog.Logger, message, key string, reason string) {
	logger.Warn(message,
		slog.String("key", key),
		slog.String("reason", reason),
		slog.String("timestamp", time.Now().Format(time.RFC3339)),
	)
}

// ErrorService تسجيل خطأ في الخدمة
func ErrorService(logger *slog.Logger, message, serviceID string, err error, additionalAttrs ...slog.Attr) {
	attrs := make([]any, 0, len(additionalAttrs)+3)
	attrs = append(attrs,
		slog.String("service_id", serviceID),
		ErrAttr(err),
		slog.String("timestamp", time.Now().Format(time.RFC3339)),
	)
	
	for _, attr := range additionalAttrs {
		attrs = append(attrs, attr)
	}
	
	logger.Error(message, attrs...)
}