package dto

import (
	"time"

	"github.com/google/uuid"
)

type UserResponse struct {
	ID        uuid.UUID `json:"id" binding:"required"`
	Fullname  string    `json:"fullname" binding:"required"`
	Email     string    `json:"email" binding:"required,email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type AssetResponse struct {
	ID           uuid.UUID        `json:"id "`
	Name         string           `json:"name"`
	Description  string           `json:"description"`
	LocationID   uuid.UUID        `json:"location_id"`
	CategoryID   uuid.UUID        `json:"category_id"`
	UserID       uuid.UUID        `json:"user_id"`
	Image        string           `json:"image"`
	PurchaseDate *time.Time       `json:"purchase_date"`
	Price        float64          `json:"price"`
	Condition    string           `json:"condition"`
	SerialNumber string           `json:"serial_number"`
	Warranty     *time.Time       `json:"warranty"`
	CreatedAt    time.Time        `json:"created_at"`
	UpdatedAt    time.Time        `json:"updated_at"`
	Location     LocationResponse `json:"location"`
	Category     CategoryResponse `json:"category"`
	User         UserResponse     `json:"user"`
}

type LocationResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CategoryResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Request structs untuk API
type CreateUserRequest struct {
	Fullname string `json:"fullname" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type CreateAssetRequest struct {
	Name         string     `json:"name" binding:"required"`
	Description  string     `json:"description"`
	LocationID   uuid.UUID  `json:"location_id" binding:"required"`
	CategoryID   uuid.UUID  `json:"category_id" binding:"required"`
	Image        string     `json:"image"`
	PurchaseDate *time.Time `json:"purchase_date"`
	Price        float64    `json:"price" binding:"required,min=0"`
	Condition    string     `json:"condition" binding:"required"`
	SerialNumber string     `json:"serial_number"`
	Warranty     *time.Time `json:"warranty"`
}

type UpdateAssetRequest struct {
	Name         string     `json:"name"`
	Description  string     `json:"description"`
	LocationID   uuid.UUID  `json:"location_id"`
	CategoryID   uuid.UUID  `json:"category_id"`
	Image        string     `json:"image"`
	PurchaseDate *time.Time `json:"purchase_date"`
	Price        float64    `json:"price" binding:"min=0"`
	Condition    string     `json:"condition"`
	SerialNumber string     `json:"serial_number"`
	Warranty     *time.Time `json:"warranty"`
}

type CreateLocationRequest struct {
	Name string `json:"name" binding:"required"`
}

type CreateCategoryRequest struct {
	Name string `json:"name" binding:"required"`
}
