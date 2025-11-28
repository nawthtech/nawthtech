package models

import (
	"time"
)

// Service نموذج الخدمة الأساسي
type Service struct {
	ID           string    `json:"id" gorm:"primaryKey"`
	Title        string    `json:"title" gorm:"not null;index"`
	Description  string    `json:"description" gorm:"type:text"`
	Category     string    `json:"category" gorm:"not null;index"`
	Price        float64   `json:"price" gorm:"not null"`
	Duration     int       `json:"duration" gorm:"not null"` // بالمدة بالأيام
	Rating       float64   `json:"rating" gorm:"default:0;index"`
	TotalOrders  int       `json:"total_orders" gorm:"default:0"`
	TotalReviews int       `json:"total_reviews" gorm:"default:0"`
	Status       string    `json:"status" gorm:"default:'active';index"` // active, inactive, suspended, pending
	Featured     bool      `json:"featured" gorm:"default:false;index"`
	SellerID     string    `json:"seller_id" gorm:"not null;index"`
	SellerName   string    `json:"seller_name,omitempty" gorm:"-"`
	Images       []string  `json:"images" gorm:"type:json;serializer:json"`
	Features     []string  `json:"features" gorm:"type:json;serializer:json"`
	Tags         []string  `json:"tags" gorm:"type:json;serializer:json"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// ServiceDetails تفاصيل الخدمة الكاملة
type ServiceDetails struct {
	Service
	Seller          *User     `json:"seller,omitempty" gorm:"-"`
	Reviews         []Review  `json:"reviews,omitempty" gorm:"-"`
	AverageRating   float64   `json:"average_rating" gorm:"-"`
	SimilarServices []Service `json:"similar_services,omitempty" gorm:"-"`
}

// ServiceCreateRequest طلب إنشاء خدمة
type ServiceCreateRequest struct {
	Title       string   `json:"title" binding:"required,min=3,max=200"`
	Description string   `json:"description" binding:"required,min=10,max=2000"`
	Category    string   `json:"category" binding:"required"`
	Price       float64  `json:"price" binding:"required,min=0"`
	Duration    int      `json:"duration" binding:"required,min=1,max=365"`
	Images      []string `json:"images"`
	Features    []string `json:"features"`
	Tags        []string `json:"tags"`
}

// ServiceUpdateRequest طلب تحديث خدمة
type ServiceUpdateRequest struct {
	Title       string   `json:"title" binding:"omitempty,min=3,max=200"`
	Description string   `json:"description" binding:"omitempty,min=10,max=2000"`
	Category    string   `json:"category"`
	Price       float64  `json:"price" binding:"omitempty,min=0"`
	Duration    int      `json:"duration" binding:"omitempty,min=1,max=365"`
	Images      []string `json:"images"`
	Features    []string `json:"features"`
	Tags        []string `json:"tags"`
}

// ServiceStatusUpdateRequest طلب تحديث حالة الخدمة
type ServiceStatusUpdateRequest struct {
	Status string `json:"status" binding:"required,oneof=active inactive suspended pending"`
}

// ServiceSearchParams معاملات البحث في الخدمات
type ServiceSearchParams struct {
	Query     string   `json:"query"`
	Category  string   `json:"category"`
	Tags      []string `json:"tags"`
	MinPrice  float64  `json:"min_price"`
	MaxPrice  float64  `json:"max_price"`
	MinRating float64  `json:"min_rating"`
	SellerID  string   `json:"seller_id"`
	Status    string   `json:"status"`
	Featured  *bool    `json:"featured"`
	SortBy    string   `json:"sort_by"`
	SortOrder string   `json:"sort_order"`
	Page      int      `json:"page"`
	Limit     int      `json:"limit"`
}

// ServiceSearchResult نتيجة البحث في الخدمات
type ServiceSearchResult struct {
	Services   []Service `json:"services"`
	Total      int64     `json:"total"`
	Page       int       `json:"page"`
	Limit      int       `json:"limit"`
	TotalPages int       `json:"total_pages"`
	HasMore    bool      `json:"has_more"`
}

// ServiceAnalytics تحليلات الخدمات
type ServiceAnalytics struct {
	ServiceID      string    `json:"service_id"`
	Views          int       `json:"views"`
	Clicks         int       `json:"clicks"`
	Conversions    int       `json:"conversions"`
	ConversionRate float64   `json:"conversion_rate"`
	Revenue        float64   `json:"revenue"`
	Period         string    `json:"period"`
	Date           time.Time `json:"date"`
}

// Review نموذج التقييم (تم تغيير الاسم من ServiceReview لمنع التعارض)
type Review struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	ServiceID   string    `json:"service_id" gorm:"not null;index"`
	UserID      string    `json:"user_id" gorm:"not null;index"`
	UserName    string    `json:"user_name" gorm:"-"`
	UserAvatar  string    `json:"user_avatar,omitempty" gorm:"-"`
	Rating      int       `json:"rating" gorm:"not null;check:rating>=1 AND rating<=5"`
	Comment     string    `json:"comment" gorm:"type:text"`
	IsVerified  bool      `json:"is_verified" gorm:"default:false"`
	Helpful     int       `json:"helpful" gorm:"default:0"`
	Reported    bool      `json:"reported" gorm:"default:false"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ServiceReviewRequest طلب إضافة تقييم
type ServiceReviewRequest struct {
	Rating  int    `json:"rating" binding:"required,min=1,max=5"`
	Comment string `json:"comment" binding:"required,min=10,max=1000"`
}

// ServiceReviewUpdateRequest طلب تحديث تقييم
type ServiceReviewUpdateRequest struct {
	Rating  int    `json:"rating" binding:"omitempty,min=1,max=5"`
	Comment string `json:"comment" binding:"omitempty,min=10,max=1000"`
}

// ServiceCategory فئة الخدمة
type ServiceCategory struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"not null;uniqueIndex"`
	Description string    `json:"description" gorm:"type:text"`
	Icon        string    `json:"icon"`
	Color       string    `json:"color"`
	SortOrder   int       `json:"sort_order" gorm:"default:0"`
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	ServiceCount int      `json:"service_count,omitempty" gorm:"-"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ServiceCategoryRequest طلب إدارة الفئة
type ServiceCategoryRequest struct {
	Name        string `json:"name" binding:"required,min=2,max=100"`
	Description string `json:"description" binding:"required,min=10,max=500"`
	Icon        string `json:"icon"`
	Color       string `json:"color"`
	SortOrder   int    `json:"sort_order"`
	IsActive    bool   `json:"is_active"`
}

// ServiceTimeSlot فترات زمنية للخدمة
type ServiceTimeSlot struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	ServiceID   string    `json:"service_id" gorm:"not null;index"`
	Date        time.Time `json:"date" gorm:"not null;index"`
	StartTime   string    `json:"start_time" gorm:"not null"` // HH:MM
	EndTime     string    `json:"end_time" gorm:"not null"`   // HH:MM
	Available   bool      `json:"available" gorm:"default:true"`
	Booked      bool      `json:"booked" gorm:"default:false"`
	MaxSlots    int       `json:"max_slots" gorm:"default:1"`
	BookedSlots int       `json:"booked_slots" gorm:"default:0"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ServiceAvailability توفر الخدمة
type ServiceAvailability struct {
	ServiceID      string            `json:"service_id"`
	Date           time.Time         `json:"date"`
	AvailableSlots []ServiceTimeSlot `json:"available_slots"`
	IsAvailable    bool              `json:"is_available"`
	Message        string            `json:"message,omitempty"`
}

// ServiceBooking حجز الخدمة
type ServiceBooking struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	ServiceID   string    `json:"service_id" gorm:"not null;index"`
	UserID      string    `json:"user_id" gorm:"not null;index"`
	TimeSlotID  string    `json:"time_slot_id" gorm:"not null;index"`
	Date        time.Time `json:"date" gorm:"not null"`
	StartTime   string    `json:"start_time" gorm:"not null"`
	EndTime     string    `json:"end_time" gorm:"not null"`
	Status      string    `json:"status" gorm:"default:'pending'"` // pending, confirmed, completed, cancelled
	Notes       string    `json:"notes" gorm:"type:text"`
	TotalAmount float64   `json:"total_amount" gorm:"not null"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ServiceFavorite المفضلة
type ServiceFavorite struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	ServiceID string    `json:"service_id" gorm:"not null;index"`
	UserID    string    `json:"user_id" gorm:"not null;index"`
	CreatedAt time.Time `json:"created_at"`
}

