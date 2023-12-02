package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"lydia-track-base/cmd/lydia-track-base/api"
	"lydia-track-base/internal/mongodb"
	"lydia-track-base/internal/service"
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
	// Initialize database
	mongodb.InitializeContainer()
	// Initialize routes
	initializeRoutes(r)
	// Initialize default user
	service.InitializeDefaultUser()
	// Initialize default role
	service.InitializeDefaultRole()
	// Initialize Casbin policy enforcer
	//auth.InitializePolicyEnforcer()

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
