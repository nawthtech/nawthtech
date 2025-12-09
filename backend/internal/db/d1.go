package db

import (
	"context"
	"fmt"
	"log"

	"github.com/nawthtech/nawthtech/backend/internal/config"

	"github.com/cloudflare/cloudflare-go/d1" // تأكد من تثبيت المكتبة الصحيحة
)

var (
	// DB هو الاتصال العالمي بقاعدة D1
	DB *d1.DB
)

// InitializeD1 تهيئة اتصال D1
func InitializeD1(cfg *config.Config) {
	if cfg.D1.DatabaseName == "" || cfg.D1.BindingName == "" {
		log.Fatal("D1 configuration missing: DatabaseName or BindingName is empty")
	}

	// فتح اتصال D1 باستخدام الـ BindingName
	db, err := d1.Open(cfg.D1.BindingName)
	if err != nil {
		log.Fatalf("فشل الاتصال بـ D1: %v", err)
	}

	DB = db
	fmt.Printf("✅ تم الاتصال بنجاح مع D1: %s\n", cfg.D1.DatabaseName)
}

// Exec تنفيذ أمر SQL (Insert/Update/Delete)
func Exec(ctx context.Context, query string, args ...interface{}) error {
	if DB == nil {
		return fmt.Errorf("D1 not initialized")
	}

	_, err := DB.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("Exec error: %w", err)
	}

	return nil
}

// QueryRow تنفيذ استعلام يعيد صف واحد
func QueryRow(ctx context.Context, query string, args ...interface{}) (map[string]interface{}, error) {
	if DB == nil {
		return nil, fmt.Errorf("D1 not initialized")
	}

	row, err := DB.QueryRow(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("QueryRow error: %w", err)
	}

	return row, nil
}

// Query تنفيذ استعلام يعيد مجموعة صفوف
func Query(ctx context.Context, query string, args ...interface{}) ([]map[string]interface{}, error) {
	if DB == nil {
		return nil, fmt.Errorf("D1 not initialized")
	}

	rows, err := DB.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("Query error: %w", err)
	}

	return rows, nil
}