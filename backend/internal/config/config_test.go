package config

import (
	"testing"
)

func TestLoad(t *testing.T) {
	t.Skip("Skipping test - requires environment setup")
	
	cfg := Load()
	if cfg == nil {
		t.Error("Expected config to be loaded")
	}
}

func TestConfigDefaults(t *testing.T) {
	t.Skip("Skipping test - requires environment setup")
	
	cfg := Load()
	if cfg.Port == "" {
		t.Error("Expected default port to be set")
	}
	if cfg.Environment == "" {
		t.Error("Expected default environment to be set")
	}
}

func TestConfigValidation(t *testing.T) {
	t.Skip("Skipping test - requires environment setup")
	
	cfg := Load()
	if cfg == nil {
		t.Fatal("Config should be loaded")
	}
	
	// Test required fields
	if cfg.Database.URL == "" {
		t.Error("Database URL should be set")
	}
	if cfg.Database.Name == "" {
		t.Error("Database name should be set")
	}
}

func TestEnvironmentSpecificConfig(t *testing.T) {
	t.Skip("Skipping test - requires environment setup")
	
	cfg := Load()
	if cfg == nil {
		t.Fatal("Config should be loaded")
	}
	
	// Test environment-specific settings
	if cfg.Environment == "production" {
		if !cfg.Security.EnableHTTPS {
			t.Error("HTTPS should be enabled in production")
		}
	}
}

func TestCORSConfig(t *testing.T) {
	t.Skip("Skipping test - requires environment setup")
	
	cfg := Load()
	if cfg == nil {
		t.Fatal("Config should be loaded")
	}
	
	if len(cfg.CORS.AllowedOrigins) == 0 {
		t.Error("CORS allowed origins should be configured")
	}
}

func TestUploadConfig(t *testing.T) {
	t.Skip("Skipping test - requires environment setup")
	
	cfg := Load()
	if cfg == nil {
		t.Fatal("Config should be loaded")
	}
	
	if cfg.Upload.MaxFileSize == 0 {
		t.Error("Max file size should be configured")
	}
	
	if cfg.Upload.AllowedTypes == nil {
		t.Error("Allowed file types should be configured")
	}
}

func TestEmailConfig(t *testing.T) {
	t.Skip("Skipping test - requires environment setup")
	
	cfg := Load()
	if cfg == nil {
		t.Fatal("Config should be loaded")
	}
	
	if cfg.Email.Host == "" {
		t.Error("Email host should be configured")
	}
	
	if cfg.Email.Port == 0 {
		t.Error("Email port should be configured")
	}
}

func TestCacheConfig(t *testing.T) {
	t.Skip("Skipping test - requires environment setup")
	
	cfg := Load()
	if cfg == nil {
		t.Fatal("Config should be loaded")
	}
	
	if cfg.Cache.RedisURL == "" {
		t.Error("Redis URL should be configured")
	}
}

func TestSecurityConfig(t *testing.T) {
	t.Skip("Skipping test - requires environment setup")
	
	cfg := Load()
	if cfg == nil {
		t.Fatal("Config should be loaded")
	}
	
	if cfg.Security.JWTSecret == "" {
		t.Error("JWT secret should be configured")
	}
	
	if cfg.Security.BCryptCost == 0 {
		t.Error("BCrypt cost should be configured")
	}
}

func TestLoggingConfig(t *testing.T) {
	t.Skip("Skipping test - requires environment setup")
	
	cfg := Load()
	if cfg == nil {
		t.Fatal("Config should be loaded")
	}
	
	if cfg.Logging.Level == "" {
		t.Error("Logging level should be configured")
	}
}