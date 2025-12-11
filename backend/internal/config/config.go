package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config يخزن جميع إعدادات التطبيق
type Config struct {
	// الإعدادات الأساسية
	AppName     string `mapstructure:"app_name"`
	Version     string `mapstructure:"version"`
	Environment string `mapstructure:"environment"`
	Port        string `mapstructure:"port"`
	Debug       bool   `mapstructure:"debug"`
	
	// URLs
	APIURL      string `mapstructure:"api_url"`
	FrontendURL string `mapstructure:"frontend_url"`
	WorkerURL   string `mapstructure:"worker_url"`
	
	// Security
	WorkerKey   string `mapstructure:"worker_key"`
	APIKey      string `mapstructure:"api_key"`
	EncryptionKey string `mapstructure:"encryption_key"`
	
	// قسم المصادقة
	Auth struct {
		JWTSecret         string        `mapstructure:"jwt_secret"`
		JWTExpiration     time.Duration `mapstructure:"jwt_expiration"`
		RefreshSecret     string        `mapstructure:"refresh_secret"`
		RefreshExpiration time.Duration `mapstructure:"refresh_expiration"`
		ResetTokenExpiry  time.Duration `mapstructure:"reset_token_expiry"`
		VerifyTokenExpiry time.Duration `mapstructure:"verify_token_expiry"`
	} `mapstructure:"auth"`
	
	// قاعدة البيانات (Cloudflare D1/SQLite)
	Database struct {
		Driver   string `mapstructure:"driver"`
		URL      string `mapstructure:"url"`
		Host     string `mapstructure:"host"`
		Port     string `mapstructure:"port"`
		Name     string `mapstructure:"name"`
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		SSLMode  string `mapstructure:"ssl_mode"`
		MaxConns int    `mapstructure:"max_conns"`
		MaxIdle  int    `mapstructure:"max_idle"`
	} `mapstructure:"database"`
	
	// CORS
	CORS struct {
		AllowedOrigins   []string `mapstructure:"allowed_origins"`
		AllowedMethods   []string `mapstructure:"allowed_methods"`
		AllowedHeaders   []string `mapstructure:"allowed_headers"`
		AllowCredentials bool     `mapstructure:"allow_credentials"`
		MaxAge           int      `mapstructure:"max_age"`
	} `mapstructure:"cors"`
	
	// البريد الإلكتروني
	Email struct {
		Enabled    bool   `mapstructure:"enabled"`
		Provider   string `mapstructure:"provider"`
		Host       string `mapstructure:"host"`
		Port       int    `mapstructure:"port"`
		Username   string `mapstructure:"username"`
		Password   string `mapstructure:"password"`
		From       string `mapstructure:"from"`
		FromName   string `mapstructure:"from_name"`
		ReplyTo    string `mapstructure:"reply_to"`
		TLS        bool   `mapstructure:"tls"`
	} `mapstructure:"email"`
	
	// الرفع (Upload)
	Upload struct {
		MaxSize        int64    `mapstructure:"max_size"`
		Path           string   `mapstructure:"path"`
		AllowedTypes   []string `mapstructure:"allowed_types"`
		CloudinaryURL  string   `mapstructure:"cloudinary_url"`
		S3Bucket       string   `mapstructure:"s3_bucket"`
		S3Region       string   `mapstructure:"s3_region"`
		S3AccessKey    string   `mapstructure:"s3_access_key"`
		S3SecretKey    string   `mapstructure:"s3_secret_key"`
	} `mapstructure:"upload"`
	
	// التخزين المؤقت (Cache)
	Cache struct {
		Enabled bool   `mapstructure:"enabled"`
		Type    string `mapstructure:"type"` // memory, redis
		Redis   string `mapstructure:"redis"`
		TTL     time.Duration `mapstructure:"ttl"`
	} `mapstructure:"cache"`
	
	// الأمان
	Security struct {
		RateLimit       int           `mapstructure:"rate_limit"`
		RateWindow      time.Duration `mapstructure:"rate_window"`
		CSPEnabled      bool          `mapstructure:"csp_enabled"`
		HSTSMaxAge      int           `mapstructure:"hsts_max_age"`
		PasswordMinLen  int           `mapstructure:"password_min_len"`
		PasswordRequire struct {
			Uppercase bool `mapstructure:"uppercase"`
			Lowercase bool `mapstructure:"lowercase"`
			Numbers   bool `mapstructure:"numbers"`
			Symbols   bool `mapstructure:"symbols"`
		} `mapstructure:"password_require"`
	} `mapstructure:"security"`
	
	// التسجيل (Logging)
	Logging struct {
		Level      string `mapstructure:"level"`
		Format     string `mapstructure:"format"` // json, text
		Output     string `mapstructure:"output"` // stdout, file
		File       string `mapstructure:"file"`
		MaxSize    int    `mapstructure:"max_size"`
		MaxBackups int    `mapstructure:"max_backups"`
		MaxAge     int    `mapstructure:"max_age"`
	} `mapstructure:"logging"`
	
	// TLS/SSL
	TLS struct {
		Enabled  bool   `mapstructure:"enabled"`
		CertFile string `mapstructure:"cert_file"`
		KeyFile  string `mapstructure:"key_file"`
	} `mapstructure:"tls"`
	
	// Cloudflare Workers & D1
	Cloudflare struct {
		AccountID  string `mapstructure:"account_id"`
		Namespace  string `mapstructure:"namespace"`
		APIToken   string `mapstructure:"api_token"`
		D1Database string `mapstructure:"d1_database"`
	} `mapstructure:"cloudflare"`
	
	// الخدمات الخارجية
	Services struct {
		Slack struct {
			Token     string `mapstructure:"token"`
			Channel   string `mapstructure:"channel"`
			AppName   string `mapstructure:"app_name"`
		} `mapstructure:"slack"`
		Stripe struct {
			SecretKey     string `mapstructure:"secret_key"`
			WebhookSecret string `mapstructure:"webhook_secret"`
			PublishableKey string `mapstructure:"publishable_key"`
		} `mapstructure:"stripe"`
		Cloudinary struct {
			CloudName string `mapstructure:"cloud_name"`
			APIKey    string `mapstructure:"api_key"`
			APISecret string `mapstructure:"api_secret"`
		} `mapstructure:"cloudinary"`
	} `mapstructure:"services"`
	
	// AI Services
	AI struct {
		OpenAI struct {
			APIKey string `mapstructure:"api_key"`
			Model  string `mapstructure:"model"`
		} `mapstructure:"openai"`
		Gemini struct {
			APIKey string `mapstructure:"api_key"`
		} `mapstructure:"gemini"`
	} `mapstructure:"ai"`
}

