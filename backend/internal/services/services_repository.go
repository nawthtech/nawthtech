package services

import (
	"context"
	"fmt"
	"time"

	"github.com/nawthtech/nawthtech/backend/internal/models"
	"github.com/nawthtech/nawthtech/backend/internal/utils"
	"gorm.io/gorm"
)

// ServicesRepository واجهة مستودع الخدمات
type ServicesRepository interface {
	GetServices(ctx context.Context, params GetServicesParams) ([]models.Service, *utils.Pagination, error)
	SearchServices(ctx context.Context, params SearchServicesParams) ([]models.Service, *utils.Pagination, error)
	GetFeaturedServices(ctx context.Context, limit int) ([]models.Service, error)
	GetServiceByID(ctx context.Context, serviceID string) (*models.Service, error)
	GetRecommendedServices(ctx context.Context, serviceID string, limit int) ([]models.Service, error)
	GetSellerServices(ctx context.Context, params GetSellerServicesParams) ([]models.Service, *utils.Pagination, error)
	CheckAvailability(ctx context.Context, params CheckAvailabilityParams) (*models.Availability, error)
	AddRating(ctx context.Context, params AddRatingParams) (*models.Rating, error)
	CreateService(ctx context.Context, service *models.Service) (*models.Service, error)
	UpdateService(ctx context.Context, service *models.Service) (*models.Service, error)
	UpdateServiceStatus(ctx context.Context, serviceID string, status string) (*models.Service, error)
	DeleteService(ctx context.Context, serviceID string) error
	GetServicesStats(ctx context.Context, userID string, timeframe string) (*models.ServicesStats, error)
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

	err := r.db.Preload("Reviews").First(&service, "id = ?", serviceID).Error
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
	err := r.db.Where("category = ? AND id != ? AND status = ?", currentService.Category, serviceID, "active").
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

func (r *servicesRepositoryImpl) CheckAvailability(ctx context.Context, params CheckAvailabilityParams) (*models.Availability, error) {
	// هذا تنفيذ مبسط - في التطبيق الحقيقي قد تحتاج للتحقق من الجدول الزمني
	availability := &models.Availability{
		Available:      true,
		ServiceID:      params.ServiceID,
		Date:           params.Date,
		Time:           params.Time,
		Guests:         params.Guests,
		Message:        "الخدمة متاحة في الوقت المحدد",
		SuggestedTimes: []string{"10:00", "14:00", "16:00"},
		CheckedAt:      time.Now(),
	}

	return availability, nil
}

func (r *servicesRepositoryImpl) AddRating(ctx context.Context, params AddRatingParams) (*models.Rating, error) {
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
	var avgRating float64
	if err := r.db.Model(&models.Rating{}).
		Where("service_id = ?", params.ServiceID).
		Select("AVG(rating)").
		Row().
		Scan(&avgRating); err != nil {
		return nil, err
	}

	// تحديث الخدمة بمتوسط التقييم الجديد
	if err := r.db.Model(&models.Service{}).
		Where("id = ?", params.ServiceID).
		Update("rating", avgRating).Error; err != nil {
		return nil, err
	}

	return rating, nil
}

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