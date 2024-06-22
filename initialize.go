package lydia_base

import (
	"github.com/LydiaTrack/lydia-base/internal/api"
	"github.com/LydiaTrack/lydia-base/internal/blocker"
	"github.com/LydiaTrack/lydia-base/internal/initializers"
	"github.com/LydiaTrack/lydia-base/internal/log"
	"github.com/LydiaTrack/lydia-base/pkg/middlewares"
	"github.com/LydiaTrack/lydia-base/pkg/mongodb"
	"github.com/LydiaTrack/lydia-base/pkg/service_initializer"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"time"
)

// Initialize initializes the Lydia base server with r as the gin Engine
func Initialize(r *gin.Engine) {
	// Initialize environment variables
	err := godotenv.Load()
	if err != nil {
		log.LogFatal("Error loading .env file")
	}
	// Initialize logging
	log.InitLogging()

	// Initialize IP Blocker
	blocker.Initialize()

	// Initialize database connection
	err = mongodb.InitializeMongoDBConnection()
	if err != nil {
		log.LogFatal("Error initializing MongoDB connection")
	}

	// Initialize services
	service_initializer.InitializeServices()

	// Initialize API routes
	initializeRoutes(r, service_initializer.GetServices())

	// Initialize default role
	err = initializers.InitializeDefaultRole()
	if err != nil {
		log.LogFatal("Error initializing default user")
	}

	// Initialize default user
	err = initializers.InitializeDefaultUser()
	if err != nil {
		log.LogFatal("Error initializing default user")
		panic(err)
	}

}

// initializeRoutes initializes routes for each API
func initializeRoutes(r *gin.Engine, services service_initializer.Services) {
	globalInterceptors := []gin.HandlerFunc{gin.Recovery(), gin.Logger()}

	r.Use(globalInterceptors...)
	r.Use(middlewares.IPBlockMiddleware())

	api.InitAuth(r, services)
	api.InitUser(r, services)
	api.InitRole(r, services)
	api.InitSwagger(r)
	api.InitHealth(r)

	go func() {
		for {
			time.Sleep(30 * time.Second)
			blocker.GlobalBlocker.RemoveExpired()
		}
	}()
}
