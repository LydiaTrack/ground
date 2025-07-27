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
	// Initialize environment variables
	fileName := ""
	envType := os.Getenv("ENV_TYPE")
	if envType == "production" {
		fileName = ".env.prod"
	} else if envType == "test" {
		fileName = ".env.test"
	} else {
		fileName = ".env.development"
	}
	err := godotenv.Load(fileName)
	if err != nil {
		log.LogFatal("Error loading .env file")
	}
	mongodb.InitializeMongoDBConnection()
}
