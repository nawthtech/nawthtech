package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nawthtech/nawthtech/backend/internal/middleware"
	"github.com/nawthtech/nawthtech/backend/internal/models"
	"github.com/nawthtech/nawthtech/backend/internal/services"
	"github.com/nawthtech/nawthtech/backend/internal/utils"
)

type AnalyticsHandler struct {
	analyticsService services.AnalyticsService
	authService      services.AuthService
}

func NewAnalyticsHandler(analyticsService services.AnalyticsService, authService services.AuthService) *AnalyticsHandler {
	return &AnalyticsHandler{
		analyticsService: analyticsService,
		authService:      authService,
	}
}

// GetOverview - الحصول على نظرة عامة على تحليلات الموقع
// @Summary الحصول على نظرة عامة على تحليلات الموقع
// @Description الحصول على نظرة عامة على تحليلات الموقع (للمشرفين فقط)
// @Tags Analytics
// @Security BearerAuth
// @Produce json
// @Param timeframe query string false "الفترة الزمنية" default(30d)
// @Param compareTo query string false "المقارنة مع" default(previous)
// @Success 200 {object} utils.Response
// @Router /api/v1/analytics/overview [get]
func (h *AnalyticsHandler) GetOverview(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "غير مصرح", "UNAUTHORIZED")
		return
	}

	timeframe := c.DefaultQuery("timeframe", "30d")
	compareTo := c.DefaultQuery("compareTo", "previous")

	overview, err := h.analyticsService.GetOverview(c, services.GetOverviewParams{
		Timeframe: timeframe,
		CompareTo: compareTo,
		UserID:    userID.(string),
	})

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في جلب نظرة عامة على التحليلات", "OVERVIEW_FETCH_FAILED")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "تم جلب نظرة عامة على التحليلات بنجاح", overview)
}

// GetPerformance - الحصول على تحليلات الأداء التفصيلية
// @Summary الحصول على تحليلات الأداء التفصيلية
// @Description الحصول على تحليلات الأداء التفصيلية (للمشرفين فقط)
// @Tags Analytics
// @Security BearerAuth
// @Produce json
// @Param timeframe query string false "الفترة الزمنية" default(30d)
// @Param metrics query string false "المقاييس" default(engagement,reach,conversion)
// @Param platform query string false "المنصة" default(all)
// @Success 200 {object} utils.Response
// @Router /api/v1/analytics/performance [get]
func (h *AnalyticsHandler) GetPerformance(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "غير مصرح", "UNAUTHORIZED")
		return
	}

	timeframe := c.DefaultQuery("timeframe", "30d")
	metrics := c.Query("metrics")
	platform := c.DefaultQuery("platform", "all")

	performance, err := h.analyticsService.GetPerformance(c, services.GetPerformanceParams{
		Timeframe: timeframe,
		Metrics:   metrics,
		Platform:  platform,
		UserID:    userID.(string),
	})

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في جلب تحليلات الأداء", "PERFORMANCE_FETCH_FAILED")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "تم جلب تحليلات الأداء بنجاح", performance)
}

// GetAIInsights - تحليل البيانات باستخدام الذكاء الاصطناعي للحصول على رؤى متقدمة
// @Summary تحليل البيانات باستخدام الذكاء الاصطناعي للحصول على رؤى متقدمة
// @Description تحليل البيانات باستخدام الذكاء الاصطناعي للحصول على رؤى متقدمة (للمشرفين فقط)
// @Tags Analytics
// @Security BearerAuth
// @Produce json
// @Param timeframe query string false "الفترة الزمنية" default(30d)
// @Param platforms query string false "المنصات" default(instagram,twitter)
// @Param insightTypes query string false "أنواع الرؤى" default(trends,predictions,recommendations)
// @Success 200 {object} utils.Response
// @Router /api/v1/analytics/ai-insights [get]
func (h *AnalyticsHandler) GetAIInsights(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "غير مصرح", "UNAUTHORIZED")
		return
	}

	if !middleware.CheckRateLimit(c, "analytics_ai_insights", 15, 5*time.Minute) {
		utils.ErrorResponse(c, http.StatusTooManyRequests, "تم تجاوز الحد المسموح", "RATE_LIMIT_EXCEEDED")
		return
	}

	timeframe := c.DefaultQuery("timeframe", "30d")
	platforms := c.DefaultQuery("platforms", "instagram,twitter")
	insightTypes := c.DefaultQuery("insightTypes", "trends,predictions,recommendations")

	insights, err := h.analyticsService.GetAIInsights(c, services.GetAIInsightsParams{
		Timeframe:    timeframe,
		Platforms:    platforms,
		InsightTypes: insightTypes,
		UserID:       userID.(string),
	})

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في توليد الرؤى باستخدام الذكاء الاصطناعي", "AI_INSIGHTS_FAILED")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "تم توليد الرؤى باستخدام الذكاء الاصطناعي بنجاح", insights)
}

