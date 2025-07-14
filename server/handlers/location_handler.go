package handlers

import (
	"github.com/fiqrioemry/asset_management_system_app/server/dto"
	"github.com/fiqrioemry/asset_management_system_app/server/services"
	"github.com/fiqrioemry/asset_management_system_app/server/utils"
	"github.com/fiqrioemry/go-api-toolkit/response"
	"github.com/gin-gonic/gin"
)

type LocationHandler struct {
	service services.LocationService
}

func NewLocationHandler(service services.LocationService) *LocationHandler {
	return &LocationHandler{service}
}

func (h *LocationHandler) GetLocations(c *gin.Context) {
	userID := utils.MustGetUserID(c)

	locationResp, err := h.service.GetLocations(userID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.OK(c, "Locations retrieved successfully", locationResp.Locations)
}

func (h *LocationHandler) CreateLocation(c *gin.Context) {
	userID := utils.MustGetUserID(c)

	var req dto.CreateLocationRequest
	if !utils.BindAndValidateJSON(c, &req) {
		return
	}

	location, err := h.service.CreateLocation(userID, &req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Created(c, "Location created successfully", location)
}

func (h *LocationHandler) UpdateLocation(c *gin.Context) {
	userID := utils.MustGetUserID(c)
	locationID := c.Param("id")

	var req dto.UpdateLocationRequest
	if !utils.BindAndValidateJSON(c, &req) {
		return
	}

	location, err := h.service.UpdateLocation(userID, locationID, &req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.OK(c, "Location updated successfully", location)

}

func (h *LocationHandler) DeleteLocation(c *gin.Context) {
	userID := utils.MustGetUserID(c)
	locationID := c.Param("id")

	if err := h.service.DeleteLocation(userID, locationID); err != nil {
		response.Error(c, err)
		return
	}

	response.OK(c, "Location deleted successfully", locationID)
}

func (h *LocationHandler) GetLocationByID(c *gin.Context) {
	userID := utils.MustGetUserID(c)
	locationID := c.Param("id")

	location, err := h.service.GetLocationByID(userID, locationID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.OK(c, "Location retrieved successfully", location)
}

func (h *LocationHandler) GetAssetsByLocation(c *gin.Context) {
	userID := utils.MustGetUserID(c)
	locationID := c.Param("id")

	result, err := h.service.GetAssetsByLocation(userID, locationID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.OK(c, "Assets in location retrieved successfully", result)
}
