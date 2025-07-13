package services

import (
	"github.com/fiqrioemry/asset_management_system_app/server/repositories"
)

type Services struct {
	UserService     UserService
	AssetService    AssetService
	LocationService LocationService
	CategoryService CategoryService
	// DashboardService DashboardService
}

func InitServices(r *repositories.Repositories) *Services {
	return &Services{
		UserService:     NewUserService(r.UserRepository),
		AssetService:    NewAssetService(r.AssetRepository, r.LocationRepository, r.CategoryRepository),
		LocationService: NewLocationService(r.LocationRepository),
		CategoryService: NewCategoryService(r.CategoryRepository),
		// DashboardService: NewDashboardService(r.DashboardRepository),
	}
}