// GetContentAnalytics - تحليل أداء المحتوى باستخدام الذكاء الاصطناعي
// @Summary تحليل أداء المحتوى باستخدام الذكاء الاصطناعي
// @Description تحليل أداء المحتوى باستخدام الذكاء الاصطناعي (للمشرفين فقط)
// @Tags Analytics
// @Security BearerAuth
// @Produce json
// @Param timeframe query string false "الفترة الزمنية" default(30d)
// @Param contentType query string false "نوع المحتوى" default(all)
// @Param platform query string false "المنصة" default(all)
// @Param sortBy query string false "ترتيب حسب" default(engagement)
// @Success 200 {object} utils.Response
// @Router /api/v1/analytics/content [get]
func (h *AnalyticsHandler) GetContentAnalytics(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "غير مصرح", "UNAUTHORIZED")
		return
	}

	if !middleware.CheckRateLimit(c, "analytics_content", 20, 10*time.Minute) {
		utils.ErrorResponse(c, http.StatusTooManyRequests, "تم تجاوز الحد المسموح", "RATE_LIMIT_EXCEEDED")
		return
	}

	timeframe := c.DefaultQuery("timeframe", "30d")
	contentType := c.DefaultQuery("contentType", "all")
	platform := c.DefaultQuery("platform", "all")
	sortBy := c.DefaultQuery("sortBy", "engagement")

	contentAnalytics, err := h.analyticsService.GetContentAnalytics(c, services.GetContentAnalyticsParams{
		Timeframe:   timeframe,
		ContentType: contentType,
		Platform:    platform,
		SortBy:      sortBy,
		UserID:      userID.(string),
	})

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في تحليل أداء المحتوى", "CONTENT_ANALYTICS_FAILED")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "تم تحليل أداء المحتوى بنجاح", contentAnalytics)
}

// GetAudienceAnalytics - تحليل الجمهور باستخدام الذكاء الاصطناعي
// @Summary تحليل الجمهور باستخدام الذكاء الاصطناعي
// @Description تحليل الجمهور باستخدام الذكاء الاصطناعي (للمشرفين فقط)
// @Tags Analytics
// @Security BearerAuth
// @Produce json
// @Param timeframe query string false "الفترة الزمنية" default(30d)
// @Param platform query string true "المنصة"
// @Param segment query string false "الشريحة" default(all)
// @Success 200 {object} utils.Response
// @Router /api/v1/analytics/audience [get]
func (h *AnalyticsHandler) GetAudienceAnalytics(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "غير مصرح", "UNAUTHORIZED")
		return
	}

	if !middleware.CheckRateLimit(c, "analytics_audience", 15, 10*time.Minute) {
		utils.ErrorResponse(c, http.StatusTooManyRequests, "تم تجاوز الحد المسموح", "RATE_LIMIT_EXCEEDED")
		return
	}

	timeframe := c.DefaultQuery("timeframe", "30d")
	platform := c.Query("platform")
	segment := c.DefaultQuery("segment", "all")

	if platform == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "المنصة مطلوبة", "PLATFORM_REQUIRED")
		return
	}

	audienceAnalytics, err := h.analyticsService.GetAudienceAnalytics(c, services.GetAudienceAnalyticsParams{
		Timeframe: timeframe,
		Platform:  platform,
		Segment:   segment,
		UserID:    userID.(string),
	})

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في تحليل بيانات الجمهور", "AUDIENCE_ANALYTICS_FAILED")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "تم تحليل بيانات الجمهور بنجاح", audienceAnalytics)
}

// GenerateCustomReportRequest - طلب إنشاء تقرير مخصص
type GenerateCustomReportRequest struct {
	Name                  string                 `json:"name" binding:"required"`
	Metrics               []string               `json:"metrics" binding:"required"`
	Timeframe             string                 `json:"timeframe" binding:"required"`
	Platforms             []string               `json:"platforms"`
	Filters               map[string]interface{} `json:"filters"`
	IncludePredictions    bool                   `json:"includePredictions"`
	IncludeRecommendations bool                  `json:"includeRecommendations"`
}

