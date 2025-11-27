package utils

import (
	"fmt"
	"math"
	"regexp"
	"strings"
	"time"
	"unicode"

	"github.com/nawthtech/nawthtech/backend/internal/models"
)

// ServiceValidator مدقق الخدمات
type ServiceValidator struct{}

// ValidateService التحقق من صحة بيانات الخدمة
func (v *ServiceValidator) ValidateService(service *models.ServiceCreateRequest) error {
	if err := v.validateTitle(service.Title); err != nil {
		return err
	}

	if err := v.validateDescription(service.Description); err != nil {
		return err
	}

	if err := v.validatePrice(service.Price); err != nil {
		return err
	}

	if err := v.validateDuration(service.Duration); err != nil {
		return err
	}

	if err := v.validateCategory(service.Category); err != nil {
		return err
	}

	if err := v.validateImages(service.Images); err != nil {
		return err
	}

	if err := v.validateFeatures(service.Features); err != nil {
		return err
	}

	if err := v.validateTags(service.Tags); err != nil {
		return err
	}

	return nil
}

// validateTitle التحقق من صحة العنوان
func (v *ServiceValidator) validateTitle(title string) error {
	if len(title) < 3 {
		return fmt.Errorf("العنوان يجب أن يكون على الأقل 3 أحرف")
	}

	if len(title) > 200 {
		return fmt.Errorf("العنوان يجب ألا يتجاوز 200 حرف")
	}

	// التحقق من وجود أحرف غير مسموحة
	if matched, _ := regexp.MatchString(`[<>{}]`, title); matched {
		return fmt.Errorf("العنوان يحتوي على أحرف غير مسموحة")
	}

	return nil
}

// validateDescription التحقق من صحة الوصف
func (v *ServiceValidator) validateDescription(description string) error {
	if len(description) < 10 {
		return fmt.Errorf("الوصف يجب أن يكون على الأقل 10 أحرف")
	}

	if len(description) > 2000 {
		return fmt.Errorf("الوصف يجب ألا يتجاوز 2000 حرف")
	}

	// حساب نسبة الأحرف العربية (اختياري)
	arabicRatio := v.calculateArabicRatio(description)
	if arabicRatio < 0.5 {
		return fmt.Errorf("الوصف يجب أن يحتوي على نسبة أعلى من النص العربي")
	}

	return nil
}

// validatePrice التحقق من صحة السعر
func (v *ServiceValidator) validatePrice(price float64) error {
	if price <= 0 {
		return fmt.Errorf("السعر يجب أن يكون أكبر من الصفر")
	}

	if price > 1000000 {
		return fmt.Errorf("السعر يتجاوز الحد المسموح")
	}

	return nil
}

// validateDuration التحقق من صحة المدة
func (v *ServiceValidator) validateDuration(duration int) error {
	if duration < 1 {
		return fmt.Errorf("المدة يجب أن تكون يوم على الأقل")
	}

	if duration > 365 {
		return fmt.Errorf("المدة تتجاوز الحد المسموح (سنة واحدة)")
	}

	return nil
}

