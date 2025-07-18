package services

import (
	"strings"

	"github.com/fiqrioemry/asset_management_system_app/server/dto"
	"github.com/fiqrioemry/asset_management_system_app/server/models"
	"github.com/fiqrioemry/asset_management_system_app/server/repositories"
	"github.com/fiqrioemry/asset_management_system_app/server/utils"
	"github.com/fiqrioemry/go-api-toolkit/response"

	"github.com/google/uuid"
)

type AssetService interface {
	DeleteAsset(userID, assetID string) error
	GetAssetByID(userID, assetID string) (*dto.AssetResponse, error)
	CreateAsset(userID string, req *dto.CreateAssetRequest) (*dto.AssetResponse, error)
	UpdateAsset(userID, assetID string, req *dto.UpdateAssetRequest) (*dto.AssetResponse, error)
	GetAssets(userID string, req *dto.GetAssetsRequest) (*[]dto.AssetResponse, int, error)
}

type assetService struct {
	assetRepo    repositories.AssetRepository
	locationRepo repositories.LocationRepository
	categoryRepo repositories.CategoryRepository
}

func NewAssetService(
	assetRepo repositories.AssetRepository,
	locationRepo repositories.LocationRepository,
	categoryRepo repositories.CategoryRepository,
) AssetService {
	return &assetService{
		assetRepo:    assetRepo,
		locationRepo: locationRepo,
		categoryRepo: categoryRepo,
	}
}

func (s *assetService) CreateAsset(userID string, req *dto.CreateAssetRequest) (*dto.AssetResponse, error) {

	// Validate location access
	location, err := s.locationRepo.GetByIDAndUserID(req.LocationID, userID)
	if err != nil {
		return nil, response.NewInternalServerError("Failed to validate location", err)
	}
	if location == nil {
		return nil, response.NewNotFound("Location not found or access denied")
	}

	// Validate category access
	category, err := s.categoryRepo.GetByIDAndUserID(req.CategoryID, userID)
	if err != nil {
		return nil, response.NewInternalServerError("Failed to validate category", err)
	}
	if category == nil {
		return nil, response.NewNotFound("Category not found or access denied")
	}

	// Parse all string IDs to UUIDs
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, response.NewBadRequest("Invalid user ID")
	}

	locationUUID, err := uuid.Parse(req.LocationID)
	if err != nil {
		return nil, response.NewBadRequest("Invalid location ID")
	}

	categoryUUID, err := uuid.Parse(req.CategoryID)
	if err != nil {
		return nil, response.NewBadRequest("Invalid category ID")
	}

	// Create asset
	asset := &models.Asset{
		Name:         strings.TrimSpace(req.Name),
		Description:  strings.TrimSpace(req.Description),
		LocationID:   locationUUID,
		CategoryID:   categoryUUID,
		UserID:       userUUID,
		Image:        req.ImageURL,
		PurchaseDate: req.PurchaseDate,
		Price:        req.Price,
		Condition:    req.Condition,
		SerialNumber: strings.TrimSpace(req.SerialNumber),
		Warranty:     req.Warranty,
	}

	if err := s.assetRepo.Create(asset); err != nil {
		return nil, response.NewInternalServerError("Failed to create asset", err)
	}

	// Load relationships for response
	asset.Location = *location

	asset.Category = *category

	response := s.convertToResponse(asset)
	return &response, nil
}

func (s *assetService) GetAssets(userID string, req *dto.GetAssetsRequest) (*[]dto.AssetResponse, int, error) {

	// validate price range
	if req.MinPrice != nil && req.MaxPrice != nil && *req.MinPrice > *req.MaxPrice {
		return nil, 0, response.NewBadRequest("Min price cannot be greater than max price")
	}

	filter := repositories.AssetFilter{
		UserID:     userID,
		Search:     strings.TrimSpace(req.Search),
		CategoryID: req.CategoryID,
		LocationID: req.LocationID,
		Condition:  req.Condition,
		MinPrice:   req.MinPrice,
		MaxPrice:   req.MaxPrice,
		SortBy:     req.SortBy,
		SortOrder:  req.SortOrder,
		Page:       req.Page,
		Limit:      req.Limit,
	}

	// get assets and total count
	assets, total, err := s.assetRepo.GetAssetsWithFilter(filter)
	if err != nil {
		return nil, 0, response.NewInternalServerError("Failed to get assets", err)
	}

	// convert assert to response
	var assetResponses []dto.AssetResponse
	for _, asset := range assets {
		assetResponses = append(assetResponses, s.convertToResponse(&asset))
	}

	return &assetResponses, int(total), nil
}

