package repository

import (
	"context"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
	"log"
	"lydia-track-base/internal/domain"
	"os"
)

// A UserMongoRepository that implements UserRepository
type UserMongoRepository struct {
	client     *mongo.Client
	collection *mongo.Collection
}

// NewUserMongoRepository creates a new UserMongoRepository
func NewUserMongoRepository() *UserMongoRepository {
	// Get connection string from env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	uri := os.Getenv("MONGODB_URI")

	// Create client
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	// Get collection
	collection := client.Database("lydia-track-base").Collection("users")

	return &UserMongoRepository{
		client:     client,
		collection: collection,
	}
}

// SaveUser saves a user
func (r *UserMongoRepository) SaveUser(user domain.User) (domain.User, error) {
	_, err := r.collection.InsertOne(context.Background(), user)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

// GetUser gets a user by id
func (r *UserMongoRepository) GetUser(id bson.ObjectId) (domain.User, error) {
	var user domain.User
	err := r.collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

// ExistsUser checks if a user exists
func (r *UserMongoRepository) ExistsUser(id bson.ObjectId) (bool, error) {
	var user domain.User
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
