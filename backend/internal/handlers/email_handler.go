package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nawthtech/nawthtech/backend/internal/email"
)

type EmailHandler struct {
	service *email.Service
}

func NewEmailHandler() (*EmailHandler, error) {
	service, err := email.NewService()
	if err != nil {
		return nil, err
	}
	return &EmailHandler{service: service}, nil
}

func (h *EmailHandler) SetupEmail(c *gin.Context) {
	if err := h.service.SetupEmailRouting(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to setup email",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Email routing setup initiated",
	})
}

func (h *EmailHandler) GetAllowList(c *gin.Context) {
	emails := h.service.GetAllowList()
	c.JSON(http.StatusOK, gin.H{
		"emails": emails,
	})
}

// ... other handler methods
