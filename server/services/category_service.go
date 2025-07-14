package services

import (
	"fmt"
	"strings"
	"time"

	"github.com/fiqrioemry/asset_management_system_app/server/dto"
	"github.com/fiqrioemry/asset_management_system_app/server/models"
	"github.com/fiqrioemry/asset_management_system_app/server/repositories"
	"github.com/fiqrioemry/asset_management_system_app/server/utils"
	"github.com/fiqrioemry/go-api-toolkit/response"
	"github.com/google/uuid"
)

type CategoryService interface {
	DeleteCategory(userID, categoryID string) error
	GetCategoriesTree(userID string) (*dto.CategoriesTreeResponse, error)
	GetCategoriesFlat(userID string) (*dto.CategoriesFlatResponse, error)
	GetParentCategories(userID string) (*dto.CategoriesTreeResponse, error)
	GetCategoryByID(userID, categoryID string) (*dto.CategoryResponse, error)
	GetChildCategories(parentID, userID string) (*dto.CategoriesTreeResponse, error)
	GetAssetsByCategory(userID, categoryID string) (*dto.CategoryWithAssetsResponse, error)
	CreateCategory(userID string, req *dto.CreateCategoryRequest) (*dto.CategoryResponse, error)
	UpdateCategory(userID, categoryID string, req *dto.UpdateCategoryRequest) (*dto.CategoryResponse, error)
}

type categoryService struct {
	categoryRepo repositories.CategoryRepository
}

func NewCategoryService(categoryRepo repositories.CategoryRepository) CategoryService {
	return &categoryService{
		categoryRepo: categoryRepo,
	}
}

func (s *categoryService) GetCategoriesTree(userID string) (*dto.CategoriesTreeResponse, error) {
	cacheKey := fmt.Sprintf("asset_app:cache:categories:tree:%s", userID)

	// Try cache first
	var cachedResponse dto.CategoriesTreeResponse
	if err := utils.GetKey(cacheKey, &cachedResponse); err == nil {
		return &cachedResponse, nil
	}

	// Get category tree from database
	categories, err := s.categoryRepo.GetCategoryTree(userID)
	if err != nil {
		return nil, response.NewInternalServerError("Failed to get categories", err)
	}

	// Convert to response format
	var categoryResponses []dto.CategoryResponse
	parentCount := 0
	childCount := 0

	for _, category := range categories {
		parentCount++
		categoryResponse := s.convertToResponse(&category, 0)

		// Convert children
		for _, child := range category.Children {
			childCount++
			childResponse := s.convertToResponse(&child, 1)
			categoryResponse.Children = append(categoryResponse.Children, childResponse)
		}

		categoryResponses = append(categoryResponses, categoryResponse)
	}

	response := &dto.CategoriesTreeResponse{
		Categories: categoryResponses,
		Total:      parentCount + childCount,
		Parents:    parentCount,
		Children:   childCount,
	}

	// Cache for 15 minutes
	go utils.AddKeys(cacheKey, response, 15*time.Minute)

	return response, nil
}

func (s *categoryService) GetCategoriesFlat(userID string) (*dto.CategoriesFlatResponse, error) {
	cacheKey := fmt.Sprintf("asset_app:cache:categories:flat:%s", userID)

	// Try cache first
	var cachedResponse dto.CategoriesFlatResponse
	if err := utils.GetKey(cacheKey, &cachedResponse); err == nil {
		return &cachedResponse, nil
	}

	// Get all categories
	categories, err := s.categoryRepo.GetAllUserCategories(userID)
	if err != nil {
		return nil, response.NewInternalServerError("Failed to get categories", err)
	}

	// Convert to flat response
	var categoryResponses []dto.CategoryFlatResponse
	for _, category := range categories {
		flatResponse := dto.CategoryFlatResponse{
			ID:        category.ID.String(),
			Name:      category.Name,
			FullName:  s.getFullCategoryName(&category),
			IsDefault: category.IsDefault,
			IsCustom:  category.UserID != nil,
			IsParent:  category.ParentID == nil,
			Level:     s.getCategoryLevel(&category),
			CreatedAt: category.CreatedAt,
			UpdatedAt: category.UpdatedAt,
		}

		if category.ParentID != nil {
			flatResponse.ParentID = &[]string{category.ParentID.String()}[0]
		}

		categoryResponses = append(categoryResponses, flatResponse)
	}

	response := &dto.CategoriesFlatResponse{
		Categories: categoryResponses,
		Total:      len(categoryResponses),
	}

	// Cache for 15 minutes
	go utils.AddKeys(cacheKey, response, 15*time.Minute)

	return response, nil
}

