package middleware

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nawthtech/nawthtech/backend/internal/config"
	"github.com/nawthtech/nawthtech/backend/internal/logger"
	"github.com/nawthtech/nawthtech/backend/internal/services"
	"github.com/nawthtech/nawthtech/backend/internal/utils"
)

// ==================== هياكل البيانات ====================

// MiddlewareContainer حاوية الوسائط
type MiddlewareContainer struct {
	AuthMiddleware      *AuthMiddlewareStruct
	AdminMiddleware     *AdminMiddlewareStruct
	CORSMiddleware      gin.HandlerFunc
	SecurityMiddleware  gin.HandlerFunc
	RateLimitMiddleware gin.HandlerFunc
}

// AuthMiddlewareStruct واجهة لمصادقة المستخدم
type AuthMiddlewareStruct struct {
	authService services.AuthService
}

// AdminMiddlewareStruct واجهة لمصادقة المسؤول
type AdminMiddlewareStruct struct{}

// SellerMiddleware واجهة لمصادقة البائعين
type SellerMiddleware struct{}

// RateLimiter بنية محدد المعدل
type rateLimiter struct {
	visits map[string][]time.Time
	mu     sync.RWMutex
}

// ==================== المتغيرات العامة ====================

var (
	limiter = &rateLimiter{
		visits: make(map[string][]time.Time),
	}
)

// ==================== الوسائط الأساسية ====================

// CORSMiddleware وسيط CORS
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// في الإنتاج، يمكن تحديد النطاقات المسموحة بدقة
		allowedOrigin := "*"
		if origin != "" {
			allowedOrigin = origin
		}

		c.Writer.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-Request-ID")
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
		c.Writer.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

		// CSP مبسطة
		if c.Request.URL.Path == "/" || strings.HasPrefix(c.Request.URL.Path, "/api/") {
			c.Writer.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'")
		}

		c.Next()
	}
}

// SecurityMiddlewareFunc middleware لإضافة رؤوس الأمان
func SecurityMiddlewareFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		// إضافة رؤوس الأمان
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
		
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

		c.Writer.Header().Set("X-Request-ID", requestID)
		c.Set("requestID", requestID)

		c.Next()
	}
}

// Logging وسيط التسجيل
func Logging() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		requestID, _ := c.Get("requestID")

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
			"ip", getClientIP(c.Request),
			"request_id", requestID,
		}

		// إضافة معرف المستخدم إذا كان متوفراً
		if userID, exists := c.Get("userID"); exists {
			fields = append(fields, "user_id", userID)
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

// Logger وظيفة logger بسيطة
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start)
		
		logger.Stdout.Info("HTTP Request",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"duration", duration.String(),
			"ip", getClientIP(c.Request),
		)
	}
}

// ==================== وسائط المصادقة ====================

// NewAuthMiddleware إنشاء وسيط مصادقة جديد
func NewAuthMiddleware(authService services.AuthService) *AuthMiddlewareStruct {
	return &AuthMiddlewareStruct{authService: authService}
}

// Handle معالجة طلب المصادقة
func (m *AuthMiddlewareStruct) Handle() gin.HandlerFunc {
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
		claims, err := m.authService.VerifyToken(c.Request.Context(), token)
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

// AuthMiddleware للاستخدام مع router
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// هذا تنفيذ مبسط - سيتم استبداله بالتنفيذ الفعلي مع AuthService
		userID := utils.GetUserIDFromGinContext(c)
		if userID == "" {
			// التحقق من رأس المصادقة
			authHeader := c.GetHeader("Authorization")
			if authHeader == "" {
				c.JSON(http.StatusUnauthorized, gin.H{
					"success": false,
					"error":   "مصادقة مطلوبة",
					"message": "يرجى تسجيل الدخول",
				})
				c.Abort()
				return
			}

			// محاولة استخراج معرف المستخدم من التوكن (مبسطة)
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && parts[0] == "Bearer" {
				// في تطبيق حقيقي، يتم فك تشفير التوكن والتحقق منه
				// هنا نستخدم تنفيذ مبسط
				c.Set("userID", "temp_user_id")
				c.Set("userRole", "user")
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{
					"success": false,
					"error":   "مصادقة مطلوبة",
					"message": "يرجى تسجيل الدخول",
				})
				c.Abort()
				return
			}
		}
		c.Next()
	}
}

