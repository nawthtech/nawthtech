package cloudflare

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/nawthtech/nawthtech/backend/internal/logger"
)

// ================================
// هياكل البيانات
// ================================

// CloudflareConfig إعدادات Cloudflare
type CloudflareConfig struct {
	ZoneID    string
	APIKey    string
	Email     string
	Enabled   bool
	CacheEnabled bool
	AnalyticsEnabled bool
}

// PurgeCacheRequest طلب مسح ذاكرة التخزين المؤقت
type PurgeCacheRequest struct {
	Files []string `json:"files,omitempty"`
	Tags  []string `json:"tags,omitempty"`
	Hosts []string `json:"hosts,omitempty"`
	Prefixes []string `json:"prefixes,omitempty"`
}

// PurgeCacheResponse استجابة مسح ذاكرة التخزين المؤقت
type PurgeCacheResponse struct {
	Success bool `json:"success"`
	Errors  []interface{} `json:"errors"`
	Messages []interface{} `json:"messages"`
}

// AnalyticsData بيانات التحليلات
type AnalyticsData struct {
	Timestamp string `json:"timestamp"`
	Requests  int64  `json:"requests"`
	Bandwidth int64  `json:"bandwidth"`
	Threats   int64  `json:"threats"`
}

// ================================
// دوال التهيئة
// ================================

// NewCloudflareConfig إنشاء إعدادات Cloudflare جديدة
func NewCloudflareConfig() *CloudflareConfig {
	return &CloudflareConfig{
		ZoneID:    os.Getenv("CLOUDFLARE_ZONE_ID"),
		APIKey:    os.Getenv("CLOUDFLARE_API_KEY"),
		Email:     os.Getenv("CLOUDFLARE_EMAIL"),
		Enabled:   getEnvBool("CLOUDFLARE_PROTECTION_ENABLED", true),
		CacheEnabled: getEnvBool("CLOUDFLARE_CACHE_ENABLED", true),
		AnalyticsEnabled: getEnvBool("CLOUDFLARE_ANALYTICS_ENABLED", true),
	}
}

// InitCloudflareService تهيئة خدمة Cloudflare
func InitCloudflareService() (*CloudflareConfig, error) {
	config := NewCloudflareConfig()
	
	if config.ZoneID == "" || config.APIKey == "" || config.Email == "" {
		return nil, fmt.Errorf("إعدادات Cloudflare غير مكتملة")
	}

	logger.Info(context.Background(), "✅ تم تهيئة خدمة Cloudflare",
		"zone_id", config.ZoneID,
		"enabled", config.Enabled,
		"cache_enabled", config.CacheEnabled,
	)

	return config, nil
}

// ================================
// دوال إدارة ذاكرة التخزين المؤقت
// ================================

// PurgeCache مسح ذاكرة التخزين المؤقت لـ Cloudflare
func PurgeCache(urls []string) error {
	config := NewCloudflareConfig()
	if !config.Enabled || !config.CacheEnabled {
		logger.Info(context.Background(), "⚠️ خدمة Cloudflare غير مفعلة - تخطي مسح الذاكرة المؤقتة")
		return nil
	}

	startTime := time.Now()

	payload := PurgeCacheRequest{
		Files: urls,
	}

	response, err := makeCloudflareRequest("POST", "/purge_cache", payload)
	if err != nil {
		logger.Error(context.Background(), "❌ فشل في مسح ذاكرة التخزين المؤقت لـ Cloudflare",
			"urls", urls,
			"duration", time.Since(startTime),
			"error", err.Error(),
		)
		return err
	}

	logger.Info(context.Background(), "✅ تم مسح ذاكرة التخزين المؤقت لـ Cloudflare بنجاح",
		"urls_count", len(urls),
		"duration", time.Since(startTime),
		"success", response.Success,
	)

	return nil
}

// PurgeEverything مسح كل ذاكرة التخزين المؤقت
func PurgeEverything() error {
	config := NewCloudflareConfig()
	if !config.Enabled || !config.CacheEnabled {
		return nil
	}

	startTime := time.Now()

	payload := PurgeCacheRequest{
		PurgeEverything: true,
	}

	response, err := makeCloudflareRequest("POST", "/purge_cache", payload)
	if err != nil {
		logger.Error(context.Background(), "❌ فشل في مسح كل ذاكرة التخزين المؤقت",
			"duration", time.Since(startTime),
			"error", err.Error(),
		)
		return err
	}

	logger.Info(context.Background(), "✅ تم مسح كل ذاكرة التخزين المؤقت بنجاح",
		"duration", time.Since(startTime),
		"success", response.Success,
	)

	return nil
}

