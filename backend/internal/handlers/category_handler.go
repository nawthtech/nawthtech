package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/nawthtech/nawthtech/backend/internal/logger"
	"github.com/nawthtech/nawthtech/backend/internal/services"

	"github.com/go-chi/chi/v5"
)

type CategoryHandler struct {
	categoryService *services.CategoryService
}

func NewCategoryHandler(categoryService *services.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
	}
}

// ==================== الفئات العامة ====================

func (h *CategoryHandler) GetCategories(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	
	includeInactive := query.Get("includeInactive") == "true"
	withStats := query.Get("withStats") == "true"
	page, _ := strconv.Atoi(query.Get("page"))
	if page == 0 {
		page = 1
	}
	limit, _ := strconv.Atoi(query.Get("limit"))
	if limit == 0 {
		limit = 50
	}

	logger.Stdout.Info("جلب الفئات", 
		"includeInactive", includeInactive, 
		"withStats", withStats, 
		"page", page, 
		"limit", limit)

	categories := []map[string]interface{}{
		{
			"id":          "cat_1",
			"name":        "وسائل التواصل الاجتماعي",
			"description": "خدمات تطوير وتحسين حسابات وسائل التواصل الاجتماعي",
			"slug":        "social-media",
			"icon":        "📱",
			"isActive":    true,
			"order":       1,
			"serviceCount": 45,
		},
		{
			"id":          "cat_2",
			"name":        "التصميم والإبداع",
			"description": "خدمات التصميم الجرافيكي والمواد الإبداعية",
			"slug":        "design-creative",
			"icon":        "🎨",
			"isActive":    true,
			"order":       2,
			"serviceCount": 32,
		},
	}

	response := map[string]interface{}{
		"success": true,
		"message": "تم جلب الفئات بنجاح",
		"data":    categories,
		"pagination": map[string]interface{}{
			"page":  page,
			"limit": limit,
			"total": len(categories),
		},
	}

	respondJSON(w, response)
}

func (h *CategoryHandler) GetCategoryTree(w http.ResponseWriter, r *http.Request) {
	includeInactive := r.URL.Query().Get("includeInactive") == "true"

	logger.Stdout.Info("جلب هيكل شجرة الفئات", "includeInactive", includeInactive)

	tree := []map[string]interface{}{
		{
			"id":       "cat_1",
			"name":     "وسائل التواصل الاجتماعي",
			"slug":     "social-media",
			"isActive": true,
			"children": []map[string]interface{}{
				{
					"id":       "cat_1_1",
					"name":     "إنستغرام",
					"slug":     "instagram",
					"isActive": true,
					"children": []map[string]interface{}{},
				},
				{
					"id":       "cat_1_2",
					"name":     "تويتر",
					"slug":     "twitter",
					"isActive": true,
					"children": []map[string]interface{}{},
				},
			},
		},
		{
			"id":       "cat_2",
			"name":     "التصميم والإبداع",
			"slug":     "design-creative",
			"isActive": true,
			"children": []map[string]interface{}{},
		},
	}

	response := map[string]interface{}{
		"success": true,
		"message": "تم جلب هيكل الفئات بنجاح",
		"data":    tree,
	}

	respondJSON(w, response)
}

func (h *CategoryHandler) GetCategoryById(w http.ResponseWriter, r *http.Request) {
	categoryID := chi.URLParam(r, "categoryId")

	logger.Stdout.Info("جلب تفاصيل فئة", "categoryID", categoryID)

	response := map[string]interface{}{
		"success": true,
		"message": "تم جلب تفاصيل الفئة بنجاح",
		"data": map[string]interface{}{
			"id":          categoryID,
			"name":        "وسائل التواصل الاجتماعي",
			"description": "خدمات تطوير وتحسين حسابات وسائل التواصل الاجتماعي",
			"slug":        "social-media",
			"icon":        "📱",
			"isActive":    true,
			"order":       1,
			"parentId":    nil,
			"createdAt":   "2024-01-01T00:00:00Z",
			"updatedAt":   "2024-01-01T00:00:00Z",
			"stats": map[string]interface{}{
				"serviceCount": 45,
				"activeServices": 42,
				"averageRating": 4.7,
			},
		},
	}

	respondJSON(w, response)
}

func (h *CategoryHandler) GetCategoryServices(w http.ResponseWriter, r *http.Request) {
	categoryID := chi.URLParam(r, "categoryId")
	
	query := r.URL.Query()
	page, _ := strconv.Atoi(query.Get("page"))
	if page == 0 {
		page = 1
	}
	limit, _ := strconv.Atoi(query.Get("limit"))
	if limit == 0 {
		limit = 12
	}

	logger.Stdout.Info("جلب خدمات الفئة", 
		"categoryID", categoryID, 
		"page", page, 
		"limit", limit)

	services := []map[string]interface{}{
		{
			"id":          "service_1",
			"name":        "متابعين إنستغرام - 1000 متابع",
			"description": "زيادة المتابعين بشكل طبيعي وآمن",
			"price":       150.00,
			"rating":      4.8,
			"reviews":     1250,
			"inStock":     true,
		},
		{
			"id":          "service_2",
			"name":        "لايكات إنستغرام - 5000 لايك",
			"description": "زيادة التفاعل على منشوراتك",
			"price":       75.00,
			"rating":      4.6,
			"reviews":     890,
			"inStock":     true,
		},
	}

	response := map[string]interface{}{
		"success": true,
		"message": "تم جلب خدمات الفئة بنجاح",
		"data":    services,
		"pagination": map[string]interface{}{
			"page":  page,
			"limit": limit,
			"total": len(services),
		},
	}

	respondJSON(w, response)
}

