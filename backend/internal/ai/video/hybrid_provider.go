package video

import (
    "fmt"
    "os"
    "sync"
    "time"
)

// HybridVideoProvider مزود هجين يدعم عدة مزودين للفيديو
type HybridVideoProvider struct {
    providers []VideoProvider
    mu        sync.RWMutex
    stats     map[string]*ProviderStats
}

// ProviderStats إحصائيات المزود
type ProviderStats struct {
    Name           string
    TotalRequests  int64
    Successful     int64
    Failed         int64
    TotalLatency   time.Duration
    LastUsed       time.Time
    LastError      string
}

// NewHybridVideoProvider إنشاء مزود فيديو هجين جديد
func NewHybridVideoProvider() *HybridVideoProvider {
    h := &HybridVideoProvider{
        stats: make(map[string]*ProviderStats),
    }
    
    fmt.Println("🔄 Initializing hybrid video provider...")
    
    // 1. أولاً: محاولة النماذج المحلية (مجانية)
    if os.Getenv("ENABLE_LOCAL_VIDEO") != "false" {
        svd := NewLocalSVDProvider()
        if svd != nil && svd.IsAvailable() {
            h.providers = append(h.providers, svd)
            h.stats[svd.Name()] = &ProviderStats{Name: svd.Name()}
            fmt.Printf("✅ Local SVD provider initialized: %s\n", svd.Name())
        } else {
            fmt.Println("⚠️  Local SVD provider not available")
        }
    }
    
    // 2. مزود وهمي للاختبار والتطوير
    if os.Getenv("ENABLE_DUMMY_VIDEO") == "true" || len(h.providers) == 0 {
        dummy := NewDummyVideoProvider()
        if dummy != nil {
            h.providers = append(h.providers, dummy)
            h.stats[dummy.Name()] = &ProviderStats{Name: dummy.Name()}
            fmt.Printf("✅ Dummy video provider initialized: %s\n", dummy.Name())
        }
    }
    
    // 3. خدمات مجانية محدودة (مثال: Stability AI)
    if apiKey := os.Getenv("STABILITY_API_KEY"); apiKey != "" {
        stability := NewStabilityVideoProvider(apiKey)
        if stability != nil && stability.IsAvailable() {
            h.providers = append(h.providers, stability)
            h.stats[stability.Name()] = &ProviderStats{Name: stability.Name()}
            fmt.Printf("✅ Stability AI provider initialized: %s\n", stability.Name())
        }
    }
    
    // 4. Google Veo (للجودة العالية)
    if os.Getenv("ENABLE_GOOGLE_VEO") == "true" {
        veo := NewGoogleVeoProvider()
        if veo != nil && veo.IsAvailable() {
            h.providers = append(h.providers, veo)
            h.stats[veo.Name()] = &ProviderStats{Name: veo.Name()}
            fmt.Printf("✅ Google Veo provider initialized: %s\n", veo.Name())
        }
    }
    
    // 5. Runway ML (خيار آخر)
    if os.Getenv("ENABLE_RUNWAY_ML") == "true" {
        runway := NewRunwayMLProvider()
        if runway != nil && runway.IsAvailable() {
            h.providers = append(h.providers, runway)
            h.stats[runway.Name()] = &ProviderStats{Name: runway.Name()}
            fmt.Printf("✅ Runway ML provider initialized: %s\n", runway.Name())
        }
    }
    
    if len(h.providers) == 0 {
        fmt.Println("⚠️  No video providers available, using dummy provider")
        dummy := NewDummyVideoProvider()
        h.providers = append(h.providers, dummy)
        h.stats[dummy.Name()] = &ProviderStats{Name: dummy.Name()}
    }
    
    fmt.Printf("✅ Hybrid provider ready with %d providers\n", len(h.providers))
    return h
}

