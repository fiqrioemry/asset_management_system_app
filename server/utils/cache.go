package utils

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/fiqrioemry/asset_management_system_app/server/config"
)

func CheckAttempts(redisKey string, maxAttempts int) error {
	attempts, err := config.RedisClient.Get(config.Ctx, redisKey).Int()
	if err != nil && err.Error() != "redis: nil" {
		return fmt.Errorf("failed to get attempts: %w", err)
	}

	if attempts >= maxAttempts {
		return fmt.Errorf("too many attempts, please try again later")
	}

	return nil
}

func CheckForgotPasswordAttempts(clientIP string, maxAttempts int) error {
	redisKey := "asset_app:forgot_password_attempts:" + clientIP
	return CheckAttempts(redisKey, maxAttempts)
}

func IncrementAttempts(redisKey string) {
	config.RedisClient.Incr(config.Ctx, redisKey)
	config.RedisClient.Expire(config.Ctx, redisKey, 30*time.Minute)
}

func AddKeys(redisKey string, data any, duration time.Duration) error {
	var value string

	switch v := data.(type) {
	case string:
		value = v
	case []byte:
		value = string(v)
	case int, int32, int64, float32, float64, bool:
		value = fmt.Sprintf("%v", v)
	default:
		jsonData, err := json.Marshal(data)
		if err != nil {
			return fmt.Errorf("failed to marshal data: %w", err)
		}
		value = string(jsonData)
	}

	err := config.RedisClient.Set(config.Ctx, redisKey, value, duration).Err()
	if err != nil {
		return fmt.Errorf("failed to set key: %w", err)
	}

	return nil
}

func DeleteKeys(redisKeys ...string) error {
	if len(redisKeys) == 0 {
		return nil
	}

	err := config.RedisClient.Del(config.Ctx, redisKeys...).Err()
	if err != nil {
		return fmt.Errorf("failed to delete keys: %w", err)
	}

	return nil
}

func GetKey(redisKey string, dest any) error {
	result, err := config.RedisClient.Get(config.Ctx, redisKey).Result()
	if err != nil {
		if err.Error() == "redis: nil" {
			return fmt.Errorf("key not found")
		}
		return fmt.Errorf("failed to get key: %w", err)
	}

	switch d := dest.(type) {
	case *string:
		*d = result
		return nil
	case *[]byte:
		*d = []byte(result)
		return nil
	case *int:
		return config.RedisClient.Get(config.Ctx, redisKey).Scan(d)
	default:
		return json.Unmarshal([]byte(result), dest)
	}
}

func KeyExists(redisKey string) bool {
	result, err := config.RedisClient.Exists(config.Ctx, redisKey).Result()
	if err != nil {
		return false
	}
	return result > 0
}

func SetKeyExpiry(redisKey string, duration time.Duration) error {
	err := config.RedisClient.Expire(config.Ctx, redisKey, duration).Err()
	if err != nil {
		return fmt.Errorf("failed to set expiry: %w", err)
	}
	return nil
}

func GetKeysByPattern(pattern string) ([]string, error) {
	keys, err := config.RedisClient.Keys(config.Ctx, pattern).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get keys by pattern: %w", err)
	}
	return keys, nil
}

func DeleteKeysByPattern(pattern string) error {
	keys, err := GetKeysByPattern(pattern)
	if err != nil {
		return err
	}

	if len(keys) > 0 {
		return DeleteKeys(keys...)
	}

	return nil
}
