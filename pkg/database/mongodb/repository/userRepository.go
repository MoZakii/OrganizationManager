package repository

import (
	"MoZaki-Organization-Manager/pkg/database/mongodb/models"
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = OpenCollection(Client, "user")

func UpdateTokens(signedToken string, signedRefreshToken string, userId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"token":         signedToken,
			"refresh_token": signedRefreshToken,
		},
	}

	filter := bson.M{"user_id": userId}

	_, err := userCollection.UpdateOne(ctx, filter, update)

	if err != nil {
		log.Panic(err)
		return err
	}

	return err
}

func CreateUser(user models.User) (err error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	_, err = userCollection.InsertOne(ctx, user)
	return
}

func CountUsersByEmail(email *string) (count int64, err error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	count, err = userCollection.CountDocuments(ctx, bson.M{"email": email})
	return
}

func GetUserByEmail(email *string) (foundUser models.User, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	err = userCollection.FindOne(ctx, bson.M{"email": *email}).Decode(&foundUser)

	return
}
