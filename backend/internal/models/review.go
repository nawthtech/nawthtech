package models

import (
	"time"
)

// Review نموذج التقييم العام
type Review struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	Type        string    `json:"type" gorm:"not null;index"` // service, seller, product
	TargetID    string    `json:"target_id" gorm:"not null;index"` // ID of the service, seller, etc.
	UserID      string    `json:"user_id" gorm:"not null;index"`
	UserName    string    `json:"user_name" gorm:"-"`
	UserAvatar  string    `json:"user_avatar,omitempty" gorm:"-"`
	Rating      int       `json:"rating" gorm:"not null;check:rating>=1 AND rating<=5"`
	Title       string    `json:"title,omitempty" gorm:"size:200"`
	Comment     string    `json:"comment" gorm:"type:text"`
	IsVerified  bool      `json:"is_verified" gorm:"default:false"`
	Helpful     int       `json:"helpful" gorm:"default:0"`
	Reported    bool      `json:"reported" gorm:"default:false"`
	ReportReason string   `json:"report_reason,omitempty"`
	Status      string    `json:"status" gorm:"default:'active';index"` // active, hidden, removed
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ReviewSummary ملخص التقييمات
type ReviewSummary struct {
	TargetID     string  `json:"target_id" gorm:"primaryKey"`
	Type         string  `json:"type" gorm:"primaryKey"` // service, seller
	TotalReviews int     `json:"total_reviews" gorm:"default:0"`
	AverageRating float64 `json:"average_rating" gorm:"default:0"`
	Rating1      int     `json:"rating_1" gorm:"default:0"` // 1 star
	Rating2      int     `json:"rating_2" gorm:"default:0"` // 2 stars
	Rating3      int     `json:"rating_3" gorm:"default:0"` // 3 stars
	Rating4      int     `json:"rating_4" gorm:"default:0"` // 4 stars
	Rating5      int     `json:"rating_5" gorm:"default:0"` // 5 stars
}

// ReviewHelpful التصويت على التقييم كمفيد
type ReviewHelpful struct {
	ID       string    `json:"id" gorm:"primaryKey"`
	ReviewID string    `json:"review_id" gorm:"not null;index"`
	UserID   string    `json:"user_id" gorm:"not null;index"`
	Helpful  bool      `json:"helpful" gorm:"not null"` // true for helpful, false for not helpful
	CreatedAt time.Time `json:"created_at"`
}

// ReviewReport بلاغ عن التقييم
type ReviewReport struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	ReviewID    string    `json:"review_id" gorm:"not null;index"`
	UserID      string    `json:"user_id" gorm:"not null;index"`
	Reason      string    `json:"reason" gorm:"not null"` // spam, inappropriate, false_info, other
	Description string    `json:"description,omitempty" gorm:"type:text"`
	Status      string    `json:"status" gorm:"default:'pending'"` // pending, reviewed, resolved
	ReviewedBy  string    `json:"reviewed_by,omitempty"`
	ReviewedAt  time.Time `json:"reviewed_at,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

// ReviewResponse رد على التقييم
type ReviewResponse struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	ReviewID  string    `json:"review_id" gorm:"not null;uniqueIndex"`
	UserID    string    `json:"user_id" gorm:"not null"` // typically the seller or admin
	Comment   string    `json:"comment" gorm:"type:text;not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ReviewCreateRequest طلب إنشاء تقييم
type ReviewCreateRequest struct {
	Type    string `json:"type" binding:"required,oneof=service seller"`
	TargetID string `json:"target_id" binding:"required"`
	Rating  int    `json:"rating" binding:"required,min=1,max=5"`
	Title   string `json:"title,omitempty" binding:"omitempty,min=5,max=200"`
	Comment string `json:"comment" binding:"required,min=10,max=1000"`
}

// ReviewUpdateRequest طلب تحديث تقييم
type ReviewUpdateRequest struct {
	Rating  int    `json:"rating,omitempty" binding:"omitempty,min=1,max=5"`
	Title   string `json:"title,omitempty" binding:"omitempty,min=5,max=200"`
	Comment string `json:"comment,omitempty" binding:"omitempty,min=10,max=1000"`
}

// ReviewReportRequest طلب بلاغ عن تقييم
type ReviewReportRequest struct {
	Reason      string `json:"reason" binding:"required,oneof=spam inappropriate false_info other"`
	Description string `json:"description,omitempty" binding:"omitempty,min=10,max=500"`
}

// ReviewResponseRequest طلب رد على تقييم
type ReviewResponseRequest struct {
	Comment string `json:"comment" binding:"required,min=10,max=1000"`
}