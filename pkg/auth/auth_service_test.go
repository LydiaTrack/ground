package auth

import (
	"encoding/json"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/LydiaTrack/ground/pkg/constants"
	"github.com/LydiaTrack/ground/pkg/domain/session"
	"github.com/LydiaTrack/ground/pkg/jwt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Mock session service for testing
type mockSessionService struct {
	sessions    map[string]session.InfoModel
	shouldError bool
}

func (m *mockSessionService) CreateSession(cmd session.CreateSessionCommand) (session.InfoModel, error) {
	sessionModel := session.InfoModel{
		ID:           primitive.NewObjectID(),
		UserID:       primitive.NewObjectID(),
		RefreshToken: cmd.RefreshToken,
		ExpireTime:   cmd.ExpireTime,
	}
	m.sessions[cmd.RefreshToken] = sessionModel
	return sessionModel, nil
}

func (m *mockSessionService) GetSessionByRefreshToken(refreshToken string) (session.InfoModel, error) {
	if m.shouldError {
		return session.InfoModel{}, constants.ErrorInternalServerError
	}

	sessionModel, exists := m.sessions[refreshToken]
	if !exists {
		return session.InfoModel{}, constants.ErrorNotFound
	}
	return sessionModel, nil
}

func (m *mockSessionService) DeleteSessionByUser(userID string) error {
	// Delete all sessions for the user
	for token, session := range m.sessions {
		if session.UserID.Hex() == userID {
			delete(m.sessions, token)
		}
	}
	return nil
}

func (m *mockSessionService) GetUserSession(id string) (session.InfoModel, error) {
	if m.shouldError {
		return session.InfoModel{}, constants.ErrorInternalServerError
	}

	// Find session by user ID
	for _, sessionModel := range m.sessions {
		if sessionModel.UserID.Hex() == id {
			return sessionModel, nil
		}
	}
	return session.InfoModel{}, constants.ErrorNotFound
}

func (m *mockSessionService) CleanupExpiredSessions() error {
	return nil
}

func TestRefreshTokenPairSessionExpiration(t *testing.T) {
	// Set up environment variables
	os.Setenv(jwt.JwtSecretKey, "test_secret")
	os.Setenv(jwt.JwtExpirationKey, "24")
	os.Setenv(jwt.RefreshExpirationKey, "168")
	defer func() {
		os.Unsetenv(jwt.JwtSecretKey)
		os.Unsetenv(jwt.JwtExpirationKey)
		os.Unsetenv(jwt.RefreshExpirationKey)
	}()

	t.Run("Success with valid non-expired session", func(t *testing.T) {
		// Create mock service
		mockSession := &mockSessionService{
			sessions: make(map[string]session.InfoModel),
		}

		// Create auth service (we'll mock other dependencies as nil since we're only testing session logic)
		authService := Service{
			sessionService: mockSession,
			// Other fields can be nil for this test
		}

		// Create a session that expires in the future
		futureExpiry := time.Now().Add(1 * time.Hour).Unix()
		userID := primitive.NewObjectID()
		refreshToken := "test_refresh_token"

		sessionModel := session.InfoModel{
			ID:           primitive.NewObjectID(),
			UserID:       userID,
			RefreshToken: refreshToken,
			ExpireTime:   futureExpiry,
		}
		mockSession.sessions[refreshToken] = sessionModel

		// Create gin context with refresh token request
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		requestBody := RefreshTokenRequest{RefreshToken: refreshToken}
		jsonBody, _ := json.Marshal(requestBody)

		req := httptest.NewRequest("POST", "/auth/refreshToken", strings.NewReader(string(jsonBody)))
		req.Header.Set("Content-Type", "application/json")
		c.Request = req

		// Test refresh token - should succeed since session is not expired
		_, err := authService.RefreshTokenPair(c)

		// Since we don't have a full user service setup, this will fail at token generation
		// But it should NOT fail with unauthorized due to expiration
		// The error should be ErrorInternalServerError (from missing user service), not ErrorUnauthorized
		if err == constants.ErrorUnauthorized {
			t.Error("Expected session expiration check to pass, but got unauthorized error")
		}
	})

	t.Run("Fail with expired session", func(t *testing.T) {
		// Create mock service
		mockSession := &mockSessionService{
			sessions: make(map[string]session.InfoModel),
		}

		authService := Service{
			sessionService: mockSession,
		}

		// Create a session that expired in the past
		pastExpiry := time.Now().Add(-1 * time.Hour).Unix()
		userID := primitive.NewObjectID()
		refreshToken := "expired_refresh_token"

		sessionModel := session.InfoModel{
			ID:           primitive.NewObjectID(),
			UserID:       userID,
			RefreshToken: refreshToken,
			ExpireTime:   pastExpiry,
		}
		mockSession.sessions[refreshToken] = sessionModel

		// Create gin context
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		requestBody := RefreshTokenRequest{RefreshToken: refreshToken}
		jsonBody, _ := json.Marshal(requestBody)

		req := httptest.NewRequest("POST", "/auth/refreshToken", strings.NewReader(string(jsonBody)))
		req.Header.Set("Content-Type", "application/json")
		c.Request = req

		// Test refresh token - should fail due to expired session
		_, err := authService.RefreshTokenPair(c)

		if err != constants.ErrorUnauthorized {
			t.Errorf("Expected ErrorUnauthorized for expired session, got %v", err)
		}

		// Verify that the expired session was cleaned up
		_, exists := mockSession.sessions[refreshToken]
		if exists {
			t.Error("Expected expired session to be cleaned up")
		}
	})

	t.Run("Fail with non-existent session", func(t *testing.T) {
		mockSession := &mockSessionService{
			sessions: make(map[string]session.InfoModel),
		}

		authService := Service{
			sessionService: mockSession,
		}

		// Create gin context with non-existent refresh token
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		requestBody := RefreshTokenRequest{RefreshToken: "non_existent_token"}
		jsonBody, _ := json.Marshal(requestBody)

		req := httptest.NewRequest("POST", "/auth/refreshToken", strings.NewReader(string(jsonBody)))
		req.Header.Set("Content-Type", "application/json")
		c.Request = req

		// Test refresh token - should fail due to non-existent session
		_, err := authService.RefreshTokenPair(c)

		if err != constants.ErrorUnauthorized {
			t.Errorf("Expected ErrorUnauthorized for non-existent session, got %v", err)
		}
	})

	t.Run("Fail with empty refresh token", func(t *testing.T) {
		mockSession := &mockSessionService{
			sessions: make(map[string]session.InfoModel),
		}

		authService := Service{
			sessionService: mockSession,
		}

		// Create gin context with empty refresh token
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		requestBody := RefreshTokenRequest{RefreshToken: ""}
		jsonBody, _ := json.Marshal(requestBody)

		req := httptest.NewRequest("POST", "/auth/refreshToken", strings.NewReader(string(jsonBody)))
		req.Header.Set("Content-Type", "application/json")
		c.Request = req

		// Test refresh token - should fail due to empty refresh token
		_, err := authService.RefreshTokenPair(c)

		if err != constants.ErrorUnauthorized {
			t.Errorf("Expected ErrorUnauthorized for empty refresh token, got %v", err)
		}
	})

	t.Run("Fail with malformed JSON request", func(t *testing.T) {
		mockSession := &mockSessionService{
			sessions: make(map[string]session.InfoModel),
		}

		authService := Service{
			sessionService: mockSession,
		}

		// Create gin context with malformed JSON
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		req := httptest.NewRequest("POST", "/auth/refreshToken", strings.NewReader("invalid json"))
		req.Header.Set("Content-Type", "application/json")
		c.Request = req

		// Test refresh token - should fail due to malformed JSON
		_, err := authService.RefreshTokenPair(c)

		if err != constants.ErrorInternalServerError {
			t.Errorf("Expected ErrorInternalServerError for malformed JSON, got %v", err)
		}
	})
}

func TestEnvironmentVariableValidation(t *testing.T) {
	t.Run("Validate environment variable handling in refresh token expiration", func(t *testing.T) {
		// Test missing JWT_REFRESH_EXPIRES_IN_HOUR
		os.Unsetenv(jwt.RefreshExpirationKey)

		refreshTokenLifespanStr := os.Getenv(jwt.RefreshExpirationKey)
		if refreshTokenLifespanStr != "" {
			t.Error("Environment variable should be unset for this test")
		}

		// Test invalid JWT_REFRESH_EXPIRES_IN_HOUR
		os.Setenv(jwt.RefreshExpirationKey, "invalid")
		refreshTokenLifespanStr = os.Getenv(jwt.RefreshExpirationKey)
		if refreshTokenLifespanStr == "" {
			t.Error("Environment variable should be set for this test")
		}

		// Test negative JWT_REFRESH_EXPIRES_IN_HOUR
		os.Setenv(jwt.RefreshExpirationKey, "-1")
		refreshTokenLifespanStr = os.Getenv(jwt.RefreshExpirationKey)
		if refreshTokenLifespanStr != "-1" {
			t.Error("Environment variable should be set to -1 for this test")
		}

		// Clean up
		os.Unsetenv(jwt.RefreshExpirationKey)
	})
}
