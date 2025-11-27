package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

// CacheService واجهة خدمة التخزين المؤقت
type CacheService interface {
	Initialize(ctx context.Context) error
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Get(ctx context.Context, key string) (interface{}, error)
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)
	Expire(ctx context.Context, key string, ttl time.Duration) error
	TTL(ctx context.Context, key string) (time.Duration, error)
	Increment(ctx context.Context, key string, value int64) (int64, error)
	LPush(ctx context.Context, key string, values ...interface{}) error
	LRange(ctx context.Context, key string, start, stop int64) ([]interface{}, error)
	HSet(ctx context.Context, key string, field string, value interface{}) error
	HGet(ctx context.Context, key string, field string) (interface{}, error)
	HGetAll(ctx context.Context, key string) (map[string]interface{}, error)
	HDel(ctx context.Context, key string, fields ...string) error
	Keys(ctx context.Context, pattern string) ([]string, error)
	Flush(ctx context.Context) error
	FlushPattern(ctx context.Context, pattern string) (int64, error)
	GetStats(ctx context.Context) (*CacheStats, error)
	HealthCheck(ctx context.Context) (*HealthStatus, error)
	Close() error
}

// CacheServiceImpl التطبيق الفعلي لخدمة التخزين المؤقت
type CacheServiceImpl struct {
	client      *redis.Client
	isConnected bool
	defaultTTL  time.Duration
	prefix      string
	isRailway   bool
	retryCount  int
	maxRetries  int
	logger      *slog.Logger
}

// CacheStats إحصائيات التخزين المؤقت
type CacheStats struct {
	Status          string `json:"status"`
	KeysCount       int64  `json:"keysCount"`
	UsedMemory      string `json:"usedMemory"`
	ConnectedClients int64  `json:"connectedClients"`
	Hits           int64  `json:"hits"`
	Misses         int64  `json:"misses"`
	HitRate        int    `json:"hitRate"`
	Uptime         int64  `json:"uptime"`
	Environment    string `json:"environment"`
	RetryCount     int    `json:"retryCount"`
}

// HealthStatus حالة صحة الخدمة
type HealthStatus struct {
	Status      string      `json:"status"`
	Message     string      `json:"message"`
	Error       string      `json:"error,omitempty"`
	Environment string      `json:"environment"`
	RetryCount  int         `json:"retryCount"`
	Stats       *CacheStats `json:"stats,omitempty"`
}

// CacheConfig تكوين خدمة التخزين المؤقت
type CacheConfig struct {
	RedisURL      string
	RedisHost     string
	RedisPort     string
	RedisPassword string
	RedisDB       int
	Prefix        string
	DefaultTTL    time.Duration
	MaxRetries    int
}

// NewCacheService إنشاء خدمة تخزين مؤقت جديدة
func NewCacheService(config CacheConfig) CacheService {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	return &CacheServiceImpl{
		isConnected: false,
		defaultTTL:  config.DefaultTTL,
		prefix:      config.Prefix,
		isRailway:   os.Getenv("RAILWAY_ENVIRONMENT") == "true",
		maxRetries:  config.MaxRetries,
		logger:      logger,
	}
}