// Load يحمل الإعدادات من متغيرات البيئة
func Load() *Config {
	config := &Config{}
	
	// ==================== الإعدادات الأساسية ====================
	config.AppName = getEnv("APP_NAME", "nawthtech")
	config.Version = getEnv("APP_VERSION", "1.0.0")
	config.Environment = getEnv("ENVIRONMENT", "development")
	config.Port = getEnv("PORT", "8080")
	config.Debug = getEnvBool("DEBUG", config.Environment == "development")
	
	// URLs
	config.APIURL = getEnv("API_URL", "http://localhost:"+config.Port)
	config.FrontendURL = getEnv("FRONTEND_URL", "http://localhost:3000")
	config.WorkerURL = getEnv("WORKER_API_URL", "https://api.nawthtech.com")
	
	// Security Keys
	config.WorkerKey = getEnv("WORKER_API_KEY", "")
	config.APIKey = getEnv("API_KEY", "")
	config.EncryptionKey = getEnv("ENCRYPTION_KEY", "dev-encryption-key-change-in-production")
	
	// ==================== المصادقة ====================
	config.Auth.JWTSecret = getEnv("JWT_SECRET", "dev-jwt-secret-change-in-production")
	config.Auth.JWTExpiration = getEnvDuration("JWT_EXPIRATION", 24*time.Hour)
	config.Auth.RefreshSecret = getEnv("REFRESH_SECRET", "dev-refresh-secret-change-in-production")
	config.Auth.RefreshExpiration = getEnvDuration("REFRESH_EXPIRATION", 7*24*time.Hour)
	config.Auth.ResetTokenExpiry = getEnvDuration("RESET_TOKEN_EXPIRY", 1*time.Hour)
	config.Auth.VerifyTokenExpiry = getEnvDuration("VERIFY_TOKEN_EXPIRY", 24*time.Hour)
	
	// ==================== قاعدة البيانات ====================
	// Cloudflare D1 (SQLite)
	config.Database.Driver = getEnv("DB_DRIVER", "sqlite3")
	config.Database.URL = getEnv("DATABASE_URL", "./data/nawthtech.db")
	config.Database.Host = getEnv("DB_HOST", "")
	config.Database.Port = getEnv("DB_PORT", "")
	config.Database.Name = getEnv("DB_NAME", "nawthtech")
	config.Database.User = getEnv("DB_USER", "")
	config.Database.Password = getEnv("DB_PASSWORD", "")
	config.Database.SSLMode = getEnv("DB_SSL_MODE", "disable")
	config.Database.MaxConns = getEnvInt("DB_MAX_CONNS", 25)
	config.Database.MaxIdle = getEnvInt("DB_MAX_IDLE", 5)
	
	// ==================== CORS ====================
	corsOrigins := getEnv("CORS_ALLOWED_ORIGINS", "*")
	if corsOrigins == "*" {
		config.CORS.AllowedOrigins = []string{"*"}
	} else {
		config.CORS.AllowedOrigins = strings.Split(corsOrigins, ",")
	}
	
	config.CORS.AllowedMethods = strings.Split(getEnv("CORS_ALLOWED_METHODS", "GET,POST,PUT,PATCH,DELETE,OPTIONS"), ",")
	config.CORS.AllowedHeaders = strings.Split(getEnv("CORS_ALLOWED_HEADERS", "Content-Type,Authorization,X-Requested-With,X-API-Key"), ",")
	config.CORS.AllowCredentials = getEnvBool("CORS_ALLOW_CREDENTIALS", true)
	config.CORS.MaxAge = getEnvInt("CORS_MAX_AGE", 300)
	
	// ==================== البريد الإلكتروني ====================
	config.Email.Enabled = getEnvBool("EMAIL_ENABLED", false)
	config.Email.Provider = getEnv("EMAIL_PROVIDER", "smtp")
	config.Email.Host = getEnv("EMAIL_HOST", "smtp.gmail.com")
	config.Email.Port = getEnvInt("EMAIL_PORT", 587)
	config.Email.Username = getEnv("EMAIL_USERNAME", "")
	config.Email.Password = getEnv("EMAIL_PASSWORD", "")
	config.Email.From = getEnv("EMAIL_FROM", "noreply@nawthtech.com")
	config.Email.FromName = getEnv("EMAIL_FROM_NAME", "NawthTech")
	config.Email.ReplyTo = getEnv("EMAIL_REPLY_TO", "support@nawthtech.com")
	config.Email.TLS = getEnvBool("EMAIL_TLS", true)
	
	// ==================== الرفع ====================
	config.Upload.MaxSize = getEnvInt64("UPLOAD_MAX_SIZE", 10*1024*1024) // 10MB
	config.Upload.Path = getEnv("UPLOAD_PATH", "./uploads")
	config.Upload.AllowedTypes = strings.Split(getEnv("UPLOAD_ALLOWED_TYPES", "image/jpeg,image/png,image/gif,image/webp,application/pdf"), ",")
	config.Upload.CloudinaryURL = getEnv("CLOUDINARY_URL", "")
	config.Upload.S3Bucket = getEnv("S3_BUCKET", "")
	config.Upload.S3Region = getEnv("S3_REGION", "us-east-1")
	config.Upload.S3AccessKey = getEnv("S3_ACCESS_KEY", "")
	config.Upload.S3SecretKey = getEnv("S3_SECRET_KEY", "")
	
	// ==================== التخزين المؤقت ====================
	config.Cache.Enabled = getEnvBool("CACHE_ENABLED", true)
	config.Cache.Type = getEnv("CACHE_TYPE", "memory")
	config.Cache.Redis = getEnv("REDIS_URL", "redis://localhost:6379")
	config.Cache.TTL = getEnvDuration("CACHE_TTL", 5*time.Minute)
	
	// ==================== الأمان ====================
	config.Security.RateLimit = getEnvInt("RATE_LIMIT", 100)
	config.Security.RateWindow = getEnvDuration("RATE_WINDOW", 1*time.Minute)
	config.Security.CSPEnabled = getEnvBool("CSP_ENABLED", false)
	config.Security.HSTSMaxAge = getEnvInt("HSTS_MAX_AGE", 31536000) // سنة واحدة
	config.Security.PasswordMinLen = getEnvInt("PASSWORD_MIN_LEN", 8)
	config.Security.PasswordRequire.Uppercase = getEnvBool("PASSWORD_REQUIRE_UPPERCASE", true)
	config.Security.PasswordRequire.Lowercase = getEnvBool("PASSWORD_REQUIRE_LOWERCASE", true)
	config.Security.PasswordRequire.Numbers = getEnvBool("PASSWORD_REQUIRE_NUMBERS", true)
	config.Security.PasswordRequire.Symbols = getEnvBool("PASSWORD_REQUIRE_SYMBOLS", false)
	
	// ==================== التسجيل ====================
	config.Logging.Level = getEnv("LOG_LEVEL", "info")
	config.Logging.Format = getEnv("LOG_FORMAT", "json")
	config.Logging.Output = getEnv("LOG_OUTPUT", "stdout")
	config.Logging.File = getEnv("LOG_FILE", "./logs/app.log")
	config.Logging.MaxSize = getEnvInt("LOG_MAX_SIZE", 100) // ميغابايت
	config.Logging.MaxBackups = getEnvInt("LOG_MAX_BACKUPS", 3)
	config.Logging.MaxAge = getEnvInt("LOG_MAX_AGE", 28) // أيام
	
	// ==================== TLS/SSL ====================
	config.TLS.Enabled = getEnvBool("TLS_ENABLED", false)
	config.TLS.CertFile = getEnv("TLS_CERT_FILE", "")
	config.TLS.KeyFile = getEnv("TLS_KEY_FILE", "")
	
	// ==================== Cloudflare ====================
	config.Cloudflare.AccountID = getEnv("CLOUDFLARE_ACCOUNT_ID", "")
	config.Cloudflare.Namespace = getEnv("CLOUDFLARE_NAMESPACE", "")
	config.Cloudflare.APIToken = getEnv("CLOUDFLARE_API_TOKEN", "")
	config.Cloudflare.D1Database = getEnv("CLOUDFLARE_D1_DATABASE", "nawthtech-db")
	
	// ==================== الخدمات الخارجية ====================
	config.Services.Slack.Token = getEnv("SLACK_TOKEN", "")
	config.Services.Slack.Channel = getEnv("SLACK_CHANNEL", "")
	config.Services.Slack.AppName = getEnv("SLACK_APP_NAME", "nawthtech-backend")
	
	config.Services.Stripe.SecretKey = getEnv("STRIPE_SECRET_KEY", "")
	config.Services.Stripe.WebhookSecret = getEnv("STRIPE_WEBHOOK_SECRET", "")
	config.Services.Stripe.PublishableKey = getEnv("STRIPE_PUBLISHABLE_KEY", "")
	
	config.Services.Cloudinary.CloudName = getEnv("CLOUDINARY_CLOUD_NAME", "")
	config.Services.Cloudinary.APIKey = getEnv("CLOUDINARY_API_KEY", "")
	config.Services.Cloudinary.APISecret = getEnv("CLOUDINARY_API_SECRET", "")
	
	// ==================== AI Services ====================
	config.AI.OpenAI.APIKey = getEnv("OPENAI_API_KEY", "")
	config.AI.OpenAI.Model = getEnv("OPENAI_MODEL", "gpt-4-turbo-preview")
	config.AI.Gemini.APIKey = getEnv("GEMINI_API_KEY", "")
	
	return config
}

