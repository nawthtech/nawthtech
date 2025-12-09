package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/caarlos0/env/v11"
	"log/slog"
)

// ========== هياكل التكوين ==========

// Cors تكوين CORS
type Cors struct {
	AllowedOrigins   []string `env:"ALLOWED_ORIGINS" envSeparator:","`
	AllowedMethods   []string `env:"ALLOWED_METHODS" envSeparator:","`
	AllowedHeaders   []string `env:"ALLOWED_HEADERS" envSeparator:","`
	ExposedHeaders   []string `env:"EXPOSED_HEADERS" envSeparator:","`
	AllowCredentials bool     `env:"ALLOW_CREDENTIALS"`
	MaxAge           int      `env:"MAX_AGE"`
}

// D1Config تكوين D1 Cloudflare
type D1Config struct {
	AccountID    string `env:"D1_ACCOUNT_ID"`
	DatabaseName string `env:"D1_DATABASE_NAME"`
	DatabaseID   string `env:"D1_DATABASE_ID"`
	BindingName  string `env:"D1_BINDING_NAME"` // الاسم المستخدم للوصول للDB
}

// Cache تكوين التخزين المؤقت
type Cache struct {
	Enabled    bool          `env:"CACHE_ENABLED"`
	Prefix     string        `env:"CACHE_PREFIX"`
	DefaultTTL time.Duration `env:"CACHE_DEFAULT_TTL"`
	MaxRetries int           `env:"CACHE_MAX_RETRIES"`
}

// ServicesConfig تكوين الخدمات
type ServicesConfig struct {
	MaxServicesPerUser     int           `env:"SERVICES_MAX_PER_USER"`
	MaxActiveServices      int           `env:"SERVICES_MAX_ACTIVE"`
	DefaultPaginationLimit int           `env:"SERVICES_PAGINATION_LIMIT"`
	MaxPaginationLimit     int           `env:"SERVICES_MAX_PAGINATION_LIMIT"`
	SearchCacheTTL         time.Duration `env:"SERVICES_SEARCH_CACHE_TTL"`
	FeaturedCacheTTL       time.Duration `env:"SERVICES_FEATURED_CACHE_TTL"`
	MaxImagesPerService    int           `env:"SERVICES_MAX_IMAGES"`
	MaxTagsPerService      int           `env:"SERVICES_MAX_TAGS"`
	MinTitleLength         int           `env:"SERVICES_MIN_TITLE_LENGTH"`
	MaxTitleLength         int           `env:"SERVICES_MAX_TITLE_LENGTH"`
	MinDescriptionLength   int           `env:"SERVICES_MIN_DESCRIPTION_LENGTH"`
	MaxDescriptionLength   int           `env:"SERVICES_MAX_DESCRIPTION_LENGTH"`
	MinPrice               float64       `env:"SERVICES_MIN_PRICE"`
	MaxPrice               float64       `env:"SERVICES_MAX_PRICE"`
	MinDuration            int           `env:"SERVICES_MIN_DURATION"`
	MaxDuration            int           `env:"SERVICES_MAX_DURATION"`
	AutoApproveServices    bool          `env:"SERVICES_AUTO_APPROVE"`
	AllowServiceEditing    bool          `env:"SERVICES_ALLOW_EDITING"`
	EnableServiceReviews   bool          `env:"SERVICES_ENABLE_REVIEWS"`
	EnableServiceRatings   bool          `env:"SERVICES_ENABLE_RATINGS"`
	RateLimitCreate        int           `env:"SERVICES_RATE_LIMIT_CREATE"`
	RateLimitUpdate        int           `env:"SERVICES_RATE_LIMIT_UPDATE"`
	RateLimitSearch        int           `env:"SERVICES_RATE_LIMIT_SEARCH"`
}

// Cloudinary تكوين Cloudinary
type Cloudinary struct {
	CloudName    string `env:"CLOUDINARY_CLOUD_NAME"`
	APIKey       string `env:"CLOUDINARY_API_KEY"`
	APISecret    string `env:"CLOUDINARY_API_SECRET"`
	UploadPreset string `env:"CLOUDINARY_UPLOAD_PRESET"`
	Folder       string `env:"CLOUDINARY_FOLDER"`
}

// Upload تكوين الرفع
type Upload struct {
	MaxFileSize    int64    `env:"UPLOAD_MAX_FILE_SIZE"`
	AllowedTypes   []string `env:"UPLOAD_ALLOWED_TYPES" envSeparator:","`
	ImageMaxWidth  int      `env:"UPLOAD_IMAGE_MAX_WIDTH"`
	ImageMaxHeight int      `env:"UPLOAD_IMAGE_MAX_HEIGHT"`
	StorageBackend string   `env:"UPLOAD_STORAGE_BACKEND"` // cloudinary أو local
}

