package handlers

import (
	"net/http"
	"time"

	"github.com/nawthtech/nawthtech/backend/internal/middleware"
	"github.com/nawthtech/nawthtech/backend/internal/services"

	"github.com/go-chi/chi/v5"
)

// Services يحتوي على جميع الخدمات
type Services struct {
	Admin   *services.AdminService
	User    *services.UserService
	Auth    *services.AuthService
	Store   *services.StoreService
	Cart    *services.CartService
	Payment *services.PaymentService
	AI      *services.AIService
	Email   *services.EmailService
	Upload  *services.UploadService
}

// Register تسجيل جميع المسارات
func Register(router chi.Router, services *Services) {
	// إنشاء المعالجات الأساسية فقط
	adminHandler := NewAdminHandler(services.Admin)
	healthHandler := NewHealthHandler()
RegisterPaymentRoutes(router, services.Payment)
RegisterUserRoutes(router, services.User, services.Admin)
RegisterStoreRoutes(router, services.Store, services.Cart)
RegisterCategoryRoutes(router, services.Category)
RegisterAIRoutes(router, services.AI)
NewCacheService(config CacheConfig)
CacheService {
logger :=
slog.New(slog.NewJSONHandler(os.Stdout,
&slog HandlerOptions{
Level: slog.LevelInfo,
RegisterAuthRoutes(router, services.Auth)
RegisterAdminRoutes(router, services.Admin)
websiteService := services.NewWebsiteService()
websiteHandler := handlers.NewWebsiteHandler(websiteService, authService)
cacheHandler := NewCacheHandler(cacheService)
RegisterCacheRoutes(router, cacheHandler, authMiddleware, adminMiddleware)
handlers.RegisterWebsiteRoutes(api, websiteHandler, authMiddleware, adminMiddleware)
analyticsService := services.NewAnalyticsService()
analyticsHandler := handlers.NewAnalyticsHandler(analyticsService, authService)
handlers.RegisterAnalyticsRoutes(api, analyticsHandler, authMiddleware, adminMiddleware)
reportsService := services.NewReportsService()
RegisterSSERoutes (router *gin. RouterGroup, authMiddleware gin. HandlerFunc, adminMiddleware
gin. HandlerFune) {sseRoutes := router. Group("/sse")
servicesRepo := services.NewServicesRepository(db)
servicesService := services.NewServicesService(servicesRepo)
RegisterHealthRoutes(router, db, healthService, config.Version, config.Environment, adminMiddleware)
handlers.RegisterServicesRoutes(apiGroup, servicesService)
reportsHandler := handlers.NewReportsHandler(reportsService, authService)
handlers.RegisterReportsRoutes(api, reportsHandler, authMiddleware, adminMiddleware)
contentService := services.NewContentService()
contentHandler := handlers.NewContentHandler(contentService, authService)
handlers.RegisterContentRoutes(api, contentHandler, authMiddleware)
notificationService := services.NewNotificationService()
ordersService := services.NewOrdersService()
strategiesService := services.NewStrategiesService()
strategiesHandler := handlers.NewStrategiesHandler(strategiesService, authService)
handlers.RegisterStrategiesRoutes(api, strategiesHandler, authMiddleware, adminMiddleware)
uploadsService := services.NewUploadsService()
uploadsHandler := handlers.NewUploadsHandler(uploadsService, authService)
handlers.RegisterUploadsRoutes(api, uploadsHandler, authMiddleware, adminMiddleware, sellerMiddleware)
ordersHandler := handlers.NewOrdersHandler(ordersService, authService)
handlers.RegisterOrdersRoutes(api, ordersHandler, authMiddleware, adminMiddleware, sellerMiddleware)
notificationHandler := handlers.NewNotificationHandler(notificationService, authService)
handlers.RegisterNotificationRoutes(api, notificationHandler, authMiddleware, adminMiddleware)

	// المسارات العامة الأساسية
	router.Route("/api", func(r chi.Router) {
		// الصحة
		r.Get("/health", healthHandler.HealthCheck)
		
		// مسارات الإدارة الأساسية
		r.Route("/admin", func(admin chi.Router) {
			admin.Get("/dashboard", adminHandler.GetDashboardData)
			admin.Get("/system/health", adminHandler.GetSystemHealth)
		})

		// TODO: إضافة المسارات الأخرى عند إنشاء الخدمات
		r.Get("/status", healthHandler.HealthCheck)
	})
}

// تعريفات المعالجات الأساسية
type AdminHandler struct {
	adminService *services.AdminService
}

type HealthHandler struct{}

func (h *HealthHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "healthy", "service": "nawthtech-backend"}`))
}

func NewAdminHandler(adminService *services.AdminService) *AdminHandler {
	return &AdminHandler{adminService: adminService}
}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// دوال placeholder للواجهات
func (h *AdminHandler) GetDashboardData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Dashboard data - under development"}`))
}

func (h *AdminHandler) GetSystemHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"system": "healthy", "timestamp": "` + time.Now().Format(time.RFC3339) + `"}`))
}