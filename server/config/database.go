package config

import (
	"fmt"
	"os"
	"time"

	"github.com/fiqrioemry/asset_management_system_app/server/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDatabase() {
	dbRootURL := os.Getenv("DB_ROOT_URL")
	dbURL := os.Getenv("DB_URL")
	dbName := os.Getenv("DB_NAME")

	// Connect to MySQL server
	dbRoot, err := gorm.Open(mysql.Open(dbRootURL), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to MySQL server: " + err.Error())
	}

	// Create DB if not exists
	sql := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", dbName)
	if err := dbRoot.Exec(sql).Error; err != nil {
		panic("Failed to create database: " + err.Error())
	}

	// Connect ke database utama
	for range 10 {
		DB, err = gorm.Open(mysql.Open(dbURL), &gorm.Config{})
		if err == nil {
			break
		}
		fmt.Println("Waiting for database to be ready...")
		time.Sleep(3 * time.Second)
	}
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}

	// Migration
	if err := DB.AutoMigrate(
		&models.User{},
		&models.Asset{},
		&models.Location{},
		&models.Category{},
	); err != nil {
		panic("Migration failed: " + err.Error())
	}

	sqlDB, err := DB.DB()
	if err != nil {
		panic("Failed to get database connection: " + err.Error())
	}
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(time.Hour)

	fmt.Println("Database connection established successfully.")
}
