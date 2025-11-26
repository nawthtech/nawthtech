package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"backend-app/internal/models"
	"backend-app/internal/services"
)

type AdminHandler struct {
	adminService *services.AdminService
}

func NewAdminHandler(adminService *services.AdminService) *AdminHandler {
	return &AdminHandler{
		adminService: adminService,
	}
}

// InitiateSystemUpdate يبدأ عملية تحديث النظام
func (h *AdminHandler) InitiateSystemUpdate(w http.ResponseWriter, r *http.Request) {
	var request models.SystemUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		respondWithError(w, http.StatusBadRequest, "طلب غير صالح")
		return
	}

	if request.UpdateType == "" || request.Version == "" {
		respondWithError(w, http.StatusBadRequest, "نوع التحديث والإصدار مطلوبان")
		return
	}

	// في الواقع، يجب الحصول على المستخدم من التوكن
	initiatedBy := "admin-user"

	updateResult, err := h.adminService.InitiateSystemUpdate(
		request.UpdateType,
		request.Version,
		request.Force,
		request.Backup,
		request.AnalyzeImpact,
		initiatedBy,
	)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "فشل في بدء التحديث")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "تم بدء عملية التحديث بنجاح",
		"data": map[string]interface{}{
			"requestId":         generateRequestID(),
			"currentVersion":    "1.0.0",
			"targetVersion":     request.Version,
			"updateId":          updateResult.UpdateID,
			"estimatedDuration": updateResult.EstimatedDuration,
			"requiresRestart":   updateResult.RequiresRestart,
			"impactAnalysis":    updateResult.ImpactAnalysis,
			"backupCreated":     request.Backup,
			"steps":             updateResult.Steps,
		},
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

// GetSystemStatus يحصل على حالة النظام
func (h *AdminHandler) GetSystemStatus(w http.ResponseWriter, r *http.Request) {
	systemStatus, err := h.adminService.GetSystemStatus()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "فشل في جلب حالة النظام")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "تم جلب حالة النظام بنجاح",
		"data":    systemStatus,
		"summary": h.generateSystemSummary(systemStatus),
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

// GetAIAnalytics يحصل على تحليلات الذكاء الاصطناعي
func (h *AdminHandler) GetAIAnalytics(w http.ResponseWriter, r *http.Request) {
	timeframe := r.URL.Query().Get("timeframe")
	if timeframe == "" {
		timeframe = "7d"
	}

	analysisType := r.URL.Query().Get("analysisType")
	if analysisType == "" {
		analysisType = "comprehensive"
	}

	analytics, err := h.adminService.GetAIAnalytics(timeframe, analysisType)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "فشل في تحليل النظام باستخدام الذكاء الاصطناعي")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"success":      true,
		"message":      "تم تحليل النظام باستخدام الذكاء الاصطناعي بنجاح",
		"data":         analytics,
		"timeframe":    timeframe,
		"analysisType": analysisType,
		"confidence":   h.calculateAIAnalysisConfidence(analytics),
		"timestamp":    time.Now().Format(time.RFC3339),
	})
}

// GetUserAnalytics يحصل على تحليلات المستخدمين
func (h *AdminHandler) GetUserAnalytics(w http.ResponseWriter, r *http.Request) {
	timeframe := r.URL.Query().Get("timeframe")
	if timeframe == "" {
		timeframe = "30d"
	}

	userSegment := r.URL.Query().Get("userSegment")
	if userSegment == "" {
		userSegment = "all"
	}

	analysisDepth := r.URL.Query().Get("analysisDepth")
	if analysisDepth == "" {
		analysisDepth = "standard"
	}

	analytics, err := h.adminService.GetUserAnalytics(timeframe, userSegment, analysisDepth)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "فشل في تحليل سلوك المستخدمين")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"success":     true,
		"message":     "تم تحليل سلوك المستخدمين بنجاح",
		"data":        analytics,
		"timeframe":   timeframe,
		"userSegment": userSegment,
		"timestamp":   time.Now().Format(time.RFC3339),
	})
}

// GetSystemHealth يحصل على صحة النظام
func (h *AdminHandler) GetSystemHealth(w http.ResponseWriter, r *http.Request) {
	healthChecks, err := h.adminService.PerformHealthChecks()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "فشل في فحص صحة النظام")
		return
	}

	aiAnalysis := h.analyzeHealthWithAI(healthChecks)
	overallStatus := h.calculateOverallHealthStatus(healthChecks, aiAnalysis)

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"status":          overallStatus.Status,
			"score":           overallStatus.Score,
			"checks":          healthChecks,
			"aiAnalysis":      aiAnalysis,
			"recommendations": overallStatus.Recommendations,
			"criticalIssues":  h.getCriticalIssues(healthChecks),
			"timestamp":       time.Now().Format(time.RFC3339),
		},
	})
}