// Initialize تهيئة خدمة التخزين المؤقت
func (c *CacheServiceImpl) Initialize(ctx context.Context) error {
	c.logger.Info("🚀 تهيئة خدمة التخزين المؤقت...")
	c.logger.Info("🏗️ بيئة Railway", "is_railway", c.isRailway)

	var redisOptions *redis.Options

	// استخدام REDIS_URL إذا كان متوفراً (مطلوب في Railway)
	redisURL := os.Getenv("REDIS_URL")
	if redisURL != "" {
		c.logger.Info("🔗 استخدام REDIS_URL للتوصيل بـ Redis")
		parsedOptions, err := c.parseRedisURL(redisURL)
		if err != nil {
			c.logger.Error("❌ خطأ في تحليل REDIS_URL", "error", err)
			return err
		}
		redisOptions = parsedOptions
	} else {
		// التكوين التقليدي للتطوير المحلي
		redisOptions = &redis.Options{
			Addr:     fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
			Password: os.Getenv("REDIS_PASSWORD"),
			DB:       0,
		}
	}

	// إعدادات إضافية لتحسين الموثوقية
	redisOptions.MaxRetries = c.maxRetries
	redisOptions.MinRetryBackoff = 1 * time.Second
	redisOptions.MaxRetryBackoff = 5 * time.Second
	redisOptions.DialTimeout = 10 * time.Second
	redisOptions.ReadTimeout = 30 * time.Second
	redisOptions.WriteTimeout = 30 * time.Second
	redisOptions.PoolSize = 100
	redisOptions.MinIdleConns = 10

	c.client = redis.NewClient(redisOptions)

	// اختبار الاتصال
	if err := c.client.Ping(ctx).Err(); err != nil {
		c.isConnected = false
		c.logger.Error("❌ فشل في الاتصال بـ Redis", "error", err)
		
		if c.isRailway {
			c.logger.Warn("⚠️ فشل الاتصال بـ Redis في Railway، سيتم العمل بدون تخزين مؤقت")
			return nil // لا نرمي خطأ في Railway
		}
		return fmt.Errorf("فشل في الاتصال بـ Redis: %v", err)
	}

	c.isConnected = true
	c.retryCount = 0
	c.logger.Info("✅ تم الاتصال بـ Redis بنجاح")
	c.logger.Info("📊 حالة الاتصال", "connected", c.isConnected)

	// تسجيل معلومات إضافية في Railway
	if c.isRailway {
		c.logger.Info("🌐 تشغيل في بيئة Railway - Redis جاهز")
	}

	return nil
}

// parseRedisURL تحليل REDIS_URL إلى إعدادات Redis
func (c *CacheServiceImpl) parseRedisURL(redisURL string) (*redis.Options, error) {
	parsedURL, err := url.Parse(redisURL)
	if err != nil {
		return nil, fmt.Errorf("تكوين Redis غير صحيح: %v", err)
	}

	password := ""
	if parsedURL.User != nil {
		password, _ = parsedURL.User.Password()
	}

	db := 0
	if parsedURL.Path != "" && len(parsedURL.Path) > 1 {
		dbStr := parsedURL.Path[1:]
		if dbInt, err := strconv.Atoi(dbStr); err == nil {
			db = dbInt
		}
	}

	options := &redis.Options{
		Addr:     parsedURL.Host,
		Password: password,
		DB:       db,
	}

	// إعداد TLS إذا كان البروتوكول rediss
	if parsedURL.Scheme == "rediss" {
		options.TLSConfig = &tls.Config{
			ServerName: parsedURL.Hostname(),
		}
	}

	return options, nil
}

// Set إضافة مفتاح وقيمة إلى التخزين المؤقت
func (c *CacheServiceImpl) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	if !c.isConnected {
		return fmt.Errorf("غير متصل بـ Redis")
	}

	prefixedKey := c.prefix + key
	serializedValue, err := c.serializeValue(value)
	if err != nil {
		return fmt.Errorf("فشل في تسلسل القيمة: %v", err)
	}

	if ttl > 0 {
		err = c.client.Set(ctx, prefixedKey, serializedValue, ttl).Err()
	} else {
		err = c.client.Set(ctx, prefixedKey, serializedValue, c.defaultTTL).Err()
	}

	if err != nil {
		c.logger.Error("❌ خطأ في تخزين البيانات في الكاش", 
			"key", prefixedKey, 
			"error", err,
			"environment", c.getEnvironment())
		return err
	}

	c.logger.Debug("✅ تم تخزين البيانات في الكاش", 
		"key", prefixedKey, 
		"ttl", ttl,
		"environment", c.getEnvironment())
	return nil
}

