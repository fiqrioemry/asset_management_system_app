package handlers

import (
	"net/http"

	"github.com/fiqrioemry/asset_management_system_app/server/dto"
	"github.com/fiqrioemry/asset_management_system_app/server/services"
	"github.com/fiqrioemry/asset_management_system_app/server/utils"

	"github.com/gin-gonic/gin"
)

type AssetHandler struct {
	service services.AssetService
}

func NewAssetHandler(service services.AssetService) *AssetHandler {
	return &AssetHandler{service}
}

func (h *AssetHandler) GetAssets(c *gin.Context) {
	userID := utils.MustGetUserID(c)

	var req dto.GetAssetsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.HandleError(c, err)
		return
	}

	assets, err := h.service.GetAssets(userID, &req)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    assets,
		"message": "Assets retrieved successfully",
	})
}

func (h *AssetHandler) GetAssetByID(c *gin.Context) {
	userID := utils.MustGetUserID(c)
	assetID := c.Param("id")

	asset, err := h.service.GetAssetByID(userID, assetID)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    asset,
		"message": "Asset retrieved successfully",
	})
}

func (h *AssetHandler) CreateAsset(c *gin.Context) {
	userID := utils.MustGetUserID(c)

	var req dto.CreateAssetRequest
	if !utils.BindAndValidateForm(c, &req) {
		return
	}

	// Handle image upload
	if req.Image != nil {
		imageURL, err := utils.UploadImageWithValidation(req.Image)
		if err != nil {
			utils.HandleError(c, err)
			return
		}
		req.ImageURL = imageURL
	}

	asset, err := h.service.CreateAsset(userID, &req)
	if err != nil {
		utils.CleanupImageOnError(req.ImageURL)
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    asset,
		"message": "Asset created successfully",
	})
}

func (h *AssetHandler) UpdateAsset(c *gin.Context) {
	userID := utils.MustGetUserID(c)
	assetID := c.Param("id")

	var req dto.UpdateAssetRequest
	if !utils.BindAndValidateForm(c, &req) {
		return
	}

	if req.Image != nil {
		imageURL, err := utils.UploadImageWithValidation(req.Image)
		if err != nil {
			utils.HandleError(c, err)
			return
		}
		req.ImageURL = imageURL
	}

	asset, err := h.service.UpdateAsset(userID, assetID, &req)
	if err != nil {
		utils.CleanupImageOnError(req.ImageURL)
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    asset,
		"message": "Asset updated successfully",
	})
}

func (h *AssetHandler) DeleteAsset(c *gin.Context) {
	userID := utils.MustGetUserID(c)
	assetID := c.Param("id")

	if err := h.service.DeleteAsset(userID, assetID); err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Asset deleted successfully",
	})
}
