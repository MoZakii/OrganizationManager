package controllers

import (
	"MoZaki-Organization-Manager/pkg/database/mongodb/models"
	"MoZaki-Organization-Manager/pkg/database/mongodb/repository"
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var organizationCollection *mongo.Collection = repository.OpenCollection(repository.Client, "organization")
var validate = validator.New()

func MatchAccessLevelOfUser(c *gin.Context, organizationID string) (err error) {
	//Find organization in database  DONE
	var organization *models.Organization

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	err = organizationCollection.FindOne(ctx, bson.M{"organization_id": organizationID}).Decode(&organization)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	//Get user token and find user email from it    NOT DONE

	var user models.User // ADD USER HERE

	var member models.Organization_Member
	// Find user email in members list DONE

	var found bool = false
	for _, item := range organization.Organization_Members {
		if item.Email == user.Email {
			found = true
			member = item
			break
		}
	}

	if !found {
		err = errors.New("unauthorized access to this recource")
		return err
	}

	// Check if user has "ADMIN" access DONE

	if *member.Access_Level == "USER" {
		err = errors.New("unauthorized access to this recource")
		return err
	}

	return nil
}

func ContainsUser(c *gin.Context, organizationID string) (organization *models.Organization, err error) {

	//Find organization in database  DONE
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	err = organizationCollection.FindOne(ctx, bson.M{"organization_id": organizationID}).Decode(&organization)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	//Get user token and find user email from it   NOT DONE

	var user models.User // ADD USER HERE
	// Find user email in members list

	var found bool = false
	for _, item := range organization.Organization_Members {
		if item.Email == user.Email {
			found = true
			break
		}
	}

	if !found {
		err = errors.New("unauthorized access to this recource")
		return nil, err
	}
	return organization, nil
}
