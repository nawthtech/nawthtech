package video

import (
    "context"
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
    
    stats := VideoStats{}
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
        
        if job.CreatedAt.After(stats.LastGeneration) {
            stats.LastGeneration = job.CreatedAt
        }
    }
    
    stats.Provider = s.provider.Name()
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