package handlers

import (
    "net/http"
    "path/filepath"
    
    "github.com/gin-gonic/gin"
    "github.com/nawthtech/nawthtech/backend/internal/ai"
    "github.com/nawthtech/nawthtech/backend/internal/services"
)

type VideoHandler struct {
    videoService *services.VideoService
    aiClient     *ai.Client
}

func NewVideoHandler(videoService *services.VideoService, aiClient *ai.Client) *VideoHandler {
    return &VideoHandler{
        videoService: videoService,
        aiClient:     aiClient,
    }
}

// GenerateVideoHandler معلب توليد فيديو
func (h *VideoHandler) GenerateVideoHandler(c *gin.Context) {
    var req struct {
        Prompt      string `json:"prompt" binding:"required"`
        Duration    int    `json:"duration" default:"30"`
        AspectRatio string `json:"aspect_ratio" default:"16:9"`
        Style       string `json:"style" default:"animated"`
        VideoType   string `json:"video_type"` // explainer, promotional, etc.
    }
    
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    videoReq := ai.VideoRequest{
        Prompt:      req.Prompt,
        Duration:    req.Duration,
        AspectRatio: req.AspectRatio,
        Style:       req.Style,
    }
    
    // استخدام النوع المخصص إذا محدد
    var videoResp *ai.VideoResponse
    var err error
    
    if req.VideoType != "" {
        videoResp, err = h.aiClient.GenerateNawthTechVideo(req.VideoType, req.Prompt)
    } else {
        videoResp, err = h.aiClient.GenerateVideo(videoReq)
    }
    
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error":   "Failed to generate video",
            "details": err.Error(),
        })
        return
    }
    
    c.JSON(http.StatusAccepted, gin.H{
        "success": true,
        "data": gin.H{
            "job_id":       videoResp.GenerationID,
            "status":       videoResp.Status,
            "video_url":    videoResp.VideoURL,
            "duration":     videoResp.Duration,
            "size":         videoResp.Size,
            "model_used":   videoResp.ModelUsed,
        },
        "message": "Video generation started. Check status with the job ID.",
    })
}

// GetVideoStatusHandler معلب حالة الفيديو
func (h *VideoHandler) GetVideoStatusHandler(c *gin.Context) {
    jobID := c.Param("jobId")
    operationID := c.Query("operation_id")
    
    var videoResp *ai.VideoResponse
    var err error
    
    if operationID != "" {
        videoResp, err = h.aiClient.GetVideoStatus(operationID)
    } else if jobID != "" {
        // الحصول من نظام jobs المحلي
        job, err := h.videoService.GetJob(jobID)
        if err != nil {
            c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
            return
        }
        
        c.JSON(http.StatusOK, gin.H{
            "job_id":    job.ID,
            "status":    job.Status,
            "progress":  job.Progress,
            "created":   job.CreatedAt,
            "updated":   job.UpdatedAt,
            "result":    job.Result,
        })
        return
    } else {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Job ID or Operation ID required"})
        return
    }
    
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, videoResp)
}

// DownloadVideoHandler معلب تحميل الفيديو
func (h *VideoHandler) DownloadVideoHandler(c *gin.Context) {
    jobID := c.Param("jobId")
    
    job, err := h.videoService.GetJob(jobID)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
        return
    }
    
    if job.Status != "completed" || job.Result == nil || len(job.Result.VideoData) == 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Video not ready for download"})
        return
    }
    
    filename := filepath.Base(job.Result.VideoURL)
    if filename == "" || filename == "." {
        filename = "nawthtech_video_" + jobID + ".mp4"
    }
    
    c.Header("Content-Type", "video/mp4")
    c.Header("Content-Disposition", "attachment; filename="+filename)
    c.Header("Content-Length", string(len(job.Result.VideoData)))
    
    c.Data(http.StatusOK, "video/mp4", job.Result.VideoData)
}

// ListVideoTypesHandler معلب أنواع الفيديوهات المتاحة
func (h *VideoHandler) ListVideoTypesHandler(c *gin.Context) {
    videoTypes := []gin.H{
        {"id": "explainer", "name": "Explainer Video", "duration": 60, "style": "animated"},
        {"id": "promotional", "name": "Promotional Video", "duration": 30, "style": "cinematic"},
        {"id": "tutorial", "name": "Tutorial Video", "duration": 120, "style": "realistic"},
        {"id": "testimonial", "name": "Testimonial Video", "duration": 45, "style": "corporate"},
        {"id": "social", "name": "Social Media Video", "duration": 15, "style": "animated"},
    }
    
    c.JSON(http.StatusOK, gin.H{
        "video_types": videoTypes,
        "max_duration": 300, // 5 دقائق
        "supported_formats": []string{"mp4", "gif"},
        "aspect_ratios": []string{"16:9", "1:1", "9:16", "4:5"},
    })
}