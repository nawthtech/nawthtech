package models

import (
	"time"
)

// ================================
// النماذج الأساسية (Core Models)
// ================================

// User نموذج المستخدم
type User struct {
	ID              string    `json:"id" gorm:"primaryKey"`
	Email           string    `json:"email" gorm:"uniqueIndex;not null"`
	Username        string    `json:"username" gorm:"uniqueIndex;not null"`
	FirstName       string    `json:"first_name" gorm:"not null"`
	LastName        string    `json:"last_name" gorm:"not null"`
	Phone           string    `json:"phone,omitempty"`
	Avatar          string    `json:"avatar,omitempty"`
	Role            string    `json:"role" gorm:"default:'user';index"`
	Status          string    `json:"status" gorm:"default:'active';index"`
	EmailVerified   bool      `json:"email_verified" gorm:"default:false"`
	PhoneVerified   bool      `json:"phone_verified" gorm:"default:false"`
	LastLogin       time.Time `json:"last_login,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// Service نموذج الخدمة
type Service struct {
	ID           string    `json:"id" gorm:"primaryKey"`
	Title        string    `json:"title" gorm:"not null;index"`
	Description  string    `json:"description" gorm:"type:text"`
	Category     string    `json:"category" gorm:"not null;index"`
	Price        float64   `json:"price" gorm:"not null"`
	Duration     int       `json:"duration" gorm:"not null"`
	Rating       float64   `json:"rating" gorm:"default:0;index"`
	TotalOrders  int       `json:"total_orders" gorm:"default:0"`
	TotalReviews int       `json:"total_reviews" gorm:"default:0"`
	Status       string    `json:"status" gorm:"default:'active';index"`
	Featured     bool      `json:"featured" gorm:"default:false;index"`
	SellerID     string    `json:"seller_id" gorm:"not null;index"`
	SellerName   string    `json:"seller_name,omitempty" gorm:"-"`
	Images       []string  `json:"images" gorm:"type:json;serializer:json"`
	Features     []string  `json:"features" gorm:"type:json;serializer:json"`
	Tags         []string  `json:"tags" gorm:"type:json;serializer:json"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Content نموذج المحتوى
type Content struct {
	ID           string                 `json:"id"`
	Topic        string                 `json:"topic"`
	Content      string                 `json:"content"`
	Platform     string                 `json:"platform"`
	ContentType  string                 `json:"content_type"`
	Tone         string                 `json:"tone"`
	Keywords     []string               `json:"keywords"`
	Language     string                 `json:"language"`
	Status       string                 `json:"status"`
	Analysis     *ContentAnalysis       `json:"analysis"`
	Optimization *PlatformOptimization  `json:"optimization"`
	Metadata     map[string]interface{} `json:"metadata"`
	Performance  *ContentPerformance    `json:"performance"`
	CreatedBy    string                 `json:"created_by"`
	CreatedAt    time.Time              `json:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at"`
}

// Notification نموذج الإشعار
type Notification struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Title     string    `json:"title"`
	Message   string    `json:"message"`
	Type      string    `json:"type"`
	Priority  string    `json:"priority"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// Review نموذج التقييم
type Review struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	Type        string    `json:"type" gorm:"not null;index"`
	TargetID    string    `json:"target_id" gorm:"not null;index"`
	UserID      string    `json:"user_id" gorm:"not null;index"`
	UserName    string    `json:"user_name" gorm:"-"`
	UserAvatar  string    `json:"user_avatar,omitempty" gorm:"-"`
	Rating      int       `json:"rating" gorm:"not null;check:rating>=1 AND rating<=5"`
	Title       string    `json:"title,omitempty" gorm:"size:200"`
	Comment     string    `json:"comment" gorm:"type:text"`
	IsVerified  bool      `json:"is_verified" gorm:"default:false"`
	Helpful     int       `json:"helpful" gorm:"default:0"`
	Reported    bool      `json:"reported" gorm:"default:false"`
	ReportReason string   `json:"report_reason,omitempty"`
	Status      string    `json:"status" gorm:"default:'active';index"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ================================
// نماذج التحليلات (Analytics Models)
// ================================

// AnalyticsOverview نظرة عامة على التحليلات
type AnalyticsOverview struct {
	Summary     *AnalyticsSummary `json:"summary"`
	Comparison  *ComparisonData   `json:"comparison"`
	Trends      *TrendsData       `json:"trends"`
	GeneratedAt time.Time         `json:"generated_at"`
}

// AnalyticsSummary ملخص التحليلات
type AnalyticsSummary struct {
	TotalVisitors     int     `json:"total_visitors"`
	TotalEngagement   float64 `json:"total_engagement"`
	TotalReach        int     `json:"total_reach"`
	ConversionRate    float64 `json:"conversion_rate"`
	GrowthRate        float64 `json:"growth_rate"`
	ActiveUsers       int     `json:"active_users"`
}

// PerformanceAnalytics تحليلات الأداء
type PerformanceAnalytics struct {
	Timeframe string              `json:"timeframe"`
	Platform  string              `json:"platform"`
	Metrics   string              `json:"metrics"`
	Data      []PerformanceMetric `json:"data"`
	Summary   *PerformanceSummary `json:"summary"`
	GeneratedAt time.Time         `json:"generated_at"`
}

// AIInsights رؤى الذكاء الاصطناعي
type AIInsights struct {
	Trends           *TrendInsights         `json:"trends"`
	Predictions      *PredictionInsights    `json:"predictions"`
	Recommendations  *RecommendationInsights `json:"recommendations"`
	OptimizationScore int                   `json:"optimization_score"`
	Confidence       int                    `json:"confidence"`
	DataSummary      *InsightsDataSummary   `json:"data_summary"`
	GeneratedAt      time.Time              `json:"generated_at"`
}

// ================================
// نماذج الإدارة (Admin Models)
// ================================

// DashboardData بيانات لوحة التحكم
type DashboardData struct {
	Stats         DashboardStats      `json:"stats"`
	StoreMetrics  StoreMetrics        `json:"store_metrics"`
	RecentOrders  []Order             `json:"recent_orders"`
	UserActivity  []UserActivity      `json:"user_activity"`
	SystemAlerts  []SystemAlert       `json:"system_alerts"`
	SalesTrends   []SalesPerformance  `json:"sales_trends"`
	Performance   []PerformanceMetric `json:"performance"`
}

// SystemStatus حالة النظام
type SystemStatus struct {
	Disk        DiskStatus        `json:"disk"`
	Database    DatabaseStatus    `json:"database"`
	System      SystemInfo        `json:"system"`
	Performance PerformanceMetrics `json:"performance"`
	Services    map[string]string `json:"services"`
	AIAnalysis  AIHealthAnalysis  `json:"ai_analysis"`
	Security    SecurityStatus    `json:"security"`
	LastChecked time.Time         `json:"last_checked"`
}

// ================================
// النماذج المشتركة (Shared Models)
// ================================

// PerformanceMetric مقياس الأداء
type PerformanceMetric struct {
	Value  float64 `json:"value"`
	Label  string  `json:"label"`
	Change float64 `json:"change"`
}

// Pagination الترقيم الصفحي
type Pagination struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
	Total int `json:"total"`
	Pages int `json:"pages"`
}

