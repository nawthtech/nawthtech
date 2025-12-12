// backend/internal/ai/verification/utils.go
package verification

import (
	"encoding/json"
	"fmt"
	"math"
	"regexp"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// BuildPrompt بناء prompt التحقق
func BuildPrompt(input string, criteria VerificationCriteria, context string) string {
	var prompt strings.Builder

	prompt.WriteString("Please verify the following content and provide a structured JSON response.\n\n")
	
	if context != "" {
		prompt.WriteString(fmt.Sprintf("Context: %s\n\n", context))
	}
	
	prompt.WriteString(fmt.Sprintf("Content to verify: \"%s\"\n\n", input))
	prompt.WriteString("Verification Criteria (check all that apply):\n")
	
	if criteria.Toxicity {
		prompt.WriteString("- Toxicity: Check for toxic, hateful, harmful, or abusive content\n")
	}
	if criteria.Factuality {
		prompt.WriteString("- Factuality: Verify factual accuracy, check for misinformation\n")
	}
	if criteria.Coherence {
		prompt.WriteString("- Coherence: Check logical flow, consistency, and clarity\n")
	}
	if criteria.Relevance {
		prompt.WriteString("- Relevance: Check if content is relevant to the intended topic\n")
	}
	if criteria.Safety {
		prompt.WriteString("- Safety: Check for dangerous, illegal, or policy-violating content\n")
	}
	if criteria.Moderation {
		prompt.WriteString("- Moderation: Check for inappropriate, explicit, or offensive content\n")
	}
	if criteria.Bias {
		prompt.WriteString("- Bias: Check for political, racial, gender, or other biases\n")
	}
	
	prompt.WriteString("\nRespond with a JSON object in this exact format:\n")
	prompt.WriteString(`{
  "isValid": boolean,
  "confidence": number between 0 and 1,
  "reason": "brief explanation",
  "issues": ["specific issue 1", "specific issue 2"],
  "suggestions": ["suggestion 1", "suggestion 2"],
  "categories": {
    "toxicity": {"passed": boolean, "score": number, "explanation": "details"},
    "factuality": {"passed": boolean, "score": number, "explanation": "details"},
    "coherence": {"passed": boolean, "score": number, "explanation": "details"},
    "relevance": {"passed": boolean, "score": number, "explanation": "details"},
    "safety": {"passed": boolean, "score": number, "explanation": "details"},
    "moderation": {"passed": boolean, "score": number, "explanation": "details"},
    "bias": {"passed": boolean, "score": number, "explanation": "details"}
  }
}`)
	
	return prompt.String()
}

// ExtractJSONFromText استخراج JSON من النص
func ExtractJSONFromText(text string) (map[string]interface{}, error) {
	// محاولة تحليل مباشر
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(text), &data); err == nil {
		return data, nil
	}

	// محاولة استخراج كائن JSON
	jsonMatch := regexp.MustCompile(`\{[\s\S]*\}`).FindString(text)
	if jsonMatch == "" {
		return nil, fmt.Errorf("no JSON found in text")
	}

	// إصلاح مشاكل JSON الشائعة
	fixedJSON := FixJSONString(jsonMatch)
	
	if err := json.Unmarshal([]byte(fixedJSON), &data); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}

	return data, nil
}

// FixJSONString إصلاح سلسلة JSON
func FixJSONString(jsonStr string) string {
	// إصلاح مفاتيح بدون علامات اقتباس
	re := regexp.MustCompile(`(['"])?([a-zA-Z0-9_]+)(['"])?:`)
	fixed := re.ReplaceAllString(jsonStr, `"$2":`)
	
	// إزالة فواصل زائدة
	fixed = regexp.MustCompile(`,\s*}`).ReplaceAllString(fixed, "}")
	fixed = regexp.MustCompile(`,\s*]`).ReplaceAllString(fixed, "]")
	
	// إصلاح علامات الاقتباس الهاربة
	fixed = strings.ReplaceAll(fixed, `\'`, "'")
	fixed = strings.ReplaceAll(fixed, `\"`, `"`)
	
	return fixed
}

// CalculateConfidence حساب الثقة
func CalculateConfidence(score interface{}) float64 {
	switch v := score.(type) {
	case float64:
		return NormalizeScore(v)
	case int:
		return NormalizeScore(float64(v))
	case string:
		var num float64
		if _, err := fmt.Sscanf(v, "%f", &num); err == nil {
			return NormalizeScore(num)
		}
	}
	return 0.5
}

// NormalizeScore تطبيع النقاط
func NormalizeScore(score float64) float64 {
	if score <= 0 {
		return 0
	}
	if score >= 1 {
		return 1
	}
	if score > 100 {
		return score / 100
	}
	if score > 10 {
		return score / 10
	}
	if score > 5 {
		return score / 5
	}
	return score
}

