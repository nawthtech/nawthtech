// backend/internal/ai/verification/middleware.go
package verification

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// VerificationMiddleware middleware للتحقق من المحتوى
type VerificationMiddleware struct {
	verifier *Verifier
	logger   *logrus.Logger
	enabled  bool
}

// NewVerificationMiddleware إنشاء middleware جديد
func NewVerificationMiddleware(verifier *Verifier, logger *logrus.Logger, enabled bool) *VerificationMiddleware {
	return &VerificationMiddleware{
		verifier: verifier,
		logger:   logger,
		enabled:  enabled,
	}
}

// ContentVerification middleware للتحقق من محتوى الطلبات
func (m *VerificationMiddleware) ContentVerification(criteria VerificationCriteria) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !m.enabled {
			c.Next()
			return
		}

		startTime := time.Now()

		// الحصول على المحتوى من الطلب
		var requestBody map[string]interface{}
		if err := c.ShouldBindJSON(&requestBody); err != nil {
			c.Next()
			return
		}

		// استخراج المحتوى للتحقق
		content := extractContentForVerification(requestBody)
		if content == "" {
			c.Next()
			return
		}

		// التحقق من المحتوى
		ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
		defer cancel()

		result, err := m.verifier.Verify(ctx, content, WithCustomCriteria(criteria))
		if err != nil {
			m.logger.WithError(err).Warn("Content verification failed")
			c.Next()
			return
		}

		// تسجيل النتيجة
		latency := time.Since(startTime)
		m.logger.WithFields(logrus.Fields{
			"isValid":    result.IsValid,
			"confidence": result.Confidence,
			"latency":    latency.Milliseconds(),
			"path":       c.Request.URL.Path,
		}).Info("Content verification completed")

		// إضافة النتيجة إلى السياق
		c.Set("verification_result", result)
		c.Set("verification_latency", latency)

		// إذا لم يكن المحتوى صالحاً، يمكن رفض الطلب
		if !result.IsValid && shouldBlockInvalidContent(criteria, result) {
			c.JSON(400, gin.H{
				"error": "Content verification failed",
				"issues": result.Issues,
				"suggestions": result.Suggestions,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// extractContentForVerification استخراج المحتوى للتحقق
func extractContentForVerification(requestBody map[string]interface{}) string {
	// محاولة استخراج المحتوى من الحقول الشائعة
	contentFields := []string{"content", "text", "message", "body", "input", "prompt"}
	
	for _, field := range contentFields {
		if val, ok := requestBody[field]; ok {
			if str, ok := val.(string); ok && str != "" {
				return str
			}
		}
	}
	
	// إذا لم يتم العثور على محتوى نصي، تحويل كامل الطلب إلى JSON
	jsonBytes, err := json.Marshal(requestBody)
	if err == nil {
		return string(jsonBytes)
	}
	
	return ""
}

// shouldBlockInvalidContent تحديد إذا كان يجب حظر المحتوى غير الصالح
func shouldBlockInvalidContent(criteria VerificationCriteria, result *VerificationResult) bool {
	// حظر إذا كانت الثقة منخفضة جداً
	if result.Confidence < 0.3 {
		return true
	}
	
	// حظر إذا كان هناك مشاكل خطيرة
	seriousIssues := []string{
		"toxic", "hateful", "harmful", "dangerous",
		"illegal", "explicit", "violent", "threatening",
	}
	
	for _, issue := range result.Issues {
		issueLower := strings.ToLower(issue)
		for _, serious := range seriousIssues {
			if strings.Contains(issueLower, serious) {
				return true
			}
		}
	}
	
	// التحقق من الفئات المهمة
	if criteria.Safety && result.Categories["safety"].Passed == false {
		return true
	}
	if criteria.Toxicity && result.Categories["toxicity"].Passed == false {
		return true
	}
	
	return false
}

// RateLimitingMiddleware middleware للتحكم في معدل الطلبات
func (m *VerificationMiddleware) RateLimitingMiddleware(maxRequests int, window time.Duration) gin.HandlerFunc {
	requests := make(map[string][]time.Time)
	
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		now := time.Now()
		
		// تنظيف الطلبات القديمة
		if timestamps, exists := requests[clientIP]; exists {
			validTimestamps := []time.Time{}
			for _, ts := range timestamps {
				if now.Sub(ts) < window {
					validTimestamps = append(validTimestamps, ts)
				}
			}
			requests[clientIP] = validTimestamps
		}
		
		// التحقق من الحد
		if len(requests[clientIP]) >= maxRequests {
			c.JSON(429, gin.H{
				"error": "Rate limit exceeded",
				"message": "Too many verification requests",
				"retry_after": window.Seconds(),
			})
			c.Abort()
			return
		}
		
		// تسجيل الطلب
		requests[clientIP] = append(requests[clientIP], now)
		
		c.Next()
	}
}

// LoggingMiddleware middleware للتسجيل
func (m *VerificationMiddleware) LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method
		
		// معالجة الطلب
		c.Next()
		
		// تسجيل بعد المعالجة
		latency := time.Since(startTime)
		status := c.Writer.Status()
		
		fields := logrus.Fields{
			"method":     method,
			"path":       path,
			"status":     status,
			"latency":    latency.Milliseconds(),
			"client_ip":  c.ClientIP(),
			"user_agent": c.Request.UserAgent(),
		}
		
		// إضافة نتيجة التحقق إذا كانت موجودة
		if result, exists := c.Get("verification_result"); exists {
			if verificationResult, ok := result.(*VerificationResult); ok {
				fields["verification_valid"] = verificationResult.IsValid
				fields["verification_confidence"] = verificationResult.Confidence
			}
		}
		
		if status >= 400 {
			m.logger.WithFields(fields).Warn("Request completed with error")
		} else {
			m.logger.WithFields(fields).Info("Request completed")
		}
	}
}

