package utils

import (
	"testing"
)

func TestGetMemoryUsageMB(t *testing.T) {
	mem := GetMemoryUsageMB()
	
	if mem.UsedMB < 0 {
		t.Errorf("Expected used memory to be non-negative, got %f", mem.UsedMB)
	}
	
	if mem.TotalMB < 0 {
		t.Errorf("Expected total memory to be non-negative, got %f", mem.TotalMB)
	}
	
	if mem.UsagePercentage < 0 || mem.UsagePercentage > 100 {
		t.Errorf("Expected usage percentage between 0 and 100, got %f", mem.UsagePercentage)
	}
}

func TestGetGoroutineCount(t *testing.T) {
	count := GetGoroutineCount()
	
	if count <= 0 {
		t.Errorf("Expected goroutine count to be positive, got %d", count)
	}
}

// إذا لم تكن generateRandomString موجودة، استبدل هذا الاختبار بـ:
func TestFormatBytes(t *testing.T) {
	// Test that FormatBytes works correctly
	tests := []struct {
		input    int64
		expected string
	}{
		{1024, "1.0 KB"},
		{1048576, "1.0 MB"},
		{1073741824, "1.0 GB"},
		{0, "0 B"},
	}
	
	for _, test := range tests {
		result := FormatBytes(test.input)
		if result != test.expected {
			t.Errorf("For %d bytes, expected '%s', got '%s'", test.input, test.expected, result)
		}
	}
}