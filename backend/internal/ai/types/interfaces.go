package types

import "time"

// MultiProviderInterface واجهة للمزود المتعدد
type MultiProviderInterface interface {
    ProviderInterface
    GetActiveProvider(providerType string) string
    SetActiveProvider(providerType, providerName string) error
    GetAvailableProviders(providerType string) []string
    GetProviderStats(providerType, providerName string) (*ProviderStats, error)
    RotateProvider(providerType string) error
    GetFallbackChain(providerType string) []string
}

// CostManagerInterface واجهة لإدارة التكاليف
type CostManagerInterface interface {
    RecordUsage(record *UsageRecord) error
    CanUseAI(userID, requestType string) (bool, string)
    GetUsageStatistics() map[string]interface{}
    GetUserQuotas(userID string) (map[string]*Quota, error)
    ResetUserQuotas(userID string) error
    SetLimits(monthly, daily float64)
    GetProviderStats(providerName string) (*ProviderStats, error)
}

// CacheManagerInterface واجهة لإدارة الذاكرة المؤقتة
type CacheManagerInterface interface {
    Get(key string) (interface{}, bool)
    Set(key string, value interface{}, ttl time.Duration)
    Delete(key string)
    Clear()
    Size() int
    GetStats() CacheStats
}

// AIClientInterface واجهة عميل AI
type AIClientInterface interface {
    // العمليات الأساسية
    GenerateText(prompt, provider string) (string, error)
    GenerateImage(prompt, provider string) (string, error)
    GenerateVideo(prompt, provider string) (string, error)
    
    // العمليات المتقدمة
    GenerateTextWithOptions(req TextRequest) (*TextResponse, error)
    GenerateImageWithOptions(req ImageRequest) (*ImageResponse, error)
    GenerateVideoWithOptions(req VideoRequest) (*VideoResponse, error)
    AnalyzeText(text, provider string) (*AnalysisResponse, error)
    AnalyzeTextWithOptions(req AnalysisRequest) (*AnalysisResponse, error)
    TranslateText(text, fromLang, toLang, provider string) (string, error)
    TranslateTextWithOptions(req TranslationRequest) (*TranslationResponse, error)
    AnalyzeImage(imageData []byte, prompt, provider string) (*AnalysisResponse, error)
    AnalyzeImageWithOptions(req AnalysisRequest) (*AnalysisResponse, error)
    
    // معلومات النظام
    GetVideoStatus(operationID string) (*VideoResponse, error)
    GetAvailableProviders() map[string][]string
    IsProviderAvailable(providerType, providerName string) bool
    GetProviderStats(providerName string) (*ProviderStats, error)
    GetUsageStatistics() map[string]interface{}
    
    // إدارة المزودين
    RegisterProvider(name string, provider ProviderInterface)
    RemoveProvider(name string)
    
    // إدارة الاتصال
    Close() error
}

// --- الأنواع من interfaces.go القديم ---

// CacheStats إحصائيات الذاكرة المؤقتة
type CacheStats struct {
    Hits        int64                  `json:"hits"`
    Misses      int64                  `json:"misses"`
    HitRate     float64                `json:"hit_rate"`
    Size        int                    `json:"size"`
    Items       int                    `json:"items"`
    MemoryUsage int64                  `json:"memory_usage"`
    Evictions   int64                  `json:"evictions,omitempty"`
}

// Quota حد الاستخدام
type Quota struct {
    UserID      string                 `json:"user_id"`
    Type        string                 `json:"type"` // daily, monthly, etc.
    Limit       float64                `json:"limit"`
    Used        float64                `json:"used"`
    ResetAt     time.Time              `json:"reset_at"`
    Period      string                 `json:"period"`
    Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// --- الأنواع المساعدة ---

// VideoOptions خيارات الفيديو (مثال)
type VideoOptions struct {
    Duration     int
    Quality      string
    AspectRatio  string
    Style        string
    OutputFormat string
    FPS          int
    Resolution   string
    Audio        bool
    Watermark    string
    Background   string
}

// TextOptions خيارات النصوص (مثال)
type TextOptions struct {
    Model         string
    Temperature   float64
    MaxTokens     int
    TopP          float64
    TopK          int
    FrequencyPenalty float64
    PresencePenalty  float64
    StopSequences []string
    SystemPrompt  string
    Language      string
}

// ImageOptions خيارات الصور (مثال)
type ImageOptions struct {
    Model         string
    Size          string
    Style         string
    Quality       string
    AspectRatio   string
    NumImages     int
    NegativePrompt string
    Seed          int64
}

// VideoInfo معلومات الفيديو (مثال)
type VideoInfo struct {
    ID          string
    Title       string
    Description string
    URL         string
    Duration    int
    Size        int64
    Status      string
    CreatedAt   time.Time
    UpdatedAt   time.Time
    Cost        float64
    Provider    string
}

// AnalysisOptions خيارات التحليل (مثال)
type AnalysisOptions struct {
    Model         string
    Task          string
    Language      string
    DetailLevel   string
    Format        string
    MaxResults    int
}