// ParseUnstructuredResponse تحليل استجابة غير مهيكلة
func ParseUnstructuredResponse(text, input string) *VerificationResult {
	result := &VerificationResult{
		IsValid:    false,
		Confidence: 0.5,
		Reason:     "Automated analysis",
		Issues:     []string{},
		Suggestions: []string{},
		Categories: make(map[string]CategoryResult),
		Metrics: VerificationMetrics{
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		},
	}

	textLower := strings.ToLower(text)
	
	// تحديد الصلاحية
	positiveWords := []string{"valid", "passed", "ok", "good", "safe", "appropriate", "acceptable"}
	negativeWords := []string{"invalid", "failed", "bad", "unsafe", "inappropriate", "toxic", "reject"}
	
	hasPositive := false
	hasNegative := false
	
	for _, word := range positiveWords {
		if strings.Contains(textLower, word) {
			hasPositive = true
			break
		}
	}
	
	for _, word := range negativeWords {
		if strings.Contains(textLower, word) {
			hasNegative = true
			break
		}
	}
	
	result.IsValid = hasPositive && !hasNegative
	
	// استخراج الثقة
	confidence := ExtractConfidenceFromText(text)
	result.Confidence = confidence
	
	// استخراج السبب
	result.Reason = ExtractReasonFromText(text)
	
	// استخراج المشاكل
	result.Issues = ExtractListFromText(text, []string{"issues?", "problems?", "concerns?"})
	if len(result.Issues) == 0 {
		result.Issues = []string{"No specific issues mentioned"}
	}
	
	// استخراج الاقتراحات
	result.Suggestions = ExtractListFromText(text, []string{"suggestions?", "recommendations?", "improvements?"})
	if len(result.Suggestions) == 0 {
		result.Suggestions = []string{"No suggestions provided"}
	}
	
	return result
}

// ExtractConfidenceFromText استخراج الثقة من النص
func ExtractConfidenceFromText(text string) float64 {
	re := regexp.MustCompile(`confidence.*?(\d+\.?\d*)`)
	matches := re.FindStringSubmatch(strings.ToLower(text))
	
	if len(matches) > 1 {
		var confidence float64
		if _, err := fmt.Sscanf(matches[1], "%f", &confidence); err == nil {
			return NormalizeScore(confidence)
		}
	}
	
	// تقدير الثقة بناءً على اللغة
	confidentWords := []string{"definitely", "certainly", "clearly", "obviously", "undoubtedly"}
	uncertainWords := []string{"maybe", "perhaps", "possibly", "likely", "probably"}
	
	textLower := strings.ToLower(text)
	confidentCount := 0
	uncertainCount := 0
	
	for _, word := range confidentWords {
		if strings.Contains(textLower, word) {
			confidentCount++
		}
	}
	
	for _, word := range uncertainWords {
		if strings.Contains(textLower, word) {
			uncertainCount++
		}
	}
	
	if confidentCount > uncertainCount {
		return 0.8
	}
	if uncertainCount > confidentCount {
		return 0.4
	}
	return 0.6
}

// ExtractReasonFromText استخراج السبب من النص
func ExtractReasonFromText(text string) string {
	re := regexp.MustCompile(`reason[:\s]+([^.\n]+)`)
	matches := re.FindStringSubmatch(strings.ToLower(text))
	
	if len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}
	
	// النسخ الاحتياطي: الجملة الأولى
	sentences := strings.Split(text, ".")
	if len(sentences) > 0 && len(sentences[0]) > 10 {
		return strings.TrimSpace(sentences[0])
	}
	
	return "Automated analysis completed"
}

// ExtractListFromText استخراج قائمة من النص
func ExtractListFromText(text string, patterns []string) []string {
	items := []string{}
	
	for _, pattern := range patterns {
		re := regexp.MustCompile(fmt.Sprintf(`%s[:\s]+([^.\n]+)`, pattern))
		matches := re.FindAllStringSubmatch(text, -1)
		
		for _, match := range matches {
			if len(match) > 1 {
				// تقسيم بالفاصلة أو الفاصلة المنقوطة أو النقاط
				parts := regexp.MustCompile(`[,;•\-]\s*`).Split(match[1], -1)
				for _, part := range parts {
					trimmed := strings.TrimSpace(part)
					if trimmed != "" {
						items = append(items, trimmed)
					}
				}
			}
		}
	}
	
	// إزالة التكرارات
	return removeDuplicates(items)
}

func removeDuplicates(items []string) []string {
	seen := make(map[string]bool)
	result := []string{}
	
	for _, item := range items {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}
	
	return result
}

