package router

import (
	"github.com/gorilla/mux"
	"github.com/zeeshanahmad0201/todo-mongo/controller"
)

func InitRouter() *mux.Router {

	r := mux.NewRouter()

	// Auth routes
	r.HandleFunc("/login", controller.Login).Methods("POST")
	r.HandleFunc("/signup", controller.SignUp).Methods("POST")

	r.HandleFunc("/todos", controller.GetAllToDos).Methods("GET")
	r.HandleFunc("/todos/{id}", controller.GetOneTodo).Methods("GET")
	r.HandleFunc("/todos/{id}", controller.DeleteOneTodo).Methods("DELETE")
	r.HandleFunc("/todos/{id}", controller.UpdateOneToDo).Methods("PUT")
	r.HandleFunc("/todos", controller.CreateOneTodo).Methods("POST")

	return r
}
