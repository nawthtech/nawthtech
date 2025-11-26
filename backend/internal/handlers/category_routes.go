package handlers

import (
	"github.com/nawthtech/nawthtech/backend/internal/middleware"
	"github.com/nawthtech/nawthtech/backend/internal/services"

	"github.com/go-chi/chi/v5"
)

// RegisterCategoryRoutes تسجيل جميع مسارات الفئات
func RegisterCategoryRoutes(router chi.Router, categoryService *services.CategoryService) {
	categoryHandler := NewCategoryHandler(categoryService)

	router.Route("/api/v1/categories", func(r chi.Router) {
		// ==================== الفئات العامة ====================
		// الحصول على جميع الفئات النشطة
		r.Get("/", categoryHandler.GetCategories)
		
		// الحصول على هيكل شجرة الفئات
		r.Get("/tree", categoryHandler.GetCategoryTree)
		
		// الحصول على تفاصيل فئة محددة
		r.Get("/{categoryId}", categoryHandler.GetCategoryById)
		
		// الحصول على خدمات فئة محددة
		r.Get("/{categoryId}/services", categoryHandler.GetCategoryServices)
		
		// الحصول على الفئات الفرعية
		r.Get("/{categoryId}/subcategories", categoryHandler.GetSubcategories)

		// ==================== إحصائيات الفئات ====================
		// الحصول على إحصائيات الفئات الشاملة
		r.Get("/stats/overview", categoryHandler.GetCategoriesStats)
		
		// الحصول على إحصائيات فئة محددة
		r.Get("/{categoryId}/stats", categoryHandler.GetCategoryStats)

		// ==================== البحث والتصفية ====================
		// البحث في الفئات
		r.Get("/search", categoryHandler.SearchCategories)

		// ==================== إدارة الفئات (للمسؤولين) ====================
		r.Route("/admin", func(admin chi.Router) {
			admin.Use(middleware.AuthMiddleware, middleware.AdminAuth)

			// إنشاء فئة جديدة
			admin.Post("/", categoryHandler.CreateCategory)
			
			// تحديث فئة موجودة
			admin.Put("/{categoryId}", categoryHandler.UpdateCategory)
			
			// حذف فئة
			admin.Delete("/{categoryId}", categoryHandler.DeleteCategory)
			
			// تحديث حالة الفئة
			admin.Patch("/{categoryId}/status", categoryHandler.UpdateCategoryStatus)
			
			// تحديث ترتيب الفئة
			admin.Patch("/{categoryId}/order", categoryHandler.UpdateCategoryOrder)

			// ==================== الاستيراد والتصدير ====================
			// استيراد فئات من ملف
			admin.Post("/import", categoryHandler.ImportCategories)
			
			// تصدير الفئات إلى ملف
			admin.Get("/export", categoryHandler.ExportCategories)
		})
	})
}