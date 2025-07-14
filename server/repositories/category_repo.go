package repositories

import (
	"errors"

	"github.com/fiqrioemry/asset_management_system_app/server/models"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	Create(data *models.Category) error
	Update(data *models.Category) error
	Delete(data *models.Category) error
	GetByID(id string) (*models.Category, error)
	GetByIDAndUserID(id, userID string) (*models.Category, error)
	GetAllUserCategories(userID string) ([]models.Category, error)
	GetParentCategories(userID string) ([]models.Category, error)
	GetChildCategories(parentID, userID string) ([]models.Category, error)
	GetCategoryTree(userID string) ([]models.Category, error)
	CheckNameExists(name, userID string, parentID *string) (bool, error)
	GetAssetsByCategory(categoryID, userID string) ([]models.Asset, error)
	CountAssetsByCategory(categoryID, userID string) (int64, error)
	CountChildCategories(parentID string) (int64, error)
	ValidateParentAccess(parentID, userID string) (bool, error)
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db}
}

func (r *categoryRepository) Create(data *models.Category) error {
	return r.db.Create(data).Error
}

func (r *categoryRepository) Update(data *models.Category) error {
	return r.db.Save(data).Error
}

func (r *categoryRepository) Delete(data *models.Category) error {
	return r.db.Delete(data).Error
}

func (r *categoryRepository) GetByID(id string) (*models.Category, error) {
	var category models.Category
	err := r.db.Preload("Parent").Preload("Children").Where("id = ?", id).First(&category).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &category, err
}

func (r *categoryRepository) GetByIDAndUserID(id, userID string) (*models.Category, error) {
	var category models.Category
	err := r.db.Preload("Parent").Where("id = ? AND user_id = ?", id, userID).First(&category).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &category, err
}

func (r *categoryRepository) GetAllUserCategories(userID string) ([]models.Category, error) {
	var categories []models.Category
	err := r.db.Preload("Parent").
		Where("user_id IS NULL OR user_id = ?", userID).
		Order("parent_id IS NULL DESC, name ASC").
		Find(&categories).Error
	return categories, err
}

func (r *categoryRepository) GetParentCategories(userID string) ([]models.Category, error) {
	var categories []models.Category
	err := r.db.Where("parent_id IS NULL AND (user_id IS NULL OR user_id = ?)", userID).
		Order("is_default DESC, name ASC").
		Find(&categories).Error
	return categories, err
}

func (r *categoryRepository) GetChildCategories(parentID, userID string) ([]models.Category, error) {
	var categories []models.Category
	err := r.db.Where("parent_id = ? AND (user_id IS NULL OR user_id = ?)", parentID, userID).
		Order("is_default DESC, name ASC").
		Find(&categories).Error
	return categories, err
}

func (r *categoryRepository) GetCategoryTree(userID string) ([]models.Category, error) {
	var categories []models.Category
	err := r.db.Preload("Children", "user_id IS NULL OR user_id = ?", userID).
		Where("parent_id IS NULL AND (user_id IS NULL OR user_id = ?)", userID).
		Order("is_default DESC, name ASC").
		Find(&categories).Error
	return categories, err
}

func (r *categoryRepository) CheckNameExists(name, userID string, parentID *string) (bool, error) {
	var count int64
	query := r.db.Model(&models.Category{}).
		Where("LOWER(name) = LOWER(?) AND (user_id IS NULL OR user_id = ?)", name, userID)

	if parentID != nil {
		query = query.Where("parent_id = ?", *parentID)
	} else {
		query = query.Where("parent_id IS NULL")
	}

	err := query.Count(&count).Error
	return count > 0, err
}

func (r *categoryRepository) GetAssetsByCategory(categoryID, userID string) ([]models.Asset, error) {
	var assets []models.Asset
	err := r.db.Where("category_id = ? AND user_id = ?", categoryID, userID).
		Order("created_at DESC").
		Find(&assets).Error
	return assets, err
}

func (r *categoryRepository) CountAssetsByCategory(categoryID, userID string) (int64, error) {
	var count int64
	err := r.db.Model(&models.Asset{}).
		Where("category_id= ? AND user_id= ?", categoryID, userID).
		Count(&count).Error
	return count, err
}

func (r *categoryRepository) CountChildCategories(parentID string) (int64, error) {
	var count int64
	err := r.db.Model(&models.Category{}).
		Where("parent_id = ?", parentID).
		Count(&count).Error
	return count, err
}

func (r *categoryRepository) ValidateParentAccess(parentID, userID string) (bool, error) {
	var count int64
	err := r.db.Model(&models.Category{}).
		Where("id = ? AND (user_id IS NULL OR user_id = ?)", parentID, userID).
		Count(&count).Error
	return count > 0, err
}
