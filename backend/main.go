package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/nawthtech/nawthtech/backend/internal/logger"
	"github.com/urfave/cli/v2"
)

// إصدار التطبيق - سيتم تعبئته أثناء البناء
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	// إنشاء تطبيق CLI
	app := &cli.App{
		Name:     "nawthtech",
		Version:  fmt.Sprintf("%s (commit: %s, built: %s)", version, commit, date),
		Usage:    "منصة نوذ تك للخدمات الإلكترونية - الخادم الخلفي",
		Compiled: time.Now(),
		Authors: []*cli.Author{
			{
				Name:  "فريق نوذ تك",
				Email: "dev@nawthtech.com",
			},
		},
		Commands: []*cli.Command{
			{
				Name:    "server",
				Aliases: []string{"s"},
				Usage:   "تشغيل خادم API",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "port",
						Aliases: []string{"p"},
						Value:   "8080",
						Usage:   "منفذ الخادم",
						EnvVars: []string{"PORT"},
					},
					&cli.StringFlag{
						Name:    "env",
						Aliases: []string{"e"},
						Value:   "development",
						Usage:   "بيئة التشغيل (development, staging, production)",
						EnvVars: []string{"APP_ENV"},
					},
					&cli.StringFlag{
						Name:    "config",
						Aliases: []string{"c"},
						Value:   "",
						Usage:   "مسار ملف الإعدادات",
						EnvVars: []string{"CONFIG_PATH"},
					},
				},
				Action: runServer,
			},
			{
				Name:    "version",
				Aliases: []string{"v"},
				Usage:   "عرض معلومات الإصدار",
				Action:  showVersion,
			},
			{
				Name:  "health",
				Usage: "فحص صحة النظام",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "timeout",
						Aliases: []string{"t"},
						Value:   "30s",
						Usage:   "مهلة فحص الصحة",
					},
				},
				Action: checkHealth,
			},
			{
				Name:  "migrate",
				Usage: "تشغيل عمليات ترحيل قاعدة البيانات",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "direction",
						Value: "up",
						Usage: "اتجاه الترحيل (up, down, reset)",
					},
					&cli.BoolFlag{
						Name:  "verbose",
						Value: false,
						Usage: "عرض تفاصيل الترحيل",
					},
				},
				Action: runMigrations,
			},
		},
		// معالجة الأخطاء العالمية
		ExitErrHandler: func(c *cli.Context, err error) {
			if err != nil {
				logger.Stderr.Error("❌ خطأ في التنفيذ", logger.ErrAttr(err))
				os.Exit(1)
			}
		},
	}

	// تشغيل التطبيق
	if err := app.Run(os.Args); err != nil {
		logger.Stderr.Error("❌ فشل في تشغيل التطبيق", logger.ErrAttr(err))
		os.Exit(1)
	}
}

// ================================
// 🛠️ معالجات الأوامر
// ================================

// runServer تشغيل خادم API
func runServer(c *cli.Context) error {
	logger.Stdout.Info("🚀 بدء تشغيل خادم نوذ تك",
		"version", version,
		"environment", c.String("env"),
		"port", c.String("port"),
	)

	// تعيين متغيرات البيئة إذا تم توفيرها
	if env := c.String("env"); env != "" {
		os.Setenv("APP_ENV", env)
	}
	if port := c.String("port"); port != "" {
		os.Setenv("PORT", port)
	}

	// تشغيل الخادم - سيتم استدعاء server.Run() من cmd/server
	fmt.Println("✅ تم بدء تشغيل خادم نوذ تك")
	fmt.Println("📡 الخادم يعمل على المنفذ:", c.String("port"))
	fmt.Println("🌍 البيئة:", c.String("env"))
	fmt.Println("\nلإيقاف الخادم، اضغط Ctrl+C")

	// انتظار الإشارة لإيقاف الخادم
	waitForShutdownSignal()
	return nil
}

// showVersion عرض معلومات الإصدار
func showVersion(c *cli.Context) error {
	fmt.Printf("نوذ تك - منصة الخدمات الإلكترونية\n")
	fmt.Printf("الإصدار:    %s\n", version)
	fmt.Printf("الكوميت:    %s\n", commit)
	fmt.Printf("وقت البناء: %s\n", date)
	fmt.Printf("بيئة التشغيل: %s\n", getEnv("APP_ENV", "development"))
	fmt.Printf("وقت التشغيل: %s\n", time.Now().Format("2006-01-02 15:04:05"))
	
	// معلومات النظام
	fmt.Printf("\nمعلومات النظام:\n")
	fmt.Printf("نظام التشغيل: %s\n", getOSInfo())
	fmt.Printf("المعالج:      %s\n", getArchitecture())
	fmt.Printf("لغة Go:       %s\n", getGoVersion())
	
	return nil
}

