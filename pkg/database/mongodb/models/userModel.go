package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID            primitive.ObjectID `bson:"_id"`
	Name          *string            `json:"user_name" validate:"required"`
	Email         *string            `json:"user_email" validate:"required"`
	Password      *string            `json:"user_password" validate:"required"`
	Token         *string            `json:"token"`
	Refresh_Token *string            `json:"refresh_token"`
	User_id       string             `json:"user_id"`
}
