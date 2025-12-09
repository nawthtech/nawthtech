package utils

import (
	"fmt"
	"os"
)

type D1Database struct {
	DBName string
}

func ConnectD1() (*D1Database, error) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}

	dbName := os.Getenv("DATABASE_NAME")
	if dbName == "" {
		dbName = "nawthtech"
	}

	return &D1Database{DBName: dbName}, nil
}

func (d *D1Database) HealthCheck() (string, error) {
	return "healthy", nil
}