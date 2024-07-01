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

func InitMongo() {
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
			PrintStackTrace: true,
			Exit:            true,
		})
		return
	}

	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		common.HandleError(err, common.ErrorHandlerConfig{
			PrintStackTrace: true,
			Exit:            true,
		})
		return
	}

	collection = client.Database(dbName).Collection(collectionName)
	fmt.Printf("Connection established with %s on %s db\n", collectionName, dbName)
}

// add new ToDo in the collection
func AddTodo(todo *model.ToDo) string {
	ctx, cancel := common.CreateContext(10 * time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, todo)

	if err != nil {
		common.HandleError(err, common.ErrorHandlerConfig{
			PrintStackTrace: true,
		})
		return "Unable to add Task"
	}
	return "Task added successfully"
}

// update existing ToDo in the collection
func UpdateToDo(todo *model.ToDo) string {
	ctx, cancel := common.CreateContext(10 * time.Second)
	defer cancel()

	filter := bson.M{"_id": todo.ID}
	update := bson.M{"$set": todo}

	_, err := collection.UpdateOne(ctx, filter, update)

	if err != nil {
		common.HandleError(err, common.ErrorHandlerConfig{
			PrintStackTrace: true,
		})
		return "Unable to update Task"
	}
	return "Task updated successfully"
}

// delete the todo by task id
func DeleteToDo(id string) string {
	ctx, cancel := common.CreateContext(10 * time.Second)
	defer cancel()

	// convert id string to ObjectID
	obID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		common.HandleError(err, common.ErrorHandlerConfig{
			PrintStackTrace: true,
		})
		return "Invalid task ID"
	}

	filter := bson.M{"_id": obID}

	_, err = collection.DeleteOne(ctx, filter)

	if err != nil {
		common.HandleError(err, common.ErrorHandlerConfig{
			PrintStackTrace: true,
		})
		return "Unable to delete Task"
	}
	return "Task deleted successfully"
}

// get task based on id
func GetTodo(id string) (*model.ToDo, error) {
	ctx, cancel := common.CreateContext(10 * time.Second)
	defer cancel()

	// Convert the string ID to ObjectID
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		common.HandleError(err, common.ErrorHandlerConfig{
			PrintStackTrace: true,
		})
		return nil, err
	}

	// Create a filter to search for the document by _id
	filter := bson.M{"_id": objID}

	// Find the document and decode it into the ToDo struct
	var todo *model.ToDo
	err = collection.FindOne(ctx, filter).Decode(&todo)
	if err != nil {
		common.HandleError(err)
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return todo, err
}

func GetAllToDos() ([]primitive.M, error) {
	ctx, cancel := common.CreateContext(10 * time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.D{{}})
	if err != nil {
		common.HandleError(err)
		return nil, err
	}

	defer cursor.Close(ctx)

	var todos []primitive.M

	for cursor.Next(ctx) {
		var todo bson.M
		err := cursor.Decode(&todo)
		if err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}

	// check for error that may have occured during iteration
	if err := cursor.Err(); err != nil {
		common.HandleError(err)
		return nil, err
	}

	return todos, nil
}

func CloseMongo() {
	if client != nil {
		err := client.Disconnect(context.TODO())
		common.HandleError(err, common.ErrorHandlerConfig{
			Message:         "Failed to disconnect with MongoDB client",
			PrintStackTrace: true,
		})
		fmt.Println("MongoDB client disconnected")
	}
}
