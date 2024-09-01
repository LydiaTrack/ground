package handlers

import (
	"github.com/LydiaTrack/lydia-base/internal/service"
	"github.com/LydiaTrack/lydia-base/pkg/auth"
	"github.com/LydiaTrack/lydia-base/pkg/domain/user"
	"github.com/LydiaTrack/lydia-base/pkg/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"os"
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

// CheckUsername
// @Summary Check username
// @Description check username.
// @Tags root
// @Accept */*
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /users/checkUsername/:username [get]
func (h UserHandler) CheckUsername(c *gin.Context) {
	username := c.Param("username")
	exists, err := h.userService.ExistsByUsername(username)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"exists": exists})
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

	authContext, err := utils.CreateAuthContext(c, h.authService, h.userService)
	if err != nil {
		utils.EvaluateError(err, c)
		return
	}

	userModel, err := h.userService.GetUser(id, authContext)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, userModel)
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

	authContext, err := utils.CreateAuthContext(c, h.authService, h.userService)
	if err != nil {
		utils.EvaluateError(err, c)
		return
	}

	createdUser, err := h.userService.CreateUser(createUserCommand, authContext)
	if err != nil {
		utils.EvaluateError(err, c)
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
	id := c.Param("id")

	authContext, err := utils.CreateAuthContext(c, h.authService, h.userService)
	if err != nil {
		utils.EvaluateError(err, c)
		return
	}

	userID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		utils.EvaluateError(err, c)
		return
	}
	deleteUserCommand := user.DeleteUserCommand{
		ID: userID,
	}
	err = h.userService.DeleteUser(deleteUserCommand, authContext)
	if err != nil {
		utils.EvaluateError(err, c)
		return
	}
	c.Status(http.StatusOK)
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

	authContext, err := utils.CreateAuthContext(c, h.authService, h.userService)
	if err != nil {
		utils.EvaluateError(err, c)
		return
	}

	err = h.userService.AddRoleToUser(addRoleToUserCommand, authContext)
	if err != nil {
		utils.EvaluateError(err, c)
		return
	}
	c.Status(http.StatusOK)
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

	authContext, err := utils.CreateAuthContext(c, h.authService, h.userService)
	if err != nil {
		utils.EvaluateError(err, c)
		return
	}

	err = h.userService.RemoveRoleFromUser(removeRoleFromUserCommand, authContext)
	if err != nil {
		utils.EvaluateError(err, c)
	}
	c.Status(http.StatusOK)
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

	authContext, err := utils.CreateAuthContext(c, h.authService, h.userService)
	if err != nil {
		utils.EvaluateError(err, c)
		return
	}

	userID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		utils.EvaluateError(err, c)
		return
	}
	roles, err := h.userService.GetUserRoles(userID, authContext)
	if err != nil {
		utils.EvaluateError(err, c)
		return
	}
	c.JSON(http.StatusOK, roles)
}

// UpdateUser godoc
// @Summary Update user
// @Description update user.
// @Tags root
// @Accept */*
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /users/:id [put]
func (h UserHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var updateUserCommand user.UpdateUserCommand
	if err := c.ShouldBindJSON(&updateUserCommand); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authContext, err := utils.CreateAuthContext(c, h.authService, h.userService)
	if err != nil {
		utils.EvaluateError(err, c)
		return
	}

	userID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		utils.EvaluateError(err, c)
		return
	}
	updatedUser, err := h.userService.UpdateUser(userID.Hex(), updateUserCommand, authContext)
	if err != nil {
		utils.EvaluateError(err, c)
		return
	}
	c.JSON(http.StatusOK, updatedUser)
}

// UpdateUserSelf godoc
// @Summary Update user self
// @Description update user self.
// @Tags root
// @Accept */*
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /users-self [put]
func (h UserHandler) UpdateUserSelf(c *gin.Context) {
	var updateUserCommand user.UpdateUserCommand
	if err := c.ShouldBindJSON(&updateUserCommand); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authContext, err := utils.CreateAuthContext(c, h.authService, h.userService)
	if err != nil {
		utils.EvaluateError(err, c)
		return
	}

	updatedUser, err := h.userService.UpdateUserSelf(updateUserCommand, authContext)
	if err != nil {
		utils.EvaluateError(err, c)
		return
	}
	c.JSON(http.StatusOK, updatedUser)
}

// UpdateUserPassword godoc
// @Summary Update user password
// @Description update user password.
// @Tags root
// @Accept */*
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /users/password:id [put]
func (h UserHandler) UpdateUserPassword(c *gin.Context) {
	id := c.Param("id")
	var updatePasswordCommand user.UpdatePasswordCommand
	if err := c.ShouldBindJSON(&updatePasswordCommand); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authContext, err := utils.CreateAuthContext(c, h.authService, h.userService)
	if err != nil {
		utils.EvaluateError(err, c)
		return
	}

	err = h.userService.UpdateUserPassword(id, updatePasswordCommand, authContext)
	if err != nil {
		utils.EvaluateError(err, c)
		return
	}
	c.Status(http.StatusOK)
}
