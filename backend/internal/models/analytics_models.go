package models

import "time"

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

// ComparisonData بيانات المقارنة
type ComparisonData struct {
	PreviousPeriod *AnalyticsSummary     `json:"previous_period"`
	Changes        map[string]float64    `json:"changes"`
}

// TrendsData بيانات الاتجاهات
type TrendsData struct {
	OverallTrend string         `json:"overall_trend"`
	MetricTrends []MetricTrend  `json:"metric_trends"`
}

// MetricTrend اتجاه المقياس
type MetricTrend struct {
	Metric    string `json:"metric"`
	Direction string `json:"direction"`
	Strength  string `json:"strength"`
	Period    string `json:"period"`
}

// PerformanceAnalytics تحليلات الأداء
type PerformanceAnalytics struct {
	Timeframe string              `json:"timeframe"`
	Platform  string              `json:"platform"`
	Metrics   string              `json:"metrics"`
	Data      []AnalyticsPerformanceMetric `json:"data"` // يستخدم الاسم المختلف
	Summary   *PerformanceSummary `json:"summary"`
	GeneratedAt time.Time         `json:"generated_at"`
}

// AnalyticsPerformanceMetric مقياس أداء التحليلات (اسم مختلف لمنع التعارض)
type AnalyticsPerformanceMetric struct {
	Metric    string  `json:"metric"`
	Value     float64 `json:"value"`
	Change    float64 `json:"change"`
	Platform  string  `json:"platform"`
	Timeframe string  `json:"timeframe"`
}

