package seeders

import (
	"fmt"
	"log"

	"github.com/fiqrioemry/asset_management_system_app/server/models"

	"gorm.io/gorm"
)

func ResetDatabase(db *gorm.DB) {
	log.Println("Dropping all tables...")

	err := db.Migrator().DropTable(
		&models.User{},
		&models.Category{},
		&models.Asset{},
		&models.Location{},
	)
	if err != nil {
		log.Fatalf("Failed to drop tables: %v", err)
	}

	log.Println("all tables dropped successfully.")

	log.Println("migrating tables...")

	err = db.AutoMigrate(
		&models.User{},
		&models.Category{},
		&models.Asset{},
		&models.Location{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate tables: %v", err)
	}

	RunAllSeeders(db)

	fmt.Println("all tables migrated successfully.")
}
