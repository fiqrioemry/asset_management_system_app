package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config holds all environment configuration
type Config struct {
	// Server settings
	ServerPort string
	ServerHost string

	// Security settings
	ApiKeys             string
	AllowedOrigins      []string
	RateLimitAttempts   int
	RateLimitDuration   time.Duration
	MaxFileSize         int64
	SkippedApiEndpoints []string
	TrustedProxies      []string
	CookieDomain        string

	// Database settings
	DatabaseRootURL string
	DatabaseName    string
	DatabaseURL     string

	// Redis settings
	RedisAddress  string
	RedisPassword string

	// JWT settings
	AccessTokenSecret  string
	RefreshTokenSecret string

	// Email settings
	SMTPHost     string
	SMTPPort     int
	SMTPEmail    string
	SMTPPassword string

	// App settings
	AppName     string
	AppEnv      string
	FrontendURL string

	// cloudinary settings
	CloudName   string
	CloudSecret string
	CloudApiKey string
	CloudFolder string

	// google oauth settings
	GoogleClientID      string
	GoogleClientSecret  string
	GoogleRedirectURL   string
	FrontendRedirectURL string
}

var AppConfig *Config

func LoadConfig() {
	AppConfig = &Config{
		// Server
		ServerPort: getEnvOrDefault("PORT", "8080"),
		ServerHost: getEnvOrDefault("HOST", "localhost"),

		// google oauth
		GoogleClientID:      getEnvOrDefault("GOOGLE_CLIENT_ID", "your-google-client-id"),
		GoogleClientSecret:  getEnvOrDefault("GOOGLE_CLIENT_SECRET", "your-google-client-secret"),
		GoogleRedirectURL:   getEnvOrDefault("GOOGLE_REDIRECT_URL", "http://localhost:5005/api/v1/users/google/callback"),
		FrontendRedirectURL: getEnvOrDefault("FRONTEND_REDIRECT_URL", "http://localhost:5173"),

		// Security
		CookieDomain:        getEnvOrDefault("COOKIE_DOMAIN", "localhost"),
		ApiKeys:             getEnvOrDefault("API_KEY", "your-api-keys"),
		RateLimitAttempts:   getEnvAsInt("RATE_LIMIT_ATTEMPTS", 100),
		RateLimitDuration:   getEnvAsDuration("RATE_LIMIT_DURATION", "60s"),
		MaxFileSize:         getEnvAsInt64("MAX_FILE_SIZE", 12<<20),
		TrustedProxies:      getEnvAsStringSlice("TRUSTED_PROXIES", []string{"localhost"}),
		SkippedApiEndpoints: getEnvAsStringSlice("SKIPPED_API_ENDPOINTS", []string{"/health", "/"}),
		AllowedOrigins:      getEnvAsStringSlice("ALLOWED_ORIGINS", []string{"http://localhost:3000"}),

		// Database
		DatabaseRootURL: getEnvOrDefault("DB_ROOT_URL", "your-db-root-url"),
		DatabaseName:    getEnvOrDefault("DB_NAME", "your-db-name"),
		DatabaseURL:     getEnvOrDefault("DB_URL", "your-db-url"),

		// Redis
		RedisAddress:  getEnvOrDefault("REDIS_ADDRESS", "localhost:6379"),
		RedisPassword: getEnvOrDefault("REDIS_PASSWORD", ""),

		// JWT
		AccessTokenSecret:  getEnvOrDefault("ACCESS_TOKEN_SECRET", "your-secret-key"),
		RefreshTokenSecret: getEnvOrDefault("REFRESH_TOKEN_SECRET", "your-refresh-token-secret"),

		// mailer configuration
		SMTPEmail:    getEnvOrDefault("SMTP_EMAIL", ""),
		SMTPPort:     getEnvAsInt("SMTP_PORT", 587),
		SMTPHost:     getEnvOrDefault("SMTP_HOST", ""),
		SMTPPassword: getEnvOrDefault("SMTP_PASSWORD", ""),

		// App
		AppName:     getEnvOrDefault("APP_NAME", "Asset Management System"),
		AppEnv:      getEnvOrDefault("APP_ENV", "development"),
		FrontendURL: getEnvOrDefault("FRONTEND_URL", "http://localhost:5173"),

		// Cloudinary
		CloudName:   getEnvOrDefault("CLOUDINARY_CLOUD_NAME", "your-cloudinary-cloud-name"),
		CloudSecret: getEnvOrDefault("CLOUDINARY_API_SECRET", "your-cloudinary-api-secret"),
		CloudApiKey: getEnvOrDefault("CLOUDINARY_API_KEY", "your-cloudinary-api-key"),
		CloudFolder: getEnvOrDefault("CLOUDINARY_FOLDER", "asset_management_app"),
	}

	fmt.Println("✅ Global configuration load complete")
}

// Helper functions for parsing environment variables

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}

func getEnvAsInt64(key string, defaultValue int64) int64 {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.ParseInt(value, 10, 64); err == nil {
			return parsed
		}
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue string) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		value = defaultValue
	}

	if duration, err := time.ParseDuration(value); err == nil {
		return duration
	}

	duration, _ := time.ParseDuration(defaultValue)
	return duration
}

func getEnvAsStringSlice(key string, defaultValue []string) []string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	var result []string
	for _, item := range strings.Split(value, ",") {
		trimmed := strings.TrimSpace(item)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func GetServerAddress() string {
	return AppConfig.ServerHost + ":" + AppConfig.ServerPort
}

func IsProduction() bool {
	return AppConfig.AppEnv == "production"
}

func IsDevelopment() bool {
	return AppConfig.AppEnv == "development"
}
