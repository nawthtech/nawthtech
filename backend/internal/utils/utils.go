package utils

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"math/big"
	"net/http"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/gin-gonic/gin"
	"github.com/nawthtech/nawthtech/backend/internal/logger"
	"golang.org/x/crypto/bcrypt"
)

// ========== هياكل البيانات ==========

// Pagination هيكل الترقيم
type Pagination struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
	HasNext    bool  `json:"has_next"`
	HasPrev    bool  `json:"has_prev"`
}

// Response هيكل الاستجابة الموحد
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

// MemoryStats إحصائيات الذاكرة
type MemoryStats struct {
	UsedMB          float64 `json:"used_mb"`
	TotalMB         float64 `json:"total_mb"`
	UsagePercentage float64 `json:"usage_percentage"`
}

// ValidationError خطأ التحقق
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ========== دوال الاستجابة ==========

// SuccessResponse إرسال استجابة ناجحة
func SuccessResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	response := Response{
		Success: true,
		Message: message,
		Data:    data,
	}

	// إضافة metadata إذا كان data يحتوي على pagination
	if pagination, ok := extractPagination(data); ok {
		response.Meta = gin.H{"pagination": pagination}
	}

	c.JSON(statusCode, response)
}

// ErrorResponse إرسال استجابة خطأ
func ErrorResponse(c *gin.Context, statusCode int, message string, errorCode string) {
	c.JSON(statusCode, Response{
		Success: false,
		Message: message,
		Error:   errorCode,
	})
}

// ValidationErrorResponse إرسال استجابة أخطاء تحقق
func ValidationErrorResponse(c *gin.Context, errors []ValidationError) {
	c.JSON(http.StatusBadRequest, Response{
		Success: false,
		Message: "أخطاء في التحقق من البيانات",
		Error:   "VALIDATION_ERROR",
		Data:    errors,
	})
}

// extractPagination استخراج معلومات الترقيم من البيانات
func extractPagination(data interface{}) (*Pagination, bool) {
	if data == nil {
		return nil, false
	}

	val := reflect.ValueOf(data)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() == reflect.Struct {
		paginationField := val.FieldByName("Pagination")
		if paginationField.IsValid() && paginationField.Type() == reflect.TypeOf(&Pagination{}) {
			if pagination, ok := paginationField.Interface().(*Pagination); ok {
				return pagination, true
			}
		}
	}

	return nil, false
}

// ========== دوال الترقيم ==========

// NewPagination إنشاء كائن ترقيم جديد
func NewPagination(page, limit int, total int64) *Pagination {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	hasNext := page < totalPages
	hasPrev := page > 1

	return &Pagination{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
		HasNext:    hasNext,
		HasPrev:    hasPrev,
	}
}

// GetPaginationParams الحصول على معاملات الترقيم من الطلب
func GetPaginationParams(c *gin.Context) (page, limit int) {
	page, _ = strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ = strconv.Atoi(c.DefaultQuery("limit", "20"))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	return page, limit
}

// ApplyPagination تطبيق الترقيم على الاستعلام
func ApplyPagination(query interface{}, page, limit int) (interface{}, *Pagination) {
	// هذه دالة عامة - يمكن تخصيصها حسب ORM المستخدم
	// في GORM يمكن استخدام: db.Offset((page - 1) * limit).Limit(limit)
	offset := (page - 1) * limit
	return map[string]interface{}{
			"query":  query,
			"offset": offset,
			"limit":  limit,
		}, &Pagination{
			Page:  page,
			Limit: limit,
		}
}

// ========== دوال التحقق والتحقق من الصحة ==========

