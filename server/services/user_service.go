package services

import (
	"fmt"
	"time"

	"github.com/fiqrioemry/asset_management_system_app/server/config"
	"github.com/fiqrioemry/asset_management_system_app/server/dto"
	"github.com/fiqrioemry/asset_management_system_app/server/errors"
	"github.com/fiqrioemry/asset_management_system_app/server/models"
	"github.com/fiqrioemry/asset_management_system_app/server/repositories"
	"github.com/fiqrioemry/asset_management_system_app/server/utils"
	"github.com/gin-gonic/gin"
)

type UserService interface {
	Login(req *dto.LoginRequest) (*dto.AuthResponse, error)
	Register(req *dto.RegisterRequest) (*dto.AuthResponse, error)
	GetMe(id string) (*dto.UserResponse, error)
	UpdateMe(id string, req *dto.UpdateUserRequest) (*dto.UserResponse, error)
	RefreshSession(c *gin.Context, token string) (*dto.UserResponse, error)
}

type userService struct {
	user repositories.UserRepository
}

func NewUserService(user repositories.UserRepository) UserService {
	return &userService{user: user}
}

func (s *userService) Login(req *dto.LoginRequest) (*dto.AuthResponse, error) {

	redisKey := fmt.Sprintf("login:attempt:%s", req.Email)
	attempts, _ := config.RedisClient.Get(config.Ctx, redisKey).Int()
	if attempts >= 5 {
		return nil, errors.NewTooManyRequests("Too many login attempts, please try again later")
	}

	user, err := s.user.GetByEmail(req.Email)
	if err != nil || user == nil || !utils.CheckPasswordHash(req.Password, user.Password) {
		config.RedisClient.Incr(config.Ctx, redisKey)
		config.RedisClient.Expire(config.Ctx, redisKey, 30*time.Minute)
		return nil, errors.NewUnauthorized("Invalid email or password")
	}

	config.RedisClient.Del(config.Ctx, redisKey)

	accessToken, err := utils.GenerateAccessToken(user.ID.String())
	if err != nil {
		return nil, errors.NewInternalServerError("Failed to generate access token", err)
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID.String())
	if err != nil {
		return nil, errors.NewInternalServerError("Failed to generate refresh token", err)
	}

	userResponse := dto.UserResponse{
		ID:       user.ID.String(),
		Email:    user.Email,
		Fullname: user.Fullname,
		Avatar:   user.Avatar,
		JoinedAt: user.CreatedAt,
	}

	return &dto.AuthResponse{
		User:         userResponse,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil

}

func (s *userService) Register(req *dto.RegisterRequest) (*dto.AuthResponse, error) {
	user, err := s.user.GetByEmail(req.Email)
	if err != nil {
		return nil, errors.NewInternalServerError("Failed to check user existence", err)
	}

	if user != nil {
		return nil, errors.NewAlreadyExists("Email already registered")
	}

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

	if err := s.user.Create(&newUser); err != nil {
		return nil, errors.NewConflict("Email is already registered")
	}

	userResponse := dto.UserResponse{
		ID:       newUser.ID.String(),
		Email:    newUser.Email,
		Fullname: newUser.Fullname,
		Avatar:   newUser.Avatar,
		JoinedAt: newUser.CreatedAt,
	}

	accessToken, err := utils.GenerateAccessToken(newUser.ID.String())
	if err != nil {
		return nil, errors.NewInternalServerError("Failed to generate access token", err)
	}
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

func (s *userService) RefreshSession(c *gin.Context, token string) (*dto.UserResponse, error) {

	userID, err := utils.DecodeRefreshToken(token)
	if err != nil {
		return nil, errors.NewUnauthorized("Invalid refresh token")
	}

	user, err := s.user.GetByID(userID)
	if err != nil || user == nil {
		return nil, errors.NewNotFound("User not found").WithContext("userID", userID)
	}

	userResponse := dto.UserResponse{
		ID:       user.ID.String(),
		Email:    user.Email,
		Fullname: user.Fullname,
		Avatar:   user.Avatar,
		JoinedAt: user.CreatedAt,
	}

	accessToken, err := utils.GenerateAccessToken(user.ID.String())
	if err != nil {
		return nil, errors.NewInternalServerError("Failed to generate access token", err)
	}

	utils.SetAccessTokenCookie(c, accessToken)

	return &userResponse, nil
}

func (s *userService) GetMe(id string) (*dto.UserResponse, error) {
	user, err := s.user.GetByID(id)
	if err != nil || user == nil {
		return nil, errors.NewNotFound("User not found").WithContext("userID", id)
	}

	return &dto.UserResponse{
		ID:       user.ID.String(),
		Email:    user.Email,
		Fullname: user.Fullname,
		Avatar:   user.Avatar,
		JoinedAt: user.CreatedAt,
	}, nil
}

func (s *userService) UpdateMe(id string, req *dto.UpdateUserRequest) (*dto.UserResponse, error) {
	user, err := s.user.GetByID(id)
	if err != nil || user == nil {
		return nil, errors.NewNotFound("User not found").WithContext("userID", id)
	}

	if req.Fullname != "" {
		user.Fullname = req.Fullname
	}

	if req.AvatarURL != "" {
		user.Avatar = req.AvatarURL
	}

	if err := s.user.Update(user); err != nil {
		return nil, errors.NewInternalServerError("Failed to update user", err)
	}

	return &dto.UserResponse{
		ID:       user.ID.String(),
		Email:    user.Email,
		Fullname: user.Fullname,
		Avatar:   user.Avatar,
		JoinedAt: user.CreatedAt,
	}, nil
}
