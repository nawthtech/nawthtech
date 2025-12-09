package main

import (
	"context"
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
)

func initLogger(env string) {
	logger.Init(env)
}

func Run() error {
	cfg := config.Load()
	initLogger(cfg.Environment)

	// Initialize SQL (development default sqlite)
	driver := os.Getenv("SQL_DRIVER")
	dsn := os.Getenv("SQL_DSN")
	if driver == "" {
		driver = "sqlite"
		dsn = "file:dev.db?_foreign_keys=1"
	}
	if err := db.InitializeSQL(cfg, driver, dsn); err != nil {
		logger.Error(context.Background(), "failed init sql", "error", err)
		return err
	}
	defer db.Close()

	// Gin app
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}
	app := gin.New()
	app.Use(gin.Recovery())

	// cors middleware
	app.Use(middleware.CORSMiddleware(cfg))

	// handlers
	hc := handlers.NewHandlerContainer(cfg, db.DB)
	handlers.RegisterAllRoutes(app, cfg, hc)

	// start server
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: app,
	}

	// graceful shutdown
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error(context.Background(), "listen error", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	return srv.Shutdown(ctx)
}

func main() {
	if err := Run(); err != nil {
		os.Exit(1)
	}
}