package config

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/caarlos0/env/v11"
)

// Cors configuration
type Cors struct {
	AllowedOrigins   []string `env:"ALLOWED_ORIGINS" envSeparator:","`
	AllowedMethods   []string `env:"ALLOWED_METHODS" envSeparator:","`
	AllowedHeaders   []string `env:"ALLOWED_HEADERS" envSeparator:","`
	ExposedHeaders   []string `env:"EXPOSED_HEADERS" envSeparator:","`
	AllowCredentials bool     `env:"ALLOW_CREDENTIALS"`
	MaxAge           int      `env:"MAX_AGE"`
}

// D1Config Cloudflare D1
type D1Config struct {
	AccountID    string `env:"D1_ACCOUNT_ID"`
	DatabaseName string `env:"D1_DATABASE_NAME"`
	DatabaseID   string `env:"D1_DATABASE_ID"`
	BindingName  string `env:"D1_BINDING_NAME"` // for Workers binding (if used)
	UseD1        bool   `env:"USE_D1"`
}

// Cache configuration
type Cache struct {
	Enabled    bool          `env:"CACHE_ENABLED"`
	Prefix     string        `env:"CACHE_PREFIX"`
	DefaultTTL time.Duration `env:"CACHE_DEFAULT_TTL"`
	MaxRetries int           `env:"CACHE_MAX_RETRIES"`
}

// ServicesConfig ...
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

// Cloudinary ...
type Cloudinary struct {
	CloudName    string `env:"CLOUDINARY_CLOUD_NAME"`
	APIKey       string `env:"CLOUDINARY_API_KEY"`
	APISecret    string `env:"CLOUDINARY_API_SECRET"`
	UploadPreset string `env:"CLOUDINARY_UPLOAD_PRESET"`
	Folder       string `env:"CLOUDINARY_FOLDER"`
}

// Upload ...
type Upload struct {
	MaxFileSize    int64    `env:"UPLOAD_MAX_FILE_SIZE"`
	AllowedTypes   []string `env:"UPLOAD_ALLOWED_TYPES" envSeparator:","`
	ImageMaxWidth  int      `env:"UPLOAD_IMAGE_MAX_WIDTH"`
	ImageMaxHeight int      `env:"UPLOAD_IMAGE_MAX_HEIGHT"`
	StorageBackend string   `env:"UPLOAD_STORAGE_BACKEND"`
}

// Email ...
type Email struct {
	Enabled   bool   `env:"EMAIL_ENABLED"`
	Host      string `env:"EMAIL_HOST"`
	Port      int    `env:"EMAIL_PORT"`
	Username  string `env:"EMAIL_USERNAME"`
	Password  string `env:"EMAIL_PASSWORD"`
	FromEmail string `env:"EMAIL_FROM_EMAIL"`
	FromName  string `env:"EMAIL_FROM_NAME"`
}

// AuthConfig ...
type AuthConfig struct {
	JWTSecret         string        `env:"JWT_SECRET"`
	JWTExpiration     time.Duration `env:"JWT_EXPIRATION"`
	RefreshExpiration time.Duration `env:"REFRESH_EXPIRATION"`
	BCryptCost        int           `env:"BCRYPT_COST"`
}

// Main Config
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

var appConfig *Config

