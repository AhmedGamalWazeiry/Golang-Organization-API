package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"org.com/org/pkg/controllers"
	"org.com/org/pkg/database/mongodb/models"
	"org.com/org/pkg/utils"
)

// Register handles the POST request for creating a new user.
func Register(c *gin.Context) {
	var user models.User

	// Bind JSON data to the 'user' variable
	if err := c.ShouldBindJSON(&user); err != nil {
		errorMessage:= utils.ExtractErrorMessage(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": errorMessage})
		return
	}

	emailValidationMessage := utils.ValidateEmail(user.Email)
	if emailValidationMessage != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": emailValidationMessage})
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
		errorMessage:= utils.ExtractErrorMessage(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": errorMessage})
		return
	}

	emailValidationMessage := utils.IsValidEmail(user.Email)

	if !emailValidationMessage  {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
		return
	}

	// Authenticate user using controller
	authenticated_user, err := controllers.AuthenticateUser(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate access and refresh tokens
	accessToken, refreshToken, err := controllers.GenerateTokenPair(*authenticated_user)
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
func RevokeRefreshToken(c *gin.Context) {
	var RefreshTokenRequest models.RefreshTokenRequest

	// Bind JSON data to the 'refreshTokenRequest' variable
	if err := c.ShouldBindJSON(&RefreshTokenRequest); err != nil {
		errorMessage:= utils.ExtractErrorMessage(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": errorMessage})
		return
	}

	// Extract refresh token from the request
	refreshToken:= RefreshTokenRequest.RefreshToken
	
	err := utils.BlacklistToken(refreshToken)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "Unable to revoke the refresh token; it appears that the refresh token may be invalid."})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Refresh token revoked successfully"})
}

// RefreshToken handles the POST request for refreshing access tokens.
func RefreshToken(c *gin.Context) {
	var RefreshTokenRequest models.RefreshTokenRequest

	// Bind JSON data to the 'refreshTokenRequest' variable
	if err := c.ShouldBindJSON(&RefreshTokenRequest); err != nil {
		errorMessage:= utils.ExtractErrorMessage(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": errorMessage})
		return
	}

	// Extract refresh token from the request
	refreshToken:= RefreshTokenRequest.RefreshToken

	ok, err := utils.IsTokenBlacklisted(refreshToken)

	if ok || err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Refresh Token"})
			return
		}
    
	// Validate and refresh the access token using controller
	accessToken, newRefreshToken, err := controllers.GenerateTokenPairByRefreshToken(refreshToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

    err = utils.BlacklistToken(refreshToken)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "Failed to create Refresh Token"})
		return
	}
	// Respond with the new tokens and success message
	c.JSON(http.StatusOK, gin.H{
		"message":         "Token refreshed successfully",
		"access_token":    accessToken,
		"refresh_token":   newRefreshToken,
	})
}


func Test(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{
		"message":         "TEST auth successfully",
	
	})
	return
}