// Get الحصول على قيمة من التخزين المؤقت
func (c *CacheServiceImpl) Get(ctx context.Context, key string) (interface{}, error) {
	if !c.isConnected {
		return nil, fmt.Errorf("غير متصل بـ Redis")
	}

	prefixedKey := c.prefix + key
	value, err := c.client.Get(ctx, prefixedKey).Result()
	if err == redis.Nil {
		return nil, nil // المفتاح غير موجود
	} else if err != nil {
		c.logger.Error("❌ خطأ في جلب البيانات من الكاش", 
			"key", prefixedKey, 
			"error", err,
			"environment", c.getEnvironment())
		return nil, err
	}

	deserializedValue, err := c.deserializeValue(value)
	if err != nil {
		return nil, fmt.Errorf("فشل في إعادة تسلسل القيمة: %v", err)
	}

	c.logger.Debug("✅ تم جلب البيانات من الكاش", 
		"key", prefixedKey,
		"environment", c.getEnvironment())
	return deserializedValue, nil
}

// Delete حذف مفتاح من التخزين المؤقت
func (c *CacheServiceImpl) Delete(ctx context.Context, key string) error {
	if !c.isConnected {
		return fmt.Errorf("غير متصل بـ Redis")
	}

	prefixedKey := c.prefix + key
	result, err := c.client.Del(ctx, prefixedKey).Result()
	if err != nil {
		c.logger.Error("❌ خطأ في حذف البيانات من الكاش", 
			"key", prefixedKey, 
			"error", err,
			"environment", c.getEnvironment())
		return err
	}

	c.logger.Debug("✅ تم حذف البيانات من الكاش", 
		"key", prefixedKey, 
		"deleted", result > 0,
		"environment", c.getEnvironment())
	return nil
}

// Exists التحقق من وجود مفتاح في التخزين المؤقت
func (c *CacheServiceImpl) Exists(ctx context.Context, key string) (bool, error) {
	if !c.isConnected {
		return false, fmt.Errorf("غير متصل بـ Redis")
	}

	prefixedKey := c.prefix + key
	result, err := c.client.Exists(ctx, prefixedKey).Result()
	if err != nil {
		c.logger.Error("❌ خطأ في التحقق من وجود المفتاح في الكاش", 
			"key", prefixedKey, 
			"error", err,
			"environment", c.getEnvironment())
		return false, err
	}

	return result > 0, nil
}

// Expire تعيين وقت انتهاء الصلاحية للمفتاح
func (c *CacheServiceImpl) Expire(ctx context.Context, key string, ttl time.Duration) error {
	if !c.isConnected {
		return fmt.Errorf("غير متصل بـ Redis")
	}

	prefixedKey := c.prefix + key
	result, err := c.client.Expire(ctx, prefixedKey, ttl).Result()
	if err != nil {
		c.logger.Error("❌ خطأ في تعيين وقت انتهاء الصلاحية", 
			"key", prefixedKey, 
			"ttl", ttl,
			"error", err,
			"environment", c.getEnvironment())
		return err
	}

	c.logger.Debug("✅ تم تعيين وقت انتهاء الصلاحية", 
		"key", prefixedKey, 
		"ttl", ttl, 
		"result", result,
		"environment", c.getEnvironment())
	return nil
}

// TTL الحصول على وقت انتهاء الصلاحية المتبقي
func (c *CacheServiceImpl) TTL(ctx context.Context, key string) (time.Duration, error) {
	if !c.isConnected {
		return -2, fmt.Errorf("غير متصل بـ Redis")
	}

	prefixedKey := c.prefix + key
	ttl, err := c.client.TTL(ctx, prefixedKey).Result()
	if err != nil {
		c.logger.Error("❌ خطأ في الحصول على وقت انتهاء الصلاحية", 
			"key", prefixedKey, 
			"error", err,
			"environment", c.getEnvironment())
		return -2, err
	}

	return ttl, nil
}

