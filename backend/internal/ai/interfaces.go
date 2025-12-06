package ai

import "time"

// VideoServiceInterface واجهة لفصل التبعيات
type VideoServiceInterface interface {
    GenerateVideo(prompt string, options VideoOptions) (*VideoResponse, error)
    GetVideoStatus(operationID string) (*VideoResponse, error)
    CancelVideoGeneration(operationID string) error
    ListVideos() ([]VideoInfo, error)
    DownloadVideo(operationID string) ([]byte, error)
}

// TextServiceInterface واجهة لخدمة النصوص
type TextServiceInterface interface {
    GenerateText(prompt string, options TextOptions) (*TextResponse, error)
    AnalyzeText(text string, options AnalysisOptions) (*AnalysisResult, error)
    TranslateText(text, sourceLang, targetLang string, options TranslationOptions) (*TranslationResult, error)
    SummarizeText(text string, options SummaryOptions) (*SummaryResult, error)
}

// ImageServiceInterface واجهة لخدمة الصور
type ImageServiceInterface interface {
    GenerateImage(prompt string, options ImageOptions) (*ImageResponse, error)
    AnalyzeImage(imageData []byte, options AnalysisOptions) (*AnalysisResult, error)
    EditImage(imageData []byte, prompt string, options EditOptions) (*ImageResponse, error)
    UpscaleImage(imageData []byte, options UpscaleOptions) (*ImageResponse, error)
}

// ProviderInterface واجهة أساسية للمزودين
type ProviderInterface interface {
    // الخصائص الأساسية
    GetName() string
    GetType() string
    IsAvailable() bool
    GetCost() float64
    GetStats() *ProviderStats
    
    // الوظائف الأساسية
    GenerateText(req TextRequest) (*TextResponse, error)
    GenerateImage(req ImageRequest) (*ImageResponse, error)
    GenerateVideo(req VideoRequest) (*VideoResponse, error)
    AnalyzeText(req AnalysisRequest) (*AnalysisResponse, error)
    AnalyzeImage(req AnalysisRequest) (*AnalysisResponse, error)
    TranslateText(req TranslationRequest) (*TranslationResponse, error)
    
    // وظائف إضافية
    SupportsStreaming() bool
    SupportsEmbedding() bool
    GetMaxTokens() int
    GetSupportedLanguages() []string
}

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

// VideoOptions خيارات الفيديو
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

// TextOptions خيارات النصوص
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

// ImageOptions خيارات الصور
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

// AnalysisOptions خيارات التحليل
type AnalysisOptions struct {
    Model         string
    Task          string
    Language      string
    DetailLevel   string
    Format        string
    MaxResults    int
}

// TranslationOptions خيارات الترجمة
type TranslationOptions struct {
    Model       string
    Formality   string
    Context     string
    GlossaryID  string
}

// SummaryOptions خيارات التلخيص
type SummaryOptions struct {
    Model       string
    Length      string // short, medium, long
    Format      string // bullet, paragraph, highlights
    Language    string
    Focus       string
}

// EditOptions خيارات التعديل
type EditOptions struct {
    Model       string
    Strength    float64
    Mask        []byte
    Instructions string
}

// UpscaleOptions خيارات التكبير
type UpscaleOptions struct {
    Scale       int
    Model       string
    Denoise     float64
    Sharpness   float64
}

// VideoInfo معلومات الفيديو
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

// Entity كيان في التحليل
type Entity struct {
    Text        string
    Type        string
    Start       int
    End         int
    Confidence  float64
}

// ProviderError خطأ المزود
type ProviderError struct {
    Code        string
    Message     string
    Timestamp   time.Time
    Retryable   bool
}

// CacheStats إحصائيات الذاكرة المؤقتة
type CacheStats struct {
    Hits        int64
    Misses      int64
    HitRate     float64
    Size        int
    Items       int
    MemoryUsage int64
}

// AnalysisResult نتيجة التحليل
type AnalysisResult struct {
    Result      string
    Confidence  float64
    Entities    []Entity
    Sentiment   string
    Categories  []string
    Cost        float64
    ModelUsed   string
}

// TranslationResult نتيجة الترجمة
type TranslationResult struct {
    TranslatedText string
    SourceText     string
    SourceLanguage string
    TargetLanguage string
    Cost           float64
    ModelUsed      string
}

