// internal/routes/auth_route.go
package routes

import (
	"github.com/fiqrioemry/asset_management_system_app/server/handlers"
	"github.com/fiqrioemry/asset_management_system_app/server/middlewares"

	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.RouterGroup, h *handlers.UserHandler) {
	user := r.Group("/users")
	user.POST("/login", h.Login)
	user.POST("/register", h.Register)
	user.POST("/logout", h.Logout)
	user.POST("/refresh-token", h.RefreshSession)

	user.Use(middlewares.AuthRequired())
	user.GET("/me", h.GetMe)
	user.PUT("/me", h.UpdateMe)
}
