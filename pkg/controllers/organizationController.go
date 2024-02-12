package controllers

import (
	"MoZaki-Organization-Manager/pkg/database/mongodb/models"
	"MoZaki-Organization-Manager/pkg/database/mongodb/repository"
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var organizationCollection *mongo.Collection = repository.OpenCollection(repository.Client, "organization")

func MatchAccessLevelOfUser(c *gin.Context, organizationID string) (organization *models.Organization, err error) {
	//Find organization in database  DONE

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	err = organizationCollection.FindOne(ctx, bson.M{"organization_id": organizationID}).Decode(&organization)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	//Get user email from context and compare it to organization author.

	userEmail, exists := c.Get("user_email")
	if !exists {
		err = errors.New("email not found")
		return nil, err
	}

	// Compare user to organization author

	if userEmail != organization.Author_Email {
		err = errors.New("unauthorized access to this recource")
		return nil, err
	}

	return organization, nil
}

func ContainsUser(c *gin.Context, organizationID string, userEmail string) (organization *models.Organization, err error) {

	//Find organization in database  DONE
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	err = organizationCollection.FindOne(ctx, bson.M{"organization_id": organizationID}).Decode(&organization)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	//Compare user email to organization members.

	var found bool = false
	for _, item := range organization.Organization_Members {
		if *item.Email == userEmail {
			found = true
			break
		}
	}

	if found {
		err = errors.New("user is already an organization member")
		return nil, err
	}
	return organization, nil
}

func GetOrganization(c *gin.Context, organizationID string) (organization *models.Organization, err error) {
	//Find organization in database  DONE
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	err = organizationCollection.FindOne(ctx, bson.M{"organization_id": organizationID}).Decode(&organization)
	if err != nil {
		return
	}
	return

}

type NeededInfo struct {
	Organization_ID      string
	Name                 string
	Description          string
	Organization_members []models.Organization_Member
}

func GetAllOrganizations(c *gin.Context) (organizations []NeededInfo, err error) {
	//Find organization in database  DONE

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	cursor, err := organizationCollection.Find(ctx, bson.M{})
	if err != nil {
		return
	}

	if err = cursor.All(ctx, &organizations); err != nil {
		return
	}

	// Prints the results of the find operation as structs
	for _, result := range organizations {
		cursor.Decode(&result)
	}

	return

}

func UpdateOrganization(c *gin.Context, organizationID string) (organization *models.Organization, err error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	if err := c.BindJSON(&organization); err != nil {
		return nil, err
	}

	var updateObj primitive.D

	updateObj = append(updateObj, bson.E{Key: "name", Value: organization.Name})
	updateObj = append(updateObj, bson.E{Key: "description", Value: organization.Description})

	filter := bson.M{"organization_id": organizationID}

	_, err = organizationCollection.UpdateOne(
		ctx,
		filter,
		bson.D{
			{Key: "$set", Value: updateObj},
		},
	)

	if err != nil {
		log.Panic(err)
		return nil, err
	}
	return organization, nil
}

func DeleteOrganization(c *gin.Context, organizationID string) (err error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	filter := bson.M{"organization_id": organizationID}
	result, err := organizationCollection.DeleteOne(ctx, filter)
	if err != nil || result.DeletedCount == 0 {
		return errors.New("organization id doesnt match")
	}
	return nil

}

func AddToOrganization(c *gin.Context, organizationID string, members *[]models.Organization_Member) (err error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var updateObj primitive.D

	updateObj = append(updateObj, bson.E{Key: "organization_members", Value: members})

	filter := bson.M{"organization_id": organizationID}

	_, err = organizationCollection.UpdateOne(
		ctx,
		filter,
		bson.D{
			{Key: "$set", Value: updateObj},
		},
	)

	if err != nil {
		log.Panic(err)
		return err
	}
	return nil
}
