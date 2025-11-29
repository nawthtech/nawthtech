package config

import (
	"os"
	"testing"
)

func TestConfigDefaults(t *testing.T) {
	// تنظيف appConfig لفرض إعادة التحميل
	appConfig = nil
	
	// تعيين بيئة اختبار بدون استخدام logger
	os.Setenv("ENVIRONMENT", "test")
	os.Setenv("DATABASE_URL", "postgres://test:test@localhost:5432/test")
	
	cfg := Load()
	
	if cfg == nil {
		t.Fatal("Expected config to be loaded, got nil")
	}
	
	if cfg.Environment != "test" {
		t.Errorf("Expected environment 'test', got '%s'", cfg.Environment)
	}
	
	if cfg.Port != "3000" {
		t.Errorf("Expected default port '3000', got '%s'", cfg.Port)
	}
}

func TestConfigMethods(t *testing.T) {
	cfg := &Config{
		Environment: "development",
		Port:        "8080",
		Version:     "1.0.0",
		Database: DatabaseConfig{
			URL: "postgres://user:pass@localhost:5432/db",
		},
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
}