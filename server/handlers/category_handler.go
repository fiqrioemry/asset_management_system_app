package handlers

import (
	"net/http"

	"github.com/fiqrioemry/asset_management_system_app/server/dto"
	"github.com/fiqrioemry/asset_management_system_app/server/services"
	"github.com/fiqrioemry/asset_management_system_app/server/utils"
	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	service services.CategoryService
}

func NewCategoryHandler(service services.CategoryService) *CategoryHandler {
	return &CategoryHandler{service}
}

// ======================== READ OPERATIONS ========================

// GetCategoriesTree returns hierarchical tree structure
func (h *CategoryHandler) GetCategoriesTree(c *gin.Context) {
	userID := utils.MustGetUserID(c)

	categories, err := h.service.GetCategoriesTree(userID)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    categories,
		"message": "Categories tree retrieved successfully",
	})
}

// GetCategoriesFlat returns flat list with full paths
func (h *CategoryHandler) GetCategoriesFlat(c *gin.Context) {
	userID := utils.MustGetUserID(c)

	categories, err := h.service.GetCategoriesFlat(userID)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    categories,
		"message": "Categories flat list retrieved successfully",
	})
}

// GetParentCategories returns only parent categories
func (h *CategoryHandler) GetParentCategories(c *gin.Context) {
	userID := utils.MustGetUserID(c)

	categories, err := h.service.GetParentCategories(userID)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    categories,
		"message": "Parent categories retrieved successfully",
	})
}

// GetChildCategories returns children of specific parent
func (h *CategoryHandler) GetChildCategories(c *gin.Context) {
	userID := utils.MustGetUserID(c)
	parentID := c.Param("id")

	categories, err := h.service.GetChildCategories(parentID, userID)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    categories,
		"message": "Child categories retrieved successfully",
	})
}

// GetCategoryByID gets specific category details
func (h *CategoryHandler) GetCategoryByID(c *gin.Context) {
	userID := utils.MustGetUserID(c)
	categoryID := c.Param("id")

	category, err := h.service.GetCategoryByID(userID, categoryID)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    category,
		"message": "Category retrieved successfully",
	})
}

// GetAssetsByCategory gets all assets in specific category
func (h *CategoryHandler) GetAssetsByCategory(c *gin.Context) {
	userID := utils.MustGetUserID(c)
	categoryID := c.Param("id")

	result, err := h.service.GetAssetsByCategory(userID, categoryID)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
		"message": "Category assets retrieved successfully",
	})
}

// CreateCategory creates new category (parent or child)
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	userID := utils.MustGetUserID(c)

	var req dto.CreateCategoryRequest
	if !utils.BindAndValidateJSON(c, &req) {
		return
	}

	category, err := h.service.CreateCategory(userID, &req)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    category,
		"message": "Category created successfully",
	})
}

// UpdateCategory updates user's own category
func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	userID := utils.MustGetUserID(c)
	categoryID := c.Param("id")

	var req dto.UpdateCategoryRequest
	if !utils.BindAndValidateJSON(c, &req) {
		return
	}

	category, err := h.service.UpdateCategory(userID, categoryID, &req)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    category,
		"message": "Category updated successfully",
	})
}

// DeleteCategory deletes user's own category
func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	userID := utils.MustGetUserID(c)
	categoryID := c.Param("id")

	if err := h.service.DeleteCategory(userID, categoryID); err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Category deleted successfully",
	})
}
