package config

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client
var Ctx = context.Background()

func InitRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     AppConfig.RedisAddress,
		Password: AppConfig.RedisPassword,
		DB:       0,
	})

	fmt.Println("✅ Redis is connected ")
}
