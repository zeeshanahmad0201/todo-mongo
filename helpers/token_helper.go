package helpers

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/zeeshanahmad0201/todo-mongo/common"
	"github.com/zeeshanahmad0201/todo-mongo/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

type SignedDetails struct {
	Name   string
	Email  string
	userId string
	jwt.StandardClaims
}

// generates the detail and refesh token
func GenerateTokens(name string, email string, userId string) (token string, refreshToken string, err error) {

	SECRET_KEY := os.Getenv("SECRET_KEY")

	claims := &SignedDetails{
		Email:  email,
		Name:   name,
		userId: userId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	refreshClaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
		},
	}

	token, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))

	if err != nil {
		common.HandleError(err)
		return "", "", err
	}

	refreshToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SECRET_KEY))

	if err != nil {
		common.HandleError(err)
		return "", "", err
	}

	return token, refreshToken, nil
}

func VerifyPassword(userPassword string, providedPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))

	if err != nil {
		return false, err
	}

	return true, nil
}

// Renews the user tokens when they login
func UpdateAllTokens(signedToken string, signedRefreshToken string, userId string) error {
	var userCollection *mongo.Collection = database.UserCollection

	ctx, cancel := common.CreateContext(10 * time.Second)
	defer cancel()

	var updateObj primitive.D

	updateObj = append(updateObj, bson.E{Key: "token", Value: signedToken})
	updateObj = append(updateObj, bson.E{Key: "refresh_token", Value: signedRefreshToken})
	updateObj = append(updateObj, bson.E{Key: "updated_on", Value: common.GetCurrentTimeStamp()})

	filter := bson.M{"user_id": userId}
	upsert := true
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	_, err := userCollection.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: updateObj}}, &opt)

	if err != nil {
		common.HandleError(err)
		return err
	}

	return nil

}
