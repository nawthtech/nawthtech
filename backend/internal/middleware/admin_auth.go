package middleware

import (
	"net/http"
	"strings"
	"time"
)

// AdminAuth وسيط مصادقة المسؤول
func AdminAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, `{"success": false, "error": "مصادقة مطلوبة"}`, http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if !isValidAdminToken(token) {
			http.Error(w, `{"success": false, "error": "صلاحيات غير كافية"}`, http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// RateLimiter وسيط تحديد المعدل
func RateLimiter(requests int, window time.Duration) func(http.Handler) http.Handler {
	// تنفيذ مبسط لتحديد المعدل
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// في الواقع، يجب استخدام خادم Redis أو ذاكرة للتتبع
			next.ServeHTTP(w, r)
		})
	}
}

// ValidateAdminAction التحقق من صحة إجراءات المسؤول
func ValidateAdminAction(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// التحقق من صحة البيانات والإجراءات
		if r.Method == "POST" || r.Method == "PUT" {
			contentType := r.Header.Get("Content-Type")
			if !strings.Contains(contentType, "application/json") {
				http.Error(w, `{"success": false, "error": "نوع المحتوى يجب أن يكون JSON"}`, http.StatusBadRequest)
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

func isValidAdminToken(token string) bool {
	// في الواقع، يجب التحقق من التوكن مع قاعدة البيانات أو خدمة المصادقة
	return token != "" && strings.HasPrefix(token, "admin_")
}

func generateRequestID() string {
	return "req_" + time.Now().Format("20060102150405")
}

func isSystemReadyForUpdate() bool {
	// في الواقع، يجب التحقق من حالة النظام ومدى ملاءمته للتحديث
	return true
}
