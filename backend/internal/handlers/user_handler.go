package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/nawthtech/nawthtech/backend/internal/logger"
	"github.com/nawthtech/nawthtech/backend/internal/services"

	"github.com/go-chi/chi/v5"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// respondJSON دالة مساعدة لإرجاع ردود JSON
func respondJSON(w http.ResponseWriter, data map[string]interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// respondError دالة مساعدة لإرجاع أخطاء
func respondError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": false,
		"error":   message,
	})
}

// ==================== الملف الشخصي والإعدادات ====================

func (h *UserHandler) GetUserProfile(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	
	logger.Stdout.Info("جلب الملف الشخصي للمستخدم", "userID", userID)

	// في الإصدار الحالي، نرجع بيانات تجريبية
	response := map[string]interface{}{
		"success": true,
		"message": "تم جلب الملف الشخصي بنجاح",
		"data": map[string]interface{}{
			"id":        userID,
			"name":      "مستخدم NawthTech",
			"email":     "user@nawthtech.com",
			"avatar":    "",
			"createdAt": "2024-01-01T00:00:00Z",
		},
	}

	respondJSON(w, response)
}

func (h *UserHandler) UpdateUserProfile(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	
	var updateData map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		respondError(w, "بيانات غير صالحة", http.StatusBadRequest)
		return
	}

	logger.Stdout.Info("تحديث الملف الشخصي للمستخدم", "userID", userID, "updateFields", updateData)

	response := map[string]interface{}{
		"success": true,
		"message": "تم تحديث الملف الشخصي بنجاح",
		"data": map[string]interface{}{
			"id":        userID,
			"updatedAt": "2024-01-01T00:00:00Z",
		},
	}

	respondJSON(w, response)
}

func (h *UserHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	
	var passwordData struct {
		CurrentPassword string `json:"currentPassword"`
		NewPassword     string `json:"newPassword"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&passwordData); err != nil {
		respondError(w, "بيانات غير صالحة", http.StatusBadRequest)
		return
	}

	logger.Stdout.Info("تغيير كلمة مرور المستخدم", "userID", userID)

	response := map[string]interface{}{
		"success": true,
		"message": "تم تغيير كلمة المرور بنجاح",
		"data":    nil,
	}

	respondJSON(w, response)
}

func (h *UserHandler) UpdateAvatar(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	
	var avatarData struct {
		AvatarURL string `json:"avatarUrl"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&avatarData); err != nil {
		respondError(w, "بيانات غير صالحة", http.StatusBadRequest)
		return
	}

	logger.Stdout.Info("تحديث صورة الملف الشخصي", "userID", userID, "avatarURL", avatarData.AvatarURL)

	response := map[string]interface{}{
		"success": true,
		"message": "تم تحديث صورة الملف الشخصي بنجاح",
		"data": map[string]interface{}{
			"avatar": avatarData.AvatarURL,
		},
	}

	respondJSON(w, response)
}

// ==================== إعدادات المستخدم ====================

func (h *UserHandler) GetUserSettings(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	
	logger.Stdout.Info("جلب إعدادات المستخدم", "userID", userID)

	response := map[string]interface{}{
		"success": true,
		"message": "تم جلب الإعدادات بنجاح",
		"data": map[string]interface{}{
			"notifications": true,
			"language":      "ar",
			"theme":         "light",
		},
	}

	respondJSON(w, response)
}

func (h *UserHandler) UpdateUserSettings(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	
	var settingsData map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&settingsData); err != nil {
		respondError(w, "بيانات غير صالحة", http.StatusBadRequest)
		return
	}

	logger.Stdout.Info("تحديث إعدادات المستخدم", "userID", userID, "updateFields", settingsData)

	response := map[string]interface{}{
		"success": true,
		"message": "تم تحديث الإعدادات بنجاح",
		"data":    settingsData,
	}

	respondJSON(w, response)
}

// ==================== الطلبات والمشتريات ====================