// GenerateCustomReport - إنشاء تقرير مخصص باستخدام الذكاء الاصطناعي
// @Summary إنشاء تقرير مخصص باستخدام الذكاء الاصطناعي
// @Description إنشاء تقرير مخصص باستخدام الذكاء الاصطناعي (للمشرفين فقط)
// @Tags Analytics
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param input body GenerateCustomReportRequest true "بيانات التقرير المخصص"
// @Success 200 {object} utils.Response
// @Router /api/v1/analytics/custom-report [post]
func (h *AnalyticsHandler) GenerateCustomReport(c *gin.Context) {
	var req GenerateCustomReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "بيانات غير صالحة", "INVALID_INPUT")
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "غير مصرح", "UNAUTHORIZED")
		return
	}

	if !middleware.CheckRateLimit(c, "analytics_custom_report", 10, 15*time.Minute) {
		utils.ErrorResponse(c, http.StatusTooManyRequests, "تم تجاوز الحد المسموح", "RATE_LIMIT_EXCEEDED")
		return
	}

	customReport, err := h.analyticsService.GenerateCustomReport(c, services.GenerateCustomReportParams{
		Name:                  req.Name,
		Metrics:               req.Metrics,
		Timeframe:             req.Timeframe,
		Platforms:             req.Platforms,
		Filters:               req.Filters,
		IncludePredictions:    req.IncludePredictions,
		IncludeRecommendations: req.IncludeRecommendations,
		UserID:                userID.(string),
	})

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في إنشاء التقرير المخصص", "CUSTOM_REPORT_FAILED")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "تم إنشاء التقرير المخصص بنجاح", customReport)
}

// GetCustomReports - الحصول على التقارير المخصصة المحفوظة
// @Summary الحصول على التقارير المخصصة المحفوظة
// @Description الحصول على التقارير المخصصة المحفوظة (للمشرفين فقط)
// @Tags Analytics
// @Security BearerAuth
// @Produce json
// @Param page query int false "الصفحة" default(1)
// @Param limit query int false "الحد" default(20)
// @Success 200 {object} utils.Response
// @Router /api/v1/analytics/custom-reports [get]
func (h *AnalyticsHandler) GetCustomReports(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "غير مصرح", "UNAUTHORIZED")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	reports, pagination, err := h.analyticsService.GetCustomReports(c, services.GetCustomReportsParams{
		UserID: userID.(string),
		Page:   page,
		Limit:  limit,
	})

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في جلب التقارير المخصصة", "CUSTOM_REPORTS_FETCH_FAILED")
		return
	}

	response := map[string]interface{}{
		"reports":    reports,
		"pagination": pagination,
	}

	utils.SuccessResponse(c, http.StatusOK, "تم جلب التقارير المخصصة بنجاح", response)
}

// GetPredictions - الحصول على توقعات الأداء المستقبلية
// @Summary الحصول على توقعات الأداء المستقبلية
// @Description الحصول على توقعات الأداء المستقبلية (للمشرفين فقط)
// @Tags Analytics
// @Security BearerAuth
// @Produce json
// @Param timeframe query string false "الفترة الزمنية" default(30d)
// @Param forecastPeriod query string false "فترة التوقع" default(7d)
// @Param metrics query string false "المقاييس" default(engagement,growth,reach)
// @Success 200 {object} utils.Response
// @Router /api/v1/analytics/predictions [get]
func (h *AnalyticsHandler) GetPredictions(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "غير مصرح", "UNAUTHORIZED")
		return
	}

	if !middleware.CheckRateLimit(c, "analytics_predictions", 30, 10*time.Minute) {
		utils.ErrorResponse(c, http.StatusTooManyRequests, "تم تجاوز الحد المسموح", "RATE_LIMIT_EXCEEDED")
		return
	}

	timeframe := c.DefaultQuery("timeframe", "30d")
	forecastPeriod := c.DefaultQuery("forecastPeriod", "7d")
	metrics := c.DefaultQuery("metrics", "engagement,growth,reach")

	predictions, err := h.analyticsService.GetPredictions(c, services.GetPredictionsParams{
		Timeframe:      timeframe,
		ForecastPeriod: forecastPeriod,
		Metrics:        metrics,
		UserID:         userID.(string),
	})

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في توليد توقعات الأداء", "PREDICTIONS_FAILED")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "تم توليد توقعات الأداء بنجاح", predictions)
}