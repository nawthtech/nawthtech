package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nawthtech/nawthtech/backend/internal/services"
	"github.com/nawthtech/nawthtech/backend/internal/utils"
)

// AuthMiddleware وسيط المصادقة
func AuthMiddleware(authService services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// استخراج التوكن من الرأس
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.ErrorResponse(c, http.StatusUnauthorized, "مطلوب توكن مصادقة", "AUTH_TOKEN_REQUIRED")
			c.Abort()
			return
		}

		// التحقق من صيغة التوكن (Bearer token)
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.ErrorResponse(c, http.StatusUnauthorized, "صيغة التوكن غير صالحة", "INVALID_TOKEN_FORMAT")
			c.Abort()
			return
		}

		token := parts[1]

		// التحقق من صحة التوكن
		claims, err := authService.VerifyToken(token)
		if err != nil {
			utils.ErrorResponse(c, http.StatusUnauthorized, "توكن غير صالح", "INVALID_TOKEN")
			c.Abort()
			return
		}

		// تخزين بيانات المستخدم في السياق
		c.Set("userID", claims.UserID)
		c.Set("userRole", claims.Role)
		c.Set("userEmail", claims.Email)

		c.Next()
	}
}

// AdminMiddleware وسيط التحقق من صلاحيات المسؤول
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("userRole")
		if !exists || userRole != "admin" {
			utils.ErrorResponse(c, http.StatusForbidden, "غير مصرح - صلاحيات مسؤول مطلوبة", "ADMIN_ACCESS_REQUIRED")
			c.Abort()
			return
		}
		c.Next()
	}
}

// SellerMiddleware وسيط التحقق من صلاحيات البائع
func SellerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("userRole")
		if !exists || (userRole != "seller" && userRole != "admin") {
			utils.ErrorResponse(c, http.StatusForbidden, "غير مصرح - صلاحيات بائع مطلوبة", "SELLER_ACCESS_REQUIRED")
			c.Abort()
			return
		}
		c.Next()
	}
}