// GetParentCategories returns only parent categories
func (s *categoryService) GetParentCategories(userID string) (*dto.CategoriesTreeResponse, error) {
	categories, err := s.categoryRepo.GetParentCategories(userID)
	if err != nil {
		return nil, response.NewInternalServerError("Failed to get parent categories", err)
	}

	var categoryResponses []dto.CategoryResponse
	for _, category := range categories {
		categoryResponses = append(categoryResponses, s.convertToResponse(&category, 0))
	}

	return &dto.CategoriesTreeResponse{
		Categories: categoryResponses,
		Total:      len(categoryResponses),
		Parents:    len(categoryResponses),
		Children:   0,
	}, nil
}

// GetChildCategories returns children of specific parent
func (s *categoryService) GetChildCategories(parentID, userID string) (*dto.CategoriesTreeResponse, error) {
	categories, err := s.categoryRepo.GetChildCategories(parentID, userID)
	if err != nil {
		return nil, response.NewInternalServerError("Failed to get child categories", err)
	}

	var categoryResponses []dto.CategoryResponse
	for _, category := range categories {
		categoryResponses = append(categoryResponses, s.convertToResponse(&category, 1))
	}

	return &dto.CategoriesTreeResponse{
		Categories: categoryResponses,
		Total:      len(categoryResponses),
		Parents:    0,
		Children:   len(categoryResponses),
	}, nil
}

// CreateCategory creates new category (parent or child)
func (s *categoryService) CreateCategory(userID string, req *dto.CreateCategoryRequest) (*dto.CategoryResponse, error) {
	// Normalize name
	req.Name = strings.TrimSpace(req.Name)

	// Validate parent if provided
	var parentUUID *uuid.UUID
	if req.ParentID != nil && *req.ParentID != "" {
		// Check if parent exists and user has access
		hasAccess, err := s.categoryRepo.ValidateParentAccess(*req.ParentID, userID)
		if err != nil {
			return nil, response.NewInternalServerError("Failed to validate parent access", err)
		}
		if !hasAccess {
			return nil, response.NewNotFound("Parent category not found or access denied")
		}

		parentID, err := uuid.Parse(*req.ParentID)
		if err != nil {
			return nil, response.NewBadRequest("Invalid parent ID")
		}
		parentUUID = &parentID
	}

	// Check if name already exists in same scope (same parent)
	exists, err := s.categoryRepo.CheckNameExists(req.Name, userID, req.ParentID)
	if err != nil {
		return nil, response.NewInternalServerError("Failed to check category name", err)
	}
	if exists {
		return nil, response.NewConflict("Category name already exists in this scope")
	}

	// Parse userID
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, response.NewBadRequest("Invalid user ID")
	}

	// Create category
	category := &models.Category{
		ParentID:  parentUUID,
		Name:      req.Name,
		UserID:    &userUUID,
		IsDefault: false,
	}

	if err := s.categoryRepo.Create(category); err != nil {
		return nil, response.NewInternalServerError("Failed to create category", err)
	}

	// Invalidate cache
	go s.invalidateUserCache(userID)

	level := 0
	if category.ParentID != nil {
		level = 1
	}

	response := s.convertToResponse(category, level)
	return &response, nil
}

func (s *categoryService) UpdateCategory(userID, categoryID string, req *dto.UpdateCategoryRequest) (*dto.CategoryResponse, error) {
	// Get category and check ownership
	category, err := s.categoryRepo.GetByIDAndUserID(categoryID, userID)
	if err != nil {
		return nil, response.NewInternalServerError("Failed to get category", err)
	}
	if category == nil {
		return nil, response.NewNotFound("Category not found or you don't have permission to update it")
	}

	// Normalize name
	req.Name = strings.TrimSpace(req.Name)

	// Validate parent if changed
	var parentUUID *uuid.UUID
	if req.ParentID != nil && *req.ParentID != "" {
		// Check if parent exists and user has access
		hasAccess, err := s.categoryRepo.ValidateParentAccess(*req.ParentID, userID)
		if err != nil {
			return nil, response.NewInternalServerError("Failed to validate parent access", err)
		}
		if !hasAccess {
			return nil, response.NewNotFound("Parent category not found or access denied")
		}

		parentID, err := uuid.Parse(*req.ParentID)
		if err != nil {
			return nil, response.NewBadRequest("Invalid parent ID")
		}
		parentUUID = &parentID

		// Cannot set self as parent
		if *req.ParentID == categoryID {
			return nil, response.NewBadRequest("Category cannot be its own parent")
		}
	}

	// Check name uniqueness in new scope if name or parent changed
	nameChanged := !strings.EqualFold(req.Name, category.Name)
	parentChanged := (req.ParentID == nil && category.ParentID != nil) ||
		(req.ParentID != nil && category.ParentID == nil) ||
		(req.ParentID != nil && category.ParentID != nil && *req.ParentID != category.ParentID.String())

	if nameChanged || parentChanged {
		exists, err := s.categoryRepo.CheckNameExists(req.Name, userID, req.ParentID)
		if err != nil {
			return nil, response.NewInternalServerError("Failed to check category name", err)
		}
		if exists {
			return nil, response.NewConflict("Category name already exists in this scope")
		}
	}

	// Update category
	category.Name = req.Name
	category.ParentID = parentUUID

	if err := s.categoryRepo.Update(category); err != nil {
		return nil, response.NewInternalServerError("Failed to update category", err)
	}

	// Invalidate cache
	go s.invalidateUserCache(userID)

	level := 0
	if category.ParentID != nil {
		level = 1
	}

	response := s.convertToResponse(category, level)
	return &response, nil
}

