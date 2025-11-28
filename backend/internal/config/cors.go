package config

import (
	"os"
	"regexp"
	"strings"
)

// Cors تكوين CORS
type Cors struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	ExposedHeaders   []string
	AllowCredentials bool
	MaxAge           int
}

// getAllowedOrigins الحصول على قائمة النطاقات المسموح بها ديناميكياً
func getAllowedOrigins() []string {
	baseDomains := []string{
		// النطاقات الرئيسية للمنصة
		"https://nawthtech.com",
		"https://www.nawthtech.com",
		"https://app.nawthtech.com",
		"https://admin.nawthtech.com",
		"https://api.nawthtech.com",
		"https://store.nawthtech.com",
		"https://dashboard.nawthtech.com",

		// نطاقات التطوير المحلي
		"http://localhost:3000",
		"http://localhost:5173",
		"http://localhost:5000",
		"http://127.0.0.1:3000",
		"http://127.0.0.1:5173",
		"http://localhost:8080",
		"http://127.0.0.1:8080",
		
		// نطاقات البيئات الأخرى
		"https://staging.nawthtech.com",
		"https://dev.nawthtech.com",

		// نطاقات خدمات Microsoft 365
		"https://outlook.office.com",
		"https://office.com",
		"https://microsoft365.com",

		// نطاقات خدمات التحليلات البديلة
		"https://plausible.io",         // Plausible Analytics
		"https://us.plausible.io",      // US region for Plausible
		"https://matomo.org",           // Matomo Analytics
		"https://fathom-analytics.com", // Fathom Analytics
		"https://app.fathom-analytics.com",

		// نطاقات الخدمات الأخرى المستخدمة
		"https://cloudinary.com",       // Cloudinary
		"https://api.cloudinary.com",
		"https://res.cloudinary.com",
		"https://stripe.com",           // Stripe Payments
		"https://api.stripe.com",
		"https://js.stripe.com",
		"https://github.com",           // GitHub
		"https://api.github.com",
		"https://railway.app",          // Railway
		"https://up.railway.app",
	}

	// إضافة نطاقات من متغيرات البيئة
	if clientURL := os.Getenv("CLIENT_URL"); clientURL != "" {
		baseDomains = append(baseDomains, clientURL)
	}
	if adminURL := os.Getenv("ADMIN_URL"); adminURL != "" {
		baseDomains = append(baseDomains, adminURL)
	}
	if storeURL := os.Getenv("STORE_URL"); storeURL != "" {
		baseDomains = append(baseDomains, storeURL)
	}
	if analyticsURL := os.Getenv("ANALYTICS_URL"); analyticsURL != "" {
		baseDomains = append(baseDomains, analyticsURL)
	}

	// إضافة نطاقات ديناميكية من متغير البيئة
	if dynamicOrigins := os.Getenv("ALLOWED_ORIGINS"); dynamicOrigins != "" {
		origins := strings.Split(dynamicOrigins, ",")
		for _, origin := range origins {
			trimmed := strings.TrimSpace(origin)
			if trimmed != "" {
				baseDomains = append(baseDomains, trimmed)
			}
		}
	}

	// إزالة التكرارات
	return removeDuplicates(baseDomains)
}

// removeDuplicates إزالة القيم المكررة من المصفوفة
func removeDuplicates(slice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, item := range slice {
		if _, value := keys[item]; !value {
			keys[item] = true
			list = append(list, item)
		}
	}
	return list
}

// isOriginAllowed التحقق من النطاق المسموح به
func isOriginAllowed(origin string, allowedOrigins []string) bool {
	if origin == "" {
		return true // طلبات بدون origin (تطبيقات محمولة، إلخ)
	}
	
	// التحقق المباشر
	for _, allowed := range allowedOrigins {
		if origin == allowed {
			return true
		}
	}
	
	// التحقق باستخدام patterns
	for _, pattern := range allowedOrigins {
		if strings.HasPrefix(pattern, ".") && strings.HasSuffix(origin, pattern) {
			return true
		}
		if strings.Contains(pattern, "*") {
			regexPattern := strings.ReplaceAll(pattern, "*", ".*")
			regexPattern = "^" + regexPattern + "$"
			if matched, _ := regexp.MatchString(regexPattern, origin); matched {
				return true
			}
		}
	}
	
	return false
}

// CORSOptions إعدادات CORS الرئيسية
type CORSOptions struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	ExposedHeaders   []string
	AllowCredentials bool
	MaxAge           int
}

