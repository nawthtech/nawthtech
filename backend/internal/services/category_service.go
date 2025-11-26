package services

// CategoryService يقدم خدمات إدارة الفئات
type CategoryService struct{}

func NewCategoryService() *CategoryService {
	return &CategoryService{}
}

// يمكن إضافة الدوال اللازمة هنا عند الحاجة