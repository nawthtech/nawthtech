package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

// ==================== وسائط وهمية للاختبار والتطوير ====================

// DummyCORS وسيط CORS وهمي للتطوير
func DummyCORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-Request-ID")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length, X-Request-ID")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// DummyLogger وسيط تسجيل وهمي للتطوير
func DummyLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		
		// معالجة الطلب
		c.Next()
		
		// تسجيل بسيط في وضع التطوير
		duration := time.Since(start)
		fmt.Printf("[Dummy Logger] %s %s - %s - %v\n",
			c.Request.Method,
			c.Request.URL.Path,
			c.Writer.Status(),
			duration,
		)
	}
}

// DummyAuth وسيط مصادقة وهمي للاختبار
func DummyAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// في وضع التطوير، نسمح بالوصول بدون مصادقة
		// ولكن نضيف بيانات مستخدم وهمية للاختبار
		
		// التحقق من وجود توكن في الرأس
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			// إذا كان هناك توكن، نحاول استخراج معلومات منه
			c.Set("userID", "test_user_123")
			c.Set("userEmail", "test@example.com")
			c.Set("userRole", "user")
		} else {
			// إذا لم يكن هناك توكن، نستخدم مستخدم ضيف
			c.Set("userID", "guest_user")
			c.Set("userRole", "guest")
		}
		
		c.Next()
	}
}

// DummyAdminAuth وسيط مصادقة مسؤول وهمي
func DummyAdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// في وضع التطوير، نضيف بيانات مسؤول وهمية
		c.Set("userID", "admin_user_123")
		c.Set("userEmail", "admin@example.com")
		c.Set("userRole", "admin")
		
		c.Next()
	}
}

// DummyRateLimit وسيط تحديد معدل وهمي
func DummyRateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		// في وضع التطوير، لا نطبق تحديد المعدل
		// ولكن نضيف الرؤوس فقط للمحاكاة
		c.Writer.Header().Set("X-RateLimit-Limit", "1000")
		c.Writer.Header().Set("X-RateLimit-Remaining", "999")
		c.Writer.Header().Set("X-RateLimit-Reset", "60")
		
		c.Next()
	}
}

// DummyRecovery وسيط استعادة أخطاء وهمي
func DummyRecovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// في وضع التطوير، نطبع الخطأ فقط
				fmt.Printf("[Dummy Recovery] Recovered from panic: %v\n", err)
				
				// إرجاع استجابة خطأ بسيطة
				c.JSON(500, gin.H{
					"success": false,
					"error":   "Internal Server Error (Dummy Mode)",
					"message": fmt.Sprintf("Panic recovered: %v", err),
				})
				
				c.Abort()
			}
		}()
		
		c.Next()
	}
}

// DummyRequestID وسيط معرف طلب وهمي
func DummyRequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := fmt.Sprintf("dummy_req_%d", time.Now().UnixNano())
		c.Writer.Header().Set("X-Request-ID", requestID)
		c.Set("requestID", requestID)
		
		c.Next()
	}
}

// DummySecurityHeaders وسيط رؤوس أمان وهمي
func DummySecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")
		c.Writer.Header().Set("X-Frame-Options", "DENY")
		c.Writer.Header().Set("X-XSS-Protection", "1; mode=block")
		
		c.Next()
	}
}

// ==================== دوال مساعدة للوضع الوهمي ====================

// IsDummyMode التحقق إذا كان النظام يعمل في الوضع الوهمي
func IsDummyMode() bool {
	// يمكن تغيير هذا بناءً على متغير بيئة
	return true // مؤقتاً، نعتبر أننا في وضع التطوير
}

// GetDummyUserID الحصول على معرف المستخدم الوهمي
func GetDummyUserID(c *gin.Context) string {
	if userID, exists := c.Get("userID"); exists {
		return userID.(string)
	}
	return "dummy_user"
}

// GetDummyUserRole الحصول على دور المستخدم الوهمي
func GetDummyUserRole(c *gin.Context) string {
	if userRole, exists := c.Get("userRole"); exists {
		return userRole.(string)
	}
	return "guest"
}

// ==================== وسائط مختلطة للاختبار ====================

// DevelopmentMiddlewares وسائط التطوير الكاملة
func DevelopmentMiddlewares() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		DummyRecovery(),
		DummyRequestID(),
		DummyLogger(),
		DummyCORS(),
		DummySecurityHeaders(),
		DummyRateLimit(),
		DummyAuth(),
	}
}

// TestingMiddlewares وسائط الاختبار
func TestingMiddlewares() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		DummyRecovery(),
		DummyRequestID(),
		DummyCORS(),
		DummySecurityHeaders(),
	}
}

// MockAuthMiddlewares وسائط مصادقة وهمية
func MockAuthMiddlewares(role string) []gin.HandlerFunc {
	middlewares := []gin.HandlerFunc{
		DummyCORS(),
		DummySecurityHeaders(),
	}
	
	// إضافة وسيط المصادقة المناسب للدور
	switch role {
	case "admin":
		middlewares = append(middlewares, DummyAdminAuth())
	case "user":
		middlewares = append(middlewares, DummyAuth())
	default:
		middlewares = append(middlewares, DummyAuth())
	}
	
	return middlewares
}

