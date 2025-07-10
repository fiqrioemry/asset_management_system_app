// routes/auth_routes.go
package routes

import (
	"github.com/fiqrioemry/asset_management_system_app/server/handlers"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.RouterGroup, h *handlers.UserHandler) {
	auth := r.Group("/auth")
	{
		// Authentication endpoints
		auth.POST("/login", h.Login)
		auth.POST("/register", h.Register)
		auth.POST("/logout", h.Logout)
		auth.POST("/refresh-token", h.RefreshSession)

		// Password reset flow
		auth.POST("/forgot-password", h.ForgotPassword)
		auth.GET("/validate-reset-token", h.ValidateResetToken)
		auth.POST("/reset-password", h.ResetPassword)
	}
}
