package logger

import (
	"log/slog"
	"testing"
)

func TestLoggerInitialization(t *testing.T) {
	// Test that loggers are initialized
	if Stdout == nil {
		t.Error("Expected Stdout logger to be initialized")
	}
	
	if Stderr == nil {
		t.Error("Expected Stderr logger to be initialized")
	}
}

func TestErrAttr(t *testing.T) {
	// Test with nil error
	attr := ErrAttr(nil)
	if attr.Key != "error" {
		t.Errorf("Expected key 'error', got '%s'", attr.Key)
	}
	
	// Test with actual error
	testErr := ErrAttr("test error")
	if testErr.Key != "error" {
		t.Errorf("Expected key 'error', got '%s'", testErr.Key)
	}
	
	// Test that it returns slog.Attr type
	var _ slog.Attr = testErr
}

func TestLoggersCanLog(t *testing.T) {
	// Test that loggers can be used without panicking
	// This is a basic smoke test to ensure they're properly initialized
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Logger panicked: %v", r)
		}
	}()
	
	// These should not panic
	Stdout.Info("test message")
	Stderr.Error("test error message")
}