// ==================== واجهات وسائط وهمية ====================

// DummyMiddlewareContainer حاوية الوسائط الوهمية
type DummyMiddlewareContainer struct {
	CORSMiddleware      gin.HandlerFunc
	LoggerMiddleware    gin.HandlerFunc
	AuthMiddleware      gin.HandlerFunc
	AdminMiddleware     gin.HandlerFunc
	SecurityMiddleware  gin.HandlerFunc
	RateLimitMiddleware gin.HandlerFunc
	RecoveryMiddleware  gin.HandlerFunc
	RequestIDMiddleware gin.HandlerFunc
}

// NewDummyMiddlewareContainer إنشاء حاوية وسائط وهمية جديدة
func NewDummyMiddlewareContainer() *DummyMiddlewareContainer {
	return &DummyMiddlewareContainer{
		CORSMiddleware:      DummyCORS(),
		LoggerMiddleware:    DummyLogger(),
		AuthMiddleware:      DummyAuth(),
		AdminMiddleware:     DummyAdminAuth(),
		SecurityMiddleware:  DummySecurityHeaders(),
		RateLimitMiddleware: DummyRateLimit(),
		RecoveryMiddleware:  DummyRecovery(),
		RequestIDMiddleware: DummyRequestID(),
	}
}

// ApplyDummyMiddlewares تطبيق جميع الوسائط الوهمية
func ApplyDummyMiddlewares(router *gin.Engine) {
	container := NewDummyMiddlewareContainer()
	
	router.Use(container.RecoveryMiddleware)
	router.Use(container.RequestIDMiddleware)
	router.Use(container.LoggerMiddleware)
	router.Use(container.CORSMiddleware)
	router.Use(container.SecurityMiddleware)
	router.Use(container.RateLimitMiddleware)
}

// ApplyDummyAuthMiddlewares تطبيق وسائط المصادقة الوهمية
func ApplyDummyAuthMiddlewares(router *gin.RouterGroup, requireAdmin bool) {
	container := NewDummyMiddlewareContainer()
	
	router.Use(container.AuthMiddleware)
	if requireAdmin {
		router.Use(container.AdminMiddleware)
	}
}

// ==================== وسائط للاختبارات الوحداتية ====================

// TestCORS وسيط CORS للاختبارات
func TestCORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "*")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}
		
		c.Next()
	}
}

// TestAuth وسيط مصادقة للاختبارات
func TestAuth(userID, userRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("userID", userID)
		c.Set("userRole", userRole)
		c.Set("userEmail", userID+"@test.com")
		
		c.Next()
	}
}

// NoOpMiddleware وسيط لا يقوم بأي عمل (للاختبارات)
func NoOpMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}

// MockRateLimit وسيط تحديد معدل وهمي للاختبارات
func MockRateLimit(limit, remaining int) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", limit))
		c.Writer.Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))
		c.Writer.Header().Set("X-RateLimit-Reset", "60")
		
		c.Next()
	}
}


// ==================== وظائف المساعدة للوضع الوهمي ====================

// SetupDummyMode إعداد النظام للعمل في الوضع الوهمي
func SetupDummyMode(router *gin.Engine) {
	fmt.Println("🚀 Running in DUMMY MODE - All middleware are mocked")
	
	// تطبيق جميع الوسائط الوهمية
	ApplyDummyMiddlewares(router)
	
	// إعداد مسارات خاصة للوضع الوهمي
	setupDummyRoutes(router)
}

// setupDummyRoutes إعداد مسارات وهمية للاختبار
func setupDummyRoutes(router *gin.Engine) {
	// مسارات معلومات الوضع الوهمي
	router.GET("/dummy/info", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"mode":         "dummy",
			"description":  "Running in dummy/testing mode",
			"features":     []string{"mocked_auth", "mocked_rate_limit", "no_real_db"},
			"available_endpoints": []string{
				"/dummy/info",
				"/dummy/auth/test",
				"/dummy/admin/test",
			},
		})
	})
	
	// مسار اختبار المصادقة الوهمية
	router.GET("/dummy/auth/test", DummyAuth(), func(c *gin.Context) {
		userID, _ := c.Get("userID")
		userRole, _ := c.Get("userRole")
		
		c.JSON(200, gin.H{
			"success":  true,
			"message":  "Dummy auth test successful",
			"user_id":  userID,
			"user_role": userRole,
			"mode":     "dummy",
		})
	})
	
	// مسار اختبار إدارة وهمي
	router.GET("/dummy/admin/test", DummyAdminAuth(), func(c *gin.Context) {
		c.JSON(200, gin.H{
			"success": true,
			"message": "Dummy admin test successful",
			"role":    "admin",
			"mode":    "dummy",
		})
	})
}