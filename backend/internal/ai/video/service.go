package video

import (
    "fmt"
    "sync"
    "time"
)

// VideoService خدمة إدارة توليد الفيديو
type VideoService struct {
    provider      VideoProvider
    jobs          map[string]*VideoJob
    mu            sync.RWMutex
}

// VideoJob مهمة فيديو
type VideoJob struct {
    ID        string         `json:"id"`
    Status    VideoJobStatus `json:"status"`
    Progress  int            `json:"progress"` // 0-100
    Result    *VideoResponse `json:"result,omitempty"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    Request   VideoRequest   `json:"request,omitempty"`
}

// NewVideoService إنشاء خدمة فيديو جديدة
func NewVideoService(provider VideoProvider) *VideoService {
    return &VideoService{
        provider: provider,
        jobs:     make(map[string]*VideoJob),
    }
}

// SubmitVideoJob تقديم طلب فيديو جديد
func (s *VideoService) SubmitVideoJob(req VideoRequest) (*VideoJob, error) {
    // التحقق من صحة الطلب
    if err := ValidateVideoRequest(req); err != nil {
        return nil, err
    }
    
    jobID := generateJobID()
    
    job := &VideoJob{
        ID:        jobID,
        Status:    VideoJobPending,
        Progress:  0,
        Request:   req,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }
    
    s.mu.Lock()
    s.jobs[jobID] = job
    s.mu.Unlock()
    
    // تشغيل في goroutine
    go s.processVideoJob(jobID, req)
    
    return job, nil
}

// processVideoJob معالجة طلب الفيديو
func (s *VideoService) processVideoJob(jobID string, req VideoRequest) {
    s.updateJobStatus(jobID, VideoJobProcessing, 10)
    
    // توليد الفيديو
    response, err := s.provider.GenerateVideo(req)
    
    s.mu.Lock()
    defer s.mu.Unlock()
    
    job, exists := s.jobs[jobID]
    if !exists {
        return
    }
    
    job.UpdatedAt = time.Now()
    
    if err != nil {
        job.Status = VideoJobFailed
        job.Progress = 0
        
        // تحويل الخطأ إلى VideoResponse
        videoErr, ok := err.(*VideoError)
        if !ok {
            videoErr = &VideoError{
                Code:    "unknown_error",
                Message: err.Error(),
            }
        }
        
        job.Result = &VideoResponse{
            Success: false,
            Error:   videoErr.Message,
            Status:  string(VideoJobFailed),
            CreatedAt: time.Now(),
        }
    } else {
        job.Status = VideoJobCompleted
        job.Progress = 100
        job.Result = response
    }
}

// GetJob الحصول على تفاصيل job
func (s *VideoService) GetJob(jobID string) (*VideoJob, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    
    job, exists := s.jobs[jobID]
    if !exists {
        return nil, fmt.Errorf("job not found: %s", jobID)
    }
    
    return job, nil
}

// ListJobs عرض جميع المهام
func (s *VideoService) ListJobs() []*VideoJob {
    s.mu.RLock()
    defer s.mu.RUnlock()
    
    jobs := make([]*VideoJob, 0, len(s.jobs))
    for _, job := range s.jobs {
        jobs = append(jobs, job)
    }
    
    return jobs
}

// CancelJob إلغاء مهمة
func (s *VideoService) CancelJob(jobID string) error {
    s.mu.Lock()
    defer s.mu.Unlock()
    
    job, exists := s.jobs[jobID]
    if !exists {
        return fmt.Errorf("job not found: %s", jobID)
    }
    
    if job.Status == VideoJobCompleted || job.Status == VideoJobFailed {
        return fmt.Errorf("cannot cancel job with status: %s", job.Status)
    }
    
    job.Status = VideoJobCancelled
    job.Progress = 0
    job.UpdatedAt = time.Now()
    
    return nil
}

// CleanupJobs تنظيف المهام القديمة
func (s *VideoService) CleanupJobs(maxAge time.Duration) int {
    s.mu.Lock()
    defer s.mu.Unlock()
    
    cutoff := time.Now().Add(-maxAge)
    deleted := 0
    
    for id, job := range s.jobs {
        if job.CreatedAt.Before(cutoff) {
            delete(s.jobs, id)
            deleted++
        }
    }
    
    return deleted
}

// GetStats الحصول على إحصائيات الخدمة
func (s *VideoService) GetStats() VideoStats {
    s.mu.RLock()
    defer s.mu.RUnlock()
    
    stats := VideoStats{
        LastGeneration: time.Time{},
        Provider:       s.provider.Name(),
    }
    
    styleCount := make(map[string]int)
    providerCount := make(map[string]int)
    
    for _, job := range s.jobs {
        stats.TotalGenerations++
        
        switch job.Status {
        case VideoJobCompleted:
            stats.Successful++
            if job.Result != nil {
                stats.TotalDuration += int64(job.Result.Duration)
                stats.TotalCost += job.Result.Cost
            }
        case VideoJobFailed:
            stats.Failed++
        }
        
        // تحديث آخر عملية توليد
        if job.CreatedAt.After(stats.LastGeneration) {
            stats.LastGeneration = job.CreatedAt
        }
        
        // عد الأنماط الأكثر استخداماً
        if job.Request.Style != "" {
            styleCount[job.Request.Style]++
        }
        
        // عد المزودين الأكثر استخداماً
        if job.Result != nil && job.Result.Provider != "" {
            providerCount[job.Result.Provider]++
        }
    }
    
    // العثور على النمط الأكثر استخداماً
    if len(styleCount) > 0 {
        maxCount := 0
        for style, count := range styleCount {
            if count > maxCount {
                maxCount = count
                stats.MostUsedStyle = style
            }
        }
    }
    
    // العثور على المزود الأكثر استخداماً
    if len(providerCount) > 0 {
        maxCount := 0
        for provider, count := range providerCount {
            if count > maxCount {
                maxCount = count
                stats.MostUsedProvider = provider
            }
        }
    }
    
    return stats
}

func (s *VideoService) updateJobStatus(jobID string, status VideoJobStatus, progress int) {
    s.mu.Lock()
    defer s.mu.Unlock()
    
    if job, exists := s.jobs[jobID]; exists {
        job.Status = status
        job.Progress = progress
        job.UpdatedAt = time.Now()
    }
}

func generateJobID() string {
    return fmt.Sprintf("video_%d", time.Now().UnixNano())
}

// GetProviderStats الحصول على إحصائيات المزود
func (s *VideoService) GetProviderStats() map[string]interface{} {
    return map[string]interface{}{
        "provider":          s.provider.Name(),
        "available":         s.provider.IsAvailable(),
        "local":             s.provider.IsLocal(),
        "free":              s.provider.IsFree(),
        "supported_resolutions": getSupportedResolutionsFromProvider(s.provider),
    }
}

// getSupportedResolutionsFromProvider الحصول على الدقات المدعومة من المزود
func getSupportedResolutionsFromProvider(provider VideoProvider) []string {
    // هذه قائمة افتراضية، يمكن للمزودين تخصيصها
    resolutions := []string{
        "512x512", "576x1024", "1024x576",
        "768x768", "1024x1024", "1280x720",
    }
    
    // تصفية الدقات المدعومة فعلياً
    var supported []string
    for _, res := range resolutions {
        if provider.SupportsResolution(res) {
            supported = append(supported, res)
        }
    }
    
    return supported
}

// GetJobSummary الحصول على ملخص المهام
func (s *VideoService) GetJobSummary() map[string]interface{} {
    s.mu.RLock()
    defer s.mu.RUnlock()
    
    summary := map[string]interface{}{
        "total_jobs":     len(s.jobs),
        "pending":        0,
        "processing":     0,
        "completed":      0,
        "failed":         0,
        "cancelled":      0,
    }
    
    for _, job := range s.jobs {
        switch job.Status {
        case VideoJobPending:
            summary["pending"] = summary["pending"].(int) + 1
        case VideoJobProcessing:
            summary["processing"] = summary["processing"].(int) + 1
        case VideoJobCompleted:
            summary["completed"] = summary["completed"].(int) + 1
        case VideoJobFailed:
            summary["failed"] = summary["failed"].(int) + 1
        case VideoJobCancelled:
            summary["cancelled"] = summary["cancelled"].(int) + 1
        }
    }
    
    return summary
}

// GetRecentJobs الحصول على المهام الحديثة
func (s *VideoService) GetRecentJobs(limit int) []*VideoJob {
    s.mu.RLock()
    defer s.mu.RUnlock()
    
    if limit <= 0 || limit > len(s.jobs) {
        limit = len(s.jobs)
    }
    
    // جمع جميع المهام
    allJobs := make([]*VideoJob, 0, len(s.jobs))
    for _, job := range s.jobs {
        allJobs = append(allJobs, job)
    }
    
    // ترتيب حسب التاريخ (الأحدث أولاً)
    for i := 0; i < len(allJobs)-1; i++ {
        for j := i + 1; j < len(allJobs); j++ {
            if allJobs[i].CreatedAt.Before(allJobs[j].CreatedAt) {
                allJobs[i], allJobs[j] = allJobs[j], allJobs[i]
            }
        }
    }
    
    // إرجاع المهام المطلوبة فقط
    if len(allJobs) > limit {
        return allJobs[:limit]
    }
    
    return allJobs
}

// GetJobsByStatus الحصول على المهام حسب الحالة
func (s *VideoService) GetJobsByStatus(status VideoJobStatus) []*VideoJob {
    s.mu.RLock()
    defer s.mu.RUnlock()
    
    var filteredJobs []*VideoJob
    for _, job := range s.jobs {
        if job.Status == status {
            filteredJobs = append(filteredJobs, job)
        }
    }
    
    return filteredJobs
}

// GetJobsByDateRange الحصول على المهام ضمن نطاق تاريخ
func (s *VideoService) GetJobsByDateRange(start, end time.Time) []*VideoJob {
    s.mu.RLock()
    defer s.mu.RUnlock()
    
    var filteredJobs []*VideoJob
    for _, job := range s.jobs {
        if !job.CreatedAt.Before(start) && !job.CreatedAt.After(end) {
            filteredJobs = append(filteredJobs, job)
        }
    }
    
    return filteredJobs
}

// IsJobOwner التحقق من ملكية المهمة
func (s *VideoService) IsJobOwner(jobID, userID string) bool {
    s.mu.RLock()
    defer s.mu.RUnlock()
    
    job, exists := s.jobs[jobID]
    if !exists {
        return false
    }
    
    // هذا مثال بسيط، في التطبيق الحقيقي قد تحتاج إلى التحقق من قاعدة البيانات
    return job.Request.UserID == userID
}