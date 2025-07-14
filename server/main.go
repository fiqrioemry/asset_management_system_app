package main

import (
	"log"

	"github.com/fiqrioemry/asset_management_system_app/server/config"
	"github.com/fiqrioemry/asset_management_system_app/server/handlers"
	"github.com/fiqrioemry/asset_management_system_app/server/middlewares"
	"github.com/fiqrioemry/asset_management_system_app/server/repositories"
	"github.com/fiqrioemry/asset_management_system_app/server/routes"
	"github.com/fiqrioemry/asset_management_system_app/server/services"
	"github.com/fiqrioemry/asset_management_system_app/server/utils"
	"github.com/fiqrioemry/go-api-toolkit/response"

	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
)

// ASSET MANAGEMENT APP SERVER
// VERSION: 1.0.0
// DEPLOYMENT: docker-compose
// PORT: 5005
// DESCRIPTION: This is a server for an asset management system that handles user registration, asset management, and payment processing.

func main() {
	// ========== Configuration =================
	config.InitConfiguration()
	utils.InitLogger()
	db := config.DB

	// seeders.ResetDatabase(db)

	// ========== Initialize response toolkit ===
	// This initializes the response toolkit with custom configurations
	// such as logging success and error responses.
	response.InitGin(response.InitConfig{
		Logger:              utils.GetLogger(),
		LogSuccessResponses: false,
		LogErrorResponses:   true,
	})

	// ========== Initialize layer ============
	repo := repositories.InitRepositories(db)
	s := services.InitServices(repo)
	h := handlers.InitHandlers(s)

	// ========== Initialize gin engine =======
	r := gin.Default()
	r.SetTrustedProxies(config.AppConfig.TrustedProxies)

	// ========== Initialize Middleware ========
	r.Use(
		ginzap.Ginzap(utils.GetLogger(), time.RFC3339, true),
		middlewares.Recovery(),
		middlewares.CORS(),
		middlewares.RateLimiterInit(),
		middlewares.LimitFileSize(config.AppConfig.MaxFileSize),
		middlewares.APIKeyGateway(config.AppConfig.SkippedApiEndpoints),
	)

	// ========== Initialize routes ===========
	routes.InitRoutes(r, h)

	port := config.AppConfig.ServerPort
	log.Println("server running on port:", port)
	log.Fatal(r.Run(":" + port))
}
