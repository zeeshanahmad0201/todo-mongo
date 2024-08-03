package service

import (
	"fmt"
	"time"

	"github.com/zeeshanahmad0201/todo-mongo/common"
	"github.com/zeeshanahmad0201/todo-mongo/database"
	"github.com/zeeshanahmad0201/todo-mongo/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// update existing ToDo in the collection
func UpdateToDo(todo *model.ToDo) (*mongo.UpdateResult, error) {
	ctx, cancel := common.CreateContext(10 * time.Second)
	defer cancel()

	todoCollection := database.GetTodoCollection()

	filter := bson.M{"_id": todo.ID, "userId": todo.UserID}
	todo.UpdatedOn = time.Now()
	update := bson.M{"$set": todo}

	result, err := todoCollection.UpdateOne(ctx, filter, update)

	if err != nil {
		common.HandleError(err, common.ErrorHandlerConfig{
			PrintStackTrace: true,
		})
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("no todo found")
		}
		return nil, err
	}

	return result, nil
}

// delete the todo by task id
func DeleteToDo(id string, userId string) (*mongo.DeleteResult, error) {
	ctx, cancel := common.CreateContext(10 * time.Second)
	defer cancel()

	// convert id string to ObjectID
	obID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		common.HandleError(err, common.ErrorHandlerConfig{
			PrintStackTrace: true,
		})
		return nil, err
	}

	filter := bson.M{"_id": obID, "userId": userId}

	todoCollection := database.GetTodoCollection()

	result, err := todoCollection.DeleteOne(ctx, filter)

	if err != nil {
		common.HandleError(err, common.ErrorHandlerConfig{
			PrintStackTrace: true,
		})
		return nil, err
	}
	return result, nil
}

// get task based on id
func GetTodo(id string, userId string) (*model.ToDo, error) {
	ctx, cancel := common.CreateContext(10 * time.Second)
	defer cancel()

	// Convert the string ID to ObjectID
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		common.HandleError(err, common.ErrorHandlerConfig{
			PrintStackTrace: true,
		})
		return nil, fmt.Errorf("invalid id")
	}

	// Create a filter to search for the document by _id
	filter := bson.M{"_id": objID, "userId": userId}

	// Find the document and decode it into the ToDo struct
	var todo *model.ToDo
	todoCollection := database.GetTodoCollection()
	err = todoCollection.FindOne(ctx, filter).Decode(&todo)
	if err != nil {
		common.HandleError(err)
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("unable to find todo")
		}
		return nil, fmt.Errorf("unable to find todo")
	}

	if todo == nil {
		common.HandleError(fmt.Errorf("no todo found"))
		return nil, fmt.Errorf("no todo found")
	}

	return todo, err
}

func GetAllToDos(userId string) ([]primitive.M, error) {
	ctx, cancel := common.CreateContext(10 * time.Second)
	defer cancel()

	todoCollection := database.GetTodoCollection()
	cursor, err := todoCollection.Find(ctx, bson.D{{}})
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

func CreateOneTodo(todo *model.ToDo) error {
	ctx, cancel := common.CreateContext(10 * time.Second)
	defer cancel()

	todo.ID = primitive.NewObjectID()
	todo.AddedOn = time.Now()

	todoCollection := database.GetTodoCollection()
	_, err := todoCollection.InsertOne(ctx, todo)

	if err != nil {
		common.HandleError(err)
		return err
	}

	return nil
}
