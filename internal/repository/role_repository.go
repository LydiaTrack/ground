package repository

import (
	"context"
	"log"
	"lydia-track-base/internal/domain"
	"lydia-track-base/internal/mongodb"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

// A RoleMongoRepository that implements RoleRepository
type RoleMongoRepository struct {
	client     *mongo.Client
	collection *mongo.Collection
}

// NewRoleMongoRepository creates a new RoleMongoRepository instance
// which implements RoleRepository
func NewRoleMongoRepository() *RoleMongoRepository {
	ctx := context.Background()
	// FIXME: Burada ileride uzaktaki bir mongodb instance'ına bağlanmak gerekecek
	// FIXME: Ortaklaştırılacak
	container, err := mongodb.StartContainer(ctx)
	if err != nil {
		return nil
	}

	host, err := container.Host(ctx)
	if err != nil {
		return nil
	}

	port, err := container.MappedPort(ctx, "27017")
	if err != nil {
		return nil
	}

	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://"+host+":"+port.Port()))
	if err != nil {
		return nil
	}

	collection := client.Database(os.Getenv("LYDIA_DB_NAME")).Collection("roles")

	return &RoleMongoRepository{
		client:     client,
		collection: collection,
	}
}

// Saverole saves a role
func (r *RoleMongoRepository) SaveRole(role domain.RoleModel) (domain.RoleModel, error) {
	_, err := r.collection.InsertOne(context.Background(), role)
	if err != nil {
		return domain.RoleModel{}, err
	}
	return role, nil
}

// GetRole gets a role by id
func (r *RoleMongoRepository) GetRole(id bson.ObjectId) (domain.RoleModel, error) {
	var role domain.RoleModel
	err := r.collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&role)
	if err != nil {
		return domain.RoleModel{}, err
	}
	return role, nil
}

// ExistsRole checks if a role exists by role id
func (r *RoleMongoRepository) ExistsRole(id bson.ObjectId) (bool, error) {
	var role domain.RoleModel
	err := r.collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&role)
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
