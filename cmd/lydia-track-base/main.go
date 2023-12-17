package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/kr/pretty"
	"log"
	"lydia-track-base/cmd/lydia-track-base/api"
	"lydia-track-base/internal/initializers"
	"lydia-track-base/internal/mongodb"
	"lydia-track-base/internal/utils"
)

// @title Lydia Track Base API
// @version 0.0.1
// @description Lydia Track Base API

// @host localhost:8080
// @BasePath /
// @schemes http
func main() {

	// Initialize environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// Initialize logging
	utils.InitLogging()
	r := gin.New()
	// Initialize database connection
	mongodb.InitializeMongoDBConnection()
	// Initialize API routes
	initializeRoutes(r)
	// Initialize default user
	err = initializers.InitializeDefaultUser()
	if err != nil {
		log.Fatal("Error initializing default user")
		panic(err)
	}
	// Initialize default role
	err = initializers.InitializeDefaultRole()
	if err != nil {
		pretty.Errorf("Error initializing default role")
		panic(err)
	}

	// Run server on port 8080
	r.Run(":8080")
}

// initializeRoutes initializes routes for each API
func initializeRoutes(r *gin.Engine) {
	globalInterceptors := []gin.HandlerFunc{gin.Recovery(), gin.Logger()}

	r.Use(globalInterceptors...)

	api.InitUser(r)
	api.InitRole(r)
	api.InitSwagger(r)
	api.InitAuth(r)
	api.InitHealth(r)
}
