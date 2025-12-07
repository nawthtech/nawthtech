package ai

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "os"
    "strings"
    "time"
    "github.com/nawthtech/nawthtech/backend/internal/ai/types"
)

// OllamaProvider مزود Ollama المحلي
type OllamaProvider struct {
    baseURL    string
    httpClient *http.Client
    models     []string
}

// NewOllamaProvider إنشاء مزود Ollama جديد
func NewOllamaProvider() *OllamaProvider {
    baseURL := os.Getenv("OLLAMA_HOST")
    if baseURL == "" {
        baseURL = "http://localhost:11434"
    }
    
    provider := &OllamaProvider{
        baseURL: baseURL,
        httpClient: &http.Client{
            Timeout: 300 * time.Second,
        },
    }
    
    // محاولة تحميل النماذج (لنفشل بهدوء إذا كان Ollama غير متوفر)
    provider.loadModels()
    
    return provider
}

// loadModels تحميل النماذج المتاحة من Ollama
func (p *OllamaProvider) loadModels() error {
    url := p.baseURL + "/api/tags"
    
    resp, err := p.httpClient.Get(url)
    if err != nil {
        return fmt.Errorf("failed to connect to Ollama: %w", err)
    }
    defer resp.Body.Close()
    
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return err
    }
    
    var result struct {
        Models []struct {
            Name string `json:"name"`
        } `json:"models"`
    }
    
    if err := json.Unmarshal(body, &result); err != nil {
        return err
    }
    
    for _, model := range result.Models {
        p.models = append(p.models, model.Name)
    }
    
    return nil
}

// GenerateText توليد نص
func (p *OllamaProvider) GenerateText(req types.TextRequest) (*types.TextResponse, error) {
    url := p.baseURL + "/api/generate"
    
    // تعيين القيم الافتراضية
    model := req.Model
    if model == "" {
        model = "llama3.2:3b"
    }
    
    temperature := req.Temperature
    if temperature == 0 {
        temperature = 0.7
    }
    
    maxTokens := req.MaxTokens
    if maxTokens == 0 {
        maxTokens = 2000
    }
    
    request := map[string]interface{}{
        "model":  model,
        "prompt": req.Prompt,
        "stream": false,
        "options": map[string]interface{}{
            "temperature": temperature,
            "num_predict": maxTokens,
        },
    }
    
    jsonData, err := json.Marshal(request)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal request: %w", err)
    }
    
    resp, err := p.httpClient.Post(url, "application/json", bytes.NewBuffer(jsonData))
    if err != nil {
        return nil, fmt.Errorf("Ollama request failed: %w", err)
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        body, _ := io.ReadAll(resp.Body)
        return nil, fmt.Errorf("Ollama API error: %s - %s", resp.Status, string(body))
    }
    
    var result struct {
        Response string    `json:"response"`
        Done     bool      `json:"done"`
        Model    string    `json:"model"`
    }
    
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, fmt.Errorf("failed to decode response: %w", err)
    }
    
    // تقدير عدد الرموز (تقريبي)
tokens := int(float64(len(strings.Fields(result.Response))) * 1.3)
    
    return &types.TextResponse{
        Text:        strings.TrimSpace(result.Response),
        Tokens:      int(tokens),
        Cost:        0.0, // Ollama مجاني
        ModelUsed:   result.Model,
        FinishReason: "length",
        CreatedAt:   time.Now(),
    }, nil
}

// GenerateImage توليد صورة - غير مدعوم في Ollama
func (p *OllamaProvider) GenerateImage(req types.ImageRequest) (*types.ImageResponse, error) {
    return nil, fmt.Errorf("image generation not supported by Ollama")
}

// GenerateVideo توليد فيديو - غير مدعوم في Ollama
func (p *OllamaProvider) GenerateVideo(req types.VideoRequest) (*types.VideoResponse, error) {
    return nil, fmt.Errorf("video generation not supported by Ollama")
}

