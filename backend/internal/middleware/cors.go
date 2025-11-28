package middleware

import (
	"net/http"
	"strings"
	
	"nawthtech/backend/internal/config"
	"nawthtech/backend/internal/logger"
)

// CORS middleware لإعدادات CORS الديناميكية
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// الحصول على إعدادات CORS بناءً على المسار
		corsConfig := config.GetCORSConfig(r.URL.Path)
		
		origin := r.Header.Get("Origin")
		
		// التحقق من النطاق المسموح به
		if !config.ValidateOrigin(origin) {
			logger.Warn("CORS request blocked", map[string]interface{}{
				"origin": origin,
				"path":   r.URL.Path,
			})
			http.Error(w, "Not allowed by CORS", http.StatusForbidden)
			return
		}
		
		// تعيين رؤوس CORS
		if origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		} else {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}
		
		w.Header().Set("Access-Control-Allow-Methods", strings.Join(corsConfig.AllowedMethods, ", "))
		w.Header().Set("Access-Control-Allow-Headers", strings.Join(corsConfig.AllowedHeaders, ", "))
		w.Header().Set("Access-Control-Expose-Headers", strings.Join(corsConfig.ExposedHeaders, ", "))
		
		if corsConfig.AllowCredentials {
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		}
		
		if corsConfig.MaxAge > 0 {
			w.Header().Set("Access-Control-Max-Age", string(rune(corsConfig.MaxAge)))
		}
		
		// معالجة طلبات Preflight
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		
		next.ServeHTTP(w, r)
	})
}

// SecurityHeaders middleware لأمان إضافي
func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		
		next.ServeHTTP(w, r)
	})
}