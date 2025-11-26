package handlers

import (
	"github.com/nawthtech/nawthtech/backend/internal/middleware"
	"github.com/nawthtech/nawthtech/backend/internal/services"

	"github.com/go-chi/chi/v5"
)

// RegisterStoreRoutes تسجيل جميع مسارات المتجر
func RegisterStoreRoutes(router chi.Router, storeService *services.StoreService, cartService *services.CartService) {
	storeHandler := NewStoreHandler(storeService)
	cartHandler := NewCartHandler(cartService)

	router.Route("/api/v1/store", func(r chi.Router) {
		// ==================== خدمات المتجر الأساسية ====================
		r.Route("/services", func(services chi.Router) {
			// الحصول على الخدمات مع التوصيات الذكية
			services.With(middleware.RateLimiter(100, 15*60*1000)).Get("/", storeHandler.GetServices)
			
			// الحصول على تفاصيل خدمة محددة
			services.Get("/{serviceId}", storeHandler.GetServiceDetails)
			
			// التحقق من توفر الخدمة
			services.Post("/{serviceId}/check-availability", storeHandler.CheckServiceAvailability)
		})

		r.Route("/categories", func(categories chi.Router) {
			// الحصول على التصنيفات مع الإحصائيات
			categories.Get("/", storeHandler.GetCategoriesWithStats)
			
			// الحصول على خدمات تصنيف محدد
			categories.Get("/{categoryId}/services", storeHandler.GetServicesByCategory)
		})

		// ==================== إدارة السلة ====================
		r.Route("/cart", func(cart chi.Router) {
			cart.Use(middleware.AuthMiddleware)

			// الحصول على سلة المستخدم
			cart.Get("/", cartHandler.GetCart)
			
			// إضافة عنصر إلى السلة
			cart.Post("/items", cartHandler.AddToCart)
			
			// تحديث كمية عنصر في السلة
			cart.Put("/items/{itemId}", cartHandler.UpdateCartItem)
			
			// إزالة عنصر من السلة
			cart.Delete("/items/{itemId}", cartHandler.RemoveFromCart)
			
			// تفريغ السلة
			cart.Delete("/", cartHandler.ClearCart)
			
			// الحصول على ملخص السلة
			cart.Get("/summary", cartHandler.GetCartSummary)
			
			// التحقق من صحة عناصر السلة
			cart.Post("/validate", cartHandler.ValidateCartItems)
		})

		// ==================== إدارة الطلبات ====================
		r.Route("/orders", func(orders chi.Router) {
			orders.Use(middleware.AuthMiddleware)

			// إنشاء طلب جديد مع التحقق بالذكاء الاصطناعي
			orders.With(middleware.RateLimiter(50, 60*60*1000)).Post("/", storeHandler.CreateAIOrder)
			
			// إنشاء طلب من السلة
			orders.Post("/from-cart", storeHandler.CreateOrderFromCart)
			
			// الحصول على قائمة طلبات المستخدم
			orders.Get("/", storeHandler.GetUserOrders)
			
			// الحصول على تفاصيل طلب محدد
			orders.Get("/{orderId}", storeHandler.GetOrderDetails)
			
			// تحديث حالة الطلب
			orders.Patch("/{orderId}/status", storeHandler.UpdateOrderStatus)
			
			// إلغاء الطلب
			orders.Post("/{orderId}/cancel", storeHandler.CancelOrder)
		})

		// ==================== التوصيات والتحليلات ====================
		r.Route("/recommendations", func(recommendations chi.Router) {
			recommendations.Use(middleware.AuthMiddleware)

			// الحصول على التوصيات الذكية المخصصة
			recommendations.With(middleware.RateLimiter(30, 10*60*1000)).Get("/", storeHandler.GetAIRecommendations)
		})

		r.Route("/stats", func(stats chi.Router) {
			stats.Use(middleware.AuthMiddleware)

			// الحصول على إحصائيات المتجر الشاملة
			stats.With(middleware.RateLimiter(60, 5*60*1000)).Get("/", storeHandler.GetStoreStats)
		})

		// ==================== إدارة الفئات والبحث ====================
		r.Route("/search", func(search chi.Router) {
			// البحث المتقدم في الخدمات
			search.With(middleware.RateLimiter(50, 5*60*1000)).Get("/", storeHandler.SearchServices)
		})

		// ==================== التقارير ====================
		r.Route("/orders/{orderId}/report", func(report chi.Router) {
			report.Use(middleware.AuthMiddleware)

			// الحصول على تقرير النمو الذكي للطلب
			report.Get("/", storeHandler.GenerateGrowthReport)
		})
	})
}