package handlers

import (
	"MoZaki-Organization-Manager/pkg/controllers"
	"MoZaki-Organization-Manager/pkg/database/mongodb/models"
	"MoZaki-Organization-Manager/pkg/database/mongodb/repository"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var organizationCollection *mongo.Collection = repository.OpenCollection(repository.Client, "organization")

func CreateOrganization() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var organization models.Organization

		defer cancel()
		if err := c.BindJSON(&organization); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		temp, exists := c.Get("user_email")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "email doesnt exist"})
			return
		}
		organization.Author_Email = temp.(string)
		validationErr := validate.Struct(organization)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}
		organization.Organization_Members = make([]models.Organization_Member, 0)
		organization.ID = primitive.NewObjectID()
		organization.Organization_ID = organization.ID.Hex()

		_, insertErr := organizationCollection.InsertOne(ctx, organization)

		if insertErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Organization Item was not created"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"organization_id": organization.Organization_ID})

	}
}

func ReadOrganization() gin.HandlerFunc {
	return func(c *gin.Context) {
		organizationID := c.Param("organization_id")
		organization, err := controllers.GetOrganization(c, organizationID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"organization_id": organization.Organization_ID,
			"name":                 organization.Name,
			"description":          organization.Description,
			"organization_members": organization.Organization_Members})
	}
}

func ReadAllOrganizations() gin.HandlerFunc {
	return func(c *gin.Context) {
		organizations, err := controllers.GetAllOrganizations(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, organizations)
	}
}

func UpdateOrganization() gin.HandlerFunc {
	return func(c *gin.Context) {
		organizationID := c.Param("organization_id")
		if _, err := controllers.MatchAccessLevelOfUser(c, organizationID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		organization, err := controllers.UpdateOrganization(c, organizationID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"organization_id": organizationID,
			"name":        organization.Name,
			"description": organization.Description,
		})

	}
}

func DeleteOrganization() gin.HandlerFunc {
	return func(c *gin.Context) {
		organizationID := c.Param("organization_id")
		if _, err := controllers.MatchAccessLevelOfUser(c, organizationID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		err := controllers.DeleteOrganization(c, organizationID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Organization deleted successfully"})

	}
}

func InviteUserToOrganization() gin.HandlerFunc {
	return func(c *gin.Context) {
		organizationID := c.Param("organization_id")
		if _, err := controllers.MatchAccessLevelOfUser(c, organizationID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var member models.Organization_Member
		if err := c.BindJSON(&member); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		organization, err := controllers.ContainsUser(c, organizationID, *member.Email)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user, err := GetUserByEmail(c, member.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Email doesn't exist"})
			return
		}
		member.Name = user.Name
		organization.Organization_Members = append(organization.Organization_Members, member)

		err = controllers.AddToOrganization(c, organizationID, &organization.Organization_Members)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Invited successfully"})
	}
}
