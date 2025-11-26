package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nawthtech/nawthtech/backend/internal/middleware"
	"github.com/nawthtech/nawthtech/backend/internal/models"
	"github.com/nawthtech/nawthtech/backend/internal/services"
	"github.com/nawthtech/nawthtech/backend/internal/utils"
)

type WebsiteHandler struct {
	websiteService services.WebsiteService
	authService    services.AuthService
}

func NewWebsiteHandler(websiteService services.WebsiteService, authService services.AuthService) *WebsiteHandler {
	return &WebsiteHandler{
		websiteService: websiteService,
		authService:    authService,
	}
}

// GetSettings - الحصول على إعدادات الموقع
// @Summary الحصول على إعدادات الموقع
// @Description الحصول على إعدادات الموقع
// @Tags Website
// @Produce json
// @Success 200 {object} utils.Response
// @Router /api/v1/website/settings [get]
func (h *WebsiteHandler) GetSettings(c *gin.Context) {
	settings, err := h.websiteService.GetSettings(c)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "لم يتم العثور على إعدادات الموقع", "SETTINGS_NOT_FOUND")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "تم جلب إعدادات الموقع بنجاح", settings)
}

// UpdateSettingsRequest - طلب تحديث إعدادات الموقع
type UpdateSettingsRequest struct {
	SiteName        string                 `json:"siteName"`
	SiteDescription string                 `json:"siteDescription"`
	SocialMedia     map[string]interface{} `json:"socialMedia"`
	SEO             map[string]interface{} `json:"seo"`
	Content         map[string]interface{} `json:"content"`
	Performance     map[string]interface{} `json:"performance"`
}

// UpdateSettings - تحديث إعدادات الموقع
// @Summary تحديث إعدادات الموقع
// @Description تحديث إعدادات الموقع (للمشرفين فقط)
// @Tags Website
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param input body UpdateSettingsRequest true "بيانات التحديث"
// @Success 200 {object} utils.Response
// @Router /api/v1/website/settings [put]
func (h *WebsiteHandler) UpdateSettings(c *gin.Context) {
	var req UpdateSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "بيانات غير صالحة", "INVALID_INPUT")
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "غير مصرح", "UNAUTHORIZED")
		return
	}

	updatedSettings, err := h.websiteService.UpdateSettings(c, services.UpdateSettingsParams{
		SiteName:        req.SiteName,
		SiteDescription: req.SiteDescription,
		SocialMedia:     req.SocialMedia,
		SEO:             req.SEO,
		Content:         req.Content,
		Performance:     req.Performance,
		UserID:          userID.(string),
	})

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في تحديث إعدادات الموقع", "SETTINGS_UPDATE_FAILED")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "تم تحديث إعدادات الموقع بنجاح", updatedSettings)
}

// GetAIOptimizedSettings - الحصول على إعدادات الموقع المحسنة بالذكاء الاصطناعي
// @Summary الحصول على إعدادات الموقع المحسنة بالذكاء الاصطناعي
// @Description الحصول على إعدادات الموقع المحسنة بالذكاء الاصطناعي (للمشرفين فقط)
// @Tags Website
// @Security BearerAuth
// @Produce json
// @Success 200 {object} utils.Response
// @Router /api/v1/website/settings/ai-optimized [get]
func (h *WebsiteHandler) GetAIOptimizedSettings(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "غير مصرح", "UNAUTHORIZED")
		return
	}

	if !middleware.CheckRateLimit(c, "website_ai_optimized", 20, 10*time.Minute) {
		utils.ErrorResponse(c, http.StatusTooManyRequests, "تم تجاوز الحد المسموح", "RATE_LIMIT_EXCEEDED")
		return
	}

	optimizedSettings, err := h.websiteService.GetAIOptimizedSettings(c, userID.(string))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في تحسين الإعدادات باستخدام الذكاء الاصطناعي", "AI_OPTIMIZATION_FAILED")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "تم تحسين الإعدادات باستخدام الذكاء الاصطناعي", optimizedSettings)
}

// GenerateContentStrategyRequest - طلب إنشاء استراتيجية محتوى
type GenerateContentStrategyRequest struct {
	Profile map[string]interface{} `json:"profile" binding:"required"`
	Goals   map[string]interface{} `json:"goals" binding:"required"`
	Options map[string]interface{} `json:"options"`
}

