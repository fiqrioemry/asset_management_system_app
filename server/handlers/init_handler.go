package handlers

import (
	"github.com/fiqrioemry/asset_management_system_app/server/services"
)

type Handlers struct {
	UserHandler     *UserHandler
	AssetHandler    *AssetHandler
	LocationHandler *LocationHandler
	CategoryHandler *CategoryHandler
	// 	DashboardHandler *DashboardHandler
	//
}

func InitHandlers(s *services.Services) *Handlers {
	return &Handlers{
		UserHandler:     NewUserHandler(s.UserService),
		AssetHandler:    NewAssetHandler(s.AssetService),
		LocationHandler: NewLocationHandler(s.LocationService),
		CategoryHandler: NewCategoryHandler(s.CategoryService),
		// DashboardHandler: NewDashboardHandler(s.DashboardService),
	}

}
