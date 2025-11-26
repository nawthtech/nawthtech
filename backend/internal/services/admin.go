package services

import (
	"time"

	"github.com/nawthtech/nawthtech/backend/internal/models"
)

type AdminService struct {
	// يمكن إضافة حقول مثل قاعدة البيانات هنا
}

func NewAdminService() *AdminService {
	return &AdminService{}
}

// GetDashboardData يحصل على بيانات لوحة التحكم
func (s *AdminService) GetDashboardData(timeRange string) (*models.DashboardData, error) {
	// محاكاة البيانات - في الواقع ستأتي من قاعدة البيانات
	stats := models.DashboardStats{
		TotalUsers:     1250,
		TotalOrders:    543,
		TotalRevenue:   125430,
		ActiveServices: 28,
		PendingOrders:  12,
		SupportTickets: 8,
		ConversionRate: 4.2,
		BounceRate:     32.1,
		StoreVisits:    3450,
		NewCustomers:   89,
	}

	storeMetrics := models.StoreMetrics{
		TotalProducts:         45,
		LowStockItems:         3,
		StoreRevenue:          89450,
		StoreOrders:           432,
		AverageOrderValue:     207,
		TopSellingCategory:    "خدمات الوسائل الاجتماعية",
		CustomerSatisfaction:  4.8,
		ReturnRate:            1.2,
	}

	recentOrders := []models.Order{
		{
			ID:       "ORD-001",
			User:     "أحمد محمد",
			Service:  "متابعين إنستغرام - 1000 متابع",
			Amount:   150,
			Status:   "completed",
			Date:     "2024-01-15",
			Type:     "store",
			Category: "وسائل اجتماعية",
		},
		// ... إضافة باقي الطلبات
	}

	userActivity := []models.UserActivity{
		{
			User:    "أحمد محمد",
			Action:  "شراء من المتجر",
			Service: "متابعين إنستغرام",
			Time:    "منذ 5 دقائق",
			IP:      "192.168.1.100",
			Type:    "purchase",
		},
		// ... إضافة باقي النشاطات
	}

	return &models.DashboardData{
		Stats:        stats,
		StoreMetrics: storeMetrics,
		RecentOrders: recentOrders,
		UserActivity: userActivity,
		SystemAlerts: []models.SystemAlert{},
	}, nil
}

// باقي الدوال تبقى كما هي...
