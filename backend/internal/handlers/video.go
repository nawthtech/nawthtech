package handlers

import (
    "fmt"
    "net/http"
    "time"
    
    "github.com/gin-gonic/gin"
    "github.com/nawthtech/nawthtech/backend/internal/ai/video"
)

type VideoHandler struct {
    videoService *video.VideoService
}

func NewVideoHandler(videoService *video.VideoService) *VideoHandler {
    return &VideoHandler{
        videoService: videoService,
    }
}

// GenerateVideoHandler معالج توليد فيديو
func (h *VideoHandler) GenerateVideoHandler(c *gin.Context) {
    var req struct {
        Prompt         string                 `json:"prompt" binding:"required"`
        Duration       int                    `json:"duration" default:"5"`
        Resolution     string                 `json:"resolution" default:"512x512"`
        Aspect         string                 `json:"aspect" default:"1:1"`
        Style          string                 `json:"style" default:"realistic"`
        NegativePrompt string                 `json:"negative_prompt,omitempty"`
        Options        map[string]interface{} `json:"options,omitempty"`
    }
    
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "success": false,
            "error":   "Invalid request format",
            "details": err.Error(),
        })
        return
    }
    
    // التحقق من صحة المدة
    if req.Duration < 1 || req.Duration > 60 {
        c.JSON(http.StatusBadRequest, gin.H{
            "success": false,
            "error":   "Duration must be between 1 and 60 seconds",
        })
        return
    }
    
    // الحصول على معلومات المستخدم
    userID := h.getUserID(c)
    
    // إنشاء طلب الفيديو
    videoReq := video.VideoRequest{
        Prompt:         req.Prompt,
        Duration:       req.Duration,
        Resolution:     req.Resolution,
        Aspect:         req.Aspect,
        Style:          req.Style,
        NegativePrompt: req.NegativePrompt,
        UserID:         userID,
        UserTier:       h.getUserTier(c),
    }
    
    // تحويل الخيارات الإضافية
    if req.Options != nil {
        videoOpts := video.VideoOptions{}
        
        if seed, ok := req.Options["seed"].(float64); ok {
            videoOpts.Seed = int64(seed)
        }
        if fps, ok := req.Options["fps"].(float64); ok {
            videoOpts.FPS = int(fps)
        }
        if quality, ok := req.Options["quality"].(string); ok {
            videoOpts.Quality = quality
        }
        if cfgScale, ok := req.Options["cfg_scale"].(float64); ok {
            videoOpts.CFGScale = cfgScale
        }
        if steps, ok := req.Options["steps"].(float64); ok {
            videoOpts.Steps = int(steps)
        }
        if model, ok := req.Options["model"].(string); ok {
            videoOpts.Model = model
        }
        
        videoReq.Options = videoOpts
    }
    
    // التحقق من صحة الطلب
    if err := video.ValidateVideoRequest(videoReq); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "success": false,
            "error":   "Invalid video request",
            "details": err.Error(),
        })
        return
    }
    
    // التحقق من إمكانية المستخدم لتوليد فيديو
    if canGenerate, message := h.videoService.CanUserGenerateVideo(userID, h.getUserTier(c)); !canGenerate {
        c.JSON(http.StatusTooManyRequests, gin.H{
            "success": false,
            "error":   message,
        })
        return
    }
    
    // إرسال طلب توليد الفيديو
    job, err := h.videoService.SubmitVideoJob(videoReq)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "success": false,
            "error":   "Failed to submit video job",
            "details": err.Error(),
        })
        return
    }
    
    // تسجيل استخدام المستخدم
    h.videoService.RecordUserGeneration(userID, h.getUserTier(c))
    
    c.JSON(http.StatusAccepted, gin.H{
        "success": true,
        "data": gin.H{
            "job_id":     job.ID,
            "status":     job.Status,
            "progress":   job.Progress,
            "created_at": job.CreatedAt.Format(time.RFC3339),
            "updated_at": job.UpdatedAt.Format(time.RFC3339),
            "prompt":     videoReq.Prompt,
            "duration":   videoReq.Duration,
            "resolution": videoReq.Resolution,
        },
        "message": "Video generation started successfully",
    })
}

// GetVideoStatusHandler معالج حالة الفيديو
func (h *VideoHandler) GetVideoStatusHandler(c *gin.Context) {
    jobID := c.Param("jobId")
    if jobID == "" {
        c.JSON(http.StatusBadRequest, gin.H{
            "success": false,
            "error":   "Job ID is required",
        })
        return
    }
    
    job, err := h.videoService.GetJob(jobID)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{
            "success": false,
            "error":   "Job not found",
            "details": err.Error(),
        })
        return
    }
    
    response := gin.H{
        "job_id":     job.ID,
        "status":     job.Status,
        "progress":   job.Progress,
        "created_at": job.CreatedAt.Format(time.RFC3339),
        "updated_at": job.UpdatedAt.Format(time.RFC3339),
        "prompt":     job.Request.Prompt,
        "duration":   job.Request.Duration,
        "resolution": job.Request.Resolution,
    }
    
    if job.Result != nil {
        response["result"] = gin.H{
            "success":    job.Result.Success,
            "video_url":  job.Result.VideoURL,
            "duration":   job.Result.Duration,
            "resolution": job.Result.Resolution,
            "format":     job.Result.Format,
            "provider":   job.Result.Provider,
            "cost":       job.Result.Cost,
            "error":      job.Result.Error,
            "created_at": job.Result.CreatedAt.Format(time.RFC3339),
        }
    }
    
    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "data":    response,
    })
}