// ==================== دوال مساعدة ====================

// getEnv يحصل على متغير بيئة أو قيمة افتراضية
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvInt يحصل على متغير بيئة كـ int
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getEnvInt64 يحصل على متغير بيئة كـ int64
func getEnvInt64(key string, defaultValue int64) int64 {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getEnvBool يحصل على متغير بيئة كـ bool
func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

// getEnvDuration يحصل على متغير بيئة كـ time.Duration
func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

// ==================== دوال الوصول ====================

// IsProduction يتحقق إذا كانت البيئة production
func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}

// IsDevelopment يتحقق إذا كانت البيئة development
func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}

// IsStaging يتحقق إذا كانت البيئة staging
func (c *Config) IsStaging() bool {
	return c.Environment == "staging"
}

// GetJWTSecret يحصل على JWT secret
func (c *Config) GetJWTSecret() string {
	return c.Auth.JWTSecret
}

// GetJWTExpiration يحصل على JWT expiration
func (c *Config) GetJWTExpiration() time.Duration {
	return c.Auth.JWTExpiration
}

// GetRefreshExpiration يحصل على refresh token expiration
func (c *Config) GetRefreshExpiration() time.Duration {
	return c.Auth.RefreshExpiration
}

// GetDatabaseURL يحصل على رابط قاعدة البيانات
func (c *Config) GetDatabaseURL() string {
	return c.Database.URL
}

