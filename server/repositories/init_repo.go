package repositories

import (
	"gorm.io/gorm"
)

type Repositories struct {
	UserRepository UserRepository
	// AssetRepository     AssetRepository
	// // LocationRepository  LocationRepository
	// // CategoryRepository  CategoryRepository
	// // DashboardRepository DashboardRepository
}

func InitRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		UserRepository: NewUserRepository(db),
		// AssetRepository:     NewAssetRepository(db),
		// LocationRepository:  NewLocationRepository(db),
		// CategoryRepository:  NewCategoryRepository(db),
		// DashboardRepository: NewDashboardRepository(db),
	}
}
