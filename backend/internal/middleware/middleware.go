package middleware

import (
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nawthtech/nawthtech/backend/internal/config"
	"github.com/nawthtech/nawthtech/backend/internal/logger"
	"github.com/nawthtech/nawthtech/backend/internal/services"
)

// ==================== هياكل البيانات ====================

// MiddlewareContainer حاوية الوسائط
type MiddlewareContainer struct {
	AuthMiddleware      gin.HandlerFunc
	AdminMiddleware     gin.HandlerFunc
	SellerMiddleware    gin.HandlerFunc
	CORSMiddleware      gin.HandlerFunc
	SecurityMiddleware  gin.HandlerFunc
	RateLimitMiddleware gin.HandlerFunc
	LoggerMiddleware    gin.HandlerFunc
}

// ==================== الوسائط الأساسية ====================

// CORS وسيط CORS
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length, X-Request-ID, X-RateLimit-Limit, X-RateLimit-Remaining, X-RateLimit-Reset")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// SecurityHeaders وسيط رؤوس الأمان
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")
		c.Writer.Header().Set("X-Frame-Options", "DENY")
		c.Writer.Header().Set("X-XSS-Protection", "1; mode=block")
		c.Writer.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		c.Writer.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Writer.Header().Set("Content-Security-Policy", "default-src 'self'")
		
		c.Next()
	}
}

// RequestID وسيط إضافة معرف الطلب
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID()
		}

		// إضافة معرف الطلب إلى الرأس
		c.Writer.Header().Set("X-Request-ID", requestID)

		// إضافة معرف الطلب إلى السياق
		c.Set("requestID", requestID)

		c.Next()
	}
}

// Logging وسيط التسجيل
func Logging() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		
		// معالجة الطلب
		c.Next()

		// حساب مدة التنفيذ
		duration := time.Since(start)

		// جمع معلومات التسجيل
		fields := []interface{}{
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"duration", duration.String(),
			"bytes", c.Writer.Size(),
			"user_agent", c.Request.UserAgent(),
			"ip", getClientIP(c.Request),
		}

		// إضافة معرف الطلب إذا كان متوفراً
		if requestID, exists := c.Get("requestID"); exists {
			fields = append(fields, "request_id", requestID)
		}

		// تسجيل بناءً على حالة الاستجابة
		status := c.Writer.Status()
		switch {
		case status >= 500:
			logger.Stderr.Error("خطأ في الخادم", fields...)
		case status >= 400:
			logger.Stdout.Warn("خطأ في العميل", fields...)
		default:
			logger.Stdout.Info("طلب معالَج", fields...)
		}
	}
}

// ==================== وسائط المصادقة ====================

// AuthMiddleware وسيط المصادقة باستخدام AuthService
func AuthMiddleware(authService services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// الحصول على التوكن من الرأس
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "مصادقة مطلوبة",
				"message": "يرجى تقديم رمز المصادقة",
			})
			c.Abort()
			return
		}

		// التحقق من صيغة التوكن (Bearer token)
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "صيغة توكن غير صالحة",
				"message": "يجب أن يكون التوكن بصيغة Bearer token",
			})
			c.Abort()
			return
		}

		token := parts[1]

		// التحقق من صحة التوكن باستخدام AuthService
		claims, err := authService.VerifyToken(c.Request.Context(), token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "توكن غير صالح",
				"message": "رمز المصادقة منتهي الصلاحية أو غير صحيح",
			})
			c.Abort()
			return
		}

		// إضافة معلومات المستخدم إلى السياق
		c.Set("userID", claims.UserID)
		c.Set("userEmail", claims.Email)
		c.Set("userRole", claims.Role)
		c.Set("token", token)

		c.Next()
	}
}

