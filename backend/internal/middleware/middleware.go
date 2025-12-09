package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nawthtech/nawthtech/backend/internal/config"
	"github.com/nawthtech/nawthtech/backend/internal/utils"
)

// Simple CORS
func CORSMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		allowed := false
		for _, o := range cfg.Cors.AllowedOrigins {
			if o == origin || o == "*" {
				allowed = true
				break
			}
		}
		if allowed {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Methods", strings.Join(cfg.Cors.AllowedMethods, ","))
			c.Header("Access-Control-Allow-Headers", strings.Join(cfg.Cors.AllowedHeaders, ","))
			if cfg.Cors.AllowCredentials {
				c.Header("Access-Control-Allow-Credentials", "true")
			}
		}
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

// AuthMiddleware verifies JWT and sets user id in context
func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"success": false, "error": "UNAUTHORIZED"})
			return
		}
		token := strings.TrimPrefix(auth, "Bearer ")
		claims, err := utils.VerifyJWT(cfg, token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"success": false, "error": "INVALID_TOKEN"})
			return
		}
		// set user in context
		c.Set("user_id", claims.UserID)
		c.Next()
	}
}