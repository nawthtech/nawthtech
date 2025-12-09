package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nawthtech/backend/internal/config"
	"github.com/nawthtech/backend/internal/handlers"
	"github.com/nawthtech/backend/internal/logger"
	"github.com/nawthtech/backend/internal/middleware"
	"github.com/nawthtech/backend/internal/services"
)

func initLogger(env string) {
	logger.Init(env)
}

func Run() error {
	cfg := config.Load()
	initLogger(cfg.Environment)

	// إزالة تهيئة قاعدة البيانات المحلية
	// بدلاً من ذلك، سنستخدم Worker API للوصول إلى قاعدة البيانات

	// Gin app
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}
	app := gin.New()
	app.Use(gin.Recovery())

	// CORS middleware
	app.Use(middleware.CORSMiddleware(cfg))

	// Initialize services
	workerService, err := services.NewWorkerService(cfg)
	if err != nil {
		logger.Error(context.Background(), "failed to initialize worker service", "error", err)
		// يمكن الاستمرار بدون Worker إذا كان اختيارياً
	}

	// Initialize handlers with worker service
	hc := handlers.NewHandlerContainer(cfg, workerService)
	handlers.RegisterAllRoutes(app, cfg, hc)

	// Health check endpoint
	app.GET("/health", func(c *gin.Context) {
		health := map[string]interface{}{
			"status":      "healthy",
			"timestamp":   time.Now().UTC(),
			"environment": cfg.Environment,
			"version":     "1.0.0",
			"services": map[string]interface{}{
				"worker_api": workerService != nil,
				"database":   false, // تم نقله إلى Worker
				"cache":      true,
			},
		}
		c.JSON(http.StatusOK, health)
	})

	// Worker proxy endpoints (اختياري - إذا كنت تريد توجيه الطلبات عبر Backend)
	if workerService != nil {
		setupWorkerProxyRoutes(app, workerService)
	}

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
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error(context.Background(), "listen error", "error", err)
			os.Exit(1)
		}
	}()

	logger.Info(context.Background(), "server started", "port", cfg.Port, "environment", cfg.Environment)

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info(context.Background(), "shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error(context.Background(), "server shutdown error", "error", err)
		return err
	}

	logger.Info(context.Background(), "server exited properly")
	return nil
}

func setupWorkerProxyRoutes(app *gin.Engine, workerService *services.WorkerService) {
	// يمكنك إضافة routes توجيهية إذا لزم الأمر
	// ولكن الأفضل هو توجيه الطلبات مباشرة إلى Worker
	workerGroup := app.Group("/worker")
	{
		workerGroup.GET("/health", func(c *gin.Context) {
			health, err := workerService.HealthCheck()
			if err != nil {
				c.JSON(http.StatusServiceUnavailable, gin.H{
					"error":   "Worker unavailable",
					"details": err.Error(),
				})
				return
			}
			c.JSON(http.StatusOK, health)
		})
	}
}

func main() {
	if err := Run(); err != nil {
		logger.Error(context.Background(), "server failed", "error", err)
		os.Exit(1)
	}
}