// validateCategory التحقق من صحة الفئة
func (v *ServiceValidator) validateCategory(category string) error {
	if category == "" {
		return fmt.Errorf("الفئة مطلوبة")
	}

	// قائمة الفئات المسموحة (يمكن جلبها من قاعدة البيانات)
	allowedCategories := []string{
		"تصميم", "تطوير", "كتابة", "ترجمة", "تسويق", 
		"تعليم", "استشارات", "برمجة", "جرافيك", "فيديو",
	}

	found := false
	for _, allowed := range allowedCategories {
		if strings.EqualFold(category, allowed) {
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("الفئة غير مسموحة")
	}

	return nil
}

// validateImages التحقق من صحة الصور
func (v *ServiceValidator) validateImages(images []string) error {
	if len(images) > 10 {
		return fmt.Errorf("لا يمكن إضافة أكثر من 10 صور")
	}

	for _, image := range images {
		if !v.isValidImageURL(image) {
			return fmt.Errorf("رابط الصورة غير صالح: %s", image)
		}
	}

	return nil
}

// validateFeatures التحقق من صحة الميزات
func (v *ServiceValidator) validateFeatures(features []string) error {
	if len(features) > 20 {
		return fmt.Errorf("لا يمكن إضافة أكثر من 20 ميزة")
	}

	for _, feature := range features {
		if len(feature) < 3 {
			return fmt.Errorf("الميزة يجب أن تكون على الأقل 3 أحرف")
		}

		if len(feature) > 100 {
			return fmt.Errorf("الميزة يجب ألا تتجاوز 100 حرف")
		}
	}

	return nil
}

// validateTags التحقق من صحة الوسوم
func (v *ServiceValidator) validateTags(tags []string) error {
	if len(tags) > 15 {
		return fmt.Errorf("لا يمكن إضافة أكثر من 15 وسم")
	}

	for _, tag := range tags {
		if len(tag) < 2 {
			return fmt.Errorf("الوسم يجب أن يكون على الأقل حرفين")
		}

		if len(tag) > 30 {
			return fmt.Errorf("الوسم يجب ألا يتجاوز 30 حرف")
		}

		// التحقق من أحرف الوسم
		if matched, _ := regexp.MatchString(`^[a-zA-Z0-9\u0600-\u06FF\s\-_]+$`, tag); !matched {
			return fmt.Errorf("الوسم يحتوي على أحرف غير مسموحة: %s", tag)
		}
	}

	return nil
}

// ServiceCalculator حاسبة الخدمات
type ServiceCalculator struct{}

// CalculateRating حساب التقييم الجديد
func (c *ServiceCalculator) CalculateRating(currentRating float64, currentCount int, newRating int) float64 {
	if currentCount == 0 {
		return float64(newRating)
	}

	totalScore := currentRating * float64(currentCount)
	totalScore += float64(newRating)
	
	return totalScore / float64(currentCount+1)
}

// CalculateDiscount حساب السعر بعد الخصم
func (c *ServiceCalculator) CalculateDiscount(originalPrice, discountPercent float64) float64 {
	if discountPercent <= 0 {
		return originalPrice
	}

	discountAmount := originalPrice * (discountPercent / 100)
	return math.Max(0, originalPrice-discountAmount)
}

// CalculateCompletionTime حساب وقت الإنجاز المتوقع
func (c *ServiceCalculator) CalculateCompletionTime(duration int, complexity string) time.Duration {
	baseDuration := time.Duration(duration) * 24 * time.Hour

	// عوامل التعقيد
	complexityFactors := map[string]float64{
		"low":      0.8,
		"medium":   1.0,
		"high":     1.5,
		"veryhigh": 2.0,
	}

	factor, exists := complexityFactors[complexity]
	if !exists {
		factor = 1.0
	}

	return time.Duration(float64(baseDuration) * factor)
}

// ServiceFormatter مُنسق الخدمات
type ServiceFormatter struct{}

// FormatPrice تنسيق السعر
func (f *ServiceFormatter) FormatPrice(price float64) string {
	if price == 0 {
		return "مجاني"
	}

	return fmt.Sprintf("%.2f ر.س", price)
}

// FormatDuration تنسيق المدة
func (f *ServiceFormatter) FormatDuration(duration int) string {
	if duration == 1 {
		return "يوم واحد"
	} else if duration == 7 {
		return "أسبوع واحد"
	} else if duration == 30 {
		return "شهر واحد"
	} else if duration < 30 {
		return fmt.Sprintf("%d أيام", duration)
	} else {
		months := duration / 30
		return fmt.Sprintf("%d أشهر", months)
	}
}

// FormatRating تنسيق التقييم
func (f *ServiceFormatter) FormatRating(rating float64) string {
	return fmt.Sprintf("%.1f/5", rating)
}

// TruncateDescription تقصير الوصف
func (f *ServiceFormatter) TruncateDescription(description string, maxLength int) string {
	if len(description) <= maxLength {
		return description
	}

	// التأكد من عدم قطع الكلمة في المنتصف
	truncated := description[:maxLength]
	lastSpace := strings.LastIndex(truncated, " ")
	if lastSpace > 0 {
		truncated = truncated[:lastSpace]
	}

	return truncated + "..."
}

// ServiceSearchOptimizer محسن البحث في الخدمات
type ServiceSearchOptimizer struct{}

// OptimizeSearchQuery تحسين استعلام البحث
func (o *ServiceSearchOptimizer) OptimizeSearchQuery(query string) string {
	// إزالة المسافات الزائدة
	query = strings.TrimSpace(query)
	query = regexp.MustCompile(`\s+`).ReplaceAllString(query, " ")

	// إزالة أحرف خاصة
	query = regexp.MustCompile(`[^\w\s\u0600-\u06FF]`).ReplaceAllString(query, "")

	// تحويل إلى أحرف صغيرة
	query = strings.ToLower(query)

	return query
}

// GenerateSearchKeywords توليد كلمات البحث
func (o *ServiceSearchOptimizer) GenerateSearchKeywords(title, description string, tags []string) []string {
	keywords := make([]string, 0)

	// إضافة كلمات من العنوان
	titleWords := strings.Fields(title)
	keywords = append(keywords, titleWords...)

	// إضافة كلمات من الوصف (الكلمات الأكثر تكراراً)
	descWords := strings.Fields(description)
	wordCount := make(map[string]int)
	for _, word := range descWords {
		if len(word) > 2 { // تجاهل الكلمات القصيرة
			wordCount[word]++
		}
	}

	// إضافة الكلمات الأكثر تكراراً
	for word, count := range wordCount {
		if count >= 2 {
			keywords = append(keywords, word)
		}
	}

	// إضافة الوسوم
	keywords = append(keywords, tags...)

	// إزالة التكرارات
	return o.removeDuplicates(keywords)
}

// ServiceRecommender مُوصي الخدمات
type ServiceRecommender struct{}

// RecommendServices توصية الخدمات
func (r *ServiceRecommender) RecommendServices(userID string, userHistory []models.Service, allServices []models.Service, limit int) []models.ServiceRecommendation {
	recommendations := make([]models.ServiceRecommendation, 0)

	// خوارزمية توصية مبسطة (يمكن تطويرها)
	for _, service := range allServices {
		if r.shouldRecommend(service, userHistory) {
			similarity := r.calculateSimilarity(service, userHistory)
			reason := r.generateRecommendationReason(service, userHistory)

			recommendations = append(recommendations, models.ServiceRecommendation{
				ServiceID:   service.ID,
				Title:       service.Title,
				Category:    service.Category,
				Price:       service.Price,
				Rating:      service.Rating,
				TotalOrders: service.TotalOrders,
				Similarity:  similarity,
				Reason:      reason,
			})
		}
	}

	// ترتيب حسب درجة التشابه
	r.sortRecommendations(recommendations)

	// إرجاع العدد المطلوب فقط
	if len(recommendations) > limit {
		recommendations = recommendations[:limit]
	}

	return recommendations
}

// ========== الدوال المساعدة ==========

// calculateArabicRatio حساب نسبة الأحرف العربية
func (v *ServiceValidator) calculateArabicRatio(text string) float64 {
	if len(text) == 0 {
		return 0
	}

	arabicCount := 0
	totalCount := 0

	for _, char := range text {
		if unicode.IsLetter(char) {
			totalCount++
			// نطاق الأحرف العربية في Unicode
			if char >= 0x0600 && char <= 0x06FF {
				arabicCount++
			}
		}
	}

	if totalCount == 0 {
		return 0
	}

	return float64(arabicCount) / float64(totalCount)
}

// isValidImageURL التحقق من صحة رابط الصورة
func (v *ServiceValidator) isValidImageURL(url string) bool {
	validExtensions := []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}
	url = strings.ToLower(url)

	for _, ext := range validExtensions {
		if strings.HasSuffix(url, ext) {
			return true
		}
	}

	return false
}