// GenerateVideo توليد فيديو باستخدام أفضل مزود متاح
func (h *HybridVideoProvider) GenerateVideo(req VideoRequest) (*VideoResponse, error) {
    h.mu.RLock()
    defer h.mu.RUnlock()
    
    startTime := time.Now()
    
    // تسجيل محاولة الاستخدام
    fmt.Printf("🎬 Generating video with prompt: %.50s...\n", req.Prompt)
    
    // اختيار أفضل مزود حسب الأولوية
    provider := h.selectBestProvider(req)
    if provider == nil {
        return nil, fmt.Errorf("no suitable video provider available")
    }
    
    // تحديث الإحصائيات
    stats := h.stats[provider.Name()]
    if stats != nil {
        stats.TotalRequests++
        stats.LastUsed = time.Now()
    }
    
    // توليد الفيديو
    resp, err := provider.GenerateVideo(req)
    
    // تسجيل النتيجة
    if stats != nil {
        latency := time.Since(startTime)
        stats.TotalLatency += latency
        
        if err != nil {
            stats.Failed++
            stats.LastError = err.Error()
            fmt.Printf("❌ Provider %s failed: %v\n", provider.Name(), err)
        } else {
            stats.Successful++
            fmt.Printf("✅ Provider %s succeeded in %v\n", provider.Name(), latency)
        }
    }
    
    return resp, err
}

// selectBestProvider اختيار أفضل مزود حسب المعايير
func (h *HybridVideoProvider) selectBestProvider(req VideoRequest) VideoProvider {
    var bestProvider VideoProvider
    var bestScore float64
    
    for _, provider := range h.providers {
        if !provider.IsAvailable() {
            continue
        }
        
        score := h.calculateProviderScore(provider, req)
        
        if score > bestScore || bestProvider == nil {
            bestScore = score
            bestProvider = provider
        }
    }
    
    return bestProvider
}

// calculateProviderScore حساب درجة المزود
func (h *HybridVideoProvider) calculateProviderScore(provider VideoProvider, req VideoRequest) float64 {
    var score float64
    
    // 1. الأولوية للمحلي (مجاني وسريع)
    if provider.IsLocal() {
        score += 100
    }
    
    // 2. الأولوية للمجاني
    if provider.IsFree() {
        score += 50
    }
    
    // 3. دعم الدقة المطلوبة
    if provider.SupportsResolution(req.Resolution) {
        score += 30
    }
    
    // 4. حسب الإحصائيات (معدل النجاح)
    stats := h.stats[provider.Name()]
    if stats != nil && stats.TotalRequests > 0 {
        successRate := float64(stats.Successful) / float64(stats.TotalRequests)
        score += successRate * 40
        
        // مفضل المزودين الأقل استخداماً مؤخراً
        if !stats.LastUsed.IsZero() {
            hoursSinceLastUse := time.Since(stats.LastUsed).Hours()
            if hoursSinceLastUse > 1 {
                score += 10
            }
        }
        
        // مفضل المزودين الأسرع
        if stats.TotalRequests > 0 {
            avgLatency := stats.TotalLatency / time.Duration(stats.TotalRequests)
            if avgLatency < 30*time.Second {
                score += 20
            }
        }
    }
    
    // 5. حسب نوع الطلب
    if req.Duration <= 10 && provider.SupportsResolution("512x512") {
        score += 15 // جيد للفيديوهات القصيرة
    }
    
    return score
}

// Name اسم المزود
func (h *HybridVideoProvider) Name() string {
    return "hybrid_video_provider"
}

// IsAvailable التحقق من توفر أي مزود
func (h *HybridVideoProvider) IsAvailable() bool {
    h.mu.RLock()
    defer h.mu.RUnlock()
    
    for _, provider := range h.providers {
        if provider.IsAvailable() {
            return true
        }
    }
    return false
}

// IsLocal التحقق إذا كان المزود محلي
func (h *HybridVideoProvider) IsLocal() bool {
    return false // الهجين ليس محلياً بحد ذاته
}

