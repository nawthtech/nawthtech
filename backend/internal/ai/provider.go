package ai

// Provider واجهة مشتركة محدثة
type Provider interface {
    // توليد نص
    GenerateText(prompt string, opts ...Option) (string, error)
    
    // توليد صور
    GenerateImage(prompt string, opts ...Option) ([]byte, error)
    
    // توليد فيديوهات ⬅️ جديد
    GenerateVideo(req VideoRequest) (*VideoResponse, error)
    
    // تحليل صور
    AnalyzeImage(imageData []byte, prompt string) (string, error)
    
    // تحقق من التوفر
    IsAvailable() bool
    
    // الحصول على التكلفة
    GetCost() float64
    
    // الحصول على اسم المزود
    GetName() string
}


// Option خيارات إضافية
type Option func(*RequestOptions)

// RequestOptions خيارات الطلب
type RequestOptions struct {
    Model       string
    Temperature float64
    MaxTokens   int
    Language    string
    ImageData   []byte
}

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