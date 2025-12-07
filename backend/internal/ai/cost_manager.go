package ai

import (
    "encoding/json"
    "fmt"
    "os"
    "path/filepath"
    "sync"
    "time"
    "github.com/nawthtech/nawthtech/backend/internal/ai/types"
)

// Quota حصة المستخدم
type Quota struct {
    Type          string    `json:"type"`           // text, image, video, audio
    Used          int64     `json:"used"`           // الكمية المستخدمة
    Limit         int64     `json:"limit"`          // الحد الأقصى
    ResetPeriod   string    `json:"reset_period"`   // daily, weekly, monthly
    LastReset     time.Time `json:"last_reset"`
}

// NewCostManager إنشاء CostManager جديد
func NewCostManager() (*CostManager, error) {
    dataPath := os.Getenv("AI_DATA_PATH")
    if dataPath == "" {
        dataPath = "./data/ai"
    }
    
    cm := &CostManager{
        dataPath:     dataPath,
        monthlyLimit: 0.0,  // 0 = لا يوجد حد (مجاني)
        dailyLimit:   0.0,
    }
    
    // تهيئة البيانات
    cm.Usage.MonthlyCost = make(map[string]float64)
    cm.Usage.DailyCost = make(map[string]float64)
    cm.Usage.UserUsage = make(map[string]*UserUsageStats)
    cm.Usage.ProviderUsage = make(map[string]*ProviderStats)
    cm.Usage.LastReset = time.Now()
    
    // تحميل البيانات المحفوظة
    if err := cm.load(); err != nil {
        fmt.Printf("Warning: Could not load cost data: %v\n", err)
    }
    
    // جدولة إعادة التعيين التلقائي
    go cm.startAutoReset()
    
    return cm, nil
}

// load تحميل البيانات المحفوظة
func (cm *CostManager) load() error {
    dataFile := filepath.Join(cm.dataPath, "cost_data.json")
    if _, err := os.Stat(dataFile); os.IsNotExist(err) {
        return nil
    }
    
    data, err := os.ReadFile(dataFile)
    if err != nil {
        return err
    }
    
    cm.mu.Lock()
    defer cm.mu.Unlock()
    
    return json.Unmarshal(data, &cm.Usage)
}

// save حفظ البيانات
func (cm *CostManager) save() error {
    cm.mu.RLock()
    defer cm.mu.RUnlock()
    
    if err := os.MkdirAll(cm.dataPath, 0755); err != nil {
        return err
    }
    
    dataFile := filepath.Join(cm.dataPath, "cost_data.json")
    data, err := json.MarshalIndent(cm.Usage, "", "  ")
    if err != nil {
        return err
    }
    
    return os.WriteFile(dataFile, data, 0644)
}

// startAutoReset بدء إعادة التعيين التلقائي
func (cm *CostManager) startAutoReset() {
    // إعادة تعيين يومية في منتصف الليل
    go func() {
        for {
            now := time.Now()
            nextMidnight := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
            durationUntilMidnight := nextMidnight.Sub(now)
            
            time.Sleep(durationUntilMidnight)
            cm.resetDailyQuotas()
        }
    }()
    
    // إعادة تعيين شهرية في أول يوم من الشهر
    go func() {
        for {
            now := time.Now()
            nextMonth := time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, now.Location())
            durationUntilNextMonth := nextMonth.Sub(now)
            
            time.Sleep(durationUntilNextMonth)
            cm.resetMonthlyQuotas()
        }
    }()
}

// resetDailyQuotas إعادة تعيين الحصص اليومية
func (cm *CostManager) resetDailyQuotas() {
    cm.mu.Lock()
    defer cm.mu.Unlock()
    
    now := time.Now()
    today := now.Format("2006-01-02")
    
    // إعادة تعيين التكلفة اليومية
    cm.Usage.DailyCost[today] = 0
    
    // إعادة تعيين حصص المستخدمين اليومية
    for _, userStats := range cm.Usage.UserUsage {
        userStats.DailyCost[today] = 0
        for _, quota := range userStats.Quotas {
            if quota.ResetPeriod == "daily" {
                quota.Used = 0
                quota.LastReset = now
            }
        }
    }
    
    cm.Usage.LastReset = now
    go cm.save()
}

// resetMonthlyQuotas إعادة تعيين الحصص الشهرية
func (cm *CostManager) resetMonthlyQuotas() {
    cm.mu.Lock()
    defer cm.mu.Unlock()
    
    now := time.Now()
    monthKey := now.Format("2006-01")
    
    // إعادة تعيين التكلفة الشهرية
    cm.Usage.MonthlyCost[monthKey] = 0
    
    // إعادة تعيين حصص المستخدمين الشهرية
    for _, userStats := range cm.Usage.UserUsage {
        userStats.MonthlyCost[monthKey] = 0
        for _, quota := range userStats.Quotas {
            if quota.ResetPeriod == "monthly" {
                quota.Used = 0
                quota.LastReset = now
            }
        }
    }
    
    cm.Usage.LastReset = now
    go cm.save()
}