// CalculateCost حساب التكلفة
func CalculateCost(inputTokens, outputTokens int, costPer1KInput, costPer1KOutput float64) float64 {
	inputCost := (float64(inputTokens) / 1000) * costPer1KInput
	outputCost := (float64(outputTokens) / 1000) * costPer1KOutput
	return inputCost + outputCost
}

// FormatDuration تنسيق المدة
func FormatDuration(d time.Duration) string {
	if d < time.Millisecond {
		return fmt.Sprintf("%dµs", d.Microseconds())
	}
	if d < time.Second {
		return fmt.Sprintf("%.2fms", float64(d.Milliseconds()))
	}
	return fmt.Sprintf("%.2fs", d.Seconds())
}

// LogVerification تسجيل التحقق
func LogVerification(result *VerificationResult, logger *logrus.Logger, fields logrus.Fields) {
	logFields := logrus.Fields{
		"isValid":    result.IsValid,
		"confidence": result.Confidence,
		"issues":     len(result.Issues),
	}
	
	for k, v := range fields {
		logFields[k] = v
	}
	
	if result.Metrics.Latency > 0 {
		logFields["latency"] = result.Metrics.Latency
	}
	if result.Metrics.TotalTokens > 0 {
		logFields["tokens"] = result.Metrics.TotalTokens
	}
	if result.Metrics.Cost > 0 {
		logFields["cost"] = result.Metrics.Cost
	}
	
	if result.IsValid {
		logger.WithFields(logFields).Info("Verification passed")
	} else {
		logger.WithFields(logFields).Warn("Verification failed")
		
		if len(result.Issues) > 0 {
			logger.WithFields(logFields).Warnf("Issues: %v", result.Issues)
		}
	}
}

// ValidateInput التحقق من صحة الإدخال
func ValidateInput(input string, maxLength int) error {
	if input == "" {
		return fmt.Errorf("input cannot be empty")
	}
	
	if len(input) > maxLength {
		return fmt.Errorf("input exceeds maximum length of %d characters", maxLength)
	}
	
	// التحقق من الأحرف غير المطبوعة
	for i, r := range input {
		if r == 0 {
			return fmt.Errorf("input contains null character at position %d", i)
		}
	}
	
	return nil
}

// CalculateBatchStats حساب إحصائيات الدفعة
func CalculateBatchStats(results []*VerificationResult) *BatchVerificationResult {
	stats := &BatchVerificationResult{
		Total:   len(results),
		Valid:   0,
		Invalid: 0,
		Summary: make(map[string]int),
		Results: results,
	}
	
	totalConfidence := 0.0
	totalCost := 0.0
	totalTokens := 0
	
	for _, result := range results {
		if result.IsValid {
			stats.Valid++
		} else {
			stats.Invalid++
		}
		
		totalConfidence += result.Confidence
		totalCost += result.Metrics.Cost
		totalTokens += result.Metrics.TotalTokens
		
		// تعداد المشاكل
		for _, issue := range result.Issues {
			stats.Summary[issue]++
		}
	}
	
	if stats.Total > 0 {
		stats.AverageConfidence = totalConfidence / float64(stats.Total)
	}
	
	stats.TotalCost = totalCost
	stats.TotalTokens = totalTokens
	
	return stats
}

// IsRetryableError التحقق إذا كان الخطأ قابلاً لإعادة المحاولة
func IsRetryableError(err error) bool {
	errStr := strings.ToLower(err.Error())
	
	retryablePatterns := []string{
		"timeout",
		"connection",
		"network",
		"rate limit",
		"too many requests",
		"server error",
		"503",
		"504",
		"temporary",
		"busy",
	}
	
	for _, pattern := range retryablePatterns {
		if strings.Contains(errStr, pattern) {
			return true
		}
	}
	
	return false
}

// ExponentialBackoff حساب الانتظار الأسي
func ExponentialBackoff(retryCount int, baseDelay time.Duration) time.Duration {
	delay := float64(baseDelay) * math.Pow(2, float64(retryCount))
	maxDelay := 30 * time.Second
	
	if delay > float64(maxDelay) {
		return maxDelay
	}
	
	return time.Duration(delay)
}

// SanitizeInput تنظيف الإدخال
func SanitizeInput(input string) string {
	// إزالة الأحرف غير المطبوعة باستثناء المسافات وعلامات الترقيم
	re := regexp.MustCompile(`[\x00-\x08\x0B\x0C\x0E-\x1F\x7F]`)
	sanitized := re.ReplaceAllString(input, "")
	
	// قص المسافات الزائدة
	sanitized = strings.TrimSpace(sanitized)
	
	// استبدال المسافات المتعددة بمسافة واحدة
	re = regexp.MustCompile(`\s+`)
	sanitized = re.ReplaceAllString(sanitized, " ")
	
	return sanitized
}