package main

import (
	"github.com/gin-gonic/gin"
	"lydia-track-base/cmd/lydia-track-base/api"
	"lydia-track-base/internal/rest"
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
	api.InitUser(r)
	api.InitSwagger(r)

	wrapper := rest.InitEndpointWrapper([]gin.HandlerFunc{}, map[string][]gin.HandlerFunc{})
	wrapper.WrapEngine(r)
}

// Client -> Server -> Wrapper -> Interceptor -> Interceptor -> API -> Service -> Repository -> DB