// Increment زيادة قيمة رقمية
func (c *CacheServiceImpl) Increment(ctx context.Context, key string, value int64) (int64, error) {
	if !c.isConnected {
		return 0, fmt.Errorf("غير متصل بـ Redis")
	}

	prefixedKey := c.prefix + key
	var result int64
	var err error

	if value == 1 {
		result, err = c.client.Incr(ctx, prefixedKey).Result()
	} else {
		result, err = c.client.IncrBy(ctx, prefixedKey, value).Result()
	}

	if err != nil {
		c.logger.Error("❌ خطأ في زيادة القيمة الرقمية", 
			"key", prefixedKey, 
			"increment", value,
			"error", err,
			"environment", c.getEnvironment())
		return 0, err
	}

	c.logger.Debug("✅ تم زيادة القيمة الرقمية", 
		"key", prefixedKey, 
		"increment", value, 
		"result", result,
		"environment", c.getEnvironment())
	return result, nil
}

// LPush تخزين بيانات في قائمة
func (c *CacheServiceImpl) LPush(ctx context.Context, key string, values ...interface{}) error {
	if !c.isConnected {
		return fmt.Errorf("غير متصل بـ Redis")
	}

	prefixedKey := c.prefix + key
	serializedValues := make([]interface{}, len(values))
	for i, v := range values {
		serialized, err := c.serializeValue(v)
		if err != nil {
			return fmt.Errorf("فشل في تسلسل القيمة: %v", err)
		}
		serializedValues[i] = serialized
	}

	result, err := c.client.LPush(ctx, prefixedKey, serializedValues...).Result()
	if err != nil {
		c.logger.Error("❌ خطأ في إضافة البيانات إلى القائمة", 
			"key", prefixedKey, 
			"error", err,
			"environment", c.getEnvironment())
		return err
	}

	c.logger.Debug("✅ تم إضافة البيانات إلى القائمة", 
		"key", prefixedKey, 
		"count", result,
		"environment", c.getEnvironment())
	return nil
}

// LRange جلب بيانات من القائمة
func (c *CacheServiceImpl) LRange(ctx context.Context, key string, start, stop int64) ([]interface{}, error) {
	if !c.isConnected {
		return nil, fmt.Errorf("غير متصل بـ Redis")
	}

	prefixedKey := c.prefix + key
	values, err := c.client.LRange(ctx, prefixedKey, start, stop).Result()
	if err != nil {
		c.logger.Error("❌ خطأ في جلب البيانات من القائمة", 
			"key", prefixedKey, 
			"error", err,
			"environment", c.getEnvironment())
		return nil, err
	}

	deserializedValues := make([]interface{}, len(values))
	for i, v := range values {
		deserialized, err := c.deserializeValue(v)
		if err != nil {
			return nil, fmt.Errorf("فشل في إعادة تسلسل القيمة: %v", err)
		}
		deserializedValues[i] = deserialized
	}

	c.logger.Debug("✅ تم جلب البيانات من القائمة", 
		"key", prefixedKey, 
		"count", len(deserializedValues),
		"environment", c.getEnvironment())
	return deserializedValues, nil
}

// HSet تخزين بيانات في هاش
func (c *CacheServiceImpl) HSet(ctx context.Context, key string, field string, value interface{}) error {
	if !c.isConnected {
		return fmt.Errorf("غير متصل بـ Redis")
	}

	prefixedKey := c.prefix + key
	serializedValue, err := c.serializeValue(value)
	if err != nil {
		return fmt.Errorf("فشل في تسلسل القيمة: %v", err)
	}

	err = c.client.HSet(ctx, prefixedKey, field, serializedValue).Err()
	if err != nil {
		c.logger.Error("❌ خطأ في تخزين البيانات في الهاش", 
			"key", prefixedKey, 
			"field", field,
			"error", err,
			"environment", c.getEnvironment())
		return err
	}

	c.logger.Debug("✅ تم تخزين البيانات في الهاش", 
		"key", prefixedKey, 
		"field", field,
		"environment", c.getEnvironment())
	return nil
}

