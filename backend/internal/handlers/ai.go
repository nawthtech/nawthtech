package handlers

import (
    "io"
    "net/http"
    "strings"
    
    "github.com/gin-gonic/gin"
    "github.com/nawthtech/nawthtech/backend/internal/ai"
)

type AIHandler struct {
    aiClient *ai.Client
}

func NewAIHandler(aiClient *ai.Client) *AIHandler {
    return &AIHandler{aiClient: aiClient}
}

// GenerateContentHandler معالج توليد المحتوى
func (h *AIHandler) GenerateContentHandler(c *gin.Context) {
    var req struct {
        Prompt      string `json:"prompt" binding:"required"`
        Language    string `json:"language" default:"en"`
        ContentType string `json:"content_type"`
        Tone        string `json:"tone"`
        Length      string `json:"length"`
        MaxTokens   int    `json:"max_tokens" default:"1000"`
        Temperature float64 `json:"temperature" default:"0.7"`
    }
    
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "success": false,
            "error":   "Invalid request format",
            "details": err.Error(),
        })
        return
    }
    
    // التحقق من طول النص
    if len(req.Prompt) > 5000 {
        c.JSON(http.StatusBadRequest, gin.H{
            "success": false,
            "error":   "Prompt is too long (max 5000 characters)",
        })
        return
    }
    
    // بناء prompt محسن
    prompt := h.buildEnhancedPrompt(req.Prompt, req.ContentType, req.Tone, req.Length)
    
    // إنشاء طلب النص
    textReq := ai.TextRequest{
        Prompt:      prompt,
        MaxTokens:   req.MaxTokens,
        Temperature: req.Temperature,
        Language:    req.Language,
        UserID:      getUserID(c),
        UserTier:    getUserTier(c),
    }
    
    // توليد المحتوى
    content, err := h.aiClient.GenerateText(textReq)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "success": false,
            "error":   "Failed to generate content",
            "details": err.Error(),
        })
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "data": gin.H{
            "content":       content.Text,
            "language":      req.Language,
            "content_type":  req.ContentType,
            "tone":          req.Tone,
            "provider":      content.Provider,
            "model_used":    content.ModelUsed,
            "cost":          content.Cost,
            "tokens_used":   content.TokensUsed,
            "created_at":    content.CreatedAt,
        },
    })
}

// AnalyzeImageHandler معالج تحليل الصور
func (h *AIHandler) AnalyzeImageHandler(c *gin.Context) {
    file, err := c.FormFile("image")
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "success": false,
            "error":   "Image file is required",
        })
        return
    }
    
    // التحقق من نوع الملف
    allowedTypes := map[string]bool{
        "image/jpeg": true,
        "image/jpg":  true,
        "image/png":  true,
        "image/gif":  true,
        "image/webp": true,
    }
    
    fileHeader, _ := file.Open()
    buffer := make([]byte, 512)
    _, err = fileHeader.Read(buffer)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "success": false,
            "error":   "Failed to read image file",
        })
        return
    }
    fileHeader.Close()
    
    contentType := http.DetectContentType(buffer)
    if !allowedTypes[contentType] {
        c.JSON(http.StatusBadRequest, gin.H{
            "success": false,
            "error":   "Invalid image format. Supported formats: JPEG, PNG, GIF, WebP",
        })
        return
    }
    
    // التحقق من حجم الملف (10MB كحد أقصى)
    if file.Size > 10*1024*1024 {
        c.JSON(http.StatusBadRequest, gin.H{
            "success": false,
            "error":   "Image size too large (max 10MB)",
        })
        return
    }
    
    prompt := c.PostForm("prompt")
    if prompt == "" {
        prompt = "Describe this image in detail"
    }
    
    analysisType := c.PostForm("analysis_type")
    if analysisType == "" {
        analysisType = "general"
    }
    
    // قراءة ملف الصورة
    uploadedFile, err := file.Open()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "success": false,
            "error":   "Failed to open image file",
            "details": err.Error(),
        })
        return
    }
    defer uploadedFile.Close()
    
    imageData, err := io.ReadAll(uploadedFile)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "success": false,
            "error":   "Failed to read image data",
            "details": err.Error(),
        })
        return
    }
    
    // إنشاء طلب تحليل الصورة
    imageReq := ai.ImageRequest{
        ImageData: imageData,
        Prompt:    prompt,
        AnalysisType: analysisType,
        UserID:    getUserID(c),
        UserTier:  getUserTier(c),
    }
    
    // تحليل الصورة
    analysis, err := h.aiClient.AnalyzeImage(imageReq)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "success": false,
            "error":   "Failed to analyze image",
            "details": err.Error(),
        })
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "data": gin.H{
            "analysis":      analysis.Result,
            "confidence":    analysis.Confidence,
            "filename":      file.Filename,
            "size":          file.Size,
            "content_type":  contentType,
            "analysis_type": analysisType,
            "provider":      analysis.Provider,
            "model_used":    analysis.ModelUsed,
            "cost":          analysis.Cost,
            "created_at":    analysis.CreatedAt,
        },
    })
}

