package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/nawthtech/nawthtech/backend/internal/logger"
	"github.com/nawthtech/nawthtech/backend/internal/services"

	"github.com/go-chi/chi/v5"
)

type PaymentHandler struct {
	paymentService *services.PaymentService
}

func NewPaymentHandler(paymentService *services.PaymentService) *PaymentHandler {
	return &PaymentHandler{
		paymentService: paymentService,
	}
}

// ==================== معالجة المدفوعات ====================

func (h *PaymentHandler) CreatePayment(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	
	var paymentData struct {
		OrderID string  `json:"orderId"`
		Amount  float64 `json:"amount"`
		Method  string  `json:"method"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&paymentData); err != nil {
		respondError(w, "بيانات غير صالحة", http.StatusBadRequest)
		return
	}

	logger.Stdout.Info("إنشاء دفعة جديدة", 
		"userID", userID, 
		"orderID", paymentData.OrderID, 
		"amount", paymentData.Amount, 
		"method", paymentData.Method)

	response := map[string]interface{}{
		"success": true,
		"message": "تم إنشاء الدفعة بنجاح",
		"data": map[string]interface{}{
			"id":         "pay_" + userID + "_" + paymentData.OrderID,
			"status":     "pending",
			"amount":     paymentData.Amount,
			"method":     paymentData.Method,
			"createdAt":  "2024-01-01T00:00:00Z",
		},
		"paymentId": "pay_" + userID + "_" + paymentData.OrderID,
		"nextSteps": []string{"complete_payment_processing"},
	}

	w.WriteHeader(http.StatusCreated)
	respondJSON(w, response)
}

func (h *PaymentHandler) ProcessPayment(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	
	var processData struct {
		PaymentID   string                 `json:"paymentId"`
		PaymentData map[string]interface{} `json:"paymentData"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&processData); err != nil {
		respondError(w, "بيانات غير صالحة", http.StatusBadRequest)
		return
	}

	logger.Stdout.Info("معالجة دفعة", 
		"userID", userID, 
		"paymentID", processData.PaymentID, 
		"method", processData.PaymentData["method"])

	response := map[string]interface{}{
		"success": true,
		"message": "تمت معالجة الدفعة بنجاح",
		"data": map[string]interface{}{
			"success":      true,
			"status":       "completed",
			"transactionId": "txn_" + processData.PaymentID,
		},
	}

	respondJSON(w, response)
}

func (h *PaymentHandler) GetPaymentById(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	paymentID := chi.URLParam(r, "paymentId")

	logger.Stdout.Info("جلب تفاصيل دفعة", "userID", userID, "paymentID", paymentID)

	response := map[string]interface{}{
		"success": true,
		"message": "تم جلب تفاصيل الدفعة بنجاح",
		"data": map[string]interface{}{
			"id":        paymentID,
			"status":    "completed",
			"amount":    150.00,
			"currency":  "SAR",
			"method":    "credit_card",
			"createdAt": "2024-01-01T00:00:00Z",
		},
	}

	respondJSON(w, response)
}

// ==================== طرق الدفع ====================

func (h *PaymentHandler) GetPaymentMethods(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)

	logger.Stdout.Info("جلب طرق الدفع المتاحة", "userID", userID)

	response := map[string]interface{}{
		"success": true,
		"message": "تم جلب طرق الدفع بنجاح",
		"data": []map[string]interface{}{
			{
				"id":          "method_1",
				"type":        "credit_card",
				"last4":       "4242",
				"brand":       "visa",
				"expiryMonth": 12,
				"expiryYear":  2025,
				"isDefault":   true,
			},
			{
				"id":        "method_2",
				"type":      "apple_pay",
				"isDefault": false,
			},
		},
	}

	respondJSON(w, response)
}

func (h *PaymentHandler) AddPaymentMethod(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	
	var methodData map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&methodData); err != nil {
		respondError(w, "بيانات غير صالحة", http.StatusBadRequest)
		return
	}

	logger.Stdout.Info("إضافة طريقة دفع جديدة", "userID", userID, "methodType", methodData["type"])

	response := map[string]interface{}{
		"success": true,
		"message": "تم إضافة طريقة الدفع بنجاح",
		"data": map[string]interface{}{
			"id":        "method_new",
			"type":      methodData["type"],
			"isDefault": false,
			"createdAt": "2024-01-01T00:00:00Z",
		},
	}

	w.WriteHeader(http.StatusCreated)
	respondJSON(w, response)
}

