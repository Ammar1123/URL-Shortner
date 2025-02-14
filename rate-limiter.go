package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// rateLimitMiddleware enforces a limit of 10 requests per minute per IP.
func rateLimitMiddleware(c *gin.Context) {
	ip := c.ClientIP()
	key := "rate:" + ip

	// Increment the counter for this IP.
	count, err := redisClient.Incr(ctx, key).Result()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Rate limiting error"})
		return
	}

	if count == 1 {
		redisClient.Expire(ctx, key, time.Minute)
	}

	if count > 10 {
		c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit exceeded. Try again later."})
		return
	}
	c.Next()
}
