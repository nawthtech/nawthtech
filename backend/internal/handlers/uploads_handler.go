package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nawthtech/nawthtech/backend/internal/middleware"
	"github.com/nawthtech/nawthtech/backend/internal/models"
	"github.com/nawthtech/nawthtech/backend/internal/services"
	"github.com/nawthtech/nawthtech/backend/internal/utils"
)

type UploadsHandler struct {
	uploadsService services.UploadsService
	authService    services.AuthService
}

func NewUploadsHandler(uploadsService services.UploadsService, authService services.AuthService) *UploadsHandler {
	return &UploadsHandler{
		uploadsService: uploadsService,
		authService:    authService,
	}
}

// UploadFile - رفع ملف واحد
// @Summary رفع ملف واحد
// @Description رفع ملف واحد
// @Tags Uploads
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "الملف المطلوب رفعه"
// @Success 200 {object} utils.Response
// @Router /api/v1/uploads [post]
func (h *UploadsHandler) UploadFile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "غير مصرح", "UNAUTHORIZED")
		return
	}

	if !middleware.CheckRateLimit(c, "uploads_single", 20, 15*time.Minute) {
		utils.ErrorResponse(c, http.StatusTooManyRequests, "تم تجاوز الحد المسموح", "RATE_LIMIT_EXCEEDED")
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "لم يتم توفير ملف", "NO_FILE_PROVIDED")
		return
	}

	result, err := h.uploadsService.UploadFile(c, services.UploadFileParams{
		File:   file,
		UserID: userID.(string),
	})

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في رفع الملف", "UPLOAD_FAILED")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "تم رفع الملف بنجاح", result)
}

// UploadMultipleFiles - رفع ملفات متعددة
// @Summary رفع ملفات متعددة
// @Description رفع ملفات متعددة
// @Tags Uploads
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param files formData file true "الملفات المطلوب رفعها"
// @Success 200 {object} utils.Response
// @Router /api/v1/uploads/multiple [post]
func (h *UploadsHandler) UploadMultipleFiles(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "غير مصرح", "UNAUTHORIZED")
		return
	}

	if !middleware.CheckRateLimit(c, "uploads_multiple", 10, 30*time.Minute) {
		utils.ErrorResponse(c, http.StatusTooManyRequests, "تم تجاوز الحد المسموح", "RATE_LIMIT_EXCEEDED")
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "بيانات غير صالحة", "INVALID_FORM_DATA")
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		utils.ErrorResponse(c, http.StatusBadRequest, "لم يتم توفير ملفات", "NO_FILES_PROVIDED")
		return
	}

	if len(files) > 5 {
		utils.ErrorResponse(c, http.StatusBadRequest, "الحد الأقصى للملفات هو 5", "MAX_FILES_EXCEEDED")
		return
	}

	result, err := h.uploadsService.UploadMultipleFiles(c, services.UploadMultipleFilesParams{
		Files:  files,
		UserID: userID.(string),
	})

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في رفع الملفات", "UPLOAD_FAILED")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "تم رفع الملفات بنجاح", result)
}

// BulkUploadFilesRequest - طلب رفع كمي للملفات
type BulkUploadFilesRequest struct {
	Purpose string `json:"purpose" binding:"required"`
}

// BulkUploadFiles - رفع كمي للملفات (للمتجر)
// @Summary رفع كمي للملفات (للمتجر)
// @Description رفع كمي للملفات (للمتجر)
// @Tags Uploads
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param files formData file true "الملفات المطلوب رفعها"
// @Param purpose formData string true "الغرض من الرفع"
// @Success 200 {object} utils.Response
// @Router /api/v1/uploads/bulk [post]
func (h *UploadsHandler) BulkUploadFiles(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "غير مصرح", "UNAUTHORIZED")
		return
	}

	if !middleware.CheckRateLimit(c, "uploads_bulk", 5, time.Hour) {
		utils.ErrorResponse(c, http.StatusTooManyRequests, "تم تجاوز الحد المسموح", "RATE_LIMIT_EXCEEDED")
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "بيانات غير صالحة", "INVALID_FORM_DATA")
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		utils.ErrorResponse(c, http.StatusBadRequest, "لم يتم توفير ملفات", "NO_FILES_PROVIDED")
		return
	}

	if len(files) > 20 {
		utils.ErrorResponse(c, http.StatusBadRequest, "الحد الأقصى للملفات هو 20", "MAX_FILES_EXCEEDED")
		return
	}

	purpose := c.PostForm("purpose")
	if purpose == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "الغرض من الرفع مطلوب", "PURPOSE_REQUIRED")
		return
	}

	result, err := h.uploadsService.BulkUploadFiles(c, services.BulkUploadFilesParams{
		Files:   files,
		UserID:  userID.(string),
		Purpose: purpose,
	})

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في الرفع الكمي", "BULK_UPLOAD_FAILED")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "تم الرفع الكمي بنجاح", result)
}

