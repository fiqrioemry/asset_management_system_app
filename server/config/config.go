package config

func InitConfiguration() {
	LoadEnv()
	InitRedis()
	InitMailer()
	InitDatabase()
	InitCloudinary()
}
