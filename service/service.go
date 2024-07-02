package service

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/zeeshanahmad0201/todo-mongo/controller"
	"github.com/zeeshanahmad0201/todo-mongo/model"
)

func GetOneTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id := vars["id"]

	todo, err := controller.GetTodo(id)
	if err != nil {
		http.Error(w, "Unable to find todo", http.StatusInternalServerError)
		return
	}

	if todo == nil {
		http.Error(w, "No Todo Found!", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todo)
}

func GetAllToDos(w http.ResponseWriter, r *http.Request) {
	todos, err := controller.GetAllToDos()
	if err != nil {
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
	// extract the id from url
	vars := mux.Vars(r)

	id := vars["id"]

	if id == "" {
		http.Error(w, "id parameter is missing!", http.StatusBadRequest)
		return
	}

	// delete todo
	result, err := controller.DeleteToDo(id)

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
