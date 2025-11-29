package logger

import (
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
	err := ErrAttr(nil)
	if err != nil {
		t.Error("Expected ErrAttr to handle nil error")
	}
	
	testErr := ErrAttr("test error")
	if testErr == nil {
		t.Error("Expected ErrAttr to return non-nil for non-nil error")
	}
}