// GenerateContentStrategy - إنشاء استراتيجية محتوى باستخدام الذكاء الاصطناعي
// @Summary إنشاء استراتيجية محتوى باستخدام الذكاء الاصطناعي
// @Description إنشاء استراتيجية محتوى باستخدام الذكاء الاصطناعي (للمشرفين فقط)
// @Tags Website
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param input body GenerateContentStrategyRequest true "بيانات الاستراتيجية"
// @Success 200 {object} utils.Response
// @Router /api/v1/website/strategy/generate [post]
func (h *WebsiteHandler) GenerateContentStrategy(c *gin.Context) {
	var req GenerateContentStrategyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "بيانات غير صالحة", "INVALID_INPUT")
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "غير مصرح", "UNAUTHORIZED")
		return
	}

	if !middleware.CheckRateLimit(c, "website_strategy_generate", 10, 5*time.Minute) {
		utils.ErrorResponse(c, http.StatusTooManyRequests, "تم تجاوز الحد المسموح", "RATE_LIMIT_EXCEEDED")
		return
	}

	strategy, err := h.websiteService.GenerateContentStrategy(c, services.GenerateContentStrategyParams{
		Profile: req.Profile,
		Goals:   req.Goals,
		Options: req.Options,
		UserID:  userID.(string),
	})

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في إنشاء استراتيجية المحتوى", "STRATEGY_GENERATION_FAILED")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "تم إنشاء استراتيجية المحتوى بنجاح", strategy)
}

// GetAIAnalyticsInsights - تحليل أداء الموقع باستخدام الذكاء الاصطناعي
// @Summary تحليل أداء الموقع باستخدام الذكاء الاصطناعي
// @Description تحليل أداء الموقع باستخدام الذكاء الاصطناعي (للمشرفين فقط)
// @Tags Website
// @Security BearerAuth
// @Produce json
// @Param timeframe query string false "الفترة الزمنية" default(30d)
// @Param platforms query string false "المنصات" default(instagram)
// @Success 200 {object} utils.Response
// @Router /api/v1/website/analytics/ai-insights [get]
func (h *WebsiteHandler) GetAIAnalyticsInsights(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "غير مصرح", "UNAUTHORIZED")
		return
	}

	if !middleware.CheckRateLimit(c, "website_ai_insights", 15, 5*time.Minute) {
		utils.ErrorResponse(c, http.StatusTooManyRequests, "تم تجاوز الحد المسموح", "RATE_LIMIT_EXCEEDED")
		return
	}

	timeframe := c.DefaultQuery("timeframe", "30d")
	platforms := c.DefaultQuery("platforms", "instagram")

	insights, err := h.websiteService.GetAIAnalyticsInsights(c, services.GetAIAnalyticsInsightsParams{
		Timeframe: timeframe,
		Platforms: platforms,
		UserID:    userID.(string),
	})

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في تحليل البيانات باستخدام الذكاء الاصطناعي", "AI_ANALYTICS_FAILED")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "تم تحليل البيانات باستخدام الذكاء الاصطناعي", insights)
}

// GenerateContentRequest - طلب إنشاء محتوى محسن
type GenerateContentRequest struct {
	Topic     string   `json:"topic" binding:"required"`
	Platform  string   `json:"platform" binding:"required"`
	Tone      string   `json:"tone"`
	Keywords  []string `json:"keywords"`
	Language  string   `json:"language"`
}

// GenerateContent - إنشاء محتوى محسن لتحسين محركات البحث باستخدام الذكاء الاصطناعي
// @Summary إنشاء محتوى محسن لتحسين محركات البحث باستخدام الذكاء الاصطناعي
// @Description إنشاء محتوى محسن لتحسين محركات البحث باستخدام الذكاء الاصطناعي (للمشرفين فقط)
// @Tags Website
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param input body GenerateContentRequest true "بيانات المحتوى"
// @Success 200 {object} utils.Response
// @Router /api/v1/website/content/generate [post]
func (h *WebsiteHandler) GenerateContent(c *gin.Context) {
	var req GenerateContentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "بيانات غير صالحة", "INVALID_INPUT")
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "غير مصرح", "UNAUTHORIZED")
		return
	}

	if !middleware.CheckRateLimit(c, "website_content_generate", 25, 10*time.Minute) {
		utils.ErrorResponse(c, http.StatusTooManyRequests, "تم تجاوز الحد المسموح", "RATE_LIMIT_EXCEEDED")
		return
	}

	content, err := h.websiteService.GenerateContent(c, services.GenerateContentParams{
		Topic:    req.Topic,
		Platform: req.Platform,
		Tone:     req.Tone,
		Keywords: req.Keywords,
		Language: req.Language,
		UserID:   userID.(string),
	})

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في إنشاء المحتوى", "CONTENT_GENERATION_FAILED")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "تم إنشاء المحتوى بنجاح", content)
}

