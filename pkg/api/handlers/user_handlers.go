package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"org.com/org/pkg/controllers"
	"org.com/org/pkg/database/mongodb/models"
)

// Register handles the POST request for creating a new user.
func Register(c *gin.Context) {
	var user models.User

	// Bind JSON data to the 'user' variable
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create user using controller
	 err := controllers.CreateUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Congratulations! You've successfully signed up and joined our community."})
}

// Login handles the POST request for user authentication.
func Login(c *gin.Context) {
	var user models.UserLoginRequest

	// Bind JSON data to the 'user' variable
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Authenticate user using controller
	response, err := controllers.AuthenticateUser(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate access and refresh tokens
	accessToken, refreshToken, err := controllers.GenerateTokenPair(*response)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Respond with tokens and success message
	c.JSON(http.StatusOK, gin.H{
		"message":         "Login successful",
		"access_token":    accessToken,
		"refresh_token":   refreshToken,
	})
}

// Logout handles the POST request for logging out a user.
func Logout(c *gin.Context) {
	// Clear cookies on logout
	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}
