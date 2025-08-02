package middlewares

import (
	"bytes"
	"encoding/json"
	"io"
	"net/url"
	"strings"
	"time"

	"github.com/LydiaTrack/ground/pkg/auth"
	"github.com/LydiaTrack/ground/pkg/log"
	"github.com/gin-gonic/gin"
)

type RequestLogMiddleware struct {
	authService auth.Service
	userService auth.UserService
}

type RequestLogData struct {
	Timestamp   time.Time              `json:"timestamp"`
	Method      string                 `json:"method"`
	Endpoint    string                 `json:"endpoint"`
	FullPath    string                 `json:"fullPath"`
	Username    string                 `json:"username,omitempty"`
	UserID      string                 `json:"userId,omitempty"`
	QueryParams map[string]interface{} `json:"queryParams,omitempty"`
	FormParams  map[string]interface{} `json:"formParams,omitempty"`
	JSONBody    map[string]interface{} `json:"jsonBody,omitempty"`
	ContentType string                 `json:"contentType,omitempty"`
	ClientIP    string                 `json:"clientIP"`
	UserAgent   string                 `json:"userAgent,omitempty"`
	Duration    string                 `json:"duration"`
	StatusCode  int                    `json:"statusCode"`
}

func NewRequestLogMiddleware(authService auth.Service, userService auth.UserService) *RequestLogMiddleware {
	return &RequestLogMiddleware{
		authService: authService,
		userService: userService,
	}
}

// RequestLoggingMiddleware logs each request with endpoint, username and parameters
func (m *RequestLogMiddleware) RequestLoggingMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// This is called after the request is processed, but we need to capture
		// request data before processing. We'll use a different approach.
		return ""
	})
}

// RequestLoggingMiddleware logs each request with endpoint, username and parameters
func (m *RequestLogMiddleware) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		// Create a copy of the request body for logging
		var bodyBytes []byte
		if c.Request.Body != nil {
			bodyBytes, _ = io.ReadAll(c.Request.Body)
			// Restore the request body for subsequent handlers
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		// Process the request
		c.Next()

		// Log the request after processing
		duration := time.Since(startTime)
		m.logRequest(c, bodyBytes, duration)
	}
}

func (m *RequestLogMiddleware) logRequest(c *gin.Context, bodyBytes []byte, duration time.Duration) {
	logData := RequestLogData{
		Timestamp:   time.Now(),
		Method:      c.Request.Method,
		Endpoint:    c.Request.URL.Path,
		FullPath:    c.FullPath(),
		ClientIP:    c.ClientIP(),
		UserAgent:   c.Request.UserAgent(),
		Duration:    duration.String(),
		StatusCode:  c.Writer.Status(),
		ContentType: c.Request.Header.Get("Content-Type"),
	}

	// Extract username and user ID from authentication context
	if currentUser, err := m.authService.GetCurrentUser(c); err == nil {
		logData.Username = currentUser.Username
		logData.UserID = currentUser.ID.Hex()
	}

	// Extract query parameters
	if len(c.Request.URL.RawQuery) > 0 {
		logData.QueryParams = m.parseQueryParams(c.Request.URL.Query())
	}

	// Extract form parameters (for form-encoded requests)
	if strings.Contains(logData.ContentType, "application/x-www-form-urlencoded") ||
		strings.Contains(logData.ContentType, "multipart/form-data") {
		if err := c.Request.ParseForm(); err == nil {
			logData.FormParams = m.parseFormParams(c.Request.PostForm)
		}
	}

	// Extract JSON body parameters (for JSON requests)
	if strings.Contains(logData.ContentType, "application/json") && len(bodyBytes) > 0 {
		var jsonBody map[string]interface{}
		if err := json.Unmarshal(bodyBytes, &jsonBody); err == nil {
			// Remove sensitive fields from logging
			logData.JSONBody = m.sanitizeJSONBody(jsonBody)
		}
	}

	// Convert to JSON string for logging
	if jsonStr, err := json.Marshal(logData); err == nil {
		log.LogRequest("%s", string(jsonStr))
	} else {
		// Fallback to simple logging if JSON marshaling fails
		username := logData.Username
		if username == "" {
			username = "anonymous"
		}
		log.LogRequest("%s %s [%s] - %d - %s",
			logData.Method, logData.Endpoint, username, logData.StatusCode, logData.Duration)
	}
}

func (m *RequestLogMiddleware) parseQueryParams(values url.Values) map[string]interface{} {
	params := make(map[string]interface{})
	for key, values := range values {
		if len(values) == 1 {
			params[key] = values[0]
		} else {
			params[key] = values
		}
	}
	return params
}

func (m *RequestLogMiddleware) parseFormParams(values url.Values) map[string]interface{} {
	params := make(map[string]interface{})
	for key, values := range values {
		if len(values) == 1 {
			params[key] = values[0]
		} else {
			params[key] = values
		}
	}
	return params
}

func (m *RequestLogMiddleware) sanitizeJSONBody(body map[string]interface{}) map[string]interface{} {
	sanitized := make(map[string]interface{})

	// List of sensitive field names to exclude from logging
	sensitiveFields := map[string]bool{
		"password":        true,
		"newPassword":     true,
		"oldPassword":     true,
		"confirmPassword": true,
		"token":           true,
		"refreshToken":    true,
		"secret":          true,
		"apiKey":          true,
		"accessToken":     true,
		"auth":            true,
		"authorization":   true,
	}

	for key, value := range body {
		lowerKey := strings.ToLower(key)
		if sensitiveFields[lowerKey] {
			sanitized[key] = "[REDACTED]"
		} else {
			// For nested objects, recursively sanitize
			if nestedMap, ok := value.(map[string]interface{}); ok {
				sanitized[key] = m.sanitizeJSONBody(nestedMap)
			} else {
				sanitized[key] = value
			}
		}
	}

	return sanitized
}

// Simple function for backward compatibility
func RequestLoggingMiddleware(authService auth.Service, userService auth.UserService) gin.HandlerFunc {
	middleware := NewRequestLogMiddleware(authService, userService)
	return middleware.Middleware()
}
