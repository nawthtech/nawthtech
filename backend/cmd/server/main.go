package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/nawthtech/nawthtech/backend/internal/config"
	"github.com/nawthtech/nawthtech/backend/internal/handlers"
	"github.com/nawthtech/nawthtech/backend/internal/utils"
	"gorm.io/gorm"
)

func main() {
	// تحميل الإعدادات
	cfg := config.Load()

	// تهيئة قاعدة البيانات
	db, err := utils.InitDatabase(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("فشل في الاتصال بقاعدة البيانات:", err)
	}

	// تهيئة تطبيق Gin
	app := gin.Default()

	// الوسائط العامة
	app.Use(gin.Logger())
	app.Use(gin.Recovery())
	app.Use(CORSMiddleware())

	// تسجيل جميع المسارات
	handlers.RegisterAllRoutes(app, db, config, router, serviceContainer )

	// بدء الخادم
	log.Printf("🚀 بدء الخادم على المنفذ %s", cfg.Port)
	log.Printf("🌍 البيئة: %s", cfg.Environment)
	log.Printf("📦 الإصدار: %s", cfg.Version)

	if err := app.Run(":" + cfg.Port); err != nil {
		log.Fatal("فشل في بدء الخادم:", err)
	}
}

// CORSMiddleware وسيط CORS
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		cfg := config.Load()
		
		c.Writer.Header().Set("Access-Control-Allow-Origin", strings.Join(cfg.Cors.AllowedOrigins, ","))
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}