func (h *CategoryHandler) GetSubcategories(w http.ResponseWriter, r *http.Request) {
	categoryID := chi.URLParam(r, "categoryId")

	logger.Stdout.Info("جلب الفئات الفرعية", "categoryID", categoryID)

	subcategories := []map[string]interface{}{
		{
			"id":          "subcat_1",
			"name":        "إنستغرام",
			"description": "خدمات خاصة بمنصة إنستغرام",
			"slug":        "instagram",
			"icon":        "📸",
			"isActive":    true,
			"order":       1,
			"serviceCount": 25,
		},
		{
			"id":          "subcat_2",
			"name":        "تويتر",
			"description": "خدمات خاصة بمنصة تويتر",
			"slug":        "twitter",
			"icon":        "🐦",
			"isActive":    true,
			"order":       2,
			"serviceCount": 15,
		},
	}

	response := map[string]interface{}{
		"success": true,
		"message": "تم جلب الفئات الفرعية بنجاح",
		"data":    subcategories,
	}

	respondJSON(w, response)
}

// ==================== إدارة الفئات (للمسؤولين) ====================

func (h *CategoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	
	var categoryData map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&categoryData); err != nil {
		respondError(w, "بيانات غير صالحة", http.StatusBadRequest)
		return
	}

	logger.Stdout.Info("إنشاء فئة جديدة", 
		"adminID", userID, 
		"categoryName", categoryData["name"], 
		"parentID", categoryData["parent"])

	response := map[string]interface{}{
		"success": true,
		"message": "تم إنشاء الفئة بنجاح",
		"data": map[string]interface{}{
			"id":          "cat_new",
			"name":        categoryData["name"],
			"description": categoryData["description"],
			"slug":        categoryData["slug"],
			"icon":        categoryData["icon"],
			"isActive":    true,
			"order":       3,
			"parentId":    categoryData["parent"],
			"createdAt":   "2024-01-01T00:00:00Z",
		},
		"categoryId": "cat_new",
	}

	w.WriteHeader(http.StatusCreated)
	respondJSON(w, response)
}

func (h *CategoryHandler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	categoryID := chi.URLParam(r, "categoryId")
	
	var updateData map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		respondError(w, "بيانات غير صالحة", http.StatusBadRequest)
		return
	}

	logger.Stdout.Info("تحديث فئة", 
		"adminID", userID, 
		"categoryID", categoryID, 
		"updateFields", updateData)

	response := map[string]interface{}{
		"success": true,
		"message": "تم تحديث الفئة بنجاح",
		"data": map[string]interface{}{
			"id":        categoryID,
			"updatedAt": "2024-01-01T00:00:00Z",
			"changes":   updateData,
		},
	}

	respondJSON(w, response)
}

func (h *CategoryHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	categoryID := chi.URLParam(r, "categoryId")

	logger.Stdout.Info("حذف فئة", "adminID", userID, "categoryID", categoryID)

	response := map[string]interface{}{
		"success": true,
		"message": "تم حذف الفئة بنجاح",
		"data": map[string]interface{}{
			"deleted":    true,
			"categoryId": categoryID,
		},
	}

	respondJSON(w, response)
}

func (h *CategoryHandler) UpdateCategoryStatus(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	categoryID := chi.URLParam(r, "categoryId")
	
	var statusData struct {
		IsActive bool `json:"isActive"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&statusData); err != nil {
		respondError(w, "بيانات غير صالحة", http.StatusBadRequest)
		return
	}

	statusText := "تفعيل"
	if !statusData.IsActive {
		statusText = "تعطيل"
	}

	logger.Stdout.Info("تحديث حالة الفئة", 
		"adminID", userID, 
		"categoryID", categoryID, 
		"newStatus", statusText)

	response := map[string]interface{}{
		"success": true,
		"message": "تم " + statusText + " الفئة بنجاح",
		"data": map[string]interface{}{
			"id":       categoryID,
			"isActive": statusData.IsActive,
			"updatedAt": "2024-01-01T00:00:00Z",
		},
	}

	respondJSON(w, response)
}

func (h *CategoryHandler) UpdateCategoryOrder(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	categoryID := chi.URLParam(r, "categoryId")
	
	var orderData struct {
		Order int `json:"order"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&orderData); err != nil {
		respondError(w, "بيانات غير صالحة", http.StatusBadRequest)
		return
	}

	logger.Stdout.Info("تحديث ترتيب الفئة", 
		"adminID", userID, 
		"categoryID", categoryID, 
		"newOrder", orderData.Order)

	response := map[string]interface{}{
		"success": true,
		"message": "تم تحديث ترتيب الفئة بنجاح",
		"data": map[string]interface{}{
			"id":    categoryID,
			"order": orderData.Order,
			"updatedAt": "2024-01-01T00:00:00Z",
		},
	}

	respondJSON(w, response)
}

