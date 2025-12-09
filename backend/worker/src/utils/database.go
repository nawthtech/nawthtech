package utils

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3" // SQLite driver required by D1
)

var (
	DB         *sql.DB
	DBDriver   = "sqlite3"
	DBFilePath string
)

// InitDatabase تهيئة اتصال قاعدة البيانات D1
func InitDatabase() error {
	// قراءة رابط قاعدة البيانات من البيئة
	DBFilePath = os.Getenv("D1_DATABASE_PATH")
	if DBFilePath == "" {
		DBFilePath = ":memory:" // افتراضيًا في الذاكرة إذا لم يكن محدد
	}

	var err error
	DB, err = sql.Open(DBDriver, DBFilePath)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	// تعيين مهلة ping للتأكد من صحة الاتصال
	DB.SetConnMaxLifetime(time.Minute * 5)
	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(5)

	// اختبار الاتصال
	if err := DB.Ping(); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}

	log.Println("✅ Connected to D1 database successfully!")
	return nil
}

// CloseDatabase إغلاق الاتصال بالقاعدة
func CloseDatabase() {
	if DB != nil {
		if err := DB.Close(); err != nil {
			log.Printf("⚠️ Failed to close database: %v", err)
		} else {
			log.Println("🔌 Database connection closed")
		}
	}
}

// HealthCheck تحقق من صحة قاعدة البيانات
func HealthCheck() (status string, err error) {
	if DB == nil {
		return "disconnected", fmt.Errorf("database not initialized")
	}

	var result int
	err = DB.QueryRow("SELECT 1").Scan(&result)
	if err != nil {
		return "unhealthy", err
	}

	if result == 1 {
		return "healthy", nil
	}
	return "unhealthy", fmt.Errorf("unexpected database response")
}

// ExecQuery تنفيذ استعلام غير إرجاعي (INSERT/UPDATE/DELETE)
func ExecQuery(query string, args ...any) (sql.Result, error) {
	if DB == nil {
		return nil, fmt.Errorf("database not initialized")
	}
	return DB.Exec(query, args...)
}

// QueryRows تنفيذ استعلام إرجاعي (SELECT)
func QueryRows(query string, args ...any) (*sql.Rows, error) {
	if DB == nil {
		return nil, fmt.Errorf("database not initialized")
	}
	return DB.Query(query, args...)
}

// QueryRow تنفيذ استعلام صف واحد
func QueryRow(query string, args ...any) *sql.Row {
	if DB == nil {
		return nil
	}
	return DB.QueryRow(query, args...)
}