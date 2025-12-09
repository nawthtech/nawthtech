package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"       // postgres driver
	_ "modernc.org/sqlite"      // sqlite pure go (بديل في التطوير)
	"github.com/nawthtech/nawthtech/backend/internal/config"
)

var DB *sql.DB

// InitializeSQL تهيئة اتصال database/sql
// driver: "postgres" أو "sqlite"
func InitializeSQL(cfg *config.Config, driver, dsn string) error {
	if driver == "" || dsn == "" {
		return fmt.Errorf("sql driver and dsn are required")
	}

	db, err := sql.Open(driver, dsn)
	if err != nil {
		return err
	}

	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(20)

	// اختبار الاتصال
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return err
	}

	DB = db
	return nil
}

func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}