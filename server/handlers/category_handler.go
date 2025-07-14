package handlers

import (
	"github.com/fiqrioemry/asset_management_system_app/server/dto"
	"github.com/fiqrioemry/asset_management_system_app/server/services"
	"github.com/fiqrioemry/asset_management_system_app/server/utils"
	"github.com/fiqrioemry/go-api-toolkit/response"
	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	service services.CategoryService
}

func NewCategoryHandler(service services.CategoryService) *CategoryHandler {
	return &CategoryHandler{service}
}

func (h *CategoryHandler) GetCategoriesTree(c *gin.Context) {
	userID := utils.MustGetUserID(c)

	categories, err := h.service.GetCategoriesTree(userID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.OK(c, "Categories tree retrieved successfully", categories.Categories)

}

func (h *CategoryHandler) GetCategoriesFlat(c *gin.Context) {
	userID := utils.MustGetUserID(c)

	categories, err := h.service.GetCategoriesFlat(userID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.OK(c, "Categories flat list retrieved successfully", categories.Categories)
}

func (h *CategoryHandler) GetParentCategories(c *gin.Context) {
	userID := utils.MustGetUserID(c)

	categories, err := h.service.GetParentCategories(userID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.OK(c, "Parent categories retrieved successfully", categories.Categories)
}

func (h *CategoryHandler) GetChildCategories(c *gin.Context) {
	parentID := c.Param("id")
	userID := utils.MustGetUserID(c)

	categories, err := h.service.GetChildCategories(parentID, userID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.OK(c, "Child categories retrieved successfully", categories.Categories)

}

func (h *CategoryHandler) GetCategoryByID(c *gin.Context) {
	userID := utils.MustGetUserID(c)
	categoryID := c.Param("id")

	category, err := h.service.GetCategoryByID(userID, categoryID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.OK(c, "Category retrieved successfully", category)
}

func (h *CategoryHandler) GetAssetsByCategory(c *gin.Context) {
	userID := utils.MustGetUserID(c)
	categoryID := c.Param("id")

	result, err := h.service.GetAssetsByCategory(userID, categoryID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.OK(c, "Assets in category retrieved successfully", result)
}

func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	userID := utils.MustGetUserID(c)

	// bind and validate request
	var req dto.CreateCategoryRequest
	if !utils.BindAndValidateJSON(c, &req) {
		return
	}

	category, err := h.service.CreateCategory(userID, &req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Created(c, "Category created successfully", category)
}

func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	userID := utils.MustGetUserID(c)
	categoryID := c.Param("id")

	// bind and validate request
	var req dto.UpdateCategoryRequest
	if !utils.BindAndValidateJSON(c, &req) {
		return
	}

	category, err := h.service.UpdateCategory(userID, categoryID, &req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.OK(c, "Category updated successfully", category)
}

func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	categoryID := c.Param("id")
	userID := utils.MustGetUserID(c)

	if err := h.service.DeleteCategory(userID, categoryID); err != nil {
		response.Error(c, err)
		return
	}

	response.OK(c, "Category deleted successfully", categoryID)
}
