package ai

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "os"
    "strings"
    "time"
    "github.com/nawthtech/nawthtech/backend/internal/ai/types"
)

// HuggingFaceProvider مزود Hugging Face
type HuggingFaceProvider struct {
    apiToken  string
    baseURL   string
    client    *http.Client
    rateLimit *RateLimiter
}

// RateLimiter محدّد معدل الطلبات
type RateLimiter struct {
    tokens    int
    capacity  int
    refill    time.Duration
    lastRefill time.Time
    mu        chan struct{}
}

// NewRateLimiter إنشاء RateLimiter جديد
func NewRateLimiter(capacity int, refill time.Duration) *RateLimiter {
    return &RateLimiter{
        tokens:    capacity,
        capacity:  capacity,
        refill:    refill,
        lastRefill: time.Now(),
        mu:        make(chan struct{}, 1),
    }
}

// Allow التحقق من السماح بالطلب
func (r *RateLimiter) Allow() bool {
    r.mu <- struct{}{}
    defer func() { <-r.mu }()
    
    now := time.Now()
    elapsed := now.Sub(r.lastRefill)
    
    // إعادة تعبئة الرموز
    if elapsed >= r.refill {
        r.tokens = r.capacity
        r.lastRefill = now
    }
    
    if r.tokens > 0 {
        r.tokens--
        return true
    }
    
    return false
}

// NewHuggingFaceProvider إنشاء مزود Hugging Face جديد
func NewHuggingFaceProvider() *HuggingFaceProvider {
    apiToken := os.Getenv("HUGGINGFACE_TOKEN")
    
    return &HuggingFaceProvider{
        apiToken: apiToken,
        baseURL:  "https://api-inference.huggingface.co/models",
        client: &http.Client{
            Timeout: 120 * time.Second,
        },
        rateLimit: NewRateLimiter(30, time.Minute), // 30 request/دقيقة
    }
}

// GenerateText توليد نص
func (p *HuggingFaceProvider) GenerateText(req types.TextRequest) (*types.TextResponse, error) {
    if p.apiToken == "" {
        return nil, fmt.Errorf("HUGGINGFACE_TOKEN environment variable is required")
    }
    
    // انتظار إذا تجاوزنا Rate limit
    if !p.rateLimit.Allow() {
        return nil, fmt.Errorf("rate limit exceeded, please try again later")
    }
    
model := "stabilityai/stable-diffusion-xl-base-1.0" // قيمة افتراضية
    if model == "" {
        model = "google/flan-t5-xl"
    }
    
    url := fmt.Sprintf("%s/%s", p.baseURL, model)
    
    payload := map[string]interface{}{
        "inputs": req.Prompt,
        "parameters": map[string]interface{}{
            "max_new_tokens":   500,
            "temperature":      0.7,
            "top_p":           0.9,
            "do_sample":       true,
            "return_full_text": false,
        },
        "options": map[string]interface{}{
            "use_cache": true,
            "wait_for_model": true,
        },
    }
    
    jsonData, err := json.Marshal(payload)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal request: %w", err)
    }
    
    httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
    if err != nil {
        return nil, fmt.Errorf("failed to create request: %w", err)
    }
    
    httpReq.Header.Set("Authorization", "Bearer "+p.apiToken)
    httpReq.Header.Set("Content-Type", "application/json")
    
    resp, err := p.client.Do(httpReq)
    if err != nil {
        return nil, fmt.Errorf("request failed: %w", err)
    }
    defer resp.Body.Close()
    
    if resp.StatusCode == http.StatusNotFound {
        return nil, fmt.Errorf("model %s not found", model)
    }
    
    if resp.StatusCode == http.StatusServiceUnavailable {
        return nil, fmt.Errorf("model is loading, please try again in a few moments")
    }
    
    if resp.StatusCode != http.StatusOK {
        body, _ := io.ReadAll(resp.Body)
        return nil, fmt.Errorf("Hugging Face API error: %s - %s", resp.Status, string(body))
    }
    
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("failed to read response: %w", err)
    }
    
    // تحليل الاستجابة
    generatedText, err := p.parseResponse(body)
    if err != nil {
        return nil, err
    }
    
    // تقدير عدد الرموز
tokens := int(float64(len(strings.Fields(generatedText))) * 1.3)
    
    return &types.TextResponse{
        Text:        strings.TrimSpace(generatedText),
        Tokens:      int(tokens),
        Cost:        0.0, // Hugging Face مجاني
        ModelUsed:   model,
        FinishReason: "length",
        CreatedAt:   time.Now(),
    }, nil
}

