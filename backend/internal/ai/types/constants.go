package types

// ثوابت الأنواع
const (
    // أنواع المزودين
    ProviderTypeText   = "text"
    ProviderTypeImage  = "image"
    ProviderTypeVideo  = "video"
    ProviderTypeMulti  = "multi"
    
    // أسماء المزودين
    ProviderNameGemini      = "gemini"
    ProviderNameOllama      = "ollama"
    ProviderNameHuggingFace = "huggingface"
    ProviderNameMulti       = "multi"
    
    // أنواع العمليات
    OperationTypeText        = "text"
    OperationTypeImage       = "image"
    OperationTypeVideo       = "video"
    OperationTypeAnalysis    = "analysis"
    OperationTypeTranslation = "translation"
    OperationTypeEmbedding   = "embedding"
    OperationTypeChat        = "chat"
    
    // طبقات المستخدمين
    UserTierFree      = "free"
    UserTierBasic     = "basic"
    UserTierPremium   = "premium"
    UserTierEnterprise = "enterprise"
    
    // حالات الفيديو
    VideoStatusPending   = "pending"
    VideoStatusProcessing = "processing"
    VideoStatusCompleted = "completed"
    VideoStatusFailed    = "failed"
    
    // أسباب إنهاء النص
    FinishReasonStop    = "stop"
    FinishReasonLength  = "length"
    FinishReasonContent = "content_filter"
    FinishReasonTimeout = "timeout"
    
    // أحجام الصور
    ImageSizeSmall  = "256x256"
    ImageSizeMedium = "512x512"
    ImageSizeLarge  = "1024x1024"
    
    // دقة الفيديو
    VideoResolutionSD   = "640x360"
    VideoResolutionHD   = "1280x720"
    VideoResolutionFullHD = "1920x1080"
    
    // جودة الصور
    ImageQualityStandard = "standard"
    ImageQualityHD       = "hd"
    
    // أنماط الصور
    ImageStyleNatural  = "natural"
    ImageStyleVivid    = "vivid"
    
    // لغات
    LanguageArabic     = "ar"
    LanguageEnglish    = "en"
    LanguageSpanish    = "es"
    LanguageFrench     = "fr"
    LanguageGerman     = "de"
    LanguageChinese    = "zh"
    LanguageJapanese   = "ja"
    LanguageKorean     = "ko"
    LanguageRussian    = "ru"
    
    // حالات الصحة
    HealthStatusHealthy   = "healthy"
    HealthStatusUnhealthy = "unhealthy"
    HealthStatusDegraded  = "degraded"
)

// ModelNames أسماء النماذج الشائعة
var ModelNames = map[string][]string{
    ProviderNameGemini: {
        "gemini-2.5-flash-exp",
        "gemini-2.0-flash-exp",
        "gemini-1.5-pro",
        "gemini-1.5-flash",
    },
    ProviderNameOllama: {
        "llama3.2:3b",
        "llama3.1:8b",
        "mistral",
        "mixtral",
        "codellama",
    },
    ProviderNameHuggingFace: {
        "gpt2",
        "bert",
        "t5",
        "flan-t5",
    },
}

// DefaultConfigs التكوينات الافتراضية
var (
    DefaultTextTemperature = 0.7
    DefaultMaxTokens       = 2048
    DefaultImageSize       = ImageSizeMedium
    DefaultImageQuality    = ImageQualityStandard
    DefaultVideoDuration   = 30 // ثانية
    DefaultVideoResolution = VideoResolutionHD
    DefaultVideoFPS        = 30
)