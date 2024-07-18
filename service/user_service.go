package service

import (
	"encoding/json"
	"net/http"

	"github.com/zeeshanahmad0201/todo-mongo/controller"
	"github.com/zeeshanahmad0201/todo-mongo/model"
)

func Login(w http.ResponseWriter, r *http.Request) {
	var user *model.User
	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if user.Email == nil || user.Password == nil || *user.Email == "" || *user.Password == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	foundUser, err := controller.Login(user)

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(foundUser)
}
