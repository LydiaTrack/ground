package handlers

import (
	"github.com/LydiaTrack/ground/pkg/mongodb/repository"
	"net/http"
	"strconv"

	"github.com/LydiaTrack/ground/internal/service"
	"github.com/LydiaTrack/ground/pkg/auth"
	"github.com/LydiaTrack/ground/pkg/domain/role"
	"github.com/LydiaTrack/ground/pkg/utils"

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

// GetRoles godoc
// @Summary Get all roles
// @Description get all roles.
// @Tags root
// @Accept */*
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /roles [get]
func (h RoleHandler) GetRoles(c *gin.Context) {
	authContext, err := utils.CreateAuthContext(c, h.authService, &h.userService)
	if err != nil {
		utils.EvaluateError(err, c)
		return
	}

	// Check if pagination parameters are present
	pageStr := c.Query("page")
	limitStr := c.Query("limit")
	searchText := c.DefaultQuery("search", "")

	if pageStr != "" && limitStr != "" {
		// Handle paginated query
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

		var roles repository.PaginatedResult[role.Model]
		roles, err = h.roleService.QueryPaginated(searchText, page, limit, authContext)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, roles)
	} else {
		// Handle non-paginated query
		var roles []role.Model
		roles, err = h.roleService.Query(searchText, authContext)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"roles": roles, "count": len(roles)})
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

	authContext, err := utils.CreateAuthContext(c, h.authService, &h.userService)
	if err != nil {
		utils.EvaluateError(err, c)
		return
	}

	getRoleResult, err := h.roleService.GetRole(id, authContext)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, getRoleResult)
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
	var createCmd role.CreateRoleCommand
	if err := c.ShouldBindJSON(&createCmd); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	authContext, err := utils.CreateAuthContext(c, h.authService, &h.userService)
	if err != nil {
		utils.EvaluateError(err, c)
		return
	}

	roleModel, err := h.roleService.CreateRole(createCmd, authContext)
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

	authContext, err := utils.CreateAuthContext(c, h.authService, &h.userService)
	if err != nil {
		utils.EvaluateError(err, c)
		return
	}

	err = h.roleService.DeleteRole(id, authContext)
	if err != nil {
		utils.EvaluateError(err, c)
		return
	}
	c.Status(http.StatusOK)
}