// SetMaintenanceMode يضبط وضع الصيانة
func (h *AdminHandler) SetMaintenanceMode(w http.ResponseWriter, r *http.Request) {
	var request models.MaintenanceRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		respondWithError(w, http.StatusBadRequest, "طلب غير صالح")
		return
	}

	if request.Enabled && request.Message == "" {
		respondWithError(w, http.StatusBadRequest, "رسالة الصيانة مطلوبة عند التفعيل")
		return
	}

	// في الواقع، يجب الحصول على المستخدم من التوكن
	initiatedBy := "admin-user"

	result, err := h.adminService.SetMaintenanceMode(request.Enabled, request.Message, request.Schedule, request.Duration, initiatedBy)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "فشل في تحديث وضع الصيانة")
		return
	}

	message := "تم تعطيل وضع الصيانة"
	if request.Enabled {
		message = "تم تفعيل وضع الصيانة"
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"success":   true,
		"message":   message,
		"data":      result,
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

// GetSystemLogs يحصل على سجلات النظام
func (h *AdminHandler) GetSystemLogs(w http.ResponseWriter, r *http.Request) {
	level := r.URL.Query().Get("level")
	if level == "" {
		level = "all"
	}

	limitStr := r.URL.Query().Get("limit")
	limit := 100
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	pageStr := r.URL.Query().Get("page")
	page := 1
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil {
			page = p
		}
	}

	analyzeStr := r.URL.Query().Get("analyze")
	analyze := true
	if analyzeStr == "false" {
		analyze = false
	}

	logs, err := h.adminService.GetSystemLogs(level, limit, page, nil, nil, analyze)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "فشل في جلب سجلات النظام")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"success":   true,
		"data":      logs,
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

// CreateSystemBackup ينشئ نسخة احتياطية
func (h *AdminHandler) CreateSystemBackup(w http.ResponseWriter, r *http.Request) {
	var request models.BackupRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		respondWithError(w, http.StatusBadRequest, "طلب غير صالح")
		return
	}

	// في الواقع، يجب الحصول على المستخدم من التوكن
	initiatedBy := "admin-user"

	backup, err := h.adminService.CreateSystemBackup(request.Type, request.IncludeLogs, request.Optimize, request.Schedule, initiatedBy)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "فشل في إنشاء النسخة الاحتياطية")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"success":   true,
		"message":   "تم إنشاء النسخة الاحتياطية بنجاح",
		"data":      backup,
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

// PerformOptimization ينفذ التحسين
func (h *AdminHandler) PerformOptimization(w http.ResponseWriter, r *http.Request) {
	var request models.OptimizationRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		respondWithError(w, http.StatusBadRequest, "طلب غير صالح")
		return
	}

	result, err := h.adminService.PerformOptimization(request.Areas, request.Intensity)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "فشل في تحسين النظام")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"success":   true,
		"message":   "تم تحسين النظام بنجاح",
		"data":      result,
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

// ==================== الدوال المساعدة ====================

func (h *AdminHandler) generateSystemSummary(status *models.SystemStatus) models.SystemSummary {
	issues := []string{}
	recommendations := []string{}

	if status.Disk.Threshold != "HEALTHY" {
		issues = append(issues, "مساحة تخزين منخفضة")
		recommendations = append(recommendations, "تفريغ السجلات القديمة", "تحسين تخزين الملفات المؤقتة")
	}
	if !status.Database.Connected {
		issues = append(issues, "مشكلة في قاعدة البيانات")
	}

	overall := "healthy"
	if len(issues) > 0 {
		overall = "degraded"
	}

	return models.SystemSummary{
		Overall:         overall,
		Issues:          issues,
		Recommendations: recommendations,
	}
}

func (h *AdminHandler) calculateAIAnalysisConfidence(analytics *models.AIAnalyticsResult) float64 {
	// محاكاة حساب الثقة في التحليل بناءً على جودة البيانات
	return 85.5
}

func (h *AdminHandler) analyzeHealthWithAI(healthChecks []models.HealthCheckResult) models.AIHealthAnalysis {
	criticalCount := 0
	warningCount := 0

	for _, check := range healthChecks {
		if check.Status == "critical" {
			criticalCount++
		} else if check.Status == "warning" {
			warningCount++
		}
	}

	healthScore := 100.0 - (float64(criticalCount)*30 + float64(warningCount)*15)
	if healthScore < 0 {
		healthScore = 0
	}

	status := "healthy"
	riskLevel := "low"

	if criticalCount > 0 {
		status = "critical"
		riskLevel = "high"
	} else if warningCount > 0 {
		status = "warning"
		riskLevel = "medium"
	}

	return models.AIHealthAnalysis{
		HealthScore: healthScore,
		Status:      status,
		RiskLevel:   riskLevel,
	}
}

func (h *AdminHandler) calculateOverallHealthStatus(healthChecks []models.HealthCheckResult, aiAnalysis models.AIHealthAnalysis) struct {
	Status         string
	Score          float64
	Recommendations []string
} {
	criticalCount := 0
	warningCount := 0

	for _, check := range healthChecks {
		if check.Status == "critical" {
			criticalCount++
		} else if check.Status == "warning" {
			warningCount++
		}
	}

	status := "healthy"
	if criticalCount > 0 {
		status = "critical"
	} else if warningCount > 0 {
		status = "warning"
	}

	score := aiAnalysis.HealthScore

	recommendations := []string{}
	if criticalCount > 0 {
		recommendations = append(recommendations, "معالجة المشاكل الحرجة فوراً")
	}
	if warningCount > 0 {
		recommendations = append(recommendations, "مراقبة المشاكل التحذيرية")
	}

	return struct {
		Status         string
		Score          float64
		Recommendations []string
	}{
		Status:         status,
		Score:          score,
		Recommendations: recommendations,
	}
}

func (h *AdminHandler) getCriticalIssues(healthChecks []models.HealthCheckResult) []models.HealthCheckResult {
	critical := []models.HealthCheckResult{}
	for _, check := range healthChecks {
		if check.Status == "critical" {
			critical = append(critical, check)
		}
	}
	return critical
}

func generateRequestID() string {
	return "req_" + time.Now().Format("20060102150405")
}

// respondWithJSON يرسل رد JSON
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// respondWithError يرسل رد خطأ
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]interface{}{
		"success": false,
		"error":   message,
		"timestamp": time.Now().Format(time.RFC3339),
	})
}
