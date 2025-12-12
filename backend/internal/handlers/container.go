package handlers

import (
	"github.com/nawthtech/nawthtech/backend/internal/services"
)

// HandlerContainer حاوية لجميع الـ handlers
type HandlerContainer struct {
	Auth         AuthHandler
	User         UserHandler
	Service      ServiceHandler
	Category     CategoryHandler
	Order        OrderHandler
	Payment      PaymentHandler
	Upload       UploadHandler
	Notification NotificationHandler
	Admin        AdminHandler
	Health       HealthHandler
}

// NewHandlerContainer إنشاء حاوية handlers جديدة
func NewHandlerContainer(serviceContainer *services.ServiceContainer) *HandlerContainer {
	return &HandlerContainer{
		Auth:         NewAuthHandler(serviceContainer.Auth),
		User:         NewUserHandler(serviceContainer.User),
		Service:      NewServiceHandler(serviceContainer.Service),
		Category:     NewCategoryHandler(serviceContainer.Category),
		Order:        NewOrderHandler(serviceContainer.Order),
		Payment:      NewPaymentHandler(serviceContainer.Payment),
		Upload:       NewUploadHandler(serviceContainer.Upload),
		Notification: NewNotificationHandler(serviceContainer.Notification),
		Admin:        NewAdminHandler(serviceContainer.Admin),
		Health:       NewHealthHandler(serviceContainer.Health),
	}
}