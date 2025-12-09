package utils

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/nawthtech/backend/internal/config"
	"github.com/nawthtech/backend/internal/db"

	"github.com/cloudflare/cloudflare-go/d1" // تأكد من تثبيت المكتبة المناسبة
)

// DatabaseManager مدير قاعدة البيانات
type DatabaseManager struct {
	cfg *config.Config
	db  *d1.DB
	mu  sync.Mutex
}

// اتصال مخبأ عالمي
var cachedDBManager *DatabaseManager
var once sync.Once

// NewDatabaseManager إنشاء مدير قاعدة بيانات جديد
func NewDatabaseManager(cfg *config.Config) *DatabaseManager {
	return &DatabaseManager{
		cfg: cfg,
	}
}

// GetDatabaseManager استرجاع أو إنشاء مدير قاعدة البيانات
func GetDatabaseManager(cfg *config.Config) *DatabaseManager {
	once.Do(func() {
		cachedDBManager = NewDatabaseManager(cfg)
	})
	return cachedDBManager
}

// Connect تهيئة الاتصال بـ D1
func (m *DatabaseManager) Connect() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.db != nil {
		return nil // متصل مسبقاً
	}

	if m.cfg.D1.DatabaseName == "" || m.cfg.D1.BindingName == "" {
		return fmt.Errorf("D1 configuration missing")
	}

	d1db, err := d1.Open(m.cfg.D1.BindingName)
	if err != nil {
		return fmt.Errorf("failed to connect to D1: %v", err)
	}

	m.db = d1db
	log.Println("✅ Connected to D1 successfully!")
	return nil
}

// GetConnection استرجاع الاتصال الحالي
func (m *DatabaseManager) GetConnection() (*d1.DB, error) {
	if m.db == nil {
		return nil, fmt.Errorf("database not connected")
	}
	return m.db, nil
}

// Disconnect إغلاق الاتصال
func (m *DatabaseManager) Disconnect() {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.db != nil {
		// في D1 لا يوجد close فعلي لأن الاتصال يتم عبر Cloudflare Workers
		m.db = nil
		log.Println("🔌 Disconnected from D1")
	}
}

// HealthCheck فحص صحة قاعدة البيانات
func (m *DatabaseManager) HealthCheck(ctx context.Context) (map[string]interface{}, error) {
	if err := m.Connect(); err != nil {
		return map[string]interface{}{
			"status": "disconnected",
			"type":   "none",
		}, err
	}

	query := `SELECT 1`
	_, err := m.db.QueryRow(ctx, query)
	if err != nil {
		return map[string]interface{}{
			"status": "unhealthy",
			"type":   "d1",
			"error":  err.Error(),
		}, err
	}

	return map[string]interface{}{
		"status": "healthy",
		"type":   "d1",
	}, nil
}

// WithDatabaseMiddleware تنفيذ handler مع قاعدة البيانات
func WithDatabaseMiddleware(cfg *config.Config, handler func(ctx context.Context, db *d1.DB) (interface{}, error)) func(ctx context.Context) (interface{}, error) {
	return func(ctx context.Context) (interface{}, error) {
		manager := GetDatabaseManager(cfg)
		if err := manager.Connect(); err != nil {
			log.Println("Database connection error:", err)
			return nil, fmt.Errorf("DATABASE_CONNECTION_FAILED: %v", err)
		}

		dbConn, err := manager.GetConnection()
		if err != nil {
			log.Println("Database not connected:", err)
			return nil, fmt.Errorf("DATABASE_CONNECTION_FAILED: %v", err)
		}

		result, err := handler(ctx, dbConn)
		if err != nil {
			log.Println("Database handler error:", err)
			return nil, err
		}

		return result, nil
	}
}