// ServiceReport تقرير الخدمة
type ServiceReport struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	ServiceID   string    `json:"service_id" gorm:"not null;index"`
	UserID      string    `json:"user_id" gorm:"not null;index"`
	Reason      string    `json:"reason" gorm:"not null"`
	Description string    `json:"description" gorm:"type:text"`
	Status      string    `json:"status" gorm:"default:'pending'"` // pending, reviewed, resolved, dismissed
	ReviewedBy  string    `json:"reviewed_by,omitempty"`
	ReviewedAt  time.Time `json:"reviewed_at,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

// ServicePromotion ترويج الخدمة
type ServicePromotion struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	ServiceID   string    `json:"service_id" gorm:"not null;index"`
	Type        string    `json:"type" gorm:"not null"` // featured, spotlight, discount, banner
	Title       string    `json:"title"`
	Description string    `json:"description" gorm:"type:text"`
	Discount    float64   `json:"discount"` // نسبة مئوية
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ServiceRecommendation توصية الخدمة
type ServiceRecommendation struct {
	ServiceID   string  `json:"service_id"`
	Title       string  `json:"title"`
	Category    string  `json:"category"`
	Price       float64 `json:"price"`
	Rating      float64 `json:"rating"`
	TotalOrders int     `json:"total_orders"`
	Similarity  float64 `json:"similarity"` // درجة التشابه
	Reason      string  `json:"reason"`     // سبب التوصية
}

// ServiceTrends اتجاهات الخدمات
type ServiceTrends struct {
	Period          string         `json:"period"`
	TopCategories   []CategoryTrend `json:"top_categories"`
	PopularServices []Service       `json:"popular_services"`
	PriceTrends     []PriceTrend    `json:"price_trends"`
	SearchTrends    []SearchTrend   `json:"search_trends"`
}

// CategoryTrend اتجاه الفئة
type CategoryTrend struct {
	Category     string  `json:"category"`
	Growth       float64 `json:"growth"` // نسبة النمو
	ServiceCount int     `json:"service_count"`
	TotalOrders  int     `json:"total_orders"`
}

// PriceTrend اتجاه الأسعار
type PriceTrend struct {
	Date string  `json:"date"`
	Min  float64 `json:"min"`
	Max  float64 `json:"max"`
	Avg  float64 `json:"avg"`
}

// SearchTrend اتجاهات البحث
type SearchTrend struct {
	Query    string `json:"query"`
	Count    int    `json:"count"`
	Category string `json:"category"`
}

// ServiceExport تصدير بيانات الخدمات
type ServiceExport struct {
	Format    string                 `json:"format"` // csv, excel, json
	StartDate time.Time              `json:"start_date"`
	EndDate   time.Time              `json:"end_date"`
	Fields    []string               `json:"fields"`
	Filters   map[string]interface{} `json:"filters"`
}

// ================================
// النماذج المطلوبة لـ admin_service.go
// ================================

// ServiceStats إحصائيات الخدمات
type ServiceStats struct {
	TotalServices     int     `json:"total_services"`
	ActiveServices    int     `json:"active_services"`
	InactiveServices  int     `json:"inactive_services"`
	SuspendedServices int     `json:"suspended_services"`
	TotalRevenue      float64 `json:"total_revenue"`
	AverageRating     float64 `json:"average_rating"`
	TotalOrders       int     `json:"total_orders"`
	PopularCategory   string  `json:"popular_category"`
}

// ServicesReport تقرير الخدمات
type ServicesReport struct {
	Timeframe string                 `json:"timeframe"`
	Summary   map[string]interface{} `json:"summary"`
	Metrics   []map[string]interface{} `json:"metrics"`
}