// PurgeByTags مسح الذاكرة المؤقتة باستخدام الوسوم
func PurgeByTags(tags []string) error {
	config := NewCloudflareConfig()
	if !config.Enabled || !config.CacheEnabled {
		return nil
	}

	startTime := time.Now()

	payload := PurgeCacheRequest{
		Tags: tags,
	}

	response, err := makeCloudflareRequest("POST", "/purge_cache", payload)
	if err != nil {
		logger.Error(context.Background(), "❌ فشل في مسح الذاكرة المؤقتة بالوسوم",
			"tags", tags,
			"duration", time.Since(startTime),
			"error", err.Error(),
		)
		return err
	}

	logger.Info(context.Background(), "✅ تم مسح الذاكرة المؤقتة بالوسوم بنجاح",
		"tags_count", len(tags),
		"duration", time.Since(startTime),
		"success", response.Success,
	)

	return nil
}

// ================================
// دوال التحليلات والإحصائيات
// ================================

// GetAnalytics الحصول على إحصائيات Cloudflare
func GetAnalytics(startTime, endTime time.Time) (*AnalyticsData, error) {
	config := NewCloudflareConfig()
	if !config.Enabled || !config.AnalyticsEnabled {
		return nil, nil
	}

	endpoint := fmt.Sprintf("/analytics/dashboard?since=%s&until=%s",
		startTime.Format(time.RFC3339),
		endTime.Format(time.RFC3339),
	)

	response, err := makeCloudflareRequest("GET", endpoint, nil)
	if err != nil {
		logger.Error(context.Background(), "❌ فشل في الحصول على إحصائيات Cloudflare",
			"error", err.Error(),
		)
		return nil, err
	}

	var analytics AnalyticsData
	if err := json.Unmarshal(response, &analytics); err != nil {
		return nil, err
	}

	return &analytics, nil
}

// GetZoneDetails الحصول على تفاصيل المنطقة
func GetZoneDetails() (map[string]interface{}, error) {
	config := NewCloudflareConfig()
	if !config.Enabled {
		return nil, nil
	}

	response, err := makeCloudflareRequest("GET", "", nil)
	if err != nil {
		return nil, err
	}

	var zoneData map[string]interface{}
	if err := json.Unmarshal(response, &zoneData); err != nil {
		return nil, err
	}

	return zoneData, nil
}

// ================================
// دوال الحماية والأمان
// ================================

// GetSecurityEvents الحصول على أحداث الأمان
func GetSecurityEvents(startTime, endTime time.Time) ([]interface{}, error) {
	config := NewCloudflareConfig()
	if !config.Enabled {
		return nil, nil
	}

	endpoint := fmt.Sprintf("/security/events?since=%s&until=%s",
		startTime.Format(time.RFC3339),
		endTime.Format(time.RFC3339),
	)

	response, err := makeCloudflareRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	var events struct {
		Result []interface{} `json:"result"`
	}
	if err := json.Unmarshal(response, &events); err != nil {
		return nil, err
	}

	return events.Result, nil
}

// CreateFirewallRule إنشاء قاعدة جدار حماية
func CreateFirewallRule(rule map[string]interface{}) error {
	config := NewCloudflareConfig()
	if !config.Enabled {
		return nil
	}

	_, err := makeCloudflareRequest("POST", "/firewall/rules", rule)
	if err != nil {
		logger.Error(context.Background(), "❌ فشل في إنشاء قاعدة جدار الحماية",
			"error", err.Error(),
		)
		return err
	}

	logger.Info(context.Background(), "✅ تم إنشاء قاعدة جدار الحماية بنجاح")
	return nil
}

// ================================
// دوال المساعدة
// ================================

// makeCloudflareRequest تنفيذ طلب إلى Cloudflare API
func makeCloudflareRequest(method, endpoint string, payload interface{}) ([]byte, error) {
	config := NewCloudflareConfig()
	
	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s%s", config.ZoneID, endpoint)
	
	var body io.Reader
	if payload != nil {
		jsonData, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		body = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Auth-Email", config.Email)
	req.Header.Set("X-Auth-Key", config.APIKey)
	req.Header.Set("User-Agent", "NawthTech-Backend/1.0")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Cloudflare API error: %s - %s", resp.Status, string(responseBody))
	}

	return responseBody, nil
}

// IsEnabled التحقق إذا كانت خدمة Cloudflare مفعلة
func IsEnabled() bool {
	config := NewCloudflareConfig()
	return config.Enabled
}

// GetConfig الحصول على إعدادات Cloudflare
func GetConfig() *CloudflareConfig {
	return NewCloudflareConfig()
}

// ================================
// دوال مساعدة للبيئة
// ================================

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

// HealthCheck فحص صحة اتصال Cloudflare
func HealthCheck() map[string]interface{} {
	config := NewCloudflareConfig()
	
	if !config.Enabled {
		return map[string]interface{}{
			"service": "cloudflare",
			"status":  "disabled",
			"enabled": false,
		}
	}

	_, err := GetZoneDetails()
	if err != nil {
		return map[string]interface{}{
			"service": "cloudflare",
			"status":  "error",
			"enabled": true,
			"error":   err.Error(),
		}
	}

	return map[string]interface{}{
		"service": "cloudflare",
		"status":  "healthy",
		"enabled": true,
		"zone_id": config.ZoneID,
	}
}