// OptionalAuth وسيط مصادقة اختياري
func OptionalAuth(authService services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Next()
			return
		}

		token := parts[1]
		claims, err := authService.VerifyToken(c.Request.Context(), token)
		if err != nil {
			c.Next()
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

// ==================== وسائط الأدوار والصلاحيات ====================

// NewAdminMiddleware إنشاء وسيط مصادقة المسؤول
func NewAdminMiddleware() *AdminMiddlewareStruct {
	return &AdminMiddlewareStruct{}
}

// Handle معالجة طلب مصادقة المسؤول
func (m *AdminMiddlewareStruct) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
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

// AdminRequired middleware للتحقق من صلاحيات المشرف
func AdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := utils.GetUserIDFromGinContext(c)
		if userID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "unauthorized",
				"message": "يجب تسجيل الدخول",
			})
			c.Abort()
			return
		}

		// هذا تنفيذ مبسط - يمكن تحديثه للتحقق من دور المستخدم من قاعدة البيانات
		userRole := c.GetString("userRole")
		if userRole != "admin" && userRole != "superadmin" {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"error":   "forbidden",
				"message": "مطلوب صلاحيات مشرف",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// AdminMiddlewareFunc للاستخدام مع router
func AdminMiddlewareFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		// هذا تنفيذ مبسط للتحقق من صلاحيات المشرف
		userRole := c.GetString("userRole")
		if userRole != "admin" && userRole != "superadmin" {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"error":   "forbidden",
				"message": "مطلوب صلاحيات مشرف",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// NewSellerMiddleware إنشاء وسيط مصادقة البائعين
func NewSellerMiddleware() *SellerMiddleware {
	return &SellerMiddleware{}
}

// Handle معالجة طلب مصادقة البائعين
func (m *SellerMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
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

// UserMiddleware وسيط مصادقة المستخدمين العاديين
func UserMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"error":   "صلاحيات غير كافية",
				"message": "يجب تسجيل الدخول للوصول إلى هذا المورد",
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
	return RateLimitWithConfig(100, time.Minute)
}

// RateLimitMiddlewareFunc middleware للحد من معدل الطلبات
func RateLimitMiddlewareFunc() gin.HandlerFunc {
	rateLimitWindow := time.Minute
	maxRequests := 60 // 60 طلب في الدقيقة

	return func(c *gin.Context) {
		clientIP := getClientIP(c.Request)
		
		limiter.mu.Lock()
		defer limiter.mu.Unlock()
		
		now := time.Now()
		windowStart := now.Add(-rateLimitWindow)
		
		// تنظيف الزيارات القديمة
		visits := limiter.visits[clientIP]
		var recentVisits []time.Time
		for _, visit := range visits {
			if visit.After(windowStart) {
				recentVisits = append(recentVisits, visit)
			}
		}
		
		// التحقق من الحد الأقصى
		if len(recentVisits) >= maxRequests {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"success": false,
				"error":   "rate_limit_exceeded",
				"message": "لقد تجاوزت الحد الأقصى للطلبات المسموح بها. الرجاء المحاولة لاحقاً.",
			})
			c.Abort()
			return
		}
		
		// إضافة الزيارة الجديدة
		recentVisits = append(recentVisits, now)
		limiter.visits[clientIP] = recentVisits
		
		c.Next()
	}
}

// RateLimitWithConfig وسيط تحديد المعدل مع تكوين مخصص
func RateLimitWithConfig(requests int, window time.Duration) gin.HandlerFunc {
	type clientLimit struct {
		count    int
		lastSeen time.Time
	}

	clients := make(map[string]*clientLimit)

	return func(c *gin.Context) {
		clientIP := getClientIP(c.Request)

		// تنظيف العملاء القدامى
		if len(clients) > 1000 {
			now := time.Now()
			for ip, limit := range clients {
				if now.Sub(limit.lastSeen) > window {
					delete(clients, ip)
				}
			}
		}

		limit, exists := clients[clientIP]
		if !exists {
			limit = &clientLimit{count: 0, lastSeen: time.Now()}
			clients[clientIP] = limit
		}

		// إعادة تعيين العداد إذا انتهت النافذة الزمنية
		if time.Since(limit.lastSeen) > window {
			limit.count = 0
			limit.lastSeen = time.Now()
		}

		// التحقق من تجاوز الحد
		if limit.count >= requests {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"success": false,
				"error":   "تم تجاوز معدل الطلبات المسموح به",
				"message": fmt.Sprintf("الحد الأقصى هو %d طلب كل %v", requests, window),
			})
			c.Abort()
			return
		}

		// زيادة العداد
		limit.count++
		limit.lastSeen = time.Now()

		// إضافة معلومات التحديد إلى الرأس
		c.Writer.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", requests))
		c.Writer.Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", requests-limit.count))
		c.Writer.Header().Set("X-RateLimit-Reset", fmt.Sprintf("%d", int(window.Seconds())))

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
				"message": fmt.Sprintf("الحجم الأقصى المسموح به هو %d ميغابايت", maxSize/(1024*1024)),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// ValidateContentType وسيط التحقق من نوع المحتوى
