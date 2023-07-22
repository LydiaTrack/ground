package service

import (
	"errors"
	"lydia-track-base/internal/domain"
	"lydia-track-base/internal/domain/commands"
	"time"

	"gopkg.in/mgo.v2/bson"
)

type RoleService struct {
	roleRepository RoleRepository
}

func NewRoleService(roleRepository RoleRepository) RoleService {
	return RoleService{
		roleRepository: roleRepository,
	}
}

// CreateRole TODO: Add permission check
func (s RoleService) CreateRole(command commands.CreateRoleCommand) (domain.RoleModel, error) {
	// TODO: These kind of operations must be done with specific requests, not by RoleModel model itself
	// Validate role
	// Map command to role
	role := domain.NewRole(bson.NewObjectId().Hex(), command.Name, command.Tags, command.Info, time.Now(), 1)
	if err := role.Validate(); err != nil {
		return role, err
	}

	roleExists := s.roleRepository.ExistsByRolename(role.Name)

	if roleExists {
		return domain.RoleModel{}, errors.New("role already exists")
	}

	role, err := s.roleRepository.SaveRole(role)
	if err != nil {
		return domain.RoleModel{}, err
	}
	return role, nil
}

func (s RoleService) GetRole(id string) (domain.RoleModel, error) {
	role, err := s.roleRepository.GetRole(bson.ObjectIdHex(id))
	if err != nil {
		return domain.RoleModel{}, err
	}
	return role, nil
}

func (s RoleService) ExistsRole(id string) (bool, error) {
	exists, err := s.roleRepository.ExistsRole(bson.ObjectIdHex(id))
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (s RoleService) DeleteRole(id string) error {
	err := s.roleRepository.DeleteRole(bson.ObjectIdHex(id))
	if err != nil {
		return err
	}
	return nil
}

type RoleRepository interface {
	// SaveRole saves a role
	SaveRole(role domain.RoleModel) (domain.RoleModel, error)
	// GetRole gets a role by id
	GetRole(id bson.ObjectId) (domain.RoleModel, error)
	// ExistsRole checks if a role exists
	ExistsRole(id bson.ObjectId) (bool, error)
	// DeleteRole deletes a role by id
	DeleteRole(id bson.ObjectId) error
	// ExistsByRolename checks if a role exists by rolename
	ExistsByRolename(rolename string) bool
}
