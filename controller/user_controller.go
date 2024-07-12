package controller

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/zeeshanahmad0201/todo-mongo/common"
	"github.com/zeeshanahmad0201/todo-mongo/database"
	"github.com/zeeshanahmad0201/todo-mongo/model"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var collection *mongo.Collection = database.UserCollection

// Encrypt the password before it is stored in the DB
func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)

	if err != nil {
		common.HandleError(err, common.ErrorHandlerConfig{
			Exit: true,
		})
	}

	return string(bytes)
}

// Validates and verifies the password in the DB
func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	isCorrect := true
	msg := ""

	if err != nil {
		common.HandleError(err)
		isCorrect = false
		msg = "email or password is incorrect"
	}

	return isCorrect, msg
}

func SignUp(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := common.CreateContext(10 * time.Second)
	defer cancel()

	var user *model.User

	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		common.HandleError(err)
		return
	}
}
