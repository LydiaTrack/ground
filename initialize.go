package lydia_base

import (
	"reflect"
	"time"

	"github.com/LydiaTrack/ground/internal/permissions"
	"github.com/LydiaTrack/ground/internal/provider"
	"github.com/LydiaTrack/ground/pkg/auth"
	"github.com/LydiaTrack/ground/pkg/domain/role"
	"github.com/LydiaTrack/ground/pkg/manager"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/LydiaTrack/ground/internal/api"
	"github.com/LydiaTrack/ground/internal/blocker"
	"github.com/LydiaTrack/ground/internal/initializers"
	"github.com/LydiaTrack/ground/internal/log"
	"github.com/LydiaTrack/ground/pkg/middlewares"
	"github.com/LydiaTrack/ground/pkg/mongodb"
	"github.com/LydiaTrack/ground/pkg/service_initializer"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// Initialize initializes the Lydia base server with r as the gin Engine
func Initialize(r *gin.Engine) {
	// Initialize environment variables
	err := godotenv.Load()
	if err != nil {
		log.LogFatal("Error loading .env file")
	}
	// Initialize logging
	log.InitLogging()

	// Initialize metrics
	initMetrics(r)

	// Initialize IP Blocker
	blocker.Initialize()

	// Initialize database connection
	err = mongodb.InitializeMongoDBConnection()
	if err != nil {
		log.LogFatal("Error initializing MongoDB connection")
	}

	// Initialize services
	service_initializer.InitializeServices()

	// Initialize API routes
	initializeRoutes(r, service_initializer.GetServices())

	// Initialize default role
	err = initializers.InitializeDefaultRole()
	if err != nil {
		log.LogFatal("Error initializing default role")
	}

	// Create default roles
	createDefaultRoles()

	// Register the self role provider
	manager.RegisterRoleProvider(provider.SelfRoleProvider{})

	// Initialize default user
	err = initializers.InitializeDefaultUser()
	if err != nil {
		log.LogFatal("Error initializing default user: " + err.Error())
		panic(err)
	}

}

// initializeRoutes initializes routes for each API
func initializeRoutes(r *gin.Engine, services service_initializer.Services) {
	globalInterceptors := []gin.HandlerFunc{gin.Recovery(), gin.Logger()}

	r.Use(globalInterceptors...)
	r.Use(middlewares.IPBlockMiddleware())

	api.InitAuth(r, services)
	api.InitUser(r, services)
	api.InitRole(r, services)
	api.InitResetPassword(r, services)
	api.InitFeedback(r, services)
	api.InitSwagger(r)
	api.InitHealth(r)

	go func() {
		for {
			time.Sleep(30 * time.Second)
			blocker.GlobalBlocker.RemoveExpired()
		}
	}()
}

func createDefaultRoles() {
	authContext := auth.PermissionContext{
		Permissions: []auth.Permission{auth.AdminPermission},
		UserId:      nil,
	}
	selfServiceRoleCmd := role.CreateRoleCommand{
		Name: "Lydia Self Service Role",
		Tags: []string{"self-service"},
		Info: "This role is for the users who can manage their profiles",
		Permissions: []auth.Permission{
			permissions.UserSelfUpdatePermission,
			permissions.UserSelfGetPermission,
		},
	}

	isExists := service_initializer.GetServices().RoleService.ExistsByName(selfServiceRoleCmd.Name, authContext)
	if isExists {
		// If exists, check if the permissions are the same
		roleModel, err := service_initializer.GetServices().RoleService.GetRoleByName(selfServiceRoleCmd.Name, authContext)
		if err != nil {
			log.LogFatal("Error creating default roles: " + err.Error())
			return
		}
		currentRolePermissions := roleModel.Permissions
		// Compare the permissions
		isSamePermissions := reflect.DeepEqual(currentRolePermissions, selfServiceRoleCmd.Permissions)
		if !isSamePermissions {
			log.Log("Permissions are different for the role: " + selfServiceRoleCmd.Name)
			// Update the role
			cmd := role.UpdateRoleCommand{
				Name:        roleModel.Name,
				Info:        roleModel.Info,
				Tags:        roleModel.Tags,
				Permissions: selfServiceRoleCmd.Permissions,
			}
			log.Log("Updating the role's permissions: " + roleModel.Name)
			_, err := service_initializer.GetServices().RoleService.UpdateRole(roleModel.ID.Hex(), cmd, authContext)
			if err != nil {
				log.LogFatal("Error creating default roles: " + err.Error())
				return
			}

		}
		return
	}

	// Create the role if it does not exist
	_, err := service_initializer.GetServices().RoleService.CreateRole(selfServiceRoleCmd, authContext)
	if err != nil {
		log.LogFatal("Error creating default roles: " + err.Error())
		return
	}
}

func initMetrics(r *gin.Engine) {
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
}