// parseResponse تحليل استجابة Hugging Face
func (p *HuggingFaceProvider) parseResponse(body []byte) (string, error) {
    // محاولة تحليل كمصفوفة
    var arrayResponse []map[string]interface{}
    if err := json.Unmarshal(body, &arrayResponse); err == nil && len(arrayResponse) > 0 {
        if generatedText, ok := arrayResponse[0]["generated_text"].(string); ok {
            return generatedText, nil
        }
    }
    
    // محاولة تحليل ككائن مفرد
    var objectResponse map[string]interface{}
    if err := json.Unmarshal(body, &objectResponse); err == nil {
        if generatedText, ok := objectResponse["generated_text"].(string); ok {
            return generatedText, nil
        }
        
        // تحقق من أنواع أخرى من الاستجابات
        if text, ok := objectResponse["text"].(string); ok {
            return text, nil
        }
        
        if textArray, ok := objectResponse["text"].([]interface{}); ok && len(textArray) > 0 {
            if text, ok := textArray[0].(string); ok {
                return text, nil
            }
        }
    }
    
    // إذا كان النص مباشراً
    var directResponse []string
    if err := json.Unmarshal(body, &directResponse); err == nil && len(directResponse) > 0 {
        return directResponse[0], nil
    }
    
    // كمحاولة أخيرة، إرجاع الجسم كنص
    return string(body), nil
}

// GenerateImage توليد صورة
func (p *HuggingFaceProvider) GenerateImage(req types.ImageRequest) (*types.ImageResponse, error) {
    if p.apiToken == "" {
        return nil, fmt.Errorf("HUGGINGFACE_TOKEN environment variable is required")
    }
    
model := "stabilityai/stable-diffusion-xl-base-1.0" // قيمة افتراضية
    if model == "" {
        model = "stabilityai/stable-diffusion-xl-base-1.0"
    }
    
    url := fmt.Sprintf("%s/%s", p.baseURL, model)
    
    payload := map[string]interface{}{
        "inputs": req.Prompt,
        "parameters": map[string]interface{}{
            "negative_prompt": "blurry, low quality, distorted",
            "num_inference_steps": 25,
            "guidance_scale": 7.5,
        },
    }
    
    jsonData, err := json.Marshal(payload)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal request: %w", err)
    }
    
    httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
    if err != nil {
        return nil, fmt.Errorf("failed to create request: %w", err)
    }
    
    httpReq.Header.Set("Authorization", "Bearer "+p.apiToken)
    httpReq.Header.Set("Content-Type", "application/json")
    
    resp, err := p.client.Do(httpReq)
    if err != nil {
        return nil, fmt.Errorf("request failed: %w", err)
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        body, _ := io.ReadAll(resp.Body)
        return nil, fmt.Errorf("image generation failed: %s - %s", resp.Status, string(body))
    }
    
    imageData, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("failed to read image data: %w", err)
    }
    
    return &types.ImageResponse{
        URL:         "", // Hugging Face لا يعيد URLs للصور
        ImageData:   imageData,
        Size:        req.Size,
        Format:      "jpeg",
        Cost:        0.0,
        ModelUsed:   model,
        CreatedAt:   time.Now(),
        Seed:        0,
    }, nil
}

// GenerateVideo توليد فيديو - غير مدعوم في Hugging Face
func (p *HuggingFaceProvider) GenerateVideo(req types.VideoRequest) (*types.VideoResponse, error) {
    return nil, fmt.Errorf("video generation not supported by Hugging Face")
}

