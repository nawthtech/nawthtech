package cloudinary

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
	"github.com/nawthtech/nawthtech/backend/internal/logger"
)

// ========== هياكل البيانات ==========

// CloudinaryService هيكل لخدمة Cloudinary
type CloudinaryService struct {
	cld *cloudinary.Cloudinary
	ctx context.Context
}

// UploadResult نتيجة الرفع
type UploadResult struct {
	PublicID     string `json:"public_id"`
	SecureURL    string `json:"secure_url"`
	Format       string `json:"format"`
	Bytes        int    `json:"bytes"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
	ResourceType string `json:"resource_type"`
}

// UploadOptions خيارات الرفع
type UploadOptions struct {
	Folder       string
	PublicID     string
	Overwrite    bool
	ResourceType string
}

// ========== دوال التهيئة ==========

// NewCloudinaryService إنشاء خدمة Cloudinary جديدة
func NewCloudinaryService() (*CloudinaryService, error) {
	// الحصول على البيانات من environment variables
	cloudName := os.Getenv("CLOUDINARY_CLOUD_NAME")
	apiKey := os.Getenv("CLOUDINARY_API_KEY")
	apiSecret := os.Getenv("CLOUDINARY_API_SECRET")

	if cloudName == "" || apiKey == "" || apiSecret == "" {
		return nil, fmt.Errorf("بيانات Cloudinary غير مكتملة - تأكد من تعيين CLOUDINARY_CLOUD_NAME, CLOUDINARY_API_KEY, CLOUDINARY_API_SECRET")
	}

	// إنشاء connection string
	connStr := fmt.Sprintf("cloudinary://%s:%s@%s", apiKey, apiSecret, cloudName)
	
	cld, err := cloudinary.NewFromURL(connStr)
	if err != nil {
		return nil, fmt.Errorf("فشل في تهيئة Cloudinary: %v", err)
	}

	cld.Config.URL.Secure = true
	ctx := context.Background()

	logger.Info(ctx, "✅ تم تهيئة خدمة Cloudinary بنجاح",
		"cloud_name", cloudName,
	)

	return &CloudinaryService{
		cld: cld,
		ctx: ctx,
	}, nil
}

// NewCloudinaryServiceWithContext إنشاء خدمة مع سياق مخصص
func NewCloudinaryServiceWithContext(ctx context.Context) (*CloudinaryService, error) {
	cloudName := os.Getenv("CLOUDINARY_CLOUD_NAME")
	apiKey := os.Getenv("CLOUDINARY_API_KEY")
	apiSecret := os.Getenv("CLOUDINARY_API_SECRET")

	if cloudName == "" || apiKey == "" || apiSecret == "" {
		return nil, fmt.Errorf("بيانات Cloudinary غير مكتملة")
	}

	connStr := fmt.Sprintf("cloudinary://%s:%s@%s", apiKey, apiSecret, cloudName)
	
	cld, err := cloudinary.NewFromURL(connStr)
	if err != nil {
		return nil, fmt.Errorf("فشل في تهيئة Cloudinary: %v", err)
	}

	cld.Config.URL.Secure = true

	logger.Info(ctx, "✅ تم تهيئة خدمة Cloudinary بنجاح",
		"cloud_name", cloudName,
	)

	return &CloudinaryService{
		cld: cld,
		ctx: ctx,
	}, nil
}

// ========== دوال الوصول ==========

// GetCredentials الحصول على بيانات الاعتماد
func (cs *CloudinaryService) GetCredentials() (*cloudinary.Cloudinary, context.Context) {
	return cs.cld, cs.ctx
}

// GetContext الحصول على السياق
func (cs *CloudinaryService) GetContext() context.Context {
	return cs.ctx
}

// ========== دوال الرفع ==========

// UploadImage رفع صورة إلى Cloudinary
func (cs *CloudinaryService) UploadImage(file interface{}, options ...UploadOptions) (*UploadResult, error) {
	startTime := time.Now()
	
	// معالجة الخيارات
	uploadParams := uploader.UploadParams{}
	if len(options) > 0 {
		opt := options[0]
		if opt.Folder != "" {
			uploadParams.Folder = opt.Folder
		} else {
			uploadParams.Folder = "nawthtech"
		}
		if opt.PublicID != "" {
			uploadParams.PublicID = opt.PublicID
		}
		uploadParams.Overwrite = &opt.Overwrite
		if opt.ResourceType != "" {
			uploadParams.ResourceType = opt.ResourceType
		}
	} else {
		uploadParams.Folder = "nawthtech"
	}

	result, err := cs.cld.Upload.Upload(cs.ctx, file, uploadParams)
	if err != nil {
		logger.Error(cs.ctx, "❌ فشل في رفع الصورة",
			"public_id", uploadParams.PublicID,
			"folder", uploadParams.Folder,
			"duration", time.Since(startTime),
			"error", err.Error(),
		)
		return nil, err
	}

	uploadResult := &UploadResult{
		PublicID:     result.PublicID,
		SecureURL:    result.SecureURL,
		Format:       result.Format,
		Bytes:        result.Bytes,
		Width:        result.Width,
		Height:       result.Height,
		ResourceType: result.ResourceType,
	}

	logger.Info(cs.ctx, "✅ تم رفع الصورة بنجاح",
		"public_id", uploadResult.PublicID,
		"folder", uploadParams.Folder,
		"url", uploadResult.SecureURL,
		"format", uploadResult.Format,
		"size_bytes", uploadResult.Bytes,
		"dimensions", fmt.Sprintf("%dx%d", uploadResult.Width, uploadResult.Height),
		"duration", time.Since(startTime),
	)

	return uploadResult, nil
}

// UploadImageFromFile رفع صورة من مسار ملف
func (cs *CloudinaryService) UploadImageFromFile(filePath string, options ...UploadOptions) (*UploadResult, error) {
	startTime := time.Now()

	// التحقق من وجود الملف
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("الملف غير موجود: %s", filePath)
	}

	// استخراج اسم الملف بدون امتداد لاستخدامه كـ public_id
	if len(options) == 0 {
		fileName := filepath.Base(filePath)
		fileExt := filepath.Ext(fileName)
		publicID := strings.TrimSuffix(fileName, fileExt)
		
		options = []UploadOptions{{
			PublicID: publicID,
			Folder:   "nawthtech",
		}}
	}

	result, err := cs.UploadImage(filePath, options...)
	if err != nil {
		logger.Error(cs.ctx, "❌ فشل في رفع الصورة من الملف",
			"file_path", filePath,
			"duration", time.Since(startTime),
			"error", err.Error(),
		)
		return nil, err
	}

	return result, nil
}

// UploadImageFromReader رفع صورة من قارئ
func (cs *CloudinaryService) UploadImageFromReader(reader io.Reader, publicID string, options ...UploadOptions) (*UploadResult, error) {
	startTime := time.Now()

	if len(options) == 0 {
		options = []UploadOptions{{
			PublicID: publicID,
			Folder:   "nawthtech",
		}}
	} else {
		options[0].PublicID = publicID
	}

	result, err := cs.UploadImage(reader, options...)
	if err != nil {
		logger.Error(cs.ctx, "❌ فشل في رفع الصورة من القارئ",
			"public_id", publicID,
			"duration", time.Since(startTime),
			"error", err.Error(),
		)
		return nil, err
	}

	return result, nil
}

// UploadImageFromGinFile رفع صورة من ملف Gin
func (cs *CloudinaryService) UploadImageFromGinFile(c *gin.Context, fieldName string, options ...UploadOptions) (*UploadResult, error) {
	startTime := time.Now()

	// الحصول على الملف من الطلب
	uploadedFile, err := c.FormFile(fieldName)
	if err != nil {
		return nil, fmt.Errorf("فشل في الحصول على الملف: %v", err)
	}

	// فتح الملف
	src, err := uploadedFile.Open()
	if err != nil {
		return nil, fmt.Errorf("فشل في فتح الملف: %v", err)
	}
	defer src.Close()

	// استخدام اسم الملف كـ public_id
	fileName := uploadedFile.Filename
	fileExt := filepath.Ext(fileName)
	publicID := strings.TrimSuffix(fileName, fileExt)

	if len(options) == 0 {
		options = []UploadOptions{{
			PublicID: publicID,
			Folder:   "nawthtech",
		}}
	} else {
		options[0].PublicID = publicID
	}

	result, err := cs.UploadImageFromReader(src, publicID, options...)
	if err != nil {
		logger.Error(cs.ctx, "❌ فشل في رفع الصورة من Gin",
			"field_name", fieldName,
			"file_name", fileName,
			"duration", time.Since(startTime),
			"error", err.Error(),
		)
		return nil, err
	}

	logger.Info(cs.ctx, "✅ تم رفع الصورة من Gin بنجاح",
		"field_name", fieldName,
		"file_name", fileName,
		"public_id", result.PublicID,
		"size_bytes", uploadedFile.Size,
		"duration", time.Since(startTime),
	)

	return result, nil
}

// ========== دوال الحذف ==========

// DeleteImage حذف صورة من Cloudinary
func (cs *CloudinaryService) DeleteImage(publicID string) error {
	startTime := time.Now()

	result, err := cs.cld.Upload.Destroy(cs.ctx, uploader.DestroyParams{
		PublicID: publicID,
	})

	if err != nil {
		logger.Error(cs.ctx, "❌ فشل في حذف الصورة",
			"public_id", publicID,
			"duration", time.Since(startTime),
			"error", err.Error(),
		)
		return err
	}

	logger.Info(cs.ctx, "✅ تم حذف الصورة بنجاح",
		"public_id", publicID,
		"result", result.Result,
		"duration", time.Since(startTime),
	)

	return nil
}

// DeleteImageByURL حذف صورة باستخدام URL
func (cs *CloudinaryService) DeleteImageByURL(imageURL string) error {
	startTime := time.Now()

	// استخراج public_id من URL
	publicID, err := cs.extractPublicIDFromURL(imageURL)
	if err != nil {
		return err
	}

	return cs.DeleteImage(publicID)
}

// ========== دوال المساعدة ==========

// extractPublicIDFromURL استخراج public_id من URL
func (cs *CloudinaryService) extractPublicIDFromURL(imageURL string) (string, error) {
	// هذا تنفيذ مبسط - يمكن تحسينه حسب احتياجاتك
	parts := strings.Split(imageURL, "/")
	if len(parts) < 2 {
		return "", fmt.Errorf("لا يمكن استخراج public_id من URL: %s", imageURL)
	}

	// آخر جزء في URL هو عادة public_id مع الامتداد
	lastPart := parts[len(parts)-1]
	publicID := strings.TrimSuffix(lastPart, filepath.Ext(lastPart))

	return publicID, nil
}

// GeneratePublicID إنشاء public_id فريد
func (cs *CloudinaryService) GeneratePublicID(prefix string) string {
	timestamp := time.Now().UnixNano()
	if prefix == "" {
		return fmt.Sprintf("img_%d", timestamp)
	}
	return fmt.Sprintf("%s_%d", prefix, timestamp)
}

// ValidateImage التحقق من صحة الصورة قبل الرفع
func (cs *CloudinaryService) ValidateImage(fileHeader *multipart.FileHeader) error {
	// التحقق من حجم الملف (5MB كحد أقصى)
	maxSize := int64(5 * 1024 * 1024) // 5MB
	if fileHeader.Size > maxSize {
		return fmt.Errorf("حجم الملف كبير جداً. الحد الأقصى المسموح به هو 5MB")
	}

	// التحقق من نوع الملف
	allowedTypes := []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}
	fileExt := strings.ToLower(filepath.Ext(fileHeader.Filename))
	
	allowed := false
	for _, allowedType := range allowedTypes {
		if fileExt == allowedType {
			allowed = true
			break
		}
	}

	if !allowed {
		return fmt.Errorf("نوع الملف غير مسموح به. الأنواع المسموحة: %v", allowedTypes)
	}

	return nil
}

// ========== دوال متقدمة ==========

// UploadMultiple رفع عدة ملفات مرة واحدة
func (cs *CloudinaryService) UploadMultiple(files []interface{}, options ...UploadOptions) ([]*UploadResult, error) {
	var results []*UploadResult
	var errors []string

	for i, file := range files {
		var currentOptions UploadOptions
		if len(options) > 0 {
			currentOptions = options[0]
		} else {
			currentOptions = UploadOptions{Folder: "nawthtech"}
		}

		// إنشاء public_id فريد لكل ملف
		if currentOptions.PublicID == "" {
			currentOptions.PublicID = cs.GeneratePublicID(fmt.Sprintf("file_%d", i))
		}

		result, err := cs.UploadImage(file, currentOptions)
		if err != nil {
			errors = append(errors, fmt.Sprintf("الملف %d: %v", i, err))
			continue
		}

		results = append(results, result)
	}

	if len(errors) > 0 {
		return results, fmt.Errorf("أخطاء في رفع بعض الملفات: %s", strings.Join(errors, "; "))
	}

	return results, nil
}

// GetResourceTypeFromExtension الحصول على نوع المورد من الامتداد
func (cs *CloudinaryService) GetResourceTypeFromExtension(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".webp", ".bmp", ".tiff":
		return "image"
	case ".mp4", ".avi", ".mov", ".wmv", ".flv":
		return "video"
	case ".pdf", ".doc", ".docx", ".txt":
		return "raw"
	default:
		return "auto"
	}
}

// HealthCheck فحص صحة خدمة Cloudinary
func (cs *CloudinaryService) HealthCheck() map[string]interface{} {
	cloudName := os.Getenv("CLOUDINARY_CLOUD_NAME")
	apiKey := os.Getenv("CLOUDINARY_API_KEY")
	
	if cloudName == "" || apiKey == "" {
		return map[string]interface{}{
			"service": "cloudinary",
			"status":  "error",
			"error":   "بيانات الاعتماد غير مكتملة",
		}
	}

	// محاولة رفع ملف تجريبي صغير
	testFile := strings.NewReader("test")
	_, err := cs.UploadImage(testFile, UploadOptions{
		PublicID: "health_check",
		Folder:   "nawthtech/tests",
	})

	if err != nil {
		return map[string]interface{}{
			"service": "cloudinary",
			"status":  "error",
			"error":   err.Error(),
		}
	}

	// حذف الملف التجريبي
	cs.DeleteImage("nawthtech/tests/health_check")

	return map[string]interface{}{
		"service":    "cloudinary",
		"status":     "healthy",
		"cloud_name": cloudName,
	}
}