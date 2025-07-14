package middlewares

import (
	"strings"

	"github.com/fiqrioemry/asset_management_system_app/server/config"
	"github.com/fiqrioemry/go-api-toolkit/response"
	"github.com/gin-gonic/gin"
)

func APIKeyGateway(skippedPaths []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		currentPath := c.Request.URL.Path

		for _, path := range skippedPaths {
			if path == "/" {
				if currentPath == "/" {
					c.Next()
					return
				}
				continue
			}

			if currentPath == path {
				c.Next()
				return
			}

			if len(path) > 1 && strings.HasPrefix(currentPath, path) {
				c.Next()
				return
			}
		}

		apiKey := c.GetHeader("X-API-KEY")
		if apiKey == "" {
			response.Error(c, response.NewUnauthorized("Unauthorized - API key is required"))
			c.Abort()
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
			response.Error(c, response.NewUnauthorized("Unauthorized - invalid API key"))
			c.Abort()
			return
		}

		c.Next()
	}
}
