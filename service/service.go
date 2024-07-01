package service

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/zeeshanahmad0201/todo-mongo/controller"
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
