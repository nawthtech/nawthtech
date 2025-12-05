package ai

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "os"
    "strings"
    "time"
)

// HuggingFaceProvider مزود Hugging Face
type HuggingFaceProvider struct {
    apiToken  string
    baseURL   string
    client    *http.Client
    rateLimit *RateLimiter
}

// NewHuggingFaceProvider إنشاء مزود Hugging Face جديد
func NewHuggingFaceProvider() (*HuggingFaceProvider, error) {
    apiToken := os.Getenv("HUGGINGFACE_TOKEN")
    if apiToken == "" {
        return nil, fmt.Errorf("HUGGINGFACE_TOKEN environment variable is required")
    }
    
    return &HuggingFaceProvider{
        apiToken: apiToken,
        baseURL:  "https://api-inference.huggingface.co/models",
        client: &http.Client{
            Timeout: 120 * time.Second,
        },
        rateLimit: NewRateLimiter(30, time.Minute), // 30 request/دقيقة
    }, nil
}

// Generate توليد نص
func (p *HuggingFaceProvider) Generate(prompt string, options ...Option) (string, error) {
    opts := &Options{
        Model:       "google/flan-t5-xl",
        Temperature: 0.7,
        MaxTokens:   500,
    }
    
    for _, opt := range options {
        opt(opts)
    }
    
    // انتظار إذا تجاوزنا Rate limit
    if !p.rateLimit.Allow() {
        return "", fmt.Errorf("rate limit exceeded, please try again later")
    }
    
    return p.generateText(prompt, opts)
}

// generateText توليد نص باستخدام Hugging Face
func (p *HuggingFaceProvider) generateText(prompt string, opts *Options) (string, error) {
    url := fmt.Sprintf("%s/%s", p.baseURL, opts.Model)
    
    payload := map[string]interface{}{
        "inputs": prompt,
        "parameters": map[string]interface{}{
            "max_new_tokens":   opts.MaxTokens,
            "temperature":      opts.Temperature,
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
    
    if resp.StatusCode == http.StatusNotFound {
        return "", fmt.Errorf("model %s not found", opts.Model)
    }
    
    if resp.StatusCode == http.StatusServiceUnavailable {
        return "", fmt.Errorf("model is loading, please try again in a few moments")
    }
    
    if resp.StatusCode != http.StatusOK {
        body, _ := io.ReadAll(resp.Body)
        return "", fmt.Errorf("Hugging Face API error: %s - %s", resp.Status, string(body))
    }
    
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return "", fmt.Errorf("failed to read response: %w", err)
    }
    
    // تحليل الاستجابة بناءً على نوع النموذج
    return p.parseResponse(body)
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
func (p *HuggingFaceProvider) GenerateImage(prompt string, options ...Option) ([]byte, error) {
    opts := &Options{
        Model: "stabilityai/stable-diffusion-xl-base-1.0",
    }
    
    for _, opt := range options {
        opt(opts)
    }
    
    url := fmt.Sprintf("%s/%s", p.baseURL, opts.Model)
    
    payload := map[string]interface{}{
        "inputs": prompt,
        "parameters": map[string]interface{}{
            "negative_prompt": "blurry, low quality, distorted",
            "num_inference_steps": 25,
            "guidance_scale": 7.5,
        },
    }
    
    jsonData, err := json.Marshal(payload)
    if err != nil {
        return nil, err
    }
    
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
    if err != nil {
        return nil, err
    }
    
    req.Header.Set("Authorization", "Bearer "+p.apiToken)
    req.Header.Set("Content-Type", "application/json")
    
    resp, err := p.client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        body, _ := io.ReadAll(resp.Body)
        return nil, fmt.Errorf("image generation failed: %s - %s", resp.Status, string(body))
    }
    
    return io.ReadAll(resp.Body)
}

// Transcribe تحويل صوت إلى نص
func (p *HuggingFaceProvider) Transcribe(audioData []byte) (string, error) {
    url := p.baseURL + "/openai/whisper-large-v3"
    
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(audioData))
    if err != nil {
        return "", err
    }
    
    req.Header.Set("Authorization", "Bearer "+p.apiToken)
    req.Header.Set("Content-Type", "audio/wav")
    
    resp, err := p.client.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()
    
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return "", err
    }
    
    var result struct {
        Text string `json:"text"`
    }
    
    if err := json.Unmarshal(body, &result); err != nil {
        return "", err
    }
    
    return result.Text, nil
}

// Translate ترجمة نص
func (p *HuggingFaceProvider) Translate(text, sourceLang, targetLang string) (string, error) {
    model := fmt.Sprintf("Helsinki-NLP/opus-mt-%s-%s", sourceLang, targetLang)
    url := fmt.Sprintf("%s/%s", p.baseURL, model)
    
    payload := map[string]interface{}{
        "inputs": text,
    }
    
    jsonData, err := json.Marshal(payload)
    if err != nil {
        return "", err
    }
    
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
    if err != nil {
        return "", err
    }
    
    req.Header.Set("Authorization", "Bearer "+p.apiToken)
    req.Header.Set("Content-Type", "application/json")
    
    resp, err := p.client.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()
    
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return "", err
    }
    
    var result []map[string]interface{}
    if err := json.Unmarshal(body, &result); err != nil {
        return "", err
    }
    
    if len(result) > 0 {
        if translation, ok := result[0]["translation_text"].(string); ok {
            return translation, nil
        }
    }
    
    return "", fmt.Errorf("translation failed")
}

// Summarize تلخيص نص
func (p *HuggingFaceProvider) Summarize(text string, maxLength int) (string, error) {
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
        return "", err
    }
    
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
    if err != nil {
        return "", err
    }
    
    req.Header.Set("Authorization", "Bearer "+p.apiToken)
    req.Header.Set("Content-Type", "application/json")
    
    resp, err := p.client.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()
    
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return "", err
    }
    
    var result []map[string]interface{}
    if err := json.Unmarshal(body, &result); err != nil {
        return "", err
    }
    
    if len(result) > 0 {
        if summary, ok := result[0]["summary_text"].(string); ok {
            return summary, nil
        }
    }
    
    return "", fmt.Errorf("summarization failed")
}

// GetAvailableModels الحصول على النماذج المتاحة
func (p *HuggingFaceProvider) GetAvailableModels() ([]string, error) {
    // هذه قائمة بالنماذج المجانية الشائعة
    return []string{
        "google/flan-t5-xl",              // نص عام
        "mistralai/Mistral-7B-Instruct-v0.2", // تعليمات
        "Qwen/Qwen2.5-7B-Instruct",       // دعم عربي
        "microsoft/phi-2",                // صغير وفعال
        "stabilityai/stable-diffusion-xl-base-1.0", // صور
        "openai/whisper-large-v3",        // صوت إلى نص
        "facebook/bart-large-cnn",        // تلخيص
        "Helsinki-NLP/opus-mt-ar-en",     // ترجمة عربي-إنجليزي
        "Helsinki-NLP/opus-mt-en-ar",     // ترجمة إنجليزي-عربي
    }, nil
}

// GetModelInfo الحصول على معلومات النموذج
func (p *HuggingFaceProvider) GetModelInfo(model string) (map[string]interface{}, error) {
    url := fmt.Sprintf("https://huggingface.co/api/models/%s", model)
    
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return nil, err
    }
    
    resp, err := p.client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    var info map[string]interface{}
    if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
        return nil, err
    }
    
    return info, nil
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

// IsAvailable التحقق من التوفر
func (p *HuggingFaceProvider) IsAvailable() bool {
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