func (h *UserHandler) GetUserOrders(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	
	query := r.URL.Query()
	status := query.Get("status")
	page, _ := strconv.Atoi(query.Get("page"))
	if page == 0 {
		page = 1
	}
	limit, _ := strconv.Atoi(query.Get("limit"))
	if limit == 0 {
		limit = 10
	}

	logger.Stdout.Info("جلب طلبات المستخدم", "userID", userID, "status", status, "page", page, "limit", limit)

	response := map[string]interface{}{
		"success": true,
		"message": "تم جلب الطلبات بنجاح",
		"data": []interface{}{
			map[string]interface{}{
				"id":     "order-1",
				"status": "completed",
				"total":  150.00,
			},
		},
		"pagination": map[string]interface{}{
			"page":  page,
			"limit": limit,
			"total": 1,
		},
	}

	respondJSON(w, response)
}

func (h *UserHandler) GetUserOrderDetails(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	orderID := chi.URLParam(r, "orderId")

	logger.Stdout.Info("جلب تفاصيل طلب المستخدم", "userID", userID, "orderID", orderID)

	response := map[string]interface{}{
		"success": true,
		"message": "تم جلب تفاصيل الطلب بنجاح",
		"data": map[string]interface{}{
			"id":     orderID,
			"status": "completed",
			"items": []interface{}{
				map[string]interface{}{
					"name":  "خدمة وسائل التواصل الاجتماعي",
					"price": 150.00,
				},
			},
		},
	}

	respondJSON(w, response)
}

// ==================== السلة والمشتريات ====================

func (h *UserHandler) GetUserCart(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	
	logger.Stdout.Info("جلب سلة المستخدم", "userID", userID)

	response := map[string]interface{}{
		"success": true,
		"message": "تم جلب السلة بنجاح",
		"data": map[string]interface{}{
			"items": []interface{}{},
			"total": 0.0,
		},
	}

	respondJSON(w, response)
}

func (h *UserHandler) GetUserWishlist(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	
	query := r.URL.Query()
	page, _ := strconv.Atoi(query.Get("page"))
	if page == 0 {
		page = 1
	}
	limit, _ := strconv.Atoi(query.Get("limit"))
	if limit == 0 {
		limit = 20
	}

	logger.Stdout.Info("جلب قائمة رغبات المستخدم", "userID", userID, "page", page, "limit", limit)

	response := map[string]interface{}{
		"success": true,
		"message": "تم جلب قائمة الرغبات بنجاح",
		"data":    []interface{}{},
		"pagination": map[string]interface{}{
			"page":  page,
			"limit": limit,
			"total": 0,
		},
	}

	respondJSON(w, response)
}

// ==================== الإحصائيات والنشاط ====================

func (h *UserHandler) GetUserStats(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	period := r.URL.Query().Get("period")
	if period == "" {
		period = "30d"
	}

	logger.Stdout.Info("جلب إحصائيات المستخدم", "userID", userID, "period", period)

	response := map[string]interface{}{
		"success": true,
		"message": "تم جلب الإحصائيات بنجاح",
		"data": map[string]interface{}{
			"totalOrders":   5,
			"totalSpent":    750.00,
			"favoriteCategory": "خدمات الوسائل الاجتماعية",
		},
		"period": period,
	}

	respondJSON(w, response)
}

func (h *UserHandler) GetUserActivity(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	
	query := r.URL.Query()
	page, _ := strconv.Atoi(query.Get("page"))
	if page == 0 {
		page = 1
	}
	limit, _ := strconv.Atoi(query.Get("limit"))
	if limit == 0 {
		limit = 20
	}

	logger.Stdout.Info("جلب نشاط المستخدم", "userID", userID, "page", page, "limit", limit)

	response := map[string]interface{}{
		"success": true,
		"message": "تم جلب النشاط بنجاح",
		"data":    []interface{}{},
		"pagination": map[string]interface{}{
			"page":  page,
			"limit": limit,
			"total": 0,
		},
	}

	respondJSON(w, response)
}

// ==================== الإشعارات ====================