// IsEmpty التحقق إذا كانت القيمة فارغة
func IsEmpty(value interface{}) bool {
	if value == nil {
		return true
	}

	switch v := value.(type) {
	case string:
		return strings.TrimSpace(v) == ""
	case int, int8, int16, int32, int64:
		return v == 0
	case uint, uint8, uint16, uint32, uint64:
		return v == 0
	case float32, float64:
		return v == 0
	case bool:
		return !v
	case []interface{}:
		return len(v) == 0
	case map[string]interface{}:
		return len(v) == 0
	default:
		return reflect.ValueOf(v).IsZero()
	}
}

// IsValidEmail التحقق من صحة البريد الإلكتروني
func IsValidEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}

// IsValidPhone التحقق من صحة رقم الهاتف
func IsValidPhone(phone string) bool {
	// نموذج مبسط لرقم الهاتف - يمكن تعديله حسب الاحتياج
	pattern := `^[\+]?[0-9]{10,15}$`
	matched, _ := regexp.MatchString(pattern, phone)
	return matched
}

// IsStrongPassword التحقق من قوة كلمة المرور
func IsStrongPassword(password string) bool {
	if len(password) < 8 {
		return false
	}

	var (
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	return hasUpper && hasLower && hasNumber && hasSpecial
}

// ValidateStruct التحقق من صحة الهيكل
func ValidateStruct(s interface{}) []ValidationError {
	var errors []ValidationError
	val := reflect.ValueOf(s).Elem()
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		// الحصول على tag التحقق
		validateTag := fieldType.Tag.Get("validate")
		if validateTag == "" {
			continue
		}

		rules := strings.Split(validateTag, ",")
		fieldName := getFieldName(fieldType)

		for _, rule := range rules {
			switch rule {
			case "required":
				if IsEmpty(field.Interface()) {
					errors = append(errors, ValidationError{
						Field:   fieldName,
						Message: "هذا الحقل مطلوب",
					})
				}
			case "email":
				if !IsEmpty(field.Interface()) && !IsValidEmail(field.String()) {
					errors = append(errors, ValidationError{
						Field:   fieldName,
						Message: "البريد الإلكتروني غير صالح",
					})
				}
			case "phone":
				if !IsEmpty(field.Interface()) && !IsValidPhone(field.String()) {
					errors = append(errors, ValidationError{
						Field:   fieldName,
						Message: "رقم الهاتف غير صالح",
					})
				}
			}
		}
	}

	return errors
}

// getFieldName الحصول على اسم الحقل للعرض
func getFieldName(field reflect.StructField) string {
	jsonTag := field.Tag.Get("json")
	if jsonTag != "" {
		parts := strings.Split(jsonTag, ",")
		if parts[0] != "" {
			return parts[0]
		}
	}
	return field.Name
}

// ========== دوال التشفير والأمان ==========

// HashPassword تشفير كلمة المرور
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPassword التحقق من كلمة المرور
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateRandomString إنشاء سلسلة عشوائية
func GenerateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// GenerateAPIKey إنشاء مفتاح API
func GenerateAPIKey() (string, error) {
	return GenerateRandomString(32)
}

// GenerateSecureToken إنشاء توكن آمن
func GenerateSecureToken() (string, error) {
	return GenerateRandomString(64)
}

// ========== دوال الوقت والتاريخ ==========

// Now الحصول على الوقت الحالي
func Now() time.Time {
	return time.Now().UTC()
}

// FormatTime تنسيق الوقت
func FormatTime(t time.Time, format string) string {
	if format == "" {
		format = time.RFC3339
	}
	return t.Format(format)
}

// ParseTime تحليل الوقت من سلسلة
func ParseTime(timeStr, format string) (time.Time, error) {
	if format == "" {
		format = time.RFC3339
	}
	return time.Parse(format, timeStr)
}

// IsExpired التحقق إذا انتهت الصلاحية
func IsExpired(expiryTime time.Time) bool {
	return Now().After(expiryTime)
}

// CalculateExpiryTime حساب وقت انتهاء الصلاحية
func CalculateExpiryTime(duration time.Duration) time.Time {
	return Now().Add(duration)
}