// GetCORSConfig الحصول على إعدادات CORS بناءً على المسار
func GetCORSConfig(path string) CORSOptions {
	allowedOrigins := getAllowedOrigins()
	
	// مسارات التحليلات البديلة
	if strings.HasPrefix(path, "/api/analytics") || 
	   strings.HasPrefix(path, "/analytics") ||
	   strings.Contains(path, "/plausible") || 
	   strings.Contains(path, "/matomo") {
		return CORSOptions{
			AllowedOrigins: []string{
				"https://plausible.io",
				"https://matomo.org",
				"https://fathom-analytics.com",
				"https://nawthtech.com",
			},
			AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
			AllowedHeaders:   []string{"Content-Type", "User-Agent", "X-Forwarded-For", "X-Plausible-Token", "X-Matomo-Token"},
			AllowCredentials: false,
			MaxAge:           3600,
		}
	}
	
	// مسارات ويب هووكs للخدمات الخارجية
	if strings.HasPrefix(path, "/webhook/stripe") || 
	   strings.HasPrefix(path, "/webhook/cloudinary") ||
	   strings.HasPrefix(path, "/api/webhooks") {
		return CORSOptions{
			AllowedOrigins: []string{
				"https://stripe.com",
				"https://api.stripe.com",
				"https://js.stripe.com",
				"https://cloudinary.com",
				"https://api.cloudinary.com",
				"https://res.cloudinary.com",
				"https://railway.app",
				"https://outlook.office.com",
				"https://office.com",
				"https://microsoft365.com",
			},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Content-Type", "Authorization", "X-API-Key", "X-Stripe-Signature"},
			AllowCredentials: false,
			MaxAge:           7200,
		}
	}
	
	// المسارات الإدارية
	if strings.HasPrefix(path, "/admin") || strings.HasPrefix(path, "/api/admin") {
		return CORSOptions{
			AllowedOrigins: []string{
				"https://admin.nawthtech.com",
				"https://dashboard.nawthtech.com",
				"http://localhost:3001",
			},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS", "HEAD"},
			AllowedHeaders:   []string{"Content-Type", "Authorization", "X-Requested-With", "X-API-Key", "Accept", "Origin", "X-Client-Version", "X-Device-ID", "X-Platform"},
			ExposedHeaders:   []string{"X-Request-ID", "X-Response-Time", "X-API-Version", "X-RateLimit-Limit", "X-RateLimit-Remaining"},
			AllowCredentials: true,
			MaxAge:           86400,
		}
	}

	// الإعدادات الافتراضية
	return CORSOptions{
		AllowedOrigins: allowedOrigins,
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS", "HEAD"},
		AllowedHeaders: []string{
			"Content-Type",
			"Authorization", 
			"X-Requested-With",
			"X-API-Key",
			"Accept",
			"Origin",
			"X-Client-Version",
			"X-Device-ID",
			"X-Platform",
			"X-Plausible-Token",
			"X-Matomo-Token",
			"X-Fathom-Key",
			"X-Request-ID",
			"Cache-Control",
		},
		ExposedHeaders: []string{
			"X-Request-ID",
			"X-Response-Time",
			"X-API-Version",
			"X-RateLimit-Limit",
			"X-RateLimit-Remaining",
			"X-Total-Count",
			"Content-Length",
		},
		AllowCredentials: true,
		MaxAge:           86400,
	}
}

// ValidateOrigin التحقق من النطاق المسموح به
func ValidateOrigin(origin string) bool {
	allowedOrigins := getAllowedOrigins()
	
	// في بيئة التطوير، السماح مع تسجيل التحذيرات
	if os.Getenv("ENVIRONMENT") == "development" || os.Getenv("ENVIRONMENT") == "" {
		if origin != "" && !isOriginAllowed(origin, allowedOrigins) {
			// سيتم التعامل مع التسجيل في middleware
			return true
		}
		return true
	}

	return origin == "" || isOriginAllowed(origin, allowedOrigins)
}

// GetCORSStats إحصائيات CORS
func GetCORSStats() map[string]interface{} {
	allowedOrigins := getAllowedOrigins()
	
	analyticsCount := 0
	microsoftCount := 0
	storageCount := 0
	paymentsCount := 0
	localCount := 0
	productionCount := 0
	
	for _, origin := range allowedOrigins {
		switch {
		case strings.Contains(origin, "plausible") || strings.Contains(origin, "matomo") || strings.Contains(origin, "fathom"):
			analyticsCount++
		case strings.Contains(origin, "microsoft") || strings.Contains(origin, "office"):
			microsoftCount++
		case strings.Contains(origin, "cloudinary"):
			storageCount++
		case strings.Contains(origin, "stripe"):
			paymentsCount++
		case strings.Contains(origin, "localhost") || strings.Contains(origin, "127.0.0.1"):
			localCount++
		case strings.Contains(origin, "nawthtech.com") && !strings.Contains(origin, "localhost"):
			productionCount++
		}
	}
	
	environment := os.Getenv("ENVIRONMENT")
	if environment == "" {
		environment = "development"
	}
	
	return map[string]interface{}{
		"totalAllowedOrigins": len(allowedOrigins),
		"services": map[string]int{
			"analytics":  analyticsCount,
			"microsoft":  microsoftCount,
			"storage":    storageCount,
			"payments":   paymentsCount,
			"local":      localCount,
			"production": productionCount,
		},
		"environment": environment,
	}
}

// GetDefaultCORSConfig الحصول على إعدادات CORS الافتراضية
func GetDefaultCORSConfig() Cors {
	allowedOrigins := getAllowedOrigins()
	
	return Cors{
		AllowedOrigins: allowedOrigins,
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS", "HEAD"},
		AllowedHeaders: []string{
			"Content-Type",
			"Authorization", 
			"X-Requested-With",
			"X-API-Key",
			"Accept",
			"Origin",
			"X-Client-Version",
			"X-Device-ID",
			"X-Platform",
			"X-Plausible-Token",
			"X-Matomo-Token",
			"X-Fathom-Key",
			"X-Request-ID",
			"Cache-Control",
			"X-CSRF-Token",
		},
		ExposedHeaders: []string{
			"X-Request-ID",
			"X-Response-Time",
			"X-API-Version",
			"X-RateLimit-Limit",
			"X-RateLimit-Remaining",
			"X-Total-Count",
			"Content-Length",
		},
		AllowCredentials: true,
		MaxAge:           86400,
	}
}

// IsDevelopmentEnvironment التحقق إذا كانت البيئة تطوير
func IsDevelopmentEnvironment() bool {
	env := os.Getenv("ENVIRONMENT")
	return env == "development" || env == ""
}

// IsProductionEnvironment التحقق إذا كانت البيئة إنتاج
func IsProductionEnvironment() bool {
	env := os.Getenv("ENVIRONMENT")
	return env == "production"
}

// IsStagingEnvironment التحقق إذا كانت البيئة تجريبية
func IsStagingEnvironment() bool {
	env := os.Getenv("ENVIRONMENT")
	return env == "staging"
}