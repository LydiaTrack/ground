package handlers

import (
	"github.com/LydiaTrack/ground/pkg/responses"
	"net/http"
	"strconv"

	"github.com/LydiaTrack/ground/internal/service"
	"github.com/LydiaTrack/ground/pkg/auth"
	"github.com/LydiaTrack/ground/pkg/domain/user"
	"github.com/LydiaTrack/ground/pkg/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	exists, err := h.userService.ExistsByUsername(username, auth.CreateAdminAuthContext())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"exists": exists})
}

// CheckEmail
// @Summary Check email
// @Description check email address is exists.
// @Tags root
// @Accept */*
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /users/checkUsername/:username [get]
func (h UserHandler) CheckEmail(c *gin.Context) {
	email := c.Param("email")
	exists, err := h.userService.ExistsByEmail(email, auth.CreateAdminAuthContext())
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

	authContext, err := auth.CreateAuthContext(c, h.authService, &h.userService)
	if err != nil {
		utils.EvaluateError(err, c)
		return
	}

	// Check if pagination parameters are present
	pageStr := c.Query("page")
	limitStr := c.Query("limit")
	searchText := c.DefaultQuery("search", "")

	if pageStr != "" && limitStr != "" {
		// Handle paginated request
		page, err := strconv.Atoi(pageStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page parameter"})
			return
		}

		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
			return
		}

		var userQueryPaginatedResult responses.PaginatedResult[user.Model]
		userQueryPaginatedResult, err = h.userService.QueryPaginated(searchText, page, limit, authContext)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, userQueryPaginatedResult)
	} else {
		// Handle non-paginated request
		var userQueryResult responses.QueryResult[user.Model]
		userQueryResult, err = h.userService.Query(searchText, authContext)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, userQueryResult)
	}
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

	authContext, err := auth.CreateAuthContext(c, h.authService, &h.userService)
	if err != nil {
		utils.EvaluateError(err, c)
		return
	}

	userModel, err := h.userService.Get(id, authContext)
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

	authContext, err := auth.CreateAuthContext(c, h.authService, &h.userService)
	if err != nil {
		utils.EvaluateError(err, c)
		return
	}

	createdUser, err := h.userService.Create(createUserCommand, authContext)
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

	authContext, err := auth.CreateAuthContext(c, h.authService, &h.userService)
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
	err = h.userService.Delete(deleteUserCommand, authContext)
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

	authContext, err := auth.CreateAuthContext(c, h.authService, &h.userService)
	if err != nil {
		utils.EvaluateError(err, c)
		return
	}

	err = h.userService.AddRole(addRoleToUserCommand, authContext)
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

	authContext, err := auth.CreateAuthContext(c, h.authService, &h.userService)
	if err != nil {
		utils.EvaluateError(err, c)
		return
	}

	err = h.userService.RemoveRole(removeRoleFromUserCommand, authContext)
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

	authContext, err := auth.CreateAuthContext(c, h.authService, &h.userService)
	if err != nil {
		utils.EvaluateError(err, c)
		return
	}

	userID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		utils.EvaluateError(err, c)
		return
	}
	result, err := h.userService.GetRolesByUserId(userID, authContext)
	if err != nil {
		utils.EvaluateError(err, c)
		return
	}
	c.JSON(http.StatusOK, result)
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

	authContext, err := auth.CreateAuthContext(c, h.authService, &h.userService)
	if err != nil {
		utils.EvaluateError(err, c)
		return
	}

	userID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		utils.EvaluateError(err, c)
		return
	}
	updatedUser, err := h.userService.Update(userID.Hex(), updateUserCommand, authContext)
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

	authContext, err := auth.CreateAuthContext(c, h.authService, &h.userService)
	if err != nil {
		utils.EvaluateError(err, c)
		return
	}

	updatedUser, err := h.userService.UpdateSelf(updateUserCommand, authContext)
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

	authContext, err := auth.CreateAuthContext(c, h.authService, &h.userService)
	if err != nil {
		utils.EvaluateError(err, c)
		return
	}

	err = h.userService.UpdatePassword(id, updatePasswordCommand, authContext)
	if err != nil {
		utils.EvaluateError(err, c)
		return
	}
	c.Status(http.StatusOK)
}

// UpdateUserSelfPassword godoc
// @Summary Update user self password
// @Description update user self password.
// @Tags root
// @Accept */*
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /users-self/password [put]
func (h UserHandler) UpdateUserSelfPassword(c *gin.Context) {
	var updatePasswordCommand user.UpdatePasswordCommand
	if err := c.ShouldBindJSON(&updatePasswordCommand); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authContext, err := auth.CreateAuthContext(c, h.authService, &h.userService)
	if err != nil {
		utils.EvaluateError(err, c)
		return
	}

	err = h.userService.UpdateSelfPassword(updatePasswordCommand, authContext)
	if err != nil {
		utils.EvaluateError(err, c)
		return
	}
	c.Status(http.StatusOK)
}
