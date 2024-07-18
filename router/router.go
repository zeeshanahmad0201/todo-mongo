package router

import (
	"github.com/gorilla/mux"
	"github.com/zeeshanahmad0201/todo-mongo/service"
)

func InitRouter() *mux.Router {
	r := mux.NewRouter()

	// Auth routes
	r.HandleFunc("/login", service.Login).Methods("GET")

	r.HandleFunc("/todos", service.GetAllToDos).Methods("GET")
	r.HandleFunc("/todos/{id}", service.GetOneTodo).Methods("GET")
	r.HandleFunc("/todos/{id}", service.DeleteOneTodo).Methods("DELETE")
	r.HandleFunc("/todos/{id}", service.UpdateOneToDo).Methods("PUT")
	r.HandleFunc("/todos", service.CreateOneTodo).Methods("POST")

	return r
}
