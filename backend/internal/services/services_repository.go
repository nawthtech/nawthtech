package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/nawthtech/nawthtech/backend/internal/models"
	"github.com/nawthtech/nawthtech/backend/internal/utils"
	"gorm.io/gorm"
)

// ServicesRepository واجهة مستودع الخدمات
type ServicesRepository interface {
	// الدوال الأساسية
	GetServices(ctx context.Context, params GetServicesParams) ([]models.Service, *utils.Pagination, error)
	SearchServices(ctx context.Context, params SearchServicesParams) ([]models.Service, *utils.Pagination, error)
	GetFeaturedServices(ctx context.Context, limit int) ([]models.Service, error)
	GetServiceByID(ctx context.Context, serviceID string) (*models.Service, error)
	GetRecommendedServices(ctx context.Context, serviceID string, limit int) ([]models.Service, error)
	GetSellerServices(ctx context.Context, params GetSellerServicesParams) ([]models.Service, *utils.Pagination, error)
	
	// إدارة الخدمات
	CreateService(ctx context.Context, service *models.Service) (*models.Service, error)
	UpdateService(ctx context.Context, service *models.Service) (*models.Service, error)
	UpdateServiceStatus(ctx context.Context, serviceID string, status string) (*models.Service, error)
	DeleteService(ctx context.Context, serviceID string) error
	BulkUpdateServices(ctx context.Context, serviceIDs []string, updates map[string]interface{}) error
	
	// التقييمات والمراجعات
	AddRating(ctx context.Context, params AddRatingParams) (*models.Rating, error)
	GetServiceRatings(ctx context.Context, serviceID string, page, limit int) ([]models.Rating, *utils.Pagination, error)
	UpdateRating(ctx context.Context, ratingID string, updates map[string]interface{}) error
	DeleteRating(ctx context.Context, ratingID string) error
	GetUserRatingForService(ctx context.Context, serviceID, userID string) (*models.Rating, error)
	
	// التوفر والجدولة
	CheckAvailability(ctx context.Context, params CheckAvailabilityParams) (*models.Availability, error)
	CreateTimeSlot(ctx context.Context, slot *models.TimeSlot) error
	GetAvailableTimeSlots(ctx context.Context, serviceID, date string) ([]models.TimeSlot, error)
	
	// الإحصائيات والتقارير
	GetServicesStats(ctx context.Context, userID string, timeframe string) (*models.ServicesStats, error)
	GetPopularServices(ctx context.Context, limit int, category string) ([]models.Service, error)
	GetServicesByCategory(ctx context.Context, category string, limit int) ([]models.Service, error)
	GetServicesGrowth(ctx context.Context, sellerID string, period string) (*models.ServicesGrowth, error)
	
	// الفئات والوسوم
	GetAllCategories(ctx context.Context) ([]string, error)
	GetPopularTags(ctx context.Context, limit int) ([]string, error)
	GetServicesByTag(ctx context.Context, tag string, page, limit int) ([]models.Service, *utils.Pagination, error)
	
	// إدارة المتجر
	IncreaseServiceOrders(ctx context.Context, serviceID string) error
	UpdateServiceRating(ctx context.Context, serviceID string) error
	GetSimilarServices(ctx context.Context, service *models.Service, limit int) ([]models.Service, error)
	
	// البحث المتقدم
	AdvancedSearch(ctx context.Context, params AdvancedSearchParams) ([]models.Service, *utils.Pagination, error)
	
	// الإدارة
	GetAllServices(ctx context.Context, page, limit int, status string) ([]models.Service, *utils.Pagination, error)
	CountServicesByStatus(ctx context.Context, sellerID string) (map[string]int, error)
}

// AdvancedSearchParams معاملات البحث المتقدم
type AdvancedSearchParams struct {
	Query     string   `json:"query"`
	Category  string   `json:"category"`
	Tags      []string `json:"tags"`
	MinPrice  float64  `json:"min_price"`
	MaxPrice  float64  `json:"max_price"`
	MinRating float64  `json:"min_rating"`
	SellerID  string   `json:"seller_id"`
	Status    string   `json:"status"`
	SortBy    string   `json:"sort_by"`
	SortOrder string   `json:"sort_order"`
	Page      int      `json:"page"`
	Limit     int      `json:"limit"`
}