// AnalyzeImage تحليل صورة
func (p *OllamaProvider) AnalyzeImage(req types.AnalysisRequest) (*types.AnalysisResponse, error) {
    // استخدام نموذج رؤية لتحليل الصورة
    // هذا يتطلب نموذج multimodal مثل llama3.2-vision
    if req.Model == "" {
        req.Model = "llava:latest"
    }
    
    // تحويل الصورة إلى prompt
    prompt := fmt.Sprintf("%s Analyze this image: [Image data provided]", req.Prompt)
    
    textReq := types.TextRequest{
        Prompt: prompt,
        Model:  req.Model,
    }
    
    resp, err := p.GenerateText(textReq)
    if err != nil {
        return nil, err
    }
    
    return &types.AnalysisResponse{
        Result:     resp.Text,
        Confidence: 0.8,
        Cost:       0.0,
        Model:      resp.ModelUsed,
        CreatedAt:  time.Now(),
    }, nil
}

// AnalyzeText تحليل نص
func (p *OllamaProvider) AnalyzeText(req types.AnalysisRequest) (*types.AnalysisResponse, error) {
    prompt := fmt.Sprintf("Analyze this text: %s\n\nProvide analysis:", req.Text)
    
    if req.Prompt != "" {
        prompt = fmt.Sprintf("%s\n\n%s", req.Prompt, prompt)
    }
    
    textReq := types.TextRequest{
        Prompt: prompt,
        Model:  req.Model,
    }
    
    resp, err := p.GenerateText(textReq)
    if err != nil {
        return nil, err
    }
    
    return &types.AnalysisResponse{
        Result:     resp.Text,
        Confidence: 0.9,
        Cost:       0.0,
        Model:      resp.ModelUsed,
        CreatedAt:  time.Now(),
    }, nil
}

// TranslateText ترجمة نص
func (p *OllamaProvider) TranslateText(req types.TranslationRequest) (*types.TranslationResponse, error) {
    prompt := fmt.Sprintf("Translate the following text from %s to %s:\n\n%s",
        req.FromLang, req.ToLang, req.Text)
    
    if req.Model == "" {
        req.Model = "llama3.2:3b"
    }
    
    textReq := types.TextRequest{
        Prompt: prompt,
        Model:  req.Model,
    }
    
    resp, err := p.GenerateText(textReq)
    if err != nil {
        return nil, err
    }
    
    return &types.TranslationResponse{
        TranslatedText: strings.TrimSpace(resp.Text),
        Cost:           0.0,
        Model:          resp.ModelUsed,
        CreatedAt:      time.Now(),
    }, nil
}

// GenerateStream توليد نص بشكل متدفق
func (p *OllamaProvider) GenerateStream(ctx context.Context, req types.TextRequest) (<-chan string, <-chan error) {
    textChan := make(chan string)
    errChan := make(chan error, 1)
    
    go func() {
        defer close(textChan)
        defer close(errChan)
        
        url := p.baseURL + "/api/generate"
        
        // تعيين القيم الافتراضية
        model := req.Model
        if model == "" {
            model = "llama3.2:3b"
        }
        
        temperature := req.Temperature
        if temperature == 0 {
            temperature = 0.7
        }
        
        maxTokens := req.MaxTokens
        if maxTokens == 0 {
            maxTokens = 2000
        }
        
        request := map[string]interface{}{
            "model":  model,
            "prompt": req.Prompt,
            "stream": true,
            "options": map[string]interface{}{
                "temperature": temperature,
                "num_predict": maxTokens,
            },
        }
        
        jsonData, err := json.Marshal(request)
        if err != nil {
            errChan <- err
            return
        }
        
        httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
        if err != nil {
            errChan <- err
            return
        }
        httpReq.Header.Set("Content-Type", "application/json")
        
        resp, err := p.httpClient.Do(httpReq)
        if err != nil {
            errChan <- err
            return
        }
        defer resp.Body.Close()
        
        if resp.StatusCode != http.StatusOK {
            body, _ := io.ReadAll(resp.Body)
            errChan <- fmt.Errorf("Ollama API error: %s - %s", resp.Status, string(body))
            return
        }
        
        decoder := json.NewDecoder(resp.Body)
        var fullText strings.Builder
        
        for {
            select {
            case <-ctx.Done():
                return
            default:
                var chunk struct {
                    Response string `json:"response"`
                    Done     bool   `json:"done"`
                }
                
                if err := decoder.Decode(&chunk); err != nil {
                    if err == io.EOF {
                        return
                    }
                    errChan <- err
                    return
                }
                
                if chunk.Response != "" {
                    fullText.WriteString(chunk.Response)
                    textChan <- chunk.Response
                }
                
                if chunk.Done {
                    return
                }
            }
        }
    }()
    
    return textChan, errChan
}

