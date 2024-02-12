package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Organization_Member struct {
	Name  *string `json:"member_name" validate:"required" bson:"member_name"`
	Email *string `json:"user_email" validate:"required" bson:"user_email"`
}

type Organization struct {
	ID                   primitive.ObjectID    `json:"_id" bson:"_id"`
	Organization_ID      string                `json:"organization_id"`
	Name                 *string               `json:"name" validate:"required,min=2,max=50"`
	Description          *string               `json:"description" validate:"required"`
	Author_Email         string                `json:"organization_author" validate:"required,email"`
	Organization_Members []Organization_Member `json:"organization_members"`
}
