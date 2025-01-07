package repository

import (
	"context"

	"github.com/LydiaTrack/ground/pkg/domain/role"
	"github.com/LydiaTrack/ground/pkg/mongodb"
	"github.com/LydiaTrack/ground/pkg/mongodb/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// A RoleMongoRepository that implements RoleRepository
type RoleMongoRepository struct {
	*repository.BaseRepository[role.Model]
}

// GetRoleMongoRepository creates a new RoleMongoRepository instance with the given collection
func GetRoleMongoRepository() *RoleMongoRepository {

	collection, err := mongodb.GetCollection("roles")
	if err != nil {
		panic(err)
	}

	return &RoleMongoRepository{
		BaseRepository: repository.NewBaseRepository[role.Model](collection),
	}
}

// ExistsByName checks if a role exists by role name
func (r *RoleMongoRepository) ExistsByName(name string) bool {
	// Check if role exists by name
	count, err := r.Collection.CountDocuments(context.Background(), bson.M{"name": name})
	if err != nil {
		panic(err)
	}
	return count > 0
}

func (r *RoleMongoRepository) GetRoleByName(name string) (role.Model, error) {
	var roleModel role.Model
	err := r.Collection.FindOne(context.Background(), primitive.M{"name": name}).Decode(&roleModel)
	if err != nil {
		return role.Model{}, err
	}
	return roleModel, nil
}
