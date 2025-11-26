package handlers

import (
	"encoding/json"
	"net/http"
	"nawthtech/backend/internal/models"
	"nawthtech/backend/internal/services"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type AdminHandler struct {
	adminService *services.AdminService
}

func NewAdminHandler(adminService *services.AdminService) *AdminHandler {
	return &AdminHandler{
		adminService: adminService,
	}
}

// GetDashboardData يحصل على بيانات لوحة التحكم
func (h *AdminHandler) GetDashboardData(w http.ResponseWriter, r *http.Request) {
	timeRange := r.URL.Query().Get("timeRange")
	if timeRange == "" {
		timeRange = "month"
	}

	dashboardData, err := h.adminService.GetDashboardData(timeRange)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "فشل في تحميل بيانات لوحة التحكم")
		return
	}

	respondWithJSON(w, http.StatusOK, dashboardData)
}

// GetStoreMetrics يحصل على مقاييس المتجر
func (h *AdminHandler) GetStoreMetrics(w http.ResponseWriter, r *http.Request) {
	timeRange := r.URL.Query().Get("timeRange")
	if timeRange == "" {
		timeRange = "month"
	}

	metrics, err := h.adminService.GetStoreMetrics(timeRange)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "فشل في تحميل مقاييس المتجر")
		return
	}

	respondWithJSON(w, http.StatusOK, metrics)
}

// GetRecentOrders يحصل على أحدث الطلبات
func (h *AdminHandler) GetRecentOrders(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	limit := 10
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	orders, err := h.adminService.GetRecentOrders(limit)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "فشل في تحميل الطلبات الحديثة")
		return
	}

	respondWithJSON(w, http.StatusOK, orders)
}

// GetUserActivity يحصل على نشاط المستخدمين
func (h *AdminHandler) GetUserActivity(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	limit := 10
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	activity, err := h.adminService.GetUserActivity(limit)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "فشل في تحميل نشاط المستخدمين")
		return
	}

	respondWithJSON(w, http.StatusOK, activity)
}

// GetSystemAlerts يحصل على تنبيهات النظام
func (h *AdminHandler) GetSystemAlerts(w http.ResponseWriter, r *http.Request) {
	alerts, err := h.adminService.GetSystemAlerts()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "فشل في تحميل تنبيهات النظام")
		return
	}

	respondWithJSON(w, http.StatusOK, alerts)
}

// ExportReport يصدر تقرير
func (h *AdminHandler) ExportReport(w http.ResponseWriter, r *http.Request) {
	reportType := r.URL.Query().Get("type")
	timeRange := r.URL.Query().Get("timeRange")

	report, err := h.adminService.ExportReport(reportType, timeRange)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "فشل في تصدير التقرير")
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment; filename=report_"+timeRange+".json")
	json.NewEncoder(w).Encode(report)
}

// UpdateOrderStatus يقوم بتحديث حالة الطلب
func (h *AdminHandler) UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID := vars["id"]

	var request struct {
		Status string `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		respondWithError(w, http.StatusBadRequest, "طلب غير صالح")
		return
	}

	err := h.adminService.UpdateOrderStatus(orderID, request.Status)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "فشل في تحديث حالة الطلب")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "تم تحديث حالة الطلب بنجاح"})
}
