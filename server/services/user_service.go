package services

import (
	"context"
	"fmt"
	"time"

	"github.com/fiqrioemry/asset_management_system_app/server/config"
	"github.com/fiqrioemry/asset_management_system_app/server/dto"
	"github.com/fiqrioemry/asset_management_system_app/server/errors"
	"github.com/fiqrioemry/asset_management_system_app/server/models"
	"github.com/fiqrioemry/asset_management_system_app/server/repositories"
	"github.com/fiqrioemry/asset_management_system_app/server/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"google.golang.org/api/idtoken"
)

type UserService interface {

	// authentication features
	Login(req *dto.LoginRequest) (*dto.AuthResponse, error)
	Register(req *dto.RegisterRequest) (*dto.AuthResponse, error)
	RefreshSession(c *gin.Context, token string) (*dto.UserSession, error)

	// user features
	GetMe(id string) (*dto.UserProfileResponse, error)
	UpdateMe(id string, req *dto.UpdateUserRequest) (*dto.UserProfileResponse, error)

	// change password features
	ChangePassword(userID string, req *dto.ChangePasswordRequest) error

	// password reset features
	ForgotPassword(c *gin.Context, req *dto.ForgotPasswordRequest) error
	ValidateToken(token string) (string, error)
	ResetPassword(req *dto.ResetPasswordRequest) error

	// Google OAuth features
	GetGoogleOAuthURL() string
	GoogleSignIn(tokenId string) (*dto.AuthResponse, error)
	HandleGoogleOAuthCallback(code string) (*dto.AuthResponse, error)
}

type userService struct {
	user repositories.UserRepository
}

func NewUserService(user repositories.UserRepository) UserService {
	return &userService{user: user}
}

