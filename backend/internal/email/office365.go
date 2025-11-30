package email

import (
	"crypto/tls"
	"fmt"
	"html/template"
	"os"
 "context"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/nawthtech/nawthtech/backend/internal/logger"
	"gopkg.in/gomail.v2"
)

// ================================
// هياكل البيانات
// ================================

// Office365Config إعدادات Office 365
type Office365Config struct {
	Host        string
	Port        int
	Username    string
	Password    string
	FromEmail   string
	FromName    string
	UseTLS      bool
	Timeout     time.Duration
	OutlookURL  string
	Enabled     bool
}

// EmailMessage رسالة البريد الإلكتروني
type EmailMessage struct {
	To          []string
	Cc          []string
	Bcc         []string
	Subject     string
	Body        string
	HTMLBody    string
	Attachments []string
	ReplyTo     string
	Priority    string // low, normal, high
}

// EmailTemplate قالب البريد الإلكتروني
type EmailTemplate struct {
	Name     string
	Subject  string
	HTMLPath string
}

// SendResult نتيجة إرسال البريد
type SendResult struct {
	Success bool
	Message string
	Error   error
}

// ================================
// دوال التهيئة
// ================================

// NewOffice365Config إنشاء إعدادات Office 365 جديدة
func NewOffice365Config() *Office365Config {
	return &Office365Config{
		Host:       getEnv("SMTP_HOST", "smtp.office365.com"),
		Port:       getEnvInt("SMTP_PORT", 587),
		Username:   getEnv("SMTP_USERNAME", ""),
		Password:   getEnv("SMTP_PASSWORD", ""),
		FromEmail:  getEnv("SMTP_FROM_EMAIL", ""),
		FromName:   getEnv("SMTP_FROM_NAME", "NawthTech"),
		UseTLS:     getEnvBool("SMTP_USE_TLS", true),
		Timeout:    time.Duration(getEnvInt("SMTP_TIMEOUT", 30)) * time.Second,
		OutlookURL: getEnv("OUTLOOK_WEB_URL", "https://outlook.office365.com"),
		Enabled:    getEnvBool("SMTP_ENABLED", true),
	}
}

// NewEmailService إنشاء خدمة بريد إلكتروني جديدة
func NewEmailService() (*Office365Config, error) {
	config := NewOffice365Config()

	if !config.Enabled {
		return config, nil
	}

	// التحقق من الإعدادات المطلوبة
	if config.Username == "" || config.Password == "" {
		return nil, fmt.Errorf("إعدادات البريد الإلكتروني غير مكتملة")
	}

	logger.Info(context.Background(), "✅ تم تهيئة خدمة البريد الإلكتروني مع Office 365",
		"host", config.Host,
		"port", config.Port,
		"from_email", config.FromEmail,
		"enabled", config.Enabled,
	)

	return config, nil
}

// ================================
// دوال إرسال البريد الأساسية
// ================================

// SendEmail إرسال بريد إلكتروني
func SendEmail(message *EmailMessage) *SendResult {
	config := NewOffice365Config()
	
	if !config.Enabled {
		return &SendResult{
			Success: false,
			Message: "خدمة البريد الإلكتروني غير مفعلة",
		}
	}

	startTime := time.Now()

	m := gomail.NewMessage()
	
	//设置发件人
	m.SetHeader("From", m.FormatAddress(config.FromEmail, config.FromName))
	
	//设置收件人
	m.SetHeader("To", message.To...)
	
	//设置抄送
	if len(message.Cc) > 0 {
		m.SetHeader("Cc", message.Cc...)
	}
	
	//设置密送
	if len(message.Bcc) > 0 {
		m.SetHeader("Bcc", message.Bcc...)
	}
	
	//设置主题
	m.SetHeader("Subject", message.Subject)
	
	//设置回复地址
	if message.ReplyTo != "" {
		m.SetHeader("Reply-To", message.ReplyTo)
	}
	
	//设置优先级
	if message.Priority != "" {
		m.SetHeader("X-Priority", getPriorityHeader(message.Priority))
	}
	
	//设置邮件正文
	if message.HTMLBody != "" {
		m.SetBody("text/html", message.HTMLBody)
		if message.Body != "" {
			m.AddAlternative("text/plain", message.Body)
		}
	} else {
		m.SetBody("text/plain", message.Body)
	}
	
	//添加附件
	for _, attachment := range message.Attachments {
		m.Attach(attachment)
	}

	//创建拨号器
	d := gomail.NewDialer(config.Host, config.Port, config.Username, config.Password)
	d.TLSConfig = &tls.Config{
		ServerName: config.Host,
	}

dialer := gomail.NewDialer(host, port, username, password)
// Timeout يتم تعيينه بشكل مختلف

	//发送邮件
	if err := d.DialAndSend(m); err != nil {
		logger.Error(context.Background(), "❌ فشل في إرسال البريد الإلكتروني",
			"to", strings.Join(message.To, ", "),
			"subject", message.Subject,
			"duration", time.Since(startTime),
			"error", err.Error(),
		)
		
		return &SendResult{
			Success: false,
			Message: "فشل في إرسال البريد الإلكتروني",
			Error:   err,
		}
	}

	logger.Info(context.Background(), "✅ تم إرسال البريد الإلكتروني بنجاح",
		"to", strings.Join(message.To, ", "),
		"subject", message.Subject,
		"duration", time.Since(startTime),
	)

	return &SendResult{
		Success: true,
		Message: "تم إرسال البريد الإلكتروني بنجاح",
	}
}

