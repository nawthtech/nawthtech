package models

import "time"

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

// NotificationStats إحصائيات الإشعارات
type NotificationStats struct {
	Overview   NotificationOverview       `json:"overview"`
	Behavior   UserNotificationBehavior   `json:"behavior"`
	AIInsights *AIInsights                `json:"ai_insights"`
	GeneratedAt time.Time                 `json:"generated_at"`
}

// NotificationOverview نظرة عامة على الإشعارات
type NotificationOverview struct {
	Total      int            `json:"total"`
	Read       int            `json:"read"`
	Unread     int            `json:"unread"`
	ByType     map[string]int `json:"by_type"`
	ByPriority map[string]int `json:"by_priority"`
}

// UserNotificationBehavior سلوك المستخدم مع الإشعارات
type UserNotificationBehavior struct {
	AverageResponseTime time.Duration       `json:"average_response_time"`
	ReadRate            float64             `json:"read_rate"`
	InteractionPatterns map[string]interface{} `json:"interaction_patterns"`
}

// NotificationInteraction تفاعل مع الإشعار
type NotificationInteraction struct {
	NotificationID string              `json:"notification_id"`
	UserID         string              `json:"user_id"`
	Action         string              `json:"action"`
	ResponseTime   time.Duration       `json:"response_time"`
	AnalyzedAt     time.Time           `json:"analyzed_at"`
	Analysis       *InteractionAnalysis `json:"analysis"`
}

// InteractionAnalysis تحليل التفاعل
type InteractionAnalysis struct {
	EngagementLevel         string   `json:"engagement_level"`
	SuggestedImprovements   []string `json:"suggested_improvements"`
}

// BulkOperationResult نتيجة عملية جماعية
type BulkOperationResult struct {
	UpdatedCount      int                      `json:"updated_count"`
	DeletedCount      int                      `json:"deleted_count"`
	Operation         string                   `json:"operation"`
	UserID            string                   `json:"user_id"`
	Type              string                   `json:"type"`
	AnalyzedAt        time.Time                `json:"analyzed_at"`
	Analysis          *BulkInteractionAnalysis `json:"analysis"`
}

// BulkInteractionAnalysis تحليل التفاعل الجماعي
type BulkInteractionAnalysis struct {
	TotalNotifications   int           `json:"total_notifications"`
	AverageResponseTime  time.Duration `json:"average_response_time"`
	Patterns             []string      `json:"patterns"`
}

// SmartNotificationResult نتيجة الإشعارات الذكية
type SmartNotificationResult struct {
	Notifications []Notification             `json:"notifications"`
	Analysis      *SmartNotificationAnalysis `json:"analysis"`
	Created       int                        `json:"created"`
	Scheduled     int                        `json:"scheduled"`
}

// SmartNotificationAnalysis تحليل الإشعارات الذكية
type SmartNotificationAnalysis struct {
	OptimalTiming        []string          `json:"optimal_timing"`
	ExpectedEngagement   float64           `json:"expected_engagement"`
	ImpactAssessment     *ImpactAssessment `json:"impact_assessment"`
}

// ImpactAssessment تقييم الأثر
type ImpactAssessment struct {
	ExpectedReach   int     `json:"expected_reach"`
	PredictedClicks int     `json:"predicted_clicks"`
	Confidence      float64 `json:"confidence"`
}

// AIRecommendations توصيات الذكاء الاصطناعي
type AIRecommendations struct {
	Recommendations []Recommendation       `json:"recommendations"`
	Summary         *RecommendationsSummary `json:"summary"`
	GeneratedAt     time.Time              `json:"generated_at"`
}

// Recommendation توصية
type Recommendation struct {
	ID              string  `json:"id"`
	Title           string  `json:"title"`
	Description     string  `json:"description"`
	Type            string  `json:"type"`
	Priority        int     `json:"priority"`
	Confidence      float64 `json:"confidence"`
	Rationale       string  `json:"rationale"`
	SuggestedAction string  `json:"suggested_action"`
}

// RecommendationsSummary ملخص التوصيات
type RecommendationsSummary struct {
	Total          int `json:"total"`
	HighPriority   int `json:"high_priority"`
	HighConfidence int `json:"high_confidence"`
}

// NotificationPreferences تفضيلات الإشعارات
type NotificationPreferences struct {
	UserID       string               `json:"user_id"`
	EmailEnabled bool                 `json:"email_enabled"`
	PushEnabled  bool                 `json:"push_enabled"`
	SMSEnabled   bool                 `json:"sms_enabled"`
	AllowedTypes []string             `json:"allowed_types"`
	QuietHours   []string             `json:"quiet_hours"`
	Language     string               `json:"language"`
	Analysis     *PreferenceAnalysis  `json:"analysis"`
	UpdatedAt    time.Time            `json:"updated_at"`
}

// PreferenceAnalysis تحليل التفضيلات
type PreferenceAnalysis struct {
	EffectivenessScore       int      `json:"effectiveness_score"`
	OptimizationSuggestions  []string `json:"optimization_suggestions"`
}

// SystemNotification إشعار النظام
type SystemNotification struct {
	ID              string             `json:"id"`
	Title           string             `json:"title"`
	Message         string             `json:"message"`
	Type            string             `json:"type"`
	Priority        string             `json:"priority"`
	TargetUsers     string             `json:"target_users"`
	ActionURL       string             `json:"action_url"`
	ExpiresAt       string             `json:"expires_at"`
	ImpactAssessment *ImpactAssessment `json:"impact_assessment"`
	CreatedBy       string             `json:"created_by"`
	CreatedAt       time.Time          `json:"created_at"`
}