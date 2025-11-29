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

func TestGenerateID(t *testing.T) {
	id1 := GenerateID()
	id2 := GenerateID()
	
	if id1 == id2 {
		t.Errorf("Expected unique IDs, got duplicates: %s", id1)
	}
	
	if len(id1) == 0 {
		t.Errorf("Expected non-empty ID, got empty string")
	}
}