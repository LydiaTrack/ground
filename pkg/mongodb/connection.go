package mongodb

import (
	"context"
	"github.com/LydiaTrack/ground/pkg/log"
	"os"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// mongodbContainer represents the mongodb container type used in the module
type mongodbContainer struct {
	testcontainers.Container
}

var (
	// mongodbContainerInstance is the instance of the mongodb container
	mongodbContainerInstance *mongodbContainer
	connected                = false
	connectionType           string
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

// InitializeMongoDBConnection initializes the mongodb container
func InitializeMongoDBConnection() error {
	ctx := context.Background()

	// Check for mongoDB connection type, if it is remote, connect to the remote host
	// else, start the container on local machine and connect to it
	connectionTypeEnv := os.Getenv("DB_CONNECTION_TYPE")
	if connectionTypeEnv == ContainerConnection {
		err := startContainer(ctx)
		if err != nil {
			return err
		}
		log.Log("Connecting to local mongodb container...")
		connectionType = ContainerConnection
	} else if connectionTypeEnv == RemoteConnection {
		log.Log("Connecting to remote host of mongodb container...")
		connectionType = RemoteConnection
	} else {
		log.LogFatal("Invalid connection type for mongodb container: " + connectionTypeEnv)
		return nil
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
func GetCollection(collectionName string) (*mongo.Collection, error) {
	ctx := context.Background()
	if connectionType == ContainerConnection {
		container := getContainer()
		host, err := container.Host(ctx)
		if err != nil {
			return nil, err

		}
		portNumber := os.Getenv("DB_PORT")
		port, err := container.MappedPort(ctx, nat.Port(portNumber))
		if err != nil {
			return nil, err
		}

		client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://"+host+":"+port.Port()))
		if err != nil {
			return nil, err
		}

		return client.Database(os.Getenv("DB_NAME")).Collection(collectionName), nil
	} else if connectionType == RemoteConnection {
		serverAPI := options.ServerAPI(options.ServerAPIVersion1)
		opts := options.Client().ApplyURI(os.Getenv("DB_URI")).SetServerAPIOptions(serverAPI)
		// Create a new client and connect to the server
		client, err := mongo.Connect(context.TODO(), opts)
		if err != nil {
			return nil, err
		}

		return client.Database(os.Getenv("DB_NAME")).Collection(collectionName), nil
	}
	log.LogFatal("Invalid connection type for mongodb container: " + connectionType)
	return nil, nil
}