// HumanizeDuration تحويل المدة إلى نص مقروء
func HumanizeDuration(duration time.Duration) string {
	if duration < time.Minute {
		return "أقل من دقيقة"
	}

	if duration < time.Hour {
		minutes := int(duration.Minutes())
		return fmt.Sprintf("%d دقيقة", minutes)
	}

	if duration < 24*time.Hour {
		hours := int(duration.Hours())
		return fmt.Sprintf("%d ساعة", hours)
	}

	days := int(duration.Hours() / 24)
	return fmt.Sprintf("%d يوم", days)
}

// ========== دوال JSON ==========

// ToJSON تحويل البيانات إلى JSON
func ToJSON(data interface{}) (string, error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// FromJSON تحويل JSON إلى بيانات
func FromJSON(jsonStr string, result interface{}) error {
	return json.Unmarshal([]byte(jsonStr), result)
}

// PrettyJSON تحويل البيانات إلى JSON منسق
func PrettyJSON(data interface{}) (string, error) {
	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// ========== دوال السياق ==========

// GetRequestIDFromContext الحصول على معرف الطلب من السياق
func GetRequestIDFromContext(ctx context.Context) string {
	if reqID, ok := ctx.Value("requestID").(string); ok {
		return reqID
	}
	return ""
}

// GetUserIDFromContext الحصول على معرف المستخدم من السياق
func GetUserIDFromContext(ctx context.Context) string {
	if userID, ok := ctx.Value("userID").(string); ok {
		return userID
	}
	return ""
}

// GetUserIDFromGinContext الحصول على معرف المستخدم من سياق Gin
func GetUserIDFromGinContext(c *gin.Context) string {
	if userID, exists := c.Get("userID"); exists {
		if id, ok := userID.(string); ok {
			return id
		}
	}
	return ""
}

// GetUserRoleFromContext الحصول على دور المستخدم من السياق
func GetUserRoleFromContext(ctx context.Context) string {
	if userRole, ok := ctx.Value("userRole").(string); ok {
		return userRole
	}
	return ""
}

// WithTimeout إنشاء سياق مع مهلة زمنية
func WithTimeout(ctx context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, timeout)
}

// ========== دوال النظام والأداء ==========

// GetMemoryUsageMB الحصول على استخدام الذاكرة بالميجابايت
func GetMemoryUsageMB() MemoryStats {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	usedMB := float64(m.Alloc) / 1024 / 1024
	totalMB := float64(m.Sys) / 1024 / 1024
	usagePercentage := (usedMB / totalMB) * 100

	return MemoryStats{
		UsedMB:          math.Round(usedMB*100) / 100,
		TotalMB:         math.Round(totalMB*100) / 100,
		UsagePercentage: math.Round(usagePercentage*100) / 100,
	}
}

// GetMemoryUsageMBFloat الحصول على استخدام الذاكرة (نسخة مبسطة)
func GetMemoryUsageMBFloat() float64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return float64(m.Alloc) / 1024 / 1024
}

// GetGoroutineCount الحصول على عدد الـ goroutines
func GetGoroutineCount() int {
	return runtime.NumGoroutine()
}

// GetCPUUsage الحصول على استخدام المعالج (محاكاة)
func GetCPUUsage() float64 {
	// محاكاة مبسطة لاستخدام المعالج
	return math.Round((float64(runtime.NumGoroutine())/1000)*10000) / 100
}

// FormatBytes تنسيق حجم الملف إلى صيغة مقروءة
func FormatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// ========== دوال الملفات والرفع ==========

// GetFileExtension الحصول على امتداد الملف
func GetFileExtension(filename string) string {
	parts := strings.Split(filename, ".")
	if len(parts) > 1 {
		return strings.ToLower(parts[len(parts)-1])
	}
	return ""
}

// IsAllowedFileType التحقق من نوع الملف المسموح
func IsAllowedFileType(filename string, allowedTypes []string) bool {
	ext := GetFileExtension(filename)
	for _, allowedType := range allowedTypes {
		if strings.EqualFold(ext, allowedType) {
			return true
		}
	}
	return false
}

