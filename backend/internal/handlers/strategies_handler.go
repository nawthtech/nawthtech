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

type StrategiesHandler struct {
	strategiesService services.StrategiesService
	authService       services.AuthService
}

func NewStrategiesHandler(strategiesService services.StrategiesService, authService services.AuthService) *StrategiesHandler {
	return &StrategiesHandler{
		strategiesService: strategiesService,
		authService:       authService,
	}
}

// CreateStrategyRequest - طلب إنشاء استراتيجية جديدة
type CreateStrategyRequest struct {
	Name           string                 `json:"name" binding:"required"`
	Description    string                 `json:"description" binding:"required"`
	Goals          []string               `json:"goals" binding:"required"`
	Platforms      []string               `json:"platforms" binding:"required"`
	TargetAudience map[string]interface{} `json:"targetAudience" binding:"required"`
	Budget         float64                `json:"budget"`
	Timeline       map[string]interface{} `json:"timeline"`
}

// CreateStrategy - إنشاء استراتيجية جديدة باستخدام الذكاء الاصطناعي
// @Summary إنشاء استراتيجية جديدة باستخدام الذكاء الاصطناعي
// @Description إنشاء استراتيجية جديدة باستخدام الذكاء الاصطناعي (للمشرفين فقط)
// @Tags Strategies
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param input body CreateStrategyRequest true "بيانات الاستراتيجية"
// @Success 200 {object} utils.Response
// @Router /api/v1/strategies [post]
func (h *StrategiesHandler) CreateStrategy(c *gin.Context) {
	var req CreateStrategyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "بيانات غير صالحة", "INVALID_INPUT")
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "غير مصرح", "UNAUTHORIZED")
		return
	}

	if !middleware.CheckRateLimit(c, "strategies_create", 10, 10*time.Minute) {
		utils.ErrorResponse(c, http.StatusTooManyRequests, "تم تجاوز الحد المسموح", "RATE_LIMIT_EXCEEDED")
		return
	}

	strategy, err := h.strategiesService.CreateStrategy(c, services.CreateStrategyParams{
		Name:           req.Name,
		Description:    req.Description,
		Goals:          req.Goals,
		Platforms:      req.Platforms,
		TargetAudience: req.TargetAudience,
		Budget:         req.Budget,
		Timeline:       req.Timeline,
		UserID:         userID.(string),
	})

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في إنشاء الاستراتيجية", "STRATEGY_CREATION_FAILED")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "تم إنشاء الاستراتيجية بنجاح", strategy)
}

// GetStrategies - الحصول على جميع الاستراتيجيات
// @Summary الحصول على جميع الاستراتيجيات
// @Description الحصول على جميع الاستراتيجيات (للمشرفين فقط)
// @Tags Strategies
// @Security BearerAuth
// @Produce json
// @Param page query int false "الصفحة" default(1)
// @Param limit query int false "الحد" default(20)
// @Param status query string false "الحالة" default(active)
// @Success 200 {object} utils.Response
// @Router /api/v1/strategies [get]
func (h *StrategiesHandler) GetStrategies(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "غير مصرح", "UNAUTHORIZED")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	status := c.DefaultQuery("status", "active")

	strategies, pagination, err := h.strategiesService.GetStrategies(c, services.GetStrategiesParams{
		Page:   page,
		Limit:  limit,
		Status: status,
		UserID: userID.(string),
	})

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في جلب الاستراتيجيات", "STRATEGIES_FETCH_FAILED")
		return
	}

	response := map[string]interface{}{
		"strategies": strategies,
		"pagination": pagination,
	}

	utils.SuccessResponse(c, http.StatusOK, "تم جلب الاستراتيجيات بنجاح", response)
}

// GetStrategyByID - الحصول على استراتيجية محددة
// @Summary الحصول على استراتيجية محددة
// @Description الحصول على استراتيجية محددة (للمشرفين فقط)
// @Tags Strategies
// @Security BearerAuth
// @Produce json
// @Param id path string true "معرف الاستراتيجية"
// @Success 200 {object} utils.Response
// @Router /api/v1/strategies/{id} [get]
func (h *StrategiesHandler) GetStrategyByID(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "غير مصرح", "UNAUTHORIZED")
		return
	}

	strategyID := c.Param("id")

	strategy, err := h.strategiesService.GetStrategyByID(c, strategyID, userID.(string))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "الاستراتيجية غير موجودة", "STRATEGY_NOT_FOUND")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "تم جلب الاستراتيجية بنجاح", strategy)
}

// UpdateStrategyRequest - طلب تحديث استراتيجية
type UpdateStrategyRequest struct {
	Name           string                 `json:"name"`
	Description    string                 `json:"description"`
	Goals          []string               `json:"goals"`
	Platforms      []string               `json:"platforms"`
	TargetAudience map[string]interface{} `json:"targetAudience"`
	Budget         float64                `json:"budget"`
	Timeline       map[string]interface{} `json:"timeline"`
	Status         string                 `json:"status"`
}

// UpdateStrategy - تحديث استراتيجية محددة
// @Summary تحديث استراتيجية محددة
// @Description تحديث استراتيجية محددة (للمشرفين فقط)
// @Tags Strategies
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "معرف الاستراتيجية"
// @Param input body UpdateStrategyRequest true "بيانات التحديث"
// @Success 200 {object} utils.Response
// @Router /api/v1/strategies/{id} [put]
func (h *StrategiesHandler) UpdateStrategy(c *gin.Context) {
	var req UpdateStrategyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "بيانات غير صالحة", "INVALID_INPUT")
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "غير مصرح", "UNAUTHORIZED")
		return
	}

	strategyID := c.Param("id")

	updatedStrategy, err := h.strategiesService.UpdateStrategy(c, services.UpdateStrategyParams{
		StrategyID:     strategyID,
		Name:           req.Name,
		Description:    req.Description,
		Goals:          req.Goals,
		Platforms:      req.Platforms,
		TargetAudience: req.TargetAudience,
		Budget:         req.Budget,
		Timeline:       req.Timeline,
		Status:         req.Status,
		UserID:         userID.(string),
	})

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في تحديث الاستراتيجية", "STRATEGY_UPDATE_FAILED")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "تم تحديث الاستراتيجية بنجاح", updatedStrategy)
}

