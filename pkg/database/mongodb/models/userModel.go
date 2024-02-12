package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID            primitive.ObjectID `bson:"_id"`
	Name          *string            `json:"name" validate:"required,min=2,max=50"`
	Email         *string            `json:"email" validate:"required,email"`
	Password      *string            `json:"password" validate:"required,min=6"`
	Token         *string            `json:"token"`
	Refresh_Token *string            `json:"refresh_token"`
	User_id       string             `json:"user_id"`
}
