package ai

import (
    "fmt"
    "log"
    "os"
    "strings"
    "sync"
    "time"
)

// ProviderType نوع المزود
type ProviderType string

// ثوابت أنواع المزودين
const (
    ProviderGemini      ProviderType = "gemini"
    ProviderOllama      ProviderType = "ollama"
    ProviderHuggingFace ProviderType = "huggingface"
)

// MultiProviderStats إحصائيات المزود المتعدد
type MultiProviderStats struct {
    TotalRequests     int64
    Successful        int64
    Failed            int64
    TotalCost         float64
    ProviderStats     map[ProviderType]*ProviderStats
    LastRotation      map[string]time.Time
    FallbackCount     map[ProviderType]int64
}

// RoutingStrategy واجهة إستراتيجية التوجيه
type RoutingStrategy interface {
    SelectProvider(userTier, promptType, providerType string) ProviderType
    GetFallbackChain(primary ProviderType, providerType string) []ProviderType
}

// MultiProvider مزود متعدد يدعم عدة مزودين AI
type MultiProvider struct {
    mu              sync.RWMutex
    providers       map[ProviderType]ProviderInterface
    textProviders   map[string]ProviderInterface
    imageProviders  map[string]ProviderInterface
    videoProviders  map[string]ProviderInterface
    strategy        RoutingStrategy
    costManager     *CostManager
    stats           *MultiProviderStats
}

// NewMultiProvider إنشاء مزود متعدد جديد
func NewMultiProvider() (*MultiProvider, error) {
    mp := &MultiProvider{
        providers:      make(map[ProviderType]ProviderInterface),
        textProviders:  make(map[string]ProviderInterface),
        imageProviders: make(map[string]ProviderInterface),
        videoProviders: make(map[string]ProviderInterface),
        strategy:       &DefaultStrategy{},
        stats: &MultiProviderStats{
            ProviderStats: make(map[ProviderType]*ProviderStats),
            LastRotation:  make(map[string]time.Time),
            FallbackCount: make(map[ProviderType]int64),
        },
    }
    
    // تهيئة مدير التكاليف
    cm, err := NewCostManager()
    if err != nil {
        log.Printf("Warning: Failed to initialize cost manager: %v", err)
    }
    mp.costManager = cm
    
    // تهيئة المزودين
    if err := mp.initProviders(); err != nil {
        return nil, fmt.Errorf("failed to initialize providers: %w", err)
    }
    
    // تهيئة الإحصائيات
    mp.updateProviderStats()
    
    log.Printf("🤖 MultiProvider initialized with %d total providers", len(mp.providers))
    
    return mp, nil
}

// initProviders تهيئة جميع المزودين
func (mp *MultiProvider) initProviders() error {
    mp.mu.Lock()
    defer mp.mu.Unlock()
    
    // 1. Ollama Provider (دائمًا متاح محليًا)
    ollama := NewOllamaProvider()
    mp.providers[ProviderOllama] = ollama
    mp.textProviders["ollama"] = ollama
    log.Println("✅ Ollama provider initialized")
    
    // 2. Hugging Face Provider
    if token := getEnvWithFallback("HUGGINGFACE_TOKEN", ""); token != "" {
        hf := NewHuggingFaceProvider()
        if hf.IsAvailable() {
            mp.providers[ProviderHuggingFace] = hf
            mp.textProviders["huggingface"] = hf
            mp.imageProviders["huggingface"] = hf
            log.Println("✅ Hugging Face provider initialized")
        }
    }
    
    // 3. Gemini Provider
    if apiKey := getEnvWithFallback("GEMINI_API_KEY", ""); apiKey != "" {
        gemini := NewGeminiProvider()
        if gemini.IsAvailable() {
            mp.providers[ProviderGemini] = gemini
            mp.textProviders["gemini"] = gemini
            mp.imageProviders["gemini"] = gemini
            log.Println("✅ Gemini provider initialized")
        }
    }
    
    if len(mp.providers) == 0 {
        return fmt.Errorf("no AI providers available")
    }
    
    return nil
}

