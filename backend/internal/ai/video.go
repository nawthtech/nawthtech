package ai

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
    "strings"
    "time"
)

// VideoProvider مزود خاص لتوليد الفيديوهات
type VideoProvider struct {
    baseURL    string
    httpClient *http.Client
    apiKey     string
    apiType    string // "gemini", "luma", "runway", "pika"
}

// NewVideoProvider إنشاء مزود فيديوهات جديد
func NewVideoProvider(apiType string) (*VideoProvider, error) {
    var apiKey string
    var baseURL string
    
    switch apiType {
    case "gemini":
        apiKey = os.Getenv("GEMINI_API_KEY")
        if apiKey == "" {
            return nil, fmt.Errorf("GEMINI_API_KEY environment variable is required for Gemini video provider")
        }
        baseURL = "https://generativelanguage.googleapis.com/v1beta"
        
    case "luma":
        apiKey = os.Getenv("LUMA_API_KEY")
        if apiKey == "" {
            return nil, fmt.Errorf("LUMA_API_KEY environment variable is required for Luma video provider")
        }
        baseURL = "https://api.lumalabs.ai/v1"
        
    case "runway":
        apiKey = os.Getenv("RUNWAY_API_KEY")
        if apiKey == "" {
            return nil, fmt.Errorf("RUNWAY_API_KEY environment variable is required for Runway video provider")
        }
        baseURL = "https://api.runwayml.com/v1"
        
    case "pika":
        apiKey = os.Getenv("PIKA_API_KEY")
        if apiKey == "" {
            return nil, fmt.Errorf("PIKA_API_KEY environment variable is required for Pika video provider")
        }
        baseURL = "https://api.pika.art/v1"
        
    default:
        return nil, fmt.Errorf("unsupported video provider type: %s", apiType)
    }
    
    provider := &VideoProvider{
        baseURL: baseURL,
        httpClient: &http.Client{
            Timeout: 600 * time.Second, // 10 دقائق للفيديوهات
        },
        apiKey:  apiKey,
        apiType: apiType,
    }
    
    return provider, nil
}

// GenerateVideo توليد فيديو باستخدام API المختار
func (p *VideoProvider) GenerateVideo(req VideoRequest) (*VideoResponse, error) {
    switch p.apiType {
    case "gemini":
        return p.generateWithGemini(req)
    case "luma":
        return p.generateWithLuma(req)
    case "runway":
        return p.generateWithRunway(req)
    case "pika":
        return p.generateWithPika(req)
    default:
        return nil, fmt.Errorf("unsupported API type: %s", p.apiType)
    }
}

// generateWithGemini توليد فيديو باستخدام Gemini Veo
func (p *VideoProvider) generateWithGemini(req VideoRequest) (*VideoResponse, error) {
    // هذا مثال مبسط - في الواقع تحتاج إلى استيراد مكتبة Google AI بشكل صحيح
    // هنا نستخدم استدعاء REST API مباشرة
    
    url := fmt.Sprintf("%s/models/veo-2.0-generate-001:generateVideo?key=%s", p.baseURL, p.apiKey)
    
    requestBody := map[string]interface{}{
        "prompt": req.Prompt,
        "video_length_seconds": req.Duration,
        "aspect_ratio": req.AspectRatio,
        "output_format": req.OutputFormat,
    }
    
    jsonData, err := json.Marshal(requestBody)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal request: %w", err)
    }
    
    resp, err := p.httpClient.Post(url, "application/json", strings.NewReader(string(jsonData)))
    if err != nil {
        return nil, fmt.Errorf("Gemini API request failed: %w", err)
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        body, _ := os.ReadAll(resp.Body)
        return nil, fmt.Errorf("Gemini API error: %s - %s", resp.Status, string(body))
    }
    
    var result struct {
        VideoURL    string `json:"videoUrl"`
        OperationID string `json:"operationId"`
        Status      string `json:"status"`
    }
    
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, fmt.Errorf("failed to decode response: %w", err)
    }
    
    return &VideoResponse{
        VideoURL:    result.VideoURL,
        Duration:    float64(req.Duration),
        Status:      result.Status,
        OperationID: result.OperationID,
        Cost:        0.1, // تقدير تكلفة Gemini Veo
    }, nil
}