// SendSimpleEmail إرسال بريد إلكتروني بسيط
func SendSimpleEmail(to, subject, body string) error {
	message := &EmailMessage{
		To:      []string{to},
		Subject: subject,
		Body:    body,
	}

	result := SendEmail(message)
	return result.Error
}

// SendHTMLEmail إرسال بريد إلكتروني HTML
func SendHTMLEmail(to, subject, htmlBody string) error {
	message := &EmailMessage{
		To:       []string{to},
		Subject:  subject,
		HTMLBody: htmlBody,
	}

	result := SendEmail(message)
	return result.Error
}

// ================================
// دوال القوالب
// ================================

// LoadTemplate تحميل قالب البريد الإلكتروني
func LoadTemplate(templateName string, data interface{}) (string, string, error) {
	templatesDir := getEnv("EMAIL_TEMPLATES_DIR", "./templates/email")
	
	htmlPath := filepath.Join(templatesDir, templateName, "template.html")
	subjectPath := filepath.Join(templatesDir, templateName, "subject.txt")

	// تحميل موضوع البريد
	subject, err := loadTemplateFile(subjectPath, data)
	if err != nil {
		return "", "", fmt.Errorf("فشل في تحميل موضوع البريد: %v", err)
	}

	// تحميل نص HTML
	htmlBody, err := loadTemplateFile(htmlPath, data)
	if err != nil {
		return "", "", fmt.Errorf("فشل في تحميل قالب HTML: %v", err)
	}


	return subject, htmlBody, nil
}

// loadTemplateFile تحميل ملف قالب
func loadTemplateFile(filePath string, data interface{}) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	tmpl, err := template.New("email").Parse(string(content))
	if err != nil {
		return "", err
	}

	var result strings.Builder
	if err := tmpl.Execute(&result, data); err != nil {
		return "", err
	}

	return result.String(), nil
}

// ================================
// دوال البريد المحددة
// ================================

// SendWelcomeEmail إرسال بريد ترحيبي
func SendWelcomeEmail(to, name string) error {
	data := map[string]interface{}{
		"Name":          name,
		"AppName":       "NawthTech",
		"SupportEmail":  "support@nawthtech.com",
		"CurrentYear":   time.Now().Year(),
		"LoginURL":      getEnv("FRONTEND_URL", "https://nawthtech.com") + "/login",
	}

	subject, htmlBody, err := LoadTemplate("welcome", data)
	if err != nil {
		// استخدام قالب افتراضي إذا فشل تحميل القالب
		subject = "مرحباً بك في NawthTech - منصة الخدمات الإلكترونية"
		htmlBody = fmt.Sprintf(`
			<!DOCTYPE html>
			<html>
			<head>
				<meta charset="UTF-8">
				<title>مرحباً بك</title>
			</head>
			<body>
				<div style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto;">
					<h1 style="color: #2563eb;">مرحباً %s!</h1>
					<p>شكراً لك على الانضمام إلى <strong>NawthTech</strong> - منصة الخدمات الإلكترونية.</p>
					<p>يمكنك الآن الاستفادة من جميع خدماتنا المميزة.</p>
					<a href="%s" style="background-color: #2563eb; color: white; padding: 12px 24px; text-decoration: none; border-radius: 5px; display: inline-block;">
						بدء الاستخدام
					</a>
					<p>إذا كان لديك أي استفسار، لا تتردد في التواصل معنا على: support@nawthtech.com</p>
					<hr>
					<p style="color: #6b7280; font-size: 12px;">
						&copy; %d NawthTech. جميع الحقوق محفوظة.
					</p>
				</div>
			</body>
			</html>
		`, name, data["LoginURL"], data["CurrentYear"])
	}

	return SendHTMLEmail(to, subject, htmlBody)
}

