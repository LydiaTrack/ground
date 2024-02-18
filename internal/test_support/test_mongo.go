package test_support

import (
	"github.com/LydiaTrack/lydia-track-base/internal/mongodb"
	"github.com/LydiaTrack/lydia-track-base/internal/utils"
	"github.com/joho/godotenv"
	"log"
)

func TestWithMongo() {
	// Initialize logging
	utils.InitLogging()
	// Initialize environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	mongodb.InitializeMongoDBConnection()
}