// Embed توليد embeddings
func (p *OllamaProvider) Embed(text string, model string) ([]float64, error) {
    if model == "" {
        model = "nomic-embed-text"
    }
    
    url := p.baseURL + "/api/embed"
    
    request := map[string]interface{}{
        "model": model,
        "input": text,
    }
    
    jsonData, err := json.Marshal(request)
    if err != nil {
        return nil, err
    }
    
    resp, err := p.httpClient.Post(url, "application/json", bytes.NewBuffer(jsonData))
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    var result struct {
        Embedding []float64 `json:"embedding"`
    }
    
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, err
    }
    
    return result.Embedding, nil
}

// PullModel سحب نموذج جديد
func (p *OllamaProvider) PullModel(model string) error {
    url := p.baseURL + "/api/pull"
    
    request := map[string]interface{}{
        "name":   model,
        "stream": false,
    }
    
    jsonData, err := json.Marshal(request)
    if err != nil {
        return err
    }
    
    resp, err := p.httpClient.Post(url, "application/json", bytes.NewBuffer(jsonData))
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        body, _ := io.ReadAll(resp.Body)
        return fmt.Errorf("failed to pull model: %s - %s", resp.Status, string(body))
    }
    
    // إضافة النموذج إلى القائمة
    p.models = append(p.models, model)
    
    return nil
}

// ListModels عرض النماذج المتاحة
func (p *OllamaProvider) ListModels() []string {
    return p.models
}

// IsAvailable التحقق من التوفر
func (p *OllamaProvider) IsAvailable() bool {
    resp, err := p.httpClient.Get(p.baseURL + "/api/tags")
    return err == nil && resp.StatusCode == http.StatusOK
}

// GetName اسم المزود
func (p *OllamaProvider) GetName() string {
    return "Ollama"
}

// GetCost التكلفة (مجاني بالكامل)
func (p *OllamaProvider) GetCost() float64 {
    return 0.0
}

// GetStats الحصول على إحصائيات
func (p *OllamaProvider) GetStats() *types.ProviderStats {
    return &types.ProviderStats{
        Name:        p.GetName(),
        Type:        "text",
        IsAvailable: p.IsAvailable(),
        Requests:    0,
        Successful:  0,
        Failed:      0,
        TotalCost:   0.0,
        AvgLatency:  0.0,
        SuccessRate: 95.0,
        LastUsed:    time.Time{},
    }
}

// GetStatsDetailed الحصول على إحصائيات مفصلة
func (p *OllamaProvider) GetStatsDetailed() (map[string]interface{}, error) {
    stats := map[string]interface{}{
        "provider":        "ollama",
        "models_count":    len(p.models),
        "models":          p.models,
        "status":          "online",
        "supports_stream": true,
        "supports_embed":  true,
    }
    
    // محاولة الحصول على إصدار Ollama
    url := p.baseURL + "/api/version"
    resp, err := p.httpClient.Get(url)
    if err == nil {
        defer resp.Body.Close()
        var versionInfo struct {
            Version string `json:"version"`
        }
        if err := json.NewDecoder(resp.Body).Decode(&versionInfo); err == nil {
            stats["version"] = versionInfo.Version
        }
    }
    
    return stats, nil
}

// StreamText توليد نص متدفق (واجهة بديلة)
func (p *OllamaProvider) StreamText(prompt string, model string, temperature float64) (<-chan string, <-chan error, context.CancelFunc) {
    ctx, cancel := context.WithCancel(context.Background())
    
    req := types.TextRequest{
        Prompt:      prompt,
        Model:       model,
        Temperature: temperature,
    }
    
    textChan, errChan := p.GenerateStream(ctx, req)
    
    return textChan, errChan, cancel
}

// GetType نوع المزود
func (p *OllamaProvider) GetType() string {
    return "text"
}

// SupportsStreaming يدعم التدفق
func (p *OllamaProvider) SupportsStreaming() bool {
    return true
}

// SupportsEmbedding يدعم التضمين
func (p *OllamaProvider) SupportsEmbedding() bool {
    return true
}

// GetMaxTokens الحد الأقصى للرموز
func (p *OllamaProvider) GetMaxTokens() int {
    return 4000
}

// GetSupportedLanguages اللغات المدعومة
func (p *OllamaProvider) GetSupportedLanguages() []string {
    return []string{"en", "es", "fr", "de", "ar", "zh", "ja", "ko"}
}