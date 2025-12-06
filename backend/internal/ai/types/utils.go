package types

import (
    "encoding/json"
    "fmt"
    "strconv"
    "strings"
    "time"
)

// Helper functions

// FormatCost تنسيق التكلفة
func FormatCost(cost float64) string {
    if cost < 0.0001 {
        return "≈ $0.0000"
    }
    return fmt.Sprintf("$%.4f", cost)
}

// CalculateTokens حساب عدد الرموز تقريبياً
func CalculateTokens(text string) int {
    // تقدير تقريبي: 1 رمز = 4 أحرف للغة الإنجليزية
    return len(text) / 4
}

// ParseSize تحليل حجم الصورة
func ParseSize(size string) (width, height int, err error) {
    if size == "" {
        return 512, 512, nil
    }
    
    parts := strings.Split(size, "x")
    if len(parts) != 2 {
        return 0, 0, fmt.Errorf("invalid size format: %s", size)
    }
    
    width, err1 := strconv.Atoi(parts[0])
    height, err2 := strconv.Atoi(parts[1])
    
    if err1 != nil || err2 != nil {
        return 0, 0, fmt.Errorf("invalid size numbers: %s", size)
    }
    
    return width, height, nil
}

// CalculateImageCost حساب تكلفة الصورة
func CalculateImageCost(size, quality string) float64 {
    baseCost := 0.02 // $0.02 لكل صورة
    
    switch size {
    case ImageSizeLarge:
        baseCost *= 2
    case ImageSizeSmall:
        baseCost *= 0.5
    }
    
    if quality == ImageQualityHD {
        baseCost *= 1.5
    }
    
    return baseCost
}

// CalculateVideoCost حساب تكلفة الفيديو
func CalculateVideoCost(duration int, resolution string) float64 {
    baseCostPerSecond := 0.001 // $0.001 لكل ثانية
    
    cost := float64(duration) * baseCostPerSecond
    
    switch resolution {
    case VideoResolutionFullHD:
        cost *= 1.5
    case VideoResolutionHD:
        cost *= 1.2
    }
    
    return cost
}

// IsValidLanguage التحقق من صحة اللغة
func IsValidLanguage(lang string) bool {
    validLanguages := []string{
        LanguageArabic, LanguageEnglish, LanguageSpanish,
        LanguageFrench, LanguageGerman, LanguageChinese,
        LanguageJapanese, LanguageKorean, LanguageRussian,
    }
    
    for _, validLang := range validLanguages {
        if strings.ToLower(lang) == validLang {
            return true
        }
    }
    
    return false
}

// GetDefaultModel الحصول على النموذج الافتراضي للمزود
func GetDefaultModel(provider string) string {
    switch strings.ToLower(provider) {
    case ProviderNameGemini:
        return "gemini-2.5-flash-exp"
    case ProviderNameOllama:
        return "llama3.2:3b"
    case ProviderNameHuggingFace:
        return "gpt2"
    default:
        return "auto"
    }
}

// MergeMetadata دمج البيانات الوصفية
func MergeMetadata(base, additional map[string]interface{}) map[string]interface{} {
    result := make(map[string]interface{})
    
    for k, v := range base {
        result[k] = v
    }
    
    for k, v := range additional {
        result[k] = v
    }
    
    return result
}

// FormatDuration تنسيق المدة
func FormatDuration(seconds int) string {
    minutes := seconds / 60
    remainingSeconds := seconds % 60
    
    if minutes > 0 {
        return fmt.Sprintf("%dm %ds", minutes, remainingSeconds)
    }
    return fmt.Sprintf("%ds", seconds)
}

// CopyRequest نسخ الطلب
func CopyRequest[T any](req T) T {
    // استخدام JSON للنسخ العميق
    data, _ := json.Marshal(req)
    var copy T
    json.Unmarshal(data, &copy)
    return copy
}

// GenerateID توليد معرف فريد
func GenerateID(prefix string) string {
    timestamp := time.Now().UnixNano()
    return fmt.Sprintf("%s_%d", prefix, timestamp)
}

// CalculateETA حساب الوقت المتوقع للانتهاء
func CalculateETA(startTime time.Time, progress float64) time.Time {
    if progress <= 0 {
        return time.Now().Add(5 * time.Minute) // تقدير افتراضي
    }
    
    elapsed := time.Since(startTime)
    estimatedTotal := time.Duration(float64(elapsed) / progress)
    return startTime.Add(estimatedTotal)
}