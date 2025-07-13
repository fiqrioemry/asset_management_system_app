package handlers

import (
	"github.com/fiqrioemry/asset_management_system_app/server/dto"
	"github.com/fiqrioemry/asset_management_system_app/server/errors"
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
		utils.HandleError(c, errors.NewBadRequest("Invalid query parameters"))
		return
	}

	if req.Page == 0 {
		req.Page = 1
	}
	if req.Limit == 0 {
		req.Limit = 10
	}

	if req.SortBy == "" {
		req.SortBy = "created_at"
	}
	if req.SortOrder == "" {
		req.SortOrder = "desc"
	}

	assets, pagination, err := h.service.GetAssets(userID, &req)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	utils.OKWithPagination(c, "Assets retrieved successfully", assets, pagination)
}

func (h *AssetHandler) GetAssetByID(c *gin.Context) {
	userID := utils.MustGetUserID(c)
	assetID := c.Param("id")

	asset, err := h.service.GetAssetByID(userID, assetID)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	utils.OK(c, "Asset retrieved successfully", asset)
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

	utils.Created(c, "Asset created successfully", asset)
}

func (h *AssetHandler) UpdateAsset(c *gin.Context) {
	assetID := c.Param("id")
	userID := utils.MustGetUserID(c)

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

	utils.OK(c, "Asset updated successfully", asset)
}

func (h *AssetHandler) DeleteAsset(c *gin.Context) {
	assetID := c.Param("id")
	userID := utils.MustGetUserID(c)

	if err := h.service.DeleteAsset(userID, assetID); err != nil {
		utils.HandleError(c, err)
		return
	}

	utils.OK(c, "Asset deleted successfully", assetID)
}
