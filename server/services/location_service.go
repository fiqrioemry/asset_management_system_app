// services/location_service.go
package services

import (
	"fmt"
	"strings"
	"time"

	"github.com/fiqrioemry/asset_management_system_app/server/dto"
	"github.com/fiqrioemry/asset_management_system_app/server/errors"
	"github.com/fiqrioemry/asset_management_system_app/server/models"
	"github.com/fiqrioemry/asset_management_system_app/server/repositories"
	"github.com/fiqrioemry/asset_management_system_app/server/utils"
	"github.com/google/uuid"
)

type LocationService interface {
	GetLocations(userID string) (*dto.LocationsResponse, error)
	CreateLocation(userID string, req *dto.CreateLocationRequest) (*dto.LocationResponse, error)
	UpdateLocation(userID, locationID string, req *dto.UpdateLocationRequest) (*dto.LocationResponse, error)
	DeleteLocation(userID, locationID string) error
	GetLocationByID(userID, locationID string) (*dto.LocationResponse, error)
	GetAssetsByLocation(userID, locationID string) (*dto.LocationWithAssetsResponse, error)
}

type locationService struct {
	locationRepo repositories.LocationRepository
}

func NewLocationService(locationRepo repositories.LocationRepository) LocationService {
	return &locationService{
		locationRepo: locationRepo,
	}
}

// GetLocations returns all locations available to user (system + user-specific) with caching
func (s *locationService) GetLocations(userID string) (*dto.LocationsResponse, error) {
	cacheKey := fmt.Sprintf("asset_app:cache:locations:all:%s", userID)

	// Try cache first
	var cachedResponse dto.LocationsResponse
	if err := utils.GetKey(cacheKey, &cachedResponse); err == nil {
		return &cachedResponse, nil
	}

	// Cache miss - get from database
	locations, err := s.locationRepo.GetAllUserLocations(userID)
	if err != nil {
		return nil, errors.NewInternalServerError("Failed to get locations", err)
	}

	// Convert to response format
	var locationResponses []dto.LocationResponse
	for _, location := range locations {
		locationResponses = append(locationResponses, dto.LocationResponse{
			ID:        location.ID.String(),
			Name:      location.Name,
			IsDefault: location.IsDefault,
			IsCustom:  location.UserID != nil,
			CreatedAt: location.CreatedAt,
			UpdatedAt: location.UpdatedAt,
		})
	}

	response := &dto.LocationsResponse{
		Locations: locationResponses,
		Total:     len(locationResponses),
	}

	// Cache for 15 minutes
	go utils.AddKeys(cacheKey, response, 15*time.Minute)

	return response, nil
}

// CreateLocation creates a new user location
func (s *locationService) CreateLocation(userID string, req *dto.CreateLocationRequest) (*dto.LocationResponse, error) {
	// Normalize name
	req.Name = strings.TrimSpace(req.Name)
	req.Name = strings.Title(strings.ToLower(req.Name))

	// Check if name already exists for this user
	exists, err := s.locationRepo.CheckNameExists(req.Name, userID)
	if err != nil {
		return nil, errors.NewInternalServerError("Failed to check location name", err)
	}

	if exists {
		return nil, errors.NewConflict("Location name already exists")
	}

	// Parse userID to UUID
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.NewBadRequest("Invalid user ID")
	}

	// Create location
	location := &models.Location{
		Name:      req.Name,
		UserID:    &userUUID,
		IsDefault: false,
	}

	if err := s.locationRepo.Create(location); err != nil {
		return nil, errors.NewInternalServerError("Failed to create location", err)
	}

	// Invalidate cache
	s.invalidateUserCache(userID)

	response := &dto.LocationResponse{
		ID:        location.ID.String(),
		Name:      location.Name,
		IsDefault: location.IsDefault,
		IsCustom:  true,
		CreatedAt: location.CreatedAt,
		UpdatedAt: location.UpdatedAt,
	}

	return response, nil
}