func (s *assetService) GetAssetByID(userID, assetID string) (*dto.AssetResponse, error) {
	asset, err := s.assetRepo.GetByIDAndUserID(assetID, userID)
	if err != nil {
		return nil, response.NewInternalServerError("Failed to get asset", err)
	}
	if asset == nil {
		return nil, response.NewNotFound("Asset not found")
	}

	response := s.convertToResponse(asset)
	return &response, nil
}

func (s *assetService) UpdateAsset(userID, assetID string, req *dto.UpdateAssetRequest) (*dto.AssetResponse, error) {
	// Get asset and check ownership
	asset, err := s.assetRepo.GetByIDAndUserID(assetID, userID)
	if err != nil {
		return nil, response.NewInternalServerError("Failed to get asset", err)
	}
	if asset == nil {
		return nil, response.NewNotFound("Asset not found or you don't have permission to update it")
	}

	// Validate location if provided
	if req.LocationID != "" {
		location, err := s.locationRepo.GetByID(req.LocationID)
		if err != nil {
			return nil, response.NewInternalServerError("Failed to validate location", err)
		}
		if location == nil {
			return nil, response.NewNotFound("Location not found or access denied")
		}

		locationUUID, err := uuid.Parse(req.LocationID)
		if err != nil {
			return nil, response.NewBadRequest("Invalid location ID")
		}
		asset.LocationID = locationUUID
	}

	// Validate category if provided
	if req.CategoryID != "" {
		category, err := s.categoryRepo.GetByID(req.CategoryID)
		if err != nil {
			return nil, response.NewInternalServerError("Failed to validate category", err)
		}
		if category == nil {
			return nil, response.NewNotFound("Category not found or access denied")
		}

		categoryUUID, err := uuid.Parse(req.CategoryID)
		if err != nil {
			return nil, response.NewBadRequest("Invalid category ID")
		}
		asset.CategoryID = categoryUUID
	}

	// Update fields
	if req.Name != "" {
		asset.Name = strings.TrimSpace(req.Name)
	}
	if req.Description != "" {
		asset.Description = strings.TrimSpace(req.Description)
	}
	if req.ImageURL != "" {
		asset.Image = req.ImageURL
	}
	if req.PurchaseDate != nil {
		asset.PurchaseDate = req.PurchaseDate
	}
	if req.Price != nil {
		asset.Price = *req.Price
	}
	if req.Condition != "" {
		asset.Condition = req.Condition
	}
	if req.SerialNumber != "" {
		asset.SerialNumber = strings.TrimSpace(req.SerialNumber)
	}
	if req.Warranty != nil {
		asset.Warranty = req.Warranty
	}

	if err := s.assetRepo.Update(asset); err != nil {
		return nil, response.NewInternalServerError("Failed to update asset", err)
	}

	response := s.convertToResponse(asset)
	return &response, nil
}

func (s *assetService) DeleteAsset(userID, assetID string) error {
	// Get asset and check ownership
	asset, err := s.assetRepo.GetByIDAndUserID(assetID, userID)
	if err != nil {
		return response.NewInternalServerError("Failed to get asset", err)
	}
	if asset == nil {
		return response.NewNotFound("Asset not found or you don't have permission to delete it")
	}

	if err := s.assetRepo.Delete(asset); err != nil {
		return response.NewInternalServerError("Failed to delete asset", err)
	}

	// cleanup image if exists
	if asset.Image != "" {
		go utils.CleanupImageOnError(asset.Image)
	}

	return nil
}

func (s *assetService) convertToResponse(asset *models.Asset) dto.AssetResponse {
	response := dto.AssetResponse{
		ID:           asset.ID.String(),
		Name:         asset.Name,
		Description:  asset.Description,
		LocationID:   asset.LocationID.String(),
		CategoryID:   asset.CategoryID.String(),
		UserID:       asset.UserID.String(),
		Image:        asset.Image,
		PurchaseDate: asset.PurchaseDate,
		Price:        asset.Price,
		Condition:    asset.Condition,
		SerialNumber: asset.SerialNumber,
		Warranty:     asset.Warranty,
		CreatedAt:    asset.CreatedAt,
		UpdatedAt:    asset.UpdatedAt,
	}

	// Add location if preloaded
	if asset.Location.ID != uuid.Nil {
		response.Location = &dto.LocationResponse{
			ID:   asset.Location.ID.String(),
			Name: asset.Location.Name,
		}
	}

	// Add category if preloaded
	if asset.Category.ID != uuid.Nil {
		response.Category = &dto.CategoryResponse{
			ID:   asset.Category.ID.String(),
			Name: asset.Category.Name,
		}
	}

	return response
}
