package audit

import (
	"encoding/json"
	"io"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// EventType represents the type of audit event
type EventType string

const (
	// Authentication events
	EventLoginSuccess      EventType = "LOGIN_SUCCESS"
	EventLoginFailed       EventType = "LOGIN_FAILED"
	EventLogout            EventType = "LOGOUT"
	EventPasswordChanged   EventType = "PASSWORD_CHANGED"
	EventPasswordResetReq  EventType = "PASSWORD_RESET_REQUEST"
	EventOAuthLogin        EventType = "OAUTH_LOGIN"
	EventOAuthBind         EventType = "OAUTH_BIND"
	EventSessionExpired    EventType = "SESSION_EXPIRED"

	// User management events
	EventUserCreated       EventType = "USER_CREATED"
	EventUserUpdated       EventType = "USER_UPDATED"
	EventUserDeleted       EventType = "USER_DELETED"
	EventUserDisabled      EventType = "USER_DISABLED"
	EventUserEnabled       EventType = "USER_ENABLED"

	// Access control events
	EventAccessDenied      EventType = "ACCESS_DENIED"
	EventRateLimited       EventType = "RATE_LIMITED"
	EventIPBanned          EventType = "IP_BANNED"
	EventInvalidToken      EventType = "INVALID_TOKEN"

	// Administrative events
	EventConfigChanged     EventType = "CONFIG_CHANGED"
	EventGroupCreated      EventType = "GROUP_CREATED"
	EventGroupUpdated      EventType = "GROUP_UPDATED"
	EventGroupDeleted      EventType = "GROUP_DELETED"

	// Data access events
	EventDataExported      EventType = "DATA_EXPORTED"
	EventBulkOperation     EventType = "BULK_OPERATION"

	// Security events
	EventSecurityAlert     EventType = "SECURITY_ALERT"
	EventSuspiciousActivity EventType = "SUSPICIOUS_ACTIVITY"
)

// Severity represents the severity level of an audit event
type Severity string

const (
	SeverityInfo     Severity = "INFO"
	SeverityWarning  Severity = "WARNING"
	SeverityCritical Severity = "CRITICAL"
)

// AuditEvent represents a structured audit log entry
type AuditEvent struct {
	Timestamp   string            `json:"timestamp"`
	EventType   EventType         `json:"event_type"`
	Severity    Severity          `json:"severity"`
	UserID      uint              `json:"user_id,omitempty"`
	Username    string            `json:"username,omitempty"`
	ClientIP    string            `json:"client_ip"`
	UserAgent   string            `json:"user_agent,omitempty"`
	RequestID   string            `json:"request_id,omitempty"`
	Method      string            `json:"method,omitempty"`
	Path        string            `json:"path,omitempty"`
	StatusCode  int               `json:"status_code,omitempty"`
	Message     string            `json:"message"`
	Details     map[string]interface{} `json:"details,omitempty"`
	Success     bool              `json:"success"`
}

// AuditLogger provides structured security audit logging
type AuditLogger struct {
	writer io.Writer
	mu     sync.Mutex
}

var (
	defaultLogger *AuditLogger
	once          sync.Once
)

// Config holds the audit logger configuration
type Config struct {
	FilePath string
	Enabled  bool
}

// Init initializes the default audit logger
func Init(config *Config) error {
	if !config.Enabled {
		return nil
	}

	var writer io.Writer = os.Stdout

	if config.FilePath != "" {
		file, err := os.OpenFile(config.FilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
		if err != nil {
			return err
		}
		// Write to both file and stdout for monitoring
		writer = io.MultiWriter(file, os.Stdout)
	}

	once.Do(func() {
		defaultLogger = &AuditLogger{
			writer: writer,
		}
	})

	return nil
}

// GetLogger returns the default audit logger
func GetLogger() *AuditLogger {
	if defaultLogger == nil {
		// Return a no-op logger if not initialized
		return &AuditLogger{writer: io.Discard}
	}
	return defaultLogger
}

// Log writes an audit event to the log
func (al *AuditLogger) Log(event *AuditEvent) {
	if al.writer == nil {
		return
	}

	event.Timestamp = time.Now().UTC().Format(time.RFC3339)

	al.mu.Lock()
	defer al.mu.Unlock()

	data, err := json.Marshal(event)
	if err != nil {
		return
	}

	al.writer.Write(append(data, '\n'))
}

// Helper functions for common audit events

// LogLoginSuccess logs a successful login attempt
func LogLoginSuccess(c *gin.Context, userID uint, username string) {
	GetLogger().Log(&AuditEvent{
		EventType: EventLoginSuccess,
		Severity:  SeverityInfo,
		UserID:    userID,
		Username:  username,
		ClientIP:  c.ClientIP(),
		UserAgent: c.Request.UserAgent(),
		Method:    c.Request.Method,
		Path:      c.Request.URL.Path,
		Message:   "User logged in successfully",
		Success:   true,
	})
}

// LogLoginFailed logs a failed login attempt
func LogLoginFailed(c *gin.Context, username, reason string) {
	GetLogger().Log(&AuditEvent{
		EventType: EventLoginFailed,
		Severity:  SeverityWarning,
		Username:  username,
		ClientIP:  c.ClientIP(),
		UserAgent: c.Request.UserAgent(),
		Method:    c.Request.Method,
		Path:      c.Request.URL.Path,
		Message:   "Login attempt failed: " + reason,
		Success:   false,
		Details: map[string]interface{}{
			"reason": reason,
		},
	})
}

// LogLogout logs a user logout
func LogLogout(c *gin.Context, userID uint, username string) {
	GetLogger().Log(&AuditEvent{
		EventType: EventLogout,
		Severity:  SeverityInfo,
		UserID:    userID,
		Username:  username,
		ClientIP:  c.ClientIP(),
		UserAgent: c.Request.UserAgent(),
		Method:    c.Request.Method,
		Path:      c.Request.URL.Path,
		Message:   "User logged out",
		Success:   true,
	})
}

// LogPasswordChanged logs a password change event
func LogPasswordChanged(c *gin.Context, userID uint, username string, changedByAdmin bool) {
	details := map[string]interface{}{
		"changed_by_admin": changedByAdmin,
	}

	GetLogger().Log(&AuditEvent{
		EventType: EventPasswordChanged,
		Severity:  SeverityInfo,
		UserID:    userID,
		Username:  username,
		ClientIP:  c.ClientIP(),
		UserAgent: c.Request.UserAgent(),
		Method:    c.Request.Method,
		Path:      c.Request.URL.Path,
		Message:   "Password changed",
		Success:   true,
		Details:   details,
	})
}

// LogUserCreated logs a user creation event
func LogUserCreated(c *gin.Context, createdUserID uint, createdUsername string, creatorID uint) {
	GetLogger().Log(&AuditEvent{
		EventType: EventUserCreated,
		Severity:  SeverityInfo,
		UserID:    creatorID,
		ClientIP:  c.ClientIP(),
		UserAgent: c.Request.UserAgent(),
		Method:    c.Request.Method,
		Path:      c.Request.URL.Path,
		Message:   "New user created: " + createdUsername,
		Success:   true,
		Details: map[string]interface{}{
			"created_user_id":   createdUserID,
			"created_username":  createdUsername,
		},
	})
}

// LogUserDeleted logs a user deletion event
func LogUserDeleted(c *gin.Context, deletedUserID uint, deletedUsername string, deletorID uint) {
	GetLogger().Log(&AuditEvent{
		EventType: EventUserDeleted,
		Severity:  SeverityWarning,
		UserID:    deletorID,
		ClientIP:  c.ClientIP(),
		UserAgent: c.Request.UserAgent(),
		Method:    c.Request.Method,
		Path:      c.Request.URL.Path,
		Message:   "User deleted: " + deletedUsername,
		Success:   true,
		Details: map[string]interface{}{
			"deleted_user_id":   deletedUserID,
			"deleted_username":  deletedUsername,
		},
	})
}

// LogAccessDenied logs an access denied event
func LogAccessDenied(c *gin.Context, userID uint, resource, reason string) {
	GetLogger().Log(&AuditEvent{
		EventType: EventAccessDenied,
		Severity:  SeverityWarning,
		UserID:    userID,
		ClientIP:  c.ClientIP(),
		UserAgent: c.Request.UserAgent(),
		Method:    c.Request.Method,
		Path:      c.Request.URL.Path,
		Message:   "Access denied to resource: " + resource,
		Success:   false,
		Details: map[string]interface{}{
			"resource": resource,
			"reason":   reason,
		},
	})
}

// LogRateLimited logs a rate limiting event
func LogRateLimited(c *gin.Context, endpoint string) {
	GetLogger().Log(&AuditEvent{
		EventType: EventRateLimited,
		Severity:  SeverityWarning,
		ClientIP:  c.ClientIP(),
		UserAgent: c.Request.UserAgent(),
		Method:    c.Request.Method,
		Path:      c.Request.URL.Path,
		Message:   "Rate limit exceeded for endpoint: " + endpoint,
		Success:   false,
		Details: map[string]interface{}{
			"endpoint": endpoint,
		},
	})
}

// LogIPBanned logs an IP ban event
func LogIPBanned(c *gin.Context, reason string, duration string) {
	GetLogger().Log(&AuditEvent{
		EventType: EventIPBanned,
		Severity:  SeverityCritical,
		ClientIP:  c.ClientIP(),
		UserAgent: c.Request.UserAgent(),
		Method:    c.Request.Method,
		Path:      c.Request.URL.Path,
		Message:   "IP address banned",
		Success:   false,
		Details: map[string]interface{}{
			"reason":   reason,
			"duration": duration,
		},
	})
}

// LogSecurityAlert logs a security alert
func LogSecurityAlert(c *gin.Context, alertType, description string) {
	GetLogger().Log(&AuditEvent{
		EventType: EventSecurityAlert,
		Severity:  SeverityCritical,
		ClientIP:  c.ClientIP(),
		UserAgent: c.Request.UserAgent(),
		Method:    c.Request.Method,
		Path:      c.Request.URL.Path,
		Message:   "Security alert: " + alertType,
		Success:   false,
		Details: map[string]interface{}{
			"alert_type":  alertType,
			"description": description,
		},
	})
}

// LogOAuthLogin logs an OAuth login event
func LogOAuthLogin(c *gin.Context, userID uint, username, provider string) {
	GetLogger().Log(&AuditEvent{
		EventType: EventOAuthLogin,
		Severity:  SeverityInfo,
		UserID:    userID,
		Username:  username,
		ClientIP:  c.ClientIP(),
		UserAgent: c.Request.UserAgent(),
		Method:    c.Request.Method,
		Path:      c.Request.URL.Path,
		Message:   "User logged in via OAuth provider: " + provider,
		Success:   true,
		Details: map[string]interface{}{
			"provider": provider,
		},
	})
}

// LogInvalidToken logs an invalid token event
func LogInvalidToken(c *gin.Context, reason string) {
	GetLogger().Log(&AuditEvent{
		EventType: EventInvalidToken,
		Severity:  SeverityWarning,
		ClientIP:  c.ClientIP(),
		UserAgent: c.Request.UserAgent(),
		Method:    c.Request.Method,
		Path:      c.Request.URL.Path,
		Message:   "Invalid authentication token",
		Success:   false,
		Details: map[string]interface{}{
			"reason": reason,
		},
	})
}

// LogBulkOperation logs a bulk data operation
func LogBulkOperation(c *gin.Context, userID uint, operation string, count int) {
	GetLogger().Log(&AuditEvent{
		EventType: EventBulkOperation,
		Severity:  SeverityInfo,
		UserID:    userID,
		ClientIP:  c.ClientIP(),
		UserAgent: c.Request.UserAgent(),
		Method:    c.Request.Method,
		Path:      c.Request.URL.Path,
		Message:   "Bulk operation performed: " + operation,
		Success:   true,
		Details: map[string]interface{}{
			"operation":     operation,
			"affected_count": count,
		},
	})
}
