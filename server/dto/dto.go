package dto

import (
	"mime/multipart"
	"time"
)

type UserSession struct {
	ID       string `json:"id"`
	Fullname string `json:"fullname"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
}

type AuthResponse struct {
	User         UserSession `json:"user"`
	AccessToken  string      `json:"accessToken"`
	RefreshToken string      `json:"refreshToken"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type RegisterRequest struct {
	Fullname string `json:"fullname" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type UserProfileResponse struct {
	ID       string    `json:"id" binding:"required"`
	Fullname string    `json:"fullname" binding:"required"`
	Email    string    `json:"email" binding:"required,email"`
	Avatar   string    `json:"avatar"`
	JoinedAt time.Time `json:"joinedAt"`
}

type UpdateUserRequest struct {
	Fullname  string                `form:"fullname" binding:"required"`
	Avatar    *multipart.FileHeader `form:"avatar" binding:"omitempty"`
	AvatarURL string                `form:"avatarUrl"`
}

type PaginationResponse struct {
	CurrentPage int `json:"currentPage"`
	TotalItems  int `json:"totalItems"`
	TotalPages  int `json:"totalPages"`
	Limit       int `json:"limit"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"currentPassword" binding:"required,min=6"`
	NewPassword     string `json:"newPassword" binding:"required,min=6"`
	ConfirmPassword string `json:"confirmPassword" binding:"required,min=6"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ResetPasswordRequest struct {
	Token           string `json:"token" binding:"required"`
	NewPassword     string `json:"newPassword" binding:"required,min=6"`
	ConfirmPassword string `json:"confirmPassword" binding:"required,min=6"`
}

type ForgotPasswordResponse struct {
	Message string `json:"message"`
	Email   string `json:"email"`
}

type ResetTokenData struct {
	UserID    string    `json:"userId"`
	Email     string    `json:"email"`
	ExpiresAt time.Time `json:"expiresAt"`
	CreatedAt time.Time `json:"createdAt"`
}

// category DTOs
type CategoryResponse struct {
	ID        string             `json:"id"`
	ParentID  *string            `json:"parentId"`
	Name      string             `json:"name"`
	IsDefault bool               `json:"isDefault"`
	IsCustom  bool               `json:"isCustom"`
	IsParent  bool               `json:"isParent"`
	Level     int                `json:"level"` // 0 for parent, 1 for child
	CreatedAt time.Time          `json:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt"`
	Children  []CategoryResponse `json:"children,omitempty"`
}

type CategoriesTreeResponse struct {
	Categories []CategoryResponse `json:"categories"`
	Total      int                `json:"total"`
	Parents    int                `json:"parents"`
	Children   int                `json:"children"`
}

type CategoryFlatResponse struct {
	ID        string    `json:"id"`
	ParentID  *string   `json:"parentId"`
	Name      string    `json:"name"`
	FullName  string    `json:"fullname"` // "Technology > Electronics"
	IsDefault bool      `json:"isDefault"`
	IsCustom  bool      `json:"isCustom"`
	IsParent  bool      `json:"isParent"`
	Level     int       `json:"level"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type CategoriesFlatResponse struct {
	Categories []CategoryFlatResponse `json:"categories"`
	Total      int                    `json:"total"`
}

type CategoryWithAssetsResponse struct {
	Category CategoryResponse        `json:"category"`
	Assets   []CategoryAssetResponse `json:"assets"`
	Total    int                     `json:"total"`
}

type CategoryAssetResponse struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Price        float64   `json:"price"`
	Condition    string    `json:"condition"`
	SerialNumber string    `json:"serialNumber"`
	CreatedAt    time.Time `json:"createdAt"`
}

// Request DTOs
type CreateCategoryRequest struct {
	Name     string  `json:"name" binding:"required,min=2,max=100"`
	ParentID *string `json:"parentId" binding:"omitempty"`
}

type UpdateCategoryRequest struct {
	Name     string  `json:"name" binding:"required,min=2,max=100"`
	ParentID *string `json:"parentId" binding:"omitempty"`
}

// location DTOs
type LocationResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	IsDefault bool      `json:"isDefault"`
	IsCustom  bool      `json:"isCustom"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type LocationsResponse struct {
	Locations []LocationResponse `json:"locations"`
	Total     int                `json:"total"`
}

type LocationAssetResponse struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Price        float64   `json:"price"`
	Condition    string    `json:"condition"`
	SerialNumber string    `json:"serialNumber"`
	CreatedAt    time.Time `json:"createdAt"`
}

type LocationWithAssetsResponse struct {
	Location LocationResponse        `json:"location"`
	Assets   []LocationAssetResponse `json:"assets"`
	Total    int                     `json:"total"`
}

type CreateLocationRequest struct {
	Name string `json:"name" binding:"required,min=2,max=100"`
}

type UpdateLocationRequest struct {
	Name string `json:"name" binding:"required,min=2,max=100"`
}

// asset DTOs
type CreateAssetRequest struct {
	Name         string                `form:"name" json:"name" binding:"required,min=1,max=100"`
	Description  string                `form:"description" json:"description" binding:"max=255"`
	LocationID   string                `form:"locationId" json:"locationId" binding:"required,uuid"`
	CategoryID   string                `form:"categoryId" json:"categoryId" binding:"required,uuid"`
	Image        *multipart.FileHeader `form:"image" json:"-"`
	PurchaseDate *time.Time            `form:"purchaseDate" json:"purchaseDate" time_format:"2006-01-02"`
	Price        float64               `form:"price" json:"price" binding:"required,min=0"`
	Condition    string                `form:"condition" json:"condition" binding:"required,oneof=new good fair poor"`
	SerialNumber string                `form:"serialNumber" json:"serialNumber" binding:"max=100"`
	Warranty     *time.Time            `form:"warranty" json:"warranty" time_format:"2006-01-02"`
	ImageURL     string                `json:"-"`
}

type UpdateAssetRequest struct {
	Name         string                `form:"name" json:"name" binding:"omitempty,min=1,max=100"`
	Description  string                `form:"description" json:"description" binding:"max=255"`
	LocationID   string                `form:"locationId" json:"locationId" binding:"omitempty,uuid"`
	CategoryID   string                `form:"categoryId" json:"categoryId" binding:"omitempty,uuid"`
	Image        *multipart.FileHeader `form:"image" json:"-"`
	PurchaseDate *time.Time            `form:"purchaseDate" json:"purchaseDate" time_format:"2006-01-02"`
	Price        *float64              `form:"price" json:"price" binding:"omitempty,min=0"`
	Condition    string                `form:"condition" json:"condition" binding:"omitempty,oneof=new good fair poor"`
	SerialNumber string                `form:"serialNumber" json:"serialNumber" binding:"max=100"`
	Warranty     *time.Time            `form:"warranty" json:"warranty" time_format:"2006-01-02"`
	ImageURL     string                `json:"-"`
}

type GetAssetsRequest struct {
	Page       int      `form:"page" json:"page" binding:"omitempty,min=1"`
	Limit      int      `form:"limit" json:"limit" binding:"omitempty,min=1,max=100"`
	Search     string   `form:"search" json:"search" binding:"omitempty,max=100"`
	CategoryID string   `form:"categoryId" json:"categoryId" binding:"omitempty,uuid"`
	LocationID string   `form:"locationId" json:"locationId" binding:"omitempty,uuid"`
	Condition  string   `form:"condition" json:"condition" binding:"omitempty,oneof=new good fair poor"`
	MinPrice   *float64 `form:"minPrice" json:"minPrice" binding:"omitempty,min=0"`
	MaxPrice   *float64 `form:"maxPrice" json:"maxPrice" binding:"omitempty,min=0"`
	SortBy     string   `form:"sortBy" json:"sortBy" binding:"omitempty,oneof=name price createdAt purchaseDate"`
	SortOrder  string   `form:"sortOrder" json:"sortOrder" binding:"omitempty,oneof=asc desc"`
}

// Response DTOs
type AssetResponse struct {
	ID           string            `json:"id"`
	Name         string            `json:"name"`
	Description  string            `json:"description"`
	LocationID   string            `json:"locationId"`
	CategoryID   string            `json:"categoryId"`
	UserID       string            `json:"userId"`
	Image        string            `json:"image"`
	PurchaseDate *time.Time        `json:"purchaseDate"`
	Price        float64           `json:"price"`
	Condition    string            `json:"condition"`
	SerialNumber string            `json:"serialNumber"`
	Warranty     *time.Time        `json:"warranty"`
	CreatedAt    time.Time         `json:"createdAt"`
	UpdatedAt    time.Time         `json:"updatedAt"`
	Location     *LocationResponse `json:"location,omitempty"`
	Category     *CategoryResponse `json:"category,omitempty"`
}
