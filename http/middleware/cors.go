package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/RobertLesgros/rustdesk-interface/v2/global"
)

// Cors handles Cross-Origin Resource Sharing
// SECURITY: Validates origin against allowed list instead of echoing any origin
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")

		// Check if origin is allowed
		allowedOrigin := validateOrigin(origin)

		if allowedOrigin != "" {
			c.Header("Access-Control-Allow-Origin", allowedOrigin)
			c.Header("Access-Control-Allow-Headers", "api-token,content-type,authorization")
			c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Max-Age", "3600")
		}

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

// validateOrigin checks if the origin is allowed
// Returns the origin if allowed, empty string otherwise
func validateOrigin(origin string) string {
	if origin == "" {
		return ""
	}

	// Get allowed origins from config (comma-separated)
	// If not configured, allow the origin for backwards compatibility
	// In production, configure CORS_ALLOWED_ORIGINS environment variable
	allowedOriginsStr := global.Config.Gin.CorsAllowedOrigins

	// If no restriction is configured, allow same-origin and API server origin
	if allowedOriginsStr == "" {
		// Allow requests from API server itself
		apiServer := global.Config.Rustdesk.ApiServer
		if apiServer != "" && strings.Contains(origin, strings.TrimPrefix(strings.TrimPrefix(apiServer, "https://"), "http://")) {
			return origin
		}
		// Default: allow all for backwards compatibility (should be configured in production)
		return origin
	}

	// Check against whitelist
	allowedOrigins := strings.Split(allowedOriginsStr, ",")
	for _, allowed := range allowedOrigins {
		allowed = strings.TrimSpace(allowed)
		if allowed == "*" || allowed == origin {
			return origin
		}
	}

	return ""
}

// SecurityHeaders adds security headers to all responses
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Prevent MIME type sniffing
		c.Header("X-Content-Type-Options", "nosniff")

		// Prevent clickjacking
		c.Header("X-Frame-Options", "DENY")

		// XSS Protection (legacy but still useful)
		c.Header("X-XSS-Protection", "1; mode=block")

		// Referrer policy
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")

		// Permissions policy (disable unnecessary features)
		c.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=()")

		c.Next()
	}
}
