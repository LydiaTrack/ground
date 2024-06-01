package service

import (
	"errors"
	"github.com/LydiaTrack/lydia-base/internal/permissions"
	"github.com/LydiaTrack/lydia-base/pkg/auth"
	"github.com/LydiaTrack/lydia-base/pkg/domain/role"
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

	roleExists := s.roleRepository.ExistsByName(roleModel.Name)

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

func (s RoleService) GetRoles(authContext auth.PermissionContext) ([]role.Model, error) {
	if auth.CheckPermission(authContext.Permissions, permissions.RoleReadPermission) != nil {
		return nil, errors.New("not permitted")
	}

	roles, err := s.roleRepository.GetRoles()
	if err != nil {
		return nil, err
	}
	return roles, nil

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

func (s RoleService) ExistsByName(name string, authContext auth.PermissionContext) bool {
	if auth.CheckPermission(authContext.Permissions, permissions.RoleReadPermission) != nil {
		return false
	}
	return s.roleRepository.ExistsByName(name)
}

func (s RoleService) GetRoleByName(name string, authContext auth.PermissionContext) (role.Model, error) {
	if auth.CheckPermission(authContext.Permissions, permissions.RoleReadPermission) != nil {
		return role.Model{}, errors.New("not permitted")
	}
	roleModel, err := s.roleRepository.GetRoleByName(name)
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
	// GetRoles gets all roles
	GetRoles() ([]role.Model, error)
	// ExistsRole checks if a role exists
	ExistsRole(id bson.ObjectId) (bool, error)
	// DeleteRole deletes a role by id
	DeleteRole(id bson.ObjectId) error
	// ExistsByName checks if a role exists by name
	ExistsByName(name string) bool
	// GetRoleByName gets a role by name
	GetRoleByName(name string) (role.Model, error)
}
