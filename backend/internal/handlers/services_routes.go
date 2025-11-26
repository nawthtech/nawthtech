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
	servicesRoutes.GET("/:serviceId", servicesHandler.GetServiceDetails)
	servicesRoutes.GET("/:serviceId/recommended", servicesHandler.GetRecommendedServices)
	servicesRoutes.GET("/seller/:sellerId", servicesHandler.GetSellerServices)
	
	// === الطرق التي تتطلب مصادقة ===
	authServicesRoutes := servicesRoutes.Group("")
	authServicesRoutes.Use(authMiddleware)
	authServicesRoutes.POST("/:serviceId/check-availability", servicesHandler.CheckAvailability)
	authServicesRoutes.POST("/:serviceId/rating", servicesHandler.AddRating)
	
	// === طرق البائعين والمسؤولين ===
	sellerServicesRoutes := servicesRoutes.Group("")
	sellerServicesRoutes.Use(sellerMiddleware)
	sellerServicesRoutes.POST("", servicesHandler.CreateService)
	sellerServicesRoutes.PUT("/:serviceId", servicesHandler.UpdateService)
	sellerServicesRoutes.PATCH("/:serviceId/status", servicesHandler.UpdateServiceStatus)
	sellerServicesRoutes.DELETE("/:serviceId", servicesHandler.DeleteService)
	sellerServicesRoutes.GET("/stats/all", servicesHandler.GetServicesStats)
}