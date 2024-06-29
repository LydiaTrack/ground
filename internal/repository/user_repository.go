package repository

import (
	"context"

	"github.com/LydiaTrack/lydia-base/pkg/domain/role"
	"github.com/LydiaTrack/lydia-base/pkg/domain/user"
	"github.com/LydiaTrack/lydia-base/pkg/mongodb"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

// A UserMongoRepository that implements UserRepository
type UserMongoRepository struct {
	collection *mongo.Collection
}

var (
	userRepository *UserMongoRepository
)

// NewUserMongoRepository creates a new UserMongoRepository instance
// which implements UserRepository
func initializeUserRepository() *UserMongoRepository {
	ctx := context.Background()
	// FIXME: Burada ileride uzaktaki bir mongodb instance'ına bağlanmak gerekecek

	collection := mongodb.GetCollection("users", ctx)
	roleRepository = GetRoleRepository()

	return &UserMongoRepository{
		collection: collection,
	}
}

// GetUserRepository returns a UserRepository instance if it is not initialized yet or
// returns the existing one
func GetUserRepository() *UserMongoRepository {
	if userRepository == nil {
		userRepository = initializeUserRepository()
	}
	return userRepository
}

// SaveUser saves a user
func (r *UserMongoRepository) SaveUser(userModel user.Model) (user.Model, error) {
	_, err := r.collection.InsertOne(context.Background(), userModel)
	if err != nil {
		return user.Model{}, err
	}
	return userModel, nil
}

// GetUsers gets all users
func (r *UserMongoRepository) GetUsers() ([]user.Model, error) {
	cursor, err := r.collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	var users []user.Model
	err = cursor.All(context.Background(), &users)
	if err != nil {
		return nil, err
	}
	return users, nil
}

// GetUser gets a user by id
func (r *UserMongoRepository) GetUser(id primitive.ObjectID) (user.Model, error) {
	var userModel user.Model
	err := r.collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&userModel)
	if err != nil {
		return user.Model{}, err
	}
	return userModel, nil
}

// ExistsUser checks if a user exists
func (r *UserMongoRepository) ExistsUser(id primitive.ObjectID) (bool, error) {
	var userModel user.Model
	err := r.collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&userModel)
	if err != nil {
		return false, err
	}
	return true, nil
}

// DeleteUser deletes a user by id
func (r *UserMongoRepository) DeleteUser(id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(context.Background(), bson.M{"_id": id})
	if err != nil {
		return err
	}
	return nil
}

func (r *UserMongoRepository) ExistsByUsernameAndEmail(username string, email string) bool {
	count, err := r.collection.CountDocuments(context.Background(), bson.M{"$or": []bson.M{{"username": username}, {"contactInfo.email": email}}})
	if err != nil {
		return false
	}
	return count > 0
}

func (r *UserMongoRepository) ExistsByUsername(username string) bool {
	count, err := r.collection.CountDocuments(context.Background(), bson.M{"username": username})
	if err != nil {
		return false
	}
	return count > 0
}

func (r *UserMongoRepository) GetUserByUsername(username string) (user.Model, error) {
	var userModel user.Model
	err := r.collection.FindOne(context.Background(), bson.M{"username": username}).Decode(&userModel)
	if err != nil {
		return user.Model{}, err
	}
	return userModel, nil
}

func (r *UserMongoRepository) AddRoleToUser(userID primitive.ObjectID, roleID primitive.ObjectID) error {
	_, err := r.collection.UpdateOne(context.Background(), bson.M{"_id": userID}, bson.M{"$push": bson.M{"roleIds": roleID}})
	if err != nil {
		return err
	}
	return nil
}

func (r *UserMongoRepository) RemoveRoleFromUser(userID primitive.ObjectID, roleID primitive.ObjectID) error {
	_, err := r.collection.UpdateOne(context.Background(), bson.M{"_id": userID}, bson.M{"$pull": bson.M{"roleIds": roleID}})
	if err != nil {
		return err
	}
	return nil
}

func (r *UserMongoRepository) GetUserRoles(userID primitive.ObjectID) ([]role.Model, error) {
	userModel, err := r.GetUser(userID)
	if err != nil {
		return nil, err
	}

	// Resolve roleIds to roles
	var roles []role.Model
	for _, roleID := range userModel.RoleIds {
		roleModel, err := GetRoleRepository().GetRole(roleID)
		if err != nil {
			return nil, err
		}
		roles = append(roles, roleModel)
	}

	return roles, nil
}
