package repositories

import (
	"fmt"
	"strings"

	"github.com/fiqrioemry/asset_management_system_app/server/models"

	"gorm.io/gorm"
)

type AssetRepository interface {
	Create(asset *models.Asset) error
	Update(asset *models.Asset) error
	Delete(asset *models.Asset) error
	GetByID(id string) (*models.Asset, error)
	GetByIDAndUserID(id, userID string) (*models.Asset, error)
	GetAssetsWithFilter(filter AssetFilter) ([]models.Asset, int, error)
}

type AssetFilter struct {
	UserID     string
	Search     string
	CategoryID string
	LocationID string
	Condition  string
	MinPrice   *float64
	MaxPrice   *float64
	SortBy     string
	SortOrder  string
	Page       int
	Limit      int
}

type assetRepository struct {
	db *gorm.DB
}

func NewAssetRepository(db *gorm.DB) AssetRepository {
	return &assetRepository{db}
}

func (r *assetRepository) Create(asset *models.Asset) error {
	return r.db.Create(asset).Error
}

func (r *assetRepository) Update(asset *models.Asset) error {
	return r.db.Save(asset).Error
}

func (r *assetRepository) Delete(asset *models.Asset) error {
	return r.db.Delete(asset).Error
}

func (r *assetRepository) GetByID(id string) (*models.Asset, error) {
	var asset models.Asset
	err := r.db.Preload("Location").Preload("Category").Preload("User").
		Where("id = ?", id).First(&asset).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &asset, nil
}

func (r *assetRepository) GetByIDAndUserID(id, userID string) (*models.Asset, error) {
	var asset models.Asset
	err := r.db.Preload("Location").Preload("Category").Preload("User").
		Where("id = ? AND user_id = ?", id, userID).First(&asset).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &asset, nil
}

func (r *assetRepository) GetAssetsWithFilter(filter AssetFilter) ([]models.Asset, int, error) {
	var assets []models.Asset
	var totalCount int64

	// Build query
	query := r.db.Model(&models.Asset{}).Where("user_id = ?", filter.UserID)

	// Apply filters
	if filter.Search != "" {
		searchTerm := "%" + strings.ToLower(filter.Search) + "%"
		query = query.Where("LOWER(name) LIKE ? OR LOWER(description) LIKE ? OR LOWER(serial_number) LIKE ?",
			searchTerm, searchTerm, searchTerm)
	}

	if filter.CategoryID != "" {
		query = query.Where("category_id = ?", filter.CategoryID)
	}

	if filter.LocationID != "" {
		query = query.Where("location_id = ?", filter.LocationID)
	}

	if filter.Condition != "" {
		query = query.Where("condition = ?", filter.Condition)
	}

	if filter.MinPrice != nil {
		query = query.Where("price >= ?", *filter.MinPrice)
	}

	if filter.MaxPrice != nil {
		query = query.Where("price <= ?", *filter.MaxPrice)
	}

	// Count total records
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	// Apply sorting
	orderBy := r.buildOrderBy(filter.SortBy, filter.SortOrder)
	if orderBy != "" {
		query = query.Order(orderBy)
	}

	// Apply pagination
	offset := (filter.Page - 1) * filter.Limit
	query = query.Offset(offset).Limit(filter.Limit)

	// Execute query with preloading
	err := query.Preload("Location").Preload("Category").Find(&assets).Error
	if err != nil {
		return nil, 0, err
	}

	return assets, int(totalCount), nil
}

func (r *assetRepository) buildOrderBy(sortBy, sortOrder string) string {
	validSortBy := map[string]string{
		"name":          "name",
		"price":         "price",
		"created_at":    "created_at",
		"purchase_date": "purchase_date",
	}

	validSortOrder := map[string]string{
		"asc":  "ASC",
		"desc": "DESC",
	}

	column, validColumn := validSortBy[sortBy]
	order, validOrder := validSortOrder[sortOrder]

	if !validColumn || !validOrder {
		return "created_at DESC" // default
	}

	return fmt.Sprintf("%s %s", column, order)
}
