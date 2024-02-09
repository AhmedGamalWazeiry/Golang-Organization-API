package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"org.com/org/pkg/controllers"
	"org.com/org/pkg/database/mongodb/models"
	"org.com/org/pkg/utils"
)

func CreateOrganizationHandler(c *gin.Context) {
	var organization models.OrganizationView

	// Bind JSON data to the 'user' variable
	if err := c.ShouldBindJSON(&organization); err != nil {
		errorMessage := utils.ExtractErrorMessage(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": errorMessage})
		return
	}
    
	userID, _ := c.Get("user_id")
	userIDString, _ := userID.(string)

	organizationID,err := controllers.InsertOrganizationController(organization,userIDString)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"organization_id":organizationID })
}

func GetOrganizationByIDHandler(c *gin.Context) {
	// Get the organization ID from the request
	organizationID := c.Param("organization_id")

	// Get the user ID from the request
	userEmail, _ := c.Get("user_email")
	userEmailString, _ := userEmail.(string)

	// Pass this data to the controller
	organization, err := controllers.GetOrganizationByIDController(organizationID, userEmailString)
	if err != nil {
		// Check if the error message contains "organization's information"
		if strings.Contains(err.Error(), "organization's information") {
			c.JSON(http.StatusForbidden,  gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the organization data
	c.JSON(http.StatusOK, organization)
}

func GetAllUserOrganizationsHandler(c *gin.Context) {
	// Get the user ID from the request
	userEmail, _ := c.Get("user_email")
	userEmailString, _ := userEmail.(string)

	// Pass this data to the controller
	organizations, err := controllers.GetAllUserOrganizationsController(userEmailString)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the organizations data
	c.JSON(http.StatusOK, organizations)
}


func UpdateOrganizationHandler(c *gin.Context) {
	// Get the organization ID from the request
	organizationID := c.Param("organization_id")

	// Get the user ID from the request
	userEmail, _ := c.Get("user_email")
	userEmailString, _ := userEmail.(string)

	// Bind JSON data to the 'organization' variable
	var organization models.OrganizationView
	if err := c.ShouldBindJSON(&organization); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Pass this data to the controller
	err := controllers.UpdateOrganizationController(organizationID, userEmailString, organization)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"organization_id": organizationID,"name": organization.Name,"description": organization.Description})
}
func DeleteOrganizationHandler(c *gin.Context) {
	// Get the organization ID from the request
	organizationID := c.Param("organization_id")

	// Get the user ID from the request
	userEmail, _ := c.Get("user_email")
	userEmailString, _ := userEmail.(string)
	// Pass this data to the controller
	err := controllers.DeleteOrganizationController(organizationID, userEmailString)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Organization deleted successfully"})
}

func InviteUserHandler(c *gin.Context) {
	// Get the organization ID from the request
	organizationID := c.Param("organization_id")

	// Get the user ID from the request
	userID, _ := c.Get("user_id")
	userIDString, _ := userID.(string)

	// Bind JSON data to the 'invite' variable
	var invite struct {
		UserEmail string `json:"user_email"`
	}
	if err := c.ShouldBindJSON(&invite); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Pass this data to the controller
	err := controllers.InviteUserController(organizationID, userIDString, invite.UserEmail)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User invited successfully"})
}