package repository

import (
	"MoZaki-Organization-Manager/pkg/database/mongodb/models"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = OpenCollection(Client, "user")

// Function that updates user's tokens in the database
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
		return err
	}

	return err
}

// Function that creates a user in the database
func CreateUser(user models.User) (err error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	_, err = userCollection.InsertOne(ctx, user)
	return
}

// Function that returns the count of users with a given email from the database
func CountUsersByEmail(email *string) (count int64, err error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	count, err = userCollection.CountDocuments(ctx, bson.M{"email": email})
	return
}

// Function that returns a user object from the database using their email
func GetUserByEmail(email *string) (foundUser models.User, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	err = userCollection.FindOne(ctx, bson.M{"email": *email}).Decode(&foundUser)

	return
}
