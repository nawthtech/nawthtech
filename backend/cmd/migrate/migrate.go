package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3" // استخدم هذا للـ SQLite
	// _ "github.com/lib/pq"          // استخدم هذا لـ PostgreSQL
)

func main() {
	// 1. افتح الاتصال بقاعدة البيانات
	db, err := sql.Open("sqlite3", "./mydatabase.db") // SQLite
	// db, err := sql.Open("postgres", "user=postgres password=yourpass dbname=mydb sslmode=disable") // PostgreSQL
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 2. جداول SQL
	sqlStmt := `
	CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		email TEXT UNIQUE NOT NULL,
		phone TEXT,
		password TEXT NOT NULL,
		status TEXT DEFAULT 'active',
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS categories (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		slug TEXT UNIQUE NOT NULL,
		image TEXT,
		is_active BOOLEAN DEFAULT TRUE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS services (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		description TEXT,
		price REAL NOT NULL,
		category_id TEXT,
		is_active BOOLEAN DEFAULT TRUE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(category_id) REFERENCES categories(id)
	);

	CREATE TABLE IF NOT EXISTS orders (
		id TEXT PRIMARY KEY,
		user_id TEXT NOT NULL,
		service_id TEXT NOT NULL,
		status TEXT DEFAULT 'pending',
		amount REAL NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(user_id) REFERENCES users(id),
		FOREIGN KEY(service_id) REFERENCES services(id)
	);

	CREATE TABLE IF NOT EXISTS payments (
		id TEXT PRIMARY KEY,
		order_id TEXT NOT NULL,
		amount REAL NOT NULL,
		status TEXT DEFAULT 'pending',
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(order_id) REFERENCES orders(id)
	);
	`

	// 3. تنفيذ SQL
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Fatalf("%q: %s\n", err, sqlStmt)
	}

	fmt.Println("✅ جميع الجداول تم إنشاؤها بنجاح!")
}