// ListVideoJobsHandler معالج قائمة مهام الفيديو
func (h *VideoHandler) ListVideoJobsHandler(c *gin.Context) {
    limitStr := c.DefaultQuery("limit", "20")
    offsetStr := c.DefaultQuery("offset", "0")
    status := c.Query("status")
    userID := h.getUserID(c)
    
    jobs := h.videoService.ListJobs()
    
    var filteredJobs []*video.VideoJob
    for _, job := range jobs {
        // تصفية حسب المستخدم
        if userID != "" && userID != "admin" && !h.videoService.IsJobOwner(job.ID, userID) {
            continue
        }
        
        // تصفية حسب الحالة
        if status != "" && string(job.Status) != status {
            continue
        }
        
        filteredJobs = append(filteredJobs, job)
    }
    
    // ترتيب حسب التاريخ (الأحدث أولاً)
    for i := 0; i < len(filteredJobs)-1; i++ {
        for j := i + 1; j < len(filteredJobs); j++ {
            if filteredJobs[i].CreatedAt.Before(filteredJobs[j].CreatedAt) {
                filteredJobs[i], filteredJobs[j] = filteredJobs[j], filteredJobs[i]
            }
        }
    }
    
    // التحديد حسب الحدود (تبسيطي)
    var result []gin.H
    for _, job := range filteredJobs {
        result = append(result, gin.H{
            "job_id":     job.ID,
            "status":     job.Status,
            "progress":   job.Progress,
            "created_at": job.CreatedAt.Format(time.RFC3339),
            "prompt":     job.Request.Prompt,
            "duration":   job.Request.Duration,
            "resolution": job.Request.Resolution,
        })
    }
    
    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "data": gin.H{
            "jobs":   result,
            "total":  len(result),
            "limit":  limitStr,
            "offset": offsetStr,
        },
    })
}

// CancelVideoJobHandler معالج إلغاء مهمة فيديو
func (h *VideoHandler) CancelVideoJobHandler(c *gin.Context) {
    jobID := c.Param("jobId")
    if jobID == "" {
        c.JSON(http.StatusBadRequest, gin.H{
            "success": false,
            "error":   "Job ID is required",
        })
        return
    }
    
    // التحقق من ملكية المهمة
    userID := h.getUserID(c)
    if userID != "admin" && !h.videoService.IsJobOwner(jobID, userID) {
        c.JSON(http.StatusForbidden, gin.H{
            "success": false,
            "error":   "You are not authorized to cancel this job",
        })
        return
    }
    
    err := h.videoService.CancelJob(jobID)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "success": false,
            "error":   "Failed to cancel job",
            "details": err.Error(),
        })
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "message": "Video job cancelled successfully",
        "job_id":  jobID,
    })
}

// DownloadVideoHandler معالج تحميل الفيديو
func (h *VideoHandler) DownloadVideoHandler(c *gin.Context) {
    jobID := c.Param("jobId")
    if jobID == "" {
        c.JSON(http.StatusBadRequest, gin.H{
            "success": false,
            "error":   "Job ID is required",
        })
        return
    }
    
    job, err := h.videoService.GetJob(jobID)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{
            "success": false,
            "error":   "Job not found",
        })
        return
    }
    
    // التحقق من ملكية المهمة
    userID := h.getUserID(c)
    if userID != "admin" && !h.videoService.IsJobOwner(jobID, userID) {
        c.JSON(http.StatusForbidden, gin.H{
            "success": false,
            "error":   "You are not authorized to download this video",
        })
        return
    }
    
    if job.Status != video.VideoJobCompleted {
        c.JSON(http.StatusBadRequest, gin.H{
            "success": false,
            "error":   fmt.Sprintf("Video is not ready for download (status: %s)", job.Status),
        })
        return
    }
    
    if job.Result == nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "success": false,
            "error":   "Video result not available",
        })
        return
    }
    
    if len(job.Result.VideoData) == 0 {
        if job.Result.VideoURL != "" {
            // إعادة توجيه إلى URL
            c.Redirect(http.StatusFound, job.Result.VideoURL)
            return
        }
        
        c.JSON(http.StatusBadRequest, gin.H{
            "success": false,
            "error":   "Video data not available",
        })
        return
    }
    
    // إنشاء اسم ملف
    filename := fmt.Sprintf("nawthtech_video_%s.%s", jobID, job.Result.Format)
    if job.Result.Format == "" {
        filename = fmt.Sprintf("nawthtech_video_%s.mp4", jobID)
    }
    
    // تعيين رؤوس الاستجابة
    contentType := "video/mp4"
    if job.Result.Format == "webm" {
        contentType = "video/webm"
    } else if job.Result.Format == "gif" {
        contentType = "image/gif"
    }
    
    c.Header("Content-Type", contentType)
    c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
    c.Header("Content-Length", fmt.Sprintf("%d", len(job.Result.VideoData)))
    c.Header("Cache-Control", "public, max-age=31536000") // تخزين لمدة سنة
    
    // إرسال البيانات
    c.Data(http.StatusOK, contentType, job.Result.VideoData)
}