// CalculateFileSizeMB حساب حجم الملف بالميجابايت
func CalculateFileSizeMB(size int64) float64 {
	return float64(size) / 1024 / 1024
}

// ValidateFileSize التحقق من حجم الملف
func ValidateFileSize(size int64, maxSizeMB int64) bool {
	return size <= maxSizeMB*1024*1024
}

// ========== دوال السلسلة النصية ==========

// TruncateString تقصير السلسلة النصية
func TruncateString(s string, maxLength int) string {
	if len(s) <= maxLength {
		return s
	}
	return s[:maxLength] + "..."
}

// ToCamelCase تحويل إلى CamelCase
func ToCamelCase(s string) string {
	words := strings.Fields(strings.ReplaceAll(s, "_", " "))
	for i, word := range words {
		if i == 0 {
			words[i] = strings.ToLower(word)
		} else {
			words[i] = strings.Title(strings.ToLower(word))
		}
	}
	return strings.Join(words, "")
}

// ToSnakeCase تحويل إلى snake_case
func ToSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 {
				result.WriteByte('_')
			}
			result.WriteRune(unicode.ToLower(r))
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}

// ContainsString التحقق إذا كانت المصفوفة تحتوي على قيمة
func ContainsString(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// RemoveDuplicates إزالة التكرارات من المصفوفة
func RemoveDuplicates(slice []string) []string {
	keys := make(map[string]bool)
	var result []string
	for _, item := range slice {
		if _, value := keys[item]; !value {
			keys[item] = true
			result = append(result, item)
		}
	}
	return result
}

// ========== دوال الرياضيات ==========

// Round التقريب إلى منازل محددة
func Round(value float64, decimals int) float64 {
	multiplier := math.Pow(10, float64(decimals))
	return math.Round(value*multiplier) / multiplier
}

// CalculatePercentage حساب النسبة المئوية
func CalculatePercentage(part, total float64) float64 {
	if total == 0 {
		return 0
	}
	return (part / total) * 100
}

// CalculateDiscount حساب الخصم
func CalculateDiscount(originalPrice, discountPercent float64) float64 {
	return originalPrice * (discountPercent / 100)
}

// CalculateTax حساب الضريبة
func CalculateTax(amount, taxRate float64) float64 {
	return amount * (taxRate / 100)
}

// ========== دوال المساعدة للخدمات ==========

// GenerateServiceSlug إنشاء slug للخدمة
func GenerateServiceSlug(title string) string {
	// إزالة الرموز الخاصة
	reg := regexp.MustCompile("[^a-zA-Z0-9\\s-]")
	slug := reg.ReplaceAllString(title, "")

	// استبدال المسافات بشرطات
	slug = strings.ReplaceAll(slug, " ", "-")

	// تحويل إلى أحرف صغيرة
	slug = strings.ToLower(slug)

	// إزالة الشرطات المكررة
	reg = regexp.MustCompile("-+")
	slug = reg.ReplaceAllString(slug, "-")

	return slug
}

// GenerateOrderNumber إنشاء رقم طلب
func GenerateOrderNumber() string {
	timestamp := time.Now().Unix()
	randomPart, _ := GenerateRandomString(6)
	return fmt.Sprintf("ORD-%d-%s", timestamp, randomPart)
}

// GenerateTrackingNumber إنشاء رقم تتبع
func GenerateTrackingNumber() string {
	timestamp := time.Now().UnixNano()
	return fmt.Sprintf("TRK%d", timestamp)
}

// CalculateOrderTotal حساب الإجمالي للطلب
func CalculateOrderTotal(subtotal, tax, shipping, discount float64) float64 {
	return subtotal + tax + shipping - discount
}

// ========== دوال التسجيل والتصحيح ==========

// LogOperation تسجيل عملية مع الوقت
func LogOperation(ctx context.Context, operation string, fn func() error) error {
	start := time.Now()

	err := fn()

	duration := time.Since(start)
	if err != nil {
		// تسجيل الخطأ
		return err
	}

	// تسجيل النجاح
	_ = duration
	return nil
}

// MeasureExecutionTime قياس وقت التنفيذ
func MeasureExecutionTime(ctx context.Context, name string, fn func()) time.Duration {
	start := time.Now()
	fn()
	return time.Since(start)
}

// ========== دوال الشبكة والـ HTTP ==========

// GetClientIP الحصول على IP العميل
func GetClientIP(r *http.Request) string {
	// التحقق من الرؤوس أولاً
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		ips := strings.Split(ip, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}

	// استخدام العنوان المباشر
	if host, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
		return host
	}
	return r.RemoteAddr
}

