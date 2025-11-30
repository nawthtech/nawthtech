package config

import (
	"os"
	"testing"
)

func TestConfigDefaults(t *testing.T) {
	// حفظ الإعدادات الحالية
	originalEnv := os.Getenv("ENVIRONMENT")
	originalDB := os.Getenv("DATABASE_URL")
	originalJWT := os.Getenv("JWT_SECRET")
	originalEncryption := os.Getenv("ENCRYPTION_KEY")
	
	defer func() {
		os.Setenv("ENVIRONMENT", originalEnv)
		os.Setenv("DATABASE_URL", originalDB)
		os.Setenv("JWT_SECRET", originalJWT)
		os.Setenv("ENCRYPTION_KEY", originalEncryption)
	}()
	
	// تنظيف appConfig لفرض إعادة التحميل
	appConfig = nil
	
	// تعيين بيئة اختبار بدون استخدام logger
	os.Setenv("ENVIRONMENT", "test")
	os.Setenv("DATABASE_URL", "postgres://test:test@localhost:5432/test")
	os.Setenv("JWT_SECRET", "test-jwt-secret")
	os.Setenv("ENCRYPTION_KEY", "test-encryption-key")
	
	cfg := Load()
	
	if cfg == nil {
		t.Fatal("Expected config to be loaded, got nil")
	}
	
	if cfg.Environment != "test" {
		t.Errorf("Expected environment 'test', got '%s'", cfg.Environment)
	}
	
	// Port should have a default value
	if cfg.Port == "" {
		t.Error("Expected port to have a default value")
	}
}

func TestConfigMethods(t *testing.T) {
	cfg := &Config{
		Environment:   "development",
		Port:          "8080",
		Version:       "1.0.0",
		DatabaseURL:   "postgres://user:pass@localhost:5432/db",
		JWTSecret:     "test-secret",
		EncryptionKey: "test-key",
		RedisURL:      "redis://localhost:6379",
		CacheEnabled:  true,
	}
	
	if !cfg.IsDevelopment() {
		t.Error("Expected IsDevelopment to return true")
	}
	
	if cfg.IsProduction() {
		t.Error("Expected IsProduction to return false for development")
	}
	
	dsn := cfg.GetDSN()
	if dsn != "postgres://user:pass@localhost:5432/db" {
		t.Errorf("Expected DSN to match, got '%s'", dsn)
	}
	
	port := cfg.GetPort()
	if port != "8080" {
		t.Errorf("Expected port '8080', got '%s'", port)
	}
	
	version := cfg.GetVersion()
	if version != "1.0.0" {
		t.Errorf("Expected version '1.0.0', got '%s'", version)
	}
	
	env := cfg.GetEnvironment()
	if env != "development" {
		t.Errorf("Expected environment 'development', got '%s'", env)
	}
	
	jwtSecret := cfg.GetJWTSecret()
	if jwtSecret != "test-secret" {
		t.Errorf("Expected JWT secret 'test-secret', got '%s'", jwtSecret)
	}
	
	encryptionKey := cfg.GetEncryptionKey()
	if encryptionKey != "test-key" {
		t.Errorf("Expected encryption key 'test-key', got '%s'", encryptionKey)
	}
	
	redisAddr := cfg.GetRedisAddress()
	if redisAddr != "redis://localhost:6379" {
		t.Errorf("Expected Redis address 'redis://localhost:6379', got '%s'", redisAddr)
	}
	
	if !cfg.IsCacheEnabled() {
		t.Error("Expected IsCacheEnabled to return true")
	}
}

func TestConfigEnvironmentDetection(t *testing.T) {
	testCases := []struct {
		env      string
		isDev    bool
		isProd   bool
		isStaging bool
	}{
		{"development", true, false, false},
		{"production", false, true, false},
		{"staging", false, false, true},
		{"test", false, false, false},
	}
	
	for _, tc := range testCases {
		cfg := &Config{Environment: tc.env}
		
		if cfg.IsDevelopment() != tc.isDev {
			t.Errorf("IsDevelopment() for env '%s': expected %v, got %v", tc.env, tc.isDev, cfg.IsDevelopment())
		}
		
		if cfg.IsProduction() != tc.isProd {
			t.Errorf("IsProduction() for env '%s': expected %v, got %v", tc.env, tc.isProd, cfg.IsProduction())
		}
		
		if cfg.IsStaging() != tc.isStaging {
			t.Errorf("IsStaging() for env '%s': expected %v, got %v", tc.env, tc.isStaging, cfg.IsStaging())
		}
	}
}

func TestConfigEmptyValues(t *testing.T) {
	cfg := &Config{}
	
	// Test empty values
	if cfg.GetDSN() != "" {
		t.Errorf("Expected empty DSN, got '%s'", cfg.GetDSN())
	}
	
	if cfg.GetPort() != "" {
		t.Errorf("Expected empty port, got '%s'", cfg.GetPort())
	}
	
	if cfg.GetVersion() != "" {
		t.Errorf("Expected empty version, got '%s'", cfg.GetVersion())
	}
	
	if cfg.GetEnvironment() != "" {
		t.Errorf("Expected empty environment, got '%s'", cfg.GetEnvironment())
	}
	
	if cfg.GetJWTSecret() != "" {
		t.Errorf("Expected empty JWT secret, got '%s'", cfg.GetJWTSecret())
	}
	
	if cfg.GetEncryptionKey() != "" {
		t.Errorf("Expected empty encryption key, got '%s'", cfg.GetEncryptionKey())
	}
	
	if cfg.GetRedisAddress() != "" {
		t.Errorf("Expected empty Redis address, got '%s'", cfg.GetRedisAddress())
	}
	
	if cfg.IsCacheEnabled() {
		t.Error("Expected IsCacheEnabled to return false for empty config")
	}
}