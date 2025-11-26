package config

import (
	"os"

	"backend-app/internal/logger"

	"github.com/caarlos0/env/v11"
)

type cors struct {
	AllowedOrigins []string `env:"ALLOWED_ORIGINS,required,notEmpty" envSeparator:","`
}

type Config struct {
	Environment    string
	DatabaseURL    string
	EncryptionKey  string
	JWTSecret      string
	Port           string
	Version        string
	Cors           *cors
}

var (
	Cors = &cors{}
	AppConfig = &Config{}
)

func init() {
	// تحليل متغيرات البيئة للـ CORS
	toParse := []any{Cors}
	errors := []error{}

	for _, v := range toParse {
		if err := env.Parse(v); err != nil {
			if er, ok := err.(env.AggregateError); ok {
				errors = append(errors, er.Errors...)
				continue
			}
			errors = append(errors, err)
		}
	}

	if len(errors) > 0 {
		logger.Stderr.Error("errors found while parsing environment variables", logger.ErrorsAttr(errors...))
		os.Exit(1)
	}
}

// Load تحميل الإعدادات
func Load() *Config {
	AppConfig = &Config{
		Environment:   getEnv("ENVIRONMENT", "development"),
		DatabaseURL:   getEnv("DATABASE_URL", ""),
		EncryptionKey: getEnv("ENCRYPTION_KEY", ""),
		JWTSecret:     getEnv("JWT_SECRET", ""),
		Port:          getEnv("PORT", "3000"),
		Version:       getEnv("APP_VERSION", "1.0.0"),
		Cors:          Cors,
	}
	return AppConfig
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
