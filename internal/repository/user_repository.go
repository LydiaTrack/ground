package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
	"lydia-track-base/internal/domain/role"
	"lydia-track-base/internal/domain/user"
	"lydia-track-base/internal/mongodb"
	"os"
)

// A UserMongoRepository that implements UserRepository
type UserMongoRepository struct {
	client     *mongo.Client
	collection *mongo.Collection
}

var (
	userRepository *UserMongoRepository
)

// NewUserMongoRepository creates a new UserMongoRepository instance
// which implements UserRepository
func newUserMongoRepository() *UserMongoRepository {
	ctx := context.Background()
	// FIXME: Burada ileride uzaktaki bir mongodb instance'ına bağlanmak gerekecek
	// FIXME: Ortaklaştırılacak
	container := mongodb.GetContainer()

	host, err := container.Host(ctx)
	if err != nil {
		return nil
	}

	port, err := container.MappedPort(ctx, "27017")
	if err != nil {
		return nil
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://"+host+":"+port.Port()))
	if err != nil {
		return nil
	}

	collection := client.Database(os.Getenv("LYDIA_DB_NAME")).Collection("users")

	return &UserMongoRepository{
		client:     client,
		collection: collection,
	}
}

// GetUserRepository returns a UserRepository instance if it is not initialized yet or
// returns the existing one
func GetUserRepository() *UserMongoRepository {
	if userRepository == nil {
		userRepository = newUserMongoRepository()
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

// GetUser gets a user by id
func (r *UserMongoRepository) GetUser(id bson.ObjectId) (user.Model, error) {
	var userModel user.Model
	err := r.collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&userModel)
	if err != nil {
		return user.Model{}, err
	}
	return userModel, nil
}

// ExistsUser checks if a user exists
func (r *UserMongoRepository) ExistsUser(id bson.ObjectId) (bool, error) {
	var userModel user.Model
	err := r.collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&userModel)
	if err != nil {
		return false, err
	}
	return true, nil
}

// DeleteUser deletes a user by id
func (r *UserMongoRepository) DeleteUser(id bson.ObjectId) error {
	_, err := r.collection.DeleteOne(context.Background(), bson.M{"_id": id})
	if err != nil {
		return err
	}
	return nil
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

func (r *UserMongoRepository) AddRoleToUser(userID bson.ObjectId, roleID bson.ObjectId) error {
	_, err := r.collection.UpdateOne(context.Background(), bson.M{"_id": userID}, bson.M{"$addToSet": bson.M{"roles": roleID}})
	if err != nil {
		return err
	}
	return nil
}

func (r *UserMongoRepository) RemoveRoleFromUser(userID bson.ObjectId, roleID bson.ObjectId) error {
	_, err := r.collection.UpdateOne(context.Background(), bson.M{"_id": userID}, bson.M{"$pull": bson.M{"roles": roleID}})
	if err != nil {
		return err
	}
	return nil
}

func (r *UserMongoRepository) GetUserRoles(userID bson.ObjectId) ([]role.Model, error) {
	userModel, err := r.GetUser(userID)
	if err != nil {
		return nil, err
	}

	return userModel.Roles, nil
}
