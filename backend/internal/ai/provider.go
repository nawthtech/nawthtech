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