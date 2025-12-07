package services

import (
	"context"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"strings"

	"github.com/nawthtech/nawthtech/backend/internal/ai/types"
)

// MediaService خدمة الوسائط المتعددة
type MediaService struct {
	imageProvider types.ImageProvider
	textProvider  types.TextProvider
}

// NewMediaService إنشاء خدمة وسائط جديدة
func NewMediaService(imageProvider types.ImageProvider, textProvider types.TextProvider) *MediaService {
	return &MediaService{
		imageProvider: imageProvider,
		textProvider:  textProvider,
	}
}

// AnalyzeImage تحليل صورة
func (s *MediaService) AnalyzeImage(ctx context.Context, imageData []byte, prompt string) (*types.AnalysisResponse, error) {
	if s.imageProvider == nil {
		return nil, fmt.Errorf("image provider not configured")
	}

	// استخراج معرف المستخدم باستخدام الدالة المساعدة
	userID := s.extractUserIDFromContext(ctx)

	// تحضير طلب التحليل
	req := types.AnalysisRequest{
		ImageData: imageData,
		Prompt:    prompt,
		UserID:    userID,
	}

	return s.imageProvider.AnalyzeImage(req)
}

// GetImageInfo الحصول على معلومات الصورة
func (s *MediaService) GetImageInfo(ctx context.Context, imageData []byte) (*types.ImageInfo, error) {
	// فك ترميز الصورة للحصول على المعلومات الأساسية
	reader := strings.NewReader(string(imageData))
	imgConfig, format, err := image.DecodeConfig(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	// الحصول على اسم نموذج الألوان
	colorModel := fmt.Sprintf("%T", imgConfig.ColorModel)

	return &types.ImageInfo{
		Width:      imgConfig.Width,
		Height:     imgConfig.Height,
		Format:     format,
		SizeBytes:  len(imageData),
		HasAlpha:   format == "png" || format == "gif",
		ColorModel: colorModel,
	}, nil
}

// ResizeImage تغيير حجم الصورة
func (s *MediaService) ResizeImage(ctx context.Context, imageData []byte, width, height int) ([]byte, error) {
	// تطبيق بسيط لتغيير حجم الصورة
	// في تطبيق حقيقي، يمكن استخدام مكتبة مثل "github.com/disintegration/imaging"

	// هنا نعيد نفس البيانات كتطبيق بسيط
	return imageData, nil
}

// CompressImage ضغط الصورة
func (s *MediaService) CompressImage(ctx context.Context, imageData []byte, quality int) ([]byte, error) {
	// تطبيق بسيط لضغط الصورة
	// في تطبيق حقيقي، يمكن استخدام مكتبة ضغط الصور

	// هنا نعيد نفس البيانات كتطبيق بسيط
	return imageData, nil
}

// ConvertImageFormat تحويل تنسيق الصورة
func (s *MediaService) ConvertImageFormat(ctx context.Context, imageData []byte, format string) ([]byte, error) {
	// تطبيق بسيط لتحويل التنسيق
	// في تطبيق حقيقي، يمكن استخدام مكتبة للتحويل

	// هنا نعيد نفس البيانات كتطبيق بسيط
	return imageData, nil
}

// GenerateImageFromText توليد صورة من نص
func (s *MediaService) GenerateImageFromText(ctx context.Context, prompt string, options types.ImageGenerationOptions) (*types.GeneratedImage, error) {
	if s.imageProvider == nil {
		return nil, fmt.Errorf("image generation not supported")
	}

	// استخراج معرف المستخدم
	userID := s.extractUserIDFromContext(ctx)

	// تحضير طلب توليد الصورة
	req := types.ImageRequest{
		Prompt:   prompt,
		Size:     options.Size,
		Style:    options.Style,
		Quality:  options.Quality,
		UserID:   userID,
	}

	resp, err := s.imageProvider.GenerateImage(req)
	if err != nil {
		return nil, err
	}

	// تحويل ImageResponse إلى GeneratedImage
	return &types.GeneratedImage{
		URL:       resp.URL,
		ImageData: resp.ImageData,
		Width:     resp.Width,
		Height:    resp.Height,
		Format:    resp.Format,
		SizeBytes: len(resp.ImageData),
		CreatedAt: resp.CreatedAt,
		Provider:  resp.ModelUsed,
		Model:     resp.ModelUsed,
		Cost:      resp.Cost,
		Prompt:    prompt,
		Seed:      resp.Seed,
	}, nil
}

// RemoveBackground إزالة خلفية الصورة
func (s *MediaService) RemoveBackground(ctx context.Context, imageData []byte) ([]byte, error) {
	// تطبيق بسيط لإزالة الخلفية
	// في تطبيق حقيقي، يمكن استخدام خدمة خارجية أو نموذج

	// هنا نعيد نفس البيانات كتطبيق بسيط
	return imageData, nil
}

// ApplyFilter تطبيق فلتر على الصورة
func (s *MediaService) ApplyFilter(ctx context.Context, imageData []byte, filterType string) ([]byte, error) {
	// تطبيق بسيط للفلتر
	// في تطبيق حقيقي، يمكن استخدام مكتبة معالجة الصور

	// هنا نعيد نفس البيانات كتطبيق بسيط
	return imageData, nil
}

// ExtractTextFromImage استخراج نص من صورة (OCR)
func (s *MediaService) ExtractTextFromImage(ctx context.Context, imageData []byte) (string, error) {
	// في تطبيق حقيقي، يمكن استخدام خدمة OCR
	// هنا نستخدم نموذج الذكاء الاصطناعي لتحليل الصورة واستخراج النص

	analysis, err := s.AnalyzeImage(ctx, imageData, "Extract all text from this image")
	if err != nil {
		return "", fmt.Errorf("failed to extract text: %w", err)
	}

	return analysis.Result, nil
}

// GetSupportedImageFormats الحصول على التنسيقات المدعومة
func (s *MediaService) GetSupportedImageFormats() []string {
	return []string{"jpeg", "jpg", "png", "gif", "webp", "bmp"}
}

// GetMaxImageSize الحصول على الحد الأقصى لحجم الصورة
func (s *MediaService) GetMaxImageSize() int64 {
	return 10 * 1024 * 1024 // 10MB
}

// ValidateImage التحقق من صحة الصورة
func (s *MediaService) ValidateImage(imageData []byte) error {
	// التحقق من الحجم
	if len(imageData) > int(s.GetMaxImageSize()) {
		return fmt.Errorf("image size exceeds maximum allowed size of 10MB")
	}

	// التحقق من التنسيق
	reader := strings.NewReader(string(imageData))
	_, format, err := image.DecodeConfig(reader)
	if err != nil {
		return fmt.Errorf("invalid image format: %w", err)
	}

	// التحقق من دعم التنسيق
	supported := s.GetSupportedImageFormats()
	format = strings.ToLower(format)
	for _, supportedFormat := range supported {
		if format == supportedFormat {
			return nil
		}
	}

	return fmt.Errorf("unsupported image format: %s", format)
}

// SaveImageToStorage حفظ الصورة في التخزين
func (s *MediaService) SaveImageToStorage(ctx context.Context, imageData []byte, filename string) (string, error) {
	// في تطبيق حقيقي، يمكن حفظ الصورة في التخزين السحابي
	// هنا نعيد معرف بسيط
	return "image_" + filename, nil
}

// LoadImageFromStorage تحميل الصورة من التخزين
func (s *MediaService) LoadImageFromStorage(ctx context.Context, imageID string) ([]byte, error) {
	// في تطبيق حقيقي، يمكن تحميل الصورة من التخزين السحابي
	// هنا نعيد بيانات فارغة
	return []byte{}, nil
}

// GetImageStats الحصول على إحصائيات الصور
func (s *MediaService) GetImageStats(ctx context.Context) map[string]interface{} {
	stats := make(map[string]interface{})

	stats["service"] = "media"
	stats["supported_formats"] = s.GetSupportedImageFormats()
	stats["max_size_mb"] = s.GetMaxImageSize() / (1024 * 1024)

	return stats
}

// BatchProcessImages معالجة دفعة من الصور
func (s *MediaService) BatchProcessImages(ctx context.Context, images [][]byte, processFunc func([]byte) ([]byte, error)) ([][]byte, error) {
	var results [][]byte

	for _, imageData := range images {
		processed, err := processFunc(imageData)
		if err != nil {
			return results, fmt.Errorf("failed to process image: %w", err)
		}
		results = append(results, processed)
	}

	return results, nil
}

// CreateThumbnail إنشاء صورة مصغرة
func (s *MediaService) CreateThumbnail(ctx context.Context, imageData []byte, maxWidth, maxHeight int) ([]byte, error) {
	// الحصول على معلومات الصورة
	info, err := s.GetImageInfo(ctx, imageData)
	if err != nil {
		return nil, err
	}

	// حساب الأبعاد الجديدة مع الحفاظ على النسبة
	width := info.Width
	height := info.Height

	if width > maxWidth || height > maxHeight {
		ratio := float64(width) / float64(height)

		if width > maxWidth {
			width = maxWidth
			height = int(float64(width) / ratio)
		}

		if height > maxHeight {
			height = maxHeight
			width = int(float64(height) * ratio)
		}
	}

	// تغيير حجم الصورة
	return s.ResizeImage(ctx, imageData, width, height)
}

// ReadImageFromReader قراءة صورة من قارئ
func (s *MediaService) ReadImageFromReader(ctx context.Context, reader io.Reader) ([]byte, error) {
	return io.ReadAll(reader)
}

// WriteImageToWriter كتابة صورة إلى كاتب
func (s *MediaService) WriteImageToWriter(ctx context.Context, imageData []byte, writer io.Writer) error {
	_, err := writer.Write(imageData)
	return err
}

// ==================== دالة مساعدة خاصة ====================

// extractUserIDFromContext استخراج معرف المستخدم من السياق
func (s *MediaService) extractUserIDFromContext(ctx context.Context) string {
	// في تطبيق حقيقي، يمكن استخراج معرف المستخدم من السياق
	// هنا نعيد قيمة افتراضية
	if userID, ok := ctx.Value("userID").(string); ok {
		return userID
	}
	return "anonymous"
}