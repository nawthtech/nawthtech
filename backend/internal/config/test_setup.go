// في internal/config/test_setup.go
package config

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// تهيئة الـ logger قبل تشغيل أي اختبار
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func setup() {
	// تهيئة الـ logger للاختبارات
}

func teardown() {
	// تنظيف بعد الاختبارات
}