// updateProviderStats تحديث إحصائيات جميع المزودين
func (mp *MultiProvider) updateProviderStats() {
    mp.mu.Lock()
    defer mp.mu.Unlock()
    
    for providerType, provider := range mp.providers {
        if _, exists := mp.stats.ProviderStats[providerType]; !exists {
            mp.stats.ProviderStats[providerType] = &ProviderStats{
                Name: provider.GetName(),
                Type: provider.GetType(),
            }
        }
        
        stats := mp.stats.ProviderStats[providerType]
        stats.IsAvailable = provider.IsAvailable()
        
        // الحصول على الإحصائيات من المزود نفسه
        if providerStats := provider.GetStats(); providerStats != nil {
            stats.Requests = providerStats.Requests
            stats.Successful = providerStats.Successful
            stats.Failed = providerStats.Failed
            stats.TotalCost = providerStats.TotalCost
            stats.AvgLatency = providerStats.AvgLatency
            stats.SuccessRate = providerStats.SuccessRate
        }
    }
}

// updateRequestStats تحديث إحصائيات الطلب
func (mp *MultiProvider) updateRequestStats(providerType ProviderType, success bool, cost float64) {
    mp.mu.Lock()
    defer mp.mu.Unlock()
    
    mp.stats.TotalRequests++
    if success {
        mp.stats.Successful++
    } else {
        mp.stats.Failed++
    }
    mp.stats.TotalCost += cost
    
    // تحديث إحصائيات المزود المحدد
    if _, exists := mp.stats.ProviderStats[providerType]; !exists {
        mp.stats.ProviderStats[providerType] = &ProviderStats{
            Name: string(providerType),
            Type: getProviderType(providerType),
        }
    }
    
    stats := mp.stats.ProviderStats[providerType]
    stats.Requests++
    if success {
        stats.Successful++
    } else {
        stats.Failed++
    }
    stats.TotalCost += cost
    stats.LastUsed = time.Now()
    
    if stats.Requests > 0 {
        stats.SuccessRate = float64(stats.Successful) / float64(stats.Requests) * 100
    }
}

// GenerateText توليد نص
func (mp *MultiProvider) GenerateText(req TextRequest) (*TextResponse, error) {
    startTime := time.Now()
    
    // تحديد المزود المناسب
    providerType := mp.strategy.SelectProvider(req.UserTier, "text", "text")
    
    // البحث عن المزود
    provider, err := mp.getProvider(providerType, "text")
    if err != nil {
        return nil, err
    }
    
    // توليد النص
    resp, err := provider.GenerateText(req)
    
    // تحديث الإحصائيات
    latency := float64(time.Since(startTime).Milliseconds())
    mp.updateRequestStats(providerType, err == nil, provider.GetCost())
    
    // تسجيل الاستخدام
    if mp.costManager != nil {
        record := &UsageRecord{
            UserID:     req.UserID,
            UserTier:   req.UserTier,
            Provider:   provider.GetName(),
            Type:       "text",
            Cost:       provider.GetCost(),
            Quantity:   int64(len(req.Prompt) / 4), // تقدير تقريبي
            Latency:    latency,
            Success:    err == nil,
            Timestamp:  time.Now(),
            Metadata: map[string]interface{}{
                "model": req.Model,
            },
        }
        mp.costManager.RecordUsage(record)
    }
    
    return resp, err
}

// GenerateImage توليد صورة
func (mp *MultiProvider) GenerateImage(req ImageRequest) (*ImageResponse, error) {
    startTime := time.Now()
    
    // تحديد المزود المناسب
    providerType := mp.strategy.SelectProvider(req.UserTier, "image", "image")
    
    // البحث عن المزود
    provider, err := mp.getProvider(providerType, "image")
    if err != nil {
        return nil, err
    }
    
    // توليد الصورة
    resp, err := provider.GenerateImage(req)
    
    // تحديث الإحصائيات
    latency := float64(time.Since(startTime).Milliseconds())
    mp.updateRequestStats(providerType, err == nil, provider.GetCost())
    
    // تسجيل الاستخدام
    if mp.costManager != nil {
        record := &UsageRecord{
            UserID:     req.UserID,
            UserTier:   req.UserTier,
            Provider:   provider.GetName(),
            Type:       "image",
            Cost:       provider.GetCost(),
            Quantity:   1,
            Latency:    latency,
            Success:    err == nil,
            Timestamp:  time.Now(),
        }
        mp.costManager.RecordUsage(record)
    }
    
    return resp, err
}

