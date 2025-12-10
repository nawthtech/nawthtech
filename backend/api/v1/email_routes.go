// backend/api/v1/email_routes.go
package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/nawthtech/nawthtech/backend/internal/handlers"
)

func RegisterEmailRoutes(router *gin.RouterGroup) {
	emailHandler, err := handlers.NewEmailHandler()
	if err != nil {
		// Log error but don't crash
		return
	}

	emailGroup := router.Group("/email")
	{
		emailGroup.POST("/setup", emailHandler.SetupEmail)
		emailGroup.GET("/allow-list", emailHandler.GetAllowList)
		emailGroup.POST("/allow-list", emailHandler.AddToAllowList)
		emailGroup.DELETE("/allow-list/:email", emailHandler.RemoveFromAllowList)
		emailGroup.POST("/test", emailHandler.TestEmail)
	}
}
