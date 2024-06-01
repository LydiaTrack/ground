package test_support

import (
	"github.com/LydiaTrack/lydia-base/internal/log"
	"github.com/LydiaTrack/lydia-base/pkg/mongodb"
	"github.com/joho/godotenv"
)

func TestWithMongo() {
	// Initialize logging
	log.InitLogging()
	// Initialize environment variables
	err := godotenv.Load()
	if err != nil {
		log.LogFatal("Error loading .env file")
	}
	mongodb.InitializeMongoDBConnection()
}
