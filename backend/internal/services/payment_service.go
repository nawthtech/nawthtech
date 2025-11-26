package services

// PaymentService يقدم خدمات إدارة المدفوعات
type PaymentService struct{}

func NewPaymentService() *PaymentService {
	return &PaymentService{}
}

// يمكن إضافة الدوال اللازمة هنا عند الحاجة