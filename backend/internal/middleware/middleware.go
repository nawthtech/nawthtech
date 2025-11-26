package middleware

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"strings"
	"time"

	"backend-app/internal/config"
	"backend-app/internal/logger"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

// trustProxyConfig تكوين وسيط الثقة بالبروكسي
type trustProxyConfig struct {
	ErrorLogger *slog.Logger
}

// TrustProxy وسيط الثقة بالبروكسي
func TrustProxy(config *trustProxyConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if config == nil {
				config = &trustProxyConfig{}
			}

			// الحصول على عنوان IP الحقيقي من الرؤوس
			if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
				r.RemoteAddr = realIP
			} else if forwardedFor := r.Header.Get("X-Forwarded-For"); forwardedFor != "" {
				// أخذ أول عنوان في القائمة
				ips := strings.Split(forwardedFor, ",")
				if len(ips) > 0 {
					r.RemoteAddr = strings.TrimSpace(ips[0])
				}
			}

			// تحديث الطلب مع معلومات البروتوكول الحقيقية
			if proto := r.Header.Get("X-Forwarded-Proto"); proto != "" {
				r.URL.Scheme = proto
			} else if r.TLS != nil {
				r.URL.Scheme = "https"
			} else {
				r.URL.Scheme = "http"
			}

			if host := r.Header.Get("X-Forwarded-Host"); host != "" {
				r.URL.Host = host
				r.Host = host
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RequestID وسيط إضافة معرف الطلب
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID()
		}

		// إضافة معرف الطلب إلى الرأس
		w.Header().Set("X-Request-ID", requestID)

		// إضافة معرف الطلب إلى السياق
		ctx := r.Context()
		ctx = context.WithValue(ctx, "requestID", requestID)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

// Logger وسيط التسجيل
func Logger() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// استخدام ResponseWriter معاد لالتقاط حالة الاستجابة
			ww := chimiddleware.NewWrapResponseWriter(w, r.ProtoMajor)

			// معالجة الطلب
			next.ServeHTTP(ww, r)

			// حساب مدة التنفيذ
			duration := time.Since(start)

			// جمع معلومات التسجيل
			fields := []interface{}{
				"method", r.Method,
				"path", r.URL.Path,
				"status", ww.Status(),
				"duration", duration.String(),
				"bytes", ww.BytesWritten(),
				"user_agent", r.UserAgent(),
				"ip", getClientIP(r),
			}

			// إضافة معرف الطلب إذا كان متوفراً
			if requestID := r.Context().Value("requestID"); requestID != nil {
				fields = append(fields, "request_id", requestID)
			}

			// تسجيل بناءً على حالة الاستجابة
			status := ww.Status()
			switch {
			case status >= 500:
				logger.Stderr.Error("خطأ في الخادم", fields...)
			case status >= 400:
				logger.Stdout.Warn("خطأ في العميل", fields...)
			default:
				logger.Stdout.Info("طلب معالَج", fields...)
			}
		})
	}
}

// AuthMiddleware وسيط المصادقة الأساسي
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// الحصول على التوكن من الرأس
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, `{"success": false, "error": "مصادقة مطلوبة"}`, http.StatusUnauthorized)
			return
		}

		// التحقق من صيغة التوكن (Bearer token)
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, `{"success": false, "error": "صيغة توكن غير صالحة"}`, http.StatusUnauthorized)
			return
		}

		token := parts[1]

		// التحقق من صحة التوكن (هنا يمكن إضافة منطق التحقق الفعلي)
		userID, err := validateToken(token)
		if err != nil {
			http.Error(w, `{"success": false, "error": "توكن غير صالح"}`, http.StatusUnauthorized)
			return
		}

		// إضافة معلومات المستخدم إلى السياق
		ctx := r.Context()
		ctx = context.WithValue(ctx, "userID", userID)
		ctx = context.WithValue(ctx, "token", token)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