// IsFree التحقق إذا كان المزود مجاني
func (h *HybridVideoProvider) IsFree() bool {
    // الهجين مجاني إذا كان فيه مزود مجاني متاح
    h.mu.RLock()
    defer h.mu.RUnlock()
    
    for _, provider := range h.providers {
        if provider.IsFree() && provider.IsAvailable() {
            return true
        }
    }
    return false
}

// SupportsResolution التحقق من دعم الدقة
func (h *HybridVideoProvider) SupportsResolution(resolution string) bool {
    h.mu.RLock()
    defer h.mu.RUnlock()
    
    for _, provider := range h.providers {
        if provider.SupportsResolution(resolution) && provider.IsAvailable() {
            return true
        }
    }
    return false
}

// GetAvailableProviders الحصول على قائمة المزودين المتاحين
func (h *HybridVideoProvider) GetAvailableProviders() []string {
    h.mu.RLock()
    defer h.mu.RUnlock()
    
    var available []string
    for _, provider := range h.providers {
        if provider.IsAvailable() {
            available = append(available, provider.Name())
        }
    }
    return available
}

// GetProviderStats الحصول على إحصائيات مزود محدد
func (h *HybridVideoProvider) GetProviderStats(name string) *ProviderStats {
    h.mu.RLock()
    defer h.mu.RUnlock()
    
    return h.stats[name]
}

// GetAllStats الحصول على جميع الإحصائيات
func (h *HybridVideoProvider) GetAllStats() map[string]ProviderStats {
    h.mu.RLock()
    defer h.mu.RUnlock()
    
    stats := make(map[string]ProviderStats)
    for name, stat := range h.stats {
        stats[name] = *stat
    }
    return stats
}

// GetCapabilities الحصول على قدرات المزود الهجين
func (h *HybridVideoProvider) GetCapabilities() map[string]interface{} {
    h.mu.RLock()
    defer h.mu.RUnlock()
    
    capabilities := make(map[string]interface{})
    var resolutions []string
    var providers []string
    
    for _, provider := range h.providers {
        if provider.IsAvailable() {
            providers = append(providers, provider.Name())
            
            // جمع جميع الدقات المدعومة
            // هذه قائمة افتراضية، يمكن تحسينها
            resolutions = append(resolutions, 
                "512x512", "576x1024", "1024x576",
                "768x768", "1024x1024",
            )
        }
    }
    
    capabilities["providers"] = providers
    capabilities["resolutions"] = removeDuplicates(resolutions)
    capabilities["hybrid"] = true
    capabilities["smart_selection"] = true
    capabilities["fallback_enabled"] = true
    
    return capabilities
}

// TestAllProviders اختبار جميع المزودين
func (h *HybridVideoProvider) TestAllProviders() map[string]bool {
    h.mu.RLock()
    defer h.mu.RUnlock()
    
    results := make(map[string]bool)
    for _, provider := range h.providers {
        results[provider.Name()] = provider.IsAvailable()
    }
    return results
}

// Helper function to remove duplicates from string slice
func removeDuplicates(slice []string) []string {
    keys := make(map[string]bool)
    list := []string{}
    for _, entry := range slice {
        if _, value := keys[entry]; !value {
            keys[entry] = true
            list = append(list, entry)
        }
    }
    return list
}

// DummyVideoProvider مزود وهمي للاختبار (لإكمال الكود)
type DummyVideoProvider struct{}

func NewDummyVideoProvider() *DummyVideoProvider {
    return &DummyVideoProvider{}
}

