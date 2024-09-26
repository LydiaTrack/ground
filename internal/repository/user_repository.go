package repository

import (
	"context"
	"errors"
	"github.com/LydiaTrack/lydia-base/pkg/domain/role"
	"github.com/LydiaTrack/lydia-base/pkg/domain/user"
	"github.com/LydiaTrack/lydia-base/pkg/mongodb"
	"github.com/LydiaTrack/lydia-base/pkg/responses"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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

	collection, err := mongodb.GetCollection("users")
	if err != nil {
		panic(err)
	}
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

// GetUsers gets all users from the repository
func (r *UserMongoRepository) GetUsers() (responses.QueryResult[user.Model], error) {
	cursor, err := r.collection.Find(context.Background(), primitive.M{})
	if err != nil {
		return responses.QueryResult[user.Model]{}, err
	}

	var users []user.Model
	err = cursor.All(context.Background(), &users)
	if err != nil {
		return responses.QueryResult[user.Model]{}, err
	}

	// Return the QueryResult by value
	return *responses.NewQueryResult(len(users), users), nil
}

// GetUser gets a user by id
func (r *UserMongoRepository) GetUser(id primitive.ObjectID) (user.Model, error) {
	var userModel user.Model
	err := r.collection.FindOne(context.Background(), primitive.M{"_id": id}).Decode(&userModel)
	if err != nil {
		return user.Model{}, err
	}
	return userModel, nil
}

// ExistsUser checks if a user exists
func (r *UserMongoRepository) ExistsUser(id primitive.ObjectID) (bool, error) {
	var userModel user.Model
	err := r.collection.FindOne(context.Background(), primitive.M{"_id": id}).Decode(&userModel)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false, nil // No documents found, so user does not exist
		}
		return false, err // An actual error occurred
	}
	return true, nil // User exists
}

// DeleteUser deletes a user by id
func (r *UserMongoRepository) DeleteUser(id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(context.Background(), primitive.M{"_id": id})
	if err != nil {
		return err
	}
	return nil
}

func (r *UserMongoRepository) ExistsByUsernameAndEmail(username string, email string) bool {
	count, err := r.collection.CountDocuments(context.Background(), primitive.M{"$or": []primitive.M{{"username": username}, {"contactInfo.email": email}}})
	if err != nil {
		return false
	}
	return count > 0
}

func (r *UserMongoRepository) ExistsByUsername(username string) bool {
	count, err := r.collection.CountDocuments(context.Background(), primitive.M{"username": username})
	if err != nil {
		return false
	}
	return count > 0
}

func (r *UserMongoRepository) ExistsByEmail(email string) bool {
	count, err := r.collection.CountDocuments(context.Background(), primitive.M{"contactInfo.email": email})
	if err != nil {
		return false
	}
	return count > 0
}

func (r *UserMongoRepository) GetUserByUsername(username string) (user.Model, error) {
	var userModel user.Model
	err := r.collection.FindOne(context.Background(), primitive.M{"username": username}).Decode(&userModel)
	if err != nil {
		return user.Model{}, err
	}
	return userModel, nil
}

func (r *UserMongoRepository) AddRoleToUser(userID primitive.ObjectID, roleID primitive.ObjectID) error {
	// roleIds can be null or empty
	_, err := r.collection.UpdateOne(context.Background(), primitive.M{"_id": userID}, primitive.M{"$addToSet": primitive.M{"roleIds": roleID}})
	if err != nil {
		return err
	}
	return nil
}

func (r *UserMongoRepository) RemoveRoleFromUser(userID primitive.ObjectID, roleID primitive.ObjectID) error {
	_, err := r.collection.UpdateOne(context.Background(), primitive.M{"_id": userID}, primitive.M{"$pull": primitive.M{"roleIds": roleID}})
	if err != nil {
		return err
	}
	return nil
}

func (r *UserMongoRepository) GetUserRoles(userID primitive.ObjectID) (responses.QueryResult[role.Model], error) {
	userModel, err := r.GetUser(userID)
	if err != nil {
		return responses.QueryResult[role.Model]{}, err
	}

	// Resolve roleIds to roles
	var roles []role.Model
	for _, roleID := range *userModel.RoleIds {
		roleModel, err := GetRoleRepository().GetRole(roleID)
		if err != nil {
			return responses.QueryResult[role.Model]{}, err
		}
		roles = append(roles, roleModel)
	}

	// Create the QueryResult with the roles and total count
	return *responses.NewQueryResult(len(roles), roles), nil
}

func (r *UserMongoRepository) UpdateUser(id primitive.ObjectID, updateCommand user.UpdateUserCommand) (user.Model, error) {
	_, err := r.collection.UpdateOne(context.Background(), primitive.M{"_id": id}, primitive.M{"$set": updateCommand})
	if err != nil {
		return user.Model{}, err
	}
	return r.GetUser(id)
}

func (r *UserMongoRepository) UpdateUserPassword(id primitive.ObjectID, password string) error {
	_, err := r.collection.UpdateOne(context.Background(), primitive.M{"_id": id}, primitive.M{"$set": primitive.M{"password": password}})
	if err != nil {
		return err
	}
	return nil
}

func (r *UserMongoRepository) GetUserByEmailAddress(email string) (user.Model, error) {
	var userModel user.Model
	err := r.collection.FindOne(context.Background(), primitive.M{"contactInfo.email": email}).Decode(&userModel)
	if err != nil {
		return user.Model{}, err
	}
	return userModel, nil
}
