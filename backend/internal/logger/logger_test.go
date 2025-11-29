package logger

import (
	"errors"
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
	
	// Test with actual error - استخدام errors.New بدلاً من string مباشرة
	testErr := errors.New("test error")
	attrWithError := ErrAttr(testErr)
	if attrWithError.Key != "error" {
		t.Errorf("Expected key 'error', got '%s'", attrWithError.Key)
	}
	
	// Test that it returns slog.Attr type
	var _ slog.Attr = attrWithError
}

func TestLoggersCanLog(t *testing.T) {
	// Test that loggers can be used without panicking
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Logger panicked: %v", r)
		}
	}()
	
	// These should not panic إذا كانت Loggers مهيئة
	if Stdout != nil {
		Stdout.Info("test message")
	}
	if Stderr != nil {
		Stderr.Error("test error message")
	}
}