// PerformanceSummary ملخص الأداء
type PerformanceSummary struct {
	AverageEngagement float64 `json:"average_engagement"`
	TotalReach        int     `json:"total_reach"`
	TotalConversions  int     `json:"total_conversions"`
	GrowthRate        float64 `json:"growth_rate"`
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

// TrendInsights رؤى الاتجاهات
type TrendInsights struct {
	PositiveTrends  []string `json:"positive_trends"`
	NegativeTrends  []string `json:"negative_trends"`
	EmergingTrends  []string `json:"emerging_trends"`
	Confidence      int      `json:"confidence"`
}

// PredictionInsights رؤى التوقعات
type PredictionInsights struct {
	NextWeek     map[string]interface{} `json:"next_week"`
	NextMonth    map[string]interface{} `json:"next_month"`
	Confidence   int                    `json:"confidence"`
	Assumptions  []string               `json:"assumptions"`
}

// RecommendationInsights رؤى التوصيات
type RecommendationInsights struct {
	HighImpact   []string `json:"high_impact"`
	MediumImpact []string `json:"medium_impact"`
	LowImpact    []string `json:"low_impact"`
}

// InsightsDataSummary ملخص بيانات الرؤى
type InsightsDataSummary struct {
	Timeframe       string `json:"timeframe"`
	Platforms       string `json:"platforms"`
	TotalDataPoints int    `json:"total_data_points"`
	AnalysisPeriod  string `json:"analysis_period"`
}

// ContentAnalytics تحليلات المحتوى
type ContentAnalytics struct {
	Performance              *ContentPerformance   `json:"performance"`
	Analysis                 *ContentAnalysis      `json:"analysis"`
	Predictions              *ContentPredictions   `json:"predictions"`
	Optimizations            []ContentOptimization `json:"optimizations"`
	ImprovementOpportunities *ContentGaps          `json:"improvement_opportunities"`
	GeneratedAt              time.Time             `json:"generated_at"`
}

// ContentPerformance أداء المحتوى
type ContentPerformance struct {
	TotalContent      int          `json:"total_content"`
	AverageEngagement float64      `json:"average_engagement"`
	TopPerforming     []ContentItem `json:"top_performing"`
}

// ContentItem عنصر المحتوى
type ContentItem struct {
	ID         string  `json:"id"`
	Title      string  `json:"title"`
	Type       string  `json:"type"`
	Engagement float64 `json:"engagement"`
	Reach      int     `json:"reach"`
	Platform   string  `json:"platform"`
}

// ContentAnalysis تحليل المحتوى
type ContentAnalysis struct {
	BestPerformingTypes []string               `json:"best_performing_types"`
	OptimalPostingTimes []string               `json:"optimal_posting_times"`
	EngagementPatterns  map[string]interface{} `json:"engagement_patterns"`
}

// ContentPredictions توقعات المحتوى
type ContentPredictions struct {
	NextWeek            map[string]interface{} `json:"next_week"`
	RecommendedContent  []string               `json:"recommended_content"`
}

// ContentOptimization تحسين المحتوى
type ContentOptimization struct {
	ContentID            string   `json:"content_id"`
	Suggestions          []string `json:"suggestions"`
	PotentialImprovement string   `json:"potential_improvement"`
}

// ContentGaps فجوات المحتوى
type ContentGaps struct {
	MissingFormats      []string `json:"missing_formats"`
	OptimalPostingTimes []string `json:"optimal_posting_times"`
	ContentThemes       []string `json:"content_themes"`
}

// AudienceAnalytics تحليلات الجمهور
type AudienceAnalytics struct {
	Demographics        *AudienceDemographics   `json:"demographics"`
	Behavior            *AudienceBehavior       `json:"behavior"`
	Analysis            *AudienceAnalysis       `json:"analysis"`
	Recommendations     *AudienceRecommendations `json:"recommendations"`
	Expansion           *AudienceExpansion      `json:"expansion"`
	Personas            []AudiencePersona       `json:"personas"`
	EngagementPatterns  *EngagementPatterns     `json:"engagement_patterns"`
	GeneratedAt         time.Time               `json:"generated_at"`
}

// AudienceDemographics ديموغرافيا الجمهور
type AudienceDemographics struct {
	AgeGroups map[string]int `json:"age_groups"`
	Genders   map[string]int `json:"genders"`
	Locations []string       `json:"locations"`
	Interests []string       `json:"interests"`
}

// AudienceBehavior سلوك الجمهور
type AudienceBehavior struct {
	ActiveTimes         []string `json:"active_times"`
	ContentPreferences  []string `json:"content_preferences"`
	EngagementLevel     string   `json:"engagement_level"`
	RetentionRate       float64  `json:"retention_rate"`
}

// AudienceAnalysis تحليل الجمهور
type AudienceAnalysis struct {
	Segments            []AudienceSegment `json:"segments"`
	GrowthOpportunities []string          `json:"growth_opportunities"`
}

// AudienceSegment شريحة الجمهور
type AudienceSegment struct {
	Name        string   `json:"name"`
	Size        int      `json:"size"`
	Engagement  float64  `json:"engagement"`
	Preferences []string `json:"preferences"`
}

// AudienceRecommendations توصيات الجمهور
type AudienceRecommendations struct {
	Targeting []string `json:"targeting"`
	Content   []string `json:"content"`
}

// AudienceExpansion توسع الجمهور
type AudienceExpansion struct {
	SimilarAudiences    []string `json:"similar_audiences"`
	GrowthOpportunities []string `json:"growth_opportunities"`
	PlatformSpecific    []string `json:"platform_specific"`
}

// AudiencePersona شخصية الجمهور
type AudiencePersona struct {
	Name         string                 `json:"name"`
	Demographics map[string]interface{} `json:"demographics"`
	Behavior     map[string]interface{} `json:"behavior"`
}

// EngagementPatterns أنماط المشاركة
type EngagementPatterns struct {
	PeakHours        []string `json:"peak_hours"`
	BestDays         []string `json:"best_days"`
	OptimalFrequency string   `json:"optimal_frequency"`
}

// CustomAnalyticsReport تقرير تحليلات مخصص
type CustomAnalyticsReport struct {
	ID              string                   `json:"id"`
	Name            string                   `json:"name"`
	Timeframe       string                   `json:"timeframe"`
	Platforms       []string                 `json:"platforms"`
	Metrics         []string                 `json:"metrics"`
	Data            []map[string]interface{} `json:"data"`
	Predictions     map[string]interface{}   `json:"predictions"`
	Recommendations []string                 `json:"recommendations"`
	Filters         map[string]interface{}   `json:"filters"`
	GeneratedAt     time.Time                `json:"generated_at"`
	GeneratedBy     string                   `json:"generated_by"`
}

// Predictions التوقعات
type Predictions struct {
	Forecasts      map[string]Forecast        `json:"forecasts"`
	Confidence     int                        `json:"confidence"`
	Assumptions    *PredictionAssumptions     `json:"assumptions"`
	Recommendations []string                  `json:"recommendations"`
	GeneratedAt    time.Time                  `json:"generated_at"`
}

// Forecast توقع
type Forecast struct {
	Value      float64 `json:"value"`
	Confidence int     `json:"confidence"`
	Timeframe  string  `json:"timeframe"`
	Trend      string  `json:"trend"`
}

// PredictionAssumptions افتراضات التوقع
type PredictionAssumptions struct {
	BasedOn           string `json:"based_on"`
	TrendContinuation string `json:"trend_continuation"`
	SeasonalFactors   string `json:"seasonal_factors"`
	MarketConditions  string `json:"market_conditions"`
}