// TranslateTextHandler معالج ترجمة النص
func (h *AIHandler) TranslateTextHandler(c *gin.Context) {
    var req struct {
        Text        string `json:"text" binding:"required"`
        SourceLang  string `json:"source_lang" default:"auto"`
        TargetLang  string `json:"target_lang" binding:"required"`
        Formal      bool   `json:"formal" default:"false"`
    }
    
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "success": false,
            "error":   "Invalid request format",
            "details": err.Error(),
        })
        return
    }
    
    // التحقق من طول النص
    if len(req.Text) > 5000 {
        c.JSON(http.StatusBadRequest, gin.H{
            "success": false,
            "error":   "Text is too long (max 5000 characters)",
        })
        return
    }
    
    // التحقق من صحة اللغة الهدف
    if !isValidLanguage(req.TargetLang) {
        c.JSON(http.StatusBadRequest, gin.H{
            "success": false,
            "error":   "Invalid target language",
        })
        return
    }
    
    // إنشاء طلب الترجمة
    translateReq := ai.TranslationRequest{
        Text:        req.Text,
        SourceLang:  req.SourceLang,
        TargetLang:  req.TargetLang,
        Formal:      req.Formal,
        UserID:      getUserID(c),
        UserTier:    getUserTier(c),
    }
    
    // ترجمة النص
    translation, err := h.aiClient.TranslateText(translateReq)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "success": false,
            "error":   "Failed to translate text",
            "details": err.Error(),
        })
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "data": gin.H{
            "original_text": req.Text,
            "translated_text": translation.TranslatedText,
            "source_lang":   translation.SourceLang,
            "target_lang":   translation.TargetLang,
            "detected_lang": translation.DetectedLang,
            "confidence":    translation.Confidence,
            "provider":      translation.Provider,
            "model_used":    translation.ModelUsed,
            "cost":          translation.Cost,
            "created_at":    translation.CreatedAt,
        },
    })
}

// SummarizeTextHandler معاخ تلخيص النص
func (h *AIHandler) SummarizeTextHandler(c *gin.Context) {
    var req struct {
        Text        string `json:"text" binding:"required"`
        Language    string `json:"language" default:"en"`
        SummaryType string `json:"summary_type" default:"paragraph"`
        MaxLength   int    `json:"max_length" default:"200"`
    }
    
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "success": false,
            "error":   "Invalid request format",
            "details": err.Error(),
        })
        return
    }
    
    // التحقق من طول النص
    if len(req.Text) > 10000 {
        c.JSON(http.StatusBadRequest, gin.H{
            "success": false,
            "error":   "Text is too long (max 10000 characters)",
        })
        return
    }
    
    // بناء prompt للتلخيص
    prompt := h.buildSummaryPrompt(req.Text, req.SummaryType, req.MaxLength)
    
    textReq := ai.TextRequest{
        Prompt:      prompt,
        MaxTokens:   500,
        Temperature: 0.3,
        Language:    req.Language,
        UserID:      getUserID(c),
        UserTier:    getUserTier(c),
    }
    
    // توليد التلخيص
    summary, err := h.aiClient.GenerateText(textReq)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "success": false,
            "error":   "Failed to generate summary",
            "details": err.Error(),
        })
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "data": gin.H{
            "original_length": len(req.Text),
            "summary":        summary.Text,
            "summary_type":   req.SummaryType,
            "language":       req.Language,
            "provider":       summary.Provider,
            "model_used":     summary.ModelUsed,
            "cost":           summary.Cost,
            "created_at":     summary.CreatedAt,
        },
    })
}

