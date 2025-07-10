package config

import (
	"fmt"
	"log"

	"github.com/cloudinary/cloudinary-go/v2"
)

var Cloud *cloudinary.Cloudinary

func InitCloudinary() {
	cloud, err := cloudinary.NewFromParams(
		AppConfig.CloudName,
		AppConfig.CloudApiKey,
		AppConfig.CloudSecret,
	)
	if err != nil {
		log.Fatalf("Failed to initialize Cloudinary: %v", err)
	}

	Cloud = cloud
	fmt.Println("✅ Cloudinary configured")
}
