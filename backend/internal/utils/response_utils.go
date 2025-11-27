package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response هيكل الاستجابة الموحد
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Code    string      `json:"code,omitempty"`
}

// SuccessResponse إرسال استجابة ناجحة
func SuccessResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	response := Response{
		Success: true,
		Message: message,
		Data:    data,
	}
	c.JSON(statusCode, response)
}

// ErrorResponse إرسال استجابة خطأ
func ErrorResponse(c *gin.Context, statusCode int, message string, errorCode string) {
	response := Response{
		Success: false,
		Message: message,
		Error:   message,
		Code:    errorCode,
	}
	c.JSON(statusCode, response)
}

// PaginationResponse إرسال استجابة مع التقسيم
func PaginationResponse(c *gin.Context, message string, data interface{}, pagination interface{}) {
	response := Response{
		Success: true,
		Message: message,
		Data: map[string]interface{}{
			"items":      data,
			"pagination": pagination,
		},
	}
	c.JSON(http.StatusOK, response)
}