// backend/internal/ai/verification/examples.go
package verification

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/nawthtech/nawthtech/backend/internal/config"
)

// ExampleBasicUsage مثال للاستخدام الأساسي
func ExampleBasicUsage() {
	// إعداد التكوين
	cfg := &config.Config{
		AI: config.AIConfig{
			OpenAIKey:     "your-openai-api-key",
			DefaultModel:  "gpt-4o-mini",
			MaxRetries:    3,
			TimeoutMS:     30000,
			Temperature:   0.7,
			MaxTokens:     1000,
		},
		Verification: config.VerificationConfig{
			CheckToxicity:   true,
			CheckFactuality: true,
			CheckSafety:     true,
			CheckCoherence:  true,
			CheckRelevance:  true,
		},
		Sentry: config.SentryConfig{
			DSN:         "your-sentry-dsn",
			Environment: "development",
		},
	}
	
	// إنشاء المسجل
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)
	
	// إنشاء المحقق
	verifier, err := NewVerifier(cfg, logger)
	if err != nil {
		log.Fatal("Failed to create verifier:", err)
	}
	
	// نص للتحقق
	content := "This is a sample content that needs verification."
	
	// التحقق من المحتوى
	ctx := context.Background()
	result, err := verifier.Verify(ctx, content)
	if err != nil {
		log.Fatal("Verification failed:", err)
	}
	
	// عرض النتيجة
	fmt.Printf("Is Valid: %v\n", result.IsValid)
	fmt.Printf("Confidence: %.2f%%\n", result.Confidence*100)
	fmt.Printf("Reason: %s\n", result.Reason)
	
	if len(result.Issues) > 0 {
		fmt.Println("Issues:")
		for _, issue := range result.Issues {
			fmt.Printf("  - %s\n", issue)
		}
	}
	
	if len(result.Suggestions) > 0 {
		fmt.Println("Suggestions:")
		for _, suggestion := range result.Suggestions {
			fmt.Printf("  - %s\n", suggestion)
		}
	}
	
	// عرض المقاييس
	fmt.Printf("Latency: %dms\n", result.Metrics.Latency)
	fmt.Printf("Model: %s\n", result.Metrics.Model)
}

// ExampleBatchVerification مثال للتحقق الدفعي
func ExampleBatchVerification() {
	cfg := &config.Config{
		AI: config.AIConfig{
			OpenAIKey:    "your-openai-api-key",
			DefaultModel: "gpt-4o-mini",
		},
	}
	
	logger := logrus.New()
	verifier, _ := NewVerifier(cfg, logger)
	
	// قائمة المحتويات للتحقق
	contents := []string{
		"First content to verify",
		"Second content that might have issues",
		"Third content with potential safety concerns",
	}
	
	// التحقق الدفعي
	ctx := context.Background()
	results, err := verifier.BatchVerify(ctx, contents)
	if err != nil {
		log.Fatal("Batch verification failed:", err)
	}
	
	// تحليل النتائج
	validCount := 0
	totalConfidence := 0.0
	
	for i, result := range results {
		fmt.Printf("\nResult %d:\n", i+1)
		fmt.Printf("  Valid: %v\n", result.IsValid)
		fmt.Printf("  Confidence: %.2f%%\n", result.Confidence*100)
		
		if result.IsValid {
			validCount++
		}
		totalConfidence += result.Confidence
		
		if len(result.Issues) > 0 {
			fmt.Printf("  Issues: %v\n", result.Issues)
		}
	}
	
	// إحصائيات الدفعة
	fmt.Printf("\nBatch Summary:\n")
	fmt.Printf("  Total: %d\n", len(results))
	fmt.Printf("  Valid: %d (%.1f%%)\n", validCount, float64(validCount)/float64(len(results))*100)
	fmt.Printf("  Average Confidence: %.2f%%\n", totalConfidence/float64(len(results))*100)
}

// ExampleWithOptions مثال مع خيارات مخصصة
func ExampleWithOptions() {
	cfg := &config.Config{
		AI: config.AIConfig{
			OpenAIKey:    "your-openai-api-key",
			DefaultModel: "gpt-4o-mini",
		},
	}
	
	logger := logrus.New()
	verifier, _ := NewVerifier(cfg, logger)
	
	content := "Content to verify with specific criteria"
	
	// استخدام خيارات مخصصة
	ctx := context.Background()
	result, err := verifier.Verify(ctx, content,
		WithModel("gpt-4"),
		WithTemperature(0.3), // أكثر تحديداً
		WithMaxTokens(500),
		WithType("fact_check"),
		WithContext("This is educational content"),
	)
	
	if err != nil {
		log.Fatal("Verification failed:", err)
	}
	
	fmt.Println("Verification with custom options:")
	fmt.Printf("Model: %s\n", result.Metrics.Model)
	fmt.Printf("Type: fact_check\n")
	fmt.Printf("Result: %v (%.1f%% confidence)\n", result.IsValid, result.Confidence*100)
}

// ExampleMiddlewareIntegration مثال تكامل الـ Middleware
func ExampleMiddlewareIntegration() {
	cfg := &config.Config{
		AI: config.AIConfig{
			OpenAIKey: "your-openai-api-key",
		},
		Verification: config.VerificationConfig{
			CheckToxicity: true,
			CheckSafety:   true,
		},
	}
	
	logger := logrus.New()
	verifier, _ := NewVerifier(cfg, logger)
	
	// إنشاء middleware
	middleware := NewVerificationMiddleware(verifier, logger, true)
	
	// معايير التحقق
	criteria := VerificationCriteria{
		Toxicity: true,
		Safety:   true,