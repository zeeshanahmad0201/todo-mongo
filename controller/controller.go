package controller

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/zeeshanahmad0201/todo-mongo/common"
	"github.com/zeeshanahmad0201/todo-mongo/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collection *mongo.Collection
var client *mongo.Client

func initMongo() {
	fmt.Println("initMongo Called")
	connectionUrl := os.Getenv("DB_URI")
	dbName := os.Getenv("DB_NAME")
	collectionName := os.Getenv("DB_COLLECTION")

	clientOptions := options.Client().ApplyURI(connectionUrl)

	// Create a context with a timeout
	ctx, cancel := common.CreateContext(10 * time.Second)
	defer cancel()

	// Connect to mongo
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		common.HandleError(err, common.ErrorHandlerConfig{
			Message:         "Failed to connect to MongoDB",
			PrintStackTrace: true,
			Exit:            true,
		})
	}

	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		common.HandleError(err, common.ErrorHandlerConfig{
			Message:         "Failed to ping MongoDB",
			PrintStackTrace: true,
			Exit:            true,
		})
	}

	collection = client.Database(dbName).Collection(collectionName)
	fmt.Printf("Connection established with %s on %s db\n", collectionName, dbName)
}

// add new ToDo in the collection
func addTodo(todo *model.ToDo) string {
	ctx, cancel := common.CreateContext(10 * time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, todo)

	if err != nil {
		common.HandleError(err, common.ErrorHandlerConfig{
			Message:         "Unable to add ToDo",
			PrintStackTrace: true,
		})
		return "Unable to add Task"
	}
	return "Task added successfully"
}

// update existing ToDo in the collection
func updateToDo(todo *model.ToDo) string {
	ctx, cancel := common.CreateContext(10 * time.Second)
	defer cancel()

	filter := bson.M{"_id": todo.ID}
	update := bson.M{"$set": todo}

	_, err := collection.UpdateOne(ctx, filter, update)

	if err != nil {
		common.HandleError(err, common.ErrorHandlerConfig{
			Message:         "Unable to update ToDo",
			PrintStackTrace: true,
		})
		return "Unable to update Task"
	}
	return "Task updated successfully"
}

func deleteToDo(id string) string {
	ctx, cancel := common.CreateContext(10 * time.Second)
	defer cancel()

	// convert id string to ObjectID
	obID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		common.HandleError(err, common.ErrorHandlerConfig{
			Message:         "Invalid ID",
			PrintStackTrace: true,
		})
		return "Invalid task ID"
	}

	filter := bson.M{"_id": obID}

	_, err = collection.DeleteOne(ctx, filter)

	if err != nil {
		common.HandleError(err, common.ErrorHandlerConfig{
			Message:         "Unable to delete Todo",
			PrintStackTrace: true,
		})
		return "Unable to delete Task"
	}
	return "Task deleted successfully"
}

func closeMongo() {
	if client != nil {
		err := client.Disconnect(context.TODO())
		common.HandleError(err, common.ErrorHandlerConfig{
			Message:         "Failed to disconnect with MongoDB client",
			PrintStackTrace: true,
		})
		fmt.Println("MongoDB client disconnected")
	}
}
