package lydia_track_base

import (
	"log"

	"github.com/Lydia/lydia-base/api"
	"github.com/Lydia/lydia-base/internal/initializers"
	"github.com/Lydia/lydia-base/internal/mongodb"
	"github.com/Lydia/lydia-base/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/kr/pretty"
)

// InitializeLydiaBase initializes the Lydia base server with r as the gin Engine
func InitializeLydiaBase(r *gin.Engine) {
	// Initialize environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// Initialize logging
	utils.InitLogging()
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