// GetDatabaseDriver يحصل على نوع قاعدة البيانات
func (c *Config) GetDatabaseDriver() string {
	return c.Database.Driver
}

// GetPort يحصل على رقم المنفذ
func (c *Config) GetPort() string {
	return c.Port
}

// GetEnvironment يحصل على البيئة
func (c *Config) GetEnvironment() string {
	return c.Environment
}

// GetVersion يحصل على إصدار التطبيق
func (c *Config) GetVersion() string {
	return c.Version
}

// GetAppName يحصل على اسم التطبيق
func (c *Config) GetAppName() string {
	return c.AppName
}

// GetFrontendURL يحصل على رابط الواجهة الأمامية
func (c *Config) GetFrontendURL() string {
	return c.FrontendURL
}

// GetAPIURL يحصل على رابط API
func (c *Config) GetAPIURL() string {
	return c.APIURL
}

// GetWorkerURL يحصل على رابط Worker
func (c *Config) GetWorkerURL() string {
	return c.WorkerURL
}

// LoadConfig اسم بديل لـ Load ليكون متوافقاً مع الكود الحالي
func LoadConfig() *Config {
	return Load()
}

// String يعرض الإعدادات كـ string (بدون معلومات حساسة)
func (c *Config) String() string {
	return fmt.Sprintf(
		"App: %s v%s\nEnvironment: %s\nPort: %s\nDebug: %v\nDatabase: %s\nJWT Exp: %v\nCache: %v",
		c.AppName, c.Version, c.Environment, c.Port, c.Debug, c.Database.Driver, 
		c.Auth.JWTExpiration, c.Cache.Enabled,
	)
}