func (s *userService) Login(req *dto.LoginRequest) (*dto.AuthResponse, error) {
	// store attempt as cache
	redisKey := fmt.Sprintf("asset_app:login:attempt:%s", req.Email)
	if err := utils.CheckAttempts(redisKey, 5); err != nil {
		return nil, errors.NewTooManyRequests("Too many login attempts, please try again later")
	}

	// check if exist
	user, err := s.user.GetByEmail(req.Email)
	if err != nil || user == nil || !utils.CheckPasswordHash(req.Password, user.Password) {
		// increment attempts cache
		utils.IncrementAttempts(redisKey)
		return nil, errors.NewUnauthorized("Invalid email or password")
	}

	// delete attempts cache
	utils.DeleteKeys(redisKey)

	// generate accessToken
	accessToken, err := utils.GenerateAccessToken(user.ID.String())
	if err != nil {
		return nil, errors.NewInternalServerError("Failed to generate access token", err)
	}

	// generate refreshToken
	refreshToken, err := utils.GenerateRefreshToken(user.ID.String())
	if err != nil {
		return nil, errors.NewInternalServerError("Failed to generate refresh token", err)
	}

	userResponse := dto.UserSession{
		ID:       user.ID.String(),
		Email:    user.Email,
		Fullname: user.Fullname,
		Avatar:   user.Avatar,
	}

	return &dto.AuthResponse{
		User:         userResponse,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil

}

func (s *userService) Register(req *dto.RegisterRequest) (*dto.AuthResponse, error) {
	// check user exist
	user, err := s.user.GetByEmail(req.Email)
	if err != nil {
		return nil, errors.NewInternalServerError("Failed to check user existence", err)
	}

	if user != nil {
		return nil, errors.NewConflict("Email already registered")
	}

	// hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, errors.NewInternalServerError("Failed to hash password", err)
	}

	newUser := models.User{
		Email:    req.Email,
		Fullname: req.Fullname,
		Password: hashedPassword,
		Avatar:   utils.RandomUserAvatar(req.Fullname),
	}

	// create new user
	if err := s.user.Create(&newUser); err != nil {
		return nil, errors.NewConflict("Email is already registered")
	}

	userResponse := dto.UserSession{
		ID:       newUser.ID.String(),
		Email:    newUser.Email,
		Fullname: newUser.Fullname,
		Avatar:   newUser.Avatar,
	}

	// generate accessToken
	accessToken, err := utils.GenerateAccessToken(newUser.ID.String())
	if err != nil {
		return nil, errors.NewInternalServerError("Failed to generate access token", err)
	}
	// generate refreshToken
	refreshToken, err := utils.GenerateRefreshToken(newUser.ID.String())
	if err != nil {
		return nil, errors.NewInternalServerError("Failed to generate refresh token", err)
	}

	return &dto.AuthResponse{
		User:         userResponse,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil

}

func (s *userService) RefreshSession(c *gin.Context, token string) (*dto.UserSession, error) {
	// decode refreshToken
	userID, err := utils.DecodeRefreshToken(token)
	if err != nil {
		return nil, errors.NewUnauthorized("Invalid refresh token")
	}

	// check user exists
	user, err := s.user.GetByID(userID)
	if err != nil || user == nil {
		return nil, errors.NewNotFound("User not found").WithContext("userID", userID)
	}

	userResponse := dto.UserSession{
		ID:       user.ID.String(),
		Email:    user.Email,
		Fullname: user.Fullname,
		Avatar:   user.Avatar,
	}

	// generate accessToken
	accessToken, err := utils.GenerateAccessToken(user.ID.String())
	if err != nil {
		return nil, errors.NewInternalServerError("Failed to generate access token", err)
	}

	// set accessToken
	utils.SetAccessTokenCookie(c, accessToken)

	return &userResponse, nil
}

func (s *userService) GetMe(id string) (*dto.UserProfileResponse, error) {

	// check user exists
	user, err := s.user.GetByID(id)
	if err != nil || user == nil {
		return nil, errors.NewNotFound("User not found")
	}

	return &dto.UserProfileResponse{
		ID:       user.ID.String(),
		Email:    user.Email,
		Fullname: user.Fullname,
		Avatar:   user.Avatar,
		JoinedAt: user.CreatedAt,
	}, nil
}

func (s *userService) UpdateMe(id string, req *dto.UpdateUserRequest) (*dto.UserProfileResponse, error) {
	// check user exists
	user, err := s.user.GetByID(id)
	if err != nil || user == nil {
		return nil, errors.NewNotFound("User not found")
	}

	if req.Fullname != "" {
		user.Fullname = req.Fullname
	}

	if req.AvatarURL != "" {
		user.Avatar = req.AvatarURL
	}
	// update user
	if err := s.user.Update(user); err != nil {
		return nil, errors.NewInternalServerError("Failed to update user", err)
	}

	return &dto.UserProfileResponse{
		ID:       user.ID.String(),
		Email:    user.Email,
		Fullname: user.Fullname,
		Avatar:   user.Avatar,
		JoinedAt: user.CreatedAt,
	}, nil
}

func (s *userService) ChangePassword(userID string, req *dto.ChangePasswordRequest) error {
	// Validate confirm password
	if req.NewPassword != req.ConfirmPassword {
		return errors.NewBadRequest("New password and confirm password don't match")
	}

	// check if user exists
	user, err := s.user.GetByID(userID)
	if err != nil {
		return errors.NewNotFound("User not found")
	}

	// Verify current password
	if !utils.CheckPasswordHash(req.CurrentPassword, user.Password) {
		return errors.NewBadRequest("Current password is incorrect")
	}

	// Hash new password
	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return errors.NewInternalServerError("Failed to hash password", err)
	}

	// Update password
	user.Password = hashedPassword
	if err := s.user.Update(user); err != nil {
		return errors.NewInternalServerError("Failed to update password", err)
	}

	return nil
}

func (s *userService) ForgotPassword(c *gin.Context, req *dto.ForgotPasswordRequest) error {
	// Check rate limit attempts
	attemptsKey := "asset_app:forgot_password_attempts:" + c.ClientIP()
	if err := utils.CheckForgotPasswordAttempts(c.ClientIP(), 3); err != nil {
		return errors.NewTooManyRequests("Too many forgot password attempts, please try again later")
	}

	// Check token existence
	existingTokenKey := "asset_app:reset_token:" + req.Email
	if utils.KeyExists(existingTokenKey) {
		return errors.NewTooManyRequests("Password reset link has already been sent. Please check your email or wait before requesting again.")
	}

	// Check if email exists
	user, err := s.user.GetByEmail(req.Email)
	if err != nil {
		// Increment attempts
		utils.IncrementAttempts(attemptsKey)
		return nil // Don't reveal if email exists for security reasons
	}

	if user == nil {
		utils.IncrementAttempts(attemptsKey)
		return nil // Don't reveal if email exists
	}

	// Generate reset token
	resetToken, err := utils.GenerateResetToken()
	if err != nil {
		return errors.NewInternalServerError("Failed to generate reset token", err)
	}

	// Prepare token data
	tokenData := map[string]any{
		"userId":    user.ID.String(),
		"email":     user.Email,
		"createdAt": time.Now().Unix(),
		"expiresAt": time.Now().Add(1 * time.Hour).Unix(),
	}

	// Store reset token data
	resetTokenKey := "asset_app:password_reset:" + resetToken
	if err := utils.AddKeys(resetTokenKey, tokenData, 1*time.Hour); err != nil {
		return errors.NewInternalServerError("Failed to store reset token", err)
	}

	// Store email -
	emailTokenKey := "asset_app:reset_token:" + user.Email
	if err := utils.AddKeys(emailTokenKey, resetToken, 1*time.Hour); err != nil {
		// Clean up reset token
		utils.DeleteKeys(resetTokenKey)
		return errors.NewInternalServerError("Failed to store email token mapping", err)
	}

	// Create reset link
	frontendURL := config.AppConfig.FrontendURL

	resetLink := fmt.Sprintf("%s/reset-password?token=%s", frontendURL, resetToken)

	// Send reset password email
	if err := utils.SendResetPasswordEmail(user.Email, user.Fullname, resetLink, 1*time.Hour); err != nil {
		// Clean up tokens
		utils.DeleteKeys(resetTokenKey, emailTokenKey)
		return errors.NewInternalServerError("Failed to send reset password email", err)
	}

	// Increment attempts
	utils.IncrementAttempts(attemptsKey)

	return nil
}

