package handlers

import (
	"net/http"

	"github.com/fiqrioemry/asset_management_system_app/server/config"
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
	if !utils.BindAndValidateForm(c, &req) {
		return
	}

	if req.Avatar != nil {
		avatarURL, err := utils.UploadImageWithValidation(req.Avatar)
		if err != nil {
			utils.HandleError(c, err)
			return
		}
		req.AvatarURL = avatarURL
	}

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

func (h *UserHandler) ChangePassword(c *gin.Context) {
	userID := utils.MustGetUserID(c)

	var req dto.ChangePasswordRequest
	if !utils.BindAndValidateJSON(c, &req) {
		return
	}

	if err := h.service.ChangePassword(userID, &req); err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Password changed successfully",
	})
}

// step 1 : User requests password reset
func (h *UserHandler) ForgotPassword(c *gin.Context) {
	var req dto.ForgotPasswordRequest
	if !utils.BindAndValidateJSON(c, &req) {
		return
	}

	if err := h.service.ForgotPassword(c, &req); err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "If an account with that email exists, we have sent a password reset link.",
	})
}

// step 2 : validate reset token and reset password
func (h *UserHandler) ValidateResetToken(c *gin.Context) {
	token := c.Query("token")

	email, err := h.service.ValidateToken(token)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"email":   email,
		"message": "Reset token is valid",
	})
}

// step 3 : reset password
func (h *UserHandler) ResetPassword(c *gin.Context) {
	var req dto.ResetPasswordRequest
	if !utils.BindAndValidateJSON(c, &req) {
		return
	}

	if err := h.service.ResetPassword(&req); err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Password has been reset successfully",
	})
}

func (h *UserHandler) GoogleOAuthRedirect(c *gin.Context) {
	url := h.service.GetGoogleOAuthURL()
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *UserHandler) GoogleOAuthCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Authorization code is missing"})
		return
	}

	tokens, err := h.service.HandleGoogleOAuthCallback(code)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	utils.SetAccessTokenCookie(c, tokens.AccessToken)

	utils.SetRefreshTokenCookie(c, tokens.RefreshToken)

	c.Redirect(http.StatusTemporaryRedirect, config.AppConfig.FrontendRedirectURL)
}
