// worker/src/utils/database.go
package utils

import (
	"database/sql"
	"fmt"
	"os"
	"sync"

	_ "github.com/mattn/go-sqlite3" // D1 يستخدم SQLite engine
)

// DatabaseManager مسؤول عن إدارة اتصال D1
type DatabaseManager struct {
	Env   map[string]string
	db    *sql.DB
	mutex sync.Mutex
}

// DatabaseHealth تمثل حالة قاعدة البيانات
type DatabaseHealth struct {
	Status   string `json:"status"`
	Database string `json:"database"`
}

// مدير قاعدة بيانات مخبأ (singleton)
var cachedDatabaseManager *DatabaseManager

// GetDatabaseManager إنشاء أو استرجاع مدير قاعدة البيانات
func GetDatabaseManager(env map[string]string) *DatabaseManager {
	if cachedDatabaseManager != nil {
		return cachedDatabaseManager
	}

	cachedDatabaseManager = &DatabaseManager{
		Env: env,
	}
	return cachedDatabaseManager
}

// Connect الاتصال بـ D1
func (m *DatabaseManager) Connect() (*sql.DB, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.db != nil {
		return m.db, nil
	}

	dsn := os.Getenv("D1_DATABASE_URL")
	if dsn == "" && m.Env != nil {
		dsn = m.Env["D1_DATABASE_URL"]
	}

	if dsn == "" {
		return nil, fmt.Errorf("D1_DATABASE_URL is required")
	}

	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open D1 connection: %v", err)
	}

	// اختبار الاتصال
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping D1: %v", err)
	}

	m.db = db
	return m.db, nil
}

// GetConnection الحصول على اتصال D1
func (m *DatabaseManager) GetConnection() (*sql.DB, error) {
	if m.db != nil {
		return m.db, nil
	}
	return m.Connect()
}

// Disconnect إغلاق الاتصال
func (m *DatabaseManager) Disconnect() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.db != nil {
		err := m.db.Close()
		m.db = nil
		return err
	}
	return nil
}

// HealthCheck فحص صحة قاعدة البيانات
func (m *DatabaseManager) HealthCheck() DatabaseHealth {
	db, err := m.GetConnection()
	if err != nil {
		return DatabaseHealth{
			Status:   "unhealthy",
			Database: "d1",
		}
	}

	if err := db.Ping(); err != nil {
		return DatabaseHealth{
			Status:   "unhealthy",
			Database: "d1",
		}
	}

	return DatabaseHealth{
		Status:   "healthy",
		Database: "d1",
	}
}