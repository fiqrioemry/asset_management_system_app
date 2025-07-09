package routes

import (
	"github.com/fiqrioemry/asset_management_system_app/server/handlers"

	"github.com/gin-gonic/gin"
)

func InitRoutes(r *gin.Engine, h *handlers.Handlers) {
	api := r.Group("/api/v1")

	// ========= Authentication & User Management ========
	UserRoutes(api, h.UserHandler)

}
