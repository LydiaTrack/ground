package repository

import (
	"context"
	"github.com/LydiaTrack/lydia-base/internal/domain/role"
	"github.com/LydiaTrack/lydia-base/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
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
	ctx := context.Background()
	// FIXME: Burada ileride uzaktaki bir mongodb instance'ına bağlanmak gerekecek
	collection := mongodb.GetCollection("roles", ctx)

	return &RoleMongoRepository{
		collection: collection,
	}
}

// GetRoleRepository returns a RoleRepository
func GetRoleRepository() *RoleMongoRepository {
	if roleRepository == nil {
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
func (r *RoleMongoRepository) GetRole(id bson.ObjectId) (role.Model, error) {
	var roleModel role.Model
	err := r.collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&roleModel)
	if err != nil {
		return role.Model{}, err
	}
	return roleModel, nil
}

// ExistsRole checks if a role exists by role id
func (r *RoleMongoRepository) ExistsRole(id bson.ObjectId) (bool, error) {
	var roleModel role.Model
	err := r.collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&roleModel)
	if err != nil {
		return false, err
	}
	return true, nil
}

// DeleteRole deletes a role by id
func (r *RoleMongoRepository) DeleteRole(id bson.ObjectId) error {
	_, err := r.collection.DeleteOne(context.Background(), bson.M{"_id": id})
	if err != nil {
		return err
	}
	return nil
}

// ExistsByRolename checks if a role exists by role name
func (r *RoleMongoRepository) ExistsByRolename(rolename string) bool {
	// Check if role exists by name
	count, err := r.collection.CountDocuments(context.Background(), bson.M{"name": rolename})
	if err != nil {
		return false
	}
	return count > 0
}
