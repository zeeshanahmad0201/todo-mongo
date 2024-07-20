package database

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/zeeshanahmad0201/todo-mongo/common"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	clientInstance      *mongo.Client
	clientInstanceError error
	mongoOnce           sync.Once
)

// InitMongoDB initializes the MongoDB client
func InitMongoDB() *mongo.Client {
	mongoOnce.Do(func() {
		fmt.Println("initMongo Called")
		connectionUrl := os.Getenv("DB_URI")
		if connectionUrl == "" {
			common.HandleError(fmt.Errorf("DB_URI environment variable is not set"), common.ErrorHandlerConfig{
				PrintStackTrace: true,
				Exit:            true,
			})
		}

		clientOptions := options.Client().ApplyURI(connectionUrl)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		clientInstance, clientInstanceError = mongo.Connect(ctx, clientOptions)
		if clientInstanceError != nil {
			common.HandleError(clientInstanceError, common.ErrorHandlerConfig{
				PrintStackTrace: true,
				Exit:            true,
			})
		}

		err := clientInstance.Ping(ctx, nil)
		if err != nil {
			common.HandleError(err, common.ErrorHandlerConfig{
				PrintStackTrace: true,
				Exit:            true,
			})
		}
	})

	return clientInstance
}

func GetDBName() string {
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		common.HandleError(fmt.Errorf("DB_NAME environment variable is not set"), common.ErrorHandlerConfig{
			PrintStackTrace: true,
			Exit:            true,
		})
	}
	return dbName
}

func GetCollection(collectionName string) *mongo.Collection {
	client := InitMongoDB()
	dbName := GetDBName()
	return client.Database(dbName).Collection(collectionName)
}

func CloseMongo() {
	if clientInstance != nil {
		err := clientInstance.Disconnect(context.TODO())
		common.HandleError(err, common.ErrorHandlerConfig{
			Message:         "Failed to disconnect with MongoDB client",
			PrintStackTrace: true,
		})
		fmt.Println("MongoDB client disconnected")
	}
}