// GenerateVideo توليد فيديو
func (mp *MultiProvider) GenerateVideo(req VideoRequest) (*VideoResponse, error) {
    startTime := time.Now()
    
    // تحديد المزود المناسب
    providerType := mp.strategy.SelectProvider(req.UserTier, "video", "video")
    
    // البحث عن المزود
    provider, err := mp.getProvider(providerType, "video")
    if err != nil {
        return nil, err
    }
    
    // توليد الفيديو
    resp, err := provider.GenerateVideo(req)
    
    // تحديث الإحصائيات
    latency := float64(time.Since(startTime).Milliseconds())
    mp.updateRequestStats(providerType, err == nil, provider.GetCost())
    
    // تسجيل الاستخدام
    if mp.costManager != nil {
        record := &UsageRecord{
            UserID:     req.UserID,
            UserTier:   req.UserTier,
            Provider:   provider.GetName(),
            Type:       "video",
            Cost:       provider.GetCost(),
            Quantity:   1,
            Latency:    latency,
            Success:    err == nil,
            Timestamp:  time.Now(),
            Metadata: map[string]interface{}{
                "duration": req.Duration,
            },
        }
        mp.costManager.RecordUsage(record)
    }
    
    return resp, err
}

// AnalyzeText تحليل نص
func (mp *MultiProvider) AnalyzeText(req AnalysisRequest) (*AnalysisResponse, error) {
    startTime := time.Now()
    
    // تحديد المزود المناسب
    providerType := mp.strategy.SelectProvider(req.UserTier, "analysis", "text")
    
    // البحث عن المزود
    provider, err := mp.getProvider(providerType, "text")
    if err != nil {
        return nil, err
    }
    
    // تحليل النص
    resp, err := provider.AnalyzeText(req)
    
    // تحديث الإحصائيات
    latency := float64(time.Since(startTime).Milliseconds())
    mp.updateRequestStats(providerType, err == nil, provider.GetCost())
    
    // تسجيل الاستخدام
    if mp.costManager != nil {
        record := &UsageRecord{
            UserID:     req.UserID,
            UserTier:   req.UserTier,
            Provider:   provider.GetName(),
            Type:       "analysis",
            Cost:       provider.GetCost(),
            Quantity:   1,
            Latency:    latency,
            Success:    err == nil,
            Timestamp:  time.Now(),
        }
        mp.costManager.RecordUsage(record)
    }
    
    return resp, err
}

// AnalyzeImage تحليل صورة
func (mp *MultiProvider) AnalyzeImage(req AnalysisRequest) (*AnalysisResponse, error) {
    startTime := time.Now()
    
    // تحديد المزود المناسب
    providerType := mp.strategy.SelectProvider(req.UserTier, "analysis", "image")
    
    // البحث عن المزود
    provider, err := mp.getProvider(providerType, "image")
    if err != nil {
        return nil, err
    }
    
    // تحليل الصورة
    resp, err := provider.AnalyzeImage(req)
    
    // تحديث الإحصائيات
    latency := float64(time.Since(startTime).Milliseconds())
    mp.updateRequestStats(providerType, err == nil, provider.GetCost())
    
    // تسجيل الاستخدام
    if mp.costManager != nil {
        record := &UsageRecord{
            UserID:     req.UserID,
            UserTier:   req.UserTier,
            Provider:   provider.GetName(),
            Type:       "image_analysis",
            Cost:       provider.GetCost(),
            Quantity:   1,
            Latency:    latency,
            Success:    err == nil,
            Timestamp:  time.Now(),
        }
        mp.costManager.RecordUsage(record)
    }
    
    return resp, err
}

// TranslateText ترجمة نص
func (mp *MultiProvider) TranslateText(req TranslationRequest) (*TranslationResponse, error) {
    startTime := time.Now()
    
    // تحديد المزود المناسب
    providerType := mp.strategy.SelectProvider(req.UserTier, "translation", "text")
    
    // البحث عن المزود
    provider, err := mp.getProvider(providerType, "text")
    if err != nil {
        return nil, err
    }
    
    // ترجمة النص
    resp, err := provider.TranslateText(req)
    
    // تحديث الإحصائيات
    latency := float64(time.Since(startTime).Milliseconds())
    mp.updateRequestStats(providerType, err == nil, provider.GetCost())
    
    // تسجيل الاستخدام
    if mp.costManager != nil {
        record := &UsageRecord{
            UserID:     req.UserID,
            UserTier:   req.UserTier,
            Provider:   provider.GetName(),
            Type:       "translation",
            Cost:       provider.GetCost(),
            Quantity:   1,
            Latency:    latency,
            Success:    err == nil,
            Timestamp:  time.Now(),
        }
        mp.costManager.RecordUsage(record)
    }
    
    return resp, err
}

