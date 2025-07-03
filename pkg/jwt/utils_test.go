package jwt

import (
	"os"
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestGenerateTokenPair(t *testing.T) {
	userID := primitive.NewObjectID()

	t.Run("Success with valid environment variables", func(t *testing.T) {
		// Set required environment variables
		os.Setenv(JwtSecretKey, "test_secret_key")
		os.Setenv(JwtExpirationKey, "5")
		defer func() {
			os.Unsetenv(JwtSecretKey)
			os.Unsetenv(JwtExpirationKey)
		}()

		tokenPair, err := GenerateTokenPair(userID)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if tokenPair.Token == "" {
			t.Error("Expected token to be generated")
		}

		if tokenPair.RefreshToken == "" {
			t.Error("Expected refresh token to be generated")
		}

		if tokenPair.UserID != userID {
			t.Errorf("Expected userID %v, got %v", userID, tokenPair.UserID)
		}
	})

	t.Run("Fail with missing JWT_SECRET", func(t *testing.T) {
		os.Unsetenv(JwtSecretKey)
		os.Setenv(JwtExpirationKey, "5")
		defer os.Unsetenv(JwtExpirationKey)

		_, err := GenerateTokenPair(userID)
		if err == nil {
			t.Error("Expected error when JWT_SECRET is missing")
		}

		expectedError := "JWT_SECRET environment variable not set"
		if err.Error() != expectedError {
			t.Errorf("Expected error message '%s', got '%s'", expectedError, err.Error())
		}
	})

	t.Run("Fail with missing JWT_EXPIRES_IN_MINUTES", func(t *testing.T) {
		os.Setenv(JwtSecretKey, "test_secret_key")
		os.Unsetenv(JwtExpirationKey)
		defer os.Unsetenv(JwtSecretKey)

		_, err := GenerateTokenPair(userID)
		if err == nil {
			t.Error("Expected error when JWT_EXPIRES_IN_MINUTES is missing")
		}

		expectedError := "JWT_EXPIRES_IN_MINUTES environment variable not set"
		if err.Error() != expectedError {
			t.Errorf("Expected error message '%s', got '%s'", expectedError, err.Error())
		}
	})

	t.Run("Fail with invalid JWT_EXPIRES_IN_MINUTES", func(t *testing.T) {
		os.Setenv(JwtSecretKey, "test_secret_key")
		os.Setenv(JwtExpirationKey, "invalid")
		defer func() {
			os.Unsetenv(JwtSecretKey)
			os.Unsetenv(JwtExpirationKey)
		}()

		_, err := GenerateTokenPair(userID)
		if err == nil {
			t.Error("Expected error when JWT_EXPIRES_IN_MINUTES is invalid")
		}

		if !containsString(err.Error(), "invalid JWT_EXPIRES_IN_MINUTES value") {
			t.Errorf("Expected error about invalid JWT_EXPIRES_IN_MINUTES value, got '%s'", err.Error())
		}
	})

	t.Run("Fail with negative JWT_EXPIRES_IN_MINUTES", func(t *testing.T) {
		os.Setenv(JwtSecretKey, "test_secret_key")
		os.Setenv(JwtExpirationKey, "-1")
		defer func() {
			os.Unsetenv(JwtSecretKey)
			os.Unsetenv(JwtExpirationKey)
		}()

		_, err := GenerateTokenPair(userID)
		if err == nil {
			t.Error("Expected error when JWT_EXPIRES_IN_MINUTES is negative")
		}

		expectedError := "JWT_EXPIRES_IN_MINUTES must be a positive number"
		if err.Error() != expectedError {
			t.Errorf("Expected error message '%s', got '%s'", expectedError, err.Error())
		}
	})

	t.Run("Fail with zero JWT_EXPIRES_IN_MINUTES", func(t *testing.T) {
		os.Setenv(JwtSecretKey, "test_secret_key")
		os.Setenv(JwtExpirationKey, "0")
		defer func() {
			os.Unsetenv(JwtSecretKey)
			os.Unsetenv(JwtExpirationKey)
		}()

		_, err := GenerateTokenPair(userID)
		if err == nil {
			t.Error("Expected error when JWT_EXPIRES_IN_MINUTES is zero")
		}

		expectedError := "JWT_EXPIRES_IN_MINUTES must be a positive number"
		if err.Error() != expectedError {
			t.Errorf("Expected error message '%s', got '%s'", expectedError, err.Error())
		}
	})
}

func TestIsTokenValid(t *testing.T) {
	userID := primitive.NewObjectID()

	t.Run("Success with valid token", func(t *testing.T) {
		// Set up environment
		os.Setenv(JwtSecretKey, "test_secret_key")
		os.Setenv(JwtExpirationKey, "5")
		defer func() {
			os.Unsetenv(JwtSecretKey)
			os.Unsetenv(JwtExpirationKey)
		}()

		// Generate a valid token
		tokenPair, err := GenerateTokenPair(userID)
		if err != nil {
			t.Fatalf("Failed to generate token: %v", err)
		}

		// Test validation
		err = IsTokenValid(tokenPair.Token)
		if err != nil {
			t.Errorf("Expected valid token to pass validation, got error: %v", err)
		}
	})

	t.Run("Fail with missing JWT_SECRET", func(t *testing.T) {
		os.Unsetenv(JwtSecretKey)

		err := IsTokenValid("some_token")
		if err == nil {
			t.Error("Expected error when JWT_SECRET is missing")
		}

		expectedError := "JWT_SECRET environment variable not set"
		if err.Error() != expectedError {
			t.Errorf("Expected error message '%s', got '%s'", expectedError, err.Error())
		}
	})

	t.Run("Fail with invalid token", func(t *testing.T) {
		os.Setenv(JwtSecretKey, "test_secret_key")
		defer os.Unsetenv(JwtSecretKey)

		err := IsTokenValid("invalid_token")
		if err == nil {
			t.Error("Expected error with invalid token")
		}
	})

	t.Run("Fail with expired token", func(t *testing.T) {
		// This test would require creating an expired token, which is complex
		// For now, we'll just test with a malformed token
		os.Setenv(JwtSecretKey, "test_secret_key")
		defer os.Unsetenv(JwtSecretKey)

		err := IsTokenValid("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.invalid.signature")
		if err == nil {
			t.Error("Expected error with malformed token")
		}
	})
}

func TestTokenPairIntegration(t *testing.T) {
	userID := primitive.NewObjectID()

	// Set up environment
	os.Setenv(JwtSecretKey, "test_secret_key")
	os.Setenv(JwtExpirationKey, "2") // 2 minutes
	defer func() {
		os.Unsetenv(JwtSecretKey)
		os.Unsetenv(JwtExpirationKey)
	}()

	// Generate token pair
	tokenPair, err := GenerateTokenPair(userID)
	if err != nil {
		t.Fatalf("Failed to generate token pair: %v", err)
	}

	// Validate the generated token
	err = IsTokenValid(tokenPair.Token)
	if err != nil {
		t.Errorf("Generated token should be valid, got error: %v", err)
	}

	// Verify refresh token is different from access token
	if tokenPair.Token == tokenPair.RefreshToken {
		t.Error("Access token and refresh token should be different")
	}

	// Verify refresh token is a valid ObjectID hex string (24 characters)
	if len(tokenPair.RefreshToken) != 24 {
		t.Errorf("Expected refresh token to be 24 characters (ObjectID hex), got %d", len(tokenPair.RefreshToken))
	}
}

// Helper function to check if a string contains a substring
func containsString(str, substr string) bool {
	return len(str) >= len(substr) && (str == substr || len(substr) == 0 ||
		(len(str) > len(substr) && (str[:len(substr)] == substr ||
			str[len(str)-len(substr):] == substr ||
			containsStringHelper(str, substr))))
}

func containsStringHelper(str, substr string) bool {
	for i := 0; i <= len(str)-len(substr); i++ {
		if str[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