// AnalyzeText تحليل نص
func (p *HuggingFaceProvider) AnalyzeText(req types.AnalysisRequest) (*types.AnalysisResponse, error) {
    if p.apiToken == "" {
        return nil, fmt.Errorf("HUGGINGFACE_TOKEN environment variable is required")
    }
    
    // يمكن استخدام نموذج تحليل المشاعر
    model := "cardiffnlp/twitter-roberta-base-sentiment-latest"
    
    url := fmt.Sprintf("%s/%s", p.baseURL, model)
    
    payload := map[string]interface{}{
        "inputs": req.Text,
    }
    
    jsonData, err := json.Marshal(payload)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal request: %w", err)
    }
    
    httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
    if err != nil {
        return nil, fmt.Errorf("failed to create request: %w", err)
    }
    
    httpReq.Header.Set("Authorization", "Bearer "+p.apiToken)
    httpReq.Header.Set("Content-Type", "application/json")
    
    resp, err := p.client.Do(httpReq)
    if err != nil {
        return nil, fmt.Errorf("request failed: %w", err)
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        body, _ := io.ReadAll(resp.Body)
        return nil, fmt.Errorf("analysis failed: %s - %s", resp.Status, string(body))
    }
    
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("failed to read response: %w", err)
    }
    
    // تحليل استجابة المشاعر
    var sentimentResult []map[string]interface{}
    if err := json.Unmarshal(body, &sentimentResult); err != nil {
        return nil, fmt.Errorf("failed to parse sentiment response: %w", err)
    }
    
    var result string
    var confidence float64
    
    if len(sentimentResult) > 0 {
        if labels, ok := sentimentResult[0]["label"].(string); ok {
            result = fmt.Sprintf("Sentiment: %s", labels)
        }
        if score, ok := sentimentResult[0]["score"].(float64); ok {
            confidence = score
        }
    }
    
    return &types.AnalysisResponse{
        Result:     result,
        Confidence: confidence,
        Cost:       0.0,
        Model:      model,
        CreatedAt:  time.Now(),
    }, nil
}

// AnalyzeImage تحليل صورة
func (p *HuggingFaceProvider) AnalyzeImage(req types.AnalysisRequest) (*types.AnalysisResponse, error) {
    return nil, fmt.Errorf("image analysis not supported by Hugging Face")
}

// TranslateText ترجمة نص
func (p *HuggingFaceProvider) TranslateText(req types.TranslationRequest) (*types.TranslationResponse, error) {
    if p.apiToken == "" {
        return nil, fmt.Errorf("HUGGINGFACE_TOKEN environment variable is required")
    }
    
    // بناء اسم النموذج بناءً على اللغات
    model := fmt.Sprintf("Helsinki-NLP/opus-mt-%s-%s", req.FromLang, req.ToLang)
    
    url := fmt.Sprintf("%s/%s", p.baseURL, model)
    
    payload := map[string]interface{}{
        "inputs": req.Text,
    }
    
    jsonData, err := json.Marshal(payload)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal request: %w", err)
    }
    
    httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
    if err != nil {
        return nil, fmt.Errorf("failed to create request: %w", err)
    }
    
    httpReq.Header.Set("Authorization", "Bearer "+p.apiToken)
    httpReq.Header.Set("Content-Type", "application/json")
    
    resp, err := p.client.Do(httpReq)
    if err != nil {
        return nil, fmt.Errorf("request failed: %w", err)
    }
    defer resp.Body.Close()
    
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("failed to read response: %w", err)
    }
    
    var result []map[string]interface{}
    if err := json.Unmarshal(body, &result); err != nil {
        return nil, fmt.Errorf("failed to parse translation response: %w", err)
    }
    
    var translatedText string
    if len(result) > 0 {
        if translation, ok := result[0]["translation_text"].(string); ok {
            translatedText = translation
        }
    }
    
    if translatedText == "" {
        return nil, fmt.Errorf("translation failed")
    }
    
    return &types.TranslationResponse{
        TranslatedText: translatedText,
        Cost:           0.0,
        Model:          model,
        CreatedAt:      time.Now(),
    }, nil
}

// Transcribe تحويل صوت إلى نص
func (p *HuggingFaceProvider) Transcribe(audioData []byte) (string, error) {
    if p.apiToken == "" {
        return "", fmt.Errorf("HUGGINGFACE_TOKEN environment variable is required")
    }
    
    url := p.baseURL + "/openai/whisper-large-v3"
    
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(audioData))
    if err != nil {
        return "", fmt.Errorf("failed to create request: %w", err)
    }
    
    req.Header.Set("Authorization", "Bearer "+p.apiToken)
    req.Header.Set("Content-Type", "audio/wav")
    
    resp, err := p.client.Do(req)
    if err != nil {
        return "", fmt.Errorf("request failed: %w", err)
    }
    defer resp.Body.Close()
    
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return "", fmt.Errorf("failed to read response: %w", err)
    }
    
    var result struct {
        Text string `json:"text"`
    }
    
    if err := json.Unmarshal(body, &result); err != nil {
        return "", fmt.Errorf("failed to parse transcription response: %w", err)
    }
    
    return result.Text, nil
}

