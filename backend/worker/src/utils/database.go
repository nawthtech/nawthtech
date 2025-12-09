package utils

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3" // Cloudflare D1 تستخدم SQLite syntax
)

// D1Database يمثل الاتصال بقاعدة D1
type D1Database struct {
	DB *sql.DB
}

var (
	dbInstance *D1Database
	once       sync.Once
)

// GetDatabase يُعيد مثيل قاعدة البيانات (singleton)
func GetDatabase() *D1Database {
	once.Do(func() {
		databaseURL := os.Getenv("D1_DATABASE_URL")
		if databaseURL == "" {
			log.Fatal("D1_DATABASE_URL is required")
		}

		// الاتصال بـ D1 (SQLite syntax)
		db, err := sql.Open("sqlite3", databaseURL)
		if err != nil {
			log.Fatalf("Failed to connect to D1: %v", err)
		}

		// تحقق من الاتصال
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := db.PingContext(ctx); err != nil {
			log.Fatalf("Failed to ping D1: %v", err)
		}

		dbInstance = &D1Database{DB: db}
		log.Println("✅ Connected to D1 database successfully!")
	})

	return dbInstance
}

// Query ينفذ استعلام SELECT ويعيد النتائج
func (d *D1Database) Query(query string, args ...interface{}) ([]map[string]interface{}, error) {
	rows, err := d.DB.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("query failed: %v", err)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("failed to get columns: %v", err)
	}

	var results []map[string]interface{}
	for rows.Next() {
		cols := make([]interface{}, len(columns))
		colPtrs := make([]interface{}, len(columns))
		for i := range cols {
			colPtrs[i] = &cols[i]
		}

		if err := rows.Scan(colPtrs...); err != nil {
			return nil, fmt.Errorf("row scan failed: %v", err)
		}

		rowMap := make(map[string]interface{})
		for i, colName := range columns {
			rowMap[colName] = cols[i]
		}

		results = append(results, rowMap)
	}

	return results, nil
}

// Exec ينفذ استعلام INSERT/UPDATE/DELETE
func (d *D1Database) Exec(query string, args ...interface{}) (sql.Result, error) {
	res, err := d.DB.Exec(query, args...)
	if err != nil {
		return nil, fmt.Errorf("exec failed: %v", err)
	}
	return res, nil
}

// HealthCheck يتحقق من صحة الاتصال بقاعدة البيانات
func (d *D1Database) HealthCheck() string {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := d.DB.PingContext(ctx); err != nil {
		return "unhealthy"
	}
	return "healthy"
}