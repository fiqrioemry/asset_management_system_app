package handlers

import (
	"github.com/fiqrioemry/asset_management_system_app/server/dto"
	"github.com/fiqrioemry/asset_management_system_app/server/services"
	"github.com/fiqrioemry/asset_management_system_app/server/utils"

	"github.com/fiqrioemry/go-api-toolkit/pagination"
	"github.com/fiqrioemry/go-api-toolkit/response"

	"github.com/gin-gonic/gin"
)

type AssetHandler struct {
	service services.AssetService
}

func NewAssetHandler(service services.AssetService) *AssetHandler {
	return &AssetHandler{service}
}

func (h *AssetHandler) CreateAsset(c *gin.Context) {
	userID := utils.MustGetUserID(c)

	var req dto.CreateAssetRequest
	if !utils.BindAndValidateForm(c, &req) {
		return
	}

	// Handle image upload
	if req.Image != nil && req.Image.Filename != "" {
		imageURL, err := utils.UploadImageWithValidation(req.Image)
		if err != nil {
			response.Error(c, err)
			return
		}
		req.ImageURL = imageURL
	}

	asset, err := h.service.CreateAsset(userID, &req)
	if err != nil {
		utils.CleanupImageOnError(req.ImageURL)
		response.Error(c, err)
		return
	}

	response.Created(c, "Asset created successfully", asset)
}

func (h *AssetHandler) GetAssets(c *gin.Context) {
	userID := utils.MustGetUserID(c)
	var req dto.GetAssetsRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, response.NewBadRequest("Invalid query parameters"))
		return
	}

	// apply pagination defaults
	if err := pagination.BindAndSetDefaults(c, &req); err != nil {
		response.Error(c, response.BadRequest(err.Error()))
		return
	}

	assets, total, err := h.service.GetAssets(userID, &req)
	if err != nil {
		response.Error(c, err)
		return
	}

	pag := pagination.Build(req.Page, req.Limit, total)

	response.OKWithPagination(c, "Assets retrieved successfully", assets, pag)
}

func (h *AssetHandler) GetAssetByID(c *gin.Context) {
	userID := utils.MustGetUserID(c)
	assetID := c.Param("id")

	asset, err := h.service.GetAssetByID(userID, assetID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.OK(c, "Asset retrieved successfully", asset)
}

func (h *AssetHandler) UpdateAsset(c *gin.Context) {
	assetID := c.Param("id")
	userID := utils.MustGetUserID(c)

	var req dto.UpdateAssetRequest
	if !utils.BindAndValidateForm(c, &req) {
		return
	}
	if req.Image != nil && req.Image.Filename != "" {
		imageURL, err := utils.UploadImageWithValidation(req.Image)
		if err != nil {
			response.Error(c, err)
			return
		}
		req.ImageURL = imageURL
	}

	asset, err := h.service.UpdateAsset(userID, assetID, &req)
	if err != nil {
		utils.CleanupImageOnError(req.ImageURL)
		response.Error(c, err)
		return
	}

	response.OK(c, "Asset updated successfully", asset)
}

func (h *AssetHandler) DeleteAsset(c *gin.Context) {
	assetID := c.Param("id")
	userID := utils.MustGetUserID(c)

	if err := h.service.DeleteAsset(userID, assetID); err != nil {
		response.Error(c, err)
		return
	}

	response.OK(c, "Asset deleted successfully", assetID)
}