// SendPasswordResetEmail إرسال بريد إعادة تعيين كلمة المرور
func SendPasswordResetEmail(to, name, resetToken string) error {
	resetURL := fmt.Sprintf("%s/reset-password?token=%s", 
		getEnv("FRONTEND_URL", "https://nawthtech.com"), 
		resetToken,
	)

	data := map[string]interface{}{
		"Name":      name,
		"ResetURL":  resetURL,
		"AppName":   "NawthTech",
		"ExpiresIn": "30 دقيقة", // يجب أن يتطابق مع PASSWORD_RESET_TOKEN_EXPIRY
	}

	subject, htmlBody, err := LoadTemplate("password_reset", data)
	if err != nil {
		subject = "إعادة تعيين كلمة المرور - NawthTech"
		htmlBody = fmt.Sprintf(`
			<!DOCTYPE html>
			<html>
			<head>
				<meta charset="UTF-8">
				<title>إعادة تعيين كلمة المرور</title>
			</head>
			<body>
				<div style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto;">
					<h1 style="color: #dc2626;">إعادة تعيين كلمة المرور</h1>
					<p>مرحباً %s,</p>
					<p>لقد طلبت إعادة تعيين كلمة المرور لحسابك في <strong>NawthTech</strong>.</p>
					<p>انقر على الزر أدناه لإعادة تعيين كلمة المرور:</p>
					<a href="%s" style="background-color: #dc2626; color: white; padding: 12px 24px; text-decoration: none; border-radius: 5px; display: inline-block;">
						إعادة تعيين كلمة المرور
					</a>
					<p style="color: #6b7280; font-size: 14px; margin-top: 20px;">
						<strong>ملاحظة:</strong> ستنتهي صلاحية هذا الرابط خلال %s.
					</p>
					<p>إذا لم تطلب إعادة تعيين كلمة المرور، يمكنك تجاهل هذا البريد.</p>
					<hr>
					<p style="color: #6b7280; font-size: 12px;">
						&copy; %d NawthTech. جميع الحقوق محفوظة.
					</p>
				</div>
			</body>
			</html>
		`, name, resetURL, data["ExpiresIn"], time.Now().Year())
	}

	return SendHTMLEmail(to, subject, htmlBody)
}

// SendVerificationEmail إرسال بريد التحقق
func SendVerificationEmail(to, name, verificationToken string) error {
	verifyURL := fmt.Sprintf("%s/verify-email?token=%s", 
		getEnv("FRONTEND_URL", "https://nawthtech.com"), 
		verificationToken,
	)

	data := map[string]interface{}{
		"Name":       name,
		"VerifyURL":  verifyURL,
		"AppName":    "NawthTech",
		"ExpiresIn":  "60 دقيقة", // يجب أن يتطابق مع EMAIL_VERIFICATION_TOKEN_EXPIRY
	}

	subject, htmlBody, err := LoadTemplate("email_verification", data)
	if err != nil {
		subject = "تحقق من بريدك الإلكتروني - NawthTech"
		htmlBody = fmt.Sprintf(`
			<!DOCTYPE html>
			<html>
			<head>
				<meta charset="UTF-8">
				<title>تحقق من البريد الإلكتروني</title>
			</head>
			<body>
				<div style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto;">
					<h1 style="color: #059669;">تحقق من بريدك الإلكتروني</h1>
					<p>مرحباً %s,</p>
					<p>شكراً لك على التسجيل في <strong>NawthTech</strong>!</p>
					<p>لتفعيل حسابك، يرجى النقر على الزر أدناه للتحقق من بريدك الإلكتروني:</p>
					<a href="%s" style="background-color: #059669; color: white; padding: 12px 24px; text-decoration: none; border-radius: 5px; display: inline-block;">
						تحقق من البريد الإلكتروني
					</a>
					<p style="color: #6b7280; font-size: 14px; margin-top: 20px;">
						<strong>ملاحظة:</strong> ستنتهي صلاحية هذا الرابط خلال %s.
					</p>
					<hr>
					<p style="color: #6b7280; font-size: 12px;">
						&copy; %d NawthTech. جميع الحقوق محفوظة.
					</p>
				</div>
			</body>
			</html>
		`, name, verifyURL, data["ExpiresIn"], time.Now().Year())
	}

	return SendHTMLEmail(to, subject, htmlBody)
}