// PrintSafe يطبع الإعدادات بدون معلومات حساسة
func (c *Config) PrintSafe() {
	fmt.Println("=== Application Configuration ===")
	fmt.Printf("App Name: %s\n", c.AppName)
	fmt.Printf("Version: %s\n", c.Version)
	fmt.Printf("Environment: %s\n", c.Environment)
	fmt.Printf("Port: %s\n", c.Port)
	fmt.Printf("Debug Mode: %v\n", c.Debug)
	fmt.Printf("Database Driver: %s\n", c.Database.Driver)
	fmt.Printf("Database Name: %s\n", c.Database.Name)
	fmt.Printf("JWT Expiration: %v\n", c.Auth.JWTExpiration)
	fmt.Printf("Refresh Expiration: %v\n", c.Auth.RefreshExpiration)
	fmt.Printf("Cache Enabled: %v\n", c.Cache.Enabled)
	fmt.Printf("Cache Type: %s\n", c.Cache.Type)
	fmt.Printf("Email Enabled: %v\n", c.Email.Enabled)
	fmt.Printf("Email Provider: %s\n", c.Email.Provider)
	fmt.Printf("Upload Max Size: %d bytes\n", c.Upload.MaxSize)
	fmt.Printf("CORS Allowed Origins: %v\n", c.CORS.AllowedOrigins)
	fmt.Printf("Rate Limit: %d requests per %v\n", c.Security.RateLimit, c.Security.RateWindow)
	fmt.Printf("Log Level: %s\n", c.Logging.Level)
	fmt.Printf("TLS Enabled: %v\n", c.TLS.Enabled)
	fmt.Println("================================")
}