// getDefaultQuotas الحصول على الحصص الافتراضية حسب الطبقة
func (cm *CostManager) getDefaultQuotas(tier string) map[string]*Quota {
    quotas := make(map[string]*Quota)
    now := time.Now()
    
    switch tier {
    case "premium":
        quotas["text"] = &Quota{
            Type:        "text",
            Used:        0,
            Limit:       100000,
            ResetPeriod: "monthly",
            LastReset:   now,
        }
        quotas["image"] = &Quota{
            Type:        "image",
            Used:        0,
            Limit:       1000,
            ResetPeriod: "monthly",
            LastReset:   now,
        }
        quotas["video"] = &Quota{
            Type:        "video",
            Used:        0,
            Limit:       100,
            ResetPeriod: "monthly",
            LastReset:   now,
        }
    case "basic":
        quotas["text"] = &Quota{
            Type:        "text",
            Used:        0,
            Limit:       10000,
            ResetPeriod: "monthly",
            LastReset:   now,
        }
        quotas["image"] = &Quota{
            Type:        "image",
            Used:        0,
            Limit:       100,
            ResetPeriod: "monthly",
            LastReset:   now,
        }
    default: // free
        quotas["text"] = &Quota{
            Type:        "text",
            Used:        0,
            Limit:       1000,
            ResetPeriod: "daily",
            LastReset:   now,
        }
        quotas["image"] = &Quota{
            Type:        "image",
            Used:        0,
            Limit:       10,
            ResetPeriod: "daily",
            LastReset:   now,
        }
    }
    
    return quotas
}

// RecordUsage تسجيل استخدام
func (cm *CostManager) RecordUsage(record *types.UsageRecord) error {
    cm.mu.Lock()
    defer cm.mu.Unlock()
    
    now := time.Now()
    monthKey := now.Format("2006-01")
    dayKey := now.Format("2006-01-02")
    
    // تحديث التكلفة الإجمالية
    cm.Usage.TotalCost += record.Cost
    cm.Usage.MonthlyCost[monthKey] += record.Cost
    cm.Usage.DailyCost[dayKey] += record.Cost
    
    // تحديث إحصائيات المستخدم
    if record.UserID != "" {
        if _, exists := cm.Usage.UserUsage[record.UserID]; !exists {
            cm.Usage.UserUsage[record.UserID] = &UserUsageStats{
                UserID:      record.UserID,
                Tier:        record.UserTier,
                MonthlyCost: make(map[string]float64),
                DailyCost:   make(map[string]float64),
                Quotas:      cm.getDefaultQuotas(record.UserTier),
            }
        }
        
        userStats := cm.Usage.UserUsage[record.UserID]
        userStats.TotalCost += record.Cost
        userStats.MonthlyCost[monthKey] += record.Cost
        userStats.DailyCost[dayKey] += record.Cost
        userStats.LastActive = now
        
        // تحديث الحصص
        if quota, exists := userStats.Quotas[record.Type]; exists {
            quota.Used += record.Quantity
        }
    }
    
    // تحديث إحصائيات المزود
    if _, exists := cm.Usage.ProviderUsage[record.Provider]; !exists {
        cm.Usage.ProviderUsage[record.Provider] = &ProviderStats{
            ProviderName: record.Provider,
        }
    }
    
    providerStats := cm.Usage.ProviderUsage[record.Provider]
    providerStats.TotalRequests++
    providerStats.TotalCost += record.Cost
    providerStats.LastUsed = now
    
    // تحديث متوسط زمن الاستجابة
    if record.Latency > 0 {
        if providerStats.AvgLatency == 0 {
            providerStats.AvgLatency = record.Latency
        } else {
            providerStats.AvgLatency = (providerStats.AvgLatency*float64(providerStats.TotalRequests-1) + record.Latency) / float64(providerStats.TotalRequests)
        }
    }
    
    // تحديث نسبة النجاح
    if record.Success {
        providerStats.SuccessRate = (providerStats.SuccessRate*float64(providerStats.TotalRequests-1) + 1.0) / float64(providerStats.TotalRequests)
    } else {
        providerStats.SuccessRate = (providerStats.SuccessRate * float64(providerStats.TotalRequests-1)) / float64(providerStats.TotalRequests)
    }
    
    // التحقق من تجاوز الحدود
    if cm.monthlyLimit > 0 && cm.Usage.MonthlyCost[monthKey] > cm.monthlyLimit {
        return fmt.Errorf("monthly cost limit exceeded: %.2f/%.2f", 
            cm.Usage.MonthlyCost[monthKey], cm.monthlyLimit)
    }
    
    if cm.dailyLimit > 0 && cm.Usage.DailyCost[dayKey] > cm.dailyLimit {
        return fmt.Errorf("daily cost limit exceeded: %.2f/%.2f", 
            cm.Usage.DailyCost[dayKey], cm.dailyLimit)
    }
    
    // الحفظ التلقائي
    go cm.save()
    
    return nil
}