// AdminAuth وسيط مصادقة المسؤول
func AdminAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// الحصول على معلومات المستخدم من السياق
		userID := r.Context().Value("userID")
		if userID == nil {
			http.Error(w, `{"success": false, "error": "صلاحيات غير كافية"}`, http.StatusForbidden)
			return
		}

		// التحقق من صلاحيات المسؤول (هنا يمكن إضافة منطق التحقق الفعلي)
		isAdmin, err := checkAdminPermissions(userID.(string))
		if err != nil || !isAdmin {
			http.Error(w, `{"success": false, "error": "صلاحيات إدارة مطلوبة"}`, http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// RateLimiter وسيط تحديد المعدل
func RateLimiter(requests int, window time.Duration) func(http.Handler) http.Handler {
	// تنفيذ مبسط لتحديد المعدل (في الواقع يجب استخدام Redis أو ذاكرة مشتركة)
	type clientLimit struct {
		count    int
		lastSeen time.Time
	}

	clients := make(map[string]*clientLimit)
	
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			clientIP := getClientIP(r)
			
			now := time.Now()
			
			// تنظيف العملاء القدامى
			if len(clients) > 10000 { // حد أقصى للذاكرة
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
				http.Error(w, `{"success": false, "error": "تم تجاوز معدل الطلبات المسموح به"}`, http.StatusTooManyRequests)
				return
			}
			
			// زيادة العداد
			limit.count++
			limit.lastSeen = now
			
			// إضافة معلومات التحديد إلى الرأس
			w.Header().Set("X-RateLimit-Limit", string(requests))
			w.Header().Set("X-RateLimit-Remaining", string(requests-limit.count))
			w.Header().Set("X-RateLimit-Reset", string(now.Add(window).Unix()))
			
			next.ServeHTTP(w, r)
		})
	}
}

// ValidateAdminAction التحقق من صحة إجراءات المسؤول
func ValidateAdminAction(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// التحقق من صحة البيانات والإجراءات للمسؤول
		if r.Method == "POST" || r.Method == "PUT" {
			contentType := r.Header.Get("Content-Type")
			if !strings.Contains(contentType, "application/json") {
				http.Error(w, `{"success": false, "error": "نوع المحتوى يجب أن يكون JSON"}`, http.StatusBadRequest)
				return
			}
			
			// التحقق من حجم الجسم للطلبات الكبيرة
			if r.ContentLength > 10*1024*1024 { // 10MB
				http.Error(w, `{"success": false, "error": "حجم البيانات كبير جداً"}`, http.StatusRequestEntityTooLarge)
				return
			}
		}
		
		next.ServeHTTP(w, r)
	})
}

// UpdateMiddleware وسيط تحديث النظام
func UpdateMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// إضافة معرف الطلب
		requestID := generateRequestID()
		w.Header().Set("X-Request-ID", requestID)

		// التحقق من حالة النظام قبل التحديث
		if !isSystemReadyForUpdate() {
			http.Error(w, `{"success": false, "error": "النظام غير جاهز للتحديث"}`, http.StatusServiceUnavailable)
			return
		}

		next.ServeHTTP(w, r)
	})
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
	return strings.ReplaceAll(time.Now().Format("20060102150405.000000"), ".", "") + "-" + randomString(8)
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

// validateToken التحقق من صحة التوكن (تنفيذ مبدئي)
func validateToken(token string) (string, error) {
	// في الواقع، يجب التحقق من التوكن مع قاعدة البيانات أو خدمة المصادقة
	// هذا تنفيذ مبسط للتوضيح
	if token == "" {
		return "", fmt.Errorf("توكن فارغ")
	}
	
	// محاكاة التحقق من التوكن
	if strings.HasPrefix(token, "valid_") {
		return strings.TrimPrefix(token, "valid_"), nil
	}
	
	return "", fmt.Errorf("توكن غير صالح")
}

// checkAdminPermissions التحقق من صلاحيات المسؤول (تنفيذ مبدئي)
func checkAdminPermissions(userID string) (bool, error) {
	// في الواقع، يجب التحقق من الصلاحيات مع قاعدة البيانات
	// هذا تنفيذ مبسط للتوضيح
	return strings.HasPrefix(userID, "admin_"), nil
}

// isSystemReadyForUpdate التحقق من جاهزية النظام للتحديث
func isSystemReadyForUpdate() bool {
	// في الواقع، يجب التحقق من حالة النظام ومدى ملاءمته للتحديث
	// هذا تنفيذ مبسط للتوضيح
	return true
}

// Register تسجيل جميع الوسائط
func Register(r *chi.Mux) {
	// الوسائط الأساسية
	r.Use(chimiddleware.Recoverer)
	r.Use(TrustProxy(&trustProxyConfig{
		ErrorLogger: logger.Stderr,
	}))
	r.Use(RequestID)
	r.Use(Logger())

	// وسيط CORS
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   config.Cors.AllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "X-Request-ID"},
		ExposedHeaders:   []string{"Link", "X-Request-ID", "X-RateLimit-Limit", "X-RateLimit-Remaining", "X-RateLimit-Reset"},
		AllowCredentials: true,
		MaxAge:           300, // 5 دقائق
	}))

	// وسائط إضافية
	r.Use(chimiddleware.Compress(5)) // ضغط GZIP
	r.Use(chimiddleware.Timeout(60 * time.Second)) // وقت انتهاء للطلبات
	r.Use(chimiddleware.Throttle(1000)) // تحديد معدل الطلبات الأساسي
	r.Use(chimiddleware.Heartbeat("/health")) // نقطة فحص الصحة
}