// SummaryResult نتيجة التلخيص
type SummaryResult struct {
    Summary     string
    OriginalLength int
    SummaryLength  int
    KeyPoints   []string
    Cost        float64
    ModelUsed   string
}

// ModelInfo معلومات النموذج
type ModelInfo struct {
    ID          string
    Name        string
    Type        string
    Provider    string
    MaxTokens   int
    Supports    []string
    CostPerToken float64
    Available   bool
}

// StreamingResponse استجابة التدفق
type StreamingResponse struct {
    Type        string // "chunk", "complete", "error"
    Content     string
    Done        bool
    Error       string
    Tokens      int
}

// EmbeddingRequest طلب التضمين
type EmbeddingRequest struct {
    Text        string
    Model       string
    Dimensions  int
}

// EmbeddingResponse استجابة التضمين
type EmbeddingResponse struct {
    Embedding   []float64
    Dimensions  int
    Model       string
    Cost        float64
}

// VoiceOptions خيارات الصوت
type VoiceOptions struct {
    Voice       string
    Language    string
    Speed       float64
    Pitch       float64
    Volume      float64
    Format      string
}

// VoiceResponse استجابة الصوت
type VoiceResponse struct {
    AudioData   []byte
    Format      string
    Duration    int
    Cost        float64
    Model       string
}

// ConversationMessage رسالة المحادثة
type ConversationMessage struct {
    Role        string // "user", "assistant", "system"
    Content     string
    Timestamp   time.Time
    Metadata    map[string]interface{}
}

// ConversationContext سياق المحادثة
type ConversationContext struct {
    Messages    []ConversationMessage
    MaxTokens   int
    Temperature float64
    Model       string
    MemorySize  int
}

// AIConfig تكوين الذكاء الاصطناعي
type AIConfig struct {
    Providers   []ProviderConfig
    CostManager CostManagerConfig
    Cache       CacheConfig
    Limits      UsageLimits
    Features    FeaturesConfig
}

// ProviderConfig تكوين المزود
type ProviderConfig struct {
    Name        string
    Type        string
    APIKey      string
    BaseURL     string
    Enabled     bool
    Priority    int
    MaxTokens   int
    Timeout     time.Duration
}

// CostManagerConfig تكوين مدير التكاليف
type CostManagerConfig struct {
    MonthlyLimit float64
    DailyLimit   float64
    DataPath     string
    AutoReset    bool
    AlertThreshold float64
}

// CacheConfig تكوين الذاكرة المؤقتة
type CacheConfig struct {
    Enabled     bool
    TTL         time.Duration
    MaxSize     int
    CleanupInterval time.Duration
}

// UsageLimits حدود الاستخدام
type UsageLimits struct {
    FreeTier    UserLimits
    BasicTier   UserLimits
    PremiumTier UserLimits
}

// UserLimits حدود المستخدم
type UserLimits struct {
    MonthlyCost  float64
    DailyRequests int
    MaxTokens    int
    MaxImages    int
    MaxVideos    int
}

// FeaturesConfig تكوين الميزات
type FeaturesConfig struct {
    TextGeneration  bool
    ImageGeneration bool
    VideoGeneration bool
    TextAnalysis    bool
    Translation     bool
    VoiceSynthesis  bool
    Embeddings      bool
    Streaming       bool
}

// ErrorResponse استجابة الخطأ
type ErrorResponse struct {
    Code        string
    Message     string
    Details     map[string]interface{}
    Timestamp   time.Time
    Suggestion  string
}

// HealthStatus حالة النظام
type HealthStatus struct {
    Status      string
    Timestamp   time.Time
    Components  map[string]ComponentStatus
    Uptime      time.Duration
    Version     string
}

// ComponentStatus حالة المكون
type ComponentStatus struct {
    Status      string
    Latency     float64
    LastCheck   time.Time
    Error       string
}

// MonitorMetric مقياس المراقبة
type MonitorMetric struct {
    Name        string
    Value       float64
    Timestamp   time.Time
    Labels      map[string]string
    Type        string
}

// RateLimiterConfig تكوين محدد المعدل
type RateLimiterConfig struct {
    Enabled     bool
    RequestsPerSecond int
    Burst       int
    PerUser     bool
}

// AuditLog سجل التدقيق
type AuditLog struct {
    ID          string
    UserID      string
    Action      string
    Resource    string
    Timestamp   time.Time
    IPAddress   string
    UserAgent   string
    Metadata    map[string]interface{}
    Status      string
}