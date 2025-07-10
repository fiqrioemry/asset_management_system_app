// routes/location_routes.go
package routes

import (
	"github.com/fiqrioemry/asset_management_system_app/server/handlers"
	"github.com/fiqrioemry/asset_management_system_app/server/middlewares"
	"github.com/gin-gonic/gin"
)

func LocationRoutes(r *gin.RouterGroup, h *handlers.LocationHandler) {
	locations := r.Group("/locations")
	locations.Use(middlewares.AuthRequired())
	{
		locations.GET("/", h.GetLocations)                  // GET /api/v1/locations
		locations.POST("/", h.CreateLocation)               // POST /api/v1/locations
		locations.GET("/:id", h.GetLocationByID)            // GET /api/v1/locations/:id
		locations.PUT("/:id", h.UpdateLocation)             // PUT /api/v1/locations/:id
		locations.DELETE("/:id", h.DeleteLocation)          // DELETE /api/v1/locations/:id
		locations.GET("/:id/assets", h.GetAssetsByLocation) // GET /api/v1/locations/:id/assets
	}
}
