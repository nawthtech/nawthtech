package handlers

import (
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

// GenerateContentHandler معلب توليد المحتوى
func (h *AIHandler) GenerateContentHandler(c *gin.Context) {
    var req struct {
        Prompt    string  `json:"prompt" binding:"required"`
        Language  string  `json:"language" default:"en"`
        ContentType string `json:"content_type"`
        Tone      string  `json:"tone"`
        Length    string  `json:"length"`
    }
    
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    // بناء prompt محسن
    prompt := h.buildEnhancedPrompt(req.Prompt, req.ContentType, req.Tone, req.Length)
    
    // توليد المحتوى
    content, err := h.aiClient.GenerateText(prompt, 
        ai.WithLanguage(req.Language),
        ai.WithModel(h.getModelForContent(req.ContentType)),
    )
    
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error":   "Failed to generate content",
            "details": err.Error(),
        })
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "data": gin.H{
            "content":      content,
            "language":     req.Language,
            "content_type": req.ContentType,
            "provider":     "gemini", // أو أي provider تم استخدامه
        },
    })
}

// AnalyzeImageHandler معلب تحليل الصور
func (h *AIHandler) AnalyzeImageHandler(c *gin.Context) {
    file, err := c.FormFile("image")
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Image file required"})
        return
    }
    
    prompt := c.PostForm("prompt")
    if prompt == "" {
        prompt = "Describe this image in detail"
    }
    
    // قراءة ملف الصورة
    uploadedFile, err := file.Open()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    defer uploadedFile.Close()
    
    imageData, err := io.ReadAll(uploadedFile)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    // تحليل الصورة
    analysis, err := h.aiClient.AnalyzeImage(imageData, prompt)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error":   "Failed to analyze image",
            "details": err.Error(),
        })
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "analysis": analysis,
        "filename": file.Filename,
        "size":     file.Size,
    })
}

// ===== دوال مساعدة =====

func (h *AIHandler) buildEnhancedPrompt(basePrompt, contentType, tone, length string) string {
    var builder strings.Builder
    
    builder.WriteString(basePrompt)
    
    if contentType != "" {
        builder.WriteString("\n\nContent Type: " + contentType)
    }
    
    if tone != "" {
        builder.WriteString("\nTone: " + tone)
    }
    
    if length != "" {
        builder.WriteString("\nLength: " + length)
    }
    
    return builder.String()
}

func (h *AIHandler) getModelForContent(contentType string) string {
    switch contentType {
    case "analysis", "strategy", "technical":
        return "gemini-1.5-pro" // للمحتوى المعقد
    default:
        return "gemini-1.5-flash" // للمحتوى البسيط
    }
}