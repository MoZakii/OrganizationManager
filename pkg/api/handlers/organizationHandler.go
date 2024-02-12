package handlers

import (
	"MoZaki-Organization-Manager/pkg/controllers"
	"MoZaki-Organization-Manager/pkg/database/mongodb/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateOrganization() gin.HandlerFunc {
	return func(c *gin.Context) {
		var organization models.Organization

		if err := c.BindJSON(&organization); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := controllers.CreateOrganization(c, &organization)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Organization Item was not created : " + err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"organization_id": organization.Organization_ID})

	}
}

func ReadOrganization() gin.HandlerFunc {
	return func(c *gin.Context) {

		organization, err := controllers.GetOrganization(c)
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
		organization, err := controllers.UpdateOrganization(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"organization_id": organization.Organization_ID,
			"name":        organization.Name,
			"description": organization.Description,
		})

	}
}

func DeleteOrganization() gin.HandlerFunc {
	return func(c *gin.Context) {

		err := controllers.DeleteOrganization(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Organization deleted successfully"})

	}
}

func InviteUserToOrganization() gin.HandlerFunc {
	return func(c *gin.Context) {

		err := controllers.AddToOrganization(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Invited successfully"})
	}
}
