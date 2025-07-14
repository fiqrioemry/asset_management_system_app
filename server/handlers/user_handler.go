package handlers

import (
	"net/http"

	"github.com/fiqrioemry/asset_management_system_app/server/config"
	"github.com/fiqrioemry/asset_management_system_app/server/dto"
	"github.com/fiqrioemry/asset_management_system_app/server/services"
	"github.com/fiqrioemry/asset_management_system_app/server/utils"
	"github.com/fiqrioemry/go-api-toolkit/response"
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

	loginResponse, err := h.service.Login(&req)
	if err != nil {
		response.Error(c, err)
		return
	}

	utils.SetAccessTokenCookie(c, loginResponse.AccessToken)

	utils.SetRefreshTokenCookie(c, loginResponse.RefreshToken)

	response.OK(c, "Asset retrieved successfully", loginResponse.User)
}

func (h *UserHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if !utils.BindAndValidateJSON(c, &req) {
		return
	}

	registerResponse, err := h.service.Register(&req)
	if err != nil {
		response.Error(c, err)
		return
	}

	utils.SetAccessTokenCookie(c, registerResponse.AccessToken)

	utils.SetRefreshTokenCookie(c, registerResponse.RefreshToken)

	response.Created(c, "Registration successful", registerResponse.User)

}

func (h *UserHandler) Logout(c *gin.Context) {
	utils.ClearAccessTokenCookie(c)
	utils.ClearRefreshTokenCookie(c)

	response.OK(c, "Logout successful", nil)
}

func (h *UserHandler) RefreshSession(c *gin.Context) {
	refreshToken, err := c.Cookie("refreshToken")

	if err != nil {
		response.Error(c, err)
		return
	}

	user, err := h.service.RefreshSession(c, refreshToken)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.OK(c, "Session refreshed successfully", user)
}

func (h *UserHandler) GetMe(c *gin.Context) {
	userID := utils.MustGetUserID(c)

	user, err := h.service.GetMe(userID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.OK(c, "User retrieved successfully", user)
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
			response.Error(c, err)
			return
		}
		req.AvatarURL = avatarURL
	}

	updatedUser, err := h.service.UpdateMe(userID, &req)
	if err != nil {
		utils.CleanupImageOnError(req.AvatarURL)
		response.Error(c, err)
		return
	}

	response.OK(c, "User updated successfully", updatedUser)
}

func (h *UserHandler) ChangePassword(c *gin.Context) {
	userID := utils.MustGetUserID(c)

	var req dto.ChangePasswordRequest
	if !utils.BindAndValidateJSON(c, &req) {
		return
	}

	if err := h.service.ChangePassword(userID, &req); err != nil {
		response.Error(c, err)
		return
	}

	response.OK(c, "Password changed successfully", nil)
}

// step 1 : User requests password reset
func (h *UserHandler) ForgotPassword(c *gin.Context) {

	var req dto.ForgotPasswordRequest
	if !utils.BindAndValidateJSON(c, &req) {
		return
	}

	if err := h.service.ForgotPassword(c, &req); err != nil {
		response.Error(c, err)
		return
	}

	response.OK(c, "Password reset link sent successfully", nil)

}

// step 2 : validate reset token and reset password
func (h *UserHandler) ValidateResetToken(c *gin.Context) {
	token := c.Query("token")

	email, err := h.service.ValidateToken(token)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.OK(c, "Reset token is valid", email)
}

// step 3 : reset password
func (h *UserHandler) ResetPassword(c *gin.Context) {
	var req dto.ResetPasswordRequest
	if !utils.BindAndValidateJSON(c, &req) {
		return
	}

	if err := h.service.ResetPassword(&req); err != nil {
		response.Error(c, err)
		return
	}

	response.OK(c, "Password has been reset successfully", nil)
}

func (h *UserHandler) GoogleOAuthRedirect(c *gin.Context) {
	url := h.service.GetGoogleOAuthURL()
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *UserHandler) GoogleOAuthCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		response.Error(c, response.NewBadRequest("Authorization code is missing"))
		return
	}

	tokens, err := h.service.HandleGoogleOAuthCallback(code)
	if err != nil {
		response.Error(c, err)
		return
	}

	utils.SetAccessTokenCookie(c, tokens.AccessToken)

	utils.SetRefreshTokenCookie(c, tokens.RefreshToken)

	c.Redirect(http.StatusTemporaryRedirect, config.AppConfig.FrontendRedirectURL)
}
