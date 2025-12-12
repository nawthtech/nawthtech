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

// ================================
// دوال Health الإضافية
// ================================

// GetDatabaseStats ترجع إحصائيات قاعدة البيانات مفصلة
func (d *DB) GetDatabaseStats(ctx context.Context) (map[string]interface{}, error) {
	if d == nil || d.DB == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	stats := d.DB.Stats()
	
	// الحصول على إحصائيات إضافية حسب نوع قاعدة البيانات
	var dbInfo map[string]interface{}
	
	// محاولة الحصول على معلومات قاعدة البيانات (يعمل مع SQLite)
	if d.GetDriver() == "sqlite3" {
		var version string
		err := d.DB.QueryRowContext(ctx, "SELECT sqlite_version()").Scan(&version)
		if err == nil {
			dbInfo = map[string]interface{}{
				"driver":    d.GetDriver(),
				"version":   version,
				"max_conns": stats.MaxOpenConnections,
				"open_conns": stats.OpenConnections,
				"in_use":    stats.InUse,
				"idle":      stats.Idle,
				"wait_count": stats.WaitCount,
				"wait_duration_ms": stats.WaitDuration.Milliseconds(),
				"max_idle_closed": stats.MaxIdleClosed,
				"max_lifetime_closed": stats.MaxLifetimeClosed,
			}
		}
	} else {
		dbInfo = map[string]interface{}{
			"driver":    d.GetDriver(),
			"max_conns": stats.MaxOpenConnections,
			"open_conns": stats.OpenConnections,
			"in_use":    stats.InUse,
			"idle":      stats.Idle,
			"wait_count": stats.WaitCount,
			"wait_duration_ms": stats.WaitDuration.Milliseconds(),
		}
	}

	return dbInfo, nil
}

// GetDriver ترجع نوع قاعدة البيانات
func (d *DB) GetDriver() string {
	// محاولة تحديد نوع قاعدة البيانات من الـ DSN
	if d == nil || d.DB == nil {
		return "unknown"
	}
	
	// يمكن تحسين هذا المنطق بناءً على تكوينك الفعلي
	return "sqlite3" // أو "postgres" حسب تكوينك
}

// TestQuery تجربة استعلام اختبار
func (d *DB) TestQuery(ctx context.Context) (time.Duration, error) {
	if d == nil || d.DB == nil {
		return 0, fmt.Errorf("database not initialized")
	}

	startTime := time.Now()
	
	var result int
	err := d.DB.QueryRowContext(ctx, "SELECT 1").Scan(&result)
	
	duration := time.Since(startTime)
	
	if err != nil {
		return duration, fmt.Errorf("test query failed: %w", err)
	}
	
	if result != 1 {
		return duration, fmt.Errorf("unexpected test result: %d", result)
	}
	
	return duration, nil
}

// CheckDatabaseHealth فحص صحة قاعدة البيانات شامل
func (d *DB) CheckDatabaseHealth(ctx context.Context) (map[string]interface{}, error) {
	health := make(map[string]interface{})
	
	// 1. التحقق من الاتصال
	startTime := time.Now()
	connected, err := IsConnected()
	pingTime := time.Since(startTime)
	
	health["connected"] = connected
	health["ping_time_ms"] = pingTime.Milliseconds()
	health["connection_error"] = nil
	
	if err != nil {
		health["connection_error"] = err.Error()
		health["overall_status"] = "unhealthy"
		return health, err
	}
	
	// 2. اختبار استعلام
	queryTime, queryErr := d.TestQuery(ctx)
	health["query_time_ms"] = queryTime.Milliseconds()
	health["query_error"] = nil
	
	if queryErr != nil {
		health["query_error"] = queryErr.Error()
		health["overall_status"] = "degraded"
	} else {
		health["overall_status"] = "healthy"
	}
	
	// 3. إحصائيات قاعدة البيانات
	if stats, err := d.GetDatabaseStats(ctx); err == nil {
		health["stats"] = stats
	}
	
	// 4. معلومات عامة
	health["driver"] = d.GetDriver()
	health["timestamp"] = time.Now().UTC().Format(time.RFC3339)
	
	return health, nil
}

// GetConnectionInfo معلومات الاتصال بقاعدة البيانات
func (d *DB) GetConnectionInfo() map[string]interface{} {
	info := make(map[string]interface{})
	
	if d == nil || d.DB == nil {
		info["status"] = "not_initialized"
		return info
	}
	
	stats := d.DB.Stats()
	info["status"] = "connected"
	info["driver"] = d.GetDriver()
	info["max_open_connections"] = stats.MaxOpenConnections
	info["open_connections"] = stats.OpenConnections
	info["in_use"] = stats.InUse
	info["idle"] = stats.Idle
	info["wait_count"] = stats.WaitCount
	info["wait_duration"] = stats.WaitDuration.String()
	
	return info
}

// LogHealthMetric تسجيل مقياس صحة قاعدة البيانات
func (d *DB) LogHealthMetric(ctx context.Context, metricName string, value float64, metadata map[string]interface{}) error {
	if d == nil || d.DB == nil {
		return fmt.Errorf("database not initialized")
	}
	
	// تحويل metadata إلى JSON
	metadataJSON := "{}"
	if metadata != nil {
		jsonBytes, err := json.Marshal(metadata)
		if err != nil {
			metadataJSON = fmt.Sprintf(`{"error": "%s"}`, err.Error())
		} else {
			metadataJSON = string(jsonBytes)
		}
	}
	
	// إدراج في جدول performance_metrics إذا كان موجوداً
	query := `INSERT INTO performance_metrics 
		(id, metric_name, metric_value, metadata, created_at)
		VALUES (?, ?, ?, ?, ?)`
	
	_, err := d.DB.ExecContext(ctx, query,
		fmt.Sprintf("metric_%d", time.Now().UnixNano()),
		metricName,
		value,
		metadataJSON,
		time.Now(),
	)
	
	return err
}

// GetHealthMetrics الحصول على مقاييس الصحة
func (d *DB) GetHealthMetrics(ctx context.Context, metricName string, hours int) ([]map[string]interface{}, error) {
	if d == nil || d.DB == nil {
		return nil, fmt.Errorf("database not initialized")
	}
	
	query := `SELECT metric_name, metric_value, metadata, created_at
			  FROM performance_metrics 
			  WHERE metric_name = ? AND created_at >= datetime('now', ?)
			  ORDER BY created_at ASC`
	
	rows, err := d.DB.QueryContext(ctx, query, metricName, fmt.Sprintf("-%d hours", hours))
	if err != nil {
		// إذا كان الجدول غير موجود، نرجع مصفوفة فارغة
		return []map[string]interface{}{}, nil
	}
	defer rows.Close()
	
	var metrics []map[string]interface{}
	for rows.Next() {
		var name string
		var value float64
		var metadataJSON string
		var createdAt time.Time
		
		err := rows.Scan(&name, &value, &metadataJSON, &createdAt)
		if err != nil {
			continue
		}
		
		metric := map[string]interface{}{
			"metric_name": name,
			"metric_value": value,
			"created_at": createdAt.Format(time.RFC3339),
		}
		
		// تحليل metadata من JSON
		if metadataJSON != "" && metadataJSON != "{}" {
			var metadata map[string]interface{}
			if err := json.Unmarshal([]byte(metadataJSON), &metadata); err == nil {
				metric["metadata"] = metadata
			}
		}
		
		metrics = append(metrics, metric)
	}
	
	return metrics, nil
}
