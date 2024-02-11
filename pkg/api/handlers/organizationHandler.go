package handlers

import (
	"MoZaki-Organization-Manager/pkg/controllers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateOrganization() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func ReadOrganization() gin.HandlerFunc {
	return func(c *gin.Context) {
		organizationID := c.Param("organization_id")
		organization, err := controllers.ContainsUser(c, organizationID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		c.JSON(http.StatusOK, gin.H{"organization": organization})
	}
}

func ReadAllOrganizations() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func UpdateOrganization() gin.HandlerFunc {
	return func(c *gin.Context) {
		organizationID := c.Param("organization_id")

		if err := controllers.MatchAccessLevelOfUser(c, organizationID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
	}
}

func DeleteOrganization() gin.HandlerFunc {
	return func(c *gin.Context) {
		organizationID := c.Param("organization_id")
		if err := controllers.MatchAccessLevelOfUser(c, organizationID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

	}
}

func InviteUserToOrganization() gin.HandlerFunc {
	return func(c *gin.Context) {
		organizationID := c.Param("organization_id")
		if err := controllers.MatchAccessLevelOfUser(c, organizationID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
	}
}
