package mongodb

import (
	"context"
	"github.com/LydiaTrack/lydia-track-base/internal/utils"
	"github.com/docker/go-connections/nat"
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
	connected                = false
)

const (
	RemoteConnection    = "REMOTE"
	ContainerConnection = "CONTAINER"
)

// startContainer creates an instance of the mongodb container type for testing
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
	connected = true
	return nil
}

// connectRemoteHost connects to the remote host of the mongodb container
func connectRemoteHost(_ctx context.Context) error {
	// TODO: Remote connection of mongodb container
	utils.LogFatal("Remote connection of mongodb container is not implemented yet!")
	panic("Remote connection of mongodb container is not implemented yet!")
	return nil
}

// InitializeMongoDBConnection initializes the mongodb container
func InitializeMongoDBConnection() error {
	ctx := context.Background()

	// Check for mongoDB connection type, if it is remote, connect to the remote host
	// else, start the container on local machine and connect to it
	connectionType := os.Getenv("LYDIA_DB_CONNECTION_TYPE")
	if connectionType == RemoteConnection {
		err := connectRemoteHost(ctx)
		if err != nil {
			return err
		}
	} else if connectionType == ContainerConnection {
		err := startContainer(ctx)
		if err != nil {
			return err
		}
	} else {
		utils.LogFatal("Cannot detect connection type for mongodb container!")
		panic("Invalid connection type for mongodb container!")
	}

	return nil
}

// getContainer returns the mongodb container instance, if instance is not initialized, it panics
func getContainer() *mongodbContainer {
	if !connected {
		panic("Mongodb container not initialized!")
	}
	return mongodbContainerInstance
}

// GetCollection returns the mongodb collection that is connected with a mongoDB container
func GetCollection(collectionName string, ctx context.Context) *mongo.Collection {
	container := getContainer()
	host, err := container.Host(ctx)
	if err != nil {
		return nil
	}

	portNumber := os.Getenv("LYDIA_DB_PORT")
	port, err := container.MappedPort(ctx, nat.Port(portNumber))
	if err != nil {
		return nil
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://"+host+":"+port.Port()))
	if err != nil {
		return nil
	}

	return client.Database(os.Getenv("LYDIA_DB_NAME")).Collection(collectionName)
}
