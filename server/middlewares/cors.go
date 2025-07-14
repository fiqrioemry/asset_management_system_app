package middlewares

import (
	"slices"
	"strings"

	"github.com/fiqrioemry/asset_management_system_app/server/config"
	"github.com/fiqrioemry/go-api-toolkit/response"
	"github.com/gin-gonic/gin"
)

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		var allowedOrigin string

		if config.AppConfig.AppEnv == "production" {
			allowedOrigin = getProductionOrigin(origin)
			if origin != "" && allowedOrigin == "" {
				err := response.NewForbidden("Origin not allowed by CORS policy")
				response.Error(c, err)
				return
			}
		} else {
			allowedOrigin = getDevelopmentOrigin(origin)
		}

		// Set CORS headers
		c.Writer.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-API-Key")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func getProductionOrigin(origin string) string {
	if slices.Contains(config.AppConfig.AllowedOrigins, origin) {
		return origin
	}

	if len(config.AppConfig.AllowedOrigins) > 0 {
		return config.AppConfig.AllowedOrigins[0]
	}

	return ""
}

func getDevelopmentOrigin(origin string) string {
	if isLocalhost(origin) {
		return origin
	}

	if slices.Contains(config.AppConfig.AllowedOrigins, origin) {
		return origin
	}

	if len(config.AppConfig.AllowedOrigins) > 0 {
		return config.AppConfig.AllowedOrigins[0]
	}

	return "*"
}

func isLocalhost(origin string) bool {
	if origin == "" {
		return false
	}

	return strings.HasPrefix(origin, "http://localhost") ||
		strings.HasPrefix(origin, "https://localhost") ||
		strings.HasPrefix(origin, "http://127.0.0.1") ||
		strings.HasPrefix(origin, "https://127.0.0.1")
}
