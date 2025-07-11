// seeders/seeder.go (Main seeder file)
package seeders

import (
	"math/rand"
	"time"

	"gorm.io/gorm"
)

func RunAllSeeders(db *gorm.DB) error {
	rand.Seed(time.Now().UnixNano())

	if err := SeedUsers(db); err != nil {
		return err
	}
	if err := SeedSystemCategories(db); err != nil {
		return err
	}
	if err := SeedSystemLocations(db); err != nil {
		return err
	}

	if err := SeedAssets(db); err != nil {
		return err
	}
	return nil
}
