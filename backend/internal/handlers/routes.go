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

	// داخل دالة Register، بعد تعريف المسارات الأخرى
RegisterUserRoutes(router, services.User, services.Admin)
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