package seeders

import (
	"github.com/fiqrioemry/asset_management_system_app/server/models"
	"gorm.io/gorm"
)

func SeedSystemCategories(db *gorm.DB) error {
	// Step 1: Create parent categories
	parentCategories := []models.Category{
		{Name: "Technology", IsDefault: true},
		{Name: "Home & Living", IsDefault: true},
		{Name: "Transportation", IsDefault: true},
		{Name: "Tools & Equipment", IsDefault: true},
		{Name: "Personal Items", IsDefault: true},
		{Name: "Entertainment", IsDefault: true},
		{Name: "Office & Business", IsDefault: true},
		{Name: "Health & Beauty", IsDefault: true},
		{Name: "Miscellaneous", IsDefault: true},
	}

	// Create parent categories first
	for _, category := range parentCategories {
		var existing models.Category
		if err := db.Where("name = ? AND user_id IS NULL AND parent_id IS NULL", category.Name).First(&existing).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&category).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		}
	}

	// Step 2: Create child categories
	childCategoriesData := map[string][]string{
		"Technology": {
			"Computers", "Laptops", "Mobile Devices", "Tablets",
			"Audio & Video", "Gaming", "Smart Home", "Electronics",
		},
		"Home & Living": {
			"Furniture", "Appliances", "Kitchen", "Bathroom",
			"Bedroom", "Living Room", "Home Decor", "Storage",
		},
		"Transportation": {
			"Vehicles", "Cars", "Motorcycles", "Bicycles", "Auto Parts",
		},
		"Tools & Equipment": {
			"Hand Tools", "Power Tools", "Garden Tools", "Construction", "Safety Equipment",
		},
		"Personal Items": {
			"Jewelry", "Watches", "Clothing", "Accessories", "Personal Care",
		},
		"Entertainment": {
			"Books", "Movies & Music", "Sports Equipment", "Musical Instruments",
			"Art & Crafts", "Photography", "Games", "Fitness Equipment",
		},
		"Office & Business": {
			"Office Equipment", "Computers & IT", "Stationery", "Documents",
		},
		"Health & Beauty": {
			"Medical Equipment", "Fitness", "Beauty Products", "Wellness",
		},
		"Miscellaneous": {
			"Collectibles", "Emergency Supplies", "Seasonal Items", "Others",
		},
	}

	// Create child categories
	for parentName, childNames := range childCategoriesData {
		// Get parent category
		var parent models.Category
		if err := db.Where("name = ? AND user_id IS NULL AND parent_id IS NULL", parentName).First(&parent).Error; err != nil {
			return err
		}

		// Create child categories
		for _, childName := range childNames {
			var existingChild models.Category
			if err := db.Where("name = ? AND parent_id = ? AND user_id IS NULL", childName, parent.ID).First(&existingChild).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					childCategory := models.Category{
						ParentID:  &parent.ID,
						Name:      childName,
						IsDefault: true,
					}
					if err := db.Create(&childCategory).Error; err != nil {
						return err
					}
				} else {
					return err
				}
			}
		}
	}

	return nil
}