// AdminMiddleware وسيط مصادقة المسؤول
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// الحصول على معلومات المستخدم من السياق
		userRole, exists := c.Get("userRole")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"error":   "صلاحيات غير كافية",
				"message": "مطلوب صلاحيات مسؤول للوصول إلى هذا المورد",
			})
			c.Abort()
			return
		}

		// التحقق من صلاحيات المسؤول
		if userRole != "admin" {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"error":   "صلاحيات إدارة مطلوبة",
				"message": "لا تملك الصلاحيات الكافية للوصول إلى هذا المورد",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// SellerMiddleware وسيط مصادقة البائعين
func SellerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// الحصول على معلومات المستخدم من السياق
		userRole, exists := c.Get("userRole")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"error":   "صلاحيات غير كافية",
				"message": "مطلوب صلاحيات بائع للوصول إلى هذا المورد",
			})
			c.Abort()
			return
		}

		// التحقق من صلاحيات البائع أو المسؤول
		if userRole != "seller" && userRole != "admin" {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"error":   "صلاحيات بائع مطلوبة",
				"message": "لا تملك الصلاحيات الكافية للوصول إلى هذا المورد",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// ==================== وسائط الأمان والأداء ====================

// RateLimit وسيط تحديد المعدل
func RateLimit() gin.HandlerFunc {
	// تنفيذ مبسط لتحديد المعدل
	type clientLimit struct {
		count    int
		lastSeen time.Time
	}

	clients := make(map[string]*clientLimit)
	requests := 100 // 100 طلب
	window := time.Minute // في الدقيقة
	
	return func(c *gin.Context) {
		clientIP := getClientIP(c.Request)
		
		now := time.Now()
		
		// تنظيف العملاء القدامى
		if len(clients) > 10000 {
			for ip, limit := range clients {
				if now.Sub(limit.lastSeen) > window {
					delete(clients, ip)
				}
			}
		}
		
		limit, exists := clients[clientIP]
		if !exists {
			limit = &clientLimit{count: 0, lastSeen: now}
			clients[clientIP] = limit
		}
		
		// إعادة تعيين العداد إذا انتهت النافذة الزمنية
		if now.Sub(limit.lastSeen) > window {
			limit.count = 0
			limit.lastSeen = now
		}
		
		// التحقق من تجاوز الحد
		if limit.count >= requests {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"success": false,
				"error":   "تم تجاوز معدل الطلبات المسموح به",
				"message": "يرجى المحاولة مرة أخرى لاحقاً",
			})
			c.Abort()
			return
		}
		
		// زيادة العداد
		limit.count++
		limit.lastSeen = now
		
		// إضافة معلومات التحديد إلى الرأس
		c.Writer.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", requests))
		c.Writer.Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", requests-limit.count))
		c.Writer.Header().Set("X-RateLimit-Reset", fmt.Sprintf("%d", now.Add(window).Unix()))
		
		c.Next()
	}
}

// ValidateContentType وسيط التحقق من نوع المحتوى
func ValidateContentType() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "PATCH" {
			contentType := c.GetHeader("Content-Type")
			if !strings.Contains(contentType, "application/json") {
				c.JSON(http.StatusBadRequest, gin.H{
					"success": false,
					"error":   "نوع المحتوى يجب أن يكون JSON",
					"message": "يرجى استخدام application/json كنوع للمحتوى",
				})
				c.Abort()
				return
			}
		}
		
		c.Next()
	}
}

// SizeLimit وسيط تحديد حجم الطلب
func SizeLimit(maxSize int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.ContentLength > maxSize {
			c.JSON(http.StatusRequestEntityTooLarge, gin.H{
				"success": false,
				"error":   "حجم البيانات كبير جداً",
				"message": fmt.Sprintf("الحجم الأقصى المسموح به هو %d بايت", maxSize),
			})
			c.Abort()
			return
		}
		
		c.Next()
	}
}

// Timeout وسيط وقت انتهاء الطلب
func Timeout(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// إعداد مهلة للطلب
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()
		
		c.Request = c.Request.WithContext(ctx)
		
		// قناة للإشارة بانتهاء المعالجة
		done := make(chan bool, 1)
		
		go func() {
			c.Next()
			done <- true
		}()
		
		select {
		case <-ctx.Done():
			if ctx.Err() == context.DeadlineExceeded {
				c.JSON(http.StatusRequestTimeout, gin.H{
					"success": false,
					"error":   "انتهت مهلة الطلب",
					"message": "تجاوز الطلب الوقت المحدد للمعالجة",
				})
				c.Abort()
			}
		case <-done:
			// الطلب اكتمل بنجاح
		}
	}
}

// Recovery وسيط استعادة الأخطاء
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// الحصول على معرف الطلب
				requestID, _ := c.Get("requestID")
				
				// تسجيل الخطأ
				logger.Stderr.Error("تعافى من حالة panic",
					"error", err,
					"request_id", requestID,
					"path", c.Request.URL.Path,
					"method", c.Request.Method,
				)
				
				// إرسال استجابة خطأ
				c.JSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"error":   "خطأ داخلي في الخادم",
					"message": "حدث خطأ غير متوقع",
				})
				
				c.Abort()
			}
		}()
		
		c.Next()
	}
}

