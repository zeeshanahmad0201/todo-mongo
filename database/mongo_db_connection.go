package database

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/zeeshanahmad0201/todo-mongo/common"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client
var ToDoCollection *mongo.Collection
var UserCollection *mongo.Collection

func InitMongoDB() {
	fmt.Println("initMongo Called")
	connectionUrl := os.Getenv("DB_URI")
	if connectionUrl == "" {
		common.HandleError(fmt.Errorf("DB_URI environment variable is not set"), common.ErrorHandlerConfig{
			PrintStackTrace: true,
			Exit:            true,
		})
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		common.HandleError(fmt.Errorf("DB_NAME environment variable is not set"), common.ErrorHandlerConfig{
			PrintStackTrace: true,
			Exit:            true,
		})
	}

	todoCollectionName := os.Getenv("DB_TODO_COLLECTION")
	if todoCollectionName == "" {
		common.HandleError(fmt.Errorf("DB_TODO_COLLECTION environment variable is not set"), common.ErrorHandlerConfig{
			PrintStackTrace: true,
			Exit:            true,
		})
	}

	userCollectionName := os.Getenv("DB_USER_COLLECTION")
	if userCollectionName == "" {
		common.HandleError(fmt.Errorf("DB_USER_COLLECTION environment variable is not set"), common.ErrorHandlerConfig{
			PrintStackTrace: true,
			Exit:            true,
		})
	}

	clientOptions := options.Client().ApplyURI(connectionUrl)

	// Create a context with a timeout
	ctx, cancel := common.CreateContext(10 * time.Second)
	defer cancel()

	// Connect to mongo
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		common.HandleError(err, common.ErrorHandlerConfig{
			PrintStackTrace: true,
			Exit:            true,
		})
	}

	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		common.HandleError(err, common.ErrorHandlerConfig{
			PrintStackTrace: true,
			Exit:            true,
		})
	}

	Client = client
	db := Client.Database(dbName)
	ToDoCollection = db.Collection(todoCollectionName)
	UserCollection = db.Collection(userCollectionName)
}

func CloseMongo() {
	if Client != nil {
		err := Client.Disconnect(context.TODO())
		common.HandleError(err, common.ErrorHandlerConfig{
			Message:         "Failed to disconnect with MongoDB client",
			PrintStackTrace: true,
		})
		fmt.Println("MongoDB client disconnected")
	}
}
