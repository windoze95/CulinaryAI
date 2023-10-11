package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/windoze95/saltybytes-api/internal/util"
	"golang.org/x/time/rate"
)

type limiterInfo struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

func RateLimitPublicOpenAIKey(publicOpenAIKeyRateLimiter *rate.Limiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve the user from the context
		user, err := util.GetUserFromContext(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if user.Settings.EncryptedOpenAIKey == "" {
			// Apply rate limiting and use shared key
			if !publicOpenAIKeyRateLimiter.Allow() {
				c.JSON(http.StatusTooManyRequests, gin.H{"error": "429: Too many requests"})
				c.Abort()
				return
			}
		}

		c.Next()
	}
}

func RateLimitByIP(rps int, cleanupInterval time.Duration, expiration time.Duration) gin.HandlerFunc {
	var limiters sync.Map

	// Cleanup goroutine
	go func() {
		for range time.Tick(cleanupInterval) {
			limiters.Range(func(key, value interface{}) bool {
				if time.Since(value.(*limiterInfo).lastSeen) > expiration {
					limiters.Delete(key)
				}
				return true
			})
		}
	}()

	return func(c *gin.Context) {
		ip := c.ClientIP()

		// Use LoadOrStore to ensure thread safety
		actual, _ := limiters.LoadOrStore(ip, &limiterInfo{
			limiter:  rate.NewLimiter(rate.Limit(rps), rps),
			lastSeen: time.Now(),
		})

		info := actual.(*limiterInfo)
		info.lastSeen = time.Now()

		if !info.limiter.Allow() {
			// Too many requests
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
			c.Abort()
			return
		}

		c.Next()
	}
}
