package repository

import (
	"MoZaki-Organization-Manager/pkg/database/mongodb/models"
	"context"
	"errors"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var organizationCollection *mongo.Collection = OpenCollection(Client, "organization")

// Function that creates organization in the database
func CreateOrganization(organization models.Organization) (err error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	_, err = organizationCollection.InsertOne(ctx, organization)
	return
}

// Function that returns an organization from the database

func GetOrganization(organizationID string) (organization *models.Organization, err error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	err = organizationCollection.FindOne(ctx, bson.M{"organization_id": organizationID}).Decode(&organization)
	if err != nil {
		return
	}
	return
}

type NeededInfo struct {
	Organization_ID     string
	Name                string
	Description         string
	OrganizationMembers []models.Organization_Member
}

// Function that returns all organizations from the database

func GetAllOrganizations() (organizations []NeededInfo, err error) {

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	cursor, err := organizationCollection.Find(ctx, bson.M{})
	if err != nil {
		return
	}

	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var organization models.Organization
		if err := cursor.Decode(&organization); err != nil {
			return nil, err
		}
		neededInfo := NeededInfo{
			Name:                *organization.Name,
			Description:         *organization.Description,
			Organization_ID:     organization.Organization_ID,
			OrganizationMembers: organization.Organization_Members,
		}
		organizations = append(organizations, neededInfo)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return
}

// Function that updates an organization's data in the database

func UpdateOrganization(organization models.Organization) (err error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	filter := bson.M{"organization_id": organization.Organization_ID}

	update := bson.M{
		"$set": bson.M{
			"name":        organization.Name,
			"description": organization.Description,
		},
	}

	_, err = organizationCollection.UpdateOne(ctx, filter, update)

	return
}

// Function that deletes an organization in the database

func DeleteOrganization(organizationID string) (err error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	filter := bson.M{"organization_id": organizationID}
	result, err := organizationCollection.DeleteOne(ctx, filter)
	if err != nil || result.DeletedCount == 0 {
		return errors.New("organization id doesn't match")
	}
	return
}

// Function that adds a member to an organization data in the database

func AddToOrganization(organizationID string, member models.Organization_Member) error {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	filter := bson.M{"organization_id": organizationID}
	update := bson.M{"$addToSet": bson.M{"organization_members": member}}

	result, err := organizationCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Println("Error adding member to organization:", err)
		return err
	}

	if result.ModifiedCount == 0 {
		log.Println("No organization found with ID:", organizationID)
		return mongo.ErrNoDocuments
	}

	return nil
}

// Function that returns an organization object from the database using its ID
func GetOrganizationByID(organizationID string) (organization *models.Organization, err error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	err = organizationCollection.FindOne(ctx, bson.M{"organization_id": organizationID}).Decode(&organization)
	return
}