// UploadServiceImageRequest - طلب رفع صورة خدمة
type UploadServiceImageRequest struct {
	ServiceID string `json:"serviceId" binding:"required"`
	IsCover   bool   `json:"isCover"`
}

// UploadServiceImage - رفع صورة خدمة (صورة رئيسية)
// @Summary رفع صورة خدمة (صورة رئيسية)
// @Description رفع صورة خدمة (صورة رئيسية)
// @Tags Uploads
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param image formData file true "صورة الخدمة"
// @Param serviceId formData string true "معرف الخدمة"
// @Param isCover formData bool false "هل هي صورة الغلاف"
// @Success 200 {object} utils.Response
// @Router /api/v1/uploads/services/images [post]
func (h *UploadsHandler) UploadServiceImage(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "غير مصرح", "UNAUTHORIZED")
		return
	}

	if !middleware.CheckRateLimit(c, "uploads_service_image", 15, 30*time.Minute) {
		utils.ErrorResponse(c, http.StatusTooManyRequests, "تم تجاوز الحد المسموح", "RATE_LIMIT_EXCEEDED")
		return
	}

	file, err := c.FormFile("image")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "لم يتم توفير صورة", "NO_IMAGE_PROVIDED")
		return
	}

	serviceID := c.PostForm("serviceId")
	if serviceID == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "معرف الخدمة مطلوب", "SERVICE_ID_REQUIRED")
		return
	}

	isCover := c.PostForm("isCover") == "true"

	result, err := h.uploadsService.UploadServiceImage(c, services.UploadServiceImageParams{
		File:      file,
		UserID:    userID.(string),
		ServiceID: serviceID,
		IsCover:   isCover,
	})

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في رفع صورة الخدمة", "SERVICE_IMAGE_UPLOAD_FAILED")
		return
	}

	message := "تم رفع صورة الخدمة بنجاح"
	if isCover {
		message = "تم رفع صورة الغلاف بنجاح"
	}

	utils.SuccessResponse(c, http.StatusOK, message, result)
}

// UploadServiceGalleryRequest - طلب رفع معرض صور للخدمة
type UploadServiceGalleryRequest struct {
	ServiceID string `json:"serviceId" binding:"required"`
}

// UploadServiceGallery - رفع معرض صور للخدمة
// @Summary رفع معرض صور للخدمة
// @Description رفع معرض صور للخدمة
// @Tags Uploads
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param gallery formData file true "صور المعرض"
// @Param serviceId formData string true "معرف الخدمة"
// @Success 200 {object} utils.Response
// @Router /api/v1/uploads/services/gallery [post]
func (h *UploadsHandler) UploadServiceGallery(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "غير مصرح", "UNAUTHORIZED")
		return
	}

	if !middleware.CheckRateLimit(c, "uploads_service_gallery", 10, 30*time.Minute) {
		utils.ErrorResponse(c, http.StatusTooManyRequests, "تم تجاوز الحد المسموح", "RATE_LIMIT_EXCEEDED")
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "بيانات غير صالحة", "INVALID_FORM_DATA")
		return
	}

	files := form.File["gallery"]
	if len(files) == 0 {
		utils.ErrorResponse(c, http.StatusBadRequest, "لم يتم توفير صور", "NO_IMAGES_PROVIDED")
		return
	}

	if len(files) > 10 {
		utils.ErrorResponse(c, http.StatusBadRequest, "الحد الأقصى للصور هو 10", "MAX_IMAGES_EXCEEDED")
		return
	}

	serviceID := c.PostForm("serviceId")
	if serviceID == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "معرف الخدمة مطلوب", "SERVICE_ID_REQUIRED")
		return
	}

	result, err := h.uploadsService.UploadServiceGallery(c, services.UploadServiceGalleryParams{
		Files:     files,
		UserID:    userID.(string),
		ServiceID: serviceID,
	})

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في رفع معرض الصور", "GALLERY_UPLOAD_FAILED")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "تم رفع معرض الصور بنجاح", result)
}