// GetVideoCapabilitiesHandler معالج قدرات توليد الفيديو
func (h *VideoHandler) GetVideoCapabilitiesHandler(c *gin.Context) {
    providerStats := h.videoService.GetProviderStats()
    
    capabilities := gin.H{
        "video_types": []gin.H{
            {
                "id":            "short",
                "name":          "Short Video",
                "description":   "Short videos for social media",
                "duration":      15,
                "max_duration":  60,
            },
            {
                "id":            "explainer",
                "name":          "Explainer Video", 
                "description":   "Educational and explanatory videos",
                "duration":      60,
                "max_duration":  300,
            },
            {
                "id":            "promotional",
                "name":          "Promotional Video",
                "description":   "Marketing and promotional content",
                "duration":      30,
                "max_duration":  120,
            },
        },
        
        "supported_resolutions": []string{
            "512x512", "576x1024", "1024x576",
            "768x768", "1024x1024", "1280x720",
        },
        
        "supported_aspects": []string{
            "1:1", "16:9", "9:16", "4:3", "21:9",
        },
        
        "supported_styles": []string{
            "realistic", "anime", "cartoon", 
            "artistic", "cinematic", "minimal",
        },
        
        "supported_formats": []string{"mp4", "webm", "gif"},
        
        "limits": gin.H{
            "max_duration":      60,    // ثواني
            "min_duration":      1,     // ثواني
            "max_resolution":    "1024x1024",
            "max_prompt_length": 500,
        },
        
        "features": []string{
            "text_to_video",
            "style_transfer",
            "aspect_ratio_conversion",
        },
        
        "provider": providerStats,
    }
    
    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "data":    capabilities,
    })
}

// GetVideoStatsHandler معالج إحصائيات الفيديو
func (h *VideoHandler) GetVideoStatsHandler(c *gin.Context) {
    stats := h.videoService.GetStats()
    
    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "data": gin.H{
            "total_generations":  stats.TotalGenerations,
            "successful":         stats.Successful,
            "failed":             stats.Failed,
            "total_duration":     stats.TotalDuration,
            "total_cost":         stats.TotalCost,
            "last_generation":    stats.LastGeneration.Format(time.RFC3339),
            "most_used_style":    stats.MostUsedStyle,
            "most_used_provider": stats.MostUsedProvider,
        },
    })
}

// GetVideoUsageHandler معالج استخدام الفيديو
func (h *VideoHandler) GetVideoUsageHandler(c *gin.Context) {
    userID := h.getUserID(c)
    userTier := h.getUserTier(c)
    
    usage := h.videoService.GetUserUsage(userID, userTier)
    
    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "data": gin.H{
            "user_id":           usage.UserID,
            "tier":              usage.Tier,
            "total_generations": usage.TotalGenerations,
            "monthly_limit":     usage.MonthlyLimit,
            "monthly_used":      usage.MonthlyUsed,
            "daily_limit":       usage.DailyLimit,
            "daily_used":        usage.DailyUsed,
            "last_generated":    usage.LastGenerated.Format(time.RFC3339),
            "last_reset":        usage.LastReset.Format(time.RFC3339),
            "can_generate":      usage.MonthlyUsed < usage.MonthlyLimit && usage.DailyUsed < usage.DailyLimit,
        },
    })
}

// UploadImageForVideoHandler معالج رفع صورة لتحويلها إلى فيديو
func (h *VideoHandler) UploadImageForVideoHandler(c *gin.Context) {
    // TODO: تنفيذ رفع صورة وتحويلها إلى فيديو
    c.JSON(http.StatusNotImplemented, gin.H{
        "success": false,
        "error":   "Image to video conversion not implemented yet",
    })
}

// Helper functions
func (h *VideoHandler) getUserID(c *gin.Context) string {
    // محاولة استخراج userID من السياق
    if userID, exists := c.Get("userID"); exists {
        if id, ok := userID.(string); ok {
            return id
        }
    }
    
    // محاولة من الرؤوس
    if userID := c.GetHeader("X-User-ID"); userID != "" {
        return userID
    }
    
    // في وضع التطوير/الاختبار، يمكن استخدام معرف وهمي
    return "test_user"
}

func (h *VideoHandler) getUserTier(c *gin.Context) string {
    // محاولة استخراج tier من السياق
    if userTier, exists := c.Get("userTier"); exists {
        if tier, ok := userTier.(string); ok {
            return tier
        }
    }
    
    // محاولة من الرؤوس
    if userTier := c.GetHeader("X-User-Tier"); userTier != "" {
        return userTier
    }
    
    // قيمة افتراضية
    return "free"
}