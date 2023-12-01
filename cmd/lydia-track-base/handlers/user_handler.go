package handlers

import (
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
	"lydia-track-base/internal/domain/user/commands"
	"lydia-track-base/internal/service"
	"net/http"
)

type UserHandler struct {
	userService service.UserService
	authService service.Service
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
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
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
	var createUserCommand commands.CreateUserCommand
	if err := c.ShouldBindJSON(&createUserCommand); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: Extract to a function
	currentUser, _ := h.authService.GetCurrentUser(c)
	currentUserPermissions, _ := h.userService.GetUserPermissions(currentUser.ID)

	createdUser, err := h.userService.CreateUser(createUserCommand, currentUserPermissions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, createdUser)
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
	var deleteUserCommand commands.DeleteUserCommand
	if err := c.ShouldBindJSON(&deleteUserCommand); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := h.userService.DeleteUser(deleteUserCommand)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

// AddRoleToUser godoc
// @Summary Add role to user
// @Description add role to user.
// @Tags root
// @Accept */*
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /users/roles [post]
func (h UserHandler) AddRoleToUser(c *gin.Context) {
	var addRoleToUserCommand commands.AddRoleToUserCommand
	if err := c.ShouldBindJSON(&addRoleToUserCommand); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := h.userService.AddRoleToUser(addRoleToUserCommand)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

// RemoveRoleFromUser godoc
// @Summary Remove role from user
// @Description remove role from user.
// @Tags root
// @Accept */*
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /users/roles [delete]
func (h UserHandler) RemoveRoleFromUser(c *gin.Context) {
	var removeRoleFromUserCommand commands.RemoveRoleFromUserCommand
	if err := c.ShouldBindJSON(&removeRoleFromUserCommand); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := h.userService.RemoveRoleFromUser(removeRoleFromUserCommand)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

// GetUserRoles godoc
// @Summary Get user roles
// @Description get user roles.
// @Tags root
// @Accept */*
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /users/roles:id [get]
func (h UserHandler) GetUserRoles(c *gin.Context) {
	id := c.Param("id")
	roles, err := h.userService.GetUserRoles(bson.ObjectIdHex(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, roles)
}