// UpdateLocation updates user's own location only
func (s *locationService) UpdateLocation(userID, locationID string, req *dto.UpdateLocationRequest) (*dto.LocationResponse, error) {
	// Get location and check ownership (only user's own locations can be updated)
	location, err := s.locationRepo.GetByIDAndUserID(locationID, userID)
	if err != nil {
		return nil, errors.NewInternalServerError("Failed to get location", err)
	}

	if location == nil {
		return nil, errors.NewNotFound("Location not found or you don't have permission to update it")
	}

	// Normalize name
	req.Name = strings.TrimSpace(req.Name)
	req.Name = strings.Title(strings.ToLower(req.Name))

	// Check if new name already exists (excluding current location)
	if strings.ToLower(req.Name) != strings.ToLower(location.Name) {
		exists, err := s.locationRepo.CheckNameExists(req.Name, userID)
		if err != nil {
			return nil, errors.NewInternalServerError("Failed to check location name", err)
		}

		if exists {
			return nil, errors.NewConflict("Location name already exists")
		}
	}

	// Update location
	location.Name = req.Name

	if err := s.locationRepo.Update(location); err != nil {
		return nil, errors.NewInternalServerError("Failed to update location", err)
	}

	// Invalidate cache
	s.invalidateUserCache(userID)

	response := &dto.LocationResponse{
		ID:        location.ID.String(),
		Name:      location.Name,
		IsDefault: location.IsDefault,
		IsCustom:  true,
		CreatedAt: location.CreatedAt,
		UpdatedAt: location.UpdatedAt,
	}

	return response, nil
}

// DeleteLocation deletes user's own location only (not system default locations)
func (s *locationService) DeleteLocation(userID, locationID string) error {
	// Get location and check ownership (only user's own locations)
	location, err := s.locationRepo.GetByIDAndUserID(locationID, userID)
	if err != nil {
		return errors.NewInternalServerError("Failed to get location", err)
	}

	if location == nil {
		return errors.NewNotFound("Location not found or you don't have permission to delete it")
	}

	// Additional check: cannot delete system default locations
	if location.IsDefault {
		return errors.NewForbidden("Cannot delete system default location")
	}

	// Check if location is being used by assets
	assetCount, err := s.locationRepo.CountAssetsByLocation(locationID, userID)
	if err != nil {
		return errors.NewInternalServerError("Failed to check location usage", err)
	}

	if assetCount > 0 {
		return errors.NewConflict("Cannot delete location that is being used by assets")
	}

	// Delete location
	if err := s.locationRepo.Delete(location); err != nil {
		return errors.NewInternalServerError("Failed to delete location", err)
	}

	// Invalidate cache
	s.invalidateUserCache(userID)

	return nil
}

// GetLocationByID gets specific location that user can access
func (s *locationService) GetLocationByID(userID, locationID string) (*dto.LocationResponse, error) {
	// Get location
	location, err := s.locationRepo.GetByID(locationID)
	if err != nil {
		return nil, errors.NewInternalServerError("Failed to get location", err)
	}

	if location == nil {
		return nil, errors.NewNotFound("Location not found")
	}

	// Check if user can access this location (system location or user's own location)
	if location.UserID != nil && location.UserID.String() != userID {
		return nil, errors.NewNotFound("Location not found")
	}

	response := &dto.LocationResponse{
		ID:        location.ID.String(),
		Name:      location.Name,
		IsDefault: location.IsDefault,
		IsCustom:  location.UserID != nil,
		CreatedAt: location.CreatedAt,
		UpdatedAt: location.UpdatedAt,
	}

	return response, nil
}

// GetAssetsByLocation gets all assets in a specific location for the user
func (s *locationService) GetAssetsByLocation(userID, locationID string) (*dto.LocationWithAssetsResponse, error) {
	// Get location first
	locationResponse, err := s.GetLocationByID(userID, locationID)
	if err != nil {
		return nil, err
	}

	// Get assets for this location
	assets, err := s.locationRepo.GetAssetsByLocation(locationID, userID)
	if err != nil {
		return nil, errors.NewInternalServerError("Failed to get assets", err)
	}

	// Convert assets to response format
	var assetResponses []dto.LocationAssetResponse
	for _, asset := range assets {
		assetResponses = append(assetResponses, dto.LocationAssetResponse{
			ID:           asset.ID.String(),
			Name:         asset.Name,
			Description:  asset.Description,
			Price:        asset.Price,
			Condition:    asset.Condition,
			SerialNumber: asset.SerialNumber,
			CreatedAt:    asset.CreatedAt,
		})
	}

	response := &dto.LocationWithAssetsResponse{
		Location: *locationResponse,
		Assets:   assetResponses,
		Total:    len(assetResponses),
	}

	return response, nil
}

// Helper function to invalidate user cache
func (s *locationService) invalidateUserCache(userID string) {
	cacheKey := fmt.Sprintf("asset_app:cache:locations:all:%s", userID)
	utils.DeleteKeys(cacheKey)
}