// GetUserAgent الحصول على user agent
func GetUserAgent(r *http.Request) string {
	return r.Header.Get("User-Agent")
}

// IsMobileRequest التحقق إذا كان الطلب من جهاز محمول
func IsMobileRequest(r *http.Request) bool {
	userAgent := strings.ToLower(GetUserAgent(r))
	mobileKeywords := []string{"mobile", "android", "iphone", "ipod", "ipad", "blackberry", "windows phone"}

	for _, keyword := range mobileKeywords {
		if strings.Contains(userAgent, keyword) {
			return true
		}
	}
	return false
}

// ========== دوال القراءة والكتابة ==========

// ReadAll قراءة كل البيانات من القارئ
func ReadAll(r io.Reader) ([]byte, error) {
	return io.ReadAll(r)
}

// CopyData نسخ البيانات من قارئ إلى كاتب
func CopyData(dst io.Writer, src io.Reader) (int64, error) {
	return io.Copy(dst, src)
}

// ========== دوال التطبيقات الخاصة ==========

// CalculateRating حساب التقييم
func CalculateRating(ratings []int) float64 {
	if len(ratings) == 0 {
		return 0
	}

	sum := 0
	for _, rating := range ratings {
		sum += rating
	}

	return Round(float64(sum)/float64(len(ratings)), 1)
}

// GenerateVerificationCode إنشاء رمز التحقق
func GenerateVerificationCode() string {
	// إنشاء رمز مكون من 6 أرقام
	random, _ := GenerateRandomNumber(0, 999999)
	return fmt.Sprintf("%06d", random)
}

// FormatCurrency تنسيق العملة
func FormatCurrency(amount float64, currency string) string {
	switch currency {
	case "USD":
		return fmt.Sprintf("$%.2f", amount)
	case "EUR":
		return fmt.Sprintf("€%.2f", amount)
	case "SAR":
		return fmt.Sprintf("%.2f ر.س", amount)
	default:
		return fmt.Sprintf("%.2f %s", amount, currency)
	}
}

// GetDefaultAvatarURL الحصول على صورة افتراضية
func GetDefaultAvatarURL() string {
	return "/assets/images/default-avatar.png"
}

// CalculateAge حساب العمر
func CalculateAge(birthDate time.Time) int {
	now := Now()
	years := now.Year() - birthDate.Year()

	// إذا لم يحن عيد الميلاد بعد هذا العام، نطرح سنة
	if now.YearDay() < birthDate.YearDay() {
		years--
	}

	return years
}

// ========== دوال مساعدة إضافية ==========

// GenerateSlug إنشاء slug من النص
func GenerateSlug(text string) string {
	// إزالة الرموز الخاصة
	reg := regexp.MustCompile("[^a-zA-Z0-9\\s]")
	slug := reg.ReplaceAllString(text, "")
	
	// استبدال المسافات بشرطات
	slug = strings.ReplaceAll(slug, " ", "-")
	
	// تحويل إلى أحرف صغيرة
	slug = strings.ToLower(slug)
	
	// إزالة الشرطات المكررة
	reg = regexp.MustCompile("-+")
	slug = reg.ReplaceAllString(slug, "-")
	
	// إزالة الشرطات من البداية والنهاية
	slug = strings.Trim(slug, "-")
	
	return slug
}

