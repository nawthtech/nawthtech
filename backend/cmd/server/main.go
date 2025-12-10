package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nawthtech/nawthtech/backend/internal/config"
	"github.com/nawthtech/nawthtech/backend/internal/db"
	"github.com/nawthtech/nawthtech/backend/internal/handlers"
	"github.com/nawthtech/nawthtech/backend/internal/logger"
	"github.com/nawthtech/nawthtech/backend/internal/middleware"
	"github.com/nawthtech/nawthtech/backend/internal/slack"
)

func main() {
	// تحميل الإعدادات
	cfg := config.Load()

	// تهيئة قاعدة البيانات
	database, err := db.InitializeFromConfig(cfg)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	// التحقق من الاتصال
	if db.IsConnected() {
		log.Println("Database is connected")
	}

	// تشغيل عمليات الترحيل
	ctx := context.Background()
	if err := db.RunMigrations(ctx); err != nil {
		log.Printf("Warning: migrations failed: %v", err)
	}

	// استخدام قاعدة البيانات
	dbInstance := db.GetDB()
	if dbInstance != nil {
		// تنفيذ استعلام
		var version string
		err := dbInstance.QueryRowContext(ctx, "SELECT version()").Scan(&version)
		if err == nil {
			log.Printf("Database version: %s", version)
		}
	}
}

func main() {
	// تهيئة Slack client من environment variables
	err := slack.Init(
		slack.WithToken(os.Getenv("SLACK_TOKEN")),
		slack.WithchannelURL(os.Getenv("SLACK_CHANNEL_URL")),
		slack.WithAppName("nawthtech-backend"),
		slack.WithEnvironment(os.Getenv("RAILWAY_ENVIRONMENT")),
	)

	if err != nil {
		log.Printf("Failed to initialize Slack client: %v", err)
	}
}

func initLogger(env string) {
	logger.Init(env)
}

func Run() error {
	cfg := config.Load()
	initLogger(cfg.Environment)

	// Gin app
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	app := gin.New()
	app.Use(gin.Recovery())
	app.Use(middleware.CORSMiddleware(cfg))

	// Initialize handlers
	hc := handlers.NewHandlerContainer(cfg)
	hc.RegisterAllRoutes(app)

	// Start server
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      app,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Graceful shutdown
	go func() {
		logger.Info(context.Background(), "server starting", "port", cfg.Port, "environment", cfg.Environment)

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error(context.Background(), "listen error", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info(context.Background(), "shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	return srv.Shutdown(ctx)
}

func main() {
	if err := Run(); err != nil {
		logger.Error(context.Background(), "server failed", "error", err)
		os.Exit(1)
	}
}