func ValidateContentType() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "PATCH" {
			contentType := c.GetHeader("Content-Type")
			if contentType == "" {
				c.JSON(http.StatusBadRequest, gin.H{
					"success": false,
					"error":   "نوع المحتوى مطلوب",
					"message": "يرجى تحديد نوع المحتوى في الرأس",
				})
				c.Abort()
				return
			}
			
			if !strings.Contains(contentType, "application/json") && 
			   !strings.Contains(contentType, "multipart/form-data") &&
			   !strings.Contains(contentType, "application/x-www-form-urlencoded") {
				c.JSON(http.StatusBadRequest, gin.H{
					"success": false,
					"error":   "نوع محتوى غير مدعوم",
					"message": "أنواع المحتوى المدعومة: application/json, multipart/form-data, application/x-www-form-urlencoded",
				})
				c.Abort()
				return
			}
		}

		c.Next()
	}
}

// ==================== وسائط استعادة الأخطاء ====================

// Recovery وسيط استعادة الأخطاء
func Recovery() gin.HandlerFunc {
	return gin.Recovery()
}

// ==================== وسائط الوقت والتحقق ====================

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

// ValidateJSON وسيط التحقق من صحة JSON
func ValidateJSON() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "PATCH" {
			if strings.Contains(c.GetHeader("Content-Type"), "application/json") {
				// التحقق من أن الجسم ليس فارغاً
				if c.Request.ContentLength == 0 {
					c.JSON(http.StatusBadRequest, gin.H{
						"success": false,
						"error":   "جسم الطلب فارغ",
						"message": "يجب أن يحتوي طلب JSON على جسم",
					})
					c.Abort()
					return
				}
			}
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

	if ip := r.Header.Get("CF-Connecting-IP"); ip != "" {
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
	return fmt.Sprintf("req_%d_%d", time.Now().Unix(), time.Now().Nanosecond())
}

// ==================== دوال مساعدة للسياق ====================

// GetUserIDFromContext الحصول على معرف المستخدم من السياق
func GetUserIDFromContext(c *gin.Context) (string, bool) {
	userID, exists := c.Get("userID")
	if !exists {
		return "", false
	}
	return userID.(string), true
}

// GetUserEmailFromContext الحصول على البريد الإلكتروني للمستخدم من السياق
func GetUserEmailFromContext(c *gin.Context) (string, bool) {
	userEmail, exists := c.Get("userEmail")
	if !exists {
		return "", false
	}
	return userEmail.(string), true
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

// GetTokenFromContext الحصول على التوكن من السياق
func GetTokenFromContext(c *gin.Context) (string, bool) {
	token, exists := c.Get("token")
	if !exists {
		return "", false
	}
	return token.(string), true
}

// ==================== تسجيل الوسائط ====================

// RegisterGlobalMiddlewares تسجيل الوسائط العامة
func RegisterGlobalMiddlewares(router *gin.Engine, cfg *config.Config) {
	// الوسائط الأساسية
	router.Use(Recovery())
	router.Use(RequestID())
	router.Use(Logging())
	router.Use(CORSMiddleware())
	router.Use(SecurityHeaders())
	
	// تحديد المعدل يختلف حسب البيئة
	if cfg.Environment == "production" {
		router.Use(RateLimitWithConfig(100, time.Minute)) // 100 طلب/دقيقة في الإنتاج
		router.Use(SizeLimit(10 * 1024 * 1024)) // 10MB في الإنتاج
		router.Use(Timeout(30 * time.Second)) // 30 ثانية في الإنتاج
	} else {
		router.Use(RateLimitWithConfig(1000, time.Minute)) // 1000 طلب/دقيقة في التطوير
		router.Use(SizeLimit(50 * 1024 * 1024)) // 50MB في التطوير
		router.Use(Timeout(60 * time.Second)) // 60 ثانية في التطوير
	}
	
	router.Use(ValidateContentType())
	router.Use(ValidateJSON())
}

// InitializeMiddlewares تهيئة حاوية الوسائط
func InitializeMiddlewares(authService services.AuthService) *MiddlewareContainer {
	return &MiddlewareContainer{
		AuthMiddleware:      NewAuthMiddleware(authService),
		AdminMiddleware:     NewAdminMiddleware(),
		CORSMiddleware:      CORSMiddleware(),
		SecurityMiddleware:  SecurityHeaders(),
		RateLimitMiddleware: RateLimit(),
	}
}

// RegisterAPIMiddlewares تسجيل وسائط API
func RegisterAPIMiddlewares(router *gin.RouterGroup, container *MiddlewareContainer) {
	// تطبيق وسائط الأمان على جميع مسارات API
	router.Use(container.CORSMiddleware)
	router.Use(container.SecurityMiddleware)
	router.Use(container.RateLimitMiddleware)
}

// RegisterProtectedMiddlewares تسجيل وسائط المسارات المحمية
func RegisterProtectedMiddlewares(router *gin.RouterGroup, container *MiddlewareContainer) {
	// تطبيق وسائط المصادقة على المسارات المحمية
	router.Use(container.AuthMiddleware.Handle())
}

// RegisterAdminMiddlewares تسجيل وسائط مسارات الإدارة
func RegisterAdminMiddlewares(router *gin.RouterGroup, container *MiddlewareContainer) {
	// تطبيق وسائط المصادقة والإدارة على مسارات الإدارة
	router.Use(container.AuthMiddleware.Handle())
	router.Use(container.AdminMiddleware.Handle())
}

// NewMiddlewareContainer إنشاء حاوية وسائط جديدة
func NewMiddlewareContainer(authService services.AuthService) *MiddlewareContainer {
	return &MiddlewareContainer{
		AuthMiddleware:      NewAuthMiddleware(authService),
		AdminMiddleware:     NewAdminMiddleware(),
		CORSMiddleware:      CORSMiddleware(),
		SecurityMiddleware:  SecurityHeaders(),
		RateLimitMiddleware: RateLimitWithConfig(100, time.Minute),
	}
}

// ==================== دوال مساعدة إضافية ====================

// ExtractUserIDFromContext استخراج معرف المستخدم من السياق
func ExtractUserIDFromContext(c *gin.Context) string {
	if userID, exists := GetUserIDFromContext(c); exists {
		return userID
	}
	return ""
}

// IsAdmin التحقق مما إذا كان المستخدم مشرفاً
func IsAdmin(c *gin.Context) bool {
	if role, exists := GetUserRoleFromContext(c); exists {
		return role == "admin"
	}
	return false
}

// IsAuthenticated التحقق مما إذا كان المستخدم مصادقاً عليه
func IsAuthenticated(c *gin.Context) bool {
	_, exists := GetUserIDFromContext(c)
	return exists
}

// GetCurrentUser الحصول على معلومات المستخدم الحالي
func GetCurrentUser(c *gin.Context) map[string]interface{} {
	user := make(map[string]interface{})
	
	if userID, exists := GetUserIDFromContext(c); exists {
		user["id"] = userID
	}
	
	if userEmail, exists := GetUserEmailFromContext(c); exists {
		user["email"] = userEmail
	}
	
	if userRole, exists := GetUserRoleFromContext(c); exists {
		user["role"] = userRole
	}
	
	return user
}

// SetUserContext تعيين معلومات المستخدم في السياق
func SetUserContext(c *gin.Context, userID, userEmail, userRole string) {
	c.Set("userID", userID)
	c.Set("userEmail", userEmail)
	c.Set("userRole", userRole)
}

// ClearUserContext مسح معلومات المستخدم من السياق
func ClearUserContext(c *gin.Context) {
	c.Set("userID", nil)
	c.Set("userEmail", nil)
	c.Set("userRole", nil)
	c.Set("token", nil)
}