// ==================== وسائط التحقق من الصلاحيات ====================

// ValidateAdminAction وسيط التحقق من إجراءات المسؤول
func ValidateAdminAction() gin.HandlerFunc {
	return func(c *gin.Context) {
		// التحقق من أن المستخدم مسؤول
		userRole, exists := c.Get("userRole")
		if !exists || userRole != "admin" {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"error":   "صلاحيات غير كافية",
				"message": "مطلوب صلاحيات مسؤول لهذا الإجراء",
			})
			c.Abort()
			return
		}
		
		// التحقق من نوع المحتوى للطلبات التي تحتوي على جسم
		if c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "PATCH" {
			contentType := c.GetHeader("Content-Type")
			if !strings.Contains(contentType, "application/json") {
				c.JSON(http.StatusBadRequest, gin.H{
					"success": false,
					"error":   "نوع المحتوى يجب أن يكون JSON",
					"message": "يرجى استخدام application/json لنوع المحتوى",
				})
				c.Abort()
				return
			}
		}
		
		c.Next()
	}
}

// OwnerOrAdmin وسيط التحقق من المالك أو المسؤول
func OwnerOrAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "غير مصرح",
				"message": "يجب تسجيل الدخول للوصول إلى هذا المورد",
			})
			c.Abort()
			return
		}
		
		userRole, _ := c.Get("userRole")
		
		// إذا كان المستخدم مسؤولاً، اسمح بالوصول
		if userRole == "admin" {
			c.Next()
			return
		}
		
		// الحصول على معرف المورد من المسار
		resourceUserID := c.Param("userID")
		if resourceUserID == "" {
			resourceUserID = c.Param("id")
		}
		
		// التحقق إذا كان المستخدم هو مالك المورد
		if resourceUserID != userID {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"error":   "صلاحيات غير كافية",
				"message": "لا يمكنك الوصول إلى هذا المورد",
			})
			c.Abort()
			return
		}
		
		c.Next()
	}
}

// ==================== الدوال المساعدة ====================

// getClientIP الحصول على عنوان IP العميل
func getClientIP(r *http.Request) string {
	// التحقق من الرؤوس أولاً
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		ips := strings.Split(ip, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}
	
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}
	
	// استخدام العنوان المباشر
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	
	return host
}

// generateRequestID إنشاء معرف طلب فريد
func generateRequestID() string {
	return fmt.Sprintf("req_%d_%s", time.Now().UnixNano(), randomString(8))
}

// randomString إنشاء سلسلة عشوائية
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}

// ==================== تسجيل الوسائط ====================

// RegisterGlobalMiddlewares تسجيل الوسائط العامة
func RegisterGlobalMiddlewares(router *gin.Engine, cfg *config.Config) {
	// الوسائط الأساسية
	router.Use(Recovery())
	router.Use(RequestID())
	router.Use(Logging())
	router.Use(CORS())
	router.Use(SecurityHeaders())
	router.Use(RateLimit())
	
	// وسائط إضافية بناءً على البيئة
	if cfg.Environment == "production" {
		router.Use(SizeLimit(10 * 1024 * 1024)) // 10MB في الإنتاج
	} else {
		router.Use(SizeLimit(50 * 1024 * 1024)) // 50MB في التطوير
	}
}

// GetUserIDFromContext الحصول على معرف المستخدم من السياق
func GetUserIDFromContext(c *gin.Context) (string, bool) {
	userID, exists := c.Get("userID")
	if !exists {
		return "", false
	}
	return userID.(string), true
}

// GetUserRoleFromContext الحصول على دور المستخدم من السياق
func GetUserRoleFromContext(c *gin.Context) (string, bool) {
	userRole, exists := c.Get("userRole")
	if !exists {
		return "", false
	}
	return userRole.(string), true
}

// GetRequestIDFromContext الحصول على معرف الطلب من السياق
func GetRequestIDFromContext(c *gin.Context) (string, bool) {
	requestID, exists := c.Get("requestID")
	if !exists {
		return "", false
	}
	return requestID.(string), true
}