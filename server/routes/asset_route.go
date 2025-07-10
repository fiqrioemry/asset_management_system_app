// routes/asset_routes.go
package routes

import (
	"github.com/fiqrioemry/asset_management_system_app/server/handlers"
	"github.com/fiqrioemry/asset_management_system_app/server/middlewares"
	"github.com/gin-gonic/gin"
)

func AssetRoutes(router *gin.RouterGroup, assetHandler *handlers.AssetHandler) {
	// Asset routes with authentication middleware
	assetRoutes := router.Group("/assets")
	assetRoutes.Use(middlewares.AuthRequired())
	{
		assetRoutes.GET("", assetHandler.GetAssets)
		assetRoutes.POST("", assetHandler.CreateAsset)
		assetRoutes.GET("/:id", assetHandler.GetAssetByID)
		assetRoutes.PUT("/:id", assetHandler.UpdateAsset)
		assetRoutes.DELETE("/:id", assetHandler.DeleteAsset)
	}
}
