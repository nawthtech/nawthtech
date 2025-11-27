package utils

// Pagination هيكل التقسيم
type Pagination struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
	Total int `json:"total"`
	Pages int `json:"pages"`
}

// GetOffset حساب الإزاحة
func (p *Pagination) GetOffset() int {
	return (p.Page - 1) * p.Limit
}

// NewPagination إنشاء تقسيم جديد
func NewPagination(page, limit int) *Pagination {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	return &Pagination{
		Page:  page,
		Limit: limit,
	}
}