// DeleteStrategy - حذف استراتيجية محددة
// @Summary حذف استراتيجية محددة
// @Description حذف استراتيجية محددة (للمشرفين فقط)
// @Tags Strategies
// @Security BearerAuth
// @Produce json
// @Param id path string true "معرف الاستراتيجية"
// @Success 200 {object} utils.Response
// @Router /api/v1/strategies/{id} [delete]
func (h *StrategiesHandler) DeleteStrategy(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "غير مصرح", "UNAUTHORIZED")
		return
	}

	strategyID := c.Param("id")

	err := h.strategiesService.DeleteStrategy(c, strategyID, userID.(string))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في حذف الاستراتيجية", "STRATEGY_DELETE_FAILED")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "تم حذف الاستراتيجية بنجاح", nil)
}

// AnalyzeStrategyRequest - طلب تحليل استراتيجية
type AnalyzeStrategyRequest struct {
	AnalysisType string `json:"analysisType"`
}

// AnalyzeStrategy - تحليل استراتيجية باستخدام الذكاء الاصطناعي
// @Summary تحليل استراتيجية باستخدام الذكاء الاصطناعي
// @Description تحليل استراتيجية باستخدام الذكاء الاصطناعي (للمشرفين فقط)
// @Tags Strategies
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "معرف الاستراتيجية"
// @Param input body AnalyzeStrategyRequest true "بيانات التحليل"
// @Success 200 {object} utils.Response
// @Router /api/v1/strategies/{id}/analyze [post]
func (h *StrategiesHandler) AnalyzeStrategy(c *gin.Context) {
	var req AnalyzeStrategyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "بيانات غير صالحة", "INVALID_INPUT")
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "غير مصرح", "UNAUTHORIZED")
		return
	}

	if !middleware.CheckRateLimit(c, "strategies_analyze", 15, 5*time.Minute) {
		utils.ErrorResponse(c, http.StatusTooManyRequests, "تم تجاوز الحد المسموح", "RATE_LIMIT_EXCEEDED")
		return
	}

	strategyID := c.Param("id")

	analysis, err := h.strategiesService.AnalyzeStrategy(c, strategyID, req.AnalysisType, userID.(string))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في تحليل الاستراتيجية", "STRATEGY_ANALYSIS_FAILED")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "تم تحليل الاستراتيجية بنجاح", analysis)
}

// GetStrategyPerformance - الحصول على أداء الاستراتيجية
// @Summary الحصول على أداء الاستراتيجية
// @Description الحصول على أداء الاستراتيجية (للمشرفين فقط)
// @Tags Strategies
// @Security BearerAuth
// @Produce json
// @Param id path string true "معرف الاستراتيجية"
// @Param timeframe query string false "الفترة الزمنية" default(30d)
// @Success 200 {object} utils.Response
// @Router /api/v1/strategies/{id}/performance [get]
func (h *StrategiesHandler) GetStrategyPerformance(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "غير مصرح", "UNAUTHORIZED")
		return
	}

	strategyID := c.Param("id")
	timeframe := c.DefaultQuery("timeframe", "30d")

	performance, err := h.strategiesService.GetStrategyPerformance(c, strategyID, timeframe, userID.(string))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في جلب أداء الاستراتيجية", "PERFORMANCE_FETCH_FAILED")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "تم جلب أداء الاستراتيجية بنجاح", performance)
}

// GetStrategyRecommendationsRequest - طلب الحصول على توصيات استراتيجية
type GetStrategyRecommendationsRequest struct {
	Goals          []string               `json:"goals" binding:"required"`
	Constraints    map[string]interface{} `json:"constraints"`
	Preferences    map[string]interface{} `json:"preferences"`
	HistoricalData map[string]interface{} `json:"historicalData"`
}

// GetStrategyRecommendations - الحصول على توصيات استراتيجية ذكية
// @Summary الحصول على توصيات استراتيجية ذكية
// @Description الحصول على توصيات استراتيجية ذكية (للمشرفين فقط)
// @Tags Strategies
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param input body GetStrategyRecommendationsRequest true "بيانات التوصيات"
// @Success 200 {object} utils.Response
// @Router /api/v1/strategies/recommend [post]
func (h *StrategiesHandler) GetStrategyRecommendations(c *gin.Context) {
	var req GetStrategyRecommendationsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "بيانات غير صالحة", "INVALID_INPUT")
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "غير مصرح", "UNAUTHORIZED")
		return
	}

	if !middleware.CheckRateLimit(c, "strategies_recommend", 20, 10*time.Minute) {
		utils.ErrorResponse(c, http.StatusTooManyRequests, "تم تجاوز الحد المسموح", "RATE_LIMIT_EXCEEDED")
		return
	}

	recommendations, err := h.strategiesService.GetStrategyRecommendations(c, services.GetStrategyRecommendationsParams{
		Goals:          req.Goals,
		Constraints:    req.Constraints,
		Preferences:    req.Preferences,
		HistoricalData: req.HistoricalData,
		UserID:         userID.(string),
	})

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في توليد التوصيات", "RECOMMENDATIONS_FAILED")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "تم توليد التوصيات بنجاح", recommendations)
}