// removeDuplicates إزالة التكرارات من المصفوفة
func (o *ServiceSearchOptimizer) removeDuplicates(items []string) []string {
	seen := make(map[string]bool)
	result := []string{}

	for _, item := range items {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}

	return result
}

// shouldRecommend التحقق إذا كان يجب التوصية بالخدمة
func (r *ServiceRecommender) shouldRecommend(service models.Service, userHistory []models.Service) bool {
	// عدم التوصية بالخدمات التي استخدمها المستخدم مسبقاً
	for _, history := range userHistory {
		if history.ID == service.ID {
			return false
		}
	}

	// شروط التوصية الأساسية
	return service.Status == "active" && 
	       service.Rating >= 3.0 && 
	       service.TotalOrders >= 1
}

// calculateSimilarity حساب درجة التشابه
func (r *ServiceRecommender) calculateSimilarity(service models.Service, userHistory []models.Service) float64 {
	if len(userHistory) == 0 {
		return 0.5 // درجة افتراضية
	}

	similarity := 0.0
	count := 0

	for _, history := range userHistory {
		if service.Category == history.Category {
			similarity += 0.3
		}

		priceDiff := math.Abs(service.Price - history.Price)
		if priceDiff <= 100 {
			similarity += 0.2
		}

		// تشابه الوسوم
		tagSimilarity := r.calculateTagSimilarity(service.Tags, history.Tags)
		similarity += tagSimilarity * 0.5

		count++
	}

	if count > 0 {
		similarity /= float64(count)
	}

	return math.Min(similarity, 1.0)
}

// calculateTagSimilarity حساب تشابه الوسوم
func (r *ServiceRecommender) calculateTagSimilarity(tags1, tags2 []string) float64 {
	if len(tags1) == 0 || len(tags2) == 0 {
		return 0
	}

	commonTags := 0
	for _, tag1 := range tags1 {
		for _, tag2 := range tags2 {
			if strings.EqualFold(tag1, tag2) {
				commonTags++
				break
			}
		}
	}

	return float64(commonTags) / float64(math.Max(float64(len(tags1)), float64(len(tags2))))
}

// generateRecommendationReason توليد سبب التوصية
func (r *ServiceRecommender) generateRecommendationReason(service models.Service, userHistory []models.Service) string {
	reasons := []string{}

	if service.Rating >= 4.5 {
		reasons = append(reasons, "تقييم ممتاز")
	}

	if service.TotalOrders >= 100 {
		reasons = append(reasons, "شائع جداً")
	}

	if len(userHistory) > 0 {
		for _, history := range userHistory {
			if service.Category == history.Category {
				reasons = append(reasons, "يشبه الخدمات التي استخدمتها")
				break
			}
		}
	}

	if len(reasons) == 0 {
		return "قد يعجبك"
	}

	return strings.Join(reasons, "، ")
}

// sortRecommendations ترتيب التوصيات
func (r *ServiceRecommender) sortRecommendations(recommendations []models.ServiceRecommendation) {
	for i := 0; i < len(recommendations)-1; i++ {
		for j := i + 1; j < len(recommendations); j++ {
			if recommendations[i].Similarity < recommendations[j].Similarity {
				recommendations[i], recommendations[j] = recommendations[j], recommendations[i]
			}
		}
	}
}