// Email تكوين البريد
type Email struct {
	Enabled   bool   `env:"EMAIL_ENABLED"`
	Host      string `env:"EMAIL_HOST"`
	Port      int    `env:"EMAIL_PORT"`
	Username  string `env:"EMAIL_USERNAME"`
	Password  string `env:"EMAIL_PASSWORD"`
	FromEmail string `env:"EMAIL_FROM_EMAIL"`
	FromName  string `env:"EMAIL_FROM_NAME"`
}

// AuthConfig تكوين المصادقة
type AuthConfig struct {
	JWTSecret         string        `env:"JWT_SECRET"`
	JWTExpiration     time.Duration `env:"JWT_EXPIRATION"`
	RefreshExpiration time.Duration `env:"REFRESH_EXPIRATION"`
	BCryptCost        int           `env:"BCRYPT_COST"`
}

// Config التكوين الرئيسي
type Config struct {
	Environment   string         `env:"ENVIRONMENT"`
	Port          string         `env:"PORT"`
	Version       string         `env:"APP_VERSION"`
	EncryptionKey string         `env:"ENCRYPTION_KEY"`
	D1            D1Config       `envPrefix:"D1_"`
	Auth          AuthConfig     `envPrefix:"AUTH_"`
	Cors          Cors           `envPrefix:"CORS_"`
	Cache         Cache          `envPrefix:"CACHE_"`
	Services      ServicesConfig `envPrefix:"SERVICES_"`
	Upload        Upload         `envPrefix:"UPLOAD_"`
	Cloudinary    Cloudinary     `envPrefix:"CLOUDINARY_"`
	Email         Email          `envPrefix:"EMAIL_"`
}

// ========== متغيرات عامة ==========

var appConfig *Config

// ========== التهيئة ==========

// Load تحميل الإعدادات من البيئة
func Load() *Config {
	if appConfig != nil {
		return appConfig
	}

	initDefaultLogger()

	appConfig = &Config{
		Environment:   getEnv("ENVIRONMENT", "development"),
		Port:          getEnv("PORT", "3000"),
		Version:       getEnv("APP_VERSION", "1.0.0"),
		EncryptionKey: getEnv("ENCRYPTION_KEY", "default-encryption-key-change-in-production"),
		D1: D1Config{
			AccountID:    getEnv("D1_ACCOUNT_ID", ""),
			DatabaseName: getEnv("D1_DATABASE_NAME", "nawthtech_d1"),
			DatabaseID:   getEnv("D1_DATABASE_ID", ""),
			BindingName:  getEnv("D1_BINDING_NAME", "DB"),
		},
		Auth: AuthConfig{
			JWTSecret:         getEnv("JWT_SECRET", "default-jwt-secret-change-in-production"),
			JWTExpiration:     getEnvDuration("JWT_EXPIRATION", 24*time.Hour),
			RefreshExpiration: getEnvDuration("REFRESH_EXPIRATION", 7*24*time.Hour),
			BCryptCost:        getEnvInt("BCRYPT_COST", 12),
		},
	}

	if err := env.Parse(appConfig); err != nil {
		slog.Error("فشل تحليل متغيرات البيئة", "error", err)
		os.Exit(1)
	}

	slog.Info("تم تحميل الإعدادات بنجاح",
		"environment", appConfig.Environment,
		"port", appConfig.Port,
		"version", appConfig.Version,
		"database", "D1 Cloudflare",
	)

	return appConfig
}

// ========== دوال مساعدة لتحويل متغيرات البيئة ==========

func getEnv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}

func getEnvInt(key string, fallback int) int {
	valStr := os.Getenv(key)
	if valStr == "" {
		return fallback
	}
	val, err := strconv.Atoi(valStr)
	if err != nil {
		return fallback
	}
	return val
}

func getEnvInt64(key string, fallback int64) int64 {
	valStr := os.Getenv(key)
	if valStr == "" {
		return fallback
	}
	val, err := strconv.ParseInt(valStr, 10, 64)
	if err != nil {
		return fallback
	}
	return val
}

func getEnvBool(key string, fallback bool) bool {
	valStr := os.Getenv(key)
	if valStr == "" {
		return fallback
	}
	valStr = strings.ToLower(valStr)
	return valStr == "true" || valStr == "1"
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	valStr := os.Getenv(key)
	if valStr == "" {
		return fallback
	}
	dur, err := time.ParseDuration(valStr)
	if err != nil {
		return fallback
	}
	return dur
}

// ========== الوصول إلى إعدادات D1 ==========

// GetD1Config الحصول على إعدادات D1
func (c *Config) GetD1Config() D1Config {
	return c.D1
}

// ========== تهيئة السجل الافتراضي ==========
func initDefaultLogger() {
	if slog.Default() == nil {
		slog.SetDefault(slog.New())
	}
}