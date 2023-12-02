package mongodb

import (
	"context"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
)

// mongodbContainer represents the mongodb container type used in the module
type mongodbContainer struct {
	testcontainers.Container
}

var (
	// mongodbContainerInstance is the instance of the mongodb container
	mongodbContainerInstance *mongodbContainer
	initilized               = false
)

// StartContainer creates an instance of the mongodb container type for testing
func startContainer(ctx context.Context) error {
	req := testcontainers.ContainerRequest{
		Image:        "mongo:latest",
		ExposedPorts: []string{"27017/tcp"},
		WaitingFor: wait.ForAll(
			wait.ForLog("Waiting for connections"),
			wait.ForListeningPort("27017/tcp"),
		),
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return err
	}
	mongodbContainerInstance = &mongodbContainer{container}
	initilized = true
	return nil
}

// InitializeContainer initializes the mongodb container
func InitializeContainer() error {
	ctx := context.Background()
	err := startContainer(ctx)
	if err != nil {
		return err
	}
	return nil
}

// GetContainer returns the mongodb container instance
func getContainer() *mongodbContainer {
	if !initilized {
		panic("Mongodb container not initialized!")
	}
	return mongodbContainerInstance
}

// GetCollection returns the mongodb collection
func GetCollection(collectionName string, ctx context.Context) *mongo.Collection {
	container := getContainer()
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

	return client.Database(os.Getenv("LYDIA_DB_NAME")).Collection(collectionName)
}
