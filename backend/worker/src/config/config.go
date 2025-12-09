package config

import (
	"os"
	"strings"
)

// Config إعدادات عامة
type Config struct {
	Environment string
	Port        string
	Version     string
	JWTSecret   string
	SessionSecret string
	D1APIKey    string
	D1DBName    string
	CORSOrigins []string
}

// LoadConfig تحميل المتغيرات من البيئة
func LoadConfig() *Config {
	c := &Config{
		Environment:   getEnv("ENVIRONMENT", "development"),
		Port:          getEnv("PORT", "3000"),
		Version:       getEnv("API_VERSION", "v1"),
		JWTSecret:     getEnv("JWT_SECRET", ""),
		SessionSecret: getEnv("SESSION_SECRET", ""),
		D1APIKey:      getEnv("D1_API_KEY", ""),
		D1DBName:      getEnv("D1_DB_NAME", ""),
		CORSOrigins:   strings.Split(getEnv("CORS_ALLOWED_ORIGINS", ""), ","),
	}
	return c
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}