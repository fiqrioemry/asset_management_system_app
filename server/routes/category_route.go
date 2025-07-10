package routes

import (
	"github.com/fiqrioemry/asset_management_system_app/server/handlers"
	"github.com/fiqrioemry/asset_management_system_app/server/middlewares"
	"github.com/gin-gonic/gin"
)

func CategoryRoutes(r *gin.RouterGroup, h *handlers.CategoryHandler) {
	categories := r.Group("/categories")
	categories.Use(middlewares.AuthRequired())
	{
		categories.GET("/tree", h.GetCategoriesTree)
		categories.GET("/flat", h.GetCategoriesFlat)
		categories.GET("/parents", h.GetParentCategories) // selected

		categories.GET("/:id", h.GetCategoryByID)
		categories.GET("/:id/children", h.GetChildCategories)
		categories.GET("/:id/assets", h.GetAssetsByCategory)

		categories.POST("/", h.CreateCategory)
		categories.PUT("/:id", h.UpdateCategory)
		categories.DELETE("/:id", h.DeleteCategory)
	}
}
