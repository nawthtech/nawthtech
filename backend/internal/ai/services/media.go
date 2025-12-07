package services

import (
    "context"
    "fmt"
    "time"
    "github.com/nawthtech/nawthtech/backend/internal/ai/types"
)

type MediaService struct {
    imageProvider types.ImageProvider
    videoProvider types.VideoProvider
}

func NewMediaService(imageProvider types.ImageProvider, videoProvider types.VideoProvider) *MediaService {
    return &MediaService{
        imageProvider: imageProvider,
        videoProvider: videoProvider,
    }
}

func (s *MediaService) GenerateSocialMediaImage(ctx context.Context, platform string, prompt string, style string) (*types.ImageResponse, error) {
    // تحديد حجم الصورة بناءً على المنصة
    size := "1024x1024" // حجم افتراضي
    switch platform {
    case "instagram":
        size = "1080x1080"
    case "facebook":
        size = "1200x630"
    case "twitter":
        size = "1200x675"
    case "linkedin":
        size = "1200x627"
    case "pinterest":
        size = "1000x1500"
    }
    
    req := types.ImageRequest{
        Prompt:  fmt.Sprintf("%s - %s style for %s platform", prompt, style, platform),
        Size:    size,
        Style:   style,
        UserID:  extractUserIDFromContext(ctx),
        UserTier: extractUserTierFromContext(ctx),
        N: variations,
    }
    
    return s.imageProvider.GenerateImage(req)
}

func (s *MediaService) GenerateMarketingImage(ctx context.Context, product string, targetAudience string, theme string) (*types.ImageResponse, error) {
    prompt := fmt.Sprintf("Marketing image for %s targeting %s with %s theme", 
        product, targetAudience, theme)
    
    req := types.ImageRequest{
        Prompt:  prompt,
        Size:    "1200x628", // حجم قياسي للتسويق
        Quality: "hd",
        UserID:  extractUserIDFromContext(ctx),
        UserTier: extractUserTierFromContext(ctx),
    }
    
    return s.imageProvider.GenerateImage(req)
}

func (s *MediaService) GenerateShortVideo(ctx context.Context, prompt string, duration int, platform string) (*types.VideoResponse, error) {
    // تحديد الإعدادات بناءً على المنصة والمدة
    var resolution string
    switch platform {
    case "tiktok":
        resolution = "1080x1920" // عمودي
    case "youtube":
        resolution = "1920x1080" // أفقي
    case "instagram":
        resolution = "1080x1920" // ريلز
    default:
        resolution = "1080x1920"
    }
    
    if duration <= 0 {
        duration = 15 // ثواني افتراضية
    }
    
    req := types.VideoRequest{
        Prompt:     fmt.Sprintf("%s - optimized for %s platform", prompt, platform),
        Duration:   duration,
        Resolution: resolution,
        UserID:     extractUserIDFromContext(ctx),
        UserTier:   extractUserTierFromContext(ctx),
    }
    
    return s.videoProvider.GenerateVideo(req)
}

func (s *MediaService) GenerateTutorialVideo(ctx context.Context, title string, steps []string, style string) (*types.VideoResponse, error) {
    prompt := fmt.Sprintf("Tutorial video: %s\n\nSteps:\n", title)
    for i, step := range steps {
        prompt += fmt.Sprintf("%d. %s\n", i+1, step)
    }
    prompt += fmt.Sprintf("\nStyle: %s", style)
    
    req := types.VideoRequest{
        Prompt:     prompt,
        Duration:   60, // دقيقة واحدة
        Resolution: "1920x1080",
        Style:      style,
        UserID:     extractUserIDFromContext(ctx),
        UserTier:   extractUserTierFromContext(ctx),
    }
    
    return s.videoProvider.GenerateVideo(req)
}

func (s *MediaService) GenerateAdVideo(ctx context.Context, product string, keyPoints []string, callToAction string) (*types.VideoResponse, error) {
    prompt := fmt.Sprintf("Advertisement video for %s\n\nKey points:\n", product)
    for _, point := range keyPoints {
        prompt += fmt.Sprintf("- %s\n", point)
    }
    prompt += fmt.Sprintf("\nCall to action: %s", callToAction)
    
    req := types.VideoRequest{
        Prompt:     prompt,
        Duration:   30, // 30 ثانية للإعلان
        Resolution: "1080x1920", // عمودي للهواتف
        Style:      "cinematic",
        UserID:     extractUserIDFromContext(ctx),
        UserTier:   extractUserTierFromContext(ctx),
    }
    
    return s.videoProvider.GenerateVideo(req)
}

func (s *MediaService) GetVideoStatus(ctx context.Context, operationID string) (*types.VideoResponse, error) {
    // هذه الدالة تتطلب أن يكون videoProvider يطبق GetVideoStatus
    // نحتاج إلى type assertion للتحقق
    if provider, ok := s.videoProvider.(interface{ GetVideoStatus(string) (*types.VideoResponse, error) }); ok {
        return provider.GetVideoStatus(operationID)
    }
    
    return nil, fmt.Errorf("video provider does not support status checking")
}