// generateWithLuma توليد فيديو باستخدام Luma AI
func (p *VideoProvider) generateWithLuma(req VideoRequest) (*VideoResponse, error) {
    url := p.baseURL + "/generations"
    
    requestBody := map[string]interface{}{
        "prompt": req.Prompt,
        "aspect_ratio": req.AspectRatio,
    }
    
    jsonData, err := json.Marshal(requestBody)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal request: %w", err)
    }
    
    httpReq, err := http.NewRequest("POST", url, strings.NewReader(string(jsonData)))
    if err != nil {
        return nil, fmt.Errorf("failed to create request: %w", err)
    }
    
    httpReq.Header.Set("Content-Type", "application/json")
    httpReq.Header.Set("Authorization", "Bearer "+p.apiKey)
    
    resp, err := p.httpClient.Do(httpReq)
    if err != nil {
        return nil, fmt.Errorf("Luma API request failed: %w", err)
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        body, _ := os.ReadAll(resp.Body)
        return nil, fmt.Errorf("Luma API error: %s - %s", resp.Status, string(body))
    }
    
    var result struct {
        ID          string `json:"id"`
        Status      string `json:"status"`
        VideoURL    string `json:"video_url"`
    }
    
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, fmt.Errorf("failed to decode response: %w", err)
    }
    
    return &VideoResponse{
        VideoURL:    result.VideoURL,
        Duration:    float64(req.Duration),
        Status:      result.Status,
        OperationID: result.ID,
        Cost:        0.05, // تقدير تكلفة Luma
    }, nil
}

// generateWithRunway توليد فيديو باستخدام Runway ML
func (p *VideoProvider) generateWithRunway(req VideoRequest) (*VideoResponse, error) {
    url := p.baseURL + "/generations"
    
    requestBody := map[string]interface{}{
        "prompt": req.Prompt,
        "duration": req.Duration,
        "aspect_ratio": req.AspectRatio,
    }
    
    jsonData, err := json.Marshal(requestBody)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal request: %w", err)
    }
    
    httpReq, err := http.NewRequest("POST", url, strings.NewReader(string(jsonData)))
    if err != nil {
        return nil, fmt.Errorf("failed to create request: %w", err)
    }
    
    httpReq.Header.Set("Content-Type", "application/json")
    httpReq.Header.Set("Authorization", "Bearer "+p.apiKey)
    
    resp, err := p.httpClient.Do(httpReq)
    if err != nil {
        return nil, fmt.Errorf("Runway API request failed: %w", err)
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        body, _ := os.ReadAll(resp.Body)
        return nil, fmt.Errorf("Runway API error: %s - %s", resp.Status, string(body))
    }
    
    var result struct {
        Generation struct {
            ID     string `json:"id"`
            Status string `json:"status"`
            URL    string `json:"url"`
        } `json:"generation"`
    }
    
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, fmt.Errorf("failed to decode response: %w", err)
    }
    
    return &VideoResponse{
        VideoURL:    result.Generation.URL,
        Duration:    float64(req.Duration),
        Status:      result.Generation.Status,
        OperationID: result.Generation.ID,
        Cost:        0.08, // تقدير تكلفة Runway
    }, nil
}

