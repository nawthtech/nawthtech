package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/nawthtech/nawthtech/backend/internal/config"
)

// DB هو نوع مغلف لـ sql.DB
type DB struct {
	*sql.DB
}

// Config يحتوي على إعدادات قاعدة البيانات
type DatabaseConfig struct {
	Driver   string
	DSN      string
	MaxConns int
	MaxIdle  int
	Timeout  time.Duration
}

var (
	// db هو الاتصال العالمي بقاعدة البيانات
	db *DB

	// defaultConfig الإعدادات الافتراضية
	defaultConfig = DatabaseConfig{
		Driver:   "postgres",
		MaxConns: 25,
		MaxIdle:  5,
		Timeout:  30 * time.Second,
	}
)

// InitializeSQL تهيئة اتصال قاعدة البيانات SQL
func InitializeSQL(cfg *config.Config, driver, dsn string) (*DB, error) {
	if dsn == "" {
		return nil, fmt.Errorf("DSN cannot be empty")
	}

	// فتح الاتصال بقاعدة البيانات
	sqlDB, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// اختبار الاتصال
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		sqlDB.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// ضبط إعدادات الاتصال
	sqlDB.SetMaxOpenConns(defaultConfig.MaxConns)
	sqlDB.SetMaxIdleConns(defaultConfig.MaxIdle)
	sqlDB.SetConnMaxLifetime(defaultConfig.Timeout)

	// إنشاء المغلف
	db = &DB{DB: sqlDB}

	log.Printf("✅ Database connected successfully (driver: %s)", driver)
	return db, nil
}

// InitializeFromConfig تهيئة قاعدة البيانات من config
func InitializeFromConfig(cfg *config.Config) (*DB, error) {
	dsn := cfg.Database.URL
	if dsn == "" {
		log.Println("⚠️ Database DSN not configured, using in-memory mode")
		return nil, nil
	}

	return InitializeSQL(cfg, "postgres", dsn)
}

// GetDB ترجع نسخة من اتصال قاعدة البيانات
func GetDB() *DB {
	return db
}

// Close يغلق اتصال قاعدة البيانات
func Close() {
	if db != nil && db.DB != nil {
		if err := db.DB.Close(); err != nil {
			log.Printf("Error closing database: %v", err)
		} else {
			log.Println("✅ Database connection closed")
		}
		db = nil
	}
}

// Ping تتحقق من أن قاعدة البيانات لا تزال متصلة
func Ping(ctx context.Context) error {
    if db == nil || db.DB == nil {
        return fmt.Errorf("database not initialized")
    }
    return db.DB.PingContext(ctx)
}

// ExecContext تنفيذ استعلام بدون إرجاع صفوف
func (d *DB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	if d == nil || d.DB == nil {
		return nil, fmt.Errorf("database not initialized")
	}
	return d.DB.ExecContext(ctx, query, args...)
}

// QueryContext تنفيذ استعلام مع إرجاع صفوف
func (d *DB) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	if d == nil || d.DB == nil {
		return nil, fmt.Errorf("database not initialized")
	}
	return d.DB.QueryContext(ctx, query, args...)
}

// QueryRowContext تنفيذ استعلام مع إرجاع صف واحد
func (d *DB) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	if d == nil || d.DB == nil {
		// إرجاع صف فارغ في حالة الخطأ
		return &sql.Row{}
	}
	return d.DB.QueryRowContext(ctx, query, args...)
}

// BeginTx بدء معاملة
func (d *DB) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	if d == nil || d.DB == nil {
		return nil, fmt.Errorf("database not initialized")
	}
	return d.DB.BeginTx(ctx, opts)
}

// Transaction تنفيذ دالة داخل معاملة
func Transaction(ctx context.Context, fn func(*sql.Tx) error) error {
	db := GetDB()
	if db == nil {
		return fmt.Errorf("database not initialized")
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}
	}()

	if err := fn(tx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("transaction error: %w, rollback error: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

// Stats ترجع إحصائيات قاعدة البيانات
func (d *DB) Stats() sql.DBStats {
	if d == nil || d.DB == nil {
		return sql.DBStats{}
	}
	return d.DB.Stats()
}

// Helper Functions

// IsConnected تتحقق مما إذا كانت قاعدة البيانات متصلة
func IsConnected() bool {
    if db == nil {
        return false
    }
    
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()

      return db.DB.PingContext(ctx) == nil
   }

// HealthCheck تتحقق من صحة قاعدة البيانات
func HealthCheck(ctx context.Context) (bool, error) {
    if db == nil {
        return false, fmt.Errorf("database not initialized")
    }
    
    if err := db.DB.PingContext(ctx); err != nil {
        return false, err
    }
    
    return true, nil
}

// RunMigrations تشغيل عمليات الترحيل
func RunMigrations(ctx context.Context) error {
	if db == nil {
		return fmt.Errorf("database not initialized")
	}

	// جدول المستخدمين
	usersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		email VARCHAR(255) UNIQUE NOT NULL,
		name VARCHAR(255),
		password_hash VARCHAR(255),
		role VARCHAR(50) DEFAULT 'user',
		status VARCHAR(50) DEFAULT 'active',
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		last_login TIMESTAMP
	)
	`

	// جدول الجلسات
	sessionsTable := `
	CREATE TABLE IF NOT EXISTS sessions (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		user_id UUID REFERENCES users(id) ON DELETE CASCADE,
		token TEXT NOT NULL,
		expires_at TIMESTAMP NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		INDEX idx_token (token),
		INDEX idx_user_id (user_id)
	)
	`

	// جدول السجلات
	auditLogTable := `
	CREATE TABLE IF NOT EXISTS audit_logs (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		user_id UUID REFERENCES users(id) ON DELETE SET NULL,
		action VARCHAR(255) NOT NULL,
		resource VARCHAR(255),
		details JSONB,
		ip_address VARCHAR(45),
		user_agent TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)
	`

	migrations := []string{usersTable, sessionsTable, auditLogTable}

	for _, migration := range migrations {
		if _, err := db.ExecContext(ctx, migration); err != nil {
			return fmt.Errorf("failed to run migration: %w", err)
		}
	}

	log.Println("✅ Database migrations completed successfully")
	return nil
}

// CreateIndexes إنشاء الفهارس
func CreateIndexes(ctx context.Context) error {
	if db == nil {
		return fmt.Errorf("database not initialized")
	}

	indexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_users_email ON users(email)",
		"CREATE INDEX IF NOT EXISTS idx_users_status ON users(status)",
		"CREATE INDEX IF NOT EXISTS idx_sessions_expires ON sessions(expires_at)",
		"CREATE INDEX IF NOT EXISTS idx_audit_logs_created ON audit_logs(created_at)",
	}

	for _, index := range indexes {
		if _, err := db.ExecContext(ctx, index); err != nil {
			return fmt.Errorf("failed to create index: %w", err)
		}
	}

	log.Println("✅ Database indexes created successfully")
	return nil
}