// AIOptimizeSettingsRequest - طلب تحسين الإعدادات بالذكاء الاصطناعي
type AIOptimizeSettingsRequest struct {
	Sections []string `json:"sections"`
}

// AIOptimizeSettings - تحديث إعدادات الموقع مع توصيات الذكاء الاصطناعي
// @Summary تحديث إعدادات الموقع مع توصيات الذكاء الاصطناعي
// @Description تحديث إعدادات الموقع مع توصيات الذكاء الاصطناعي (للمشرفين فقط)
// @Tags Website
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param input body AIOptimizeSettingsRequest true "أقسام التحسين"
// @Success 200 {object} utils.Response
// @Router /api/v1/website/settings/ai-optimize [patch]
func (h *WebsiteHandler) AIOptimizeSettings(c *gin.Context) {
	var req AIOptimizeSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "بيانات غير صالحة", "INVALID_INPUT")
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "غير مصرح", "UNAUTHORIZED")
		return
	}

	if !middleware.CheckRateLimit(c, "website_ai_optimize", 5, 15*time.Minute) {
		utils.ErrorResponse(c, http.StatusTooManyRequests, "تم تجاوز الحد المسموح", "RATE_LIMIT_EXCEEDED")
		return
	}

	optimizedSettings, err := h.websiteService.AIOptimizeSettings(c, services.AIOptimizeSettingsParams{
		Sections: req.Sections,
		UserID:   userID.(string),
	})

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في تحسين الإعدادات باستخدام الذكاء الاصطناعي", "AI_OPTIMIZATION_FAILED")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "تم تحسين الإعدادات باستخدام الذكاء الاصطناعي", optimizedSettings)
}

// GetPerformancePredictions - الحصول على توقعات الأداء باستخدام الذكاء الاصطناعي
// @Summary الحصول على توقعات الأداء باستخدام الذكاء الاصطناعي
// @Description الحصول على توقعات الأداء باستخدام الذكاء الاصطناعي (للمشرفين فقط)
// @Tags Website
// @Security BearerAuth
// @Produce json
// @Param timeframe query string false "الفترة الزمنية" default(7d)
// @Param metrics query string false "المقاييس" default(engagement,growth,reach)
// @Success 200 {object} utils.Response
// @Router /api/v1/website/predictions/performance [get]
func (h *WebsiteHandler) GetPerformancePredictions(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "غير مصرح", "UNAUTHORIZED")
		return
	}

	if !middleware.CheckRateLimit(c, "website_predictions", 30, 10*time.Minute) {
		utils.ErrorResponse(c, http.StatusTooManyRequests, "تم تجاوز الحد المسموح", "RATE_LIMIT_EXCEEDED")
		return
	}

	timeframe := c.DefaultQuery("timeframe", "7d")
	metrics := c.DefaultQuery("metrics", "engagement,growth,reach")

	predictions, err := h.websiteService.GetPerformancePredictions(c, services.GetPerformancePredictionsParams{
		Timeframe: timeframe,
		Metrics:   metrics,
		UserID:    userID.(string),
	})

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في توليد توقعات الأداء", "PREDICTIONS_FAILED")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "تم توليد توقعات الأداء بنجاح", predictions)
}

// GetAudienceInsights - تحليل الجمهور باستخدام الذكاء الاصطناعي
// @Summary تحليل الجمهور باستخدام الذكاء الاصطناعي
// @Description تحليل الجمهور باستخدام الذكاء الاصطناعي (للمشرفين فقط)
// @Tags Website
// @Security BearerAuth
// @Produce json
// @Param platform query string false "المنصة" default(instagram)
// @Param timeframe query string false "الفترة الزمنية" default(30d)
// @Success 200 {object} utils.Response
// @Router /api/v1/website/audience/insights [get]
func (h *WebsiteHandler) GetAudienceInsights(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "غير مصرح", "UNAUTHORIZED")
		return
	}

	if !middleware.CheckRateLimit(c, "website_audience_insights", 20, 10*time.Minute) {
		utils.ErrorResponse(c, http.StatusTooManyRequests, "تم تجاوز الحد المسموح", "RATE_LIMIT_EXCEEDED")
		return
	}

	platform := c.DefaultQuery("platform", "instagram")
	timeframe := c.DefaultQuery("timeframe", "30d")

	insights, err := h.websiteService.GetAudienceInsights(c, services.GetAudienceInsightsParams{
		Platform:  platform,
		Timeframe: timeframe,
		UserID:    userID.(string),
	})

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في تحليل الجمهور باستخدام الذكاء الاصطناعي", "AUDIENCE_INSIGHTS_FAILED")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "تم تحليل الجمهور باستخدام الذكاء الاصطناعي", insights)
}