func (s *MediaService) BatchGenerateImages(ctx context.Context, prompts []string, size string, style string) ([]*types.ImageResponse, error) {
    var responses []*types.ImageResponse
    
    for _, prompt := range prompts {
        req := types.ImageRequest{
            Prompt:  prompt,
            Size:    size,
            Style:   style,
            UserID:  extractUserIDFromContext(ctx),
            UserTier: extractUserTierFromContext(ctx),
        }
        
        resp, err := s.imageProvider.GenerateImage(req)
        if err != nil {
            // يمكنك اختيار إما إرجاع الأخطاء أو الاستمرار
            // هنا سنستمر ونعيد الأخطاء في النهاية
            return responses, fmt.Errorf("failed to generate image for prompt '%s': %w", prompt, err)
        }
        
        responses = append(responses, resp)
        
        // تأخير بسيط بين الطلبات لتجنب rate limiting
        select {
        case <-ctx.Done():
            return responses, ctx.Err()
        case <-time.After(100 * time.Millisecond):
            // استمر
        }
    }
    
    return responses, nil
}

func (s *MediaService) AnalyzeImage(ctx context.Context, imageData []byte, analysisType string) (*types.AnalysisResponse, error) {
    req := types.AnalysisRequest{
        ImageData: imageData,
        Prompt:    fmt.Sprintf("Analyze this image for %s", analysisType),
        UserID:    extractUserIDFromContext(ctx),
        UserTier:  extractUserTierFromContext(ctx),
    }
    
    // تحتاج إلى imageProvider يدعم AnalyzeImage
    if provider, ok := s.imageProvider.(interface{ AnalyzeImage(types.AnalysisRequest) (*types.AnalysisResponse, error) }); ok {
        return provider.AnalyzeImage(req)
    }
    
    return nil, fmt.Errorf("image provider does not support image analysis")
}

func (s *MediaService) GenerateVariations(ctx context.Context, style string) ([]*types.ImageResponse, error) {
    // هذه وظيفة متقدمة قد لا تدعمها جميع المزودين
    prompt := fmt.Sprintf("Generate %d variations of this image in %s style", variations, style)
    
    req := types.ImageRequest{
        Prompt:       prompt,
        Style:        style,
        UserID:       extractUserIDFromContext(ctx),
        UserTier:     extractUserTierFromContext(ctx),
    }
    
    // نحتاج إلى معرفة إذا كان المزود يدعم NumImages
    resp, err := s.imageProvider.GenerateImage(req)
    if err != nil {
        return nil, err
    }
    
    // إذا كان المزود يدعم إعادة مصفوفة من الصور
    return []*types.ImageResponse{resp}, nil
}

func (s *MediaService) UpscaleImage(ctx context.Context, imageData []byte, scaleFactor int) (*types.ImageResponse, error) {
    prompt := "Upscale this image while maintaining quality and details"
    
    req := types.ImageRequest{
        Prompt:    prompt,
        ImageData: imageData,
        Size:      fmt.Sprintf("%dx", scaleFactor), // مثال: "2x" أو "4x"
        Quality:   "highest",
        UserID:    extractUserIDFromContext(ctx),
        UserTier:  extractUserTierFromContext(ctx),
    }
    
    return s.imageProvider.GenerateImage(req)
}

func (s *MediaService) RemoveBackground(ctx context.Context, imageData []byte) (*types.ImageResponse, error) {
    prompt := "Remove background from this image, keep only the main subject with transparent background"
    
    req := types.ImageRequest{
        Prompt:       prompt,
        ImageData:    imageData,
        ResponseFormat: "png", // لحفظ الشفافية
        UserID:       extractUserIDFromContext(ctx),
        UserTier:     extractUserTierFromContext(ctx),
    }
    
    return s.imageProvider.GenerateImage(req)
}

func (s *MediaService) GenerateThumbnail(ctx context.Context, videoDescription string, style string) (*types.ImageResponse, error) {
    prompt := fmt.Sprintf("YouTube thumbnail for video about: %s. Style: %s. Eye-catching, high click-through rate", 
        videoDescription, style)
    
    req := types.ImageRequest{
        Prompt:  prompt,
        Size:    "1280x720", // حجم ثامنة يوتيوب قياسي
        Style:   style,
        UserID:  extractUserIDFromContext(ctx),
        UserTier: extractUserTierFromContext(ctx),
    }
    
    return s.imageProvider.GenerateImage(req)
}

func (s *MediaService) GetServiceStats(ctx context.Context) map[string]interface{} {
    stats := make(map[string]interface{})
    
    // الحصول على إحصائيات من المزودين إذا كانت متوفرة
    if imgProvider, ok := s.imageProvider.(interface{ GetStats() *types.ProviderStats }); ok {
        stats["image_provider"] = imgProvider.GetStats()
    }
    
    if vidProvider, ok := s.videoProvider.(interface{ GetStats() *types.ProviderStats }); ok {
        stats["video_provider"] = vidProvider.GetStats()
    }
    
    stats["service"] = "media"
    stats["timestamp"] = time.Now()
    
    return stats
}