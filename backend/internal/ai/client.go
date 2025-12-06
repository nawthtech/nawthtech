package ai

import (
    "fmt"
    "log"
    "os"
    "sync"
    
    "github.com/nawthtech/nawthtech/backend/internal/ai/providers"
    "github.com/nawthtech/nawthtech/backend/internal/ai/services"
)

// Client عميل AI متكامل يدعم النصوص والصور والفيديوهات
type Client struct {
    mu                sync.RWMutex
    textProviders     map[string]providers.TextProvider
    imageProviders    map[string]providers.ImageProvider
    videoProviders    map[string]providers.VideoProvider
    
    // Services
    ContentService    *services.ContentService
    AnalysisService   *services.AnalysisService
    StrategyService   *services.StrategyService
    MediaService      *services.MediaService
    TranslationService *services.TranslationService
    VideoService      *services.VideoService
}

// NewClient إنشاء عميل AI جديد
func NewClient() (*Client, error) {
    c := &Client{
        textProviders:  make(map[string]providers.TextProvider),
        imageProviders: make(map[string]providers.ImageProvider),
        videoProviders: make(map[string]providers.VideoProvider),
    }
    
    // تهيئة مزودي النصوص
    if err := c.initTextProviders(); err != nil {
        log.Printf("Warning: Text providers init failed: %v", err)
    }
    
    // تهيئة مزودي الصور
    if err := c.initImageProviders(); err != nil {
        log.Printf("Warning: Image providers init failed: %v", err)
    }
    
    // تهيئة مزودي الفيديو
    if err := c.initVideoProviders(); err != nil {
        log.Printf("Warning: Video providers init failed: %v", err)
    }
    
    // إنشاء الخدمات
    if err := c.initServices(); err != nil {
        log.Printf("Warning: Services init failed: %v", err)
    }
    
    log.Printf("🤖 AI Client initialized with %d text, %d image, %d video providers",
        len(c.textProviders), len(c.imageProviders), len(c.videoProviders))
    
    return c, nil
}

// initTextProviders تهيئة مزودي النصوص
func (c *Client) initTextProviders() error {
    // 1. Gemini (مجاني - 60 request/دقيقة)
    if apiKey := os.Getenv("GEMINI_API_KEY"); apiKey != "" {
        gemini, err := providers.NewGeminiProvider()
        if err == nil {
            c.textProviders["gemini"] = gemini
            log.Println("✅ Gemini text provider initialized")
        }
    }
    
    // 2. Ollama (محلي - مجاني بالكامل)
    ollama := providers.NewOllamaProvider()
    c.textProviders["ollama"] = ollama
    log.Println("✅ Ollama text provider initialized")
    
    // 3. Hugging Face (مجاني - 30k tokens/شهر)
    if token := os.Getenv("HUGGINGFACE_TOKEN"); token != "" {
        hf := providers.NewHuggingFaceProvider()
        c.textProviders["huggingface"] = hf
        log.Println("✅ Hugging Face text provider initialized")
    }
    
    if len(c.textProviders) == 0 {
        return fmt.Errorf("no text providers available")
    }
    
    return nil
}

// initImageProviders تهيئة مزودي الصور
func (c *Client) initImageProviders() error {
    // 1. Gemini Image Generation
    if apiKey := os.Getenv("GEMINI_API_KEY"); apiKey != "" {
        gemini, err := providers.NewGeminiProvider()
        if err == nil {
            c.imageProviders["gemini"] = gemini
            log.Println("✅ Gemini image provider initialized")
        }
    }
    
    // 2. Hugging Face (مجاني - 1000 صورة/شهر)
    if token := os.Getenv("HUGGINGFACE_TOKEN"); token != "" {
        hf := providers.NewHuggingFaceProvider()
        c.imageProviders["huggingface"] = hf
        log.Println("✅ Hugging Face image provider initialized")
    }
    
    return nil
}

