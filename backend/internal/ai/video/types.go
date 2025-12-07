package video

import (
    "time"
)

// VideoProvider واجهة لمزود خدمة الفيديو
type VideoProvider interface {
    GenerateVideo(req VideoRequest) (*VideoResponse, error)
    Name() string
    IsAvailable() bool
    IsLocal() bool
    IsFree() bool
    SupportsResolution(resolution string) bool
}

// VideoRequest طلب توليد فيديو
type VideoRequest struct {
    Prompt         string            `json:"prompt"`
    NegativePrompt string            `json:"negative_prompt,omitempty"`
    Duration       int               `json:"duration"`        // بالثواني
    Resolution     string            `json:"resolution"`      // مثال: "1920x1080"
    Aspect         string            `json:"aspect,omitempty"` // مثال: "16:9", "9:16", "1:1"
    Style          string            `json:"style,omitempty"`  // فني، واقعي، كرتوني، إلخ.
    Options        VideoOptions      `json:"options,omitempty"`
    UserID         string            `json:"user_id,omitempty"`
    UserTier       string            `json:"user_tier,omitempty"`
}

// VideoOptions خيارات الفيديو
type VideoOptions struct {
    FPS         int     `json:"fps,omitempty"`
    Seed        int64   `json:"seed,omitempty"`
    CFGScale    float64 `json:"cfg_scale,omitempty"`
    Steps       int     `json:"steps,omitempty"`
    Model       string  `json:"model,omitempty"`
    Quality     string  `json:"quality,omitempty"` // low, medium, high
    MotionScale float64 `json:"motion_scale,omitempty"`
}

// VideoResponse استجابة توليد الفيديو
type VideoResponse struct {
    Success     bool        `json:"success"`
    VideoURL    string      `json:"video_url,omitempty"`
    VideoData   []byte      `json:"-"`
    Duration    int         `json:"duration"`
    Width       int         `json:"width,omitempty"`
    Height      int         `json:"height,omitempty"`
    Resolution  string      `json:"resolution"`
    Aspect      string      `json:"aspect,omitempty"`
    Format      string      `json:"format,omitempty"` // mp4, webm, gif
    Provider    string      `json:"provider"`
    Cost        float64     `json:"cost,omitempty"`
    Status      string      `json:"status"`
    Error       string      `json:"error,omitempty"`
    Timestamp   int64       `json:"timestamp"`
    CreatedAt   time.Time   `json:"created_at"`
}

// VideoJobStatus حالة مهمة الفيديو
type VideoJobStatus string

const (
    VideoJobPending    VideoJobStatus = "pending"
    VideoJobProcessing VideoJobStatus = "processing"
    VideoJobCompleted  VideoJobStatus = "completed"
    VideoJobFailed     VideoJobStatus = "failed"
    VideoJobCancelled  VideoJobStatus = "cancelled"
)

// VideoQuality جودة الفيديو
type VideoQuality string

const (
    QualityLow    VideoQuality = "low"
    QualityMedium VideoQuality = "medium"
    QualityHigh   VideoQuality = "high"
)

// VideoAspect نسبة العرض إلى الارتفاع
type VideoAspect string

const (
    Aspect16_9  VideoAspect = "16:9"
    Aspect9_16  VideoAspect = "9:16"
    Aspect1_1   VideoAspect = "1:1"
    Aspect4_3   VideoAspect = "4:3"
    Aspect21_9  VideoAspect = "21:9"
)

// VideoStyle نمط الفيديو
type VideoStyle string

const (
    StyleRealistic   VideoStyle = "realistic"
    StyleAnime       VideoStyle = "anime"
    StyleCartoon     VideoStyle = "cartoon"
    StyleArtistic    VideoStyle = "artistic"
    StyleCinematic   VideoStyle = "cinematic"
    StyleMinimal     VideoStyle = "minimal"
)

// VideoFormat صيغة الفيديو
type VideoFormat string

const (
    FormatMP4  VideoFormat = "mp4"
    FormatWEBM VideoFormat = "webm"
    FormatGIF  VideoFormat = "gif"
    FormatMOV  VideoFormat = "mov"
)

// ResolutionToDimensions تحويل دقة الفيديو إلى أبعاد
func ResolutionToDimensions(resolution string) (width, height int, ok bool) {
    switch resolution {
    case "1920x1080":
        return 1920, 1080, true
    case "1080x1920":
        return 1080, 1920, true
    case "1280x720":
        return 1280, 720, true
    case "720x1280":
        return 720, 1280, true
    case "1024x1024":
        return 1024, 1024, true
    case "512x512":
        return 512, 512, true
    case "256x256":
        return 256, 256, true
    default:
        return 0, 0, false
    }
}

// AspectToDimensions تحويل النسبة إلى أبعاد
func AspectToDimensions(aspect VideoAspect, height int) (width int) {
    switch aspect {
    case Aspect16_9:
        return height * 16 / 9
    case Aspect9_16:
        return height * 9 / 16
    case Aspect1_1:
        return height
    case Aspect4_3:
        return height * 4 / 3
    case Aspect21_9:
        return height * 21 / 9
    default:
        return height * 16 / 9 // افتراضي 16:9
    }
}

