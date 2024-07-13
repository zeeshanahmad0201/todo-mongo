package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/zeeshanahmad0201/todo-mongo/common"
	"github.com/zeeshanahmad0201/todo-mongo/database"
	"github.com/zeeshanahmad0201/todo-mongo/helpers"
	"github.com/zeeshanahmad0201/todo-mongo/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection = database.UserCollection

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

func SignUp(w http.ResponseWriter, r *http.Request) error {
	ctx, cancel := common.CreateContext(10 * time.Second)
	defer cancel()

	var user *model.User

	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		common.HandleError(err)
		return err
	}

	validate := validator.New()

	err = validate.Struct(user)

	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			common.HandleError(err)
			return err
		}
	}

	emailFilter := bson.M{"email": user.Email}
	count, err := todoCollection.CountDocuments(ctx, emailFilter)

	if err != nil {
		common.HandleError(err)
		return err
	}

	if count > 0 {
		return fmt.Errorf("Email already exists!")
	}

	password := HashPassword(*user.Password)
	user.Password = &password

	user.AddedOn, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.UpdatedOn, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.ID = primitive.NewObjectID()
	user.UserID = user.ID.Hex()
	token, refreshToken, err := helpers.GenerateTokens(*user.Name, *user.Email, user.UserID)
	user.Token = &token
	user.RefreshToken = &refreshToken

}
