package handlers

import (
	"net/http"

	"github.com/fiqrioemry/asset_management_system_app/server/dto"
	"github.com/fiqrioemry/asset_management_system_app/server/services"
	"github.com/fiqrioemry/asset_management_system_app/server/utils"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service services.UserService
}

func NewUserHandler(service services.UserService) *UserHandler {
	return &UserHandler{service}
}

func (h *UserHandler) Login(c *gin.Context) {

	var req dto.LoginRequest
	if !utils.BindAndValidateJSON(c, &req) {
		return
	}

	response, err := h.service.Login(&req)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	utils.SetAccessTokenCookie(c, response.AccessToken)

	utils.SetRefreshTokenCookie(c, response.RefreshToken)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"user":    response.User,
		"message": "Login successful",
	})
}

func (h *UserHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if !utils.BindAndValidateJSON(c, &req) {
		return
	}

	response, err := h.service.Register(&req)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	utils.SetAccessTokenCookie(c, response.AccessToken)

	utils.SetRefreshTokenCookie(c, response.RefreshToken)

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"user":    response.User,
		"message": "Registration successful",
	})

}

func (h *UserHandler) Logout(c *gin.Context) {
	utils.ClearAccessTokenCookie(c)
	utils.ClearRefreshTokenCookie(c)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Logout successful",
	})
}

func (h *UserHandler) RefreshSession(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")

	if err != nil {
		utils.HandleError(c, err)
		return
	}

	user, err := h.service.RefreshSession(c, refreshToken)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"user":    user,
		"message": "Session refreshed successfully",
	})
}

func (h *UserHandler) GetMe(c *gin.Context) {
	userID := utils.MustGetUserID(c)

	user, err := h.service.GetMe(userID)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"user":    user,
		"message": "User retrieved successfully",
	})
}

func (h *UserHandler) UpdateMe(c *gin.Context) {
	userID := utils.MustGetUserID(c)
	var req dto.UpdateUserRequest
	if !utils.BindAndValidateJSON(c, &req) {
		return
	}

	avatarURL, err := utils.UploadImageWithValidation(req.Avatar)
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	req.AvatarURL = avatarURL

	updatedUser, err := h.service.UpdateMe(userID, &req)
	if err != nil {
		utils.CleanupImageOnError(req.AvatarURL)
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"user":    updatedUser,
		"message": "User updated successfully",
	})
}
