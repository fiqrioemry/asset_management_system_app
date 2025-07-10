// repositories/location_repo.go
package repositories

import (
	"errors"

	"github.com/fiqrioemry/asset_management_system_app/server/models"
	"gorm.io/gorm"
)

type LocationRepository interface {
	Create(data *models.Location) error
	Update(data *models.Location) error
	Delete(data *models.Location) error
	GetByID(id string) (*models.Location, error)
	GetByIDAndUserID(id, userID string) (*models.Location, error)
	GetAllUserLocations(userID string) ([]models.Location, error)
	CheckNameExists(name, userID string) (bool, error)
	GetAssetsByLocation(locationID, userID string) ([]models.Asset, error)
	CountAssetsByLocation(locationID, userID string) (int64, error)
}

type locationRepository struct {
	db *gorm.DB
}

func NewLocationRepository(db *gorm.DB) LocationRepository {
	return &locationRepository{db}
}

func (r *locationRepository) Create(data *models.Location) error {
	return r.db.Create(data).Error
}

func (r *locationRepository) Update(data *models.Location) error {
	return r.db.Save(data).Error
}

func (r *locationRepository) Delete(data *models.Location) error {
	return r.db.Delete(data).Error
}

func (r *locationRepository) GetByID(id string) (*models.Location, error) {
	var location models.Location
	err := r.db.Where("id = ?", id).First(&location).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &location, err
}

func (r *locationRepository) GetByIDAndUserID(id, userID string) (*models.Location, error) {
	var location models.Location
	err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&location).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &location, err
}

func (r *locationRepository) GetAllUserLocations(userID string) ([]models.Location, error) {
	var locations []models.Location
	err := r.db.Where("user_id IS NULL OR user_id = ?", userID).
		Order("is_default DESC, name ASC").
		Find(&locations).Error
	return locations, err
}

func (r *locationRepository) CheckNameExists(name, userID string) (bool, error) {
	var count int64
	err := r.db.Model(&models.Location{}).
		Where("LOWER(name) = LOWER(?) AND (user_id IS NULL OR user_id = ?)", name, userID).
		Count(&count).Error
	return count > 0, err
}

func (r *locationRepository) GetAssetsByLocation(locationID, userID string) ([]models.Asset, error) {
	var assets []models.Asset
	err := r.db.Where("location_id = ? AND user_id = ?", locationID, userID).
		Order("created_at DESC").
		Find(&assets).Error
	return assets, err
}

func (r *locationRepository) CountAssetsByLocation(locationID, userID string) (int64, error) {
	var count int64
	err := r.db.Model(&models.Asset{}).
		Where("location_id = ? AND user_id = ?", locationID, userID).
		Count(&count).Error
	return count, err
}