// generateWithPika توليد فيديو باستخدام Pika Labs
func (p *VideoProvider) generateWithPika(req VideoRequest) (*VideoResponse, error) {
    url := p.baseURL + "/generate"
    
    requestBody := map[string]interface{}{
        "prompt": req.Prompt,
        "duration": req.Duration,
        "aspect_ratio": req.AspectRatio,
    }
    
    jsonData, err := json.Marshal(requestBody)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal request: %w", err)
    }
    
    httpReq, err := http.NewRequest("POST", url, strings.NewReader(string(jsonData)))
    if err != nil {
        return nil, fmt.Errorf("failed to create request: %w", err)
    }
    
    httpReq.Header.Set("Content-Type", "application/json")
    httpReq.Header.Set("Authorization", "Bearer "+p.apiKey)
    
    resp, err := p.httpClient.Do(httpReq)
    if err != nil {
        return nil, fmt.Errorf("Pika API request failed: %w", err)
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        body, _ := os.ReadAll(resp.Body)
        return nil, fmt.Errorf("Pika API error: %s - %s", resp.Status, string(body))
    }
    
    var result struct {
        ID     string `json:"id"`
        Status string `json:"status"`
        Video  struct {
            URL string `json:"url"`
        } `json:"video"`
    }
    
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, fmt.Errorf("failed to decode response: %w", err)
    }
    
    return &VideoResponse{
        VideoURL:    result.Video.URL,
        Duration:    float64(req.Duration),
        Status:      result.Status,
        OperationID: result.ID,
        Cost:        0.03, // تقدير تكلفة Pika
    }, nil
}

// AnalyzeImage تحليل صورة - غير مدعوم في معظم مزودي الفيديو
func (p *VideoProvider) AnalyzeImage(req AnalysisRequest) (*AnalysisResponse, error) {
    return nil, fmt.Errorf("image analysis not supported by video provider %s", p.apiType)
}

// AnalyzeText تحليل نص - غير مدعوم في معظم مزودي الفيديو
func (p *VideoProvider) AnalyzeText(req AnalysisRequest) (*AnalysisResponse, error) {
    return nil, fmt.Errorf("text analysis not supported by video provider %s", p.apiType)
}

// TranslateText ترجمة نص - غير مدعوم في معظم مزودي الفيديو
func (p *VideoProvider) TranslateText(req TranslationRequest) (*TranslationResponse, error) {
    return nil, fmt.Errorf("text translation not supported by video provider %s", p.apiType)
}

// GetVideoStatus الحصول على حالة فيديو
func (p *VideoProvider) GetVideoStatus(operationID string) (*VideoResponse, error) {
    switch p.apiType {
    case "gemini":
        return p.getGeminiStatus(operationID)
    case "luma":
        return p.getLumaStatus(operationID)
    case "runway":
        return p.getRunwayStatus(operationID)
    case "pika":
        return p.getPikaStatus(operationID)
    default:
        return nil, fmt.Errorf("unsupported API type: %s", p.apiType)
    }
}

// getGeminiStatus الحصول على حالة فيديو Gemini
func (p *VideoProvider) getGeminiStatus(operationID string) (*VideoResponse, error) {
    url := fmt.Sprintf("%s/operations/%s?key=%s", p.baseURL, operationID, p.apiKey)
    
    resp, err := p.httpClient.Get(url)
    if err != nil {
        return nil, fmt.Errorf("failed to get Gemini status: %w", err)
    }
    defer resp.Body.Close()
    
    var result struct {
        Done   bool `json:"done"`
        Result struct {
            VideoURL string `json:"videoUrl"`
        } `json:"result"`
        Error *struct {
            Message string `json:"message"`
        } `json:"error"`
    }
    
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, fmt.Errorf("failed to decode Gemini status: %w", err)
    }
    
    response := &VideoResponse{
        OperationID: operationID,
        Status:      "pending",
    }
    
    if result.Done {
        if result.Error != nil {
            response.Status = "failed"
        } else {
            response.Status = "completed"
            response.VideoURL = result.Result.VideoURL
        }
    }
    
    return response, nil
}

// getLumaStatus الحصول على حالة فيديو Luma
func (p *VideoProvider) getLumaStatus(operationID string) (*VideoResponse, error) {
    url := fmt.Sprintf("%s/generations/%s", p.baseURL, operationID)
    
    httpReq, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to create request: %w", err)
    }
    
    httpReq.Header.Set("Authorization", "Bearer "+p.apiKey)
    
    resp, err := p.httpClient.Do(httpReq)
    if err != nil {
        return nil, fmt.Errorf("failed to get Luma status: %w", err)
    }
    defer resp.Body.Close()
    
    var result struct {
        Status   string `json:"status"`
        VideoURL string `json:"video_url"`
    }
    
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, fmt.Errorf("failed to decode Luma status: %w", err)
    }
    
    return &VideoResponse{
        OperationID: operationID,
        Status:      result.Status,
        VideoURL:    result.VideoURL,
    }, nil
}