func (h *PaymentHandler) UpdatePaymentMethod(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	methodID := chi.URLParam(r, "methodId")

	logger.Stdout.Info("تحديث طريقة دفع", "userID", userID, "methodID", methodID)

	response := map[string]interface{}{
		"success": true,
		"message": "تم تحديث طريقة الدفع بنجاح",
		"data": map[string]interface{}{
			"id":        methodID,
			"updatedAt": "2024-01-01T00:00:00Z",
		},
	}

	respondJSON(w, response)
}

func (h *PaymentHandler) DeletePaymentMethod(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	methodID := chi.URLParam(r, "methodId")

	logger.Stdout.Info("حذف طريقة دفع", "userID", userID, "methodID", methodID)

	response := map[string]interface{}{
		"success": true,
		"message": "تم حذف طريقة الدفع بنجاح",
		"data": map[string]interface{}{
			"deleted": true,
		},
	}

	respondJSON(w, response)
}

func (h *PaymentHandler) GetDefaultPaymentMethod(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)

	logger.Stdout.Info("جلب طريقة الدفع الافتراضية", "userID", userID)

	response := map[string]interface{}{
		"success": true,
		"message": "تم جلب طريقة الدفع الافتراضية بنجاح",
		"data": map[string]interface{}{
			"id":          "method_1",
			"type":        "credit_card",
			"last4":       "4242",
			"brand":       "visa",
			"expiryMonth": 12,
			"expiryYear":  2025,
			"isDefault":   true,
		},
	}

	respondJSON(w, response)
}

func (h *PaymentHandler) SetDefaultPaymentMethod(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	methodID := chi.URLParam(r, "methodId")

	logger.Stdout.Info("تعيين طريقة دفع كافتراضية", "userID", userID, "methodID", methodID)

	response := map[string]interface{}{
		"success": true,
		"message": "تم تعيين طريقة الدفع الافتراضية بنجاح",
		"data": map[string]interface{}{
			"id":        methodID,
			"isDefault": true,
		},
	}

	respondJSON(w, response)
}

// ==================== الاستردادات ====================

