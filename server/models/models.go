package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID      `json:"id" gorm:"type:varchar(36);primaryKey"`
	Fullname  string         `json:"fullname" gorm:"type:varchar(100);not null"`
	Avatar    string         `json:"avatar" gorm:"type:varchar(255)"`
	Email     string         `json:"email" gorm:"type:varchar(100);unique;not null"`
	Password  string         `json:"password" gorm:"type:varchar(100);not null"`
	CreatedAt time.Time      `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updatedAt" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"index"`

	Assets     []Asset    `json:"assets" gorm:"foreignKey:UserID"`
	Categories []Category `json:"categories" gorm:"foreignKey:UserID"`
	Locations  []Location `json:"locations" gorm:"foreignKey:UserID"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

// Location model
type Location struct {
	ID        uuid.UUID      `json:"id" gorm:"type:varchar(36);primaryKey"`
	Name      string         `json:"name" gorm:"type:varchar(100);not null"`
	UserID    *uuid.UUID     `json:"userId" gorm:"type:varchar(36);index"`
	IsDefault bool           `json:"isDefault" gorm:"default:false"`
	CreatedAt time.Time      `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updatedAt" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"index"`

	Assets []Asset `json:"assets" gorm:"foreignKey:LocationID"`
	User   *User   `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

func (l *Location) BeforeCreate(tx *gorm.DB) error {
	if l.ID == uuid.Nil {
		l.ID = uuid.New()
	}
	return nil
}

type Category struct {
	ID        uuid.UUID      `json:"id" gorm:"type:varchar(36);primaryKey"`
	ParentID  *uuid.UUID     `json:"parentId" gorm:"type:varchar(36);index"`
	Name      string         `json:"name" gorm:"type:varchar(100);not null"`
	UserID    *uuid.UUID     `json:"userId" gorm:"type:varchar(36);index"`
	IsDefault bool           `json:"isDefault" gorm:"default:false"`
	CreatedAt time.Time      `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updatedAt" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"index"`

	// Relationships
	Assets   []Asset    `json:"assets" gorm:"foreignKey:CategoryID"`
	User     *User      `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Parent   *Category  `json:"parent,omitempty" gorm:"foreignKey:ParentID"`
	Children []Category `json:"children,omitempty" gorm:"foreignKey:ParentID"`
}

func (c *Category) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

// Helper methods
func (c *Category) IsParent() bool {
	return c.ParentID == nil
}

func (c *Category) IsChild() bool {
	return c.ParentID != nil
}

func (c *Category) IsSystemCategory() bool {
	return c.UserID == nil
}

func (c *Category) IsUserCategory() bool {
	return c.UserID != nil
}

type Asset struct {
	ID           uuid.UUID      `json:"id" gorm:"type:varchar(36);primaryKey"`
	Name         string         `json:"name" gorm:"type:varchar(100);not null"`
	Description  string         `json:"description" gorm:"type:varchar(255)"`
	LocationID   uuid.UUID      `json:"locationId" gorm:"type:varchar(36);not null"`
	CategoryID   uuid.UUID      `json:"categoryId" gorm:"type:varchar(36);not null"`
	UserID       uuid.UUID      `json:"userId" gorm:"type:varchar(36);not null"`
	Image        string         `json:"image" gorm:"type:varchar(255)"`
	PurchaseDate *time.Time     `json:"purchaseDate" gorm:"type:date"`
	Price        float64        `json:"price" gorm:"type:decimal(10,2);not null"`
	Condition    string         `json:"condition" gorm:"type:varchar(50);not null"`
	SerialNumber string         `json:"serialNumber" gorm:"type:varchar(100)"`
	Warranty     *time.Time     `json:"warranty" gorm:"type:date"`
	CreatedAt    time.Time      `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt    time.Time      `json:"updatedAt" gorm:"autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `json:"deletedAt" gorm:"index"`

	Location Location `json:"location" gorm:"foreignKey:LocationID"`
	Category Category `json:"category" gorm:"foreignKey:CategoryID"`
	User     User     `json:"user" gorm:"foreignKey:UserID"`
}

func (a *Asset) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}
