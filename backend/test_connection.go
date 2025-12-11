package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	fmt.Println("🧪 Testing Cloudflare D1 Database Connection")

	// تحديد مسار قاعدة البيانات
	dbPath := os.Getenv("D1_DB_PATH")
	if dbPath == "" {
		dbPath = "file:./data/nawthtech.db?cache=shared&mode=rwc"
	}

	// الاتصال
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("❌ Failed to connect: %v", err)
	}
	defer db.Close()

	// اختبار الاتصال
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("❌ Database ping failed: %v", err)
	}

	fmt.Println("✅ Connected to Cloudflare D1 database successfully!")

	// اختبار استعلام بسيط
	var version string
	err = db.QueryRowContext(ctx, "SELECT sqlite_version()").Scan(&version)
	if err != nil {
		log.Fatalf("❌ Failed to get SQLite version: %v", err)
	}

	fmt.Printf("📊 SQLite Version: %s\n", version)

	// التحقق من الجداول
	rows, err := db.QueryContext(ctx, "SELECT name FROM sqlite_master WHERE type='table'")
	if err != nil {
		log.Printf("⚠️ Could not list tables: %v", err)
	} else {
		defer rows.Close()
		
		fmt.Println("📋 Database Tables:")
		count := 0
		for rows.Next() {
			var tableName string
			if err := rows.Scan(&tableName); err == nil {
				fmt.Printf("  - %s\n", tableName)
				count++
			}
		}
		fmt.Printf("📊 Total Tables: %d\n", count)
	}

	fmt.Println("🎉 Database connection test completed successfully!")
}