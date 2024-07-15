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

func SignUp(w http.ResponseWriter, r *http.Request) (*mongo.InsertOneResult, error) {
	ctx, cancel := common.CreateContext(10 * time.Second)
	defer cancel()

	var user *model.User

	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		common.HandleError(err)
		return nil, err
	}

	validate := validator.New()

	err = validate.Struct(user)

	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			common.HandleError(err)
			return nil, err
		}
	}

	emailFilter := bson.M{"email": user.Email}
	count, err := todoCollection.CountDocuments(ctx, emailFilter)

	if err != nil {
		common.HandleError(err)
		return nil, err
	}

	if count > 0 {
		err := fmt.Errorf("Email already exists!")
		common.HandleError(err)
		return nil, err
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

	result, err := userCollection.InsertOne(ctx, user)
	if err != nil {
		common.HandleError(err)
		return nil, err
	}

	return result, nil
}

func Login(w http.ResponseWriter, r *http.Request) error {
	ctx, cancel := common.CreateContext(10 * time.Second)
	defer cancel()

	var user *model.User
	var foundUser *model.User

	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		common.HandleError(err)
		return err
	}

	filter := bson.M{"email": user.Email}
	err = userCollection.FindOne(ctx, filter).Decode(&foundUser)

	if err != nil {
		common.HandleError(err)
		return fmt.Errorf("Invalid email/password")
	}

	validPass, err := helpers.VerifyPassword(*user.Password, *foundUser.Password)

	if !validPass {
		common.HandleError(err)
		return fmt.Errorf("Invalid email/password")
	}

	token, refreshToken, err := helpers.GenerateTokens(*foundUser.Name, *foundUser.Email, *&foundUser.UserID)

	if err != nil {
		common.HandleError(err)
		return fmt.Errorf("Something went wrong. Please try again later!")
	}

	helpers.U

	return nil

}