func (p *DummyVideoProvider) GenerateVideo(req VideoRequest) (*VideoResponse, error) {
    // محاكاة وقت التوليد
    time.Sleep(2 * time.Second)
    
    return &VideoResponse{
        Success:    true,
        VideoURL:   "https://example.com/dummy-video.mp4",
        Duration:   req.Duration,
        Width:      512,
        Height:     512,
        Resolution: "512x512",
        Format:     "mp4",
        Provider:   "dummy",
        Cost:       0.0,
        Status:     "completed",
        CreatedAt:  time.Now(),
        Timestamp:  time.Now().Unix(),
    }, nil
}

func (p *DummyVideoProvider) Name() string {
    return "dummy_video"
}

func (p *DummyVideoProvider) IsAvailable() bool {
    return true
}

func (p *DummyVideoProvider) IsLocal() bool {
    return true
}

func (p *DummyVideoProvider) IsFree() bool {
    return true
}

func (p *DummyVideoProvider) SupportsResolution(resolution string) bool {
    return resolution == "512x512" || resolution == "256x256"
}

// StabilityVideoProvider مزود Stability AI (مثال)
type StabilityVideoProvider struct {
    apiKey string
}

func NewStabilityVideoProvider(apiKey string) *StabilityVideoProvider {
    return &StabilityVideoProvider{apiKey: apiKey}
}

func (p *StabilityVideoProvider) GenerateVideo(req VideoRequest) (*VideoResponse, error) {
    // تنفيذ Stability AI هنا
    return nil, fmt.Errorf("Stability AI provider not implemented yet")
}

func (p *StabilityVideoProvider) Name() string {
    return "stability_ai"
}

func (p *StabilityVideoProvider) IsAvailable() bool {
    return p.apiKey != ""
}

func (p *StabilityVideoProvider) IsLocal() bool {
    return false
}

func (p *StabilityVideoProvider) IsFree() bool {
    return false // Stability AI له حدود مجانية ثم مدفوعة
}

func (p *StabilityVideoProvider) SupportsResolution(resolution string) bool {
    supported := []string{"512x512", "576x1024", "1024x576"}
    for _, res := range supported {
        if res == resolution {
            return true
        }
    }
    return false
}

// GoogleVeoProvider مزود Google Veo (مثال)
type GoogleVeoProvider struct{}

func NewGoogleVeoProvider() *GoogleVeoProvider {
    return &GoogleVeoProvider{}
}

func (p *GoogleVeoProvider) GenerateVideo(req VideoRequest) (*VideoResponse, error) {
    // تنفيذ Google Veo هنا
    return nil, fmt.Errorf("Google Veo provider not implemented yet")
}

func (p *GoogleVeoProvider) Name() string {
    return "google_veo"
}

func (p *GoogleVeoProvider) IsAvailable() bool {
    return os.Getenv("GOOGLE_API_KEY") != ""
}

func (p *GoogleVeoProvider) IsLocal() bool {
    return false
}

func (p *GoogleVeoProvider) IsFree() bool {
    return false // Google Veo مدفوع
}

func (p *GoogleVeoProvider) SupportsResolution(resolution string) bool {
    return resolution == "1920x1080" || resolution == "1080x1920"
}

// RunwayMLProvider مزود Runway ML (مثال)
type RunwayMLProvider struct{}

func NewRunwayMLProvider() *RunwayMLProvider {
    return &RunwayMLProvider{}
}

func (p *RunwayMLProvider) GenerateVideo(req VideoRequest) (*VideoResponse, error) {
    // تنفيذ Runway ML هنا
    return nil, fmt.Errorf("Runway ML provider not implemented yet")
}

func (p *RunwayMLProvider) Name() string {
    return "runway_ml"
}

func (p *RunwayMLProvider) IsAvailable() bool {
    return os.Getenv("RUNWAYML_API_KEY") != ""
}

func (p *RunwayMLProvider) IsLocal() bool {
    return false
}

func (p *RunwayMLProvider) IsFree() bool {
    return false // Runway ML له حدود مجانية
}

func (p *RunwayMLProvider) SupportsResolution(resolution string) bool {
    return resolution == "512x512" || resolution == "768x768"
}