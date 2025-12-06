package ai

import (
    "fmt"
    "log"
    "os"
    "sync"
)

// Client عميل AI متكامل
type Client struct {
    mu                sync.RWMutex
    providers         map[string]ProviderInterface
    multiProvider     *MultiProvider
    costManager       *CostManager
}

// NewClient إنشاء عميل AI جديد
func NewClient() (*Client, error) {
    c := &Client{
        providers: make(map[string]ProviderInterface),
    }
    
    // إنشاء مدير التكاليف
    costManager, err := NewCostManager()
    if err != nil {
        log.Printf("Warning: Failed to initialize cost manager: %v", err)
    }
    c.costManager = costManager
    
    // إنشاء مزود متعدد
    mp, err := NewMultiProvider()
    if err != nil {
        log.Printf("Warning: Failed to create multi-provider: %v", err)
        // استمرار بدون multi-provider
    } else {
        c.multiProvider = mp
    }
    
    // تهيئة مزود Ollama (دائمًا متاح محليًا)
    ollama := NewOllamaProvider()
    if ollama != nil {
        c.providers["ollama"] = ollama
        log.Println("✅ Ollama provider initialized")
    }
    
    // محاولة تهيئة مزود Hugging Face إذا كان هناك API key
    if token := os.Getenv("HUGGINGFACE_TOKEN"); token != "" {
        hf := NewHuggingFaceProvider()
        if hf != nil && hf.IsAvailable() {
            c.providers["huggingface"] = hf
            log.Println("✅ Hugging Face provider initialized")
        }
    }
    
    // محاولة تهيئة مزود Gemini إذا كان هناك API key
    if apiKey := os.Getenv("GEMINI_API_KEY"); apiKey != "" {
        gemini := NewGeminiProvider()
        if gemini != nil && gemini.IsAvailable() {
            c.providers["gemini"] = gemini
            log.Println("✅ Gemini provider initialized")
        }
    }
    
    if len(c.providers) == 0 {
        log.Println("⚠️ No AI providers available")
    } else {
        log.Printf("🤖 AI Client initialized with %d providers", len(c.providers))
    }
    
    return c, nil
}

// GenerateText توليد نص
func (c *Client) GenerateText(prompt, provider string) (string, error) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    
    var resp *TextResponse
    var err error
    
    if provider == "" || provider == "auto" {
        // استخدام MultiProvider للاختيار التلقائي إذا كان متاحاً
        if c.multiProvider != nil && c.multiProvider.IsAvailable() {
            req := TextRequest{
                Prompt: prompt,
                Model:  "llama3.2:3b",
            }
            
            resp, err = c.multiProvider.GenerateText(req)
        } else {
            // استخدام أول مزود متاح
            for _, p := range c.providers {
                if p.IsAvailable() {
                    req := TextRequest{
                        Prompt: prompt,
                    }
                    resp, err = p.GenerateText(req)
                    break
                }
            }
        }
    } else {
        // استخدام مزود محدد
        p, exists := c.providers[provider]
        if !exists {
            return "", fmt.Errorf("provider %s not found", provider)
        }
        
        req := TextRequest{
            Prompt: prompt,
        }
        resp, err = p.GenerateText(req)
    }
    
    if err != nil {
        return "", err
    }
    
    // تسجيل الاستخدام إذا كان هناك مدير تكاليف
    if c.costManager != nil && resp != nil {
        record := &UsageRecord{
            Provider:   provider,
            Type:       "text",
            Cost:       resp.Cost,
            Quantity:   int64(resp.Tokens),
            Success:    true,
            Timestamp:  resp.CreatedAt,
        }
        c.costManager.RecordUsage(record)
    }
    
    if resp == nil {
        return "", fmt.Errorf("no response generated")
    }
    
    return resp.Text, nil
}

// GenerateTextWithOptions توليد نص مع خيارات متقدمة
func (c *Client) GenerateTextWithOptions(req TextRequest) (*TextResponse, error) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    
    var provider ProviderInterface
    var err error
    
    if req.Model == "" || req.Model == "auto" {
        // استخدام MultiProvider للاختيار التلقائي إذا كان متاحاً
        if c.multiProvider != nil && c.multiProvider.IsAvailable() {
            return c.multiProvider.GenerateText(req)
        }
        
        // استخدام أول مزود نص متاح
        for _, p := range c.providers {
            if p.IsAvailable() && p.GetType() == "text" {
                provider = p
                break
            }
        }
    } else {
        // استخدام مزود محدد
        provider, err = c.getProviderByModel(req.Model)
        if err != nil {
            return nil, err
        }
    }
    
    if provider == nil {
        return nil, fmt.Errorf("no available text provider")
    }
    
    resp, err := provider.GenerateText(req)
    if err != nil {
        return nil, err
    }
    
    // تسجيل الاستخدام
    if c.costManager != nil {
        record := &UsageRecord{
            Provider:   provider.GetName(),
            Type:       "text",
            Cost:       resp.Cost,
            Quantity:   int64(resp.Tokens),
            Success:    true,
            Timestamp:  resp.CreatedAt,
        }
        c.costManager.RecordUsage(record)
    }
    
    return resp, nil
}

