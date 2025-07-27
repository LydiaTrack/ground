package test_support

import (
	"os"

	"github.com/LydiaTrack/ground/pkg/log"
	"github.com/LydiaTrack/ground/pkg/mongodb"
	"github.com/joho/godotenv"
)

func TestWithMongo() {
	// Initialize logging
	log.InitLogging()

	// Set default test environment variables if not already set
	setDefaultTestEnvVars()

	// Try to load environment file, but don't fail if it doesn't exist
	fileName := ""
	envType := os.Getenv("ENV_TYPE")
	if envType == "production" {
		fileName = ".env.production"
	} else if envType == "test" {
		fileName = ".env.test"
	} else {
		fileName = ".env.development"
	}

	// Try to load the env file, but continue if it doesn't exist
	_ = godotenv.Load(fileName)

	mongodb.InitializeMongoDBConnection()
}

// setDefaultTestEnvVars sets default environment variables for testing
func setDefaultTestEnvVars() {
	// Set defaults only if not already set
	if os.Getenv("DB_CONNECTION_TYPE") == "" {
		os.Setenv("DB_CONNECTION_TYPE", "CONTAINER")
	}
	if os.Getenv("DB_PORT") == "" {
		os.Setenv("DB_PORT", "27017")
	}
	if os.Getenv("DB_NAME") == "" {
		os.Setenv("DB_NAME", "test_db")
	}
	if os.Getenv("JWT_SECRET") == "" {
		os.Setenv("JWT_SECRET", "test_jwt_secret_for_testing_only")
	}
	if os.Getenv("JWT_EXPIRES_IN_MINUTES") == "" {
		os.Setenv("JWT_EXPIRES_IN_MINUTES", "60")
	}
	if os.Getenv("JWT_REFRESH_EXPIRES_IN_HOUR") == "" {
		os.Setenv("JWT_REFRESH_EXPIRES_IN_HOUR", "72")
	}
	if os.Getenv("DEFAULT_USER_USERNAME") == "" {
		os.Setenv("DEFAULT_USER_USERNAME", "test_admin")
	}
	if os.Getenv("DEFAULT_USER_PASSWORD") == "" {
		os.Setenv("DEFAULT_USER_PASSWORD", "test_password")
	}
	if os.Getenv("DEFAULT_ROLE_NAME") == "" {
		os.Setenv("DEFAULT_ROLE_NAME", "TEST_ADMIN")
	}
	if os.Getenv("DEFAULT_ROLE_TAGS") == "" {
		os.Setenv("DEFAULT_ROLE_TAGS", "TEST_ROLE")
	}
	if os.Getenv("DEFAULT_ROLE_INFO") == "" {
		os.Setenv("DEFAULT_ROLE_INFO", "Test administrator role")
	}
	if os.Getenv("ENV_TYPE") == "" {
		os.Setenv("ENV_TYPE", "test")
	}
}