// CanUseAI التحقق من إمكانية استخدام AI
func (cm *CostManager) CanUseAI(userID, requestType string) (bool, string) {
    cm.mu.RLock()
    defer cm.mu.RUnlock()
    
    // إذا لم يكن هناك حدود، السماح دائماً
    if cm.monthlyLimit == 0 && cm.dailyLimit == 0 {
        return true, ""
    }
    
    now := time.Now()
    monthKey := now.Format("2006-01")
    dayKey := now.Format("2006-01-02")
    
    // التحقق من الحدود العامة
    if cm.monthlyLimit > 0 && cm.Usage.MonthlyCost[monthKey] >= cm.monthlyLimit {
        return false, "تم تجاوز الحد الشهري للتكاليف"
    }
    
    if cm.dailyLimit > 0 && cm.Usage.DailyCost[dayKey] >= cm.dailyLimit {
        return false, "تم تجاوز الحد اليومي للتكاليف"
    }
    
    // التحقق من حصص المستخدم
    if userID != "" {
        if userStats, exists := cm.Usage.UserUsage[userID]; exists {
            // التحقق من التكلفة الشهرية للمستخدم
            if monthlyCost, ok := userStats.MonthlyCost[monthKey]; ok {
                // حساب الحد الشهري حسب الطبقة
                var userMonthlyLimit float64
                switch userStats.Tier {
                case "premium":
                    userMonthlyLimit = 100.0 // 100 دولار
                case "basic":
                    userMonthlyLimit = 10.0  // 10 دولار
                default:
                    userMonthlyLimit = 1.0   // 1 دولار
                }
                
                if monthlyCost >= userMonthlyLimit {
                    return false, "تم تجاوز حد التكلفة الشهرية للمستخدم"
                }
            }
            
            // التحقق من حصص النوع المحدد
            if quota, exists := userStats.Quotas[requestType]; exists {
                if quota.Used >= quota.Limit {
                    return false, fmt.Sprintf("تم تجاوز حصة %s لهذا المستخدم", requestType)
                }
            }
        }
    }
    
    return true, ""
}

// GetUsageStatistics الحصول على إحصائيات الاستخدام
func (cm *CostManager) GetUsageStatistics() map[string]interface{} {
    cm.mu.RLock()
    defer cm.mu.RUnlock()
    
    stats := make(map[string]interface{})
    stats["total_cost"] = cm.Usage.TotalCost
    stats["monthly_cost"] = cm.Usage.MonthlyCost
    stats["daily_cost"] = cm.Usage.DailyCost
    stats["total_users"] = len(cm.Usage.UserUsage)
    stats["providers"] = len(cm.Usage.ProviderUsage)
    stats["last_reset"] = cm.Usage.LastReset
    
    // حساب التكلفة الشهرية الحالية
    monthKey := time.Now().Format("2006-01")
    if monthly, ok := cm.Usage.MonthlyCost[monthKey]; ok {
        stats["current_month_cost"] = monthly
        if cm.monthlyLimit > 0 {
            stats["monthly_usage_percentage"] = (monthly / cm.monthlyLimit) * 100
        }
    }
    
    // حساب التكلفة اليومية الحالية
    dayKey := time.Now().Format("2006-01-02")
    if daily, ok := cm.Usage.DailyCost[dayKey]; ok {
        stats["current_day_cost"] = daily
        if cm.dailyLimit > 0 {
            stats["daily_usage_percentage"] = (daily / cm.dailyLimit) * 100
        }
    }
    
    return stats
}

// SetLimits تعيين الحدود
func (cm *CostManager) SetLimits(monthly, daily float64) {
    cm.mu.Lock()
    defer cm.mu.Unlock()
    
    cm.monthlyLimit = monthly
    cm.dailyLimit = daily
    
    go cm.save()
}

// ResetUsage إعادة تعيين جميع الإحصائيات
func (cm *CostManager) ResetUsage() error {
    cm.mu.Lock()
    defer cm.mu.Unlock()
    
    cm.Usage.TotalCost = 0
    cm.Usage.MonthlyCost = make(map[string]float64)
    cm.Usage.DailyCost = make(map[string]float64)
    cm.Usage.UserUsage = make(map[string]*UserUsageStats)
    cm.Usage.ProviderUsage = make(map[string]*ProviderStats)
    cm.Usage.LastReset = time.Now()
    
    return cm.save()
}