// Summarize تلخيص نص
func (p *HuggingFaceProvider) Summarize(text string, maxLength int) (string, error) {
    if p.apiToken == "" {
        return "", fmt.Errorf("HUGGINGFACE_TOKEN environment variable is required")
    }
    
    url := p.baseURL + "/facebook/bart-large-cnn"
    
    payload := map[string]interface{}{
        "inputs": text,
        "parameters": map[string]interface{}{
            "max_length": maxLength,
            "min_length": 30,
            "do_sample": false,
        },
    }
    
    jsonData, err := json.Marshal(payload)
    if err != nil {
        return "", fmt.Errorf("failed to marshal request: %w", err)
    }
    
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
    if err != nil {
        return "", fmt.Errorf("failed to create request: %w", err)
    }
    
    req.Header.Set("Authorization", "Bearer "+p.apiToken)
    req.Header.Set("Content-Type", "application/json")
    
    resp, err := p.client.Do(req)
    if err != nil {
        return "", fmt.Errorf("request failed: %w", err)
    }
    defer resp.Body.Close()
    
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return "", fmt.Errorf("failed to read response: %w", err)
    }
    
    var result []map[string]interface{}
    if err := json.Unmarshal(body, &result); err != nil {
        return "", fmt.Errorf("failed to parse summary response: %w", err)
    }
    
    if len(result) > 0 {
        if summary, ok := result[0]["summary_text"].(string); ok {
            return summary, nil
        }
    }
    
    return "", fmt.Errorf("summarization failed")
}

// GetAvailableModels الحصول على النماذج المتاحة
func (p *HuggingFaceProvider) GetAvailableModels() []string {
    return []string{
        "google/flan-t5-xl",
        "mistralai/Mistral-7B-Instruct-v0.2",
        "Qwen/Qwen2.5-7B-Instruct",
        "microsoft/phi-2",
        "stabilityai/stable-diffusion-xl-base-1.0",
        "openai/whisper-large-v3",
        "facebook/bart-large-cnn",
        "Helsinki-NLP/opus-mt-ar-en",
        "Helsinki-NLP/opus-mt-en-ar",
        "cardiffnlp/twitter-roberta-base-sentiment-latest",
    }
}

// GetModelInfo الحصول على معلومات النموذج
func (p *HuggingFaceProvider) GetModelInfo(model string) (map[string]interface{}, error) {
    url := fmt.Sprintf("https://huggingface.co/api/models/%s", model)
    
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to create request: %w", err)
    }
    
    resp, err := p.client.Do(req)
    if err != nil {
        return nil, fmt.Errorf("request failed: %w", err)
    }
    defer resp.Body.Close()
    
    var info map[string]interface{}
    if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
        return nil, fmt.Errorf("failed to parse model info: %w", err)
    }
    
    return info, nil
}

// IsAvailable التحقق من التوفر
func (p *HuggingFaceProvider) IsAvailable() bool {
    if p.apiToken == "" {
        return false
    }
    
    // اختبار بسيط للاتصال
    resp, err := p.client.Get("https://huggingface.co/api/whoami")
    return err == nil && resp.StatusCode == http.StatusOK
}

// GetName اسم المزود
func (p *HuggingFaceProvider) GetName() string {
    return "Hugging Face"
}

// GetCost التكلفة (مجاني محدود)
func (p *HuggingFaceProvider) GetCost() float64 {
    return 0.0
}

// GetStats الحصول على إحصائيات
func (p *HuggingFaceProvider) GetStats() *types.ProviderStats {
    return &types.ProviderStats{
        Name:        p.GetName(),
        Type:        "text",
        IsAvailable: p.IsAvailable(),
        Requests:    0,
        Successful:  0,
        Failed:      0,
        TotalCost:   0.0,
        AvgLatency:  0.0,
        LastUsed:    time.Time{},
        SuccessRate: 85.0,
    }
}

// GetType نوع المزود
func (p *HuggingFaceProvider) GetType() string {
    return "text"
}

// SupportsStreaming يدعم التدفق
func (p *HuggingFaceProvider) SupportsStreaming() bool {
    return false
}

// SupportsEmbedding يدعم التضمين
func (p *HuggingFaceProvider) SupportsEmbedding() bool {
    return false
}

// GetMaxTokens الحد الأقصى للرموز
func (p *HuggingFaceProvider) GetMaxTokens() int {
    return 500
}

// GetSupportedLanguages اللغات المدعومة
func (p *HuggingFaceProvider) GetSupportedLanguages() []string {
    return []string{"ar", "en", "es", "fr", "de"}
}