package main

import (
	"context"
	"database/sql"
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
	"github.com/nawthtech/nawthtech/backend/internal/router"
	"github.com/nawthtech/nawthtech/backend/internal/services"
	"github.com/nawthtech/nawthtech/backend/internal/slack"
)

func main() {
	if err := Run(); err != nil {
		logger.Error(context.Background(), "server failed", "error", err)
		os.Exit(1)
	}
}

func Run() error {
	// تحميل الإعدادات
	cfg := config.Load()

	// تهيئة الـ logger
	initLogger(cfg.Environment)

	logger.Info(context.Background(), "starting nawthtech backend server",
		"environment", cfg.Environment,
		"port", cfg.Port,
		"version", cfg.Version)

	// تهيئة Slack client
	initSlack()

	// تهيئة قاعدة البيانات
	database, err := initDatabase(cfg)
	if err != nil {
		logger.Error(context.Background(), "failed to initialize database", "error", err)
		return err
	}
	defer closeDatabase(database)

	// إنشاء service container
	serviceContainer := initServices(database, cfg)

	// تكوين تطبيق Gin
	app := setupGinApp(cfg, database, serviceContainer)

	// بدء الخادم مع graceful shutdown
	return startServer(app, cfg)
}

func initLogger(env string) {
	logger.Init(env)
	logger.Info(context.Background(), "logger initialized", "environment", env)
}

func initSlack() {
	slackToken := os.Getenv("SLACK_TOKEN")
	slackChannel := os.Getenv("SLACK_CHANNEL")
	appName := os.Getenv("APP_NAME")
	if appName == "" {
		appName = "nawthtech-backend"
	}
	environment := os.Getenv("ENVIRONMENT")
	if environment == "" {
		environment = "development"
	}

	if slackToken != "" && slackChannel != "" {
		client, err := slack.New(
			slack.WithToken(slackToken),
			slack.WithChannelURL(slackChannel),
			slack.WithAppName(appName),
			slack.WithEnvironment(environment),
		)
		if err != nil {
			logger.Warn(context.Background(), "failed to initialize Slack client", "error", err)
		} else {
			logger.Info(context.Background(), "Slack client initialized successfully")
			// إرسال إشعار بدء التشغيل
			go func() {
				_, _, err := client.SendAlert("info", "Backend Server Started", 
					"nawthtech backend server has started successfully")
				if err != nil {
					logger.Warn(context.Background(), "failed to send Slack notification", "error", err)
				}
			}()
		}
	} else {
		logger.Info(context.Background(), "Slack not configured, running without notifications")
	}
}

func initDatabase(cfg *config.Config) (*sql.DB, error) {
	logger.Info(context.Background(), "initializing database connection")

	// تهيئة اتصال قاعدة البيانات
	database, err := db.InitializeFromConfig(cfg)
	if err != nil {
		return nil, err
	}

	// التحقق من الاتصال
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := database.PingContext(ctx); err != nil {
		return nil, err
	}

	logger.Info(context.Background(), "database connected successfully")

	// تشغيل عمليات الترحيل
	logger.Info(context.Background(), "running database migrations")
	if err := db.RunMigrations(ctx, database); err != nil {
		logger.Warn(context.Background(), "migrations failed or already applied", "error", err)
	} else {
		logger.Info(context.Background(), "migrations completed successfully")
	}

	// التحقق من الجداول الأساسية
	checkTables(ctx, database)

	return database, nil
}

func checkTables(ctx context.Context, db *sql.DB) {
	query := `SELECT name FROM sqlite_master WHERE type='table' ORDER BY name`
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		logger.Warn(context.Background(), "failed to check tables", "error", err)
		return
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			continue
		}
		tables = append(tables, tableName)
	}

	logger.Info(context.Background(), "database tables", "count", len(tables), "tables", tables)
}

func closeDatabase(db *sql.DB) {
	if db != nil {
		if err := db.Close(); err != nil {
			logger.Error(context.Background(), "failed to close database", "error", err)
		} else {
			logger.Info(context.Background(), "database connection closed")
		}
	}
}

func initServices(db *sql.DB, cfg *config.Config) *services.ServiceContainer {
	logger.Info(context.Background(), "initializing services")

	serviceContainer := services.NewServiceContainerWithConfig(db, cfg, logger.GetLogger())

	// اختبار الخدمات الأساسية
	testBasicServices(serviceContainer)

	return serviceContainer
}

