package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"worker/src/handlers"
	"worker/src/utils"
)

// envVariables تُخزن إعدادات البيئة
var envVariables map[string]string

func init() {
	envVariables = map[string]string{
		"ENVIRONMENT": getEnv("ENVIRONMENT", "development"),
		"API_VERSION": getEnv("API_VERSION", "v1"),
	}
}

func main() {
	// تهيئة اتصال D1
	if err := utils.InitDatabase(); err != nil {
		log.Fatalf("❌ Failed to initialize database: %v", err)
	}

	// مسارات الخدمة
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		handlers.CheckHealthHandler(w, r, envVariables)
	})

	http.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
		handlers.ReadyHandler(w, r, envVariables)
	})

	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		handlers.TestHandler(w, r, envVariables)
	})

	port := getEnv("PORT", "8787")
	log.Printf("🚀 Worker running on port %s in %s mode", port, envVariables["ENVIRONMENT"])
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("❌ Server failed: %v", err)
	}
}

// getEnv يقرأ متغيرات البيئة مع قيمة افتراضية
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return strings.TrimSpace(value)
	}
	return defaultValue
}