// HGet جلب بيانات من الهاش
func (c *CacheServiceImpl) HGet(ctx context.Context, key string, field string) (interface{}, error) {
	if !c.isConnected {
		return nil, fmt.Errorf("غير متصل بـ Redis")
	}

	prefixedKey := c.prefix + key
	value, err := c.client.HGet(ctx, prefixedKey, field).Result()
	if err == redis.Nil {
		return nil, nil // الحقل غير موجود
	} else if err != nil {
		c.logger.Error("❌ خطأ في جلب البيانات من الهاش", 
			"key", prefixedKey, 
			"field", field,
			"error", err,
			"environment", c.getEnvironment())
		return nil, err
	}

	deserializedValue, err := c.deserializeValue(value)
	if err != nil {
		return nil, fmt.Errorf("فشل في إعادة تسلسل القيمة: %v", err)
	}

	c.logger.Debug("✅ تم جلب البيانات من الهاش", 
		"key", prefixedKey, 
		"field", field,
		"environment", c.getEnvironment())
	return deserializedValue, nil
}

// HGetAll جلب جميع بيانات الهاش
func (c *CacheServiceImpl) HGetAll(ctx context.Context, key string) (map[string]interface{}, error) {
	if !c.isConnected {
		return nil, fmt.Errorf("غير متصل بـ Redis")
	}

	prefixedKey := c.prefix + key
	hash, err := c.client.HGetAll(ctx, prefixedKey).Result()
	if err != nil {
		c.logger.Error("❌ خطأ في جلب جميع بيانات الهاش", 
			"key", prefixedKey, 
			"error", err,
			"environment", c.getEnvironment())
		return nil, err
	}

	deserializedHash := make(map[string]interface{})
	for field, value := range hash {
		deserialized, err := c.deserializeValue(value)
		if err != nil {
			return nil, fmt.Errorf("فشل في إعادة تسلسل القيمة: %v", err)
		}
		deserializedHash[field] = deserialized
	}

	c.logger.Debug("✅ تم جلب جميع بيانات الهاش", 
		"key", prefixedKey, 
		"fieldCount", len(deserializedHash),
		"environment", c.getEnvironment())
	return deserializedHash, nil
}

// HDel حذف حقل من الهاش
func (c *CacheServiceImpl) HDel(ctx context.Context, key string, fields ...string) error {
	if !c.isConnected {
		return fmt.Errorf("غير متصل بـ Redis")
	}

	prefixedKey := c.prefix + key
	result, err := c.client.HDel(ctx, prefixedKey, fields...).Result()
	if err != nil {
		c.logger.Error("❌ خطأ في حذف الحقل من الهاش", 
			"key", prefixedKey, 
			"fields", fields,
			"error", err,
			"environment", c.getEnvironment())
		return err
	}

	c.logger.Debug("✅ تم حذف الحقل من الهاش", 
		"key", prefixedKey, 
		"fields", fields, 
		"deleted", result > 0,
		"environment", c.getEnvironment())
	return nil
}

// Keys البحث عن المفاتيح باستخدام النمط
func (c *CacheServiceImpl) Keys(ctx context.Context, pattern string) ([]string, error) {
	if !c.isConnected {
		return nil, fmt.Errorf("غير متصل بـ Redis")
	}

	prefixedPattern := c.prefix + pattern
	keys, err := c.client.Keys(ctx, prefixedPattern).Result()
	if err != nil {
		c.logger.Error("❌ خطأ في البحث عن المفاتيح", 
			"pattern", prefixedPattern, 
			"error", err,
			"environment", c.getEnvironment())
		return nil, err
	}

	// إزالة البادئة من النتائج
	cleanKeys := make([]string, len(keys))
	for i, key := range keys {
		cleanKeys[i] = strings.TrimPrefix(key, c.prefix)
	}

	c.logger.Debug("✅ تم البحث عن المفاتيح", 
		"pattern", prefixedPattern, 
		"count", len(cleanKeys),
		"environment", c.getEnvironment())
	return cleanKeys, nil
}

