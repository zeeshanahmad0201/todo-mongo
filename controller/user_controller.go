package controller

import (
	"fmt"
	"time"

	"github.com/zeeshanahmad0201/todo-mongo/common"
	"github.com/zeeshanahmad0201/todo-mongo/database"
	"github.com/zeeshanahmad0201/todo-mongo/helpers"
	"github.com/zeeshanahmad0201/todo-mongo/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func SignUp(user *model.User) error {
	ctx, cancel := common.CreateContext(10 * time.Second)
	defer cancel()

	userCollection := database.GetCollection("users")

	// Check if email already exists
	emailFilter := bson.M{"email": user.Email}
	count, err := userCollection.CountDocuments(ctx, emailFilter)

	if err != nil {
		common.HandleError(err)
		return fmt.Errorf("error checking email existence: %w", err)
	}

	if count > 0 {
		err := fmt.Errorf("email already exists")
		common.HandleError(err)
		return err
	}

	// Hash password
	hashedPassword := helpers.HashPassword(*user.Password)
	user.Password = &hashedPassword

	// Set timestamps
	user.AddedOn = *common.GetCurrentTimeStamp()
	user.UpdatedOn = *common.GetCurrentTimeStamp()

	// Generate IDs and tokens
	user.ID = primitive.NewObjectID()
	user.UserID = user.ID.Hex()
	token, refreshToken, err := helpers.GenerateTokens(*user.Name, *user.Email, user.UserID)
	if err != nil {

		return fmt.Errorf("error generating tokens: %w", err)
	}
	user.Token = &token
	user.RefreshToken = &refreshToken

	// Insert user into database
	_, err = userCollection.InsertOne(ctx, user)
	if err != nil {
		return fmt.Errorf("error inserting user: %w", err)
	}

	return nil
}

func Login(user *model.User) (*model.User, error) {
	ctx, cancel := common.CreateContext(10 * time.Second)
	defer cancel()

	userCollection := database.GetCollection("users")

	var foundUser *model.User

	filter := bson.M{"email": user.Email}
	err := userCollection.FindOne(ctx, filter).Decode(&foundUser)

	if err != nil {
		common.HandleError(err)
		return nil, fmt.Errorf("invalid email/password")
	}

	validPass, err := helpers.VerifyPassword(user.Password, foundUser.Password)

	if !validPass {
		common.HandleError(err)
		return nil, fmt.Errorf("invalid email/password")
	}

	token, refreshToken, err := helpers.GenerateTokens(*foundUser.Name, *foundUser.Email, foundUser.UserID)

	if err != nil {
		common.HandleError(err)
		return nil, fmt.Errorf("something went wrong, please try again later")
	}

	helpers.UpdateAllTokens(token, refreshToken, foundUser.UserID)

	return foundUser, nil

}