func initDefaultLogger() {
	if slog.Default().Handler() == nil {
		handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
		slog.SetDefault(slog.New(handler))
	}
}

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
			UseD1:        getEnvBool("USE_D1", false),
		},
		Auth: AuthConfig{
			JWTSecret:         getEnv("JWT_SECRET", "default-jwt-secret-change-in-production"),
			JWTExpiration:     getEnvDuration("JWT_EXPIRATION", 24*time.Hour),
			RefreshExpiration: getEnvDuration("REFRESH_EXPIRATION", 7*24*time.Hour),
			BCryptCost:        getEnvInt("BCRYPT_COST", 12),
		},
		Cors: Cors{
			AllowedOrigins:   getEnvSlice("ALLOWED_ORIGINS", []string{}, ","),
			AllowedMethods:   getEnvSlice("ALLOWED_METHODS", []string{}, ","),
			AllowedHeaders:   getEnvSlice("ALLOWED_HEADERS", []string{}, ","),
			ExposedHeaders:   getEnvSlice("EXPOSED_HEADERS", []string{}, ","),
			AllowCredentials: getEnvBool("ALLOW_CREDENTIALS", true),
			MaxAge:           getEnvInt("MAX_AGE", 86400),
		},
		Cache: Cache{
			Enabled:    getEnvBool("CACHE_ENABLED", true),
			Prefix:     getEnv("CACHE_PREFIX", "nawthtech:"),
			DefaultTTL: getEnvDuration("CACHE_DEFAULT_TTL", 1*time.Hour),
			MaxRetries: getEnvInt("CACHE_MAX_RETRIES", 3),
		},
		Services: ServicesConfig{
			MaxServicesPerUser:     getEnvInt("SERVICES_MAX_PER_USER", 50),
			MaxActiveServices:      getEnvInt("SERVICES_MAX_ACTIVE", 20),
			DefaultPaginationLimit: getEnvInt("SERVICES_PAGINATION_LIMIT", 20),
			MaxPaginationLimit:     getEnvInt("SERVICES_MAX_PAGINATION_LIMIT", 100),
			SearchCacheTTL:         getEnvDuration("SERVICES_SEARCH_CACHE_TTL", 5*time.Minute),
			FeaturedCacheTTL:       getEnvDuration("SERVICES_FEATURED_CACHE_TTL", 30*time.Minute),
			MaxImagesPerService:    getEnvInt("SERVICES_MAX_IMAGES", 10),
			MaxTagsPerService:      getEnvInt("SERVICES_MAX_TAGS", 15),
			MinTitleLength:         getEnvInt("SERVICES_MIN_TITLE_LENGTH", 3),
			MaxTitleLength:         getEnvInt("SERVICES_MAX_TITLE_LENGTH", 200),
			MinDescriptionLength:   getEnvInt("SERVICES_MIN_DESCRIPTION_LENGTH", 10),
			MaxDescriptionLength:   getEnvInt("SERVICES_MAX_DESCRIPTION_LENGTH", 2000),
			MinPrice:               getEnvFloat("SERVICES_MIN_PRICE", 0),
			MaxPrice:               getEnvFloat("SERVICES_MAX_PRICE", 1000000),
			MinDuration:            getEnvInt("SERVICES_MIN_DURATION", 1),
			MaxDuration:            getEnvInt("SERVICES_MAX_DURATION", 365),
			AutoApproveServices:    getEnvBool("SERVICES_AUTO_APPROVE", true),
			AllowServiceEditing:    getEnvBool("SERVICES_ALLOW_EDITING", true),
			EnableServiceReviews:   getEnvBool("SERVICES_ENABLE_REVIEWS", true),
			EnableServiceRatings:   getEnvBool("SERVICES_ENABLE_RATINGS", true),
			RateLimitCreate:        getEnvInt("SERVICES_RATE_LIMIT_CREATE", 10),
			RateLimitUpdate:        getEnvInt("SERVICES_RATE_LIMIT_UPDATE", 30),
			RateLimitSearch:        getEnvInt("SERVICES_RATE_LIMIT_SEARCH", 60),
		},
		Upload: Upload{
			MaxFileSize:    getEnvInt64("UPLOAD_MAX_FILE_SIZE", 10*1024*1024),
			AllowedTypes:   getEnvSlice("UPLOAD_ALLOWED_TYPES", []string{"image/jpeg", "image/png", "image/gif", "image/webp", "application/pdf"}, ","),
			ImageMaxWidth:  getEnvInt("UPLOAD_IMAGE_MAX_WIDTH", 1920),
			ImageMaxHeight: getEnvInt("UPLOAD_IMAGE_MAX_HEIGHT", 1080),
			StorageBackend: getEnv("UPLOAD_STORAGE_BACKEND", "cloudinary"),
		},
		Cloudinary: Cloudinary{
			CloudName:    getEnv("CLOUDINARY_CLOUD_NAME", ""),
			APIKey:       getEnv("CLOUDINARY_API_KEY", ""),
			APISecret:    getEnv("CLOUDINARY_API_SECRET", ""),
			UploadPreset: getEnv("CLOUDINARY_UPLOAD_PRESET", "nawthtech_uploads"),
			Folder:       getEnv("CLOUDINARY_FOLDER", "nawthtech"),
		},
		Email: Email{
			Enabled:   getEnvBool("EMAIL_ENABLED", false),
			Host:      getEnv("EMAIL_HOST", ""),
			Port:      getEnvInt("EMAIL_PORT", 587),
			Username:  getEnv("EMAIL_USERNAME", ""),
			Password:  getEnv("EMAIL_PASSWORD", ""),
			FromEmail: getEnv("EMAIL_FROM_EMAIL", "noreply@nawthtech.com"),
			FromName:  getEnv("EMAIL_FROM_NAME", "نوذ تك"),
		},
	}

	setCorsDefaults()

	if err := validateConfig(); err != nil {
		slog.Error("فشل التحقق من صحة الإعدادات", "error", err)
		os.Exit(1)
	}

	if err := env.Parse(appConfig); err != nil {
		slog.Error("فشل تحليل متغيرات البيئة", "error", err)
		os.Exit(1)
	}

	slog.Info("تم تحميل الإعدادات بنجاح",
		"environment", appConfig.Environment,
		"port", appConfig.Port,
		"version", appConfig.Version,
		"database", func() string {
			if appConfig.D1.UseD1 {
				return "D1 Cloudflare"
			}
			return "sqlite (local dev)"
		}(),
		"storage", appConfig.Upload.StorageBackend,
	)

	return appConfig
}

// helpers (getEnv, getEnvInt, etc.) — استخدم نفس دوالك القديمة (أنسخها من ملفك السابق).
// لتوفير المساحة هنا، افترض وجود دوال:
// getEnv, getEnvInt, getEnvBool, getEnvDuration, getEnvSlice, getEnvFloat, getEnvInt64
// و setCorsDefaults, validateConfig, validateCloudinaryConfig, validateAuthConfig, etc.
//
// كما أضفت دالة مفيدة:
func (c *Config) GetD1Config() map[string]interface{} {
	return map[string]interface{}{
		"account_id":    c.D1.AccountID,
		"database_name": c.D1.DatabaseName,
		"database_id":   c.D1.DatabaseID,
		"binding_name":  c.D1.BindingName,
		"use_d1":        c.D1.UseD1,
	}
}