// Flush مسح جميع البيانات من التخزين المؤقت
func (c *CacheServiceImpl) Flush(ctx context.Context) error {
	if !c.isConnected {
		return fmt.Errorf("غير متصل بـ Redis")
	}

	err := c.client.FlushDB(ctx).Err()
	if err != nil {
		c.logger.Error("❌ خطأ في مسح بيانات التخزين المؤقت", 
			"error", err,
			"environment", c.getEnvironment())
		return err
	}

	c.logger.Info("✅ تم مسح جميع بيانات التخزين المؤقت", 
		"environment", c.getEnvironment())
	return nil
}

// FlushPattern مسح البيانات باستخدام النمط
func (c *CacheServiceImpl) FlushPattern(ctx context.Context, pattern string) (int64, error) {
	if !c.isConnected {
		return 0, fmt.Errorf("غير متصل بـ Redis")
	}

	keysToDelete, err := c.Keys(ctx, pattern)
	if err != nil {
		return 0, err
	}

	if len(keysToDelete) == 0 {
		return 0, nil
	}

	prefixedKeys := make([]string, len(keysToDelete))
	for i, key := range keysToDelete {
		prefixedKeys[i] = c.prefix + key
	}

	result, err := c.client.Del(ctx, prefixedKeys...).Result()
	if err != nil {
		c.logger.Error("❌ خطأ في مسح البيانات باستخدام النمط", 
			"pattern", pattern,
			"error", err,
			"environment", c.getEnvironment())
		return 0, err
	}

	c.logger.Info("✅ تم مسح البيانات باستخدام النمط", 
		"pattern", pattern, 
		"deletedCount", result,
		"environment", c.getEnvironment())
	return result, nil
}

// GetStats الحصول على إحصائيات التخزين المؤقت
func (c *CacheServiceImpl) GetStats(ctx context.Context) (*CacheStats, error) {
	if !c.isConnected {
		return &CacheStats{
			Status:       "disconnected",
			Environment:  c.getEnvironment(),
			RetryCount:   c.retryCount,
		}, nil
	}

	info, err := c.client.Info(ctx).Result()
	if err != nil {
		c.logger.Error("❌ خطأ في الحصول على إحصائيات التخزين المؤقت", 
			"error", err,
			"environment", c.getEnvironment())
		return nil, err
	}

	keysCount, err := c.client.DBSize(ctx).Result()
	if err != nil {
		return nil, err
	}

	stats := &CacheStats{
		Status:          "connected",
		KeysCount:       keysCount,
		UsedMemory:      c.extractUsedMemory(info),
		ConnectedClients: c.extractConnectedClients(info),
		Hits:           c.extractHits(info),
		Misses:         c.extractMisses(info),
		HitRate:        c.calculateHitRate(info),
		Uptime:         c.extractUptime(info),
		Environment:    c.getEnvironment(),
		RetryCount:     c.retryCount,
	}

	return stats, nil
}

// HealthCheck فحص صحة الخدمة
func (c *CacheServiceImpl) HealthCheck(ctx context.Context) (*HealthStatus, error) {
	if !c.isConnected {
		return &HealthStatus{
			Status:      "disconnected",
			Message:     "غير متصل بـ Redis",
			Environment: c.getEnvironment(),
			RetryCount:  c.retryCount,
		}, nil
	}

	if err := c.client.Ping(ctx).Err(); err != nil {
		return &HealthStatus{
			Status:      "unhealthy",
			Message:     "فشل في فحص صحة خدمة التخزين المؤقت",
			Error:       err.Error(),
			Environment: c.getEnvironment(),
			RetryCount:  c.retryCount,
		}, nil
	}

	stats, err := c.GetStats(ctx)
	if err != nil {
		return &HealthStatus{
			Status:      "degraded",
			Message:     "خدمة التخزين المؤقت تعمل ولكن فشل في جلب الإحصائيات",
			Error:       err.Error(),
			Environment: c.getEnvironment(),
			RetryCount:  c.retryCount,
		}, nil
	}

	return &HealthStatus{
		Status:      "healthy",
		Message:     "خدمة التخزين المؤقت تعمل بشكل طبيعي",
		Environment: c.getEnvironment(),
		RetryCount:  c.retryCount,
		Stats:       stats,
	}, nil
}