// ValidateVideoRequest التحقق من صحة طلب الفيديو
func ValidateVideoRequest(req VideoRequest) error {
    if req.Prompt == "" {
        return ErrEmptyPrompt
    }
    
    if req.Duration <= 0 || req.Duration > 60 {
        return ErrInvalidDuration
    }
    
    if !IsSupportedResolution(req.Resolution) {
        return ErrUnsupportedResolution
    }
    
    if req.Aspect != "" {
        if !IsSupportedAspect(VideoAspect(req.Aspect)) {
            return ErrUnsupportedAspect
        }
    }
    
    if req.Style != "" {
        if !IsSupportedStyle(VideoStyle(req.Style)) {
            return ErrUnsupportedStyle
        }
    }
    
    return nil
}

// IsSupportedResolution التحقق من دعم الدقة
func IsSupportedResolution(resolution string) bool {
    supported := []string{
        "1920x1080", "1080x1920",
        "1280x720", "720x1280",
        "1024x1024", "512x512", "256x256",
    }
    
    for _, res := range supported {
        if res == resolution {
            return true
        }
    }
    return false
}

// IsSupportedAspect التحقق من دعم النسبة
func IsSupportedAspect(aspect VideoAspect) bool {
    switch aspect {
    case Aspect16_9, Aspect9_16, Aspect1_1, Aspect4_3, Aspect21_9:
        return true
    default:
        return false
    }
}

// IsSupportedStyle التحقق من دعم النمط
func IsSupportedStyle(style VideoStyle) bool {
    switch style {
    case StyleRealistic, StyleAnime, StyleCartoon, StyleArtistic, StyleCinematic, StyleMinimal:
        return true
    default:
        return false
    }
}

// أخطاء الفيديو
var (
    ErrEmptyPrompt            = &VideoError{"empty_prompt", "Prompt cannot be empty"}
    ErrInvalidDuration        = &VideoError{"invalid_duration", "Duration must be between 1 and 60 seconds"}
    ErrUnsupportedResolution  = &VideoError{"unsupported_resolution", "Unsupported video resolution"}
    ErrUnsupportedAspect      = &VideoError{"unsupported_aspect", "Unsupported aspect ratio"}
    ErrUnsupportedStyle       = &VideoError{"unsupported_style", "Unsupported video style"}
    ErrProviderUnavailable    = &VideoError{"provider_unavailable", "Video provider is currently unavailable"}
    ErrGenerationFailed       = &VideoError{"generation_failed", "Video generation failed"}
    ErrQuotaExceeded         = &VideoError{"quota_exceeded", "Video generation quota exceeded"}
)

// VideoError خطأ الفيديو
type VideoError struct {
    Code    string `json:"code"`
    Message string `json:"message"`
}

func (e *VideoError) Error() string {
    return e.Message
}

// VideoStats إحصائيات الفيديو
type VideoStats struct {
    TotalGenerations  int64     `json:"total_generations"`
    Successful        int64     `json:"successful"`
    Failed            int64     `json:"failed"`
    TotalDuration     int64     `json:"total_duration"` // إجمالي الثواني
    TotalCost         float64   `json:"total_cost"`
    LastGeneration    time.Time `json:"last_generation"`
    MostUsedStyle     string    `json:"most_used_style"`
    MostUsedProvider  string    `json:"most_used_provider"`
}

// VideoUsage استخدام الفيديو للمستخدم
type VideoUsage struct {
    UserID           string    `json:"user_id"`
    Tier             string    `json:"tier"`
    TotalGenerations int       `json:"total_generations"`
    MonthlyLimit     int       `json:"monthly_limit"`
    MonthlyUsed      int       `json:"monthly_used"`
    DailyLimit       int       `json:"daily_limit"`
    DailyUsed        int       `json:"daily_used"`
    LastGenerated    time.Time `json:"last_generated"`
    LastReset        time.Time `json:"last_reset"`
}

// CanGenerateVideo التحقق من إمكانية توليد فيديو
func (u *VideoUsage) CanGenerateVideo() (bool, string) {
    now := time.Now()
    
    // التحقق من الحد الشهري
    if u.LastReset.Month() != now.Month() || u.LastReset.Year() != now.Year() {
        u.MonthlyUsed = 0
        u.LastReset = now
    }
    
    if u.MonthlyUsed >= u.MonthlyLimit {
        return false, "Monthly video generation limit exceeded"
    }
    
    // التحقق من الحد اليومي
    if u.LastGenerated.Day() != now.Day() || 
       u.LastGenerated.Month() != now.Month() || 
       u.LastGenerated.Year() != now.Year() {
        u.DailyUsed = 0
    }
    
    if u.DailyUsed >= u.DailyLimit {
        return false, "Daily video generation limit exceeded"
    }
    
    return true, ""
}

// RecordGeneration تسجيل عملية توليد
func (u *VideoUsage) RecordGeneration() {
    u.TotalGenerations++
    u.MonthlyUsed++
    u.DailyUsed++
    u.LastGenerated = time.Now()
}

// GetDefaultLimits الحصول على الحدود الافتراضية حسب الطبقة
func GetDefaultLimits(tier string) (monthly, daily int) {
    switch tier {
    case "premium":
        return 100, 10
    case "basic":
        return 30, 3
    default: // free
        return 10, 1
    }
}