func (s *categoryService) DeleteCategory(userID, categoryID string) error {
	// Get category and check ownership
	category, err := s.categoryRepo.GetByIDAndUserID(categoryID, userID)
	if err != nil {
		return response.NewInternalServerError("Failed to get category", err)
	}
	if category == nil {
		return response.NewNotFound("Category not found or you don't have permission to delete it")
	}

	// Cannot delete system default categories
	if category.IsDefault {
		return response.NewForbidden("Cannot delete system default category")
	}

	// Check if category has children
	childCount, err := s.categoryRepo.CountChildCategories(categoryID)
	if err != nil {
		return response.NewInternalServerError("Failed to check child categories", err)
	}
	if childCount > 0 {
		return response.NewConflict("Cannot delete category that has subcategories")
	}

	// Check if category is being used by assets
	assetCount, err := s.categoryRepo.CountAssetsByCategory(categoryID, userID)
	if err != nil {
		return response.NewInternalServerError("Failed to check category usage", err)
	}
	if assetCount > 0 {
		return response.NewConflict("Cannot delete category that is being used by assets")
	}

	// Delete category
	if err := s.categoryRepo.Delete(category); err != nil {
		return response.NewInternalServerError("Failed to delete category", err)
	}

	// Invalidate cache
	go s.invalidateUserCache(userID)

	return nil
}

func (s *categoryService) GetCategoryByID(userID, categoryID string) (*dto.CategoryResponse, error) {
	category, err := s.categoryRepo.GetByID(categoryID)
	if err != nil {
		return nil, response.NewInternalServerError("Failed to get category", err)
	}
	if category == nil {
		return nil, response.NewNotFound("Category not found")
	}

	// Check if user can access this category
	if category.UserID != nil && category.UserID.String() != userID {
		return nil, response.NewNotFound("Category not found")
	}

	level := 0
	if category.ParentID != nil {
		level = 1
	}

	// Convert to response format
	resp := s.convertToResponse(category, level)
	return &resp, nil
}

func (s *categoryService) GetAssetsByCategory(userID, categoryID string) (*dto.CategoryWithAssetsResponse, error) {
	// Get category first
	categoryResponse, err := s.GetCategoryByID(userID, categoryID)
	if err != nil {
		return nil, err
	}

	// Get assets for this category
	assets, err := s.categoryRepo.GetAssetsByCategory(categoryID, userID)
	if err != nil {
		return nil, response.NewInternalServerError("Failed to get assets", err)
	}

	// Convert assets to response format
	var assetResponses []dto.CategoryAssetResponse
	for _, asset := range assets {
		assetResponses = append(assetResponses, dto.CategoryAssetResponse{
			ID:           asset.ID.String(),
			Name:         asset.Name,
			Description:  asset.Description,
			Price:        asset.Price,
			Condition:    asset.Condition,
			SerialNumber: asset.SerialNumber,
			CreatedAt:    asset.CreatedAt,
		})
	}

	resp := &dto.CategoryWithAssetsResponse{
		Category: *categoryResponse,
		Assets:   assetResponses,
		Total:    len(assetResponses),
	}

	return resp, nil
}

// Helper methods
func (s *categoryService) convertToResponse(category *models.Category, level int) dto.CategoryResponse {
	resp := dto.CategoryResponse{
		ID:        category.ID.String(),
		Name:      category.Name,
		IsDefault: category.IsDefault,
		IsCustom:  category.UserID != nil,
		IsParent:  category.ParentID == nil,
		Level:     level,
		CreatedAt: category.CreatedAt,
		UpdatedAt: category.UpdatedAt,
	}

	if category.ParentID != nil {
		resp.ParentID = &[]string{category.ParentID.String()}[0]
	}

	return resp
}

func (s *categoryService) getFullCategoryName(category *models.Category) string {
	if category.Parent != nil {
		return category.Parent.Name + " > " + category.Name
	}
	return category.Name
}

func (s *categoryService) getCategoryLevel(category *models.Category) int {
	if category.ParentID == nil {
		return 0
	}
	return 1
}

func (s *categoryService) invalidateUserCache(userID string) {
	cacheKeys := []string{
		fmt.Sprintf("asset_app:cache:categories:tree:%s", userID),
		fmt.Sprintf("asset_app:cache:categories:flat:%s", userID),
	}
	utils.DeleteKeys(cacheKeys...)
}
