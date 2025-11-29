package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	// حفظ الإعدادات الحالية
	originalEnv := os.Getenv("ENVIRONMENT")
	defer os.Setenv("ENVIRONMENT", originalEnv)
	
	// تعيين بيئة اختبار
	os.Setenv("ENVIRONMENT", "test")
	
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