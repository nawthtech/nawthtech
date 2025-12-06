package models

// RequestOptions خيارات الطلب (محدث)
type RequestOptions struct {
    Model       string
    Temperature float64
    MaxTokens   int
    Language    string
    ImageData   []byte
}

// Option تحديث الوظيفة لتقبل RequestOptions
type Option struct {
    Model       string
    Temperature float64
    MaxTokens   int
    Language    string
    ImageData   []byte
}

// NewOption إنشاء option جديد
func NewOption() *Option {
    return &Option{
        Model:       "gemini-2.5-flash",
        Temperature: 0.7,
        MaxTokens:   1000,
        Language:    "en",
    }
}

// WithModel تحديث الدالة المساعدة
func WithModel(model string) func(*Option) {
    return func(o *Option) {
        o.Model = model
    }
}

// WithLanguage تحديث الدالة المساعدة
func WithLanguage(lang string) func(*Option) {
    return func(o *Option) {
        o.Language = lang
    }
}