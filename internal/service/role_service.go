package service

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"

	"github.com/LydiaTrack/lydia-base/internal/permissions"
	"github.com/LydiaTrack/lydia-base/pkg/auth"
	"github.com/LydiaTrack/lydia-base/pkg/constants"
	"github.com/LydiaTrack/lydia-base/pkg/domain/role"
)

type RoleService struct {
	roleRepository RoleRepository
}

func NewRoleService(roleRepository RoleRepository) *RoleService {
	return &RoleService{
		roleRepository: roleRepository,
	}
}

type RoleRepository interface {
	// SaveRole saves a role
	SaveRole(role role.Model) (role.Model, error)
	// GetRole gets a role by id
	GetRole(id primitive.ObjectID) (role.Model, error)
	// GetRoles gets all roles
	GetRoles() ([]role.Model, error)
	// ExistsRole checks if a role exists
	ExistsRole(id primitive.ObjectID) (bool, error)
	// DeleteRole deletes a role by id
	DeleteRole(id primitive.ObjectID) error
	// ExistsByName checks if a role exists by name
	ExistsByName(name string) bool
	// GetRoleByName gets a role by name
	GetRoleByName(name string) (role.Model, error)
	// UpdateRole updates a role
	UpdateRole(id primitive.ObjectID, updateCommand role.UpdateRoleCommand) (role.Model, error)
}

func (s RoleService) CreateRole(command role.CreateRoleCommand, authContext auth.PermissionContext) (role.Model, error) {
	if auth.CheckPermission(authContext.Permissions, permissions.RoleCreatePermission) != nil {
		return role.Model{}, constants.ErrorPermissionDenied
	}

	// Validate role
	// Map command to role
	roleModel := role.NewRole(primitive.NewObjectID().Hex(), command.Name, command.Permissions, command.Tags, command.Info, time.Now(), 1)
	if err := roleModel.Validate(); err != nil {
		return roleModel, err
	}

	roleExists := s.roleRepository.ExistsByName(roleModel.Name)

	if roleExists {
		return role.Model{}, constants.ErrorConflict
	}

	roleModel, err := s.roleRepository.SaveRole(roleModel)
	if err != nil {
		return role.Model{}, constants.ErrorInternalServerError
	}
	return roleModel, nil
}

func (s RoleService) GetRole(id string, authContext auth.PermissionContext) (role.Model, error) {
	if auth.CheckPermission(authContext.Permissions, permissions.RoleReadPermission) != nil {
		return role.Model{}, constants.ErrorPermissionDenied
	}

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return role.Model{}, constants.ErrorBadRequest
	}
	roleModel, err := s.roleRepository.GetRole(objID)
	if err != nil {
		return role.Model{}, err
	}
	return roleModel, nil
}

func (s RoleService) GetRoles(authContext auth.PermissionContext) ([]role.Model, error) {
	if auth.CheckPermission(authContext.Permissions, permissions.RoleReadPermission) != nil {
		return nil, constants.ErrorPermissionDenied
	}

	roles, err := s.roleRepository.GetRoles()
	if err != nil {
		return nil, err
	}
	return roles, nil

}

func (s RoleService) ExistsRole(id string, authContext auth.PermissionContext) (bool, error) {
	if auth.CheckPermission(authContext.Permissions, permissions.RoleReadPermission) != nil {
		return false, constants.ErrorPermissionDenied
	}

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false, constants.ErrorBadRequest
	}
	exists, err := s.roleRepository.ExistsRole(objID)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (s RoleService) DeleteRole(id string, authContext auth.PermissionContext) error {
	if auth.CheckPermission(authContext.Permissions, permissions.RoleDeletePermission) != nil {
		return constants.ErrorPermissionDenied
	}

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return constants.ErrorBadRequest
	}
	err = s.roleRepository.DeleteRole(objID)
	if err != nil {
		return err
	}
	return nil
}

func (s RoleService) ExistsByName(name string, authContext auth.PermissionContext) bool {
	if auth.CheckPermission(authContext.Permissions, permissions.RoleReadPermission) != nil {
		return false
	}
	return s.roleRepository.ExistsByName(name)
}

func (s RoleService) GetRoleByName(name string, authContext auth.PermissionContext) (role.Model, error) {
	if auth.CheckPermission(authContext.Permissions, permissions.RoleReadPermission) != nil {
		return role.Model{}, constants.ErrorPermissionDenied
	}
	roleModel, err := s.roleRepository.GetRoleByName(name)
	if err != nil {
		return role.Model{}, err
	}
	return roleModel, nil
}

func (s RoleService) UpdateRole(id string, command role.UpdateRoleCommand, authContext auth.PermissionContext) (role.Model, error) {
	if auth.CheckPermission(authContext.Permissions, permissions.RoleUpdatePermission) != nil {
		return role.Model{}, constants.ErrorPermissionDenied
	}

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return role.Model{}, constants.ErrorBadRequest
	}

	roleModel, err := s.roleRepository.UpdateRole(objID, command)
	if err != nil {
		return role.Model{}, err
	}
	return roleModel, nil
}
