package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/RobertLesgros/rustdesk-interface/v2/global"
	"github.com/RobertLesgros/rustdesk-interface/v2/http/response"
	"github.com/RobertLesgros/rustdesk-interface/v2/lib/audit"
)

// Limiter checks if the client IP is banned (for login attempts)
func Limiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		loginLimiter := global.LoginLimiter
		clientIp := c.ClientIP()
		banned, _ := loginLimiter.CheckSecurityStatus(clientIp)
		if banned {
			audit.LogIPBanned(c, "Too many failed login attempts", "30 minutes")
			response.Fail(c, http.StatusLocked, response.TranslateMsg(c, "Banned"))
			c.Abort()
			return
		}
		c.Next()
	}
}

// RateLimiter provides general rate limiting for sensitive endpoints
type RateLimiter struct {
	requests map[string][]time.Time
	mu       sync.Mutex
	limit    int           // Maximum requests
	window   time.Duration // Time window
}

// Global rate limiter instance for sensitive operations
var sensitiveRateLimiter = &RateLimiter{
	requests: make(map[string][]time.Time),
	limit:    10,              // 10 requests
	window:   1 * time.Minute, // per minute
}

// SensitiveOperationLimiter rate limits sensitive operations like password changes, user creation, etc.
func SensitiveOperationLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()

		if !sensitiveRateLimiter.allow(clientIP) {
			audit.LogRateLimited(c, c.Request.URL.Path)
			response.Fail(c, http.StatusTooManyRequests, response.TranslateMsg(c, "TooManyRequests"))
			c.Abort()
			return
		}

		c.Next()
	}
}

// allow checks if the request should be allowed based on rate limiting rules
func (rl *RateLimiter) allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	windowStart := now.Add(-rl.window)

	// Get existing requests for this key
	requests := rl.requests[key]

	// Filter out old requests outside the window
	var validRequests []time.Time
	for _, t := range requests {
		if t.After(windowStart) {
			validRequests = append(validRequests, t)
		}
	}

	// Check if limit exceeded
	if len(validRequests) >= rl.limit {
		rl.requests[key] = validRequests
		return false
	}

	// Add current request
	validRequests = append(validRequests, now)
	rl.requests[key] = validRequests

	return true
}

// CleanupRateLimiter periodically cleans up old entries
func init() {
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()

		for range ticker.C {
			sensitiveRateLimiter.cleanup()
		}
	}()
}

func (rl *RateLimiter) cleanup() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	windowStart := now.Add(-rl.window)

	for key, requests := range rl.requests {
		var validRequests []time.Time
		for _, t := range requests {
			if t.After(windowStart) {
				validRequests = append(validRequests, t)
			}
		}
		if len(validRequests) == 0 {
			delete(rl.requests, key)
		} else {
			rl.requests[key] = validRequests
		}
	}
}
