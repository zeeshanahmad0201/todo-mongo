package service

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/zeeshanahmad0201/todo-mongo/controller"
	"github.com/zeeshanahmad0201/todo-mongo/helpers"
	"github.com/zeeshanahmad0201/todo-mongo/model"
)

func GetOneTodo(w http.ResponseWriter, r *http.Request) {

	claims, err := helpers.ExtractAndValidateToken(r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)

	id := vars["id"]

	todo, err := controller.GetTodo(id, claims.UserId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if todo == nil {
		http.Error(w, "no todo found!", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todo)
}

func GetAllToDos(w http.ResponseWriter, r *http.Request) {

	claims, err := helpers.ExtractAndValidateToken(r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	todos, err := controller.GetAllToDos(claims.UserId)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			http.Error(w, "No todos added yet!", http.StatusNotFound)
			return
		}
		http.Error(w, "Unable to fetch data at this time", http.StatusInternalServerError)
		return
	}

	if len(todos) == 0 {
		http.Error(w, "No todos added yet!", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(todos); err != nil {
		http.Error(w, "Error parsing data", http.StatusInternalServerError)
	}
}

func UpdateOneToDo(w http.ResponseWriter, r *http.Request) {
	// extract id from url
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		http.Error(w, "id parameter is missing!", http.StatusBadRequest)
		return
	}

	// extract data from body
	var todo *model.ToDo
	err := json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// update todo
	result, err := controller.UpdateToDo(todo)

	if err != nil {
		http.Error(w, "Unable to update ToDo", http.StatusInternalServerError)
		return
	}

	if result.ModifiedCount == 0 {
		http.Error(w, "No matching todos found!", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func DeleteOneTodo(w http.ResponseWriter, r *http.Request) {

	claims, err := helpers.ExtractAndValidateToken(r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// extract the id from url
	vars := mux.Vars(r)

	id := vars["id"]

	if id == "" {
		http.Error(w, "id parameter is missing!", http.StatusBadRequest)
		return
	}

	// delete todo
	result, err := controller.DeleteToDo(id, claims.UserId)

	if err != nil {
		http.Error(w, "Unable to delete todo!", http.StatusInternalServerError)
		return
	}

	if result.DeletedCount == 0 {
		http.Error(w, "No matching todo found!", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(result)
}

func CreateOneTodo(w http.ResponseWriter, r *http.Request) {

	claims, err := helpers.ExtractAndValidateToken(r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	var todo *model.ToDo

	err = json.NewDecoder(r.Body).Decode(&todo)

	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if todo.Title == "" {
		http.Error(w, "Title can't be empty", http.StatusBadRequest)
		return
	}

	todo.UserID = claims.UserId

	err = controller.CreateOneTodo(todo)

	if err != nil {
		http.Error(w, "Unable to add todo", http.StatusInternalServerError)
		return
	}
}