func (h *UserHandler) GetUserNotifications(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	
	query := r.URL.Query()
	page, _ := strconv.Atoi(query.Get("page"))
	if page == 0 {
		page = 1
	}
	limit, _ := strconv.Atoi(query.Get("limit"))
	if limit == 0 {
		limit = 20
	}

	logger.Stdout.Info("جلب إشعارات المستخدم", "userID", userID, "page", page, "limit", limit)

	response := map[string]interface{}{
		"success": true,
		"message": "تم جلب الإشعارات بنجاح",
		"data":    []interface{}{},
		"pagination": map[string]interface{}{
			"page":  page,
			"limit": limit,
			"total": 0,
		},
		"unreadCount": 0,
	}

	respondJSON(w, response)
}

func (h *UserHandler) MarkNotificationsAsRead(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	
	var requestData struct {
		NotificationIDs []string `json:"notificationIds"`
		MarkAll         bool     `json:"markAll"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		respondError(w, "بيانات غير صالحة", http.StatusBadRequest)
		return
	}

	logger.Stdout.Info("تعليم الإشعارات كمقروءة", "userID", userID, "notificationIDs", requestData.NotificationIDs, "markAll", requestData.MarkAll)

	message := "تم تعليم الإشعارات كمقروءة"
	if requestData.MarkAll {
		message = "تم تعليم جميع الإشعارات كمقروءة"
	}

	response := map[string]interface{}{
		"success": true,
		"message": message,
		"data": map[string]interface{}{
			"marked": len(requestData.NotificationIDs),
		},
	}

	respondJSON(w, response)
}

func (h *UserHandler) DeleteNotifications(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	
	var requestData struct {
		NotificationIDs []string `json:"notificationIds"`
		DeleteAll       bool     `json:"deleteAll"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		respondError(w, "بيانات غير صالحة", http.StatusBadRequest)
		return
	}

	logger.Stdout.Info("حذف إشعارات المستخدم", "userID", userID, "notificationIDs", requestData.NotificationIDs, "deleteAll", requestData.DeleteAll)

	message := "تم حذف الإشعارات"
	if requestData.DeleteAll {
		message = "تم حذف جميع الإشعارات"
	}

	response := map[string]interface{}{
		"success": true,
		"message": message,
		"data": map[string]interface{}{
			"deleted": len(requestData.NotificationIDs),
		},
	}

	respondJSON(w, response)
}

// ==================== مسارات البائعين ====================

func (h *UserHandler) GetSellerServices(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	
	query := r.URL.Query()
	status := query.Get("status")
	page, _ := strconv.Atoi(query.Get("page"))
	if page == 0 {
		page = 1
	}
	limit, _ := strconv.Atoi(query.Get("limit"))
	if limit == 0 {
		limit = 10
	}

	logger.Stdout.Info("جلب خدمات البائع", "userID", userID, "status", status, "page", page, "limit", limit)

	response := map[string]interface{}{
		"success": true,
		"message": "تم جلب خدمات البائع بنجاح",
		"data":    []interface{}{},
		"pagination": map[string]interface{}{
			"page":  page,
			"limit": limit,
			"total": 0,
		},
	}

	respondJSON(w, response)
}

func (h *UserHandler) GetSellerStats(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	period := r.URL.Query().Get("period")
	if period == "" {
		period = "30d"
	}

	logger.Stdout.Info("جلب إحصائيات البائع", "userID", userID, "period", period)

	response := map[string]interface{}{
		"success": true,
		"message": "تم جلب إحصائيات البائع بنجاح",
		"data": map[string]interface{}{
			"totalServices": 0,
			"totalSales":    0,
			"rating":        0.0,
		},
		"period": period,
	}

	respondJSON(w, response)
}

func (h *UserHandler) GetSellerOrders(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	
	query := r.URL.Query()
	status := query.Get("status")
	page, _ := strconv.Atoi(query.Get("page"))
	if page == 0 {
		page = 1
	}
	limit, _ := strconv.Atoi(query.Get("limit"))
	if limit == 0 {
		limit = 10
	}

	logger.Stdout.Info("جلب طلبات البائع", "userID", userID, "status", status, "page", page, "limit", limit)

	response := map[string]interface{}{
		"success": true,
		"message": "تم جلب طلبات البائع بنجاح",
		"data":    []interface{}{},
		"pagination": map[string]interface{}{
			"page":  page,
			"limit": limit,
			"total": 0,
		},
	}

	respondJSON(w, response)
}