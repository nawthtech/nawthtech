package handlers

import (
	"net/http"
	"strconv"

	"github.com/nawthtech/nawthtech/backend/internal/middleware"
	"github.com/nawthtech/nawthtech/backend/internal/services"

	"github.com/go-chi/chi/v5"
)

// RegisterUserRoutes تسجيل جميع مسارات المستخدمين
func RegisterUserRoutes(router chi.Router, userService *services.UserService, adminService *services.AdminService) {
	userHandler := NewUserHandler(userService)
	adminHandler := NewAdminHandler(adminService)

	router.Route("/api/v1/users", func(r chi.Router) {
		// تطبيق rate limiting على جميع مسارات المستخدمين
		r.Use(middleware.RateLimiter(100, 15*60*1000)) // 100 طلب كل 15 دقيقة

		// ==================== الملف الشخصي والإعدادات ====================
		r.Route("/profile", func(profile chi.Router) {
			profile.Use(middleware.AuthMiddleware)

			// الحصول على بيانات الملف الشخصي للمستخدم الحالي
			profile.Get("/", userHandler.GetUserProfile)

			// تحديث بيانات الملف الشخصي للمستخدم الحالي
			profile.Put("/", userHandler.UpdateUserProfile)

			// تغيير كلمة مرور المستخدم الحالي
			profile.Put("/password", userHandler.ChangePassword)

			// تحديث صورة الملف الشخصي
			profile.Put("/avatar", userHandler.UpdateAvatar)
		})

		// ==================== إعدادات المستخدم ====================
		r.Route("/settings", func(settings chi.Router) {
			settings.Use(middleware.AuthMiddleware)

			// الحصول على إعدادات المستخدم الحالي
			settings.Get("/", userHandler.GetUserSettings)

			// تحديث إعدادات المستخدم الحالي
			settings.Put("/", userHandler.UpdateUserSettings)
		})

		// ==================== الطلبات والمشتريات ====================
		r.Route("/orders", func(orders chi.Router) {
			orders.Use(middleware.AuthMiddleware)

			// الحصول على طلبات المستخدم الحالي
			orders.Get("/", userHandler.GetUserOrders)

			// الحصول على تفاصيل طلب محدد للمستخدم الحالي
			orders.Get("/{orderId}", userHandler.GetUserOrderDetails)
		})

		// ==================== السلة والمشتريات ====================
		r.Route("/cart", func(cart chi.Router) {
			cart.Use(middleware.AuthMiddleware)

			// الحصول على سلة المستخدم الحالي
			cart.Get("/", userHandler.GetUserCart)
		})

		r.Route("/wishlist", func(wishlist chi.Router) {
			wishlist.Use(middleware.AuthMiddleware)

			// الحصول على قائمة رغبات المستخدم
			wishlist.Get("/", userHandler.GetUserWishlist)
		})

		// ==================== الإحصائيات والنشاط ====================
		r.Route("/stats", func(stats chi.Router) {
			stats.Use(middleware.AuthMiddleware)

			// الحصول على إحصائيات المستخدم الحالي
			stats.Get("/", userHandler.GetUserStats)
		})

		r.Route("/activity", func(activity chi.Router) {
			activity.Use(middleware.AuthMiddleware)

			// الحصول على نشاط المستخدم الحالي
			activity.Get("/", userHandler.GetUserActivity)
		})

		// ==================== الإشعارات ====================
		r.Route("/notifications", func(notifications chi.Router) {
			notifications.Use(middleware.AuthMiddleware)

			// الحصول على إشعارات المستخدم الحالي
			notifications.Get("/", userHandler.GetUserNotifications)

			// تعليم الإشعارات كمقروءة
			notifications.Put("/read", userHandler.MarkNotificationsAsRead)

			// حذف إشعارات المستخدم
			notifications.Delete("/", userHandler.DeleteNotifications)
		})

		// ==================== مسارات البائعين ====================
		r.Route("/seller", func(seller chi.Router) {
			seller.Use(middleware.AuthMiddleware, middleware.AdminAuth)

			// الحصول على خدمات البائع الحالي
			seller.Get("/services", userHandler.GetSellerServices)

			// الحصول على إحصائيات البائع الحالي
			seller.Get("/stats", userHandler.GetSellerStats)

			// الحصول على طلبات البائع الحالي
			seller.Get("/orders", userHandler.GetSellerOrders)
		})

		// ==================== مسارات الإدارة ====================
		r.Route("/admin", func(admin chi.Router) {
			admin.Use(middleware.AuthMiddleware, middleware.AdminAuth)

			// الحصول على جميع المستخدمين
			admin.Get("/users", adminHandler.GetAllUsers)

			// الحصول على مستخدم بواسطة المعرف
			admin.Get("/users/{userId}", adminHandler.GetUserById)

			// تحديث بيانات مستخدم
			admin.Put("/users/{userId}", adminHandler.UpdateUser)

			// حذف مستخدم
			admin.Delete("/users/{userId}", adminHandler.DeleteUser)

			// تحديث دور مستخدم
			admin.Put("/users/{userId}/role", adminHandler.UpdateUserRole)

			// تعطيل حساب مستخدم
			admin.Post("/users/{userId}/deactivate", adminHandler.DeactivateUser)

			// تفعيل حساب مستخدم معطل
			admin.Post("/users/{userId}/activate", adminHandler.ActivateUser)
		})

		// ==================== البحث والتصفية ====================
		r.Route("/search", func(search chi.Router) {
			search.Use(middleware.AuthMiddleware, middleware.AdminAuth)

			// البحث عن المستخدمين
			search.Get("/", adminHandler.SearchUsers)
		})
	})
}