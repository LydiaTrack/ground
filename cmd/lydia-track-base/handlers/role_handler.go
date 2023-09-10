package handlers

import (
	"lydia-track-base/internal/domain/role/commands"
	"lydia-track-base/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RoleHandler struct {
	roleService service.RoleService
}

func NewRoleHandler(roleService service.RoleService) RoleHandler {
	return RoleHandler{roleService: roleService}
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
	role, err := h.roleService.GetRole(id)
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
	var role commands.CreateRoleCommand
	if err := c.ShouldBindJSON(&role); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	roleModel, err := h.roleService.CreateRole(role)
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
	err := h.roleService.DeleteRole(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