// UploadStoreAssetRequest - طلب رفع أصول المتجر
type UploadStoreAssetRequest struct {
	AssetType string `json:"assetType" binding:"required"`
	StoreID   string `json:"storeId" binding:"required"`
}

// UploadStoreAsset - رفع أصول المتجر (لوجو، بانرات، إلخ)
// @Summary رفع أصول المتجر (لوجو، بانرات، إلخ)
// @Description رفع أصول المتجر (لوجو، بانرات، إلخ)
// @Tags Uploads
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param asset formData file true "أصل المتجر"
// @Param assetType formData string true "نوع الأصل"
// @Param storeId formData string true "معرف المتجر"
// @Success 200 {object} utils.Response
// @Router /api/v1/uploads/store/assets [post]
func (h *UploadsHandler) UploadStoreAsset(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "غير مصرح", "UNAUTHORIZED")
		return
	}

	if !middleware.CheckRateLimit(c, "uploads_store_asset", 10, time.Hour) {
		utils.ErrorResponse(c, http.StatusTooManyRequests, "تم تجاوز الحد المسموح", "RATE_LIMIT_EXCEEDED")
		return
	}

	file, err := c.FormFile("asset")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "لم يتم توفير ملف", "NO_FILE_PROVIDED")
		return
	}

	assetType := c.PostForm("assetType")
	if assetType == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "نوع الأصل مطلوب", "ASSET_TYPE_REQUIRED")
		return
	}

	storeID := c.PostForm("storeId")
	if storeID == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "معرف المتجر مطلوب", "STORE_ID_REQUIRED")
		return
	}

	result, err := h.uploadsService.UploadStoreAsset(c, services.UploadStoreAssetParams{
		File:      file,
		UserID:    userID.(string),
		AssetType: assetType,
		StoreID:   storeID,
	})

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في رفع أصل المتجر", "STORE_ASSET_UPLOAD_FAILED")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "تم رفع أصل المتجر بنجاح", result)
}

// UploadAndAnalyzeRequest - طلب رفع وتحليل بالذكاء الاصطناعي
type UploadAndAnalyzeRequest struct {
	AnalysisType string                 `json:"analysisType"`
	Options      map[string]interface{} `json:"options"`
}

// UploadAndAnalyze - رفع وتحليل بالذكاء الاصطناعي
// @Summary رفع وتحليل بالذكاء الاصطناعي
// @Description رفع وتحليل بالذكاء الاصطناعي
// @Tags Uploads
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "الملف المطلوب تحليله"
// @Param analysisType formData string false "نوع التحليل"
// @Param options formData string false "خيارات التحليل"
// @Success 200 {object} utils.Response
// @Router /api/v1/uploads/ai-analysis [post]
func (h *UploadsHandler) UploadAndAnalyze(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "غير مصرح", "UNAUTHORIZED")
		return
	}

	if !middleware.CheckRateLimit(c, "uploads_ai_analysis", 10, time.Hour) {
		utils.ErrorResponse(c, http.StatusTooManyRequests, "تم تجاوز الحد المسموح", "RATE_LIMIT_EXCEEDED")
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "لم يتم توفير ملف", "NO_FILE_PROVIDED")
		return
	}

	analysisType := c.PostForm("analysisType")
	options := make(map[string]interface{})
	// يمكن تحليل options من JSON string إذا كانت مرسلة

	result, err := h.uploadsService.UploadAndAnalyze(c, services.UploadAndAnalyzeParams{
		File:         file,
		UserID:       userID.(string),
		AnalysisType: analysisType,
		Options:      options,
	})

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في رفع الملف والتحليل", "AI_ANALYSIS_FAILED")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "تم رفع الملف والتحليل بنجاح", result)
}

