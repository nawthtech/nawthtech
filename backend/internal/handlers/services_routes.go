package handlers

import (
	"github.com/gin-gonic/gin"
)

func RegisterServicesRoutes(router *gin.RouterGroup, servicesHandler *ServicesHandler, authMiddleware gin.HandlerFunc, sellerMiddleware gin.HandlerFunc, adminMiddleware gin.HandlerFunc) {
	servicesRoutes := router.Group("/services")
	
	// === الطرق العامة (لا تتطلب مصادقة) ===
	servicesRoutes.GET("", servicesHandler.GetServices)
	servicesRoutes.GET("/search", servicesHandler.SearchServices)
	servicesRoutes.GET("/featured", servicesHandler.GetFeaturedServices)
	servicesRoutes.GET("/categories", servicesHandler.GetAllCategories)
	servicesRoutes.GET("/tags/popular", servicesHandler.GetPopularTags)
	servicesRoutes.GET("/popular", servicesHandler.GetPopularServices)
	servicesRoutes.GET("/category/:category", servicesHandler.GetServicesByCategory)
	servicesRoutes.GET("/tag/:tag", servicesHandler.GetServicesByTag)
	servicesRoutes.GET("/:serviceId", servicesHandler.GetServiceDetails)
	servicesRoutes.GET("/:serviceId/recommended", servicesHandler.GetRecommendedServices)
	servicesRoutes.GET("/:serviceId/similar", servicesHandler.GetSimilarServices)
	servicesRoutes.GET("/seller/:sellerId", servicesHandler.GetSellerServices)
	servicesRoutes.GET("/:serviceId/ratings", servicesHandler.GetServiceRatings)
	
	// === الطرق التي تتطلب مصادقة ===
	authServicesRoutes := servicesRoutes.Group("")
	authServicesRoutes.Use(authMiddleware)
	{
		authServicesRoutes.POST("/:serviceId/check-availability", servicesHandler.CheckAvailability)
		authServicesRoutes.POST("/:serviceId/ratings", servicesHandler.AddRating)
		authServicesRoutes.PUT("/ratings/:ratingId", servicesHandler.UpdateRating)
		authServicesRoutes.DELETE("/ratings/:ratingId", servicesHandler.DeleteRating)
		authServicesRoutes.GET("/my/ratings", servicesHandler.GetMyRatings)
	}
	
	// === طرق البائعين ===
	sellerServicesRoutes := servicesRoutes.Group("")
	sellerServicesRoutes.Use(sellerMiddleware)
	{
		// إدارة الخدمات
		sellerServicesRoutes.POST("", servicesHandler.CreateService)
		sellerServicesRoutes.PUT("/:serviceId", servicesHandler.UpdateService)
		sellerServicesRoutes.PATCH("/:serviceId/status", servicesHandler.UpdateServiceStatus)
		sellerServicesRoutes.DELETE("/:serviceId", servicesHandler.DeleteService)
		
		// الخدمات الشخصية
		sellerServicesRoutes.GET("/my/services", servicesHandler.GetMyServices)
		sellerServicesRoutes.GET("/my/stats", servicesHandler.GetMyServicesStats)
		sellerServicesRoutes.GET("/my/stats/status", servicesHandler.GetServicesStatusCount)
		sellerServicesRoutes.GET("/my/growth", servicesHandler.GetServicesGrowth)
		
		// إدارة الفترات الزمنية
		sellerServicesRoutes.POST("/:serviceId/time-slots", servicesHandler.CreateTimeSlot)
		sellerServicesRoutes.GET("/:serviceId/time-slots", servicesHandler.GetTimeSlots)
		sellerServicesRoutes.PUT("/time-slots/:slotId", servicesHandler.UpdateTimeSlot)
		sellerServicesRoutes.DELETE("/time-slots/:slotId", servicesHandler.DeleteTimeSlot)
		
		// البحث المتقدم للبائعين
		sellerServicesRoutes.POST("/search/advanced", servicesHandler.AdvancedSearch)
	}
	
	// === طرق المسؤولين ===
	adminServicesRoutes := servicesRoutes.Group("/admin")
	adminServicesRoutes.Use(adminMiddleware)
	{
		// إدارة جميع الخدمات
		adminServicesRoutes.GET("", servicesHandler.GetAllServices)
		adminServicesRoutes.PUT("/:serviceId/status", servicesHandler.AdminUpdateServiceStatus)
		adminServicesRoutes.PUT("/:serviceId/featured", servicesHandler.AdminUpdateFeaturedStatus)
		adminServicesRoutes.DELETE("/:serviceId", servicesHandler.AdminDeleteService)
		
		// الإحصائيات والتقارير
		adminServicesRoutes.GET("/stats/overview", servicesHandler.GetAdminServicesStats)
		adminServicesRoutes.GET("/stats/categories", servicesHandler.GetCategoriesStats)
		adminServicesRoutes.GET("/stats/growth", servicesHandler.GetAdminServicesGrowth)
		adminServicesRoutes.GET("/reports/popular", servicesHandler.GetPopularServicesReport)
		
		// إدارة التقييمات
		adminServicesRoutes.DELETE("/ratings/:ratingId", servicesHandler.AdminDeleteRating)
		adminServicesRoutes.GET("/ratings/reported", servicesHandler.GetReportedRatings)
		
		// إدارة الفئات والوسوم
		adminServicesRoutes.POST("/categories", servicesHandler.CreateCategory)
		adminServicesRoutes.PUT("/categories/:categoryId", servicesHandler.UpdateCategory)
		adminServicesRoutes.DELETE("/categories/:categoryId", servicesHandler.DeleteCategory)
		adminServicesRoutes.GET("/tags/management", servicesHandler.ManageTags)
	}
	
	// === طرق مختلطة (تتطلب مصادقة ولكن ليست حصرية للبائعين) ===
	mixedServicesRoutes := servicesRoutes.Group("")
	mixedServicesRoutes.Use(authMiddleware)
	{
		// يمكن للمستخدمين العاديين والبائعين الوصول لهذه المسارات
		mixedServicesRoutes.GET("/user/recommendations", servicesHandler.GetPersonalizedRecommendations)
		mixedServicesRoutes.GET("/user/history", servicesHandler.GetServiceHistory)
		mixedServicesRoutes.POST("/:serviceId/favorite", servicesHandler.AddToFavorites)
		mixedServicesRoutes.DELETE("/:serviceId/favorite", servicesHandler.RemoveFromFavorites)
		mixedServicesRoutes.GET("/user/favorites", servicesHandler.GetFavoriteServices)
	}
}