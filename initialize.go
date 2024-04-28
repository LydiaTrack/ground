package lydia_base

import (
	"github.com/LydiaTrack/lydia-base/internal/api"
	"github.com/LydiaTrack/lydia-base/internal/initializers"
	"github.com/LydiaTrack/lydia-base/internal/utils"
	"github.com/LydiaTrack/lydia-base/mongodb"
	"github.com/LydiaTrack/lydia-base/service_initializer"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
)

// Initialize initializes the Lydia base server with r as the gin Engine
func Initialize(r *gin.Engine) {
	// Initialize environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// Initialize logging
	utils.InitLogging()

	// Initialize database connection
	err = mongodb.InitializeMongoDBConnection()
	if err != nil {
		log.Fatal("Error initializing MongoDB connection")
	}

	// Initialize services
	service_initializer.InitializeServices()

	// Initialize API routes
	initializeRoutes(r, service_initializer.GetServices())

	// Initialize default user
	err = initializers.InitializeDefaultUser()
	if err != nil {
		log.Fatal("Error initializing default user")
	}

	// Initialize default role
	err = initializers.InitializeDefaultRole()
	if err != nil {
		log.Fatal("Error initializing default role")
	}

}

// initializeRoutes initializes routes for each API
func initializeRoutes(r *gin.Engine, services service_initializer.Services) {
	globalInterceptors := []gin.HandlerFunc{gin.Recovery(), gin.Logger()}

	r.Use(globalInterceptors...)

	api.InitAuth(r, services)
	api.InitUser(r, services)
	api.InitRole(r, services)
	api.InitSwagger(r)
	api.InitHealth(r)
}