// OptimizeImageRequest - طلب تحسين صورة
type OptimizeImageRequest struct {
	Quality int                    `json:"quality"`
	Format  string                 `json:"format"`
	Resize  map[string]interface{} `json:"resize"`
}

// OptimizeImage - تحسين صورة
// @Summary تحسين صورة
// @Description تحسين صورة
// @Tags Uploads
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param fileId path string true "معرف الملف"
// @Param input body OptimizeImageRequest true "بيانات التحسين"
// @Success 200 {object} utils.Response
// @Router /api/v1/uploads/optimize/{fileId} [post]
func (h *UploadsHandler) OptimizeImage(c *gin.Context) {
	var req OptimizeImageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "بيانات غير صالحة", "INVALID_INPUT")
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "غير مصرح", "UNAUTHORIZED")
		return
	}

	fileID := c.Param("fileId")

	result, err := h.uploadsService.OptimizeImage(c, services.OptimizeImageParams{
		FileID:  fileID,
		UserID:  userID.(string),
		Quality: req.Quality,
		Format:  req.Format,
		Resize:  req.Resize,
	})

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في تحسين الصورة", "IMAGE_OPTIMIZATION_FAILED")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "تم تحسين الصورة بنجاح", result)
}

// GetFileInfo - الحصول على معلومات الملف
// @Summary الحصول على معلومات الملف
// @Description الحصول على معلومات الملف
// @Tags Uploads
// @Security BearerAuth
// @Produce json
// @Param fileId path string true "معرف الملف"
// @Success 200 {object} utils.Response
// @Router /api/v1/uploads/info/{fileId} [get]
func (h *UploadsHandler) GetFileInfo(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "غير مصرح", "UNAUTHORIZED")
		return
	}

	fileID := c.Param("fileId")

	fileInfo, err := h.uploadsService.GetFileInfo(c, fileID, userID.(string))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "الملف غير موجود", "FILE_NOT_FOUND")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "تم جلب معلومات الملف بنجاح", fileInfo)
}

// GetUserFiles - الحصول على ملفات المستخدم
// @Summary الحصول على ملفات المستخدم
// @Description الحصول على ملفات المستخدم
// @Tags Uploads
// @Security BearerAuth
// @Produce json
// @Param page query int false "الصفحة" default(1)
// @Param limit query int false "الحد" default(20)
// @Param type query string false "نوع الملف"
// @Success 200 {object} utils.Response
// @Router /api/v1/uploads/user [get]
func (h *UploadsHandler) GetUserFiles(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "غير مصرح", "UNAUTHORIZED")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	fileType := c.Query("type")

	files, pagination, err := h.uploadsService.GetUserFiles(c, services.GetUserFilesParams{
		UserID: userID.(string),
		Page:   page,
		Limit:  limit,
		Type:   fileType,
	})

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في جلب ملفات المستخدم", "FILES_FETCH_FAILED")
		return
	}

	response := map[string]interface{}{
		"files":      files,
		"pagination": pagination,
	}

	utils.SuccessResponse(c, http.StatusOK, "تم جلب ملفات المستخدم بنجاح", response)
}

// UpdateFileInfoRequest - طلب تحديث معلومات الملف
type UpdateFileInfoRequest struct {
	FileName string                 `json:"fileName"`
	Metadata map[string]interface{} `json:"metadata"`
	IsPublic bool                   `json:"isPublic"`
}

