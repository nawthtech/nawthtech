package main

import (
    "log"
    "os"
    
    "github.com/gin-gonic/gin"
    "github.com/nawthtech/nawthtech/backend/internal/ai/video"
)

func main() {
    // إنشاء مزود فيديو هجين
    provider := video.NewHybridVideoProvider()
    
    // إنشاء خدمة فيديو
    videoService := video.NewVideoService(provider)
    
    r := gin.Default()
    
    // مسارات API
    api := r.Group("/api")
    
    // مسار توليد الفيديو
    api.POST("/video/generate", func(c *gin.Context) {
        var req struct {
            Prompt         string `json:"prompt" binding:"required"`
            Duration       int    `json:"duration" default:"5"`
            Resolution     string `json:"resolution" default:"512x512"`
            Aspect         string `json:"aspect" default:"1:1"`
            Style          string `json:"style" default:"realistic"`
            UserID         string `json:"user_id"`
            Tier           string `json:"tier" default:"free"`
        }
        
        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(400, gin.H{"error": err.Error()})
            return
        }
        
        // إنشاء طلب الفيديو
        videoReq := video.VideoRequest{
            Prompt:     req.Prompt,
            Duration:   req.Duration,
            Resolution: req.Resolution,
            Aspect:     req.Aspect,
            Style:      req.Style,
            UserID:     req.UserID,
            UserTier:   req.Tier,
        }
        
        // التحقق من صحة الطلب
        if err := video.ValidateVideoRequest(videoReq); err != nil {
            c.JSON(400, gin.H{"error": err.Error()})
            return
        }
        
        // التحقق من إمكانية المستخدم لتوليد فيديو
        if canGenerate, message := videoService.CanUserGenerateVideo(req.UserID, req.Tier); !canGenerate {
            c.JSON(403, gin.H{
                "error": message,
            })
            return
        }
        
        // إرسال طلب توليد الفيديو
        job, err := videoService.SubmitVideoJob(videoReq)
        if err != nil {
            c.JSON(500, gin.H{"error": err.Error()})
            return
        }
        
        // تسجيل استخدام المستخدم
        videoService.RecordUserGeneration(req.UserID, req.Tier)
        
        c.JSON(202, gin.H{
            "success": true,
            "job_id":  job.ID,
            "status":  job.Status,
            "message": "Video generation started",
        })
    })
    
    // مسار حالة الفيديو
    api.GET("/video/status/:jobId", func(c *gin.Context) {
        jobID := c.Param("jobId")
        
        job, err := videoService.GetJob(jobID)
        if err != nil {
            c.JSON(404, gin.H{"error": "Job not found"})
            return
        }
        
        c.JSON(200, gin.H{
            "job_id":    job.ID,
            "status":    job.Status,
            "progress":  job.Progress,
            "result":    job.Result,
        })
    })
    
    // مسار تحميل الفيديو
    api.GET("/video/download/:jobId", func(c *gin.Context) {
        jobID := c.Param("jobId")
        
        job, err := videoService.GetJob(jobID)
        if err != nil {
            c.JSON(404, gin.H{"error": "Job not found"})
            return
        }
        
        if job.Status != video.VideoJobCompleted || job.Result == nil {
            c.JSON(400, gin.H{"error": "Video not ready for download"})
            return
        }
        
        if job.Result.VideoURL != "" {
            c.Redirect(302, job.Result.VideoURL)
            return
        }
        
        if len(job.Result.VideoData) > 0 {
            c.Header("Content-Type", "video/mp4")
            c.Data(200, "video/mp4", job.Result.VideoData)
            return
        }
        
        c.JSON(400, gin.H{"error": "No video data available"})
    })
    
    // مسار إحصائيات
    api.GET("/video/stats", func(c *gin.Context) {
        stats := videoService.GetStats()
        c.JSON(200, gin.H{"stats": stats})
    })
    
    // مسار قدرات المزود
    api.GET("/video/capabilities", func(c *gin.Context) {
        capabilities := gin.H{
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
            "max_duration": 60,
        }
        
        c.JSON(200, gin.H{"capabilities": capabilities})
    })
    
    port := os.Getenv("PORT")
    if port == "" {
        port = "8081"
    }
    
    log.Printf("🎬 Free Video Generation Server running on port %s", port)
    log.Fatal(r.Run(":" + port))
}