package handlers

import (
	"net/http"
	"os"

	"github.com/LydiaTrack/lydia-base/auth"
	"github.com/LydiaTrack/lydia-base/auth_utils"

	"github.com/LydiaTrack/lydia-base/internal/domain/user"
	"github.com/LydiaTrack/lydia-base/internal/service"

	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
)

type UserHandler struct {
	userService service.UserService
	authService auth.Service
}

func NewUserHandler(userService service.UserService, authService auth.Service) UserHandler {
	return UserHandler{
		userService: userService,
		authService: authService,
	}
}

// GetUsers godoc
// @Summary Get users
// @Description get users.
// @Tags root
// @Accept */*
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /users [get]
func (h UserHandler) GetUsers(c *gin.Context) {
	currentUser, err := h.authService.GetCurrentUser(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	currentUserPermissions, err := h.userService.GetUserPermissionList(currentUser.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var users []user.Model
	if currentUser.Username == os.Getenv("DEFAULT_USER_USERNAME") {
		users, err = h.userService.GetUsers(auth.PermissionContext{
			Permissions: []auth.Permission{auth.AdminPermission},
			UserId:      &currentUser.ID,
		})
	} else {
		users, err = h.userService.GetUsers(auth.PermissionContext{
			Permissions: currentUserPermissions,
			UserId:      &currentUser.ID,
		})
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)

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

	authContext, err := auth_utils.CreateAuthContext(c, h.authService, h.userService)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.GetUser(id, authContext)
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
	var createUserCommand user.CreateUserCommand
	if err := c.ShouldBindJSON(&createUserCommand); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authContext, err := auth_utils.CreateAuthContext(c, h.authService, h.userService)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	createdUser, err := h.userService.CreateUser(createUserCommand, authContext)
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
	var deleteUserCommand user.DeleteUserCommand
	if err := c.ShouldBindJSON(&deleteUserCommand); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authContext, err := auth_utils.CreateAuthContext(c, h.authService, h.userService)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = h.userService.DeleteUser(deleteUserCommand, authContext)
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
	var addRoleToUserCommand user.AddRoleToUserCommand
	if err := c.ShouldBindJSON(&addRoleToUserCommand); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authContext, err := auth_utils.CreateAuthContext(c, h.authService, h.userService)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = h.userService.AddRoleToUser(addRoleToUserCommand, authContext)
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
	var removeRoleFromUserCommand user.RemoveRoleFromUserCommand
	if err := c.ShouldBindJSON(&removeRoleFromUserCommand); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	authContext, err := auth_utils.CreateAuthContext(c, h.authService, h.userService)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = h.userService.RemoveRoleFromUser(removeRoleFromUserCommand, authContext)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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

	authContext, err := auth_utils.CreateAuthContext(c, h.authService, h.userService)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	roles, err := h.userService.GetUserRoles(bson.ObjectIdHex(id), authContext)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, roles)
}