// getProvider الحصول على مزود من النوع المحدد
func (mp *MultiProvider) getProvider(providerType ProviderType, requestedType string) (ProviderInterface, error) {
    mp.mu.RLock()
    defer mp.mu.RUnlock()
    
    // محاولة الحصول على المزود المحدد
    provider, exists := mp.providers[providerType]
    if !exists || !provider.IsAvailable() {
        // استخدام التسلسل الاحتياطي
        fallbackChain := mp.strategy.GetFallbackChain(providerType, requestedType)
        for _, fbType := range fallbackChain {
            if fbProvider, fbExists := mp.providers[fbType]; fbExists && fbProvider.IsAvailable() {
                mp.stats.FallbackCount[fbType]++
                log.Printf("🔄 Fallback from %s to %s", providerType, fbType)
                return fbProvider, nil
            }
        }
        return nil, fmt.Errorf("no available %s provider", requestedType)
    }
    
    return provider, nil
}

// GetTextProvider الحصول على مزود نصوص محدد
func (mp *MultiProvider) GetTextProvider(name string) ProviderInterface {
    mp.mu.RLock()
    defer mp.mu.RUnlock()
    
    return mp.textProviders[name]
}

// GetImageProvider الحصول على مزود صور محدد
func (mp *MultiProvider) GetImageProvider(name string) ProviderInterface {
    mp.mu.RLock()
    defer mp.mu.RUnlock()
    
    return mp.imageProviders[name]
}

// GetVideoProvider الحصول على مزود فيديوهات محدد
func (mp *MultiProvider) GetVideoProvider(name string) ProviderInterface {
    mp.mu.RLock()
    defer mp.mu.RUnlock()
    
    return mp.videoProviders[name]
}

// GetAvailableProviders الحصول على المزودين المتاحين
func (mp *MultiProvider) GetAvailableProviders() map[string][]string {
    mp.mu.RLock()
    defer mp.mu.RUnlock()
    
    result := make(map[string][]string)
    
    // مزودي النصوص
    textProviders := make([]string, 0, len(mp.textProviders))
    for name, provider := range mp.textProviders {
        if provider.IsAvailable() {
            textProviders = append(textProviders, name)
        }
    }
    if len(textProviders) > 0 {
        result["text"] = textProviders
    }
    
    // مزودي الصور
    imageProviders := make([]string, 0, len(mp.imageProviders))
    for name, provider := range mp.imageProviders {
        if provider.IsAvailable() {
            imageProviders = append(imageProviders, name)
        }
    }
    if len(imageProviders) > 0 {
        result["image"] = imageProviders
    }
    
    // مزودي الفيديو
    videoProviders := make([]string, 0, len(mp.videoProviders))
    for name, provider := range mp.videoProviders {
        if provider.IsAvailable() {
            videoProviders = append(videoProviders, name)
        }
    }
    if len(videoProviders) > 0 {
        result["video"] = videoProviders
    }
    
    return result
}

// SetRoutingStrategy تعيين إستراتيجية التوجيه
func (mp *MultiProvider) SetRoutingStrategy(strategy RoutingStrategy) {
    mp.mu.Lock()
    defer mp.mu.Unlock()
    
    mp.strategy = strategy
}

// GetStats الحصول على إحصائيات المزود المتعدد
func (mp *MultiProvider) GetStats() *ProviderStats {
    mp.mu.RLock()
    defer mp.mu.RUnlock()
    
    stats := &ProviderStats{
        Name:        "MultiProvider",
        Type:        "multi",
        IsAvailable: len(mp.providers) > 0,
        Requests:    mp.stats.TotalRequests,
        Successful:  mp.stats.Successful,
        Failed:      mp.stats.Failed,
        TotalCost:   mp.stats.TotalCost,
        AvgLatency:  0.0,
        LastUsed:    time.Time{},
    }
    
    if mp.stats.TotalRequests > 0 {
        stats.SuccessRate = float64(mp.stats.Successful) / float64(mp.stats.TotalRequests) * 100
    }
    
    return stats
}

// GetProviderStats الحصول على إحصائيات مزود محدد
func (mp *MultiProvider) GetProviderStats(providerType ProviderType) (*ProviderStats, error) {
    mp.mu.RLock()
    defer mp.mu.RUnlock()
    
    if stats, exists := mp.stats.ProviderStats[providerType]; exists {
        return stats, nil
    }
    
    return nil, fmt.Errorf("provider stats not found: %s", providerType)
}

