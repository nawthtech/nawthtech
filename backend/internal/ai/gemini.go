package ai

import (
    "context"
    "fmt"
    "log"
    "os"
    "strings"
    
    "google.golang.org/genai"
)

// GeminiProvider تنفيذ باستخدام الحزمة الرسمية
type GeminiProvider struct {
    client *genai.Client
    apiKey string
}

// NewGeminiProvider إنشاء مزود Gemini جديد
func NewGeminiProvider() (*GeminiProvider, error) {
    apiKey := os.Getenv("GEMINI_API_KEY")
    if apiKey == "" {
        return nil, fmt.Errorf("GEMINI_API_KEY environment variable is required")
    }
    
    ctx := context.Background()
    
    // إنشاء client بالطريقة الصحيحة
    client, err := genai.NewClient(ctx, &genai.ClientConfig{
        APIKey: apiKey,
    })
    if err != nil {
        return nil, fmt.Errorf("failed to create Gemini client: %w", err)
    }
    
    return &GeminiProvider{
        client: client,
        apiKey: apiKey,
    }, nil
}

// GenerateText توليد نص باستخدام Gemini
func (p *GeminiProvider) GenerateText(prompt string, opts ...Option) (string, error) {
    ctx := context.Background()
    
    // تحديد النموذج (flash مجاني، pro للجودة الأعلى)
    modelName := "gemini-2.5-flash"
    for _, opt := range opts {
        if opt.Model != "" {
            modelName = opt.Model
        }
    }
    
    log.Printf("Using Gemini model: %s", modelName)
    
    // توليد المحتوى كما في الكود الأصلي
    result, err := p.client.Models.GenerateContent(
        ctx,
        modelName,
        genai.Text(prompt),
        nil, // options إضافية
    )
    if err != nil {
        return "", fmt.Errorf("Gemini generation failed: %w", err)
    }
    
    // استخراج النص
    return extractTextFromResult(result), nil
}

// GenerateImage توليد صور باستخدام Gemini Vision
func (p *GeminiProvider) GenerateImage(prompt string, opts ...Option) ([]byte, error) {
    ctx := context.Background()
    
    // استخدام نموذج Vision
    result, err := p.client.Models.GenerateContent(
        ctx,
        "gemini-2.5-flash-image", // النموذج المخصص للصور
        genai.Text(prompt),
        nil,
    )
    if err != nil {
        return nil, fmt.Errorf("Gemini image generation failed: %w", err)
    }
    
    // استخراج بيانات الصورة
    return extractImageFromResult(result), nil
}

// AnalyzeImage تحليل صور
func (p *GeminiProvider) AnalyzeImage(imageData []byte, prompt string) (string, error) {
    ctx := context.Background()
    
    // إنشاء parts مع الصورة والنص
    parts := []genai.Part{
        genai.ImageData{
            MIMEType: "image/png",
            Data:     imageData,
        },
        genai.Text(prompt),
    }
    
    result, err := p.client.Models.GenerateContent(
        ctx,
        "gemini-2.5-flash", // أو gemini-2.5-pro-vision
        parts...,
    )
    if err != nil {
        return "", fmt.Errorf("Gemini image analysis failed: %w", err)
    }
    
    return extractTextFromResult(result), nil
}

// ===== دوال مساعدة =====

// extractTextFromResult استخراج النص من النتيجة
func extractTextFromResult(result *genai.GenerateContentResponse) string {
    if result == nil || len(result.Candidates) == 0 {
        return ""
    }
    
    var textParts []string
    for _, part := range result.Candidates[0].Content.Parts {
        if part.Text != "" {
            textParts = append(textParts, part.Text)
        }
    }
    
    return strings.Join(textParts, "\n")
}

// extractImageFromResult استخراج الصورة من النتيجة
func extractImageFromResult(result *genai.GenerateContentResponse) []byte {
    if result == nil || len(result.Candidates) == 0 {
        return nil
    }
    
    for _, part := range result.Candidates[0].Content.Parts {
        if part.InlineData != nil {
            return part.InlineData.Data
        }
    }
    
    return nil
}

// IsAvailable التحقق من التوفر
func (p *GeminiProvider) IsAvailable() bool {
    return p.client != nil
}

// GetCost التكلفة (Gemini flash مجاني)
func (p *GeminiProvider) GetCost() float64 {
    return 0.0
}

// GetName اسم المزود
func (p *GeminiProvider) GetName() string {
    return "Google Gemini (Official)"
}