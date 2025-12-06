package ai

// RequestOptions خيارات الطلب
type RequestOptions struct {
    Model       string
    Temperature float64
    MaxTokens   int
    Language    string
    ImageData   []byte
}

// VideoRequest طلب توليد فيديو
type VideoRequest struct {
    Prompt      string
    Model       string
    Duration    int    // مدة الفيديو بالثواني
    AspectRatio string // نسبة الأبعاد
    Style       string // النمط
}

// VideoResponse استجابة توليد فيديو
type VideoResponse struct {
    VideoURL    string
    Duration    int
    Cost        float64
    Status      string
    OperationID string
}

// TextRequest طلب توليد نص
type TextRequest struct {
    Prompt      string
    Model       string
    Temperature float64
    MaxTokens   int
    Language    string
}

// TextResponse استجابة توليد نص
type TextResponse struct {
    Text        string
    Tokens      int
    Cost        float64
    Model       string
}

// ImageRequest طلب توليد صورة
type ImageRequest struct {
    Prompt      string
    Model       string
    Size        string
    Style       string
    Quality     string
}

// ImageResponse استجابة توليد صورة
type ImageResponse struct {
    ImageURL    string
    ImageData   []byte
    Size        string
    Cost        float64
    Model       string
}

// AnalysisRequest طلب تحليل
type AnalysisRequest struct {
    Text        string
    ImageData   []byte
    Prompt      string
    Model       string
}

// AnalysisResponse استجابة تحليل
type AnalysisResponse struct {
    Result      string
    Confidence  float64
    Cost        float64
    Model       string
}

// TranslationRequest طلب ترجمة
type TranslationRequest struct {
    Text        string
    FromLang    string
    ToLang      string
    Model       string
}

// TranslationResponse استجابة ترجمة
type TranslationResponse struct {
    TranslatedText string
    Cost           float64
    Model          string
}

// ProviderStats إحصائيات المزود
type ProviderStats struct {
    Name         string
    Type         string
    IsAvailable  bool
    Requests     int64
    Errors       int64
    LastUsed     string
    TotalCost    float64
    SuccessRate  float64
}

// Provider واجهة مشتركة محدثة
type Provider interface {
    // توليد نص
    GenerateText(req TextRequest) (*TextResponse, error)
    
    // توليد صور
    GenerateImage(req ImageRequest) (*ImageResponse, error)
    
    // توليد فيديوهات
    GenerateVideo(req VideoRequest) (*VideoResponse, error)
    
    // تحليل صور
    AnalyzeImage(req AnalysisRequest) (*AnalysisResponse, error)
    
    // تحليل نص
    AnalyzeText(req AnalysisRequest) (*AnalysisResponse, error)
    
    // ترجمة نص
    TranslateText(req TranslationRequest) (*TranslationResponse, error)
    
    // تحقق من التوفر
    IsAvailable() bool
    
    // الحصول على التكلفة
    GetCost() float64
    
    // الحصول على اسم المزود
    GetName() string
    
    // الحصول على الإحصائيات (اختياري)
    GetStats() *ProviderStats
}

// Option خيارات إضافية
type Option func(*RequestOptions)

// WithModel تحديد النموذج
func WithModel(model string) Option {
    return func(o *RequestOptions) {
        o.Model = model
    }
}

// WithLanguage تحديد اللغة
func WithLanguage(lang string) Option {
    return func(o *RequestOptions) {
        o.Language = lang
    }
}

// WithImage إضافة صورة للتحليل
func WithImage(imageData []byte) Option {
    return func(o *RequestOptions) {
        o.ImageData = imageData
    }
}

// WithTemperature تحديد درجة الحرارة
func WithTemperature(temp float64) Option {
    return func(o *RequestOptions) {
        o.Temperature = temp
    }
}

// WithMaxTokens تحديد الحد الأقصى للرموز
func WithMaxTokens(tokens int) Option {
    return func(o *RequestOptions) {
        o.MaxTokens = tokens
    }
}

// OptionHandler معالج الخيارات
type OptionHandler interface {
    Apply(options *RequestOptions)
}

// applyOptions تطبيق الخيارات
func applyOptions(opts []Option, defaultOpts *RequestOptions) *RequestOptions {
    if defaultOpts == nil {
        defaultOpts = &RequestOptions{
            Temperature: 0.7,
            MaxTokens:   1024,
            Language:    "ar",
        }
    }
    
    for _, opt := range opts {
        opt(defaultOpts)
    }
    
    return defaultOpts
}

// ProviderBase قاعدة للمزودين
type ProviderBase struct {
    Name      string
    BaseCost  float64
    Available bool
}

// NewProviderBase إنشاء مزود قاعدة
func NewProviderBase(name string, baseCost float64) *ProviderBase {
    return &ProviderBase{
        Name:      name,
        BaseCost:  baseCost,
        Available: true,
    }
}

// GetName الحصول على اسم المزود
func (p *ProviderBase) GetName() string {
    return p.Name
}

// GetCost الحصول على التكلفة الأساسية
func (p *ProviderBase) GetCost() float64 {
    return p.BaseCost
}

// IsAvailable التحقق من التوفر
func (p *ProviderBase) IsAvailable() bool {
    return p.Available
}

// SetAvailability تعيين حالة التوفر
func (p *ProviderBase) SetAvailability(available bool) {
    p.Available = available
}

// CalculateCost حساب التكلفة
func (p *ProviderBase) CalculateCost(tokens int, complexity float64) float64 {
    return p.BaseCost * float64(tokens) * complexity
}

// GetStatsDefault الحصول على إحصائيات افتراضية
func (p *ProviderBase) GetStatsDefault() *ProviderStats {
    return &ProviderStats{
        Name:        p.Name,
        Type:        "unknown",
        IsAvailable: p.Available,
        Requests:    0,
        Errors:      0,
        LastUsed:    "",
        TotalCost:   0,
        SuccessRate: 100.0,
    }
}