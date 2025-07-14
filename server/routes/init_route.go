package routes

import (
	"time"

	"github.com/fiqrioemry/asset_management_system_app/server/handlers"

	"github.com/gin-gonic/gin"
)

func InitRoutes(r *gin.Engine, h *handlers.Handlers) {
	// Health check route (no API key required)
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":    "OK",
			"message":   "Server is healthy",
			"timestamp": time.Now().Format(time.RFC3339),
		})
	})

	v1 := r.Group("/api/v1")
	AuthRoutes(v1, h.UserHandler)
	UserRoutes(v1, h.UserHandler)
	CategoryRoutes(v1, h.CategoryHandler)
	AssetRoutes(v1, h.AssetHandler)
	LocationRoutes(v1, h.LocationHandler)
}
