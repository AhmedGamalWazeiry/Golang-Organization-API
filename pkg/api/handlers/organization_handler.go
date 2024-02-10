package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"org.com/org/pkg/controllers"
	"org.com/org/pkg/database/mongodb/models"
	"org.com/org/pkg/utils"
)

// CreateOrganizationHandler handles the POST request for creating a new organization.
func CreateOrganizationHandler(c *gin.Context) {
	var organization models.OrganizationView

	statusCode, errorMessage := utils.BindAndValidate(c, &organization)
	if statusCode != http.StatusOK {
		c.JSON(statusCode, gin.H{"error": errorMessage})
		return
	}
    
	userID, _ := c.Get("user_id")
	userIDString, _ := userID.(string)

	statusCode, organizationID, err := controllers.InsertOrganizationController(organization, userIDString)

	if err != nil {
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}
	c.JSON(statusCode, gin.H{"organization_id": organizationID })
}

// GetOrganizationByIDHandler handles the GET request for retrieving an organization by its ID.
func GetOrganizationByIDHandler(c *gin.Context) {
	organizationID := c.Param("organization_id")
	userEmail, _ := c.Get("user_email")
	userEmailString, _ := userEmail.(string)

	statusCode, organization, err := controllers.GetOrganizationByIDController(organizationID, userEmailString)
	if err != nil {
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(statusCode, organization)
}

// GetAllUserOrganizationsHandler handles the GET request for retrieving all organizations a user is part of.
func GetAllUserOrganizationsHandler(c *gin.Context) {
	userEmail, _ := c.Get("user_email")
	userEmailString, _ := userEmail.(string)

	statusCode, organizations, err := controllers.GetAllUserOrganizationsController(userEmailString)
	if err != nil {
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(statusCode, organizations)
}

// UpdateOrganizationHandler handles the PUT request for updating an organization's information.
func UpdateOrganizationHandler(c *gin.Context) {
	organizationID := c.Param("organization_id")
	userEmail, _ := c.Get("user_email")
	userEmailString, _ := userEmail.(string)

	var organization models.OrganizationView
	statusCode, errorMessage := utils.BindAndValidate(c, &organization)
	if statusCode != http.StatusOK {
		c.JSON(statusCode, gin.H{"error": errorMessage})
		return
	}

	statusCode, err := controllers.UpdateOrganizationController(organizationID, userEmailString, organization)
	if err != nil {
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(statusCode, gin.H{"organization_id": organizationID,"name": organization.Name,"description": organization.Description})
}

// DeleteOrganizationHandler handles the DELETE request for deleting an organization.
func DeleteOrganizationHandler(c *gin.Context) {
	organizationID := c.Param("organization_id")
	userEmail, _ := c.Get("user_email")
	userEmailString, _ := userEmail.(string)

	statusCode, err := controllers.DeleteOrganizationController(organizationID, userEmailString)
	if err != nil {
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(statusCode, gin.H{"status": "Organization deleted successfully"})
}

// InviteUserHandler handles the POST request for inviting a user to an organization.
func InviteUserHandler(c *gin.Context) {
	var invite models.Invite
	organizationID := c.Param("organization_id")
	userID, _ := c.Get("user_id")
	userIDString, _ := userID.(string)

	statusCode, errorMessage := utils.BindAndValidate(c, &invite)
	if statusCode != http.StatusOK {
		c.JSON(statusCode, gin.H{"error": errorMessage})
		return
	}

	statusCode, err := controllers.InviteUserController(organizationID, userIDString, invite.Email)
	if err != nil {
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(statusCode, gin.H{"message": "User invited successfully"})
}