// ErrorHandlingMiddleware middleware لمعالجة الأخطاء
func (m *VerificationMiddleware) ErrorHandlingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				m.logger.WithFields(logrus.Fields{
					"error":   err,
					"path":    c.Request.URL.Path,
					"method":  c.Request.Method,
				}).Error("Panic recovered")
				
				c.JSON(500, gin.H{
					"error":   "Internal server error",
					"message": "An unexpected error occurred",
				})
			}
		}()
		
		c.Next()
		
		// معالجة أخطاء التحقق
		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				m.logger.WithError(err).Error("Request error")
			}
			
			// إرجاع استجابة منسقة للأخطاء
			if !c.Writer.Written() {
				c.JSON(400, gin.H{
					"error":   "Validation error",
					"message": c.Errors.String(),
				})
			}
		}
	}
}

// MetricsMiddleware middleware للمقاييس
func (m *VerificationMiddleware) MetricsMiddleware() gin.HandlerFunc {
	type metrics struct {
		TotalRequests   int64
		FailedRequests  int64
		TotalLatency    time.Duration
		Verifications   int64
		FailedVerifications int64
	}
	
	stats := &metrics{}
	
	return func(c *gin.Context) {
		startTime := time.Now()
		
		c.Next()
		
		// تحديث الإحصائيات
		latency := time.Since(startTime)
		stats.TotalRequests++
		stats.TotalLatency += latency
		
		if c.Writer.Status() >= 400 {
			stats.FailedRequests++
		}
		
		// تحديث إحصائيات التحقق
		if result, exists := c.Get("verification_result"); exists {
			stats.Verifications++
			if verificationResult, ok := result.(*VerificationResult); ok && !verificationResult.IsValid {
				stats.FailedVerifications++
			}
		}
		
		// تسجيل المقاييس كل 100 طلب
		if stats.TotalRequests%100 == 0 {
			avgLatency := stats.TotalLatency / time.Duration(stats.TotalRequests)
			
			m.logger.WithFields(logrus.Fields{
				"total_requests":          stats.TotalRequests,
				"failed_requests":         stats.FailedRequests,
				"average_latency":         avgLatency.Milliseconds(),
				"total_verifications":     stats.Verifications,
				"failed_verifications":    stats.FailedVerifications,
				"verification_success_rate": calculateSuccessRate(stats.Verifications, stats.FailedVerifications),
			}).Info("Request metrics")
		}
	}
}

func calculateSuccessRate(total, failed int64) float64 {
	if total == 0 {
		return 0
	}
	successful := total - failed
	return float64(successful) / float64(total) * 100
}