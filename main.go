package main

import (
	"github.com/LydiaTrack/lydia-base/helper"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	// Initialize the Lydia Base
	helper.Initialize(r)
	r.Run(":8080")
}
