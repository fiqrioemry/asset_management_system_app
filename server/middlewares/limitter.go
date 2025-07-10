package middlewares

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/fiqrioemry/asset_management_system_app/server/config"
	"github.com/gin-gonic/gin"
)

func RateLimiterInit() gin.HandlerFunc {
	return RateLimiter(config.AppConfig.RateLimitAttempts, config.AppConfig.RateLimitDuration)
}

func RateLimiter(maxAttempts int, duration time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := GetClientIP(c)
		key := fmt.Sprintf("ratelimit:%s", ip)

		count, _ := config.RedisClient.Get(config.Ctx, key).Int()

		if count >= maxAttempts {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":       "Too many requests",
				"message":     "Rate limit exceeded. Please try again later.",
				"retry_after": int(duration.Seconds()),
			})
			c.Abort()
			return
		}

		// Increment counter with expiration
		pipe := config.RedisClient.TxPipeline()
		pipe.Incr(config.Ctx, key)
		pipe.Expire(config.Ctx, key, duration)
		_, _ = pipe.Exec(config.Ctx)

		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", maxAttempts))
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", maxAttempts-count-1))
		c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", time.Now().Add(duration).Unix()))

		c.Next()
	}
}

func LimitFileSize(maxSize int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxSize)
		c.Next()
	}
}

func GetClientIP(c *gin.Context) string {
	if forwarded := c.Request.Header.Get("X-Forwarded-For"); forwarded != "" {
		ips := strings.Split(forwarded, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	if realIP := c.Request.Header.Get("X-Real-IP"); realIP != "" {
		return strings.TrimSpace(realIP)
	}

	return c.ClientIP()
}