// Order طلب
type Order struct {
	ID          string    `json:"id"`
	User        string    `json:"user"`
	Service     string    `json:"service"`
	Amount      float64   `json:"amount"`
	Status      string    `json:"status"`
	Date        string    `json:"date"`
	Type        string    `json:"type"`
	Category    string    `json:"category"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// UserActivity نشاط المستخدم
type UserActivity struct {
	User      string    `json:"user"`
	Action    string    `json:"action"`
	Service   string    `json:"service,omitempty"`
	Time      string    `json:"time"`
	IP        string    `json:"ip"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
}

// ================================
// نماذج الخدمات (Service Models)
// ================================

// ServiceStats إحصائيات الخدمات
type ServiceStats struct {
	TotalServices     int     `json:"total_services"`
	ActiveServices    int     `json:"active_services"`
	InactiveServices  int     `json:"inactive_services"`
	SuspendedServices int     `json:"suspended_services"`
	TotalRevenue      float64 `json:"total_revenue"`
	AverageRating     float64 `json:"average_rating"`
	TotalOrders       int     `json:"total_orders"`
	PopularCategory   string  `json:"popular_category"`
}

// ServiceDetails تفاصيل الخدمة
type ServiceDetails struct {
	Service
	Seller          *User    `json:"seller,omitempty" gorm:"-"`
	Reviews         []Review `json:"reviews,omitempty" gorm:"-"`
	AverageRating   float64  `json:"average_rating" gorm:"-"`
	SimilarServices []Service `json:"similar_services,omitempty" gorm:"-"`
}

// ================================
// نماذج المحتوى (Content Models)
// ================================

// BatchContent محتوى جماعي
type BatchContent struct {
	ID        string        `json:"id"`
	Content   []Content     `json:"content"`
	Summary   *BatchSummary `json:"summary"`
	CreatedBy string        `json:"created_by"`
	CreatedAt time.Time     `json:"created_at"`
}

// ContentAnalysis تحليل المحتوى
type ContentAnalysis struct {
	ContentID        string   `json:"content_id"`
	AnalysisType     string   `json:"analysis_type"`
	SentimentScore   int      `json:"sentiment_score"`
	SEOScore         int      `json:"seo_score"`
	EngagementScore  int      `json:"engagement_score"`
	ReadabilityScore int      `json:"readability_score"`
	OverallScore     int      `json:"overall_score"`
	Recommendations  []string `json:"recommendations"`
	GeneratedAt      time.Time `json:"generated_at"`
}

// ================================
// نماذج الإشعارات (Notification Models)
// ================================

// NotificationStats إحصائيات الإشعارات
type NotificationStats struct {
	Overview   NotificationOverview       `json:"overview"`
	Behavior   UserNotificationBehavior   `json:"behavior"`
	AIInsights *AIInsights                `json:"ai_insights"`
	GeneratedAt time.Time                 `json:"generated_at"`
}

// ================================
// باقي النماذج المساعدة
// ================================

// ... (يتم إضافة جميع النماذج الأخرى هنا مع تنظيمها في أقسام)