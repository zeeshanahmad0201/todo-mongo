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
