package config

func InitConfiguration() {
	LoadConfig()
	InitRedis()
	InitMailer()
	InitDatabase()
	InitCloudinary()
	InitGoogleOAuthConfig()
}
