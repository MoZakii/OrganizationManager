package utils

import (
	"MoZaki-Organization-Manager/pkg/database/mongodb/repository"
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = repository.OpenCollection(repository.Client, "user")

func GenerateAllTokens(email string, name string, uid string) (signedToken string, signedRefreshToken string, err error) {

	err = godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	var SECRET_KEY string = os.Getenv("SECRET_KEY")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"Email": email,
		"Name":  name,
		"Uid":   uid,
		"exp":   time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
	})

	refresh := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"Email": email,
		"Name":  name,
		"Uid":   uid,
		"exp":   time.Now().Local().Add(time.Hour * time.Duration(169)).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(SECRET_KEY))

	if err != nil {
		log.Panic(err)
		return
	}
	refreshToken, err := refresh.SignedString([]byte(SECRET_KEY))

	if err != nil {
		log.Panic(err)
		return
	}

	return tokenString, refreshToken, err
}

func ValidateToken(tokenString string) (jwt.MapClaims, error) {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	var SECRET_KEY string = os.Getenv("SECRET_KEY")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(SECRET_KEY), nil
	})
	if err != nil {
		log.Fatal(err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, err
	}

	tokenTime := claims["exp"]

	if int64(tokenTime.(float64)) < time.Now().Local().Unix() {
		err = errors.New("token is expired")
		return nil, err
	}

	return claims, nil
}

func UpdateAllTokens(signedToken string, signedRefreshToken string, userId string) (string, string) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var err error

	_, err = ValidateToken(signedRefreshToken)
	if err != nil {
		log.Panic(err)
		return "", ""
	}

	var updateObj primitive.D

	updateObj = append(updateObj, bson.E{Key: "token", Value: signedToken})
	updateObj = append(updateObj, bson.E{Key: "refresh_token", Value: signedRefreshToken})

	filter := bson.M{"user_id": userId}

	_, err = userCollection.UpdateOne(
		ctx,
		filter,
		bson.D{
			{Key: "$set", Value: updateObj},
		},
	)

	if err != nil {
		log.Panic(err)
		return "", ""
	}

	return signedToken, signedRefreshToken

}
