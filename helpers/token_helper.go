package helpers

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/zeeshanahmad0201/todo-mongo/common"
	"github.com/zeeshanahmad0201/todo-mongo/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SignedDetails struct {
	Name   string
	Email  string
	UserId string
	jwt.StandardClaims
}

// generates the detail and refesh token
func GenerateTokens(name string, email string, userId string) (token string, refreshToken string, err error) {

	SECRET_KEY := os.Getenv("SECRET_KEY")

	claims := &SignedDetails{
		Email:  email,
		Name:   name,
		UserId: userId,
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

// Renews the user tokens when they login
func UpdateAllTokens(signedToken string, signedRefreshToken string, userId string) error {
	userCollection := database.GetCollection("users")

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

func ValidateToken(clientToken string) (claims *SignedDetails, err error) {
	SECRET_KEY := os.Getenv("SECRET_KEY")
	token, err := jwt.ParseWithClaims(clientToken, &SignedDetails{}, func(t *jwt.Token) (interface{}, error) { return []byte(SECRET_KEY), nil })

	if err != nil {
		common.HandleError(err)
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(*SignedDetails)

	if !ok {
		return nil, fmt.Errorf("invalid token")
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		return nil, fmt.Errorf("token is expired")
	}

	return claims, nil
}

// Extracts the token from the request header
func ExtractToken(r *http.Request) (string, error) {
	// validate header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("auth header is missing")
	}
	fmt.Println(authHeader)

	// extract token
	parts := strings.Split(authHeader, " ")
	fmt.Println(parts)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", fmt.Errorf("auth header must be in format `Bearer {token}`")
	}

	return parts[1], nil
}

// Extract and Validate the token
func ExtractAndValidateToken(r *http.Request) (*SignedDetails, error) {
	token, err := ExtractToken(r)
	if err != nil {
		return nil, err
	}

	claims, err := ValidateToken(token)
	if err != nil {
		return nil, err
	}

	return claims, nil
}
