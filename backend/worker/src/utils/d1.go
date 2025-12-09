package utils

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type D1Manager struct {
	db   *sql.DB
	once sync.Once
}

var instance *D1Manager
var dbOnce sync.Once

func GetD1() *D1Manager {
	dbOnce.Do(func() {
		instance = &D1Manager{}
	})
	return instance
}

func (d *D1Manager) Connect() error {
	var err error
	d.once.Do(func() {
		dsn := os.Getenv("D1_DATABASE_URL")
		if dsn == "" {
			err = fmt.Errorf("D1_DATABASE_URL is required")
			return
		}

		d.db, err = sql.Open("sqlite3", dsn)
		if err != nil {
			return
		}

		d.db.SetConnMaxLifetime(time.Minute * 5)
		d.db.SetMaxOpenConns(10)
		d.db.SetMaxIdleConns(5)

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

func (d *D1Manager) Disconnect(ctx context.Context) error {
	if d.db != nil {
		fmt.Println("🔌 Disconnecting D1...")
		return d.db.Close()
	}
	return nil
}

func (d *D1Manager) GetDB() (*sql.DB, error) {
	if d.db == nil {
		return nil, fmt.Errorf("D1 database not connected")
	}
	return d.db, nil
}

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