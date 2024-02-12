package handlers

import (
	"MoZaki-Organization-Manager/pkg/controllers"
	"MoZaki-Organization-Manager/pkg/database/mongodb/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Function that handles Create Organization Route
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

// Function that handles read organization route
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

// Function that handles read all organizations route
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

// Function that handles update organization route
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

// Function that handles delete organization route
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

// Function that handles inviting users to organization route
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
