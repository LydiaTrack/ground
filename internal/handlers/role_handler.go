package handlers

import (
	"github.com/LydiaTrack/lydia-base/auth"
	"github.com/LydiaTrack/lydia-base/internal/domain/role"
	"github.com/LydiaTrack/lydia-base/internal/service"
	"github.com/LydiaTrack/lydia-base/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RoleHandler struct {
	roleService service.RoleService
	authService auth.Service
	userService service.UserService
}

func NewRoleHandler(roleService service.RoleService, authService auth.Service, userService service.UserService) RoleHandler {
	return RoleHandler{
		roleService: roleService,
		authService: authService,
		userService: userService,
	}
}

// GetRole godoc
// @Summary Get role by ID
// @Description get the status of server.
// @Tags root
// @Accept */*
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /roles:id [get]
func (h RoleHandler) GetRole(c *gin.Context) {
	id := c.Param("id")

	authContext, err := utils.CreateAuthContext(c, h.authService, h.userService)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	role, err := h.roleService.GetRole(id, authContext)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, role)
}

// CreateRole godoc
// @Summary Create role
// @Description create role.
// @Tags root
// @Accept */*
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /roles [post]
func (h RoleHandler) CreateRole(c *gin.Context) {
	var role role.CreateRoleCommand
	if err := c.ShouldBindJSON(&role); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	authContext, err := utils.CreateAuthContext(c, h.authService, h.userService)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	roleModel, err := h.roleService.CreateRole(role, authContext)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, roleModel)
}

// DeleteRole godoc
// @Summary Delete role
// @Description delete role.
// @Tags root
// @Accept */*
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /roles:id [delete]
func (h RoleHandler) DeleteRole(c *gin.Context) {
	id := c.Param("id")

	authContext, err := utils.CreateAuthContext(c, h.authService, h.userService)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = h.roleService.DeleteRole(id, authContext)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