// Close إغلاق اتصال التخزين المؤقت
func (c *CacheServiceImpl) Close() error {
	if c.client != nil {
		err := c.client.Close()
		c.isConnected = false
		if err != nil {
			c.logger.Error("❌ خطأ في إغلاق اتصال التخزين المؤقت", "error", err)
			return err
		}
		c.logger.Info("✅ تم إغلاق اتصال التخزين المؤقت بنجاح")
		c.logger.Info("🏗️ البيئة", "environment", c.getEnvironment())
	}
	return nil
}

// ========== الدوال المساعدة ==========

func (c *CacheServiceImpl) serializeValue(value interface{}) (string, error) {
	switch v := value.(type) {
	case string:
		return v, nil
	case []byte:
		return string(v), nil
	default:
		jsonData, err := json.Marshal(v)
		if err != nil {
			c.logger.Error("❌ خطأ في تسلسل القيمة", 
				"error", err,
				"environment", c.getEnvironment())
			return "", err
		}
		return string(jsonData), nil
	}
}

func (c *CacheServiceImpl) deserializeValue(value string) (interface{}, error) {
	// محاولة تحليل JSON
	var result interface{}
	if err := json.Unmarshal([]byte(value), &result); err == nil {
		return result, nil
	}
	// إذا فشل التحليل، إرجاع القيمة كما هي
	return value, nil
}

func (c *CacheServiceImpl) extractUsedMemory(info string) string {
	lines := strings.Split(info, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "used_memory_human:") {
			parts := strings.Split(line, ":")
			if len(parts) > 1 {
				return strings.TrimSpace(parts[1])
			}
		}
	}
	return "unknown"
}

func (c *CacheServiceImpl) extractConnectedClients(info string) int64 {
	lines := strings.Split(info, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "connected_clients:") {
			parts := strings.Split(line, ":")
			if len(parts) > 1 {
				if val, err := strconv.ParseInt(strings.TrimSpace(parts[1]), 10, 64); err == nil {
					return val
				}
			}
		}
	}
	return 0
}

func (c *CacheServiceImpl) extractHits(info string) int64 {
	lines := strings.Split(info, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "keyspace_hits:") {
			parts := strings.Split(line, ":")
			if len(parts) > 1 {
				if val, err := strconv.ParseInt(strings.TrimSpace(parts[1]), 10, 64); err == nil {
					return val
				}
			}
		}
	}
	return 0
}

func (c *CacheServiceImpl) extractMisses(info string) int64 {
	lines := strings.Split(info, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "keyspace_misses:") {
			parts := strings.Split(line, ":")
			if len(parts) > 1 {
				if val, err := strconv.ParseInt(strings.TrimSpace(parts[1]), 10, 64); err == nil {
					return val
				}
			}
		}
	}
	return 0
}

func (c *CacheServiceImpl) calculateHitRate(info string) int {
	hits := c.extractHits(info)
	misses := c.extractMisses(info)
	total := hits + misses

	if total == 0 {
		return 0
	}
	return int((hits * 100) / total)
}

func (c *CacheServiceImpl) extractUptime(info string) int64 {
	lines := strings.Split(info, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "uptime_in_seconds:") {
			parts := strings.Split(line, ":")
			if len(parts) > 1 {
				if val, err := strconv.ParseInt(strings.TrimSpace(parts[1]), 10, 64); err == nil {
					return val
				}
			}
		}
	}
	return 0
}

func (c *CacheServiceImpl) getEnvironment() string {
	if c.isRailway {
		return "railway"
	}
	return "local"
}