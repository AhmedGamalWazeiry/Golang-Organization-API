package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"org.com/org/pkg/controllers"
	"org.com/org/pkg/database/mongodb/models"
)

// GetUserByID handles the GET request for retrieving a user by ID.
func GetUserByID(c *gin.Context) {
	userID := c.Param("id") // Get the user ID from the URL parameter
	fmt.Println("User IDdsdsadsaaaaaaaaaaaaaaaaaa:", userID)
	// For testing purposes, create a dummy user with the provided ID
	dummyUser := models.User{
		
		Name:     "Dummy User",
		Email:    "dummy@example.com",
		Password: "dummyPassword",
	}

	c.JSON(http.StatusOK, dummyUser)
}

// CreateUser handles the POST request for creating a new user.
func CreateUser(c *gin.Context) {
	fmt.Println("User IDdsdsadsaaaaaaaaaaaaaaaaaa3:")
	var user models.User
	
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println("User IDdsdsadsaaaaaaaaaaaaaaaaaa1:")
	createdID, err := controllers.CreateUser(user)
	fmt.Println("User IDdsdsadsaaaaaaaaaaaaaaaaaa2:")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": createdID})
}

// UpdateUser handles the PUT request for updating a user by ID.
func UpdateUser(c *gin.Context) {
	// Implementation similar to GetUsers
}

// DeleteUser handles the DELETE request for deleting a user by ID.
func DeleteUser(c *gin.Context) {
	// Implementation similar to GetUsers
}