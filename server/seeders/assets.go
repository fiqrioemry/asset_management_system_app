package seeders

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/fiqrioemry/asset_management_system_app/server/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func SeedAssets(db *gorm.DB) error {
	// Get all users
	var users []models.User
	if err := db.Find(&users).Error; err != nil {
		return err
	}

	// Get all categories (only child categories for realistic data)
	var categories []models.Category
	if err := db.Where("parent_id IS NOT NULL").Find(&categories).Error; err != nil {
		return err
	}

	// Get all locations
	var locations []models.Location
	if err := db.Find(&locations).Error; err != nil {
		return err
	}

	if len(categories) == 0 || len(locations) == 0 {
		return fmt.Errorf("categories or locations not found, please run category and location seeders first")
	}

	// Asset templates with realistic data
	assetTemplates := []AssetTemplate{
		// Technology Assets
		{Name: "MacBook Pro 16-inch", Description: "High-performance laptop for development work", Price: 2499.99, Condition: "good", SerialNumber: "MBP2023001", CategoryName: "Laptops"},
		{Name: "iPhone 14 Pro", Description: "Latest smartphone with advanced camera system", Price: 999.99, Condition: "new", SerialNumber: "IP14P001", CategoryName: "Mobile Devices"},
		{Name: "iPad Air", Description: "Versatile tablet for work and entertainment", Price: 599.99, Condition: "good", SerialNumber: "IPA2023001", CategoryName: "Tablets"},
		{Name: "AirPods Pro", Description: "Wireless earbuds with noise cancellation", Price: 249.99, Condition: "good", SerialNumber: "APP2023001", CategoryName: "Audio & Video"},
		{Name: "Gaming Desktop PC", Description: "High-end gaming computer with RTX 4080", Price: 3499.99, Condition: "new", SerialNumber: "GPC2023001", CategoryName: "Computers"},

		// Home & Living Assets
		{Name: "Samsung 55\" Smart TV", Description: "4K QLED smart television", Price: 899.99, Condition: "good", SerialNumber: "SAM55Q001", CategoryName: "Electronics"},
		{Name: "Dyson V15 Vacuum", Description: "Cordless vacuum cleaner with laser detection", Price: 449.99, Condition: "good", SerialNumber: "DYS15001", CategoryName: "Appliances"},
		{Name: "Herman Miller Chair", Description: "Ergonomic office chair", Price: 1200.00, Condition: "good", SerialNumber: "HM2023001", CategoryName: "Furniture"},
		{Name: "KitchenAid Mixer", Description: "Stand mixer for baking and cooking", Price: 379.99, Condition: "good", SerialNumber: "KA2023001", CategoryName: "Kitchen"},
		{Name: "Instant Pot Duo", Description: "Multi-use pressure cooker", Price: 99.99, Condition: "good", SerialNumber: "IP2023001", CategoryName: "Kitchen"},

		// Transportation Assets
		{Name: "Toyota Camry 2022", Description: "Reliable sedan with hybrid engine", Price: 28999.99, Condition: "good", SerialNumber: "TC2022001", CategoryName: "Cars"},
		{Name: "Trek Mountain Bike", Description: "All-terrain mountain bicycle", Price: 899.99, Condition: "good", SerialNumber: "TRK2023001", CategoryName: "Bicycles"},
		{Name: "Honda Civic 2023", Description: "Compact car with excellent fuel economy", Price: 24999.99, Condition: "new", SerialNumber: "HC2023001", CategoryName: "Cars"},

		// Tools & Equipment
		{Name: "DeWalt Drill Set", Description: "Cordless drill with battery pack", Price: 199.99, Condition: "good", SerialNumber: "DW2023001", CategoryName: "Power Tools"},
		{Name: "Craftsman Tool Box", Description: "Large toolbox with multiple compartments", Price: 299.99, Condition: "good", SerialNumber: "CT2023001", CategoryName: "Hand Tools"},

		// Entertainment
		{Name: "Canon EOS R5", Description: "Professional mirrorless camera", Price: 3899.99, Condition: "good", SerialNumber: "CR5001", CategoryName: "Photography"},
		{Name: "Yamaha Piano", Description: "Digital piano with weighted keys", Price: 1299.99, Condition: "good", SerialNumber: "YP2023001", CategoryName: "Musical Instruments"},
		{Name: "Nintendo Switch OLED", Description: "Gaming console with OLED screen", Price: 349.99, Condition: "new", SerialNumber: "NSW2023001", CategoryName: "Gaming"},

		// Office & Business
		{Name: "HP LaserJet Printer", Description: "Wireless laser printer for office use", Price: 299.99, Condition: "good", SerialNumber: "HP2023001", CategoryName: "Office Equipment"},
		{Name: "Dell Monitor 27-inch", Description: "4K monitor for productivity", Price: 399.99, Condition: "good", SerialNumber: "DM27001", CategoryName: "Computers & IT"},

		// Personal Items
		{Name: "Rolex Submariner", Description: "Luxury diving watch", Price: 8999.99, Condition: "good", SerialNumber: "RLX2023001", CategoryName: "Watches"},
		{Name: "Diamond Ring", Description: "Engagement ring with 1 carat diamond", Price: 5999.99, Condition: "new", SerialNumber: "DR2023001", CategoryName: "Jewelry"},

		// Health & Beauty
		{Name: "Peloton Bike", Description: "Indoor exercise bike with screen", Price: 1899.99, Condition: "good", SerialNumber: "PEL2023001", CategoryName: "Fitness Equipment"},
		{Name: "Air Purifier", Description: "HEPA air purifier for home use", Price: 299.99, Condition: "good", SerialNumber: "AP2023001", CategoryName: "Wellness"},

		// Additional varied assets
		{Name: "Sony PlayStation 5", Description: "Latest gaming console", Price: 499.99, Condition: "new", SerialNumber: "PS5001", CategoryName: "Gaming"},
		{Name: "Espresso Machine", Description: "Semi-automatic espresso maker", Price: 699.99, Condition: "good", SerialNumber: "EM2023001", CategoryName: "Kitchen"},
		{Name: "Smart Watch", Description: "Fitness tracker with heart rate monitor", Price: 299.99, Condition: "good", SerialNumber: "SW2023001", CategoryName: "Electronics"},
		{Name: "Leather Sofa", Description: "3-seater leather sofa", Price: 1599.99, Condition: "good", SerialNumber: "LS2023001", CategoryName: "Furniture"},
		{Name: "Projector", Description: "4K home theater projector", Price: 1299.99, Condition: "good", SerialNumber: "PJ2023001", CategoryName: "Audio & Video"},
		{Name: "Electric Scooter", Description: "Foldable electric scooter", Price: 599.99, Condition: "good", SerialNumber: "ES2023001", CategoryName: "Transportation"},
	}

	// Create 10 assets for each user
	for _, user := range users {
		// Shuffle asset templates to get random selection
		shuffledAssets := make([]AssetTemplate, len(assetTemplates))
		copy(shuffledAssets, assetTemplates)

		// Simple shuffle
		for i := range shuffledAssets {
			j := rand.Intn(i + 1)
			shuffledAssets[i], shuffledAssets[j] = shuffledAssets[j], shuffledAssets[i]
		}

		// Create 10 assets for this user
		for i := 0; i < 10; i++ {
			template := shuffledAssets[i%len(shuffledAssets)]

			// Find category by name
			var category models.Category
			found := false
			for _, cat := range categories {
				if cat.Name == template.CategoryName {
					category = cat
					found = true
					break
				}
			}
			if !found {
				category = categories[rand.Intn(len(categories))] // fallback to random category
			}

			// Random location
			location := locations[rand.Intn(len(locations))]

			// Generate random dates
			purchaseDate := generateRandomPastDate(365)   // within last year
			warrantyDate := generateRandomFutureDate(730) // up to 2 years from now

			// Create unique asset name for each user
			assetName := fmt.Sprintf("%s (%s)", template.Name, user.Fullname)

			asset := models.Asset{
				ID:           uuid.New(),
				Name:         assetName,
				Description:  template.Description,
				LocationID:   location.ID,
				CategoryID:   category.ID,
				UserID:       user.ID,
				Image:        generateAssetImageURL(template.Name),
				PurchaseDate: &purchaseDate,
				Price:        addPriceVariation(template.Price),
				Condition:    template.Condition,
				SerialNumber: fmt.Sprintf("%s-%s", template.SerialNumber, user.ID.String()[:8]),
				Warranty:     &warrantyDate,
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			}

			// Check if asset already exists
			var existing models.Asset
			if err := db.Where("name = ? AND user_id = ?", asset.Name, asset.UserID).First(&existing).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					if err := db.Create(&asset).Error; err != nil {
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

// Helper struct for asset templates
type AssetTemplate struct {
	Name         string
	Description  string
	Price        float64
	Condition    string
	SerialNumber string
	CategoryName string
}

// Helper functions
func generateRandomPastDate(maxDaysAgo int) time.Time {
	daysAgo := rand.Intn(maxDaysAgo)
	return time.Now().AddDate(0, 0, -daysAgo)
}

func generateRandomFutureDate(maxDaysFromNow int) time.Time {
	daysFromNow := rand.Intn(maxDaysFromNow)
	return time.Now().AddDate(0, 0, daysFromNow)
}

func addPriceVariation(basePrice float64) float64 {
	// Add ±10% price variation
	variation := (rand.Float64() - 0.5) * 0.2 // -0.1 to +0.1
	return basePrice * (1 + variation)
}

func generateAssetImageURL(assetName string) string {
	// Generate placeholder image URLs based on asset type
	imageMap := map[string]string{
		"MacBook":     "https://images.unsplash.com/photo-1517336714731-489689fd1ca8?w=400&h=300&fit=crop",
		"iPhone":      "https://images.unsplash.com/photo-1592750475338-74b7b21085ab?w=400&h=300&fit=crop",
		"iPad":        "https://images.unsplash.com/photo-1544244015-0df4b3ffc6b0?w=400&h=300&fit=crop",
		"AirPods":     "https://images.unsplash.com/photo-1606220588913-b3aacb4d2f46?w=400&h=300&fit=crop",
		"Gaming":      "https://images.unsplash.com/photo-1587202372634-32705e3bf49c?w=400&h=300&fit=crop",
		"Samsung":     "https://images.unsplash.com/photo-1593359677879-a4bb92f829d1?w=400&h=300&fit=crop",
		"Dyson":       "https://images.unsplash.com/photo-1558618666-fbd1c1d6c7e2?w=400&h=300&fit=crop",
		"Chair":       "https://images.unsplash.com/photo-1506439773649-6e0eb8cfb237?w=400&h=300&fit=crop",
		"Mixer":       "https://images.unsplash.com/photo-1556909114-4df7eca6531a?w=400&h=300&fit=crop",
		"Toyota":      "https://images.unsplash.com/photo-1549924231-f129b911e442?w=400&h=300&fit=crop",
		"Honda":       "https://images.unsplash.com/photo-1552519507-da3b142c6e3d?w=400&h=300&fit=crop",
		"Bike":        "https://images.unsplash.com/photo-1558618666-fcd25c85cd64?w=400&h=300&fit=crop",
		"Drill":       "https://images.unsplash.com/photo-1504148455328-a87fb251851e?w=400&h=300&fit=crop",
		"Camera":      "https://images.unsplash.com/photo-1502920917128-1aa500764cbd?w=400&h=300&fit=crop",
		"Piano":       "https://images.unsplash.com/photo-1493225457124-a3eb161ffa5f?w=400&h=300&fit=crop",
		"Nintendo":    "https://images.unsplash.com/photo-1606144042614-b2417e99c4e3?w=400&h=300&fit=crop",
		"Printer":     "https://images.unsplash.com/photo-1612198188060-c7c2a3b66eae?w=400&h=300&fit=crop",
		"Monitor":     "https://images.unsplash.com/photo-1527443224154-c4a3942d3acf?w=400&h=300&fit=crop",
		"Watch":       "https://images.unsplash.com/photo-1524592094714-0f0654e20314?w=400&h=300&fit=crop",
		"Ring":        "https://images.unsplash.com/photo-1605100804763-247f67b3557e?w=400&h=300&fit=crop",
		"Peloton":     "https://images.unsplash.com/photo-1571019613454-1cb2f99b2d8b?w=400&h=300&fit=crop",
		"PlayStation": "https://images.unsplash.com/photo-1606144042614-b2417e99c4e3?w=400&h=300&fit=crop",
		"Sofa":        "https://images.unsplash.com/photo-1555041469-a586c61ea9bc?w=400&h=300&fit=crop",
	}

	// Find matching image based on asset name
	for key, url := range imageMap {
		if containsIgnoreCase(assetName, key) {
			return url
		}
	}

	// Default image if no match found
	return "https://images.unsplash.com/photo-1560472354-b33ff0c44a43?w=400&h=300&fit=crop"
}

func containsIgnoreCase(str, substr string) bool {
	return len(str) >= len(substr) &&
		(str[:len(substr)] == substr ||
			strings.ToLower(str[:len(substr)]) == strings.ToLower(substr))
}
