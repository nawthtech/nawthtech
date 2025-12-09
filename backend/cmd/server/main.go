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
)

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