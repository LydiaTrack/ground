package test_support

import (
	"github.com/joho/godotenv"
	"log"
	"lydia-track-base/internal/mongodb"
	"lydia-track-base/internal/utils"
)

func TestWithMongo() {
	// Initialize logging
	utils.InitLogging()
	// Initialize environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	mongodb.InitializeContainer()
}
