package controllers

import (
	"MoZaki-Organization-Manager/pkg/database/mongodb/models"
	"MoZaki-Organization-Manager/pkg/database/mongodb/repository"
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var Validate = validator.New()

// Function that checks if the User has owner level access on organization
func MatchAccessLevelOfUser(c *gin.Context, organizationID string) (organization *models.Organization, err error) {
	//Find organization in database

	organization, err = repository.GetOrganizationByID(organizationID)
	if err != nil {
		return
	}

	//Get user email from context and compare it to organization author.

	userEmail, exists := c.Get("user_email")
	if !exists {
		err = errors.New("email not found")
		return
	}

	// Compare user to organization author

	if userEmail != organization.Author_Email {
		err = errors.New("unauthorized access to this recource")
		return
	}

	return organization, nil
}

// Function that checks whether user is present in the organization members
func ContainsUser(c *gin.Context, organizationID string, userEmail string) (organization *models.Organization, found bool, err error) {

	found = false
	//Find organization in database

	organization, err = repository.GetOrganizationByID(organizationID)
	if err != nil {
		return
	}

	//Compare user email to organization members.

	for _, item := range organization.Organization_Members {
		if *item.Email == userEmail {
			found = true
			break
		}
	}
	return organization, found, err
}

// Function that extracts organization data from context and creates an organization with given data
func CreateOrganization(c *gin.Context, organization *models.Organization) (err error) {

	temp, exists := c.Get("user_email")
	if !exists {
		err = errors.New("email doesnt exist")
		return
	}
	organization.Author_Email = temp.(string)
	err = Validate.Struct(organization)
	if err != nil {
		return
	}
	organization.Organization_Members = make([]models.Organization_Member, 0)
	organization.ID = primitive.NewObjectID()
	organization.Organization_ID = organization.ID.Hex()
	err = repository.CreateOrganization(*organization)
	return
}

// Function that returns organization data given the id
func GetOrganization(c *gin.Context) (organization *models.Organization, err error) {
	organizationID := c.Param("organization_id")
	organization, err = repository.GetOrganization(organizationID)

	return
}

// Function that returns all organizations' data
func GetAllOrganizations(c *gin.Context) (organizations []repository.NeededInfo, err error) {

	organizations, err = repository.GetAllOrganizations()

	return

}

// Function that updates organization data after extracting them from context
func UpdateOrganization(c *gin.Context) (organization *models.Organization, err error) {

	organizationID := c.Param("organization_id")
	_, err = MatchAccessLevelOfUser(c, organizationID)
	if err != nil {
		return
	}

	err = c.BindJSON(&organization)
	if err != nil {
		return
	}
	organization.Organization_ID = organizationID
	err = repository.UpdateOrganization(*organization)

	return organization, err
}

// Function that deletes an organization given an id
func DeleteOrganization(c *gin.Context) (err error) {
	organizationID := c.Param("organization_id")
	_, err = MatchAccessLevelOfUser(c, organizationID)
	if err != nil {
		return
	}
	err = repository.DeleteOrganization(organizationID)
	return

}

// Function that adds a user to the members of a given organization
func AddToOrganization(c *gin.Context) (err error) {

	organizationID := c.Param("organization_id")
	_, err = MatchAccessLevelOfUser(c, organizationID) // Check if user is allowed to invite
	if err != nil {
		return
	}

	var member models.Organization_Member
	err = c.BindJSON(&member)
	if err != nil {
		return
	}
	_, found, err := ContainsUser(c, organizationID, *member.Email) // Check if organization contains the invited user
	if err != nil {
		return
	}
	if found {
		err = errors.New("user is already an organization member")
		return
	}
	user, err := repository.GetUserByEmail(member.Email) //Fetch invited user data
	if err != nil {
		return
	}
	member.Name = user.Name

	err = repository.AddToOrganization(organizationID, member) // Add user to the organization

	return err
}
