package repository

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
	"log"
	"lydia-track-base/internal/domain"
	"lydia-track-base/internal/mongodb"
	"os"
)

// A UserMongoRepository that implements UserRepository
type UserMongoRepository struct {
	client     *mongo.Client
	collection *mongo.Collection
}

// NewUserMongoRepository creates a new UserMongoRepository instance
// which implements UserRepository
func NewUserMongoRepository() *UserMongoRepository {
	ctx := context.Background()
	// FIXME: Burada ileride uzaktaki bir mongodb instance'ına bağlanmak gerekecek
	container, err := mongodb.StartContainer(ctx)
	if err != nil {
		return nil
	}

	endpoint, err := container.Endpoint(ctx, "mongodb")
	if err != nil {
		log.Fatal(fmt.Errorf("failed to get endpoint: %w", err))
	}

	mongoClient, err := mongo.NewClient(options.Client().ApplyURI(endpoint))
	if err != nil {
		log.Fatal(fmt.Errorf("error creating mongo client: %w", err))
	}

	err = mongoClient.Connect(ctx)
	if err != nil {
		log.Fatal("Error connecting to mongo: ", err)
	}

	err = godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	collection := mongoClient.Database(os.Getenv("LYDIA_DB_NAME")).Collection("users")

	return &UserMongoRepository{
		client:     mongoClient,
		collection: collection,
	}
}

// SaveUser saves a user
func (r *UserMongoRepository) SaveUser(user domain.UserModel) (domain.UserModel, error) {
	_, err := r.collection.InsertOne(context.Background(), user)
	if err != nil {
		return domain.UserModel{}, err
	}
	return user, nil
}

// GetUser gets a user by id
func (r *UserMongoRepository) GetUser(id bson.ObjectId) (domain.UserModel, error) {
	var user domain.UserModel
	err := r.collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return domain.UserModel{}, err
	}
	return user, nil
}

// ExistsUser checks if a user exists
func (r *UserMongoRepository) ExistsUser(id bson.ObjectId) (bool, error) {
	var user domain.UserModel
	err := r.collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&user)
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
	var user domain.UserModel
	err := r.collection.FindOne(context.Background(), bson.M{"username": username}).Decode(&user)
	if err != nil {
		return false
	}
	return true
}
