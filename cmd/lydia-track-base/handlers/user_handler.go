package handlers

import (
	"github.com/gin-gonic/gin"
	"lydia-track-base/internal/domain"
	"lydia-track-base/internal/service"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) UserHandler {
	return UserHandler{userService: userService}
}

// GetUser godoc
// @Summary Get user by ID
// @Description get the status of server.
// @Tags root
// @Accept */*
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /users:id [get]
func (h UserHandler) GetUser(c *gin.Context) {
	id := c.Param("id")
	user, err := h.userService.GetUser(id)
	if err != nil {
		c.JSON(404, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, user)
}

// CreateUser godoc
// @Summary Create user
// @Description create user.
// @Tags root
// @Accept */*
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /users [post]
func (h UserHandler) CreateUser(c *gin.Context) {
	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	user, err := h.userService.CreateUser(user)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, user)
}

// DeleteUser godoc
// @Summary Delete user
// @Description delete user.
// @Tags root
// @Accept */*
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /users:id [delete]
func (h UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	err := h.userService.DeleteUser(id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": "success"})
}
