package repository

import (
	"context"
	"github.com/LydiaTrack/lydia-base/pkg/domain/role"
	"github.com/LydiaTrack/lydia-base/pkg/mongodb"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// A RoleMongoRepository that implements RoleRepository
type RoleMongoRepository struct {
	collection *mongo.Collection
}

var (
	roleRepository *RoleMongoRepository
)

// NewRoleMongoRepository creates a new RoleMongoRepository instance
// which implements RoleRepository
func newRoleMongoRepository() *RoleMongoRepository {

	// FIXME: Burada ileride uzaktaki bir mongodb instance'ına bağlanmak gerekecek
	collection, err := mongodb.GetCollection("roles")
	if err != nil {
		panic(err)
	}

	return &RoleMongoRepository{
		collection: collection,
	}
}

// GetRoleRepository returns a RoleRepository
func GetRoleRepository() *RoleMongoRepository {
	if roleRepository == nil || roleRepository.collection == nil {
		roleRepository = newRoleMongoRepository()
	}
	return roleRepository
}

// SaveRole saves a role
func (r *RoleMongoRepository) SaveRole(roleModel role.Model) (role.Model, error) {
	_, err := r.collection.InsertOne(context.Background(), roleModel)
	if err != nil {
		return role.Model{}, err
	}
	return roleModel, nil
}

// GetRole gets a role by id
func (r *RoleMongoRepository) GetRole(id primitive.ObjectID) (role.Model, error) {
	var roleModel role.Model
	err := r.collection.FindOne(context.Background(), primitive.M{"_id": id}).Decode(&roleModel)
	if err != nil {
		return role.Model{}, err
	}
	return roleModel, nil
}

// GetRoles gets all roles
func (r *RoleMongoRepository) GetRoles() ([]role.Model, error) {
	var roles []role.Model
	cursor, err := r.collection.Find(context.Background(), primitive.M{})
	if err != nil {
		return nil, err
	}
	err = cursor.All(context.Background(), &roles)
	if err != nil {
		return nil, err
	}
	return roles, nil

}

// ExistsRole checks if a role exists by role id
func (r *RoleMongoRepository) ExistsRole(id primitive.ObjectID) (bool, error) {
	var roleModel role.Model
	err := r.collection.FindOne(context.Background(), primitive.M{"_id": id}).Decode(&roleModel)
	if err != nil {
		return false, err
	}
	return true, nil
}

// DeleteRole deletes a role by id
func (r *RoleMongoRepository) DeleteRole(id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(context.Background(), primitive.M{"_id": id})
	if err != nil {
		return err
	}
	return nil
}

// ExistsByName checks if a role exists by role name
func (r *RoleMongoRepository) ExistsByName(name string) bool {
	// Check if role exists by name
	count, err := r.collection.CountDocuments(context.TODO(), primitive.M{"name": name})
	if err != nil {
		panic(err)
	}
	return count > 0
}

func (r *RoleMongoRepository) GetRoleByName(name string) (role.Model, error) {
	var roleModel role.Model
	err := r.collection.FindOne(context.Background(), primitive.M{"name": name}).Decode(&roleModel)
	if err != nil {
		return role.Model{}, err
	}
	return roleModel, nil

}