// initVideoProviders تهيئة مزودي الفيديو
func (c *Client) initVideoProviders() error {
    // 1. Luma AI (مجاني - 30 فيديو/شهر)
    if apiKey := os.Getenv("LUMA_API_KEY"); apiKey != "" {
        luma, err := providers.NewLumaProvider()
        if err == nil {
            c.videoProviders["luma"] = luma
            log.Println("✅ Luma video provider initialized")
        }
    }
    
    // 2. Runway ML (مجاني - 125 ثانية/شهر)
    if apiKey := os.Getenv("RUNWAY_API_KEY"); apiKey != "" {
        runway, err := providers.NewRunwayProvider()
        if err == nil {
            c.videoProviders["runway"] = runway
            log.Println("✅ Runway video provider initialized")
        }
    }
    
    // 3. Pika Labs (مجاني - 100 فيديو/شهر)
    if apiKey := os.Getenv("PIKA_API_KEY"); apiKey != "" {
        pika, err := providers.NewPikaProvider()
        if err == nil {
            c.videoProviders["pika"] = pika
            log.Println("✅ Pika video provider initialized")
        }
    }
    
    return nil
}

// initServices تهيئة الخدمات
func (c *Client) initServices() error {
    // ContentService
    c.ContentService = services.NewContentService(c)
    
    // AnalysisService
    c.AnalysisService = services.NewAnalysisService(c)
    
    // StrategyService  
    c.StrategyService = services.NewStrategyService(c)
    
    // MediaService
    c.MediaService = services.NewMediaService(c)
    
    // TranslationService
    c.TranslationService = services.NewTranslationService(c)
    
    // VideoService
    c.VideoService = services.NewVideoService(c)
    
    return nil
}

// GetTextProvider الحصول على مزود النصوص
func (c *Client) GetTextProvider(name string) providers.TextProvider {
    c.mu.RLock()
    defer c.mu.RUnlock()
    
    if name == "" || name == "auto" {
        // اختيار تلقائي: Gemini أولاً، ثم Ollama
        if provider, ok := c.textProviders["gemini"]; ok {
            return provider
        }
        return c.textProviders["ollama"]
    }
    
    return c.textProviders[name]
}

// GetImageProvider الحصول على مزود الصور
func (c *Client) GetImageProvider(name string) providers.ImageProvider {
    c.mu.RLock()
    defer c.mu.RUnlock()
    
    if name == "" || name == "auto" {
        // اختيار تلقائي: Gemini أولاً، ثم Hugging Face
        if provider, ok := c.imageProviders["gemini"]; ok {
            return provider
        }
        if provider, ok := c.imageProviders["huggingface"]; ok {
            return provider
        }
    }
    
    return c.imageProviders[name]
}

// GetVideoProvider الحصول على مزود الفيديو
func (c *Client) GetVideoProvider(name string) providers.VideoProvider {
    c.mu.RLock()
    defer c.mu.RUnlock()
    
    if name == "" || name == "auto" {
        // اختيار تلقائي: Luma أولاً، ثم Runway، ثم Pika
        if provider, ok := c.videoProviders["luma"]; ok {
            return provider
        }
        if provider, ok := c.videoProviders["runway"]; ok {
            return provider
        }
        if provider, ok := c.videoProviders["pika"]; ok {
            return provider
        }
    }
    
    return c.videoProviders[name]
}

// GenerateText توليد نص
func (c *Client) GenerateText(prompt, provider string) (string, error) {
    p := c.GetTextProvider(provider)
    if p == nil {
        return "", fmt.Errorf("text provider %s not found", provider)
    }
    
    req := providers.TextRequest{
        Prompt: prompt,
    }
    
    resp, err := p.GenerateText(req)
    if err != nil {
        return "", err
    }
    
    return resp.Text, nil
}

// GenerateImage توليد صورة
func (c *Client) GenerateImage(prompt, provider string) (string, error) {
    p := c.GetImageProvider(provider)
    if p == nil {
        return "", fmt.Errorf("image provider %s not found", provider)
    }
    
    req := providers.ImageRequest{
        Prompt: prompt,
    }
    
    resp, err := p.GenerateImage(req)
    if err != nil {
        return "", err
    }
    
    return resp.ImageURL, nil
}

// GenerateVideo توليد فيديو
func (c *Client) GenerateVideo(prompt, provider string) (string, error) {
    p := c.GetVideoProvider(provider)
    if p == nil {
        return "", fmt.Errorf("video provider %s not found", provider)
    }
    
    req := providers.VideoRequest{
        Prompt: prompt,
    }
    
    resp, err := p.GenerateVideo(req)
    if err != nil {
        return "", err
    }
    
    return resp.VideoURL, nil
}

// AnalyzeText تحليل نص
func (c *Client) AnalyzeText(text, provider string) (*providers.AnalysisResponse, error) {
    p := c.GetTextProvider(provider)
    if p == nil {
        return nil, fmt.Errorf("text provider %s not found", provider)
    }
    
    req := providers.TextRequest{
        Prompt: text,
    }
    
    return p.AnalyzeText(req)
}

