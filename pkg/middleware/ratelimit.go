package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

)

type RateLimitConfig struct {
	Limit  int
	Window time.Duration
}

func RateLimit(rdb *redis.Client, config RateLimitConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		ip := c.ClientIP()
		key := fmt.Sprintf("rl:%s", ip)

		count, err := rdb.Incr(ctx, key).Result()
		if err != nil {
			log.Printf("⚠️ Redis Rate Limit Error: %v", err)
			c.Next()
			return
		}

		if count == 1 {
			rdb.Expire(ctx, key, config.Window)
		}

		remaining := int64(config.Limit) - count
		if remaining < 0 {
			remaining = 0
		}
		
		c.Header("X-RateLimit-Limit", strconv.Itoa(config.Limit))
		c.Header("X-RateLimit-Remaining", strconv.FormatInt(remaining, 10))

		if count > int64(config.Limit) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many requests. Please slow down.",
			})
			return
		}

		c.Next()
	}
}