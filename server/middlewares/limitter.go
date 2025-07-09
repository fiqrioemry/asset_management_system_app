package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/fiqrioemry/asset_management_system_app/server/config"

	"github.com/gin-gonic/gin"
)

const (
	RateLimitKeyPrefix = "rl:"
	RateLimitDuration  = time.Minute
	RateLimitMaxReq    = 5
)

func RateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := getClientIP(c)
		key := fmt.Sprintf("%s%s", RateLimitKeyPrefix, ip)

		count, err := config.RedisClient.Get(config.Ctx, key).Int()
		if err != nil && err.Error() != "redis: nil" {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Rate limiter error"})
			return
		}

		if count >= RateLimitMaxReq {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
			return
		}

		pipe := config.RedisClient.TxPipeline()
		pipe.Incr(config.Ctx, key)
		pipe.Expire(config.Ctx, key, RateLimitDuration)
		_, _ = pipe.Exec(config.Ctx)

		c.Next()
	}
}

// getClientIP extracts IP from request
func getClientIP(c *gin.Context) string {
	ip := c.ClientIP()
	if ip == "" {
		ip = "unknown"
	}
	return strings.TrimSpace(ip)
}