// GenerateRandomNumber إنشاء رقم عشوائي
func GenerateRandomNumber(min, max int) int {
	if max <= min {
		return min
	}
	
	nBig, err := rand.Int(rand.Reader, big.NewInt(int64(max-min+1)))
	if err != nil {
		// Fallback بسيط
		return min + int(time.Now().UnixNano()%int64(max-min+1))
	}
	return min + int(nBig.Int64())
}

// GenerateRandomColor إنشاء لون عشوائي
func GenerateRandomColor() string {
	colors := []string{
		"#FF6B6B", "#4ECDC4", "#45B7D1", "#96CEB4", "#FFEAA7",
		"#DDA0DD", "#98D8C8", "#F7DC6F", "#BB8FCE", "#85C1E9",
	}
	return colors[GenerateRandomNumber(0, len(colors)-1)]
}

// ExtractUsernameFromEmail استخراج اسم المستخدم من البريد الإلكتروني
func ExtractUsernameFromEmail(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) > 0 {
		return parts[0]
	}
	return email
}

// MaskEmail إخفاء جزء من البريد الإلكتروني
func MaskEmail(email string) string {
	if !IsValidEmail(email) {
		return email
	}
	
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return email
	}
	
	username := parts[0]
	domain := parts[1]
	
	if len(username) <= 2 {
		return username + "@" + domain
	}
	
	maskedUsername := username[:2] + strings.Repeat("*", len(username)-2)
	return maskedUsername + "@" + domain
}

// MaskPhone إخفاء جزء من رقم الهاتف
func MaskPhone(phone string) string {
	if len(phone) <= 4 {
		return strings.Repeat("*", len(phone))
	}
	
	visiblePart := phone[len(phone)-4:]
	maskedPart := strings.Repeat("*", len(phone)-4)
	return maskedPart + visiblePart
}

// GenerateOTP إنشاء OTP
func GenerateOTP(length int) string {
	otp := ""
	for i := 0; i < length; i++ {
		otp += strconv.Itoa(GenerateRandomNumber(0, 9))
	}
	return otp
}

// ParseQueryParams تحليل معاملات الاستعلام
func ParseQueryParams(c *gin.Context, params []string) map[string]string {
	result := make(map[string]string)
	for _, param := range params {
		value := c.Query(param)
		if value != "" {
			result[param] = value
		}
	}
	return result
}

// GetQueryInt الحصول على قيمة عددية من الاستعلام
func GetQueryInt(c *gin.Context, key string, defaultValue int) int {
	value := c.Query(key)
	if value == "" {
		return defaultValue
	}
	
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	
	return intValue
}

// GetQueryFloat الحصول على قيمة عشرية من الاستعلام
func GetQueryFloat(c *gin.Context, key string, defaultValue float64) float64 {
	value := c.Query(key)
	if value == "" {
		return defaultValue
	}
	
	floatValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return defaultValue
	}
	
	return floatValue
}

// GetQueryBool الحصول على قيمة منطقية من الاستعلام
func GetQueryBool(c *gin.Context, key string, defaultValue bool) bool {
	value := c.Query(key)
	if value == "" {
		return defaultValue
	}
	
	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}
	
	return boolValue
}

// GenerateSessionID إنشاء معرف جلسة
func GenerateSessionID() (string, error) {
	timestamp := time.Now().UnixNano()
	random, err := GenerateRandomString(16)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("sess_%d_%s", timestamp, random), nil
}

// GenerateUploadToken إنشاء توكن رفع
func GenerateUploadToken() (string, error) {
	timestamp := time.Now().Unix()
	random, err := GenerateRandomString(32)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("upl_%d_%s", timestamp, random), nil
}