func testBasicServices(sc *services.ServiceContainer) {
	// اختبار خدمة التخزين المؤقت
	ctx := context.Background()
	testKey := "server_start_test"
	testValue := "nawthtech_backend_" + time.Now().Format(time.RFC3339)

	if err := sc.Cache.Set(testKey, testValue, 1*time.Minute); err != nil {
		logger.Warn(context.Background(), "cache service test failed", "error", err)
	} else {
		logger.Info(context.Background(), "cache service test passed")
	}
}

func setupGinApp(cfg *config.Config, db *sql.DB, serviceContainer *services.ServiceContainer) *gin.Engine {
	// تكوين Gin
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	app := gin.New()

	// Middleware الأساسية
	app.Use(gin.Recovery())
	app.Use(middleware.CORSMiddleware(cfg))
	app.Use(middleware.RequestIDMiddleware())
	app.Use(middleware.LoggerMiddleware())
	app.Use(middleware.RateLimitMiddleware(cfg))

	// مسارات الصحة
	setupHealthRoutes(app, db)

	// API Routes
	setupAPIRoutes(app, serviceContainer)

	// مسارات غير موجودة
	app.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "endpoint not found",
			"path":    c.Request.URL.Path,
		})
	})

	return app
}

func setupHealthRoutes(app *gin.Engine, db *sql.DB) {
	// مسارات الصحة الأساسية
	app.GET("/health", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		if err := db.PingContext(ctx); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status":  "unhealthy",
				"message": "database connection failed",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"timestamp": time.Now().UTC(),
			"service":   "nawthtech-backend",
			"database":  "connected",
		})
	})

	app.GET("/health/ready", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "ready",
			"timestamp": time.Now().UTC(),
			"message":   "service is ready to accept requests",
		})
	})

	app.GET("/health/live", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "live",
			"timestamp": time.Now().UTC(),
			"message":   "service is alive",
		})
	})
}

func setupAPIRoutes(app *gin.Engine, serviceContainer *services.ServiceContainer) {
	// Group لـ API v1
	apiV1 := app.Group("/api/v1")

	// إنشاء Handlers
	handlerContainer := handlers.NewHandlerContainer(serviceContainer)

	// Register routes
	routes.RegisterV1Routes(apiV1, handlerContainer, middleware.AuthMiddleware())
}

func startServer(app *gin.Engine, cfg *config.Config) error {
	// إعداد الخادم
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      app,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Graceful shutdown channel
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// بدء الخادم في goroutine
	serverErr := make(chan error, 1)
	go func() {
		logger.Info(context.Background(), "server starting",
			"address", srv.Addr,
			"environment", cfg.Environment,
			"read_timeout", srv.ReadTimeout,
			"write_timeout", srv.WriteTimeout)

		if cfg.IsProduction() && cfg.TLS.CertFile != "" && cfg.TLS.KeyFile != "" {
			logger.Info(context.Background(), "starting server with TLS")
			if err := srv.ListenAndServeTLS(cfg.TLS.CertFile, cfg.TLS.KeyFile); err != nil && err != http.ErrServerClosed {
				serverErr <- err
			}
		} else {
			logger.Info(context.Background(), "starting server without TLS")
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				serverErr <- err
			}
		}
	}()

	// انتظار إشارة الإغلاق
	select {
	case err := <-serverErr:
		return err
	case sig := <-quit:
		logger.Info(context.Background(), "received shutdown signal", "signal", sig.String())
		
		// إعداد context للإغلاق
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		
		// محاولة إرسال إشعار Slack
		sendShutdownNotification()
		
		// إغلاق الخادم
		if err := srv.Shutdown(ctx); err != nil {
			logger.Error(context.Background(), "server shutdown failed", "error", err)
			return err
		}
		
		logger.Info(context.Background(), "server shutdown completed")
		return nil
	}
}

func sendShutdownNotification() {
	slackToken := os.Getenv("SLACK_TOKEN")
	slackChannel := os.Getenv("SLACK_CHANNEL")
	
	if slackToken != "" && slackChannel != "" {
		client, err := slack.New(
			slack.WithToken(slackToken),
			slack.WithChannelURL(slackChannel),
			slack.WithAppName("nawthtech-backend"),
			slack.WithEnvironment(os.Getenv("ENVIRONMENT")),
		)
		if err == nil {
			_, _, err := client.SendAlert("warning", "Backend Server Shutdown",
				"nawthtech backend server is shutting down gracefully")
			if err != nil {
				logger.Warn(context.Background(), "failed to send shutdown notification", "error", err)
			}
		}
	}
}