// ==================== إحصائيات الفئات ====================

func (h *CategoryHandler) GetCategoriesStats(w http.ResponseWriter, r *http.Request) {
	logger.Stdout.Info("جلب إحصائيات الفئات الشاملة")

	response := map[string]interface{}{
		"success": true,
		"message": "تم جلب إحصائيات الفئات بنجاح",
		"data": map[string]interface{}{
			"totalCategories":     15,
			"activeCategories":    12,
			"totalServices":       543,
			"averageServicesPerCategory": 36.2,
			"mostPopularCategory": map[string]interface{}{
				"id":    "cat_1",
				"name":  "وسائل التواصل الاجتماعي",
				"count": 45,
			},
			"categoryDistribution": []map[string]interface{}{
				{
					"category": "وسائل التواصل الاجتماعي",
					"count":    45,
					"percentage": 28.3,
				},
				{
					"category": "التصميم والإبداع",
					"count":    32,
					"percentage": 20.1,
				},
			},
		},
	}

	respondJSON(w, response)
}

func (h *CategoryHandler) GetCategoryStats(w http.ResponseWriter, r *http.Request) {
	categoryID := chi.URLParam(r, "categoryId")

	logger.Stdout.Info("جلب إحصائيات فئة", "categoryID", categoryID)

	response := map[string]interface{}{
		"success": true,
		"message": "تم جلب إحصائيات الفئة بنجاح",
		"data": map[string]interface{}{
			"categoryId":   categoryID,
			"serviceCount": 45,
			"activeServices": 42,
			"inactiveServices": 3,
			"averageRating": 4.7,
			"totalReviews": 1250,
			"totalRevenue": 125430.00,
			"popularServices": []map[string]interface{}{
				{
					"id":    "service_1",
					"name":  "متابعين إنستغرام - 1000 متابع",
					"sales": 890,
				},
				{
					"id":    "service_2",
					"name":  "لايكات إنستغرام - 5000 لايك",
					"sales": 543,
				},
			},
			"monthlyGrowth": 12.5,
		},
	}

	respondJSON(w, response)
}

// ==================== البحث والتصفية ====================

func (h *CategoryHandler) SearchCategories(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	
	searchQuery := query.Get("q")
	searchType := query.Get("type")
	page, _ := strconv.Atoi(query.Get("page"))
	if page == 0 {
		page = 1
	}
	limit, _ := strconv.Atoi(query.Get("limit"))
	if limit == 0 {
		limit = 20
	}

	logger.Stdout.Info("البحث في الفئات", 
		"query", searchQuery, 
		"type", searchType, 
		"page", page, 
		"limit", limit)

	categories := []map[string]interface{}{
		{
			"id":          "search_cat_1",
			"name":        "نتيجة البحث: " + searchQuery,
			"description": "فئة متوافقة مع بحثك",
			"slug":        "search-result",
			"isActive":    true,
			"serviceCount": 15,
		},
	}

	response := map[string]interface{}{
		"success": true,
		"message": "تم البحث في الفئات بنجاح",
		"data":    categories,
		"pagination": map[string]interface{}{
			"page":  page,
			"limit": limit,
			"total": len(categories),
		},
		"searchQuery": searchQuery,
	}

	respondJSON(w, response)
}

// ==================== الاستيراد والتصدير ====================

func (h *CategoryHandler) ImportCategories(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	
	var importData struct {
		Data   []map[string]interface{} `json:"data"`
		Format string                   `json:"format"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&importData); err != nil {
		respondError(w, "بيانات غير صالحة", http.StatusBadRequest)
		return
	}

	logger.Stdout.Info("استيراد فئات", 
		"adminID", userID, 
		"format", importData.Format, 
		"itemsCount", len(importData.Data))

	response := map[string]interface{}{
		"success": true,
		"message": "تم استيراد الفئات بنجاح",
		"data": map[string]interface{}{
			"imported": len(importData.Data),
			"failed":   0,
		},
		"importedCount": len(importData.Data),
		"failedCount":   0,
	}

	respondJSON(w, response)
}

func (h *CategoryHandler) ExportCategories(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	format := r.URL.Query().Get("format")
	if format == "" {
		format = "json"
	}

	logger.Stdout.Info("تصدير الفئات", "adminID", userID, "format", format)

	response := map[string]interface{}{
		"success": true,
		"message": "تم تصدير الفئات بنجاح",
		"data": map[string]interface{}{
			"format":     format,
			"itemCount":  15,
			"exportedAt": "2024-01-01T00:00:00Z",
		},
		"downloadUrl": "/api/v1/categories/export/download?format=" + format,
	}

	respondJSON(w, response)
}