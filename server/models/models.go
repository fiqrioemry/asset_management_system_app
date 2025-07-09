package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User model
type User struct {
	ID        uuid.UUID      `json:"id" gorm:"type:varchar(36);primaryKey"`
	Fullname  string         `json:"fullname" gorm:"type:varchar(100);not null"`
	Email     string         `json:"email" gorm:"type:varchar(100);unique;not null"`
	Password  string         `json:"password" gorm:"type:varchar(100);not null"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`

	// Relationships
	Assets []Asset `json:"assets" gorm:"foreignKey:UserID"`
}

// Location model
type Location struct {
	ID        uuid.UUID      `json:"id" gorm:"type:varchar(36);primaryKey"`
	Name      string         `json:"name" gorm:"type:varchar(100);not null"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`

	// Relationships
	Assets []Asset `json:"assets" gorm:"foreignKey:LocationID"`
}

// Category model
type Category struct {
	ID        uuid.UUID      `json:"id" gorm:"type:varchar(36);primaryKey"`
	Name      string         `json:"name" gorm:"type:varchar(100);not null"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`

	// Relationships
	Assets []Asset `json:"assets" gorm:"foreignKey:CategoryID"`
}

// Asset model
type Asset struct {
	ID           uuid.UUID      `json:"id" gorm:"type:varchar(36);primaryKey"`
	Name         string         `json:"name" gorm:"type:varchar(100);not null"`
	Description  string         `json:"description" gorm:"type:varchar(255)"`
	LocationID   uuid.UUID      `json:"location_id" gorm:"type:varchar(36);not null"`
	CategoryID   uuid.UUID      `json:"category_id" gorm:"type:varchar(36);not null"`
	UserID       uuid.UUID      `json:"user_id" gorm:"type:varchar(36);not null"`
	Image        string         `json:"image" gorm:"type:varchar(255)"`
	PurchaseDate *time.Time     `json:"purchase_date" gorm:"type:date"`
	Price        float64        `json:"price" gorm:"type:decimal(10,2);not null"`
	Condition    string         `json:"condition" gorm:"type:varchar(50);not null"`
	SerialNumber string         `json:"serial_number" gorm:"type:varchar(100)"`
	Warranty     *time.Time     `json:"warranty" gorm:"type:date"`
	CreatedAt    time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at" gorm:"index"`

	// Relationships
	Location Location `json:"location" gorm:"foreignKey:LocationID"`
	Category Category `json:"category" gorm:"foreignKey:CategoryID"`
	User     User     `json:"user" gorm:"foreignKey:UserID"`
}

// BeforeCreate hook untuk generate UUID
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

func (l *Location) BeforeCreate(tx *gorm.DB) error {
	if l.ID == uuid.Nil {
		l.ID = uuid.New()
	}
	return nil
}

func (c *Category) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

func (a *Asset) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}