// GenerateImage توليد صورة
func (c *Client) GenerateImage(prompt, provider string) (string, error) {
    req := ImageRequest{
        Prompt: prompt,
    }
    
    resp, err := c.GenerateImageWithOptions(req)
    if err != nil {
        return "", err
    }
    
    return resp.URL, nil
}

// GenerateImageWithOptions توليد صورة مع خيارات متقدمة
func (c *Client) GenerateImageWithOptions(req ImageRequest) (*ImageResponse, error) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    
    // البحث عن مزود صور
    for _, p := range c.providers {
        if p.IsAvailable() {
            resp, err := p.GenerateImage(req)
            if err == nil {
                // تسجيل الاستخدام
                if c.costManager != nil {
                    record := &UsageRecord{
                        Provider:   p.GetName(),
                        Type:       "image",
                        Cost:       resp.Cost,
                        Quantity:   1,
                        Success:    true,
                        Timestamp:  resp.CreatedAt,
                    }
                    c.costManager.RecordUsage(record)
                }
                return resp, nil
            }
        }
    }
    
    return nil, fmt.Errorf("no available image provider")
}

// GenerateVideo توليد فيديو
func (c *Client) GenerateVideo(prompt, provider string) (string, error) {
    req := VideoRequest{
        Prompt:   prompt,
        Duration: 30, // 30 ثانية افتراضياً
    }
    
    resp, err := c.GenerateVideoWithOptions(req)
    if err != nil {
        return "", err
    }
    
    return resp.URL, nil
}

// GenerateVideoWithOptions توليد فيديو مع خيارات متقدمة
func (c *Client) GenerateVideoWithOptions(req VideoRequest) (*VideoResponse, error) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    
    // البحث عن مزود فيديو
    for _, p := range c.providers {
        if p.IsAvailable() {
            resp, err := p.GenerateVideo(req)
            if err == nil {
                // تسجيل الاستخدام
                if c.costManager != nil {
                    record := &UsageRecord{
                        Provider:   p.GetName(),
                        Type:       "video",
                        Cost:       resp.Cost,
                        Quantity:   1,
                        Success:    true,
                        Timestamp:  time.Now(),
                    }
                    c.costManager.RecordUsage(record)
                }
                return resp, nil
            }
        }
    }
    
    return nil, fmt.Errorf("no available video provider")
}

// AnalyzeText تحليل نص
func (c *Client) AnalyzeText(text, provider string) (*AnalysisResponse, error) {
    req := AnalysisRequest{
        Text: text,
    }
    
    return c.AnalyzeTextWithOptions(req)
}

// AnalyzeTextWithOptions تحليل نص مع خيارات متقدمة
func (c *Client) AnalyzeTextWithOptions(req AnalysisRequest) (*AnalysisResponse, error) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    
    // البحث عن مزود يدعم تحليل النصوص
    for _, p := range c.providers {
        if p.IsAvailable() {
            resp, err := p.AnalyzeText(req)
            if err == nil {
                // تسجيل الاستخدام
                if c.costManager != nil {
                    record := &UsageRecord{
                        Provider:   p.GetName(),
                        Type:       "analysis",
                        Cost:       resp.Cost,
                        Quantity:   1,
                        Success:    true,
                        Timestamp:  time.Now(),
                    }
                    c.costManager.RecordUsage(record)
                }
                return resp, nil
            }
        }
    }
    
    return nil, fmt.Errorf("no available text analysis provider")
}

// TranslateText ترجمة نص
func (c *Client) TranslateText(text, fromLang, toLang, provider string) (string, error) {
    req := TranslationRequest{
        Text:     text,
        FromLang: fromLang,
        ToLang:   toLang,
    }
    
    resp, err := c.TranslateTextWithOptions(req)
    if err != nil {
        return "", err
    }
    
    return resp.TranslatedText, nil
}

// TranslateTextWithOptions ترجمة نص مع خيارات متقدمة
func (c *Client) TranslateTextWithOptions(req TranslationRequest) (*TranslationResponse, error) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    
    // البحث عن مزود يدعم الترجمة
    for _, p := range c.providers {
        if p.IsAvailable() {
            resp, err := p.TranslateText(req)
            if err == nil {
                // تسجيل الاستخدام
                if c.costManager != nil {
                    record := &UsageRecord{
                        Provider:   p.GetName(),
                        Type:       "translation",
                        Cost:       resp.Cost,
                        Quantity:   1,
                        Success:    true,
                        Timestamp:  time.Now(),
                    }
                    c.costManager.RecordUsage(record)
                }
                return resp, nil
            }
        }
    }
    
    return nil, fmt.Errorf("no available translation provider")
}

// AnalyzeImage تحليل صورة
func (c *Client) AnalyzeImage(imageData []byte, prompt, provider string) (*AnalysisResponse, error) {
    req := AnalysisRequest{
        ImageData: imageData,
        Prompt:    prompt,
    }
    
    return c.AnalyzeImageWithOptions(req)
}

