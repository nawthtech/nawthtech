package services

import (
	"context"
	"fmt"
	"time"

	"github.com/nawthtech/nawthtech/backend/internal/models"
	"github.com/nawthtech/nawthtech/backend/internal/utils"
)

// ServicesService واجهة خدمة الخدمات
type ServicesService interface {
	GetServices(ctx context.Context, params GetServicesParams) ([]models.Service, *utils.Pagination, error)
	SearchServices(ctx context.Context, params SearchServicesParams) ([]models.Service, *utils.Pagination, error)
	GetFeaturedServices(ctx context.Context, limit int) ([]models.Service, error)
	GetServiceDetails(ctx context.Context, serviceID string) (*models.Service, error)
	GetRecommendedServices(ctx context.Context, serviceID string, limit int) ([]models.Service, error)
	GetSellerServices(ctx context.Context, params GetSellerServicesParams) ([]models.Service, *utils.Pagination, error)
	CheckAvailability(ctx context.Context, params CheckAvailabilityParams) (*models.Availability, error)
	AddRating(ctx context.Context, params AddRatingParams) (*models.Rating, error)
	CreateService(ctx context.Context, params CreateServiceParams) (*models.Service, error)
	UpdateService(ctx context.Context, params UpdateServiceParams) (*models.Service, error)
	UpdateServiceStatus(ctx context.Context, params UpdateServiceStatusParams) (*models.Service, error)
	DeleteService(ctx context.Context, serviceID string, userID string) error
	GetServicesStats(ctx context.Context, userID string, timeframe string) (*models.ServicesStats, error)
}

// GetServicesParams معاملات جلب الخدمات
type GetServicesParams struct {
	Page      int
	Limit     int
	Category  string
	SortBy    string
	SortOrder string
	MinPrice  float64
	MaxPrice  float64
}

// SearchServicesParams معاملات البحث في الخدمات
type SearchServicesParams struct {
	Query    string
	Page     int
	Limit    int
	Category string
}

// GetSellerServicesParams معاملات جلب خدمات البائع
type GetSellerServicesParams struct {
	SellerID string
	Page     int
	Limit    int
	Status   string
}

// CheckAvailabilityParams معاملات التحقق من التوفر
type CheckAvailabilityParams struct {
	ServiceID string
	Date      string
	Time      string
	Guests    int
}

// AddRatingParams معاملات إضافة تقييم
type AddRatingParams struct {
	ServiceID string
	UserID    string
	Rating    int
	Comment   string
}

// CreateServiceParams معاملات إنشاء خدمة
type CreateServiceParams struct {
	Title       string
	Description string
	Category    string
	Price       float64
	Duration    int
	Images      []string
	Features    []string
	Tags        []string
	SellerID    string
}

// UpdateServiceParams معاملات تحديث الخدمة
type UpdateServiceParams struct {
	ServiceID   string
	SellerID    string
	Title       string
	Description string
	Category    string
	Price       float64
	Duration    int
	Images      []string
	Features    []string
	Tags        []string
}

// UpdateServiceStatusParams معاملات تحديث حالة الخدمة
type UpdateServiceStatusParams struct {
	ServiceID string
	SellerID  string
	Status    string
}

// servicesServiceImpl التطبيق الفعلي لخدمة الخدمات
type servicesServiceImpl struct {
	repo ServicesRepository
}

// NewServicesService إنشاء خدمة خدمات جديدة
func NewServicesService(repo ServicesRepository) ServicesService {
	return &servicesServiceImpl{
		repo: repo,
	}
}


func (s *servicesServiceImpl) GetServices(ctx context.Context, params GetServicesParams) ([]models.Service, *utils.Pagination, error) {
	// التحقق من الصلاحيات والمعاملات
	if params.Page < 1 {
		params.Page = 1
	}
	if params.Limit < 1 || params.Limit > 100 {
		params.Limit = 10
	}

	return s.repo.GetServices(ctx, params)
}

func (s *servicesServiceImpl) SearchServices(ctx context.Context, params SearchServicesParams) ([]models.Service, *utils.Pagination, error) {
	if params.Query == "" {
		return nil, nil, fmt.Errorf("استعلام البحث مطلوب")
	}
	if params.Page < 1 {
		params.Page = 1
	}
	if params.Limit < 1 || params.Limit > 100 {
		params.Limit = 10
	}

	return s.repo.SearchServices(ctx, params)
}

func (s *servicesServiceImpl) GetFeaturedServices(ctx context.Context, limit int) ([]models.Service, error) {
	if limit < 1 || limit > 20 {
		limit = 10
	}

	return s.repo.GetFeaturedServices(ctx, limit)
}

func (s *servicesServiceImpl) GetServiceDetails(ctx context.Context, serviceID string) (*models.Service, error) {
	if serviceID == "" {
		return nil, fmt.Errorf("معرف الخدمة مطلوب")
	}

	return s.repo.GetServiceByID(ctx, serviceID)
}

func (s *servicesServiceImpl) GetRecommendedServices(ctx context.Context, serviceID string, limit int) ([]models.Service, error) {
	if serviceID == "" {
		return nil, fmt.Errorf("معرف الخدمة مطلوب")
	}
	if limit < 1 || limit > 20 {
		limit = 5
	}

	return s.repo.GetRecommendedServices(ctx, serviceID, limit)
}

