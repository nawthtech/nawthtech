package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
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
	
	// تعيين بيئة اختبار
	os.Setenv("ENVIRONMENT", "test")
	os.Setenv("DATABASE_URL", "postgres://test:test@localhost:5432/test")
	os.Setenv("JWT_SECRET", "test-jwt-secret")
	os.Setenv("ENCRYPTION_KEY", "test-encryption-key")
	
	// إعادة تعيين appConfig لفرض إعادة التحميل
	appConfig = nil
	
	cfg := Load()
	
	if cfg == nil {
		t.Fatal("Expected config to be loaded, got nil")
	}
	
	if cfg.Environment != "test" {
		t.Errorf("Expected environment 'test', got '%s'", cfg.Environment)
	}
	
	if cfg.Port == "" {
		t.Error("Expected port to be set")
	}
	
	if cfg.Version == "" {
		t.Error("Expected version to be set")
	}
}

func TestIsDevelopment(t *testing.T) {
	cfg := &Config{Environment: "development"}
	
	if !cfg.IsDevelopment() {
		t.Error("Expected IsDevelopment to return true for 'development' environment")
	}
	
	cfg.Environment = "production"
	if cfg.IsDevelopment() {
		t.Error("Expected IsDevelopment to return false for 'production' environment")
	}
}

func TestIsProduction(t *testing.T) {
	cfg := &Config{Environment: "production"}
	
	if !cfg.IsProduction() {
		t.Error("Expected IsProduction to return true for 'production' environment")
	}
	
	cfg.Environment = "development"
	if cfg.IsProduction() {
		t.Error("Expected IsProduction to return false for 'development' environment")
	}
}

func TestGetDSN(t *testing.T) {
	cfg := &Config{
		Database: DatabaseConfig{
			URL: "postgres://user:pass@localhost:5432/db",
		},
	}
	
	dsn := cfg.GetDSN()
	if dsn != "postgres://user:pass@localhost:5432/db" {
		t.Errorf("Expected DSN to match, got '%s'", dsn)
	}
}

func TestGetPort(t *testing.T) {
	cfg := &Config{Port: "8080"}
	
	port := cfg.GetPort()
	if port != "8080" {
		t.Errorf("Expected port '8080', got '%s'", port)
	}
}

func TestGetVersion(t *testing.T) {
	cfg := &Config{Version: "1.0.0"}
	
	version := cfg.GetVersion()
	if version != "1.0.0" {
		t.Errorf("Expected version '1.0.0', got '%s'", version)
	}
}

func TestGetEnvironment(t *testing.T) {
	cfg := &Config{Environment: "staging"}
	
	env := cfg.GetEnvironment()
	if env != "staging" {
		t.Errorf("Expected environment 'staging', got '%s'", env)
	}
}

func TestIsStaging(t *testing.T) {
	cfg := &Config{Environment: "staging"}
	
	if !cfg.IsStaging() {
		t.Error("Expected IsStaging to return true for 'staging' environment")
	}
	
	cfg.Environment = "development"
	if cfg.IsStaging() {
		t.Error("Expected IsStaging to return false for 'development' environment")
	}
}

func TestGetJWTSecret(t *testing.T) {
	cfg := &Config{
		Auth: AuthConfig{
			JWTSecret: "test-secret",
		},
	}
	
	secret := cfg.GetJWTSecret()
	if secret != "test-secret" {
		t.Errorf("Expected JWT secret 'test-secret', got '%s'", secret)
	}
}

func TestGetEncryptionKey(t *testing.T) {
	cfg := &Config{EncryptionKey: "test-key"}
	
	key := cfg.GetEncryptionKey()
	if key != "test-key" {
		t.Errorf("Expected encryption key 'test-key', got '%s'", key)
	}
}

func TestGetRedisAddress(t *testing.T) {
	// Test with URL
	cfg1 := &Config{
		Redis: Redis{
			URL: "redis://localhost:6379",
		},
	}
	
	addr1 := cfg1.GetRedisAddress()
	if addr1 != "redis://localhost:6379" {
		t.Errorf("Expected Redis URL 'redis://localhost:6379', got '%s'", addr1)
	}
	
	// Test with host and port
	cfg2 := &Config{
		Redis: Redis{
			Host: "localhost",
			Port: "6380",
		},
	}
	
	addr2 := cfg2.GetRedisAddress()
	if addr2 != "localhost:6380" {
		t.Errorf("Expected Redis address 'localhost:6380', got '%s'", addr2)
	}
}

func TestIsCacheEnabled(t *testing.T) {
	// Test when cache is enabled
	cfg1 := &Config{
		Cache: Cache{
			Enabled: true,
		},
	}
	
	if !cfg1.IsCacheEnabled() {
		t.Error("Expected IsCacheEnabled to return true when cache is enabled")
	}
	
	// Test when cache is disabled
	cfg2 := &Config{
		Cache: Cache{
			Enabled: false,
		},
	}
	
	if cfg2.IsCacheEnabled() {
		t.Error("Expected IsCacheEnabled to return false when cache is disabled")
	}
}

func TestConfigValidation(t *testing.T) {
	// Test that Load doesn't panic with minimal required config
	originalEnv := os.Getenv("ENVIRONMENT")
	originalDB := os.Getenv("DATABASE_URL")
	originalJWT := os.Getenv("JWT_SECRET")
	originalEncryption := os.Getenv("ENCRYPTION_KEY")
	
	defer func() {
		os.Setenv("ENVIRONMENT", originalEnv)
		os.Setenv("DATABASE_URL", originalDB)
		os.Setenv("JWT_SECRET", originalJWT)
		os.Setenv("ENCRYPTION_KEY", originalEncryption)
		
		// Reset appConfig
		appConfig = nil
	}()
	
	// Set minimal required environment variables
	os.Setenv("ENVIRONMENT", "test")
	os.Setenv("DATABASE_URL", "postgres://user:pass@localhost:5432/db")
	os.Setenv("JWT_SECRET", "test-jwt-secret")
	os.Setenv("ENCRYPTION_KEY", "test-encryption-key")
	
	// This should not panic
	cfg := Load()
	
	if cfg == nil {
		t.Fatal("Expected config to be loaded without panicking")
	}
}