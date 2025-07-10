package middlewares

import (
	"strings"

	"github.com/fiqrioemry/asset_management_system_app/server/config"
	"github.com/gin-gonic/gin"
)

func APIKeyGateway(skippedPaths []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		currentPath := c.Request.URL.Path

		for _, path := range skippedPaths {
			if currentPath == path || strings.HasPrefix(currentPath, path) {
				c.Next()
				return
			}
		}

		apiKey := c.GetHeader("X-API-KEY")
		if apiKey == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized - API key is required"})
			return
		}

		validKeys := strings.Split(config.AppConfig.ApiKeys, ",")
		isValid := false
		for _, validKey := range validKeys {
			if strings.TrimSpace(validKey) == apiKey {
				isValid = true
				break
			}
		}

		if !isValid {
			c.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized - invalid API key"})
			return
		}

		c.Next()
	}
}