func (s *servicesServiceImpl) GetSellerServices(ctx context.Context, params GetSellerServicesParams) ([]models.Service, *utils.Pagination, error) {
	if params.SellerID == "" {
		return nil, nil, fmt.Errorf("معرف البائع مطلوب")
	}
	if params.Page < 1 {
		params.Page = 1
	}
	if params.Limit < 1 || params.Limit > 100 {
		params.Limit = 10
	}

	return s.repo.GetSellerServices(ctx, params)
}

func (s *servicesServiceImpl) CheckAvailability(ctx context.Context, params CheckAvailabilityParams) (*models.Availability, error) {
	if params.ServiceID == "" {
		return nil, fmt.Errorf("معرف الخدمة مطلوب")
	}
	if params.Date == "" {
		return nil, fmt.Errorf("التاريخ مطلوب")
	}

	return s.repo.CheckAvailability(ctx, params)
}

func (s *servicesServiceImpl) AddRating(ctx context.Context, params AddRatingParams) (*models.Rating, error) {
	if params.ServiceID == "" {
		return nil, fmt.Errorf("معرف الخدمة مطلوب")
	}
	if params.UserID == "" {
		return nil, fmt.Errorf("معرف المستخدم مطلوب")
	}
	if params.Rating < 1 || params.Rating > 5 {
		return nil, fmt.Errorf("التقييم يجب أن يكون بين 1 و 5")
	}

	return s.repo.AddRating(ctx, params)
}

func (s *servicesServiceImpl) CreateService(ctx context.Context, params CreateServiceParams) (*models.Service, error) {
	// التحقق من البيانات المدخلة
	if params.Title == "" {
		return nil, fmt.Errorf("عنوان الخدمة مطلوب")
	}
	if params.Description == "" {
		return nil, fmt.Errorf("وصف الخدمة مطلوب")
	}
	if params.Category == "" {
		return nil, fmt.Errorf("فئة الخدمة مطلوبة")
	}
	if params.Price <= 0 {
		return nil, fmt.Errorf("سعر الخدمة يجب أن يكون أكبر من الصفر")
	}
	if params.SellerID == "" {
		return nil, fmt.Errorf("معرف البائع مطلوب")
	}

	service := &models.Service{
		ID:          fmt.Sprintf("service_%d", time.Now().UnixNano()),
		Title:       params.Title,
		Description: params.Description,
		Category:    params.Category,
		Price:       params.Price,
		Duration:    params.Duration,
		Images:      params.Images,
		Features:    params.Features,
		Tags:        params.Tags,
		SellerID:    params.SellerID,
		Status:      "active",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	return s.repo.CreateService(ctx, service)
}

func (s *servicesServiceImpl) UpdateService(ctx context.Context, params UpdateServiceParams) (*models.Service, error) {
	if params.ServiceID == "" {
		return nil, fmt.Errorf("معرف الخدمة مطلوب")
	}
	if params.SellerID == "" {
		return nil, fmt.Errorf("معرف البائع مطلوب")
	}

	// التحقق من ملكية الخدمة
	existingService, err := s.repo.GetServiceByID(ctx, params.ServiceID)
	if err != nil {
		return nil, err
	}
	if existingService.SellerID != params.SellerID {
		return nil, fmt.Errorf("غير مصرح بتحديث هذه الخدمة")
	}

	updateData := &models.Service{
		ID:          params.ServiceID,
		Title:       params.Title,
		Description: params.Description,
		Category:    params.Category,
		Price:       params.Price,
		Duration:    params.Duration,
		Images:      params.Images,
		Features:    params.Features,
		Tags:        params.Tags,
		UpdatedAt:   time.Now(),
	}

	return s.repo.UpdateService(ctx, updateData)
}

func (s *servicesServiceImpl) UpdateServiceStatus(ctx context.Context, params UpdateServiceStatusParams) (*models.Service, error) {
	if params.ServiceID == "" {
		return nil, fmt.Errorf("معرف الخدمة مطلوب")
	}
	if params.SellerID == "" {
		return nil, fmt.Errorf("معرف البائع مطلوب")
	}
	if params.Status != "active" && params.Status != "inactive" && params.Status != "suspended" {
		return nil, fmt.Errorf("حالة الخدمة غير صالحة")
	}

	// التحقق من ملكية الخدمة
	existingService, err := s.repo.GetServiceByID(ctx, params.ServiceID)
	if err != nil {
		return nil, err
	}
	if existingService.SellerID != params.SellerID {
		return nil, fmt.Errorf("غير مصرح بتحديث حالة هذه الخدمة")
	}

	return s.repo.UpdateServiceStatus(ctx, params.ServiceID, params.Status)
}

func (s *servicesServiceImpl) DeleteService(ctx context.Context, serviceID string, userID string) error {
	if serviceID == "" {
		return fmt.Errorf("معرف الخدمة مطلوب")
	}
	if userID == "" {
		return fmt.Errorf("معرف المستخدم مطلوب")
	}

	// التحقق من ملكية الخدمة
	existingService, err := s.repo.GetServiceByID(ctx, serviceID)
	if err != nil {
		return err
	}
	if existingService.SellerID != userID {
		return fmt.Errorf("غير مصرح بحذف هذه الخدمة")
	}

	return s.repo.DeleteService(ctx, serviceID)
}

func (s *servicesServiceImpl) GetServicesStats(ctx context.Context, userID string, timeframe string) (*models.ServicesStats, error) {
	if userID == "" {
		return nil, fmt.Errorf("معرف المستخدم مطلوب")
	}

	if timeframe == "" {
		timeframe = "month"
	}

	return s.repo.GetServicesStats(ctx, userID, timeframe)
}