package handlers

import (
	"github.com/nawthtech/nawthtech/backend/internal/middleware"
	"github.com/nawthtech/nawthtech/backend/internal/services"

	"github.com/go-chi/chi/v5"
)

// RegisterPaymentRoutes تسجيل جميع مسارات المدفوعات
func RegisterPaymentRoutes(router chi.Router, paymentService *services.PaymentService) {
	paymentHandler := NewPaymentHandler(paymentService)

	router.Route("/api/v1/payments", func(r chi.Router) {
		// ==================== معالجة المدفوعات ====================
		r.Route("/", func(payments chi.Router) {
			payments.Use(middleware.AuthMiddleware)

			// إنشاء دفعة جديدة
			payments.With(middleware.RateLimiter(15, 10*60*1000)).Post("/", paymentHandler.CreatePayment)
			
			// معالجة دفعة
			payments.With(middleware.RateLimiter(10, 5*60*1000)).Post("/process", paymentHandler.ProcessPayment)
			
			// الحصول على تفاصيل دفعة
			payments.Get("/{paymentId}", paymentHandler.GetPaymentById)
		})

		// ==================== طرق الدفع ====================
		r.Route("/methods", func(methods chi.Router) {
			methods.Use(middleware.AuthMiddleware)

			// الحصول على طرق الدفع المتاحة
			methods.Get("/", paymentHandler.GetPaymentMethods)
			
			// إضافة طريقة دفع جديدة
			methods.Post("/", paymentHandler.AddPaymentMethod)
			
			// تحديث طريقة دفع
			methods.Put("/{methodId}", paymentHandler.UpdatePaymentMethod)
			
			// حذف طريقة دفع
			methods.Delete("/{methodId}", paymentHandler.DeletePaymentMethod)
			
			// الحصول على طريقة الدفع الافتراضية
			methods.Get("/default", paymentHandler.GetDefaultPaymentMethod)
			
			// تعيين طريقة دفع كافتراضية
			methods.Put("/{methodId}/default", paymentHandler.SetDefaultPaymentMethod)
		})

		// ==================== الاستردادات ====================
		r.Route("/refunds", func(refunds chi.Router) {
			refunds.Use(middleware.AuthMiddleware)

			// طلب استرداد أموال
			refunds.Post("/{paymentId}/refund", paymentHandler.RequestRefund)
			
			// الحصول على تفاصيل استرداد
			refunds.Get("/{refundId}", paymentHandler.GetRefundById)
		})

		// ==================== سجل المعاملات ====================
		r.Route("/transactions", func(transactions chi.Router) {
			transactions.Use(middleware.AuthMiddleware)

			// الحصول على سجل المعاملات
			transactions.Get("/", paymentHandler.GetTransactions)
			
			// الحصول على تفاصيل معاملة
			transactions.Get("/{transactionId}", paymentHandler.GetTransactionById)
		})

		// ==================== webhooks ====================
		r.Route("/webhooks", func(webhooks chi.Router) {
			// webhook لاستقبال تحديثات Stripe
			webhooks.Post("/stripe", paymentHandler.HandleStripeWebhook)
			
			// webhook لاستقبال تحديثات PayPal
			webhooks.Post("/paypal", paymentHandler.HandlePayPalWebhook)
		})

		// ==================== الإحصائيات والتقارير ====================
		r.Route("/stats", func(stats chi.Router) {
			stats.Use(middleware.AuthMiddleware, middleware.AdminAuth)

			// الحصول على إحصائيات المدفوعات الشاملة
			stats.Get("/overview", paymentHandler.GetPaymentStats)
			
			// الحصول على إحصائيات الإيرادات
			stats.Get("/revenue", paymentHandler.GetRevenueStats)
		})

		// ==================== إدارة المدفوعات للمسؤولين ====================
		r.Route("/admin", func(admin chi.Router) {
			admin.Use(middleware.AuthMiddleware, middleware.AdminAuth)

			// الحصول على المدفوعات قيد الانتظار
			admin.Get("/pending", paymentHandler.GetPendingPayments)
			
			// التحقق اليدوي من دفعة
			admin.Post("/{paymentId}/manual-verify", paymentHandler.ManualVerifyPayment)
		})
	})
}