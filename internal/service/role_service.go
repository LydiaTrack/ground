package service

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/LydiaTrack/ground/internal/permissions"
	"github.com/LydiaTrack/ground/pkg/auth"
	"github.com/LydiaTrack/ground/pkg/constants"
	"github.com/LydiaTrack/ground/pkg/domain/role"
	"github.com/LydiaTrack/ground/pkg/mongodb/repository"
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
	repository.Repository[role.Model]
	// ExistsByName checks if a role exists by name
	ExistsByName(name string) bool
	// GetRoleByName gets a role by name
	GetRoleByName(name string) (role.Model, error)
}

func (s RoleService) CreateRole(command role.CreateRoleCommand, authContext auth.PermissionContext) (role.Model, error) {
	if auth.CheckPermission(authContext.Permissions, permissions.RoleCreatePermission) != nil {
		return role.Model{}, constants.ErrorPermissionDenied
	}

	// Validate role
	// Map command to role
	roleModel, err := role.NewRole(
		role.WithName(command.Name),
		role.WithPermissions(command.Permissions),
		role.WithTags(command.Tags),
		role.WithInfo(command.Info),
	)

	if err != nil {
		return role.Model{}, constants.ErrorBadRequest
	}

	if err := roleModel.Validate(); err != nil {
		return role.Model{}, constants.ErrorBadRequest
	}

	roleExists := s.roleRepository.ExistsByName(roleModel.Name)

	if roleExists {
		return role.Model{}, constants.ErrorConflict
	}

	createResult, err := s.roleRepository.Create(context.Background(), *roleModel)
	if err != nil {
		return role.Model{}, constants.ErrorInternalServerError
	}

	insertedId := createResult.InsertedID.(primitive.ObjectID)

	roleAfterSave, err := s.roleRepository.GetByID(context.Background(), insertedId)
	if err != nil {
		return role.Model{}, err
	}
	return roleAfterSave, nil
}

func (s RoleService) GetRole(id string, authContext auth.PermissionContext) (role.Model, error) {
	if auth.CheckPermission(authContext.Permissions, permissions.RoleReadPermission) != nil {
		return role.Model{}, constants.ErrorPermissionDenied
	}

	roleModel, err := s.roleRepository.GetByID(context.Background(), id)
	if err != nil {
		return role.Model{}, err
	}
	return roleModel, nil
}

func (s RoleService) Query(searchText string, authContext auth.PermissionContext) ([]role.Model, error) {
	if auth.CheckPermission(authContext.Permissions, permissions.RoleReadPermission) != nil {
		return nil, constants.ErrorPermissionDenied
	}

	searchFields := []string{"name", "info"}
	roles, err := s.roleRepository.Query(context.Background(), nil, searchFields, searchText)
	if err != nil {
		return nil, err
	}
	return roles, nil
}

func (s RoleService) QueryPaginated(searchText string, page int, limit int, authContext auth.PermissionContext) (repository.PaginatedResult[role.Model], error) {
	if auth.CheckPermission(authContext.Permissions, permissions.RoleReadPermission) != nil {
		return repository.PaginatedResult[role.Model]{}, constants.ErrorPermissionDenied
	}

	searchFields := []string{"name", "info"}
	roles, err := s.roleRepository.QueryPaginate(context.Background(), nil, searchFields, searchText, page, limit, nil)
	if err != nil {
		return repository.PaginatedResult[role.Model]{}, err
	}
	return roles, nil
}

// Exists checks if a user exists by ID
func (s RoleService) Exists(id string, authContext auth.PermissionContext) (bool, error) {
	if auth.CheckPermission(authContext.Permissions, permissions.RoleReadPermission) != nil {
		return false, constants.ErrorPermissionDenied
	}

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false, constants.ErrorBadRequest
	}

	exists, err := s.roleRepository.ExistsByID(context.Background(), objID)
	if err != nil {
		return false, constants.ErrorInternalServerError
	}

	return exists, nil
}

func (s RoleService) DeleteRole(id string, authContext auth.PermissionContext) error {
	if auth.CheckPermission(authContext.Permissions, permissions.RoleDeletePermission) != nil {
		return constants.ErrorPermissionDenied
	}

	_, err := s.roleRepository.Delete(context.Background(), id)
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

	exists, err := s.Exists(id, authContext)
	if err != nil {
		return role.Model{}, constants.ErrorInternalServerError
	}

	if !exists {
		return role.Model{}, constants.ErrorNotFound
	}

	_, err = s.roleRepository.Update(context.Background(), id, command)
	if err != nil {
		return role.Model{}, err
	}

	roleAfterUpdate, err := s.roleRepository.GetByID(context.Background(), id)
	if err != nil {
		return role.Model{}, err
	}
	return roleAfterUpdate, nil
}
