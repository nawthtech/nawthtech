package services

// UserService يقدم خدمات إدارة المستخدمين
type UserService struct{}

func NewUserService() *UserService {
	return &UserService{}
}

// يمكن إضافة الدوال اللازمة هنا عند الحاجة