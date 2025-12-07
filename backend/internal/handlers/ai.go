package handlers

import (
    "io"
    "net/http"
    "strings"
    
    "github.com/gin-gonic/gin"
    "github.com/nawthtech/nawthtech/backend/internal/ai"
    "github.com/nawthtech/nawthtech/backend/internal/utils"
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
        Prompt      string  `json:"prompt" binding:"required"`
        Language    string  `json:"language" default:"en"`
        ContentType string  `json:"content_type"`
        Tone        string  `json:"tone"`
        Length      string  `json:"length"`
        MaxTokens   int     `json:"max_tokens" default:"1000"`
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
    
    // توليد المحتوى باستخدام الواجهة الصحيحة
    content, err := h.aiClient.GenerateText(prompt, h.getUserID(c))
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
            "content":       content,
            "language":      req.Language,
            "content_type":  req.ContentType,
            "tone":          req.Tone,
            "provider":      "ai_provider",
            "model_used":    "default_model",
            "cost":          0.0,
            "tokens_used":   len(content),
            "created_at":    "now",
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
    
    // تحليل الصورة باستخدام الواجهة الصحيحة
    analysis, err := h.aiClient.AnalyzeImage(imageData, prompt, h.getUserID(c))
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
            "analysis":      analysis,
            "confidence":    0.85,
            "filename":      file.Filename,
            "size":          file.Size,
            "content_type":  contentType,
            "analysis_type": "general",
            "provider":      "ai_provider",
            "model_used":    "vision_model",
            "cost":          0.0,
            "created_at":    "now",
        },
    })
}

// TranslateTextHandler معالج ترجمة النص
func (h *AIHandler) TranslateTextHandler(c *gin.Context) {
    var req struct {
        Text       string `json:"text" binding:"required"`
        SourceLang string `json:"source_lang" default:"auto"`
        TargetLang string `json:"target_lang" binding:"required"`
        Formal     bool   `json:"formal" default:"false"`
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
    if !h.isValidLanguage(req.TargetLang) {
        c.JSON(http.StatusBadRequest, gin.H{
            "success": false,
            "error":   "Invalid target language",
        })
        return
    }
    
    // ترجمة النص باستخدام الواجهة الصحيحة
    translatedText, err := h.aiClient.TranslateText(req.Text, req.SourceLang, req.TargetLang, h.getUserID(c))
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
            "original_text":   req.Text,
            "translated_text": translatedText,
            "source_lang":     req.SourceLang,
            "target_lang":     req.TargetLang,
            "detected_lang":   "auto",
            "confidence":      0.9,
            "provider":        "translation_provider",
            "model_used":      "translation_model",
            "cost":            0.0,
            "created_at":      "now",
        },
    })
}

// SummarizeTextHandler معالج تلخيص النص
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
    
    // توليد التلخيص باستخدام الواجهة الصحيحة
    summary, err := h.aiClient.GenerateText(prompt, h.getUserID(c))
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
            "summary":         summary,
            "summary_type":    req.SummaryType,
            "language":        req.Language,
            "provider":        "ai_provider",
            "model_used":      "summarization_model",
            "cost":            0.0,
            "created_at":      "now",
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
                "supported_languages": h.getSupportedLanguagesSimple(),
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

func (h *AIHandler) isValidLanguage(lang string) bool {
    supportedLangs := h.getSupportedLanguagesSimple()
    for _, supportedLang := range supportedLangs {
        if supportedLang == lang {
            return true
        }
    }
    return false
}

func (h *AIHandler) getSupportedLanguagesSimple() []string {
    return []string{
        "en", "ar", "fr", "es", "de", "zh",
        "ru", "ja", "ko", "pt", "it", "nl",
        "tr", "fa", "hi", "bn", "ur",
    }
}

// Helper functions
func (h *AIHandler) getUserID(c *gin.Context) string {
    return utils.GetUserIDFromContext(c)
}