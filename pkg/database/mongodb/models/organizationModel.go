package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Organization_Member struct {
	Name         *string `json:"member_name" validate:"required, min=2, max=50" bson:"member_name"`
	Email        *string `json:"user_email" validate:"required, email" bson:"user_email"`
	Access_Level *string `json:"member_access_level" validate:"required" bson:"member_access_level"`
}

type Organization struct {
	ID                   primitive.ObjectID    `bson:"_id"`
	Organization_ID      *string               `json:"organization_id"`
	Name                 *string               `json:"organization_name" validate:"required, min=2, max=50"`
	Description          *string               `json:"organization_description" validate:"required"`
	Author_Email         *string               `json:"organization_author" validate:"required, email"`
	Organization_Members []Organization_Member `json:"organization_members"`
}
