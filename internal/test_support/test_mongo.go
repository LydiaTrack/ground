package test_support

import (
	"github.com/Lydia/lydia-base/internal/mongodb"
	"github.com/Lydia/lydia-base/internal/utils"
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