// SendOrderConfirmationEmail إرسال بريد تأكيد الطلب
func SendOrderConfirmationEmail(to, name, orderID string, amount float64, serviceName string) error {
	data := map[string]interface{}{
		"Name":        name,
		"OrderID":     orderID,
		"Amount":      fmt.Sprintf("%.2f ر.س", amount),
		"ServiceName": serviceName,
		"AppName":     "NawthTech",
		"OrderDate":   time.Now().Format("2006-01-02 15:04"),
	}

	subject, htmlBody, err := LoadTemplate("order_confirmation", data)
	if err != nil {
		subject = "تم تأكيد طلبك - NawthTech"
		htmlBody = fmt.Sprintf(`
			<!DOCTYPE html>
			<html>
			<head>
				<meta charset="UTF-8">
				<title>تأكيد الطلب</title>
			</head>
			<body>
				<div style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto;">
					<h1 style="color: #059669;">تم تأكيد طلبك بنجاح!</h1>
					<p>مرحباً %s,</p>
					<p>شكراً لك على طلبك في <strong>NawthTech</strong>.</p>
					
					<div style="background-color: #f3f4f6; padding: 20px; border-radius: 5px; margin: 20px 0;">
						<h3 style="margin-top: 0;">تفاصيل الطلب:</h3>
						<p><strong>رقم الطلب:</strong> %s</p>
						<p><strong>الخدمة:</strong> %s</p>
						<p><strong>المبلغ:</strong> %s</p>
						<p><strong>تاريخ الطلب:</strong> %s</p>
					</div>
					
					<p>سنقوم بتحديثك على حالة طلبك عبر البريد الإلكتروني.</p>
					<p>إذا كان لديك أي استفسار، لا تتردد في التواصل معنا.</p>
					
					<hr>
					<p style="color: #6b7280; font-size: 12px;">
						&copy; %d NawthTech. جميع الحقوق محفوظة.
					</p>
				</div>
			</body>
			</html>
		`, name, orderID, serviceName, data["Amount"], data["OrderDate"], time.Now().Year())
	}

	return SendHTMLEmail(to, subject, htmlBody)
}

// ================================
// دوال المساعدة
// ================================

// getPriorityHeader الحصول على رأس الأولوية
func getPriorityHeader(priority string) string {
	switch strings.ToLower(priority) {
	case "high":
		return "1"
	case "low":
		return "5"
	default:
		return "3" // normal
	}
}

// HealthCheck فحص صحة خدمة البريد الإلكتروني
func HealthCheck() map[string]interface{} {
	config := NewOffice365Config()
	
	if !config.Enabled {
		return map[string]interface{}{
			"service": "email",
			"status":  "disabled",
			"enabled": false,
		}
	}

	// اختبار الاتصال بمخدم SMTP
	d := gomail.NewDialer(config.Host, config.Port, config.Username, config.Password)
	d.TLSConfig = &tls.Config{ServerName: config.Host}

	conn, err := d.Dial()
	if err != nil {
		return map[string]interface{}{
			"service": "email",
			"status":  "error",
			"enabled": true,
			"error":   err.Error(),
		}
	}
	defer conn.Close()

	return map[string]interface{}{
		"service":   "email",
		"status":    "healthy",
		"enabled":   true,
		"host":      config.Host,
		"port":      config.Port,
		"from_email": config.FromEmail,
	}
}

// IsEnabled التحقق إذا كانت خدمة البريد الإلكتروني مفعلة
func IsEnabled() bool {
	config := NewOffice365Config()
	return config.Enabled
}

// ================================
// دوال مساعدة للبيئة
// ================================

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}