func (s *userService) ResetPassword(req *dto.ResetPasswordRequest) error {

	// check password match
	if req.NewPassword != req.ConfirmPassword {
		return errors.NewBadRequest("New password and confirm password do not match")
	}

	resetTokenKey := "asset_app:password_reset:" + req.Token
	var tokenData map[string]any

	// get token from cache
	if err := utils.GetKey(resetTokenKey, &tokenData); err != nil {
		return errors.NewBadRequest("Invalid or expired reset token")
	}

	email, ok := tokenData["email"].(string)
	if !ok {
		return errors.NewBadRequest("Invalid reset token data")
	}

	userID, ok := tokenData["userId"].(string)
	if !ok {
		return errors.NewBadRequest("Invalid reset token data")
	}

	// check user exists
	user, err := s.user.GetByID(userID)
	if err != nil {
		return errors.NewNotFound("User not found")
	}

	// hash new password
	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return errors.NewInternalServerError("Failed to hash password", err)
	}

	// update user password
	user.Password = hashedPassword

	if err := s.user.Update(user); err != nil {
		return errors.NewInternalServerError("Failed to update password", err)
	}

	// delete related cache keys
	utils.DeleteKeys(resetTokenKey)
	utils.DeleteKeys("asset_app:reset_token:" + email)
	utils.DeleteKeys("asset_app:forgot_password_attempts:" + email)

	return nil
}

func (s *userService) ValidateToken(token string) (string, error) {

	resetTokenKey := "asset_app:password_reset:" + token
	var tokenData map[string]any
	// check token existence
	if err := utils.GetKey(resetTokenKey, &tokenData); err != nil {
		return "", errors.NewBadRequest("Invalid or expired reset token")
	}

	// Check token age
	createdAt, ok := tokenData["createdAt"].(float64)
	if !ok {
		return "", errors.NewBadRequest("Invalid token data")
	}

	tokenAge := time.Since(time.Unix(int64(createdAt), 0))
	if tokenAge > 1*time.Hour {
		utils.DeleteKeys(resetTokenKey)
		return "", errors.NewBadRequest("Reset token has expired")
	}

	// email for display
	email, _ := tokenData["email"].(string)

	return email, nil

}

func (s *userService) GoogleSignIn(tokenId string) (*dto.AuthResponse, error) {
	payload, err := idtoken.Validate(context.Background(), tokenId, config.AppConfig.GoogleClientID)
	if err != nil {
		return nil, errors.NewUnauthorized("Invalid Google ID token")
	}

	email, ok := payload.Claims["email"].(string)
	if !ok || email == "" {
		return nil, errors.NewNotFound("Email not found in token")
	}

	name, _ := payload.Claims["name"].(string)

	user, err := s.user.GetByEmail(email)
	if err != nil {
		user = &models.User{
			Email:    email,
			Avatar:   utils.RandomUserAvatar(name),
			Fullname: name,
			Password: "-",
		}

		if err := s.user.Create(user); err != nil {
			return nil, err
		}

		if user.ID == uuid.Nil {
			return nil, errors.NewInternalServerError("Failed to assign UUID to user", err)
		}
	}

	accessToken, err := utils.GenerateAccessToken(user.ID.String())
	if err != nil {
		return nil, err
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID.String())
	if err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *userService) GetGoogleOAuthURL() string {
	return config.GoogleOAuthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
}

func (s *userService) HandleGoogleOAuthCallback(code string) (*dto.AuthResponse, error) {
	token, err := config.GoogleOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, errors.NewUnauthorized("Failed to exchange Google OAuth code")
	}

	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return nil, errors.NewUnauthorized("ID token not found in Google OAuth response")
	}

	return s.GoogleSignIn(rawIDToken)
}