func (h *PaymentHandler) RequestRefund(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	paymentID := chi.URLParam(r, "paymentId")
	
	var refundData struct {
		Amount float64 `json:"amount"`
		Reason string  `json:"reason"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&refundData); err != nil {
		respondError(w, "بيانات غير صالحة", http.StatusBadRequest)
		return
	}

	logger.Stdout.Info("طلب استرداد أموال", 
		"userID", userID, 
		"paymentID", paymentID, 
		"amount", refundData.Amount, 
		"reason", refundData.Reason)

	response := map[string]interface{}{
		"success": true,
		"message": "تم تقديم طلب الاسترداد بنجاح",
		"data": map[string]interface{}{
			"id":                   "ref_" + paymentID,
			"status":               "pending",
			"amount":               refundData.Amount,
			"reason":               refundData.Reason,
			"estimatedProcessing":  "5-7 أيام عمل",
		},
		"refundId":            "ref_" + paymentID,
		"estimatedProcessing": "5-7 أيام عمل",
	}

	respondJSON(w, response)
}

func (h *PaymentHandler) GetRefundById(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	refundID := chi.URLParam(r, "refundId")

	logger.Stdout.Info("جلب تفاصيل استرداد", "userID", userID, "refundID", refundID)

	response := map[string]interface{}{
		"success": true,
		"message": "تم جلب تفاصيل الاسترداد بنجاح",
		"data": map[string]interface{}{
			"id":     refundID,
			"status": "completed",
			"amount": 150.00,
			"reason": "طلب العميل",
		},
	}

	respondJSON(w, response)
}

// ==================== سجل المعاملات ====================

func (h *PaymentHandler) GetTransactions(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	
	query := r.URL.Query()
	page, _ := strconv.Atoi(query.Get("page"))
	if page == 0 {
		page = 1
	}
	limit, _ := strconv.Atoi(query.Get("limit"))
	if limit == 0 {
		limit = 20
	}

	logger.Stdout.Info("جلب سجل المعاملات", 
		"userID", userID, 
		"page", page, 
		"limit", limit)

	response := map[string]interface{}{
		"success": true,
		"message": "تم جلب سجل المعاملات بنجاح",
		"data": []map[string]interface{}{
			{
				"id":     "txn_1",
				"type":   "payment",
				"amount": 150.00,
				"status": "completed",
				"date":   "2024-01-01T00:00:00Z",
			},
		},
		"pagination": map[string]interface{}{
			"page":  page,
			"limit": limit,
			"total": 1,
		},
	}

	respondJSON(w, response)
}

func (h *PaymentHandler) GetTransactionById(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	transactionID := chi.URLParam(r, "transactionId")

	logger.Stdout.Info("جلب تفاصيل معاملة", "userID", userID, "transactionID", transactionID)

	response := map[string]interface{}{
		"success": true,
		"message": "تم جلب تفاصيل المعاملة بنجاح",
		"data": map[string]interface{}{
			"id":        transactionID,
			"type":      "payment",
			"amount":    150.00,
			"status":    "completed",
			"currency":  "SAR",
			"createdAt": "2024-01-01T00:00:00Z",
		},
	}

	respondJSON(w, response)
}

// ==================== webhooks ====================

func (h *PaymentHandler) HandleStripeWebhook(w http.ResponseWriter, r *http.Request) {
	var webhookData map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&webhookData); err != nil {
		respondError(w, "بيانات غير صالحة", http.StatusBadRequest)
		return
	}

	logger.Stdout.Info("معالجة Stripe webhook", 
		"eventType", webhookData["type"], 
		"webhookID", webhookData["id"])

	response := map[string]interface{}{
		"success": true,
		"message": "تم معالجة webhook بنجاح",
		"data": map[string]interface{}{
			"processed": true,
			"eventId":   webhookData["id"],
		},
	}

	respondJSON(w, response)
}

func (h *PaymentHandler) HandlePayPalWebhook(w http.ResponseWriter, r *http.Request) {
	var webhookData map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&webhookData); err != nil {
		respondError(w, "بيانات غير صالحة", http.StatusBadRequest)
		return
	}

	logger.Stdout.Info("معالجة PayPal webhook", 
		"eventType", webhookData["event_type"], 
		"webhookID", webhookData["id"])

	response := map[string]interface{}{
		"success": true,
		"message": "تم معالجة webhook بنجاح",
		"data": map[string]interface{}{
			"processed": true,
			"eventId":   webhookData["id"],
		},
	}

	respondJSON(w, response)
}

// ==================== الإحصائيات والتقارير ====================

func (h *PaymentHandler) GetPaymentStats(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	period := r.URL.Query().Get("period")
	if period == "" {
		period = "30d"
	}

	logger.Stdout.Info("جلب إحصائيات المدفوعات الشاملة", "userID", userID, "period", period)

	response := map[string]interface{}{
		"success": true,
		"message": "تم جلب إحصائيات المدفوعات بنجاح",
		"data": map[string]interface{}{
			"totalRevenue":   125430.00,
			"totalTransactions": 543,
			"successRate":    98.5,
			"averageOrderValue": 231.0,
		},
		"period": period,
	}

	respondJSON(w, response)
}

func (h *PaymentHandler) GetRevenueStats(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	period := r.URL.Query().Get("period")
	if period == "" {
		period = "30d"
	}
	groupBy := r.URL.Query().Get("groupBy")
	if groupBy == "" {
		groupBy = "day"
	}

	logger.Stdout.Info("جلب إحصائيات الإيرادات", "userID", userID, "period", period, "groupBy", groupBy)

	response := map[string]interface{}{
		"success": true,
		"message": "تم جلب إحصائيات الإيرادات بنجاح",
		"data": []map[string]interface{}{
			{
				"date":   "2024-01-01",
				"revenue": 4500.00,
				"orders":  23,
			},
		},
		"period":  period,
		"groupBy": groupBy,
	}

	respondJSON(w, response)
}

// ==================== إدارة المدفوعات للمسؤولين ====================

func (h *PaymentHandler) GetPendingPayments(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	
	query := r.URL.Query()
	page, _ := strconv.Atoi(query.Get("page"))
	if page == 0 {
		page = 1
	}
	limit, _ := strconv.Atoi(query.Get("limit"))
	if limit == 0 {
		limit = 20
	}

	logger.Stdout.Info("جلب المدفوعات قيد الانتظار للمسؤول", "userID", userID, "page", page, "limit", limit)

	response := map[string]interface{}{
		"success": true,
		"message": "تم جلب المدفوعات قيد الانتظار بنجاح",
		"data":    []map[string]interface{}{},
		"pagination": map[string]interface{}{
			"page":  page,
			"limit": limit,
			"total": 0,
		},
	}

	respondJSON(w, response)
}

func (h *PaymentHandler) ManualVerifyPayment(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	paymentID := chi.URLParam(r, "paymentId")
	
	var verifyData struct {
		Notes string `json:"notes"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&verifyData); err != nil {
		respondError(w, "بيانات غير صالحة", http.StatusBadRequest)
		return
	}

	logger.Stdout.Info("التحقق اليدوي من دفعة", "userID", userID, "paymentID", paymentID, "notes", verifyData.Notes)

	response := map[string]interface{}{
		"success": true,
		"message": "تم التحقق من الدفعة بنجاح",
		"data": map[string]interface{}{
			"verified": true,
			"verifiedBy": userID,
			"notes":    verifyData.Notes,
		},
	}

	respondJSON(w, response)
}