// GenerateText توليد نص - غير مدعوم في مزود الفيديو
func (p *VideoProvider) GenerateText(req TextRequest) (*TextResponse, error) {
    return nil, fmt.Errorf("text generation not supported by video provider %s", p.apiType)
}

// GenerateImage توليد صورة - غير مدعوم في مزود الفيديو
func (p *VideoProvider) GenerateImage(req ImageRequest) (*ImageResponse, error) {
    return nil, fmt.Errorf("image generation not supported by video provider %s", p.apiType)
}

// IsAvailable التحقق من التوفر
func (p *VideoProvider) IsAvailable() bool {
    // محاولة الاتصال بالـ API للتحقق من التوفر
    testURL := p.baseURL + "/"
    if p.apiType == "gemini" {
        testURL = fmt.Sprintf("%s/models?key=%s", p.baseURL, p.apiKey)
    }
    
    resp, err := p.httpClient.Get(testURL)
    if err != nil {
        return false
    }
    defer resp.Body.Close()
    
    return resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusUnauthorized
}

// GetName اسم المزود
func (p *VideoProvider) GetName() string {
    switch p.apiType {
    case "gemini":
        return "Gemini Video (Veo)"
    case "luma":
        return "Luma AI"
    case "runway":
        return "Runway ML"
    case "pika":
        return "Pika Labs"
    default:
        return "Video Provider"
    }
}

// GetCost التكلفة
func (p *VideoProvider) GetCost() float64 {
    switch p.apiType {
    case "gemini":
        return 0.1 // دولار لكل فيديو
    case "luma":
        return 0.05
    case "runway":
        return 0.08
    case "pika":
        return 0.03
    default:
        return 0.1
    }
}

// GetStats الحصول على إحصائيات
func (p *VideoProvider) GetStats() *ProviderStats {
    return &ProviderStats{
        Name:        p.GetName(),
        Type:        "video",
        IsAvailable: p.IsAvailable(),
        Requests:    0,
        Errors:      0,
        LastUsed:    "",
        TotalCost:   0,
        SuccessRate: 85.0, // تقديري
    }
}

// getRunwayStatus الحصول على حالة فيديو Runway
func (p *VideoProvider) getRunwayStatus(operationID string) (*VideoResponse, error) {
    url := fmt.Sprintf("%s/generations/%s", p.baseURL, operationID)
    
    httpReq, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to create request: %w", err)
    }
    
    httpReq.Header.Set("Authorization", "Bearer "+p.apiKey)
    
    resp, err := p.httpClient.Do(httpReq)
    if err != nil {
        return nil, fmt.Errorf("failed to get Runway status: %w", err)
    }
    defer resp.Body.Close()
    
    var result struct {
        Generation struct {
            Status string `json:"status"`
            URL    string `json:"url"`
        } `json:"generation"`
    }
    
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, fmt.Errorf("failed to decode Runway status: %w", err)
    }
    
    return &VideoResponse{
        OperationID: operationID,
        Status:      result.Generation.Status,
        VideoURL:    result.Generation.URL,
    }, nil
}

// getPikaStatus الحصول على حالة فيديو Pika
func (p *VideoProvider) getPikaStatus(operationID string) (*VideoResponse, error) {
    url := fmt.Sprintf("%s/generations/%s", p.baseURL, operationID)
    
    httpReq, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to create request: %w", err)
    }
    
    httpReq.Header.Set("Authorization", "Bearer "+p.apiKey)
    
    resp, err := p.httpClient.Do(httpReq)
    if err != nil {
        return nil, fmt.Errorf("failed to get Pika status: %w", err)
    }
    defer resp.Body.Close()
    
    var result struct {
        Status string `json:"status"`
        Video  struct {
            URL string `json:"url"`
        } `json:"video"`
    }
    
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, fmt.Errorf("failed to decode Pika status: %w", err)
    }
    
    return &VideoResponse{
        OperationID: operationID,
        Status:      result.Status,
        VideoURL:    result.Video.URL,
    }, nil
}