// GetAICapabilitiesHandler معالج قدرات الذكاء الاصطناعي
func (h *AIHandler) GetAICapabilitiesHandler(c *gin.Context) {
    capabilities := gin.H{
        "features": []gin.H{
            {
                "name": "text_generation",
                "description": "Generate text content",
                "supported_languages": []string{"en", "ar", "fr", "es", "de", "zh"},
                "max_tokens": 4000,
            },
            {
                "name": "image_analysis",
                "description": "Analyze and describe images",
                "supported_formats": []string{"jpeg", "jpg", "png", "gif", "webp"},
                "max_size_mb": 10,
            },
            {
                "name": "translation",
                "description": "Translate text between languages",
                "supported_languages": getSupportedLanguages(),
                "max_text_length": 5000,
            },
            {
                "name": "summarization",
                "description": "Summarize long texts",
                "max_input_length": 10000,
                "summary_types": []string{"paragraph", "bullet_points", "keywords"},
            },
        },
        
        "content_types": []string{
            "blog_post", "article", "social_media", "email",
            "product_description", "ad_copy", "story", "poem",
            "code", "technical_documentation", "business_plan",
        },
        
        "tones": []string{
            "professional", "casual", "friendly", "formal",
            "persuasive", "informative", "creative", "humorous",
        },
        
        "limits": gin.H{
            "daily_text_generation": 10000,
            "daily_image_analysis":  50,
            "daily_translations":    100,
            "max_concurrent_requests": 5,
        },
    }
    
    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "data":    capabilities,
    })
}

// ===== دوال مساعدة =====

func (h *AIHandler) buildEnhancedPrompt(basePrompt, contentType, tone, length string) string {
    var builder strings.Builder
    
    builder.WriteString(basePrompt)
    
    // إضافة توجيهات بناءً على نوع المحتوى
    if contentType != "" {
        builder.WriteString("\n\nPlease generate " + contentType + " content.")
    }
    
    // إضافة النبرة المطلوبة
    if tone != "" {
        builder.WriteString("\nTone: " + tone)
    }
    
    // إضافة الطول المطلوب
    if length != "" {
        switch length {
        case "short":
            builder.WriteString("\nPlease be concise (1-2 paragraphs).")
        case "medium":
            builder.WriteString("\nPlease provide detailed content (3-5 paragraphs).")
        case "long":
            builder.WriteString("\nPlease provide comprehensive content (6+ paragraphs).")
        }
    }
    
    return builder.String()
}

func (h *AIHandler) buildSummaryPrompt(text, summaryType string, maxLength int) string {
    var prompt strings.Builder
    
    prompt.WriteString("Please summarize the following text")
    
    switch summaryType {
    case "paragraph":
        prompt.WriteString(" in a single paragraph")
    case "bullet_points":
        prompt.WriteString(" using bullet points")
    case "keywords":
        prompt.WriteString(" by extracting key keywords and phrases")
    }
    
    prompt.WriteString(". Maximum length: ")
    prompt.WriteString(string(rune(maxLength)))
    prompt.WriteString(" characters.\n\nText to summarize:\n")
    prompt.WriteString(text)
    
    return prompt.String()
}

func isValidLanguage(lang string) bool {
    supportedLangs := getSupportedLanguages()
    for _, supportedLang := range supportedLangs {
        if supportedLang["code"] == lang {
            return true
        }
    }
    return false
}

func getSupportedLanguages() []gin.H {
    return []gin.H{
        {"code": "en", "name": "English"},
        {"code": "ar", "name": "Arabic"},
        {"code": "fr", "name": "French"},
        {"code": "es", "name": "Spanish"},
        {"code": "de", "name": "German"},
        {"code": "zh", "name": "Chinese"},
        {"code": "ru", "name": "Russian"},
        {"code": "ja", "name": "Japanese"},
        {"code": "ko", "name": "Korean"},
        {"code": "pt", "name": "Portuguese"},
        {"code": "it", "name": "Italian"},
        {"code": "nl", "name": "Dutch"},
        {"code": "tr", "name": "Turkish"},
        {"code": "fa", "name": "Persian"},
        {"code": "hi", "name": "Hindi"},
        {"code": "bn", "name": "Bengali"},
        {"code": "ur", "name": "Urdu"},
    }
}

// Helper functions from video.go (يجب أن تكون في ملف مشترك)
func getUserID(c *gin.Context) string {
    if userID, exists := c.Get("userID"); exists {
        return userID.(string)
    }
    return ""
}

func getUserTier(c *gin.Context) string {
    if userTier, exists := c.Get("userTier"); exists {
        return userTier.(string)
    }
    return "free"
}