// GetName اسم المزود
func (mp *MultiProvider) GetName() string {
    return "MultiProvider"
}

// GetType نوع المزود
func (mp *MultiProvider) GetType() string {
    return "multi"
}

// IsAvailable التحقق من التوفر
func (mp *MultiProvider) IsAvailable() bool {
    mp.mu.RLock()
    defer mp.mu.RUnlock()
    return len(mp.providers) > 0
}

// GetCost التكلفة
func (mp *MultiProvider) GetCost() float64 {
    return 0.0 // سيتم حسابها بناءً على الاستخدام الفعلي
}

// SupportsStreaming يدعم التدفق
func (mp *MultiProvider) SupportsStreaming() bool {
    mp.mu.RLock()
    defer mp.mu.RUnlock()
    
    // التحقق إذا كان أي مزود يدعم التدفق
    for _, provider := range mp.providers {
        if provider.SupportsStreaming() {
            return true
        }
    }
    return false
}

// SupportsEmbedding يدعم التضمين
func (mp *MultiProvider) SupportsEmbedding() bool {
    mp.mu.RLock()
    defer mp.mu.RUnlock()
    
    // التحقق إذا كان أي مزود يدعم التضمين
    for _, provider := range mp.providers {
        if provider.SupportsEmbedding() {
            return true
        }
    }
    return false
}

// GetMaxTokens الحد الأقصى للرموز
func (mp *MultiProvider) GetMaxTokens() int {
    mp.mu.RLock()
    defer mp.mu.RUnlock()
    
    maxTokens := 0
    for _, provider := range mp.providers {
        if tokens := provider.GetMaxTokens(); tokens > maxTokens {
            maxTokens = tokens
        }
    }
    
    if maxTokens == 0 {
        return 2048 // القيمة الافتراضية
    }
    return maxTokens
}

// GetSupportedLanguages اللغات المدعومة
func (mp *MultiProvider) GetSupportedLanguages() []string {
    mp.mu.RLock()
    defer mp.mu.RUnlock()
    
    languages := make(map[string]bool)
    for _, provider := range mp.providers {
        for _, lang := range provider.GetSupportedLanguages() {
            languages[lang] = true
        }
    }
    
    result := make([]string, 0, len(languages))
    for lang := range languages {
        result = append(result, lang)
    }
    
    if len(result) == 0 {
        return []string{"ar", "en", "es", "fr", "de"}
    }
    return result
}

// DefaultStrategy إستراتيجية افتراضية
type DefaultStrategy struct{}

func (s *DefaultStrategy) SelectProvider(userTier, promptType, providerType string) ProviderType {
    // إستراتيجية بسيطة حسب طبقة المستخدم
    switch userTier {
    case "premium", "enterprise":
        // المستخدمين المميزين يحصلون على Gemini
        if providerType == "text" || providerType == "" {
            return ProviderGemini
        }
    case "free":
        fallback:
        // المستخدمين المجانيين يحصلون على Ollama أو HuggingFace
        if providerType == "text" || providerType == "" {
            return ProviderOllama
        }
    default:
        goto fallback
    }
    
    // للأنواع الأخرى
    switch providerType {
    case "image":
        return ProviderHuggingFace
    case "video":
        return ProviderOllama // Ollama لا يدعم الفيديو، لكن نستخدمه كاحتياطي
    default:
        return ProviderOllama
    }
}

func (s *DefaultStrategy) GetFallbackChain(primary ProviderType, providerType string) []ProviderType {
    chains := map[ProviderType][]ProviderType{
        ProviderGemini:      {ProviderHuggingFace, ProviderOllama},
        ProviderHuggingFace: {ProviderOllama, ProviderGemini},
        ProviderOllama:      {ProviderHuggingFace, ProviderGemini},
    }
    
    if chain, exists := chains[primary]; exists {
        return chain
    }
    
    // سلسلة احتياطية افتراضية
    return []ProviderType{ProviderOllama, ProviderHuggingFace, ProviderGemini}
}

// Helper functions

func getEnvWithFallback(key, fallback string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return fallback
}

func getProviderType(providerType ProviderType) string {
    switch providerType {
    case ProviderGemini, ProviderOllama, ProviderHuggingFace:
        return "text"
    default:
        return "mixed"
    }
}