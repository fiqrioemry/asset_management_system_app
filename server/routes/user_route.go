// internal/routes/auth_route.go
package routes

import (
	"github.com/fiqrioemry/asset_management_system_app/server/handlers"
	"github.com/fiqrioemry/asset_management_system_app/server/middlewares"

	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.RouterGroup, h *handlers.UserHandler) {
	users := r.Group("/users")
	users.Use(middlewares.AuthRequired())
	{
		users.GET("/me", h.GetMe)
		users.PUT("/me", h.UpdateMe)
		users.POST("/change-password", h.ChangePassword)
	}
}
