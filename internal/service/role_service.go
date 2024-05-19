package service

import (
	"errors"
	"github.com/LydiaTrack/lydia-base/auth"
	"github.com/LydiaTrack/lydia-base/internal/domain/role"
	"github.com/LydiaTrack/lydia-base/internal/permissions"
	"time"

	"gopkg.in/mgo.v2/bson"
)

type RoleService struct {
	roleRepository RoleRepository
}

func NewRoleService(roleRepository RoleRepository) *RoleService {
	return &RoleService{
		roleRepository: roleRepository,
	}
}

func (s RoleService) CreateRole(command role.CreateRoleCommand, authContext auth.PermissionContext) (role.Model, error) {
	if auth.CheckPermission(authContext.Permissions, permissions.RoleCreatePermission) != nil {
		return role.Model{}, errors.New("not permitted")
	}

	// Validate role
	// Map command to role
	roleModel := role.NewRole(bson.NewObjectId().Hex(), command.Name, command.Permissions, command.Tags, command.Info, time.Now(), 1)
	if err := roleModel.Validate(); err != nil {
		return roleModel, err
	}

	roleExists := s.roleRepository.ExistsByRolename(roleModel.Name)

	if roleExists {
		return role.Model{}, errors.New("role already exists")
	}

	roleModel, err := s.roleRepository.SaveRole(roleModel)
	if err != nil {
		return role.Model{}, err
	}
	return roleModel, nil
}

func (s RoleService) GetRole(id string, authContext auth.PermissionContext) (role.Model, error) {
	if auth.CheckPermission(authContext.Permissions, permissions.RoleReadPermission) != nil {
		return role.Model{}, errors.New("not permitted")
	}

	roleModel, err := s.roleRepository.GetRole(bson.ObjectIdHex(id))
	if err != nil {
		return role.Model{}, err
	}
	return roleModel, nil
}

func (s RoleService) ExistsRole(id string, authContext auth.PermissionContext) (bool, error) {
	if auth.CheckPermission(authContext.Permissions, permissions.RoleReadPermission) != nil {
		return false, errors.New("not permitted")
	}

	exists, err := s.roleRepository.ExistsRole(bson.ObjectIdHex(id))
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (s RoleService) DeleteRole(id string, authContext auth.PermissionContext) error {
	if auth.CheckPermission(authContext.Permissions, permissions.RoleDeletePermission) != nil {
		return errors.New("not permitted")
	}

	err := s.roleRepository.DeleteRole(bson.ObjectIdHex(id))
	if err != nil {
		return err
	}
	return nil
}

func (s RoleService) ExistsByRolename(rolename string, authContext auth.PermissionContext) bool {
	if auth.CheckPermission(authContext.Permissions, permissions.RoleReadPermission) != nil {
		return false
	}
	return s.roleRepository.ExistsByRolename(rolename)
}

func (s RoleService) GetRoleByRolename(rolename string, authContext auth.PermissionContext) (role.Model, error) {
	if auth.CheckPermission(authContext.Permissions, permissions.RoleReadPermission) != nil {
		return role.Model{}, errors.New("not permitted")
	}
	roleModel, err := s.roleRepository.GetRoleByRolename(rolename)
	if err != nil {
		return role.Model{}, err
	}
	return roleModel, nil
}

type RoleRepository interface {
	// SaveRole saves a role
	SaveRole(role role.Model) (role.Model, error)
	// GetRole gets a role by id
	GetRole(id bson.ObjectId) (role.Model, error)
	// ExistsRole checks if a role exists
	ExistsRole(id bson.ObjectId) (bool, error)
	// DeleteRole deletes a role by id
	DeleteRole(id bson.ObjectId) error
	// ExistsByRolename checks if a role exists by rolename
	ExistsByRolename(rolename string) bool
	// GetRoleByRolename gets a role by rolename
	GetRoleByRolename(rolename string) (role.Model, error)
}
