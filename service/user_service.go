package service

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
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

	validate := validator.New()
	if err := validate.Struct(user); err != nil {
		var validationErrors []string
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, err.Error())
		}
		http.Error(w, "Validation failed: "+strings.Join(validationErrors, ", "), http.StatusBadRequest)
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

func SignUp(w http.ResponseWriter, r *http.Request) {
	var user model.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	validate := validator.New()

	if err := validate.Struct(user); err != nil {
		var validationErrors []string
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, err.Error())
		}
		http.Error(w, "Validation failed: "+strings.Join(validationErrors, ", "), http.StatusBadRequest)
		return
	}

	err = controller.SignUp(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User created successfully"})
}