// TranslateText ترجمة نص
func (c *Client) TranslateText(text, fromLang, toLang, provider string) (string, error) {
    p := c.GetTextProvider(provider)
    if p == nil {
        return "", fmt.Errorf("text provider %s not found", provider)
    }
    
    req := providers.TranslationRequest{
        Text:     text,
        FromLang: fromLang,
        ToLang:   toLang,
    }
    
    resp, err := p.TranslateText(req)
    if err != nil {
        return "", err
    }
    
    return resp.TranslatedText, nil
}

// GetVideoStatus الحصول على حالة فيديو
func (c *Client) GetVideoStatus(operationID string) (*providers.VideoResponse, error) {
    if c.VideoService != nil {
        return c.VideoService.GetStatus(operationID)
    }
    return nil, fmt.Errorf("video service not available")
}

// GetAvailableProviders الحصول على المزودين المتاحين
func (c *Client) GetAvailableProviders() map[string][]string {
    c.mu.RLock()
    defer c.mu.RUnlock()
    
    providers := make(map[string][]string)
    
    // مزودي النصوص
    textProviders := make([]string, 0, len(c.textProviders))
    for name := range c.textProviders {
        textProviders = append(textProviders, name)
    }
    providers["text"] = textProviders
    
    // مزودي الصور
    imageProviders := make([]string, 0, len(c.imageProviders))
    for name := range c.imageProviders {
        imageProviders = append(imageProviders, name)
    }
    providers["image"] = imageProviders
    
    // مزودي الفيديو
    videoProviders := make([]string, 0, len(c.videoProviders))
    for name := range c.videoProviders {
        videoProviders = append(videoProviders, name)
    }
    providers["video"] = videoProviders
    
    return providers
}

// IsProviderAvailable التحقق من توفر مزود
func (c *Client) IsProviderAvailable(providerType, providerName string) bool {
    c.mu.RLock()
    defer c.mu.RUnlock()
    
    switch providerType {
    case "text":
        _, ok := c.textProviders[providerName]
        return ok
    case "image":
        _, ok := c.imageProviders[providerName]
        return ok
    case "video":
        _, ok := c.videoProviders[providerName]
        return ok
    default:
        return false
    }
}

// GetProviderStats الحصول على إحصائيات المزود
func (c *Client) GetProviderStats(providerType, providerName string) (*providers.ProviderStats, error) {
    var provider interface{}
    
    c.mu.RLock()
    switch providerType {
    case "text":
        if p, ok := c.textProviders[providerName]; ok {
            provider = p
        }
    case "image":
        if p, ok := c.imageProviders[providerName]; ok {
            provider = p
        }
    case "video":
        if p, ok := c.videoProviders[providerName]; ok {
            provider = p
        }
    }
    c.mu.RUnlock()
    
    if provider == nil {
        return nil, fmt.Errorf("provider %s/%s not found", providerType, providerName)
    }
    
    // محاولة الحصول على الإحصائيات إذا كانت الطريقة متوفرة
    switch p := provider.(type) {
    case interface{ GetStats() *providers.ProviderStats }:
        return p.GetStats(), nil
    default:
        return &providers.ProviderStats{
            Name:         providerName,
            Type:         providerType,
            IsAvailable:  true,
            Requests:     0,
            Errors:       0,
            LastUsed:     "",
        }, nil
    }
}

// Close إغلاق العميل وتحرير الموارد
func (c *Client) Close() error {
    c.mu.Lock()
    defer c.mu.Unlock()
    
    log.Println("Closing AI client...")
    
    // إغلاق جميع المزودين
    for name, provider := range c.textProviders {
        if closer, ok := provider.(interface{ Close() error }); ok {
            if err := closer.Close(); err != nil {
                log.Printf("Error closing text provider %s: %v", name, err)
            }
        }
    }
    
    for name, provider := range c.imageProviders {
        if closer, ok := provider.(interface{ Close() error }); ok {
            if err := closer.Close(); err != nil {
                log.Printf("Error closing image provider %s: %v", name, err)
            }
        }
    }
    
    for name, provider := range c.videoProviders {
        if closer, ok := provider.(interface{ Close() error }); ok {
            if err := closer.Close(); err != nil {
                log.Printf("Error closing video provider %s: %v", name, err)
            }
        }
    }
    
    return nil
}