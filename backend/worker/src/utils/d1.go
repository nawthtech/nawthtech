package utils

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3" // استخدم SQLite للـ D1 driver
)

// D1Manager هو مدير قاعدة D1
type D1Manager struct {
	db   *sql.DB
	once sync.Once
}

var instance *D1Manager
var dbOnce sync.Once

// GetD1 يعيد مثيل D1Manager (singleton)
func GetD1() *D1Manager {
	dbOnce.Do(func() {
		instance = &D1Manager{}
	})
	return instance
}

// Connect يفتح اتصال D1
func (d *D1Manager) Connect() error {
	var err error
	d.once.Do(func() {
		dsn := os.Getenv("D1_DATABASE_URL")
		if dsn == "" {
			err = fmt.Errorf("D1_DATABASE_URL is required")
			return
		}

		// فتح قاعدة البيانات
		d.db, err = sql.Open("sqlite3", dsn)
		if err != nil {
			return
		}

		// تعيين timeout
		d.db.SetConnMaxLifetime(time.Minute * 5)
		d.db.SetMaxOpenConns(10)
		d.db.SetMaxIdleConns(5)

		// اختبار الاتصال
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err = d.db.PingContext(ctx)
		if err != nil {
			d.db.Close()
			d.db = nil
			return
		}
		fmt.Println("✅ Connected to D1 successfully!")
	})
	return err
}

// Disconnect يغلق اتصال D1
func (d *D1Manager) Disconnect(ctx context.Context) error {
	if d.db != nil {
		fmt.Println("🔌 Disconnecting D1...")
		return d.db.Close()
	}
	return nil
}

// GetDB يعيد *sql.DB
func (d *D1Manager) GetDB() (*sql.DB, error) {
	if d.db == nil {
		return nil, fmt.Errorf("D1 database not connected")
	}
	return d.db, nil
}

// HealthCheck يتحقق من حالة قاعدة البيانات
func (d *D1Manager) HealthCheck(ctx context.Context) (string, error) {
	if d.db == nil {
		return "disconnected", fmt.Errorf("D1 database not connected")
	}

	err := d.db.PingContext(ctx)
	if err != nil {
		return "unhealthy", err
	}
	return "healthy", nil
}

// ExecuteQuery تنفيذ استعلام بدون نتائج
func (d *D1Manager) ExecuteQuery(ctx context.Context, query string, args ...interface{}) error {
	if d.db == nil {
		return fmt.Errorf("D1 database not connected")
	}
	_, err := d.db.ExecContext(ctx, query, args...)
	return err
}

// QueryRows تنفيذ استعلام وإرجاع الصفوف
func (d *D1Manager) QueryRows(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	if d.db == nil {
		return nil, fmt.Errorf("D1 database not connected")
	}
	return d.db.QueryContext(ctx, query, args...)
}

// QueryRow تنفيذ استعلام وإرجاع صف واحد
func (d *D1Manager) QueryRow(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return d.db.QueryRowContext(ctx, query, args...)
}