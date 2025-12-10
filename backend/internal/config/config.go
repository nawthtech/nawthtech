package config

import (
	"os"
)

type Config struct {
	Environment string `mapstructure:"environment"`
	Port        string `mapstructure:"port"`
	WorkerURL   string `mapstructure:"worker_url"`
	WorkerKey   string `mapstructure:"worker_key"`

	// قسم المصادقة
	Auth struct {
		JWTSecret string `mapstructure:"jwt_secret"`
		// يمكنك إضافة المزيد من الحقول هنا
	} `mapstructure:"auth"`

	// أقسام إضافية (اختيارية)
	CORS struct {
		AllowedOrigins []string `mapstructure:"allowed_origins"`
	} `mapstructure:"cors"`

	Database struct {
		URL string `mapstructure:"url"`
	} `mapstructure:"database"`

	Email struct {
		From     string `mapstructure:"from"`
		Provider string `mapstructure:"provider"`
	} `mapstructure:"email"`

	Upload struct {
		MaxSize int64  `mapstructure:"max_size"`
		Path    string `mapstructure:"path"`
	} `mapstructure:"upload"`

	Cache struct {
		Enabled bool   `mapstructure:"enabled"`
		Redis   string `mapstructure:"redis"`
	} `mapstructure:"cache"`

	Security struct {
		RateLimit int `mapstructure:"rate_limit"`
	} `mapstructure:"security"`

	Logging struct {
		Level string `mapstructure:"level"`
	} `mapstructure:"logging"`
}

// Load يحمل الإعدادات من متغيرات البيئة
func Load() *Config {
	config := &Config{}

	// Environment
	config.Environment = os.Getenv("ENVIRONMENT")
	if config.Environment == "" {
		config.Environment = "development"
	}

	// Port
	config.Port = os.Getenv("PORT")
	if config.Port == "" {
		config.Port = "8080"
	}

	// Worker URL
	config.WorkerURL = os.Getenv("WORKER_API_URL")
	if config.WorkerURL == "" {
		config.WorkerURL = "https://api.nawthtech.com"
	}

	// Worker Key
	config.WorkerKey = os.Getenv("WORKER_API_KEY")

	// Auth - JWT Secret
	config.Auth.JWTSecret = os.Getenv("JWT_SECRET")
	if config.Auth.JWTSecret == "" {
		// قيمة افتراضية للتطوير فقط
		config.Auth.JWTSecret = "dev-secret-key-change-in-production"
	}

	// CORS
	corsOrigins := os.Getenv("CORS_ALLOWED_ORIGINS")
	if corsOrigins != "" {
		// يمكنك معالجة القائمة هنا إذا احتجت
		config.CORS.AllowedOrigins = []string{corsOrigins}
	} else {
		config.CORS.AllowedOrigins = []string{"*"}
	}

	// Database
	config.Database.URL = os.Getenv("DATABASE_URL")

	// Email
	config.Email.From = os.Getenv("EMAIL_FROM")
	config.Email.Provider = os.Getenv("EMAIL_PROVIDER")

	// Upload
	config.Upload.Path = os.Getenv("UPLOAD_PATH")
	if config.Upload.Path == "" {
		config.Upload.Path = "./uploads"
	}

	// Cache
	cacheEnabled := os.Getenv("CACHE_ENABLED")
	config.Cache.Enabled = cacheEnabled == "true"
	config.Cache.Redis = os.Getenv("REDIS_URL")

	// Logging
	config.Logging.Level = os.Getenv("LOG_LEVEL")
	if config.Logging.Level == "" {
		config.Logging.Level = "info"
	}

	return config
}

func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}

func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}

func (c *Config) IsStaging() bool {
	return c.Environment == "staging"
}

// طريقة للحصول على JWTSecret
func (c *Config) GetJWTSecret() string {
	return c.Auth.JWTSecret
}

// LoadConfig اسم بديل لـ Load ليكون متوافقاً مع الكود الحالي
func LoadConfig() *Config {
	return Load()
}
