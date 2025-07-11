// seeders/user_seeder.go
package seeders

import (
	"time"

	"github.com/fiqrioemry/asset_management_system_app/server/models"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SeedUsers(db *gorm.DB) error {
	// Hash password for all users
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	users := []models.User{
		{
			ID:       uuid.New(),
			Email:    "john.doe@example.com",
			Password: string(hashedPassword),
			Fullname: "John Doe",
			Avatar:   "https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=150&h=150&fit=crop&crop=face",
		},
		{
			ID:       uuid.New(),
			Email:    "jane.smith@example.com",
			Password: string(hashedPassword),
			Fullname: "Jane Smith",
			Avatar:   "https://images.unsplash.com/photo-1494790108755-2616b612b786?w=150&h=150&fit=crop&crop=face",
		},
		{
			ID:       uuid.New(),
			Email:    "mike.johnson@example.com",
			Password: string(hashedPassword),
			Fullname: "Mike Johnson",
			Avatar:   "https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?w=150&h=150&fit=crop&crop=face",
		},
	}

	for _, user := range users {
		var existing models.User
		if err := db.Where("email = ?", user.Email).First(&existing).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				user.CreatedAt = time.Now()
				user.UpdatedAt = time.Now()
				if err := db.Create(&user).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		}
	}

	return nil
}