// UpdateFileInfo - تحديث معلومات الملف
// @Summary تحديث معلومات الملف
// @Description تحديث معلومات الملف
// @Tags Uploads
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param fileId path string true "معرف الملف"
// @Param input body UpdateFileInfoRequest true "بيانات التحديث"
// @Success 200 {object} utils.Response
// @Router /api/v1/uploads/{fileId} [put]
func (h *UploadsHandler) UpdateFileInfo(c *gin.Context) {
	var req UpdateFileInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "بيانات غير صالحة", "INVALID_INPUT")
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "غير مصرح", "UNAUTHORIZED")
		return
	}

	fileID := c.Param("fileId")

	updatedFile, err := h.uploadsService.UpdateFileInfo(c, services.UpdateFileInfoParams{
		FileID:   fileID,
		UserID:   userID.(string),
		FileName: req.FileName,
		Metadata: req.Metadata,
		IsPublic: req.IsPublic,
	})

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في تحديث معلومات الملف", "FILE_UPDATE_FAILED")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "تم تحديث معلومات الملف بنجاح", updatedFile)
}

// DeleteFile - حذف ملف
// @Summary حذف ملف
// @Description حذف ملف
// @Tags Uploads
// @Security BearerAuth
// @Produce json
// @Param fileId path string true "معرف الملف"
// @Success 200 {object} utils.Response
// @Router /api/v1/uploads/{fileId} [delete]
func (h *UploadsHandler) DeleteFile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "غير مصرح", "UNAUTHORIZED")
		return
	}

	fileID := c.Param("fileId")

	err := h.uploadsService.DeleteFile(c, fileID, userID.(string))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في حذف الملف", "FILE_DELETE_FAILED")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "تم حذف الملف بنجاح", nil)
}

// GetStorageUsage - الحصول على استخدام التخزين
// @Summary الحصول على استخدام التخزين
// @Description الحصول على استخدام التخزين
// @Tags Uploads
// @Security BearerAuth
// @Produce json
// @Success 200 {object} utils.Response
// @Router /api/v1/uploads/storage/usage [get]
func (h *UploadsHandler) GetStorageUsage(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "غير مصرح", "UNAUTHORIZED")
		return
	}

	usage, err := h.uploadsService.GetStorageUsage(c, userID.(string))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في جلب معلومات استخدام التخزين", "STORAGE_USAGE_FAILED")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "تم جلب معلومات استخدام التخزين بنجاح", usage)
}

// GetUploadStats - إحصائيات الرفع (للمدراء فقط)
// @Summary إحصائيات الرفع (للمدراء فقط)
// @Description إحصائيات الرفع (للمدراء فقط)
// @Tags Uploads
// @Security BearerAuth
// @Produce json
// @Param period query string false "الفترة" default(30d)
// @Param type query string false "النوع" default(overview)
// @Success 200 {object} utils.Response
// @Router /api/v1/uploads/stats [get]
func (h *UploadsHandler) GetUploadStats(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "غير مصرح", "UNAUTHORIZED")
		return
	}

	period := c.DefaultQuery("period", "30d")
	statsType := c.DefaultQuery("type", "overview")

	stats, err := h.uploadsService.GetUploadStats(c, services.GetUploadStatsParams{
		UserID: userID.(string),
		Period: period,
		Type:   statsType,
	})

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في جلب إحصائيات الرفع", "STATS_FETCH_FAILED")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "تم جلب إحصائيات الرفع بنجاح", stats)
}

// CleanupFilesRequest - طلب تنظيف الملفات المؤقتة
type CleanupFilesRequest struct {
	OlderThan string `json:"olderThan"`
	Type      string `json:"type"`
}

// CleanupFiles - تنظيف الملفات المؤقتة (للمدراء فقط)
// @Summary تنظيف الملفات المؤقتة (للمدراء فقط)
// @Description تنظيف الملفات المؤقتة (للمدراء فقط)
// @Tags Uploads
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param input body CleanupFilesRequest true "بيانات التنظيف"
// @Success 200 {object} utils.Response
// @Router /api/v1/uploads/cleanup [post]
func (h *UploadsHandler) CleanupFiles(c *gin.Context) {
	var req CleanupFilesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "بيانات غير صالحة", "INVALID_INPUT")
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "غير مصرح", "UNAUTHORIZED")
		return
	}

	result, err := h.uploadsService.CleanupFiles(c, services.CleanupFilesParams{
		UserID:   userID.(string),
		OlderThan: req.OlderThan,
		Type:     req.Type,
	})

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في تنظيف الملفات", "CLEANUP_FAILED")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "تم تنظيف الملفات بنجاح", result)
}