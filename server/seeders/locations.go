package seeders

import (
	"github.com/fiqrioemry/asset_management_system_app/server/models"
	"gorm.io/gorm"
)

func SeedSystemLocations(db *gorm.DB) error {
	systemLocations := []models.Location{
		{Name: "Living Room", IsDefault: true},
		{Name: "Bedroom", IsDefault: true},
		{Name: "Kitchen", IsDefault: true},
		{Name: "Bathroom", IsDefault: true},
		{Name: "Garage", IsDefault: true},
		{Name: "Office", IsDefault: true},
		{Name: "Storage", IsDefault: true},
		{Name: "Basement", IsDefault: true},
		{Name: "Attic", IsDefault: true},
		{Name: "Garden", IsDefault: true},
	}

	for _, location := range systemLocations {
		var existing models.Location
		if err := db.Where("name = ? AND user_id IS NULL", location.Name).First(&existing).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&location).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		}
	}

	return nil
}