// AnalyzeImageWithOptions تحليل صورة مع خيارات متقدمة
func (c *Client) AnalyzeImageWithOptions(req AnalysisRequest) (*AnalysisResponse, error) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    
    // البحث عن مزود يدعم تحليل الصور
    for _, p := range c.providers {
        if p.IsAvailable() {
            resp, err := p.AnalyzeImage(req)
            if err == nil {
                // تسجيل الاستخدام
                if c.costManager != nil {
                    record := &UsageRecord{
                        Provider:   p.GetName(),
                        Type:       "image_analysis",
                        Cost:       resp.Cost,
                        Quantity:   1,
                        Success:    true,
                        Timestamp:  time.Now(),
                    }
                    c.costManager.RecordUsage(record)
                }
                return resp, nil
            }
        }
    }
    
    return nil, fmt.Errorf("no available image analysis provider")
}

// GetVideoStatus الحصول على حالة فيديو
func (c *Client) GetVideoStatus(operationID string) (*VideoResponse, error) {
    // هذه وظيفة تحتاج إلى VideoService
    // سنعود إليها لاحقاً
    return nil, fmt.Errorf("video service not available yet")
}

// GetAvailableProviders الحصول على المزودين المتاحين
func (c *Client) GetAvailableProviders() map[string][]string {
    c.mu.RLock()
    defer c.mu.RUnlock()
    
    providers := make(map[string][]string)
    
    // تصنيف المزودين حسب النوع
    textProviders := []string{}
    imageProviders := []string{}
    videoProviders := []string{}
    
    for name, provider := range c.providers {
        if provider.IsAvailable() {
            providerType := provider.GetType()
            switch providerType {
            case "text":
                textProviders = append(textProviders, name)
            case "image":
                imageProviders = append(imageProviders, name)
            case "video":
                videoProviders = append(videoProviders, name)
            default:
                textProviders = append(textProviders, name)
            }
        }
    }
    
    if len(textProviders) > 0 {
        providers["text"] = textProviders
    }
    if len(imageProviders) > 0 {
        providers["image"] = imageProviders
    }
    if len(videoProviders) > 0 {
        providers["video"] = videoProviders
    }
    
    // إضافة "auto" كخيار
    if len(textProviders) > 0 {
        providers["text"] = append(providers["text"], "auto")
    }
    
    return providers
}

// IsProviderAvailable التحقق من توفر مزود
func (c *Client) IsProviderAvailable(providerType, providerName string) bool {
    c.mu.RLock()
    defer c.mu.RUnlock()
    
    if providerName == "auto" {
        // التحقق إذا كان هناك أي مزود من النوع المطلوب متاح
        for _, p := range c.providers {
            if p.GetType() == providerType && p.IsAvailable() {
                return true
            }
        }
        return false
    }
    
    if p, exists := c.providers[providerName]; exists {
        return p.IsAvailable() && p.GetType() == providerType
    }
    
    return false
}

// GetProviderStats الحصول على إحصائيات مزود
func (c *Client) GetProviderStats(providerName string) (*ProviderStats, error) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    
    p, exists := c.providers[providerName]
    if !exists {
        return nil, fmt.Errorf("provider %s not found", providerName)
    }
    
    return p.GetStats(), nil
}

// GetUsageStatistics الحصول على إحصائيات الاستخدام
func (c *Client) GetUsageStatistics() map[string]interface{} {
    if c.costManager != nil {
        return c.costManager.GetUsageStatistics()
    }
    
    return map[string]interface{}{
        "total_cost": 0.0,
        "providers":  len(c.providers),
        "message":    "cost manager not available",
    }
}

// Close إغلاق العميل وتحرير الموارد
func (c *Client) Close() error {
    c.mu.Lock()
    defer c.mu.Unlock()
    
    log.Println("Closing AI client...")
    
    // إغلاق جميع المزودين
    for name, provider := range c.providers {
        if closer, ok := provider.(interface{ Close() error }); ok {
            if err := closer.Close(); err != nil {
                log.Printf("Error closing provider %s: %v", name, err)
            }
        }
    }
    
    return nil
}

// Helper functions

func (c *Client) getProviderByModel(model string) (ProviderInterface, error) {
    // بحث مبسط عن المزود المناسب للنموذج
    for _, provider := range c.providers {
        if provider.IsAvailable() {
            // يمكن إضافة منطق أكثر تعقيداً هنا
            return provider, nil
        }
    }
    
    return nil, fmt.Errorf("no provider available for model %s", model)
}

// RegisterProvider تسجيل مزود جديد
func (c *Client) RegisterProvider(name string, provider ProviderInterface) {
    c.mu.Lock()
    defer c.mu.Unlock()
    
    c.providers[name] = provider
    log.Printf("✅ Registered provider: %s", name)
}

// RemoveProvider إزالة مزود
func (c *Client) RemoveProvider(name string) {
    c.mu.Lock()
    defer c.mu.Unlock()
    
    if provider, exists := c.providers[name]; exists {
        if closer, ok := provider.(interface{ Close() error }); ok {
            if err := closer.Close(); err != nil {
                log.Printf("Error closing provider %s: %v", name, err)
            }
        }
        delete(c.providers, name)
        log.Printf("Removed provider: %s", name)
    }
}