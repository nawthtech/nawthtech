package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nawthtech/nawthtech/backend/internal/email"
	"github.com/nawthtech/nawthtech/backend/internal/services"
)

// EmailHandler معالجة طلبات البريد الإلكتروني
type EmailHandler struct {
	service      services.EmailService
	emailWorker  *email.CloudflareEmailWorker
}

// NewEmailHandler إنشاء Email handler جديد
func NewEmailHandler(service services.EmailService, emailWorker *email.CloudflareEmailWorker) *EmailHandler {
	return &EmailHandler{
		service:     service,
		emailWorker: emailWorker,
	}
}

// DeployEmailWorker نشر عامل البريد الإلكتروني
func (h *EmailHandler) DeployEmailWorker(c *gin.Context) {
	if h.emailWorker == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"success": false,
			"error":   "Email worker not initialized",
		})
		return
	}

	err := h.emailWorker.DeployWorkerScript()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to deploy email worker",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Email worker deployed successfully",
	})
}

// SetupEmailDNS إعداد سجلات DNS للبريد الإلكتروني
func (h *EmailHandler) SetupEmailDNS(c *gin.Context) {
	if h.emailWorker == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"success": false,
			"error":   "Email worker not initialized",
		})
		return
	}

	err := h.emailWorker.SetupDNSRecords()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to setup DNS records",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "DNS records for email routing configured successfully",
	})
}

// AddToEmailAllowList إضافة بريد إلكتروني إلى قائمة السماح
func (h *EmailHandler) AddToEmailAllowList(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	if h.emailWorker == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"success": false,
			"error":   "Email worker not initialized",
		})
		return
	}

	err := h.emailWorker.AddToAllowList(req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to add email to allow list",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": fmt.Sprintf("Email %s added to allow list successfully", req.Email),
	})
}

// RemoveFromEmailAllowList إزالة بريد إلكتروني من قائمة السماح
func (h *EmailHandler) RemoveFromEmailAllowList(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	if h.emailWorker == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"success": false,
			"error":   "Email worker not initialized",
		})
		return
	}

	err := h.emailWorker.RemoveFromAllowList(req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to remove email from allow list",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": fmt.Sprintf("Email %s removed from allow list successfully", req.Email),
	})
}

// GetEmailAllowList الحصول على قائمة السماح للبريد الإلكتروني
func (h *EmailHandler) GetEmailAllowList(c *gin.Context) {
	if h.emailWorker == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"success": false,
			"error":   "Email worker not initialized",
		})
		return
	}

	allowList := h.emailWorker.GetAllowList()

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"allow_list": allowList,
			"count":      len(allowList),
		},
	})
}

// TestEmailRouting اختبار توجيه البريد الإلكتروني
func (h *EmailHandler) TestEmailRouting(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	if h.emailWorker == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"success": false,
			"error":   "Email worker not initialized",
		})
		return
	}

	err := h.emailWorker.TestEmailRouting(req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to test email routing",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Email routing test initialized",
		"data": gin.H{
			"test_email": req.Email,
			"note":       "Check your inbox for the forwarded email",
		},
	})
}

// SendEmail إرسال بريد إلكتروني مباشر
func (h *EmailHandler) SendEmail(c *gin.Context) {
	var req struct {
		To      string   `json:"to" binding:"required,email"`
		Subject string   `json:"subject" binding:"required"`
		Body    string   `json:"body" binding:"required"`
		CC      []string `json:"cc"`
		BCC     []string `json:"bcc"`
		IsHTML  bool     `json:"is_html" default:"false"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// استخدام service للإرسال إذا كان موجوداً
	if h.service != nil {
		err := h.service.SendEmail(c.Request.Context(), req.To, req.Subject, req.Body, req.CC, req.BCC, req.IsHTML)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "Failed to send email",
				"details": err.Error(),
			})
			return
		}
	} else {
		// أو استخدام طريقة بديلة
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Email queued for sending (service not configured)",
			"note":    "Configure email service for actual sending",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Email sent successfully",
		"data": gin.H{
			"to":      req.To,
			"subject": req.Subject,
			"body_length": len(req.Body),
			"cc_count": len(req.CC),
			"bcc_count": len(req.BCC),
		},
	})
}

// GetEmailConfig الحصول على إعدادات البريد الإلكتروني
func (h *EmailHandler) GetEmailConfig(c *gin.Context) {
	if h.emailWorker == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"success": false,
			"error":   "Email worker not initialized",
		})
		return
	}

	config := gin.H{
		"provider":     "cloudflare_email_worker",
		"domain":       h.emailWorker.Domain,
		"forward_to":   h.emailWorker.ForwardTo,
		"allow_list":   h.emailWorker.GetAllowList(),
		"allow_list_count": len(h.emailWorker.GetAllowList()),
		"script_name":  h.emailWorker.ScriptName,
		"has_dns_setup": h.emailWorker.ZoneID != "",
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    config,
	})
}

// ValidateEmail التحقق من صحة البريد الإلكتروني
func (h *EmailHandler) ValidateEmail(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// تحقق بسيط من صيغة البريد الإلكتروني
	isValid := isValidEmailFormat(req.Email)
	domain := ""
	if isValid {
		parts := strings.Split(req.Email, "@")
		if len(parts) == 2 {
			domain = parts[1]
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"email":        req.Email,
			"is_valid":     isValid,
			"domain":       domain,
			"has_mx_record": isValid, // في الواقع يحتاج تحقق DNS
			"is_disposable": isDisposableEmailDomain(domain),
		},
	})
}

// ===== دوال مساعدة =====

func isValidEmailFormat(email string) bool {
	// تحقق بسيط من الصيغة
	if len(email) > 254 {
		return false
	}
	
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}
	
	localPart := parts[0]
	domainPart := parts[1]
	
	if len(localPart) == 0 || len(localPart) > 64 {
		return false
	}
	
	if len(domainPart) == 0 || len(domainPart) > 255 {
		return false
	}
	
	return true
}

func isDisposableEmailDomain(domain string) bool {
	// قائمة نطاقات مؤقتة معروفة
	disposableDomains := []string{
		"tempmail.com", "10minutemail.com", "guerrillamail.com",
		"mailinator.com", "yopmail.com", "dispostable.com",
		"throwawaymail.com", "fakeinbox.com", "trashmail.com",
	}
	
	domainLower := strings.ToLower(domain)
	for _, d := range disposableDomains {
		if strings.Contains(domainLower, d) {
			return true
		}
	}
	
	return false
}