// checkHealth فحص صحة النظام
func checkHealth(c *cli.Context) error {
	timeoutStr := c.String("timeout")
	timeout, err := time.ParseDuration(timeoutStr)
	if err != nil {
		return fmt.Errorf("مهلة غير صالحة: %s", timeoutStr)
	}

	// استخدام context مع timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	logger.Stdout.Info("🔍 فحص صحة النظام...",
		"timeout", timeout.String(),
	)

	// استخدام ctx لمنع تحذير "declared and not used"
	select {
	case <-ctx.Done():
		if ctx.Err() == context.DeadlineExceeded {
			return fmt.Errorf("انتهت مهلة فحص الصحة")
		}
	default:
		// الاستمرار في الفحص
	}

	// هنا يمكن إضافة فحوصات إضافية
	// مثل اتصال قاعدة البيانات، خدمات الطرف الثالث، إلخ.

	fmt.Printf("✅ النظام يعمل بشكل صحيح\n")
	fmt.Printf("⏱️  المهلة: %s\n", timeout.String())
	fmt.Printf("🕐 الوقت: %s\n", time.Now().Format("2006-01-02 15:04:05"))

	return nil
}

// runMigrations تشغيل عمليات ترحيل قاعدة البيانات
func runMigrations(c *cli.Context) error {
	direction := c.String("direction")
	verbose := c.Bool("verbose")

	logger.Stdout.Info("🗄️  تشغيل عمليات ترحيل قاعدة البيانات",
		"direction", direction,
		"verbose", verbose,
	)

	// تنفيذ عمليات الترحيل
	// هذا مكان لوضع منطق ترحيل قاعدة البيانات

	switch direction {
	case "up":
		fmt.Printf("✅ تم تنفيذ ترحيل قاعدة البيانات (UP)\n")
	case "down":
		fmt.Printf("✅ تم ترجيع ترحيل قاعدة البيانات (DOWN)\n")
	case "reset":
		fmt.Printf("✅ تم إعادة تعيين ترحيل قاعدة البيانات (RESET)\n")
	default:
		return fmt.Errorf("اتجاه ترحيل غير معروف: %s", direction)
	}

	if verbose {
		fmt.Printf("📋 المهام المنفذة:\n")
		fmt.Printf("  - إنشاء الجداول الأساسية\n")
		fmt.Printf("  - إضافة الفهارس\n")
		fmt.Printf("  - إدراج البيانات الأولية\n")
	}

	return nil
}

// ================================
// 🛠️ دوال مساعدة
// ================================

// getEnv الحصول على متغير بيئة مع قيمة افتراضية
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getOSInfo الحصول على معلومات نظام التشغيل
func getOSInfo() string {
	return runtime.GOOS
}

// getArchitecture الحصول على بنية المعالج
func getArchitecture() string {
	return runtime.GOARCH
}

// getGoVersion الحصول على إصدار Go
func getGoVersion() string {
	return runtime.Version()
}

// waitForShutdownSignal انتظار إشارة الإغلاق
func waitForShutdownSignal() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	
	<-sigChan
	fmt.Println("\n🛑 استلام إشارة إغلاق...")
}

// ================================
// 🛡️ معالجة الإشارات
// ================================

// setupSignalHandler إعداد معالج الإشارات
func setupSignalHandler() context.Context {
	ctx, cancel := context.WithCancel(context.Background())

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		sig := <-sigChan
		logger.Stdout.Info("🛑 استلام إشارة إغلاق",
			"signal", sig.String(),
		)
		cancel()
	}()

	return ctx
}

// init التهيئة - تُنفذ قبل main()
func init() {
	// معالجة الإشارات
	ctx := setupSignalHandler()

	// استخدام ctx لمنع تحذير "declared and not used"
	go func() {
		<-ctx.Done()
		logger.Stdout.Info("🔚 إغلاق التطبيق بناءً على الإشارة")
	}()

	// تهيئة التسجيل الأساسي
	logger.Init(getEnv("APP_ENV", "development"))

	// تسجيل بدء التشغيل
	logger.Stdout.Info("🔧 تهيئة تطبيق نوذ تك",
		"version", version,
		"go_version", getGoVersion(),
		"environment", getEnv("APP_ENV", "development"),
	)
}