// ValidateUploadToken التحقق من توكن الرفع
func ValidateUploadToken(token string, maxAge time.Duration) bool {
	if !strings.HasPrefix(token, "upl_") {
		return false
	}
	
	parts := strings.Split(token, "_")
	if len(parts) != 3 {
		return false
	}
	
	timestamp, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return false
	}
	
	uploadTime := time.Unix(timestamp, 0)
	return time.Since(uploadTime) <= maxAge
}

// ========== دوال مساعدة للـ Context ==========

// SetContextValue تعيين قيمة في السياق
func SetContextValue(ctx context.Context, key, value interface{}) context.Context {
	return context.WithValue(ctx, key, value)
}

// GetContextValue الحصول على قيمة من السياق
func GetContextValue(ctx context.Context, key interface{}) interface{} {
	return ctx.Value(key)
}

// CreateRequestContext إنشاء سياق طلب جديد
func CreateRequestContext() (context.Context, error) {
	ctx := context.Background()
	randomStr, err := GenerateRandomString(16)
	if err != nil {
		return ctx, err
	}
	ctx = SetContextValue(ctx, "requestID", randomStr)
	ctx = SetContextValue(ctx, "requestTime", time.Now())
	return ctx, nil
}

// GetRequestDuration الحصول على مدة الطلب
func GetRequestDuration(ctx context.Context) time.Duration {
	startTime := GetContextValue(ctx, "requestTime")
	if startTime == nil {
		return 0
	}
	
	if t, ok := startTime.(time.Time); ok {
		return time.Since(t)
	}
	
	return 0
}

// ========== دوال التطبيق الأساسية ==========

// GetCurrentUserID الحصول على معرف المستخدم الحالي
func GetCurrentUserID(c *gin.Context) string {
	// حاول الحصول من سياق Gin أولاً
	if userID := GetUserIDFromGinContext(c); userID != "" {
		return userID
	}
	
	// حاول الحصول من سياق HTTP
	if userID := c.GetString("userID"); userID != "" {
		return userID
	}
	
	// حاول الحصول من الرؤوس
	if userID := c.GetHeader("X-User-ID"); userID != "" {
		return userID
	}
	
	return ""
}

// GetCurrentUserRole للحصول على دور المستخدم الحالي
func GetCurrentUserRole(c *gin.Context) string {
	// حاول الحصول من سياق Gin
	if userRole, exists := c.Get("userRole"); exists {
		if role, ok := userRole.(string); ok {
			return role
		}
	}
	
	// حاول الحصول من الرؤوس
	if userRole := c.GetHeader("X-User-Role"); userRole != "" {
		return userRole
	}
	
	return "user" // القيمة الافتراضية
}

// IsAdminUser التحقق إذا كان المستخدم مشرفاً
func IsAdminUser(c *gin.Context) bool {
	role := GetCurrentUserRole(c)
	return role == "admin" || role == "superadmin"
}

// IsAuthenticatedUser التحقق إذا كان المستخدم مصادقاً عليه
func IsAuthenticatedUser(c *gin.Context) bool {
	return GetCurrentUserID(c) != ""
}

// ValidateAdminAccess التحقق من صلاحيات المشرف
func ValidateAdminAccess(c *gin.Context) bool {
	if !IsAuthenticatedUser(c) {
		return false
	}
	return IsAdminUser(c)
}

// GenerateResponseData إنشاء بيانات الاستجابة
func GenerateResponseData(data interface{}, pagination *Pagination) gin.H {
	response := gin.H{
		"success": true,
		"data":    data,
	}
	
	if pagination != nil {
		response["pagination"] = pagination
	}
	
	return response
}

// GenerateErrorResponse إنشاء استجابة خطأ
func GenerateErrorResponse(message, errorCode string) gin.H {
	return gin.H{
		"success": false,
		"message": message,
		"error":   errorCode,
	}
}