// TimeSlot نموذج الفترة الزمنية
type TimeSlot struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	ServiceID string    `json:"service_id" gorm:"not null;index"`
	Date      string    `json:"date" gorm:"not null"` // YYYY-MM-DD
	StartTime string    `json:"start_time" gorm:"not null"` // HH:MM
	EndTime   string    `json:"end_time" gorm:"not null"`   // HH:MM
	Available bool      `json:"available" gorm:"default:true"`
	Booked    bool      `json:"booked" gorm:"default:false"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ServicesGrowth نموذج نمو الخدمات
type ServicesGrowth struct {
	SellerID      string    `json:"seller_id"`
	Period        string    `json:"period"`
	NewServices   int       `json:"new_services"`
	TotalOrders   int       `json:"total_orders"`
	Revenue       float64   `json:"revenue"`
	GrowthRate    float64   `json:"growth_rate"`
	CalculatedAt  time.Time `json:"calculated_at"`
}

// servicesRepositoryImpl التطبيق الفعلي لمستودع الخدمات
type servicesRepositoryImpl struct {
	db *gorm.DB
}

// NewServicesRepository إنشاء مستودع خدمات جديد
func NewServicesRepository(db *gorm.DB) ServicesRepository {
	return &servicesRepositoryImpl{
		db: db,
	}
}

// الدوال الأساسية (تم تنفيذها سابقاً)
func (r *servicesRepositoryImpl) GetServices(ctx context.Context, params GetServicesParams) ([]models.Service, *utils.Pagination, error) {
	var services []models.Service
	var total int64

	query := r.db.Model(&models.Service{}).Where("status = ?", "active")

	// تطبيق الفلاتر
	if params.Category != "" {
		query = query.Where("category = ?", params.Category)
	}
	if params.MinPrice > 0 {
		query = query.Where("price >= ?", params.MinPrice)
	}
	if params.MaxPrice > 0 {
		query = query.Where("price <= ?", params.MaxPrice)
	}

	// حساب العدد الإجمالي
	if err := query.Count(&total).Error; err != nil {
		return nil, nil, err
	}

	// تطبيق الترتيب
	if params.SortBy != "" {
		order := params.SortBy
		if params.SortOrder != "" {
			order = order + " " + params.SortOrder
		}
		query = query.Order(order)
	} else {
		query = query.Order("created_at DESC")
	}

	// تطبيق التقسيم
	offset := (params.Page - 1) * params.Limit
	if err := query.Offset(offset).Limit(params.Limit).Find(&services).Error; err != nil {
		return nil, nil, err
	}

	// حساب معلومات التقسيم
	pages := int(total) / params.Limit
	if int(total)%params.Limit > 0 {
		pages++
	}

	pagination := &utils.Pagination{
		Page:  params.Page,
		Limit: params.Limit,
		Total: int(total),
		Pages: pages,
	}

	return services, pagination, nil
}

func (r *servicesRepositoryImpl) SearchServices(ctx context.Context, params SearchServicesParams) ([]models.Service, *utils.Pagination, error) {
	var services []models.Service
	var total int64

	query := r.db.Model(&models.Service{}).Where("status = ?", "active")

	// تطبيق البحث
	if params.Query != "" {
		searchQuery := "%" + params.Query + "%"
		query = query.Where("title LIKE ? OR description LIKE ? OR tags LIKE ?", searchQuery, searchQuery, searchQuery)
	}

	if params.Category != "" {
		query = query.Where("category = ?", params.Category)
	}

	// حساب العدد الإجمالي
	if err := query.Count(&total).Error; err != nil {
		return nil, nil, err
	}

	// تطبيق التقسيم
	offset := (params.Page - 1) * params.Limit
	if err := query.Offset(offset).Limit(params.Limit).Order("created_at DESC").Find(&services).Error; err != nil {
		return nil, nil, err
	}

	// حساب معلومات التقسيم
	pages := int(total) / params.Limit
	if int(total)%params.Limit > 0 {
		pages++
	}

	pagination := &utils.Pagination{
		Page:  params.Page,
		Limit: params.Limit,
		Total: int(total),
		Pages: pages,
	}

	return services, pagination, nil
}

func (r *servicesRepositoryImpl) GetFeaturedServices(ctx context.Context, limit int) ([]models.Service, error) {
	var services []models.Service

	err := r.db.Where("featured = ? AND status = ?", true, "active").
		Order("rating DESC, total_orders DESC").
		Limit(limit).
		Find(&services).Error

	if err != nil {
		return nil, err
	}

	return services, nil
}

func (r *servicesRepositoryImpl) GetServiceByID(ctx context.Context, serviceID string) (*models.Service, error) {
	var service models.Service

	err := r.db.Preload("Reviews", func(db *gorm.DB) *gorm.DB {
		return db.Order("created_at DESC")
	}).First(&service, "id = ?", serviceID).Error
	
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("الخدمة غير موجودة")
		}
		return nil, err
	}

	return &service, nil
}

func (r *servicesRepositoryImpl) GetRecommendedServices(ctx context.Context, serviceID string, limit int) ([]models.Service, error) {
	var services []models.Service

	// الحصول على الخدمة الحالية لمعرفة فئتها
	var currentService models.Service
	if err := r.db.First(&currentService, "id = ?", serviceID).Error; err != nil {
		return nil, err
	}

	// جلب خدمات من نفس الفئة (باستثناء الخدمة الحالية)
	err := r.db.Where("category = ? AND id != ? AND status = ?", 
		currentService.Category, serviceID, "active").
		Order("rating DESC, total_orders DESC").
		Limit(limit).
		Find(&services).Error

	if err != nil {
		return nil, err
	}

	return services, nil
}

func (r *servicesRepositoryImpl) GetSellerServices(ctx context.Context, params GetSellerServicesParams) ([]models.Service, *utils.Pagination, error) {
	var services []models.Service
	var total int64

	query := r.db.Model(&models.Service{}).Where("seller_id = ?", params.SellerID)

	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}

	// حساب العدد الإجمالي
	if err := query.Count(&total).Error; err != nil {
		return nil, nil, err
	}

	// تطبيق التقسيم
	offset := (params.Page - 1) * params.Limit
	if err := query.Offset(offset).Limit(params.Limit).Order("created_at DESC").Find(&services).Error; err != nil {
		return nil, nil, err
	}

	// حساب معلومات التقسيم
	pages := int(total) / params.Limit
	if int(total)%params.Limit > 0 {
		pages++
	}

	pagination := &utils.Pagination{
		Page:  params.Page,
		Limit: params.Limit,
		Total: int(total),
		Pages: pages,
	}

	return services, pagination, nil
}

// الدوال الجديدة المكتملة

func (r *servicesRepositoryImpl) CreateService(ctx context.Context, service *models.Service) (*models.Service, error) {
	if err := r.db.Create(service).Error; err != nil {
		return nil, err
	}
	return service, nil
}

func (r *servicesRepositoryImpl) UpdateService(ctx context.Context, service *models.Service) (*models.Service, error) {
	if err := r.db.Save(service).Error; err != nil {
		return nil, err
	}
	return service, nil
}

func (r *servicesRepositoryImpl) UpdateServiceStatus(ctx context.Context, serviceID string, status string) (*models.Service, error) {
	var service models.Service
	if err := r.db.Model(&service).Where("id = ?", serviceID).Update("status", status).Error; err != nil {
		return nil, err
	}

	// جلب الخدمة المحدثة
	if err := r.db.First(&service, "id = ?", serviceID).Error; err != nil {
		return nil, err
	}

	return &service, nil
}

func (r *servicesRepositoryImpl) DeleteService(ctx context.Context, serviceID string) error {
	return r.db.Where("id = ?", serviceID).Delete(&models.Service{}).Error
}

func (r *servicesRepositoryImpl) BulkUpdateServices(ctx context.Context, serviceIDs []string, updates map[string]interface{}) error {
	return r.db.Model(&models.Service{}).
		Where("id IN ?", serviceIDs).
		Updates(updates).Error
}

// التقييمات والمراجعات

func (r *servicesRepositoryImpl) AddRating(ctx context.Context, params AddRatingParams) (*models.Rating, error) {
	// التحقق من عدم وجود تقييم سابق من نفس المستخدم
	var existingRating models.Rating
	err := r.db.Where("service_id = ? AND user_id = ?", params.ServiceID, params.UserID).
		First(&existingRating).Error
	
	if err == nil {
		return nil, fmt.Errorf("لقد قمت بتقييم هذه الخدمة مسبقاً")
	} else if err != gorm.ErrRecordNotFound {
		return nil, err
	}

	rating := &models.Rating{
		ID:        fmt.Sprintf("rating_%d", time.Now().UnixNano()),
		ServiceID: params.ServiceID,
		UserID:    params.UserID,
		Rating:    params.Rating,
		Comment:   params.Comment,
		CreatedAt: time.Now(),
	}

	if err := r.db.Create(rating).Error; err != nil {
		return nil, err
	}

	// تحديث متوسط التقييم للخدمة
	if err := r.UpdateServiceRating(ctx, params.ServiceID); err != nil {
		return nil, err
	}

	return rating, nil
}

func (r *servicesRepositoryImpl) GetServiceRatings(ctx context.Context, serviceID string, page, limit int) ([]models.Rating, *utils.Pagination, error) {
	var ratings []models.Rating
	var total int64

	query := r.db.Model(&models.Rating{}).Where("service_id = ?", serviceID)

	// حساب العدد الإجمالي
	if err := query.Count(&total).Error; err != nil {
		return nil, nil, err
	}

	// تطبيق التقسيم
	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).
		Order("created_at DESC").
		Find(&ratings).Error; err != nil {
		return nil, nil, err
	}

	// حساب معلومات التقسيم
	pages := int(total) / limit
	if int(total)%limit > 0 {
		pages++
	}

	pagination := &utils.Pagination{
		Page:  page,
		Limit: limit,
		Total: int(total),
		Pages: pages,
	}

	return ratings, pagination, nil
}

func (r *servicesRepositoryImpl) UpdateRating(ctx context.Context, ratingID string, updates map[string]interface{}) error {
	return r.db.Model(&models.Rating{}).
		Where("id = ?", ratingID).
		Updates(updates).Error
}

func (r *servicesRepositoryImpl) DeleteRating(ctx context.Context, ratingID string) error {
	return r.db.Where("id = ?", ratingID).Delete(&models.Rating{}).Error
}

func (r *servicesRepositoryImpl) GetUserRatingForService(ctx context.Context, serviceID, userID string) (*models.Rating, error) {
	var rating models.Rating
	err := r.db.Where("service_id = ? AND user_id = ?", serviceID, userID).
		First(&rating).Error
	
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &rating, nil
}

// التوفر والجدولة

func (r *servicesRepositoryImpl) CheckAvailability(ctx context.Context, params CheckAvailabilityParams) (*models.Availability, error) {
	// هذا تنفيذ مبسط - في التطبيق الحقيقي قد تحتاج للتحقق من الجدول الزمني
	var existingBookings int64
	err := r.db.Model(&TimeSlot{}).
		Where("service_id = ? AND date = ? AND start_time = ? AND booked = ?", 
			params.ServiceID, params.Date, params.Time, true).
		Count(&existingBookings).Error
	
	if err != nil {
		return nil, err
	}

	available := existingBookings == 0
	message := "الخدمة متاحة في الوقت المحدد"
	if !available {
		message = "الخدمة غير متاحة في الوقت المحدد"
	}

	availability := &models.Availability{
		Available:      available,
		ServiceID:      params.ServiceID,
		Date:           params.Date,
		Time:           params.Time,
		Guests:         params.Guests,
		Message:        message,
		SuggestedTimes: []string{"10:00", "14:00", "16:00"},
		CheckedAt:      time.Now(),
	}

	return availability, nil
}

func (r *servicesRepositoryImpl) CreateTimeSlot(ctx context.Context, slot *TimeSlot) error {
	slot.ID = fmt.Sprintf("slot_%d", time.Now().UnixNano())
	return r.db.Create(slot).Error
}

func (r *servicesRepositoryImpl) GetAvailableTimeSlots(ctx context.Context, serviceID, date string) ([]TimeSlot, error) {
	var slots []TimeSlot
	err := r.db.Where("service_id = ? AND date = ? AND available = ? AND booked = ?",
		serviceID, date, true, false).
		Order("start_time ASC").
		Find(&slots).Error
	
	return slots, err
}

// الإحصائيات والتقارير

func (r *servicesRepositoryImpl) GetServicesStats(ctx context.Context, userID string, timeframe string) (*models.ServicesStats, error) {
	stats := &models.ServicesStats{
		UserID:    userID,
		Timeframe: timeframe,
	}

	// حساب الإحصائيات الأساسية
	var totalServices int64
	if err := r.db.Model(&models.Service{}).Where("seller_id = ?", userID).Count(&totalServices).Error; err != nil {
		return nil, err
	}
	stats.TotalServices = int(totalServices)

	var activeServices int64
	if err := r.db.Model(&models.Service{}).Where("seller_id = ? AND status = ?", userID, "active").Count(&activeServices).Error; err != nil {
		return nil, err
	}
	stats.ActiveServices = int(activeServices)

	// حساب متوسط التقييم
	var avgRating float64
	if err := r.db.Model(&models.Service{}).Where("seller_id = ?", userID).Select("AVG(rating)").Row().Scan(&avgRating); err != nil {
		avgRating = 0
	}
	stats.AverageRating = avgRating

	// حساب إجمالي الطلبات
	var totalOrders int64
	if err := r.db.Model(&models.Service{}).Where("seller_id = ?", userID).Select("SUM(total_orders)").Row().Scan(&totalOrders); err != nil {
		totalOrders = 0
	}
	stats.TotalOrders = int(totalOrders)

	stats.CalculatedAt = time.Now()

	return stats, nil
}

func (r *servicesRepositoryImpl) GetPopularServices(ctx context.Context, limit int, category string) ([]models.Service, error) {
	var services []models.Service

	query := r.db.Where("status = ?", "active")
	if category != "" {
		query = query.Where("category = ?", category)
	}

	err := query.Order("total_orders DESC, rating DESC").
		Limit(limit).
		Find(&services).Error

	return services, err
}

func (r *servicesRepositoryImpl) GetServicesByCategory(ctx context.Context, category string, limit int) ([]models.Service, error) {
	var services []models.Service

	err := r.db.Where("category = ? AND status = ?", category, "active").
		Order("created_at DESC").
		Limit(limit).
		Find(&services).Error

	return services, err
}

func (r *servicesRepositoryImpl) GetServicesGrowth(ctx context.Context, sellerID string, period string) (*ServicesGrowth, error) {
	growth := &ServicesGrowth{
		SellerID: sellerID,
		Period:   period,
	}

	// حساب الخدمات الجديدة
	var newServices int64
	startDate := time.Now().AddDate(0, 0, -30) // آخر 30 يوم كمثال
	if err := r.db.Model(&models.Service{}).
		Where("seller_id = ? AND created_at >= ?", sellerID, startDate).
		Count(&newServices).Error; err != nil {
		return nil, err
	}
	growth.NewServices = int(newServices)

	growth.CalculatedAt = time.Now()
	return growth, nil
}

// الفئات والوسوم

func (r *servicesRepositoryImpl) GetAllCategories(ctx context.Context) ([]string, error) {
	var categories []string
	err := r.db.Model(&models.Service{}).
		Where("status = ?", "active").
		Distinct("category").
		Pluck("category", &categories).Error
	
	return categories, err
}

func (r *servicesRepositoryImpl) GetPopularTags(ctx context.Context, limit int) ([]string, error) {
	// هذا تنفيذ مبسط - في التطبيق الحقيقي قد تحتاج استعلام أكثر تعقيداً
	var tags []string
	err := r.db.Model(&models.Service{}).
		Where("status = ?", "active").
		Select("tags").
		Limit(limit).
		Find(&tags).Error
	
	return tags, err
}

func (r *servicesRepositoryImpl) GetServicesByTag(ctx context.Context, tag string, page, limit int) ([]models.Service, *utils.Pagination, error) {
	var services []models.Service
	var total int64

	query := r.db.Model(&models.Service{}).
		Where("status = ? AND tags LIKE ?", "active", "%"+tag+"%")

	// حساب العدد الإجمالي
	if err := query.Count(&total).Error; err != nil {
		return nil, nil, err
	}

	// تطبيق التقسيم
	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).
		Order("created_at DESC").
		Find(&services).Error; err != nil {
		return nil, nil, err
	}

	pagination := &utils.Pagination{
		Page:  page,
		Limit: limit,
		Total: int(total),
		Pages: (int(total) + limit - 1) / limit,
	}

	return services, pagination, nil
}

// إدارة المتجر

func (r *servicesRepositoryImpl) IncreaseServiceOrders(ctx context.Context, serviceID string) error {
	return r.db.Model(&models.Service{}).
		Where("id = ?", serviceID).
		Update("total_orders", gorm.Expr("total_orders + ?", 1)).Error
}

func (r *servicesRepositoryImpl) UpdateServiceRating(ctx context.Context, serviceID string) error {
	var avgRating float64
	err := r.db.Model(&models.Rating{}).
		Where("service_id = ?", serviceID).
		Select("AVG(rating)").
		Row().
		Scan(&avgRating)
	
	if err != nil {
		return err
	}

	// حساب عدد التقييمات
	var totalReviews int64
	if err := r.db.Model(&models.Rating{}).
		Where("service_id = ?", serviceID).
		Count(&totalReviews).Error; err != nil {
		return err
	}

	return r.db.Model(&models.Service{}).
		Where("id = ?", serviceID).
		Updates(map[string]interface{}{
			"rating":        avgRating,
			"total_reviews": totalReviews,
		}).Error
}

func (r *servicesRepositoryImpl) GetSimilarServices(ctx context.Context, service *models.Service, limit int) ([]models.Service, error) {
	var similarServices []models.Service

	// البحث عن خدمات مشابهة بناءً على الفئة والوسوم
	query := r.db.Where("category = ? AND id != ? AND status = ?", 
		service.Category, service.ID, "active")

	// إضافة شروط الوسوم المشتركة إذا وجدت
	if len(service.Tags) > 0 {
		tagConditions := []string{}
		for _, tag := range service.Tags {
			tagConditions = append(tagConditions, "tags LIKE '%"+tag+"%'")
		}
		query = query.Where(strings.Join(tagConditions, " OR "))
	}

	err := query.Order("rating DESC, total_orders DESC").
		Limit(limit).
		Find(&similarServices).Error

	return similarServices, err
}

// البحث المتقدم

func (r *servicesRepositoryImpl) AdvancedSearch(ctx context.Context, params AdvancedSearchParams) ([]models.Service, *utils.Pagination, error) {
	var services []models.Service
	var total int64

	query := r.db.Model(&models.Service{}).Where("status = ?", "active")

	// تطبيق شروط البحث
	if params.Query != "" {
		searchQuery := "%" + params.Query + "%"
		query = query.Where("title LIKE ? OR description LIKE ?", searchQuery, searchQuery)
	}

	if params.Category != "" {
		query = query.Where("category = ?", params.Category)
	}

	if len(params.Tags) > 0 {
		tagConditions := []string{}
		for _, tag := range params.Tags {
			tagConditions = append(tagConditions, "tags LIKE '%"+tag+"%'")
		}
		query = query.Where(strings.Join(tagConditions, " OR "))
	}

	if params.MinPrice > 0 {
		query = query.Where("price >= ?", params.MinPrice)
	}

	if params.MaxPrice > 0 {
		query = query.Where("price <= ?", params.MaxPrice)
	}

	if params.MinRating > 0 {
		query = query.Where("rating >= ?", params.MinRating)
	}

	if params.SellerID != "" {
		query = query.Where("seller_id = ?", params.SellerID)
	}

	// حساب العدد الإجمالي
	if err := query.Count(&total).Error; err != nil {
		return nil, nil, err
	}

	// تطبيق الترتيب
	if params.SortBy != "" {
		order := params.SortBy
		if params.SortOrder != "" {
			order = order + " " + params.SortOrder
		}
		query = query.Order(order)
	} else {
		query = query.Order("created_at DESC")
	}

	// تطبيق التقسيم
	offset := (params.Page - 1) * params.Limit
	if err := query.Offset(offset).Limit(params.Limit).Find(&services).Error; err != nil {
		return nil, nil, err
	}

	pagination := &utils.Pagination{
		Page:  params.Page,
		Limit: params.Limit,
		Total: int(total),
		Pages: (int(total) + params.Limit - 1) / params.Limit,
	}

	return services, pagination, nil
}

// الإدارة

func (r *servicesRepositoryImpl) GetAllServices(ctx context.Context, page, limit int, status string) ([]models.Service, *utils.Pagination, error) {
	var services []models.Service
	var total int64

	query := r.db.Model(&models.Service{})
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// حساب العدد الإجمالي
	if err := query.Count(&total).Error; err != nil {
		return nil, nil, err
	}

	// تطبيق التقسيم
	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).
		Order("created_at DESC").
		Find(&services).Error; err != nil {
		return nil, nil, err
	}

	pagination := &utils.Pagination{
		Page:  page,
		Limit: limit,
		Total: int(total),
		Pages: (int(total) + limit - 1) / limit,
	}

	return services, pagination, nil
}

func (r *servicesRepositoryImpl) CountServicesByStatus(ctx context.Context, sellerID string) (map[string]int, error) {
	results := make(map[string]int)
	
	var counts []struct {
		Status string
		Count  int
	}
	
	err := r.db.Model(&models.Service{}).
		Select("status, COUNT(*) as count").
		Where("seller_id = ?", sellerID).
		Group("status").
		Find(&counts).Error
	
	if err != nil {
		return nil, err
	}

	for _, count := range counts {
		results[count.Status] = count.Count
	}

	return results, nil
}