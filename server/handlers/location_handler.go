// handlers/location_handler.go
package handlers

import (
	"net/http"

	"github.com/fiqrioemry/asset_management_system_app/server/dto"
	"github.com/fiqrioemry/asset_management_system_app/server/services"
	"github.com/fiqrioemry/asset_management_system_app/server/utils"
	"github.com/gin-gonic/gin"
)

type LocationHandler struct {
	service services.LocationService
}

func NewLocationHandler(service services.LocationService) *LocationHandler {
	return &LocationHandler{service}
}

// GetLocations returns all locations (system + user's custom locations)
func (h *LocationHandler) GetLocations(c *gin.Context) {
	userID := utils.MustGetUserID(c)

	locations, err := h.service.GetLocations(userID)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    locations,
		"message": "Locations retrieved successfully",
	})
}

// CreateLocation creates a new user location
func (h *LocationHandler) CreateLocation(c *gin.Context) {
	userID := utils.MustGetUserID(c)

	var req dto.CreateLocationRequest
	if !utils.BindAndValidateJSON(c, &req) {
		return
	}

	location, err := h.service.CreateLocation(userID, &req)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    location,
		"message": "Location created successfully",
	})
}

// UpdateLocation updates user's own location
func (h *LocationHandler) UpdateLocation(c *gin.Context) {
	userID := utils.MustGetUserID(c)
	locationID := c.Param("id")

	var req dto.UpdateLocationRequest
	if !utils.BindAndValidateJSON(c, &req) {
		return
	}

	location, err := h.service.UpdateLocation(userID, locationID, &req)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    location,
		"message": "Location updated successfully",
	})
}

// DeleteLocation deletes user's own location (not system defaults)
func (h *LocationHandler) DeleteLocation(c *gin.Context) {
	userID := utils.MustGetUserID(c)
	locationID := c.Param("id")

	if err := h.service.DeleteLocation(userID, locationID); err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Location deleted successfully",
	})
}

// GetLocationByID gets specific location details
func (h *LocationHandler) GetLocationByID(c *gin.Context) {
	userID := utils.MustGetUserID(c)
	locationID := c.Param("id")

	location, err := h.service.GetLocationByID(userID, locationID)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    location,
		"message": "Location retrieved successfully",
	})
}

// GetAssetsByLocation gets all assets in a specific location
func (h *LocationHandler) GetAssetsByLocation(c *gin.Context) {
	userID := utils.MustGetUserID(c)
	locationID := c.Param("id")

	result, err := h.service.GetAssetsByLocation(userID, locationID)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
		"message": "Location assets retrieved successfully",
	})
}
