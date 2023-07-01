package main

import (
	"github.com/gin-gonic/gin"
	"lydia-track-base/cmd/lydia-track-base/api"
)

// @title Lydia Track Base API
// @version 0.0.1
// @description Lydia Track Base API

// @host localhost:8080
// @BasePath /
// @schemes http
func main() {

	r := gin.New()

	// Initialize routes
	initializeRoutes(r)

	// Run server on port 8080
	r.Run(":8080")

}

// initializeRoutes initializes routes for each API
func initializeRoutes(r *gin.Engine) {
	globalInterceptors := []gin.HandlerFunc{gin.Recovery(), gin.Logger()}

	r.Use(globalInterceptors...)

	api.InitUser(r)
	api.InitSwagger(r)
	api.InitAuth(r)
}
